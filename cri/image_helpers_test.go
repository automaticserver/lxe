package cri

import (
	"testing"
)

func Test_convertDockerImageNameToLXC(t *testing.T) {
	t.Parallel()

	tests := []struct {
		inputName string
		want      string
	}{
		{"busybox", "busybox"},
		{"busybox:other", "busybox%other"},
		{"hub.example.io/busybox:other", "hub.example.io:busybox%other"},
		{"hub.example.io/someuser/images/busybox:other", "hub.example.io:someuser/images/busybox%other"},
		{"images/ubuntu/14.04", "images:ubuntu/14.04"},
		{"images/ubuntu/14.04:latest", "images:ubuntu/14.04"},
		{"missingremote/example/ubuntu/14.04", "missingremote:example/ubuntu/14.04"},
		{"docker.io/library/nginx", "nginx"},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.inputName, func(t *testing.T) {
			t.Parallel()

			got := convertDockerImageNameToLXD(tt.inputName)

			if got != tt.want {
				t.Errorf("convertDockerImageNameToLXC() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convertLXEAliasNameToDocker(t *testing.T) {
	t.Parallel()

	tests := []struct {
		inputName string
		want      string
	}{
		{"busybox", "busybox:latest"},
		{"busybox%other", "busybox:other"},
		{"hub.example.io/busybox%other", "hub.example.io/busybox:other"},
		{"hub.example.io/someuser/images/busybox%other", "hub.example.io/someuser/images/busybox:other"},
		{"images/ubuntu/14.04", "images/ubuntu/14.04:latest"},
		{"images/ubuntu/14.04", "images/ubuntu/14.04:latest"},
		{"missingremote/example/ubuntu/14.04", "missingremote/example/ubuntu/14.04:latest"},
		{"nginx", "nginx:latest"},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.inputName, func(t *testing.T) {
			t.Parallel()

			got := convertLXEAliasNameToDocker(tt.inputName)

			if got != tt.want {
				t.Errorf("convertLXEAliasNameToDocker() = %v, want %v", got, tt.want)
			}
		})
	}
}
