package lxo

import (
	lxd "github.com/lxc/lxd/client"
	"github.com/lxc/lxd/shared/api"
)

// CopyImage copies an image from the specified server and wait till operation is done or
// return an error
func CopyImage(server lxd.ContainerServer, source lxd.ImageServer, image api.Image, args *lxd.ImageCopyArgs) error {
	op, err := server.CopyImage(source, image, args)
	if err != nil {
		return err
	}

	err = op.Wait()
	return err
}

// DeleteImage deletes an image and wait till operation is done or
// return an error
func DeleteImage(server lxd.ContainerServer, hash string) error {
	op, err := server.DeleteImage(hash)
	if err != nil {
		return err
	}

	err = op.Wait()
	return err
}
