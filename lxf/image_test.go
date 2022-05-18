package lxf

import (
	"testing"
)

func Test_convertDockerImageNameToLXC(t *testing.T) {
	tests := []struct {
		inputName string
		want      string
		wantErr   bool
	}{
		{"busybox", "busybox", false},
		{"busybox:other", "busybox", false},
		{"hub.example.io/busybox:other", "hub.example.io:busybox", false},
		{"hub.example.io/someuser/images/busybox:other", "hub.example.io:someuser/images/busybox", false},
		{"images/ubuntu/14.04", "images:ubuntu/14.04", false},
		{"images/ubuntu/14.04:latest", "images:ubuntu/14.04", false},
		{"missingremote/example/ubuntu/14.04", "missingremote:example/ubuntu/14.04", false},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.inputName, func(t *testing.T) {
			got, err := convertDockerImageNameToLXC(tt.inputName)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertDockerImageNameToLXC() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("convertDockerImageNameToLXC() = %v, want %v", got, tt.want)
			}
		})
	}
}
