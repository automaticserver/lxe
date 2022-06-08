package lxf

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	lxd "github.com/lxc/lxd/client"
	"github.com/lxc/lxd/shared/api"
)

const lxeAliasPrefix = "lxe/"

// Image is here to translate the relevant data from lxd image to cri image
type Image struct {
	Hash    string
	Aliases []string
	Size    int64
}

// PullImage copies the given image from the remote server. The image is remembered by setting a specific alias.
func (l *client) PullImage(image string) (string, error) {
	orig := image

	if l.critestMode {
		image = replaceCRITestImage(image, l.config.DefaultRemote)
	}

	remote, aliasOrFingerprint, err := l.config.ParseRemote(image)
	if err != nil {
		return "", err
	}

	// get the image server for the remote also when it's the explicitly defined default (local) remote
	imgServer, err := l.config.GetImageServer(remote)
	if err != nil {
		return "", err
	}

	lxdImg, err := getRemoteImageFromAliasOrFingerprint(imgServer, aliasOrFingerprint)
	if err != nil {
		return "", err
	}

	// copy only if it is a foreign remote
	if remote != l.config.DefaultRemote {
		args := lxd.ImageCopyArgs{
			CopyAliases: false,
		}

		err = l.opwait.CopyImage(imgServer, *lxdImg, &args)
		if err != nil {
			return "", err
		}
	}

	setAlias := lxeAlias(image)

	if l.critestMode {
		setAlias = lxeAlias(orig)
	} else {
		err = l.ensureCRIImage(lxdImg.Fingerprint)
		if err != nil {
			return "", err
		}
	}

	err = l.ensureImageAlias(setAlias, lxdImg.Fingerprint)
	if err != nil {
		return "", err
	}

	return lxdImg.Fingerprint, nil
}

// RemoveImage will remove a pulled image
func (l *client) RemoveImage(image string) error {
	if l.critestMode {
		return l.removeCRITestImage(image)
	}

	lxdImg, err := l.getLocalImageFromAliasOrFingerprint(image)
	if err != nil {
		return err
	}

	err = l.opwait.DeleteImage(lxdImg.Fingerprint)
	if err != nil {
		return err
	}

	return nil
}

// Return the alias how lxe is remembering pulled images
func lxeAlias(image string) string {
	return fmt.Sprintf("%s%s", lxeAliasPrefix, strings.Replace(image, ":", "/", 1))
}

// Reverts the lxe specific alias modification
func revLxeAlias(alias string) string {
	return strings.TrimPrefix(alias, lxeAliasPrefix)
}

// Set a cri field property so we know which images are part of the CRI
func (l *client) ensureCRIImage(fingerprint string) error {
	lxdImg, err := l.getLocalImageFromAliasOrFingerprint(fingerprint)
	if err != nil {
		return err
	}

	if lxdImg.Properties == nil {
		lxdImg.Properties = map[string]string{}
	}

	lxdImg.Properties[cfgIsCRI] = strconv.FormatBool(true)

	return l.server.UpdateImage(lxdImg.Fingerprint, lxdImg.Writable(), "")
}

// Create the specified image alias, update if already exist
// from github.com/lxc/lxd/lxc/image.go:172 + changes
func (l *client) ensureImageAlias(alias string, fingerprint string) error {
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
				return fmt.Errorf("failed to delete alias for update: %v, %w", alias, err)
			}
		}
	}

	// we are done if it points already to correct hash
	if alreadyCorrect {
		return nil
	}

	// Create new alias
	aliasPost := api.ImageAliasesPost{}
	aliasPost.Name = alias
	aliasPost.Target = fingerprint

	err = l.server.CreateImageAlias(aliasPost)
	if err != nil {
		return fmt.Errorf("failed to create alias: %v, %w", alias, err)
	}

	return nil
}

// ListImages will list all pulled images
func (l *client) ListImages(filter string) ([]*Image, error) {
	response := []*Image{}

	imglist, err := l.server.GetImages()
	if err != nil {
		return nil, fmt.Errorf("unable to list images: %w", err)
	}

	for _, lxdImg := range imglist {
		lxdImg := lxdImg

		if !l.IsCRI(lxdImg) {
			continue
		}

		if filter != "" && filter != lxdImg.Fingerprint {
			continue
		}

		response = append(response, toImage(&lxdImg))
	}

	return response, nil
}

// GetImage will fetch information about a pulled image
func (l *client) GetImage(image string) (*Image, error) {
	lxdImg, err := l.getLocalImageFromAliasOrFingerprint(image)
	if err != nil {
		return nil, err
	}

	if !l.IsCRI(lxdImg) {
		return nil, ErrNotFound
	}

	return toImage(lxdImg), nil
}

func getImageFingerprint(imgServer lxd.ImageServer, alias string) (string, error) {
	lxdAlias, _, err := imgServer.GetImageAlias(alias)
	if err != nil {
		return "", err
	}

	return lxdAlias.Target, err
}

// similar to l.getLocalImageFromAliasOrFingerprint but remote aliases wont have the lxe specific prefix
func getRemoteImageFromAliasOrFingerprint(imgServer lxd.ImageServer, aliasOrFingerprint string) (*api.Image, error) {
	// try to find out if aliasOrFingerprint is a known alias on the server
	fingerprint, err := getImageFingerprint(imgServer, aliasOrFingerprint)
	if err != nil {
		// if the alias is not found, try to find it by fingerprint
		if IsNotFoundError(err) {
			lxdImg, _, err := imgServer.GetImage(aliasOrFingerprint)

			return lxdImg, err
		}

		return nil, err
	}

	// get the image from alias' target fingerprint
	lxdImg, _, err := imgServer.GetImage(fingerprint)
	if err != nil {
		return nil, err
	}

	return lxdImg, nil
}

func (l *client) getLocalImageFromAliasOrFingerprint(aliasOrFingerprint string) (*api.Image, error) {
	// try to find out if aliasOrFingerprint is a known alias on the server
	fingerprint, err := getImageFingerprint(l.server, lxeAlias(aliasOrFingerprint))
	if err != nil {
		// if the alias is not found, try to find it by fingerprint
		if IsNotFoundError(err) {
			lxdImg, _, err := l.server.GetImage(aliasOrFingerprint)

			return lxdImg, err
		}

		return nil, err
	}

	// get the image from alias' target fingerprint
	lxdImg, _, err := l.server.GetImage(fingerprint)
	if err != nil {
		return nil, err
	}

	return lxdImg, nil
}

// FSPoolUsage contains fields to describe the usage of a filesystem / storagepool
type FSPoolUsage struct {
	Timestamp  int64
	FsID       string
	UsedBytes  uint64
	InodesUsed uint64
}

// GetFSPoolUsage returns a list of usage information about the used storage pools
func (l *client) GetFSPoolUsage() ([]FSPoolUsage, error) {
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

func toImage(lxdImg *api.Image) *Image {
	img := &Image{
		Hash: lxdImg.Fingerprint,
		Size: lxdImg.Size,
	}

	for _, a := range lxdImg.Aliases {
		if strings.HasPrefix(a.Name, lxeAliasPrefix) {
			img.Aliases = append(img.Aliases, revLxeAlias(a.Name))
		}
	}

	return img
}
