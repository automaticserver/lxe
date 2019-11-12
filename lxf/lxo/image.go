package lxo

import (
	lxd "github.com/lxc/lxd/client"
	"github.com/lxc/lxd/shared/api"
)

// CopyImage copies an image from the specified server and wait till operation is done or
// return an error
func (l *LXO) CopyImage(source lxd.ImageServer, image api.Image, args *lxd.ImageCopyArgs) error {
	op, err := l.server.CopyImage(source, image, args)
	if err != nil {
		return err
	}

	return op.Wait()
}

// DeleteImage deletes an image and wait till operation is done or
// return an error
func (l *LXO) DeleteImage(hash string) error {
	op, err := l.server.DeleteImage(hash)
	if err != nil {
		return err
	}

	return op.Wait()
}
