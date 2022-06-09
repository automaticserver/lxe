package lxf

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/automaticserver/lxe/third_party/ioutils"
	lxd "github.com/lxc/lxd/client"
	"github.com/lxc/lxd/shared/api"
	"gopkg.in/yaml.v3"
)

var (
	critestDefaultImageSource = "images:alpine/edge/cloud"
	critestOtherImageSources  = map[string]string{
		"gcr.io:k8s-staging-cri-tools/test-image-2": "images:alpine/edge/cloud/arm64",
		"gcr.io:k8s-staging-cri-tools/test-image-3": "images:alpine/edge/cloud/armhf",
	}

	critestDefaultAlias   = "critest/default"
	critestWebserverAlias = "critest/webserver"
)

const (
	criTestTimeout = 5 // seconds
)

func isCRITestImage(image api.Image) bool {
	for _, a := range image.Aliases {
		if strings.HasPrefix(a.Name, lxeAliasPrefix) {
			return true
		}
	}

	return false
}

// Replace requested image internally for critest purposes
func replaceCRITestImage(image, defaultRemote string) string {
	// don't manipulate locally referenced images
	if strings.HasPrefix(image, fmt.Sprintf("%s:", defaultRemote)) {
		return image
	}

	ret := critestDefaultImageSource

	replacer, exists := critestOtherImageSources[image]
	if exists {
		ret = replacer
	}

	log.Infof("CRITest: replacing requested image '%s' with '%s'", image, ret)

	return ret
}

// Remove only lxe aliases for critest purposes since the provided base images for critesting are also local and we don't want to remove them
func (l *client) removeCRITestImage(image string) error {
	lxdImg, err := l.getLocalImageFromAliasOrFingerprint(image)
	if err != nil {
		return err
	}

	for _, a := range lxdImg.Aliases {
		if strings.HasPrefix(a.Name, lxeAliasPrefix) {
			log.Infof("CRITest: Deleting alias '%s' for image fingerprint '%s'", a.Name, lxdImg.Fingerprint)

			err = l.server.DeleteImageAlias(a.Name)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (l *client) createCRITestImages() error { // nolint: gocognit, cyclop
	// check whether the default image already exists
	defaultFingerprint, err := getImageFingerprint(l.server, critestDefaultAlias)
	if err != nil && !IsNotFoundError(err) {
		return err
	}

	// create the default image by downloading the predifined source image and alias it accordingly
	if IsNotFoundError(err) { // nolint: nestif
		log.Info("downloading default image")

		remote, alias, err := l.config.ParseRemote(critestDefaultImageSource)
		if err != nil {
			return err
		}

		imgServer, err := l.config.GetImageServer(remote)
		if err != nil {
			return err
		}

		lxdImg, err := getRemoteImageFromAliasOrFingerprint(imgServer, alias)
		if err != nil {
			return err
		}

		err = l.opwait.CopyImage(imgServer, *lxdImg, &lxd.ImageCopyArgs{CopyAliases: false})
		if err != nil {
			return err
		}

		err = l.ensureImageAlias(critestDefaultAlias, lxdImg.Fingerprint)
		if err != nil {
			return err
		}

		log.Infof("default image ready: 'local:%s'", critestDefaultAlias)

		defaultFingerprint = lxdImg.Fingerprint
	}

	// check whether the webserver image already exists
	_, err = getImageFingerprint(l.server, critestWebserverAlias)
	if err != nil && !IsNotFoundError(err) {
		return err
	}

	cName := "critest-nginx"

	// create the webserver image based on the default image, install nginx automatically with cloud-init and alias it accordingly
	if IsNotFoundError(err) { // nolint: nestif
		// in case of a previous error clean up
		err = l.opwait.StopContainer(cName, criTestTimeout, 0)
		if err != nil && !IsNotFoundError(err) {
			return err
		}

		err = l.opwait.DeleteContainer(cName)
		if err != nil && !IsNotFoundError(err) {
			return err
		}

		log.Infof("creating webserver image")

		err = l.opwait.CreateContainer(api.ContainersPost{
			Name: cName,
			Source: api.ContainerSource{
				Fingerprint: defaultFingerprint,
				Type:        "image",
			},
			ContainerPut: api.ContainerPut{
				Config: map[string]string{
					"user.user-data": `#cloud-config
runcmd:
  - apk add nginx
  - rc-update add nginx default
  - rc-service nginx start
`,
				},
			},
		})
		if err != nil {
			return err
		}

		err = l.opwait.StartContainer(cName)
		if err != nil {
			return err
		}

		// wait for cloud-init to complete
		stdin := bytes.NewReader(nil)
		stdinR := ioutil.NopCloser(stdin)
		stdout := bytes.NewBuffer(nil)
		stdoutW := ioutils.WriteCloserWrapper(stdout)
		stderr := bytes.NewBuffer(nil)
		stderrW := ioutils.WriteCloserWrapper(stderr)

		code, err := l.Exec(cName, []string{"cloud-init", "status", "-w"}, stdinR, stdoutW, stderrW, false, false, 10, nil)
		if err != nil {
			return err
		}

		if code != 0 {
			return fmt.Errorf("cloud-init did not complete successfully") // nolint: goerr113
		}

		// clean up cloud-init before we create an image
		_, err = l.Exec(cName, []string{"cloud-init", "clean", "--logs", "--seed"}, stdinR, stdoutW, stderrW, false, false, 10, nil)
		if err != nil {
			return err
		}

		// create the image based on this container
		err = l.opwait.StopContainer(cName, criTestTimeout, 0)
		if err != nil {
			return err
		}

		fingerprint, err := l.opwait.CreateImage(api.ImagesPost{
			Source: &api.ImagesPostSource{
				Type: "instance",
				Name: cName,
			},
		}, nil)
		if err != nil {
			return err
		}

		err = l.ensureImageAlias(critestWebserverAlias, fingerprint)
		if err != nil {
			return err
		}

		log.Infof("webserver image ready: 'local:%s'", critestWebserverAlias)
	}

	imagesFile := fmt.Sprintf("%s/lxe-critest-images-file.yaml", os.TempDir())

	b, err := yaml.Marshal(TestImageList{
		DefaultTestContainerImage: fmt.Sprintf("local/%s", critestDefaultAlias),
		WebServerTestImage:        fmt.Sprintf("local/%s", critestWebserverAlias),
	})
	if err != nil {
		return err
	}

	err = os.WriteFile(imagesFile, b, 0644) // nolint: gomnd
	if err != nil {
		return err
	}

	log.Warnf("CRITest ready: Use --test-images-file=%s in your critest command", imagesFile)

	err = l.opwait.DeleteContainer(cName)
	if err != nil && !IsNotFoundError(err) {
		return err
	}

	return nil
}

// TestImageList aggregates references to the images used in tests.
// Borrowed from github.com/kubernetes-sigs/cri-tools/pkg/framework/test_context.go
type TestImageList struct {
	DefaultTestContainerImage string `yaml:"defaultTestContainerImage"`
	WebServerTestImage        string `yaml:"webServerTestImage"`
}
