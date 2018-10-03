package lxf

import (
	"fmt"

	"github.com/lxc/lxd/shared/api"
	"github.com/lxc/lxe/lxf/device"
	"github.com/lxc/lxe/lxf/lxo"
)

func (l *LXF) createSleeper(s *Sandbox) error {
	devs := map[string]map[string]string{}

	err := device.AddDisksToMap(devs, device.Disk{
		Path: "/",
		Pool: "default",
	})

	if err != nil {
		return err
	}

	// create container
	op, err := l.server.CreateContainer(api.ContainersPost{
		Name: sleeperName(s.ID),
		ContainerPut: api.ContainerPut{
			Profiles: []string{},
			Devices:  devs,
		},
		Source: api.ContainerSource{
			Fingerprint: l.sleeperHash,
			Type:        "image",
		},
	})

	if err != nil {
		return err
	}

	err = op.Wait()
	if err != nil {
		return err
	}

	return lxo.StartContainer(l.server, sleeperName(s.ID))
}

func (l *LXF) stopSleeper(name string) error {
	return lxo.StopContainer(l.server, sleeperName(name))
}

func (l *LXF) deleteSleeper(name string) error {
	// just to make sure it really is stopped
	err := l.stopSleeper(name)
	if err != nil {
		return err
	}

	op, err := l.server.DeleteContainer(sleeperName(name))
	if err != nil {
		if err.Error() == "not found" { // it's already gone, that's ok with us
			return nil
		}
		return fmt.Errorf("delete sleeper container '%v', %v", sleeperName(name), err)
	}
	return op.Wait()
}

func sleeperName(name string) string {
	return "sleeper-" + name
}
