package cri

import (
	"fmt"
	"strings"
)

// Convert docker image names to lxd image names.
// images/ubuntu/14.04:tagname -> images:ubuntu/14.04%tagname
func convertDockerImageNameToLXD(name string) string {
	// ignore docker hub prefix
	name = strings.TrimPrefix(name, "docker.io/library/")

	// ignore latest tag
	name = strings.TrimSuffix(name, ":latest")

	// mask docker tag separator with a percent sign
	name = strings.Replace(name, ":", "%", 1)

	// everything else is remote/path where path is a full pathed alias or a fingerprint
	before, after, found := strings.Cut(name, "/")

	// except it contains no path separator. This is most probably a fingerprint or possibly an alias without path separator
	if !found {
		return name
	}

	return fmt.Sprintf("%s:%s", before, after)
}

// Convert lxe alias names to docker names.
// images/ubuntu/14.04%tagname -> images/ubuntu/14.04:tagname
func convertLXEAliasNameToDocker(name string) string {
	// unmask docker tag separator back to colon
	name = strings.Replace(name, "%", ":", 1)

	// add latest if no tag is present
	if !strings.Contains(name, ":") {
		name = fmt.Sprintf("%s:latest", name)
	}

	return name
}
