package lxo

import (
	"strings"

	lxd "github.com/canonical/lxd/client"
	"github.com/canonical/lxd/shared/api"
)

// StopInstance will try to stop the instance and waits till operation is done
func (l *LXO) StopInstance(id string, timeout, retries int) error {
	var (
		err  error
		etag string
	)

	for i := 0; i <= retries; i++ {
		lxdReq := api.InstanceStatePut{
			Action:  "stop", // note: can't get constant within lxd/lxd/*, first licencing and we don't want to build dqlite
			Timeout: timeout,
			Force:   i == retries,
		}

		var op lxd.Operation

		op, err = l.server.UpdateInstanceState(id, lxdReq, etag)
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

// StartInstance will start the instance and waits till operation is done
func (l *LXO) StartInstance(id string) error {
	ETag := ""
	lxdReq := api.InstanceStatePut{
		Action:  "start", // note: can't get constant within lxd/lxd/*, first licencing and we don't want to build dqlite
		Timeout: -1,
	}

	op, err := l.server.UpdateInstanceState(id, lxdReq, ETag)
	if err != nil {
		return err
	}

	return op.Wait()
}

// CreateInstance will create the instance and waits till operation is done
func (l *LXO) CreateInstance(instance api.InstancesPost) error {
	op, err := l.server.CreateInstance(instance)
	if err != nil {
		return err
	}

	return op.Wait()
}

// UpdateInstance will create the instance and waits till operation is done
func (l *LXO) UpdateInstance(id string, instance api.InstancePut, etag string) error {
	op, err := l.server.UpdateInstance(id, instance, etag)
	if err != nil {
		return err
	}

	return op.Wait()
}

// DeleteInstance will delete the instance and waits till operation is done
func (l *LXO) DeleteInstance(id string) error {
	op, err := l.server.DeleteInstance(id)
	if err != nil {
		return err
	}

	return op.Wait()
}
