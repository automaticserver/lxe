// Offer compatibility with CRI interface v1alpha2
// Use unsafe Pointer as cri-o did: https://github.com/cri-o/cri-o/commit/96679844e96b9235813580215325ca5f0a0a27a6
// nolint: nlreturn
package cri

import (
	"context"
	"unsafe"

	rtApi "k8s.io/cri-api/pkg/apis/runtime/v1"
	rtApiOld "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

type oldRuntimeServer struct {
	current rtApi.RuntimeServiceServer
}

func NewOldRuntimeServer(current rtApi.RuntimeServiceServer) rtApiOld.RuntimeServiceServer {
	return &oldRuntimeServer{current: current}
}

type oldImageServer struct {
	current rtApi.ImageServiceServer
}

func NewOldImageServer(current rtApi.ImageServiceServer) rtApiOld.ImageServiceServer {
	return &oldImageServer{current: current}
}

func (s *oldRuntimeServer) Attach(ctx context.Context, req *rtApiOld.AttachRequest) (*rtApiOld.AttachResponse, error) {
	resp, err := s.current.Attach(ctx, (*rtApi.AttachRequest)(unsafe.Pointer(req)))
	return (*rtApiOld.AttachResponse)(unsafe.Pointer(resp)), err
}

func (s *oldRuntimeServer) ContainerStats(ctx context.Context, req *rtApiOld.ContainerStatsRequest) (*rtApiOld.ContainerStatsResponse, error) {
	resp, err := s.current.ContainerStats(ctx, (*rtApi.ContainerStatsRequest)(unsafe.Pointer(req)))
	return (*rtApiOld.ContainerStatsResponse)(unsafe.Pointer(resp)), err
}

func (s *oldRuntimeServer) ContainerStatus(ctx context.Context, req *rtApiOld.ContainerStatusRequest) (*rtApiOld.ContainerStatusResponse, error) {
	resp, err := s.current.ContainerStatus(ctx, (*rtApi.ContainerStatusRequest)(unsafe.Pointer(req)))
	return (*rtApiOld.ContainerStatusResponse)(unsafe.Pointer(resp)), err
}

func (s *oldRuntimeServer) CreateContainer(ctx context.Context, req *rtApiOld.CreateContainerRequest) (*rtApiOld.CreateContainerResponse, error) {
	resp, err := s.current.CreateContainer(ctx, (*rtApi.CreateContainerRequest)(unsafe.Pointer(req)))
	return (*rtApiOld.CreateContainerResponse)(unsafe.Pointer(resp)), err
}

func (s *oldRuntimeServer) Exec(ctx context.Context, req *rtApiOld.ExecRequest) (*rtApiOld.ExecResponse, error) {
	resp, err := s.current.Exec(ctx, (*rtApi.ExecRequest)(unsafe.Pointer(req)))
	return (*rtApiOld.ExecResponse)(unsafe.Pointer(resp)), err
}

func (s *oldRuntimeServer) ExecSync(ctx context.Context, req *rtApiOld.ExecSyncRequest) (*rtApiOld.ExecSyncResponse, error) {
	resp, err := s.current.ExecSync(ctx, (*rtApi.ExecSyncRequest)(unsafe.Pointer(req)))
	return (*rtApiOld.ExecSyncResponse)(unsafe.Pointer(resp)), err
}

func (s *oldImageServer) ImageFsInfo(ctx context.Context, req *rtApiOld.ImageFsInfoRequest) (*rtApiOld.ImageFsInfoResponse, error) {
	resp, err := s.current.ImageFsInfo(ctx, (*rtApi.ImageFsInfoRequest)(unsafe.Pointer(req)))
	return (*rtApiOld.ImageFsInfoResponse)(unsafe.Pointer(resp)), err
}

func (s *oldImageServer) ImageStatus(ctx context.Context, req *rtApiOld.ImageStatusRequest) (*rtApiOld.ImageStatusResponse, error) {
	resp, err := s.current.ImageStatus(ctx, (*rtApi.ImageStatusRequest)(unsafe.Pointer(req)))
	return (*rtApiOld.ImageStatusResponse)(unsafe.Pointer(resp)), err
}

func (s *oldRuntimeServer) ListContainerStats(ctx context.Context, req *rtApiOld.ListContainerStatsRequest) (*rtApiOld.ListContainerStatsResponse, error) {
	resp, err := s.current.ListContainerStats(ctx, (*rtApi.ListContainerStatsRequest)(unsafe.Pointer(req)))
	return (*rtApiOld.ListContainerStatsResponse)(unsafe.Pointer(resp)), err
}

func (s *oldRuntimeServer) ListContainers(ctx context.Context, req *rtApiOld.ListContainersRequest) (*rtApiOld.ListContainersResponse, error) {
	resp, err := s.current.ListContainers(ctx, (*rtApi.ListContainersRequest)(unsafe.Pointer(req)))
	return (*rtApiOld.ListContainersResponse)(unsafe.Pointer(resp)), err
}

func (s *oldImageServer) ListImages(ctx context.Context, req *rtApiOld.ListImagesRequest) (*rtApiOld.ListImagesResponse, error) {
	resp, err := s.current.ListImages(ctx, (*rtApi.ListImagesRequest)(unsafe.Pointer(req)))
	return (*rtApiOld.ListImagesResponse)(unsafe.Pointer(resp)), err
}

func (s *oldRuntimeServer) ListPodSandbox(ctx context.Context, req *rtApiOld.ListPodSandboxRequest) (*rtApiOld.ListPodSandboxResponse, error) {
	resp, err := s.current.ListPodSandbox(ctx, (*rtApi.ListPodSandboxRequest)(unsafe.Pointer(req)))
	return (*rtApiOld.ListPodSandboxResponse)(unsafe.Pointer(resp)), err
}

func (s *oldRuntimeServer) ListPodSandboxStats(ctx context.Context, req *rtApiOld.ListPodSandboxStatsRequest) (*rtApiOld.ListPodSandboxStatsResponse, error) {
	resp, err := s.current.ListPodSandboxStats(ctx, (*rtApi.ListPodSandboxStatsRequest)(unsafe.Pointer(req)))
	return (*rtApiOld.ListPodSandboxStatsResponse)(unsafe.Pointer(resp)), err
}

func (s *oldRuntimeServer) PodSandboxStats(ctx context.Context, req *rtApiOld.PodSandboxStatsRequest) (*rtApiOld.PodSandboxStatsResponse, error) {
	resp, err := s.current.PodSandboxStats(ctx, (*rtApi.PodSandboxStatsRequest)(unsafe.Pointer(req)))
	return (*rtApiOld.PodSandboxStatsResponse)(unsafe.Pointer(resp)), err
}

func (s *oldRuntimeServer) PodSandboxStatus(ctx context.Context, req *rtApiOld.PodSandboxStatusRequest) (*rtApiOld.PodSandboxStatusResponse, error) {
	resp, err := s.current.PodSandboxStatus(ctx, (*rtApi.PodSandboxStatusRequest)(unsafe.Pointer(req)))
	return (*rtApiOld.PodSandboxStatusResponse)(unsafe.Pointer(resp)), err
}

func (s *oldRuntimeServer) PortForward(ctx context.Context, req *rtApiOld.PortForwardRequest) (*rtApiOld.PortForwardResponse, error) {
	resp, err := s.current.PortForward(ctx, (*rtApi.PortForwardRequest)(unsafe.Pointer(req)))
	return (*rtApiOld.PortForwardResponse)(unsafe.Pointer(resp)), err
}

func (s *oldImageServer) PullImage(ctx context.Context, req *rtApiOld.PullImageRequest) (*rtApiOld.PullImageResponse, error) {
	resp, err := s.current.PullImage(ctx, (*rtApi.PullImageRequest)(unsafe.Pointer(req)))
	return (*rtApiOld.PullImageResponse)(unsafe.Pointer(resp)), err
}

func (s *oldRuntimeServer) RemoveContainer(ctx context.Context, req *rtApiOld.RemoveContainerRequest) (*rtApiOld.RemoveContainerResponse, error) {
	resp, err := s.current.RemoveContainer(ctx, (*rtApi.RemoveContainerRequest)(unsafe.Pointer(req)))
	return (*rtApiOld.RemoveContainerResponse)(unsafe.Pointer(resp)), err
}

func (s *oldImageServer) RemoveImage(ctx context.Context, req *rtApiOld.RemoveImageRequest) (*rtApiOld.RemoveImageResponse, error) {
	resp, err := s.current.RemoveImage(ctx, (*rtApi.RemoveImageRequest)(unsafe.Pointer(req)))
	return (*rtApiOld.RemoveImageResponse)(unsafe.Pointer(resp)), err
}

func (s *oldRuntimeServer) RemovePodSandbox(ctx context.Context, req *rtApiOld.RemovePodSandboxRequest) (*rtApiOld.RemovePodSandboxResponse, error) {
	resp, err := s.current.RemovePodSandbox(ctx, (*rtApi.RemovePodSandboxRequest)(unsafe.Pointer(req)))
	return (*rtApiOld.RemovePodSandboxResponse)(unsafe.Pointer(resp)), err
}

func (s *oldRuntimeServer) ReopenContainerLog(ctx context.Context, req *rtApiOld.ReopenContainerLogRequest) (*rtApiOld.ReopenContainerLogResponse, error) {
	resp, err := s.current.ReopenContainerLog(ctx, (*rtApi.ReopenContainerLogRequest)(unsafe.Pointer(req)))
	return (*rtApiOld.ReopenContainerLogResponse)(unsafe.Pointer(resp)), err
}

func (s *oldRuntimeServer) RunPodSandbox(ctx context.Context, req *rtApiOld.RunPodSandboxRequest) (*rtApiOld.RunPodSandboxResponse, error) {
	resp, err := s.current.RunPodSandbox(ctx, (*rtApi.RunPodSandboxRequest)(unsafe.Pointer(req)))
	return (*rtApiOld.RunPodSandboxResponse)(unsafe.Pointer(resp)), err
}

func (s *oldRuntimeServer) StartContainer(ctx context.Context, req *rtApiOld.StartContainerRequest) (*rtApiOld.StartContainerResponse, error) {
	resp, err := s.current.StartContainer(ctx, (*rtApi.StartContainerRequest)(unsafe.Pointer(req)))
	return (*rtApiOld.StartContainerResponse)(unsafe.Pointer(resp)), err
}

func (s *oldRuntimeServer) Status(ctx context.Context, req *rtApiOld.StatusRequest) (*rtApiOld.StatusResponse, error) {
	resp, err := s.current.Status(ctx, (*rtApi.StatusRequest)(unsafe.Pointer(req)))
	return (*rtApiOld.StatusResponse)(unsafe.Pointer(resp)), err
}

func (s *oldRuntimeServer) StopContainer(ctx context.Context, req *rtApiOld.StopContainerRequest) (*rtApiOld.StopContainerResponse, error) {
	resp, err := s.current.StopContainer(ctx, (*rtApi.StopContainerRequest)(unsafe.Pointer(req)))
	return (*rtApiOld.StopContainerResponse)(unsafe.Pointer(resp)), err
}

func (s *oldRuntimeServer) StopPodSandbox(ctx context.Context, req *rtApiOld.StopPodSandboxRequest) (*rtApiOld.StopPodSandboxResponse, error) {
	resp, err := s.current.StopPodSandbox(ctx, (*rtApi.StopPodSandboxRequest)(unsafe.Pointer(req)))
	return (*rtApiOld.StopPodSandboxResponse)(unsafe.Pointer(resp)), err
}

func (s *oldRuntimeServer) UpdateContainerResources(ctx context.Context, req *rtApiOld.UpdateContainerResourcesRequest) (*rtApiOld.UpdateContainerResourcesResponse, error) {
	resp, err := s.current.UpdateContainerResources(ctx, (*rtApi.UpdateContainerResourcesRequest)(unsafe.Pointer(req)))
	return (*rtApiOld.UpdateContainerResourcesResponse)(unsafe.Pointer(resp)), err
}

func (s *oldRuntimeServer) UpdateRuntimeConfig(ctx context.Context, req *rtApiOld.UpdateRuntimeConfigRequest) (*rtApiOld.UpdateRuntimeConfigResponse, error) {
	resp, err := s.current.UpdateRuntimeConfig(ctx, (*rtApi.UpdateRuntimeConfigRequest)(unsafe.Pointer(req)))
	return (*rtApiOld.UpdateRuntimeConfigResponse)(unsafe.Pointer(resp)), err
}

func (s *oldRuntimeServer) Version(ctx context.Context, req *rtApiOld.VersionRequest) (*rtApiOld.VersionResponse, error) {
	resp, err := s.current.Version(ctx, (*rtApi.VersionRequest)(unsafe.Pointer(req)))
	return (*rtApiOld.VersionResponse)(unsafe.Pointer(resp)), err
}
