package lxf_test

import (
	"os"
	"testing"

	"github.com/automaticserver/lxe/lxf"
)

func TestConnection(t *testing.T) {
	_, err := lxf.New("", os.Getenv("HOME")+"/.config/lxc/config.yml")
	if err != nil {
		t.Errorf("failed to set up connection %v", err)
	}
}

func TestConnectionWithInvalidSocket(t *testing.T) {
	_, err := lxf.New("/var/lib/lxd/unix.invalidsocket",
		os.Getenv("HOME")+"/.config/lxc/config.yml")
	if err == nil {
		t.Errorf("invalid socket should return an error")
	}
}
