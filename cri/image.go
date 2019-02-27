package cri

import (
	"github.com/lxc/lxd/lxc/config"
	"github.com/lxc/lxd/shared/logger"
	"github.com/lxc/lxe/lxf"
	"golang.org/x/net/context"

	rtApi "k8s.io/kubernetes/pkg/kubelet/apis/cri/runtime/v1alpha2"
)

// ImageServer is the PoC implementation of the CRI ImageServer
type ImageServer struct {
	rtApi.ImageServiceServer
	lxdConfig     *config.Config
	criConfig     *LXEConfig
	runtimeRemote string
	lxf           *lxf.Client
}

// NewImageServer returns a new ImageServer backed by LXD
// we only need one connection — until we start distinguishing runtime & image service
func NewImageServer(s *RuntimeServer, lxf *lxf.Client) (*ImageServer, error) {
	i := ImageServer{
		lxdConfig: s.lxdConfig,
		criConfig: s.criConfig,
		lxf:       lxf,
	}
	// apply default image remote
	i.runtimeRemote = i.lxdConfig.DefaultRemote

	configPath, err := getLXDConfigPath(i.criConfig)
	if err != nil {
		return nil, err
	}
	i.lxdConfig, err = config.LoadConfig(configPath)
	if err != nil {
		return nil, err
	}
	i.lxdConfig.DefaultRemote = s.criConfig.LXDImageRemote
	return &i, nil
}

// ListImages lists existing images.
func (s ImageServer) ListImages(ctx context.Context, req *rtApi.ListImagesRequest) (*rtApi.ListImagesResponse, error) {
	logger.Debugf("ListImages(%v) triggered", req)

	var response rtApi.ListImagesResponse
	imglist, err := s.lxf.ListImages(req.GetFilter().GetImage().GetImage())
	if err != nil {
		logger.Errorf("Unable to list images, %v", err)
		return nil, err
	}

	for _, imgInfo := range imglist {
		rspImage := &rtApi.Image{
			Id:    imgInfo.Hash,
			Size_: uint64(imgInfo.Size),
			RepoDigests: []string{
				imgInfo.Hash,
			},
			RepoTags: imgInfo.Aliases,
		}
		response.Images = append(response.Images, rspImage)
	}
	return &response, nil
}

// ImageStatus returns the status of the image. If the image is not
// present, returns a response with ImageStatusResponse.Image set to
// nil.
func (s ImageServer) ImageStatus(ctx context.Context, req *rtApi.ImageStatusRequest) (*rtApi.ImageStatusResponse, error) {
	logger.Debugf("ImageStatus(%v) triggered", req)

	img, err := s.lxf.GetImage(req.GetImage().GetImage())
	if err != nil {
		// If the image can't be found, return no error with empty result
		if lxf.IsImageNotFound(err) {
			return &rtApi.ImageStatusResponse{}, nil
		}
		logger.Errorf("failed to get image status %v, %v", req.GetImage().GetImage(), err)
		return nil, err
	}
	return &rtApi.ImageStatusResponse{Image: &rtApi.Image{
		Id:    img.Hash,
		Size_: uint64(img.Size),
		RepoDigests: []string{
			img.Hash,
		},
		RepoTags: img.Aliases,
	}}, nil
}

// TODO
// 1. not impl: auth
// 1b. Authentication is provided in the pull request

// PullImage pulls an image with authentication config.
func (s ImageServer) PullImage(ctx context.Context, req *rtApi.PullImageRequest) (*rtApi.PullImageResponse, error) {
	logger.Debugf("PullImage(%v) triggered", req)

	hash, err := s.lxf.PullImage(req.GetImage().GetImage())
	if err != nil {
		logger.Errorf("failed to pull image %v, %v", req.GetImage().GetImage(), err)
		return nil, err
	}

	return &rtApi.PullImageResponse{
		ImageRef: hash,
	}, nil
}

// RemoveImage removes the image.
// This call is idempotent, and must not return an error if the image has
// already been removed.
func (s ImageServer) RemoveImage(ctx context.Context, req *rtApi.RemoveImageRequest) (*rtApi.RemoveImageResponse, error) {
	logger.Debugf("RemoveImage(%q) triggered", req)

	err := s.lxf.RemoveImage(req.GetImage().GetImage())
	if err != nil {
		logger.Errorf("failed to remove image %v, %v", req.GetImage().GetImage(), err)
		return nil, err
	}

	return &rtApi.RemoveImageResponse{}, nil
}

// ImageFsInfo returns information of the filesystem that is used to store images.
func (s ImageServer) ImageFsInfo(ctx context.Context, req *rtApi.ImageFsInfoRequest) (*rtApi.ImageFsInfoResponse, error) {
	logger.Debugf("ImageFsInfo(%v) triggered", req)

	poolUsage, err := s.lxf.GetFSPoolUsage()
	if err != nil {
		logger.Errorf("ImageFsInfo: GetFSPoolUsage(): %v", err)
		return nil, err
	}
	response := rtApi.ImageFsInfoResponse{}
	for _, i := range poolUsage {
		fs := &rtApi.FilesystemUsage{
			Timestamp:  i.Timestamp,
			FsId:       &rtApi.FilesystemIdentifier{Mountpoint: i.FsID},
			UsedBytes:  &rtApi.UInt64Value{Value: i.UsedBytes},
			InodesUsed: &rtApi.UInt64Value{Value: i.InodesUsed},
		}
		response.ImageFilesystems = append(response.ImageFilesystems, fs)
	}
	return &response, nil
}
