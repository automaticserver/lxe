package cri // import "github.com/automaticserver/lxe/cri"

import (
	"time"

	"github.com/automaticserver/lxe/lxf"
	"github.com/automaticserver/lxe/shared"
	"github.com/lxc/lxd/lxc/config"
	sharedLXD "github.com/lxc/lxd/shared"
	"golang.org/x/net/context"
	rtApi "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

// ImageServer is the PoC implementation of the CRI ImageServer
type ImageServer struct {
	rtApi.ImageServiceServer
	lxdConfig     *config.Config
	criConfig     *Config
	runtimeRemote string
	lxf           lxf.Client
}

// NewImageServer returns a new ImageServer backed by LXD
// we only need one connection — until we start distinguishing runtime & image service
func NewImageServer(s *RuntimeServer, lxf lxf.Client) (*ImageServer, error) {
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
	log := log.WithContext(ctx).WithField("filter", req.GetFilter().GetImage())

	response := &rtApi.ListImagesResponse{}

	imglist, err := s.lxf.ListImages(req.GetFilter().GetImage().GetImage())
	if err != nil {
		return nil, AnnErr(log, err, "Unable to list images")
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

	return response, nil
}

// ImageStatus returns the status of the image. If the image is not
// present, returns a response with ImageStatusResponse.Image set to
// nil.
func (s ImageServer) ImageStatus(ctx context.Context, req *rtApi.ImageStatusRequest) (*rtApi.ImageStatusResponse, error) {
	log := log.WithContext(ctx).WithField("image", req.GetImage().GetImage())

	img, err := s.lxf.GetImage(req.GetImage().GetImage())
	if err != nil {
		// If the image can't be found, return no error with empty result
		if shared.IsErrNotFound(err) {
			return &rtApi.ImageStatusResponse{}, nil
		}

		return nil, AnnErr(log, err, "failed to get image status")
	}

	response := &rtApi.ImageStatusResponse{Image: &rtApi.Image{
		Id:    img.Hash,
		Size_: uint64(img.Size),
		RepoDigests: []string{
			img.Hash,
		},
		RepoTags: img.Aliases,
	}}

	return response, nil
}

// TODO
// 1. not impl: auth
// 1b. Authentication is provided in the pull request

// PullImage pulls an image with authentication config.
func (s ImageServer) PullImage(ctx context.Context, req *rtApi.PullImageRequest) (*rtApi.PullImageResponse, error) {
	log := log.WithContext(ctx).WithField("image", req.GetImage().GetImage())

	hash, err := s.lxf.PullImage(req.GetImage().GetImage())
	if err != nil {
		return nil, AnnErr(log, err, "failed to pull image")
	}

	response := &rtApi.PullImageResponse{
		ImageRef: hash,
	}

	return response, nil
}

// RemoveImage removes the image.
// This call is idempotent, and must not return an error if the image has
// already been removed.
func (s ImageServer) RemoveImage(ctx context.Context, req *rtApi.RemoveImageRequest) (*rtApi.RemoveImageResponse, error) {
	log := log.WithContext(ctx).WithField("image", req.GetImage().GetImage())

	err := s.lxf.RemoveImage(req.GetImage().GetImage())
	if err != nil {
		return nil, AnnErr(log, err, "failed to remove image")
	}

	return &rtApi.RemoveImageResponse{}, nil
}

// ImageFsInfo returns information of the filesystem that is used to store images.
func (s ImageServer) ImageFsInfo(ctx context.Context, req *rtApi.ImageFsInfoRequest) (*rtApi.ImageFsInfoResponse, error) {
	// log := log.WithContext(ctx)
	// Images are not saved in pools (for now?)
	// poolUsage, err := s.lxf.GetFSPoolUsage()
	// if err != nil {
	// 	return nil, err
	// }
	response := &rtApi.ImageFsInfoResponse{}
	// for _, i := range poolUsage {
	// 	fs := &rtApi.FilesystemUsage{
	// 		Timestamp:  i.Timestamp,
	// 		FsId:       &rtApi.FilesystemIdentifier{Mountpoint: i.FsID},
	// 		UsedBytes:  &rtApi.UInt64Value{Value: i.UsedBytes},
	// 		InodesUsed: &rtApi.UInt64Value{Value: i.InodesUsed},
	// 	}
	// 	response.ImageFilesystems = append(response.ImageFilesystems, fs)
	// }

	// TODO: UsedBytes, InodesUsed
	response.ImageFilesystems = append(response.ImageFilesystems, &rtApi.FilesystemUsage{
		Timestamp:  time.Now().UnixNano(),
		FsId:       &rtApi.FilesystemIdentifier{Mountpoint: sharedLXD.VarPath("images")},
		UsedBytes:  &rtApi.UInt64Value{Value: 0},
		InodesUsed: &rtApi.UInt64Value{Value: 0},
	})

	return response, nil
}
