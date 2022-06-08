package lxo

import (
	"strings"

	lxd "github.com/lxc/lxd/client"
	"github.com/lxc/lxd/shared"
	"github.com/lxc/lxd/shared/api"
)

// StopContainer will try to stop the container and waits till operation is done
func (l *LXO) StopContainer(id string, timeout, retries int) error {
	var (
		err  error
		etag string
	)

	for i := 0; i <= retries; i++ {
		lxdReq := api.ContainerStatePut{
			Action:  string(shared.Stop),
			Timeout: timeout,
			Force:   i == retries,
		}

		var op lxd.Operation

		op, err = l.server.UpdateContainerState(id, lxdReq, etag)
		if err != nil {
			return err
		}

		err = op.Wait()
		if err != nil {
			if strings.Contains(err.Error(), "is already stopped") {
				return nil
			}
		} else {
			return nil
		}
	}

	return err
}

// StartContainer will start the container and waits till operation is done
func (l *LXO) StartContainer(id string) error {
	ETag := ""
	lxdReq := api.ContainerStatePut{
		Action:  string(shared.Start),
		Timeout: -1,
	}

	op, err := l.server.UpdateContainerState(id, lxdReq, ETag)
	if err != nil {
		return err
	}

	return op.Wait()
}

// CreateContainer will create the container and waits till operation is done
func (l *LXO) CreateContainer(container api.ContainersPost) error {
	op, err := l.server.CreateContainer(container)
	if err != nil {
		return err
	}

	return op.Wait()
}

// UpdateContainer will create the container and waits till operation is done
func (l *LXO) UpdateContainer(id string, container api.ContainerPut, etag string) error {
	op, err := l.server.UpdateContainer(id, container, etag)
	if err != nil {
		return err
	}

	return op.Wait()
}

// DeleteContainer will delete the container and waits till operation is done
func (l *LXO) DeleteContainer(id string) error {
	op, err := l.server.DeleteContainer(id)
	if err != nil {
		return err
	}

	return op.Wait()
}
