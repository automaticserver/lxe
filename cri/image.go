package cri

import (
	"strings"
	"time"

	"github.com/automaticserver/lxe/lxf"
	"github.com/lxc/lxd/lxc/config"
	sharedLXD "github.com/lxc/lxd/shared"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	rtApi "k8s.io/cri-api/pkg/apis/runtime/v1"
)

// ImageServer is the PoC implementation of the CRI ImageServer
type ImageServer struct {
	lxdConfig     *config.Config
	criConfig     *Config
	runtimeRemote string
	lxf           lxf.Client
}

// NewImageServer returns a new ImageServer backed by LXD
// we only need one connection — until we start distinguishing runtime & image service
func NewImageServer(s *RuntimeServer, lxf lxf.Client) (*ImageServer, error) {
	var err error

	i := ImageServer{
		lxdConfig: s.lxdConfig,
		criConfig: s.criConfig,
		lxf:       lxf,
	}
	// apply default image remote
	i.runtimeRemote = i.lxdConfig.DefaultRemote

	i.lxdConfig, err = config.LoadConfig(s.criConfig.LXDRemoteConfig)
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
		return nil, AnnErr(log, codes.Unknown, err, "Unable to list images")
	}

	for _, imgInfo := range imglist {
		rspImage := &rtApi.Image{
			Id:          imgInfo.Hash,
			Size_:       uint64(imgInfo.Size),
			RepoDigests: []string{},
			RepoTags:    []string{},
		}

		for _, a := range imgInfo.Aliases {
			if strings.Contains(a, "@sha256:") {
				rspImage.RepoDigests = append(rspImage.RepoDigests, convertLXEAliasNameToDocker(a))
			} else {
				rspImage.RepoTags = append(rspImage.RepoTags, convertLXEAliasNameToDocker(a))
			}
		}

		response.Images = append(response.Images, rspImage)
	}

	return response, nil
}

// ImageStatus returns the status of the image. If the image is not
// present, returns a response with ImageStatusResponse.Image set to
// nil.
func (s ImageServer) ImageStatus(ctx context.Context, req *rtApi.ImageStatusRequest) (*rtApi.ImageStatusResponse, error) {
	image := convertDockerImageNameToLXD(req.GetImage().GetImage())
	log := log.WithContext(ctx).WithField("image", image)

	imgInfo, err := s.lxf.GetImage(image)
	if err != nil {
		// If the image can't be found, return no error with empty result
		if lxf.IsNotFoundError(err) {
			return &rtApi.ImageStatusResponse{}, nil
		}

		return nil, AnnErr(log, codes.Unknown, err, "failed to get image status")
	}

	rspImage := &rtApi.Image{
		Id:          imgInfo.Hash,
		Size_:       uint64(imgInfo.Size),
		RepoDigests: []string{},
		RepoTags:    []string{},
	}

	for _, a := range imgInfo.Aliases {
		if strings.Contains(a, "@sha256") {
			rspImage.RepoDigests = append(rspImage.RepoDigests, convertLXEAliasNameToDocker(a))
		} else {
			rspImage.RepoTags = append(rspImage.RepoTags, convertLXEAliasNameToDocker(a))
		}
	}

	return &rtApi.ImageStatusResponse{Image: rspImage}, nil
}

// TODO
// 1. not impl: auth
// 1b. Authentication is provided in the pull request

// PullImage pulls an image with authentication config.
func (s ImageServer) PullImage(ctx context.Context, req *rtApi.PullImageRequest) (*rtApi.PullImageResponse, error) {
	image := convertDockerImageNameToLXD(req.GetImage().GetImage())
	log := log.WithContext(ctx).WithField("image", image)

	hash, err := s.lxf.PullImage(image)
	if err != nil {
		return nil, AnnErr(log, codes.Unknown, err, "failed to pull image")
	}

	return &rtApi.PullImageResponse{ImageRef: hash}, nil
}

// RemoveImage removes the image.
// This call is idempotent, and must not return an error if the image has
// already been removed.
func (s ImageServer) RemoveImage(ctx context.Context, req *rtApi.RemoveImageRequest) (*rtApi.RemoveImageResponse, error) {
	image := convertDockerImageNameToLXD(req.GetImage().GetImage())
	log := log.WithContext(ctx).WithField("image", image)

	err := s.lxf.RemoveImage(image)
	if err != nil && !lxf.IsNotFoundError(err) {
		return nil, AnnErr(log, codes.Unknown, err, "failed to remove image")
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
