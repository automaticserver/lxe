// Package lxo abstracts some of the lxd calls with additional functionality like
// retrying, idempotency and some level of error recovery
package lxo

import (
	"fmt"

	lxd "github.com/lxc/lxd/client"
	"github.com/lxc/lxd/shared/api"
	"github.com/lxc/lxd/shared/logger"
)

// StopContainer will try to stop the container with provided name.
// It will retry for half a minute and return success when it's stopped.
// It will also return success when the container does not exist.
func StopContainer(server lxd.ContainerServer, id string) error {
	tries := 10
	var lastErr error
	for i := 1; i <= tries; i++ {
		lxdReq := api.ContainerStatePut{
			Action:  "stop",
			Timeout: 3,
			Force:   i == tries,
		}
		op, err := server.UpdateContainerState(id, lxdReq, "")
		if err != nil {
			if err.Error() == "not found" { // it's not around, that's ok with us
				return nil
			}
			return fmt.Errorf("failed to stop container %v, %v", id, err)
		}

		err = op.Wait()
		if err != nil && err.Error() == "The container is already stopped" {
			logger.Debugf("container is stopped")
			// done…
			return nil
		}
		lastErr = err
		// we try again with or without err
	}
	return lastErr
}

// StartContainer will start the container and wait till operation is done or
// return an error
func StartContainer(server lxd.ContainerServer, id string) error {
	lxdReq := api.ContainerStatePut{
		Action:  "start",
		Timeout: -1,
	}
	op, err := server.UpdateContainerState(id, lxdReq, "")
	if err != nil {
		return err
	}

	return op.Wait()
}

// CreateContainer will create the container and wait till operation is done or
// return an error
func CreateContainer(server lxd.ContainerServer, container api.ContainersPost) error {
	op, err := server.CreateContainer(container)
	if err != nil {
		return err
	}
	return op.Wait()
}

// UpdateContainer will create the container and wait till operation is done or
// return an error
func UpdateContainer(server lxd.ContainerServer, id string, container api.ContainerPut) error {
	op, err := server.UpdateContainer(id, container, "")
	if err != nil {
		return err
	}
	return op.Wait()
}

// DeleteContainer will delete the container and wait till operation is done or
// return an error
func DeleteContainer(server lxd.ContainerServer, id string) error {
	op, err := server.DeleteContainer(id)
	if err != nil {
		return err
	}
	return op.Wait()
}

// MoveContainer will rename the container and wait till operation is done or
// return an error
func MoveContainer(server lxd.ContainerServer, id string, post api.ContainerPost) error {
	op, err := server.RenameContainer(id, post)
	if err != nil {
		return err
	}

	return op.Wait()
}

// CreateContainerSnapshot creates a snapshot for the container and wait till operation is done or
// return an error
func CreateContainerSnapshot(server lxd.ContainerServer, id string, snapshot api.ContainerSnapshotsPost) error {
	op, err := server.CreateContainerSnapshot(id, snapshot)
	if err != nil {
		return err
	}
	return op.Wait()
}

// CopyContainerSnapshot copies a snapshot into a container an wait till operation is done or
// return an error
func CopyContainerSnapshot(server lxd.ContainerServer, s api.ContainerSnapshot, args *lxd.ContainerSnapshotCopyArgs) error {
	op, err := server.CopyContainerSnapshot(server, s, args)
	if err != nil {
		return err
	}
	return op.Wait()
}
