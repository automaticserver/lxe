package lxo

import (
	"errors"
	"fmt"

	lxd "github.com/lxc/lxd/client"
	"github.com/lxc/lxd/shared/api"
)

// CopyImage will copy an image from the specified server waits till operation is done
func (l *LXO) CopyImage(source lxd.ImageServer, image api.Image, args *lxd.ImageCopyArgs) error {
	op, err := l.server.CopyImage(source, image, args)
	if err != nil {
		return err
	}

	return op.Wait()
}

// DeleteImage will delete an image and waits till operation is done
func (l *LXO) DeleteImage(hash string) error {
	op, err := l.server.DeleteImage(hash)
	if err != nil {
		return err
	}

	return op.Wait()
}

var (
	ErrParse = errors.New("parse error")
)

// CreateImage will create an image and waits till operation is done. Returns resulting fingerprint
func (l *LXO) CreateImage(image api.ImagesPost, args *lxd.ImageCreateArgs) (string, error) {
	op, err := l.server.CreateImage(image, args)
	if err != nil {
		return "", err
	}

	err = op.Wait()
	if err != nil {
		return "", err
	}

	opAPI := op.Get()

	fingerprint, ok := opAPI.Metadata["fingerprint"].(string)
	if !ok {
		return "", fmt.Errorf("%w: %#v", ErrParse, opAPI.Metadata["fingerprint"])
	}

	return fingerprint, nil
}
