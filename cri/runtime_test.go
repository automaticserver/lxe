package cri

import (
	rtApi "k8s.io/cri-api/pkg/apis/runtime/v1"
)

// implements interface
var _ rtApi.RuntimeServiceServer = &RuntimeServer{}
