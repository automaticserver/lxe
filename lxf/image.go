package lxf

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	lxd "github.com/lxc/lxd/client"
	lxdApi "github.com/lxc/lxd/shared/api"
	"github.com/lxc/lxe/lxf/lxo"
)

// Image is here to translate the relevant data from lxd image to cri image
type Image struct {
	Hash    string
	Aliases []string
	Size    int64
}

// PullImage copies the given image from the remote server
func (l *LXF) PullImage(name string) (string, error) {
	imageID, err := l.parseImage(name)
	if err != nil {
		return "", err
	}

	// we will cretae an image server for the remote.
	// we will also create one when it's the default remote, because the default does not always
	// need to be the local.
	imgServer, err := l.config.GetImageServer(imageID.Remote)
	if err != nil {
		return "", err
	}

	imageRef := dereferenceAlias(imgServer, imageID.Alias)
	image, _, err := imgServer.GetImage(imageRef)
	if err != nil {
		return "", err
	}

	args := lxd.ImageCopyArgs{
		CopyAliases: false, // We shouldn't rely on default aliases, as aliases are unique per remote
		AutoUpdate:  true,  // Maybe bug: currently NOT a technical requirement to know where the source is
	}

	err = lxo.CopyImage(l.server, imgServer, *image, &args)
	if err != nil {
		return "", fmt.Errorf("unable to pull requested image (%v) from server %v, %v",
			image, imageID.Remote, err)
	}

	return image.Fingerprint, l.ensureImageAlias(imageID.Tag(), image.Fingerprint)
}

// RemoveImage will remove the given alias
func (l *LXF) RemoveImage(name string) error {
	imageID, err := l.parseImage(name)
	if err != nil {
		return err
	}

	hash, found, err := imageID.Hash(l)
	if !found {
		return nil
	}
	if err != nil {
		return err
	}

	err = lxo.DeleteImage(l.server, hash)
	if err != nil {
		if IsErrorNotFound(err) {
			return nil
		}
		return err
	}

	return nil
}

// Create the specified image alis, update if already exist
// from github.com/lxc/lxd/lxc/image.go:172 + changes
func (l *LXF) ensureImageAlias(alias string, fingerprint string) error {
	current, err := l.server.GetImageAliases()
	if err != nil {
		return err
	}

	alreadyCorrect := false

	// Delete existing aliases that match provided one but has not already correct fingerprint
	for _, ca := range current {
		if ca.Name == alias {
			if ca.Target == fingerprint {
				alreadyCorrect = true
				break
			}
			err = l.server.DeleteImageAlias(ca.Name)
			if err != nil {
				return fmt.Errorf("failed to delete alias %v for update, %v", alias, err)
			}
		}
	}

	// we are done if it points already to correct hash
	if alreadyCorrect {
		return nil
	}

	// Create new alias
	aliasPost := lxdApi.ImageAliasesPost{}
	aliasPost.Name = alias
	aliasPost.Target = fingerprint
	err = l.server.CreateImageAlias(aliasPost)
	if err != nil {
		return fmt.Errorf("failed to create alias %v, %v", alias, err)
	}
	return nil
}

// ListImages will list all local images from the lxd server
func (l *LXF) ListImages(filter string) ([]Image, error) {
	var response = []Image{}
	imglist, err := l.server.GetImages()
	if err != nil {
		return nil, fmt.Errorf("unable to list images, %v", err)
	}

	for _, imgInfo := range imglist {
		if filter != "" && filter != imgInfo.Fingerprint {
			continue
		}
		aliases := []string{}
		for _, ali := range imgInfo.Aliases {
			aliases = append(aliases, ali.Name+":latest")
		}
		response = append(response, Image{
			Hash:    imgInfo.Fingerprint,
			Aliases: aliases,
			Size:    imgInfo.Size,
		})
	}

	return response, nil
}

// GetImage will fetch information about the image identified by name
// It will only work on local images.
func (l *LXF) GetImage(name string) (*Image, error) {
	imageID, err := l.parseImage(name)
	if err != nil {
		if strings.HasSuffix(err.Error(), "doesn't exist") {
			return nil, fmt.Errorf(ErrorNotFound)
		}
		return nil, err
	}

	hash, found, err := imageID.Hash(l)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve image %v, %v", name, err)
	}
	if !found {
		return nil, fmt.Errorf(ErrorNotFound)
	}

	img, _, err := l.server.GetImage(hash)
	if err != nil {
		return nil, fmt.Errorf("unable to get image %v, %v", name, err)
	}
	aliases := []string{}
	for _, ali := range img.Aliases {
		aliases = append(aliases, ali.Name+":latest")
	}
	return &Image{
		Hash:    img.Fingerprint,
		Aliases: aliases,
		Size:    img.Size,
	}, nil
}

// FSPoolUsage contains fields to describe the usage of a filesystem / storagepool
type FSPoolUsage struct {
	Timestamp  int64
	FsID       string
	UsedBytes  uint64
	InodesUsed uint64
}

// GetFSPoolUsage returns a list of usage information about the used storage pools
func (l *LXF) GetFSPoolUsage() ([]FSPoolUsage, error) {
	pools, err := l.server.GetStoragePools()
	if err != nil {
		return nil, err
	}

	rval := []FSPoolUsage{}
	for _, pool := range pools {
		pRcs, err := l.server.GetStoragePoolResources(pool.Name)
		if err != nil {
			return nil, err
		}
		rval = append(rval, FSPoolUsage{
			Timestamp:  time.Now().UnixNano(),
			FsID:       pool.Config["source"],
			UsedBytes:  pRcs.Space.Used,
			InodesUsed: pRcs.Inodes.Used,
		})
	}

	return rval, nil
}

// ImageID contains the remote and alias of an image identifier.
type ImageID struct {
	Remote string
	Alias  string
}

// Tag builds from remote and alias an alias for local
func (i ImageID) Tag() string {
	return i.Remote + "/" + i.Alias
}

// Hash returns the hash from the combined local tag or if it's
// already a hash this one.
// It it's not found second return will be false and error will be zero.
func (i ImageID) Hash(l *LXF) (string, bool, error) {
	exists, _, err := l.server.GetImageAlias(i.Tag())
	if err != nil {
		if IsErrorNotFound(err) {
			// it still might be a hash, check that
			_, _, err = l.server.GetImage(i.Alias)
			if err != nil {
				if IsErrorNotFound(err) {
					return "", false, nil
				}
				return "", false, err
			}
			// it worked so it must be a hash
			return i.Alias, true, nil
		}
		return "", false, err
	}
	// we could resolve the local alias
	return exists.Target, true, nil
}

// parseImage will take an external image and split it up into
// remote and tag
func (l *LXF) parseImage(name string) (ImageID, error) {
	img, err := convertDockerImageNameToLXC(name)
	if err != nil {
		return ImageID{}, err
	}
	remote, tag, err := l.config.ParseRemote(img)
	if err != nil {
		return ImageID{}, err
	}
	return ImageID{Remote: remote, Alias: tag}, nil
}

func convertDockerImageNameToLXC(inputName string) (string, error) {
	// always remove docker tags
	var re = regexp.MustCompile(`(.*)(:.*)`)
	if re.MatchString(inputName) {
		match := re.FindStringSubmatch(inputName)
		inputName = match[1]
	}

	// no path, easy everything is the imageName
	if !strings.Contains(inputName, "/") {
		return inputName, nil
	}

	// with remote/path, use first part as remote only if it contains a dot
	re = regexp.MustCompile(`(.+?)/(.*)`)
	if re.MatchString(inputName) {
		match := re.FindStringSubmatch(inputName)
		return match[1] + ":" + match[2], nil
	}

	return "", fmt.Errorf("could not parse image name %v", inputName)
}

// dereferenceAlias from github.com/lxc/lxd/lxc/image.go:102
// default tag handling removed, that can not happen with our docker-> lxc
// conversion.
func dereferenceAlias(d lxd.ImageServer, inName string) string {
	result, _, err := d.GetImageAlias(inName)
	if result == nil || err != nil {
		return inName
	}
	return result.Target
}
