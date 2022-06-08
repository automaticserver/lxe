package cri

import (
	"bytes"
	"errors"
	"io/ioutil"
	"path"
	"strconv"
	"strings"

	"github.com/automaticserver/lxe/cli/version"
	"github.com/automaticserver/lxe/lxf"
	"github.com/automaticserver/lxe/lxf/device"
	"github.com/automaticserver/lxe/network"
	"github.com/automaticserver/lxe/third_party/ioutils"
	"github.com/lxc/lxd/lxc/config"
	opencontainers "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	utilNet "k8s.io/apimachinery/pkg/util/net"
	rtApi "k8s.io/cri-api/pkg/apis/runtime/v1"
)

const (
	criVersion = "0.1.0"
)

var (
	ErrNotImplemented       = errors.New("not implemented")
	ErrUnknownNetworkPlugin = errors.New("unknown network plugin")
)

// RuntimeServer is the PoC implementation of the CRI RuntimeServer
type RuntimeServer struct {
	rtApi.RuntimeServiceServer
	lxf       lxf.Client
	stream    *streamService
	lxdConfig *config.Config
	criConfig *Config
	network   network.Plugin
}

// NewRuntimeServer returns a new RuntimeServer backed by LXD
func NewRuntimeServer(criConfig *Config, lxf lxf.Client, network network.Plugin) (*RuntimeServer, error) {
	var err error

	runtime := RuntimeServer{
		criConfig: criConfig,
		network:   network,
	}

	runtime.lxdConfig, err = config.LoadConfig(criConfig.LXDRemoteConfig)
	if err != nil {
		return nil, err
	}

	runtime.lxf = lxf

	return &runtime, nil
}

// Version returns the runtime name, runtime version, and runtime API version.
func (s RuntimeServer) Version(ctx context.Context, req *rtApi.VersionRequest) (*rtApi.VersionResponse, error) {
	log := log.WithContext(ctx).WithField("version", req.GetVersion())

	// According to containerd CRI implementation RuntimeName=ShimName, RuntimeVersion=ShimVersion,
	// RuntimeApiVersion=someAPIVersion. The actual runtime name and version is not present
	info, err := s.lxf.GetRuntimeInfo()
	if err != nil {
		return nil, AnnErr(log, codes.Unknown, err, "unable to get server environment")
	}

	response := &rtApi.VersionResponse{
		Version:           criVersion,
		RuntimeName:       Domain,
		RuntimeVersion:    version.Version,
		RuntimeApiVersion: info.Version,
	}

	return response, nil
}

// RunPodSandbox creates and starts a pod-level sandbox. Runtimes must ensure the sandbox is in the ready state on
// success
func (s RuntimeServer) RunPodSandbox(ctx context.Context, req *rtApi.RunPodSandboxRequest) (*rtApi.RunPodSandboxResponse, error) { // nolint: gocognit, cyclop
	log := log.WithContext(ctx).WithFields(logrus.Fields{
		"podname":   req.GetConfig().GetMetadata().GetName(),
		"namespace": req.GetConfig().GetMetadata().GetNamespace(),
		"poduid":    req.GetConfig().GetMetadata().GetUid(),
	})
	log.Info("run pod")

	var err error

	sb := s.lxf.NewSandbox()

	sb.Hostname = req.GetConfig().GetHostname()
	sb.LogDirectory = req.GetConfig().GetLogDirectory()
	meta := req.GetConfig().GetMetadata()
	sb.Metadata = lxf.SandboxMetadata{
		Attempt:   meta.GetAttempt(),
		Name:      meta.GetName(),
		Namespace: meta.GetNamespace(),
		UID:       meta.GetUid(),
	}
	sb.Labels = req.GetConfig().GetLabels()
	sb.Annotations = req.GetConfig().GetAnnotations()

	if req.GetConfig().GetDnsConfig() != nil {
		sb.NetworkConfig.Nameservers = req.GetConfig().GetDnsConfig().GetServers()
		sb.NetworkConfig.Searches = req.GetConfig().GetDnsConfig().GetSearches()
	}

	// Find out which network mode should be used
	if strings.ToLower(req.GetConfig().GetLinux().GetSecurityContext().GetNamespaceOptions().GetNetwork().String()) == string(lxf.NetworkHost) {
		// host network explicitly requested
		sb.NetworkConfig.Mode = lxf.NetworkHost
		lxf.AppendIfSet(&sb.Config, "raw.lxc", "lxc.include = "+s.criConfig.LXEHostnetworkFile)
	} else {
		// manage network according to selected network plugin
		// TODO: we could omit these since we use network plugin, but we still need to remember if it is HostNetwork
		switch s.criConfig.LXENetworkPlugin {
		case NetworkPluginBridge:
			sb.NetworkConfig.Mode = lxf.NetworkBridged
		case NetworkPluginCNI:
			sb.NetworkConfig.Mode = lxf.NetworkCNI
		default:
			// unknown plugin name provided
			return nil, AnnErr(log, codes.Unknown, ErrUnknownNetworkPlugin, s.criConfig.LXENetworkPlugin)
		}
	}

	// If HostPort is defined, set forwardings from that port to the container. In lxd, we can use proxy devices for that.
	// This can be applied to all NetworkModes except HostNetwork.
	if sb.NetworkConfig.Mode != lxf.NetworkHost {
		for _, portMap := range req.Config.PortMappings {
			// both HostPort and ContainerPort must be defined, otherwise invalid
			if portMap.GetHostPort() == 0 || portMap.GetContainerPort() == 0 {
				continue
			}

			hostPort := int(portMap.GetHostPort())
			containerPort := int(portMap.GetContainerPort())

			var protocol device.Protocol

			switch portMap.GetProtocol() { // nolint: exhaustive
			case rtApi.Protocol_UDP:
				protocol = device.ProtocolUDP
			case rtApi.Protocol_TCP:
				fallthrough
			default:
				protocol = device.ProtocolTCP
			}

			hostIP := portMap.GetHostIp()
			if hostIP == "" {
				hostIP = "0.0.0.0"
			}

			containerIP := "127.0.0.1"

			sb.Devices.Upsert(&device.Proxy{
				Listen: &device.ProxyEndpoint{
					Protocol: protocol,
					Address:  hostIP,
					Port:     hostPort,
				},
				Destination: &device.ProxyEndpoint{
					Protocol: protocol,
					Address:  containerIP,
					Port:     containerPort,
				},
			})
		}
	}

	// TODO: Refactor...
	if req.Config.Linux != nil { // nolint: nestif
		lxf.SetIfSet(&sb.Config, "user.linux.cgroup_parent", req.Config.Linux.CgroupParent)

		for key, value := range req.Config.Linux.Sysctls {
			sb.Config["user.linux.sysctls."+key] = value
		}

		if req.Config.Linux.SecurityContext != nil {
			privileged := req.Config.Linux.SecurityContext.Privileged
			sb.Config["user.linux.security_context.privileged"] = strconv.FormatBool(privileged)
			sb.Config["security.privileged"] = strconv.FormatBool(privileged)

			if req.Config.Linux.SecurityContext.NamespaceOptions != nil {
				nsi := "user.linux.security_context.namespace_options"
				nso := req.Config.Linux.SecurityContext.NamespaceOptions

				sb.Config[nsi+".ipc"] = nameSpaceOptionToString(nso.Ipc)
				sb.Config[nsi+".network"] = nameSpaceOptionToString(nso.Network)
				sb.Config[nsi+".pid"] = nameSpaceOptionToString(nso.Pid)
			}

			if req.Config.Linux.SecurityContext.ReadonlyRootfs {
				sb.Devices.Upsert(&device.Disk{
					Path:     "/",
					Readonly: true,
					// TODO magic constant, and also, is it always default?
					Pool: "default",
				})
			}

			if req.Config.Linux.SecurityContext.RunAsUser != nil {
				sb.Config["user.linux.security_context.run_as_user"] =
					strconv.FormatInt(req.Config.Linux.SecurityContext.RunAsUser.Value, 10)
			}

			lxf.SetIfSet(&sb.Config, "user.linux.security_context.seccomp_profile_path",
				req.Config.Linux.SecurityContext.SeccompProfilePath)

			if req.Config.Linux.SecurityContext.SelinuxOptions != nil {
				sci := "user.linux.security_context.namespace_options"
				sco := req.Config.Linux.SecurityContext.SelinuxOptions
				lxf.SetIfSet(&sb.Config, sci+".role", sco.Role)
				lxf.SetIfSet(&sb.Config, sci+".level", sco.Level)
				lxf.SetIfSet(&sb.Config, sci+".user", sco.User)
				lxf.SetIfSet(&sb.Config, sci+".type", sco.Type)
			}
		}
	}

	err = sb.Apply()
	if err != nil {
		return nil, AnnErr(log, codes.Unknown, err, "failed to create pod")
	}

	log = log.WithField("podid", sb.ID)

	// create network
	if sb.NetworkConfig.Mode != lxf.NetworkHost { // nolint: nestif
		podNet, err := s.network.PodNetwork(sb.ID, sb.Annotations)
		if err != nil {
			return nil, AnnErr(log, codes.Unknown, err, "can't enter pod network context")
		}

		res, err := podNet.WhenCreated(ctx, &network.Properties{})
		if err != nil {
			return nil, AnnErr(log, codes.Unknown, err, "can't create pod network")
		}

		err = s.handleNetworkResult(sb, res)
		if err != nil {
			return nil, AnnErr(log, codes.Unknown, err, "unable to save pod network result")
		}

		// Since a PodSandbox is created "started", also fire started network
		res, err = podNet.WhenStarted(ctx, &network.PropertiesRunning{
			Properties: network.Properties{
				Data: sb.NetworkConfig.ModeData,
			},
			Pid: 0, // if we had real 1:n pod:container we would add here the pid of the pod process
		})
		if err != nil {
			return nil, AnnErr(log, codes.Unknown, err, "can't start pod network")
		}

		err = s.handleNetworkResult(sb, res)
		if err != nil {
			return nil, AnnErr(log, codes.Unknown, err, "unable to save start pod network result")
		}
	}

	log.Info("run pod successful")

	return &rtApi.RunPodSandboxResponse{PodSandboxId: sb.ID}, nil
}

// StopPodSandbox stops any running process that is part of the sandbox and reclaims network resources (e.g. IP
// addresses) allocated to the sandbox. If there are any running containers in the sandbox, they must be forcibly
// terminated. This call is idempotent, and must not return an error if all relevant resources have already been
// reclaimed. kubelet will call StopPodSandbox at least once before calling RemovePodSandbox. It will also attempt to
// reclaim resources eagerly, as soon as a sandbox is not needed. Hence, multiple StopPodSandbox calls are expected.
func (s RuntimeServer) StopPodSandbox(ctx context.Context, req *rtApi.StopPodSandboxRequest) (*rtApi.StopPodSandboxResponse, error) {
	log := log.WithContext(ctx).WithField("podid", req.GetPodSandboxId())
	log.Info("stop pod")

	sb, err := s.lxf.GetSandbox(req.GetPodSandboxId())
	if err != nil {
		// If the sandbox can't be found, return no error with empty result
		if lxf.IsNotFoundError(err) {
			return &rtApi.StopPodSandboxResponse{}, nil
		}

		return nil, AnnErr(log, codes.Unknown, err, "unable to get pod")
	}

	err = s.stopContainers(sb)
	if err != nil {
		return nil, AnnErr(log, codes.Unknown, err, "unable to stop containers")
	}

	err = s.stopSandbox(ctx, sb)
	if err != nil {
		return nil, AnnErr(log, codes.Unknown, err, "unable to stop pod")
	}

	log.Info("stop pod successful")

	return &rtApi.StopPodSandboxResponse{}, nil
}

// RemovePodSandbox removes the sandbox. This is pretty much the same as StopPodSandbox but also removes the sandbox and
// the containers
func (s RuntimeServer) RemovePodSandbox(ctx context.Context, req *rtApi.RemovePodSandboxRequest) (*rtApi.RemovePodSandboxResponse, error) {
	log := log.WithContext(ctx).WithField("podid", req.GetPodSandboxId())
	log.Info("remove pod")

	sb, err := s.lxf.GetSandbox(req.GetPodSandboxId())
	if err != nil {
		// If the sandbox can't be found, return no error with empty result
		if lxf.IsNotFoundError(err) {
			return &rtApi.RemovePodSandboxResponse{}, nil
		}

		return nil, AnnErr(log, codes.Unknown, err, "unable to get pod")
	}

	err = s.stopContainers(sb)
	if err != nil {
		return nil, AnnErr(log, codes.Unknown, err, "unable to stop containers")
	}

	err = s.deleteContainers(ctx, sb)
	if err != nil {
		return nil, AnnErr(log, codes.Unknown, err, "unable to delete containers")
	}

	err = s.deleteSandbox(ctx, sb)
	if err != nil {
		return nil, AnnErr(log, codes.Unknown, err, "unable to delete pod")
	}

	log.Info("remove pod successful")

	return &rtApi.RemovePodSandboxResponse{}, nil
}

// PodSandboxStatus returns the status of the PodSandbox. If the PodSandbox is not present, returns an error.
func (s RuntimeServer) PodSandboxStatus(ctx context.Context, req *rtApi.PodSandboxStatusRequest) (*rtApi.PodSandboxStatusResponse, error) {
	log := log.WithContext(ctx).WithField("podid", req.GetPodSandboxId())

	sb, err := s.lxf.GetSandbox(req.GetPodSandboxId())
	if err != nil {
		if lxf.IsNotFoundError(err) {
			return nil, AnnErr(log, codes.NotFound, err, "pod not found")
		}

		return nil, AnnErr(log, codes.Unknown, err, "unable to get pod")
	}

	response := &rtApi.PodSandboxStatusResponse{
		Status: &rtApi.PodSandboxStatus{
			Id: sb.ID,
			Metadata: &rtApi.PodSandboxMetadata{
				Attempt:   sb.Metadata.Attempt,
				Name:      sb.Metadata.Name,
				Namespace: sb.Metadata.Namespace,
				Uid:       sb.Metadata.UID,
			},
			Linux:       &rtApi.LinuxPodSandboxStatus{},
			Labels:      sb.Labels,
			Annotations: sb.Annotations,
			CreatedAt:   sb.CreatedAt.UnixNano(),
			State:       stateSandboxAsCri(sb.State),
			Network: &rtApi.PodSandboxNetworkStatus{
				Ip: "",
			},
		},
	}

	for k, v := range sb.Config {
		if strings.HasPrefix(k, "user.linux.security_context.namespace_options.") {
			key := strings.TrimPrefix(k, "user.linux.security_context.namespace_options.")

			if response.Status.Linux.Namespaces == nil {
				response.Status.Linux.Namespaces = &rtApi.Namespace{Options: &rtApi.NamespaceOption{}}
			}

			switch key {
			case "ipc":
				response.Status.Linux.Namespaces.Options.Ipc = stringToNamespaceOption(v)
			case "pid":
				response.Status.Linux.Namespaces.Options.Pid = stringToNamespaceOption(v)
			case "network":
				response.Status.Linux.Namespaces.Options.Network = stringToNamespaceOption(v)
			}
		}
	}

	ip := s.getInetAddress(ctx, sb)
	if ip != "" {
		response.Status.Network.Ip = ip
	}

	return response, nil
}

// getInetAddress returns the ip address of the sandbox. empty string if nothing was found
func (s RuntimeServer) getInetAddress(ctx context.Context, sb *lxf.Sandbox) string { // nolint: cyclop
	log := log.WithContext(ctx).WithField("podid", sb.ID)

	switch sb.NetworkConfig.Mode {
	case lxf.NetworkHost:
		ip, err := utilNet.ChooseHostInterface()
		if err != nil {
			log.WithError(err).Error("Couldn't choose host interface")

			return ""
		}

		return ip.String()
	case lxf.NetworkNone:
		return ""
	case lxf.NetworkBridged:
		fallthrough
	case lxf.NetworkCNI:
		podNet, err := s.network.PodNetwork(sb.ID, sb.Annotations)
		if err != nil {
			log.WithError(err).Error("Couldn't get cni pod network")

			return ""
		}

		status, err := podNet.Status(ctx, &network.PropertiesRunning{Properties: network.Properties{Data: sb.NetworkConfig.ModeData}, Pid: 0})
		if err != nil {
			log.WithError(err).Error("Couldn't get status of cni pod network")

			return ""
		}

		if len(status.IPs) > 0 {
			return status.IPs[0].String()
		}
	}

	// If not yet returned, look into the containers interface list and select the address from the default interface
	// TODO: is this still needed? Look into network.Bridge as well
	cl, err := sb.Containers()
	if err != nil {
		log.WithError(err).Error("Couldn't list containers while trying to get inet address")

		return ""
	}

	for _, c := range cl {
		// ignore any non-running containers
		if c.StateName != lxf.ContainerStateRunning {
			continue
		}

		// get the ipv4 address of eth0
		ip := c.GetInetAddress([]string{network.DefaultInterface})
		if ip != "" {
			return ip
		}
	}

	return ""
}

// ListPodSandbox returns a list of PodSandboxes.
func (s RuntimeServer) ListPodSandbox(ctx context.Context, req *rtApi.ListPodSandboxRequest) (*rtApi.ListPodSandboxResponse, error) {
	log := log.WithContext(ctx).WithField("filter", req.GetFilter().String())

	sandboxes, err := s.lxf.ListSandboxes()
	if err != nil {
		return nil, AnnErr(log, codes.Unknown, err, "unable to list pods")
	}

	response := &rtApi.ListPodSandboxResponse{}

	for _, sb := range sandboxes {
		if req.GetFilter() != nil {
			filter := req.GetFilter()
			if filter.GetId() != "" && filter.GetId() != sb.ID {
				continue
			}

			if filter.GetState() != nil && filter.GetState().GetState() != stateSandboxAsCri(sb.State) {
				continue
			}

			if !CompareFilterMap(sb.Labels, filter.GetLabelSelector()) {
				continue
			}
		}

		// TODO: toSandboxCRI()
		pod := rtApi.PodSandbox{
			Id:        sb.ID,
			CreatedAt: sb.CreatedAt.UnixNano(),
			Metadata: &rtApi.PodSandboxMetadata{
				Attempt:   sb.Metadata.Attempt,
				Name:      sb.Metadata.Name,
				Namespace: sb.Metadata.Namespace,
				Uid:       sb.Metadata.UID,
			},
			State:       stateSandboxAsCri(sb.State),
			Labels:      sb.Labels,
			Annotations: sb.Annotations,
		}
		response.Items = append(response.Items, &pod)
	}

	return response, nil
}

// CreateContainer creates a new container in specified PodSandbox
func (s RuntimeServer) CreateContainer(ctx context.Context, req *rtApi.CreateContainerRequest) (*rtApi.CreateContainerResponse, error) { // nolint: cyclop
	image := convertDockerImageNameToLXD(req.GetConfig().GetImage().GetImage())
	log := log.WithContext(ctx).WithFields(logrus.Fields{
		"containername": req.GetConfig().GetMetadata().GetName(),
		"attempt":       req.GetConfig().GetMetadata().GetAttempt(),
		"podid":         req.GetPodSandboxId(),
		"image":         image,
	})
	log.Info("create container")

	img, err := s.lxf.GetImage(image)
	if err != nil {
		return nil, AnnErr(log, codes.Unknown, err, "failed to find image hash")
	}

	c := s.lxf.NewContainer(req.GetPodSandboxId(), s.criConfig.LXDProfiles...)
	c.Image = img.Hash
	c.Labels = req.GetConfig().GetLabels()
	c.Annotations = req.GetConfig().GetAnnotations()
	meta := req.GetConfig().GetMetadata()
	c.Metadata = lxf.ContainerMetadata{
		Attempt: meta.GetAttempt(),
		Name:    meta.GetName(),
	}
	c.LogPath = req.GetConfig().GetLogPath()

	for _, mnt := range req.GetConfig().GetMounts() {
		hostPath := mnt.GetHostPath()
		containerPath := mnt.GetContainerPath()
		// cannot use /var/run as most distros symlink that to /run and lxd doesn't like mounts there because of that
		if strings.HasPrefix(containerPath, "/var/run") {
			containerPath = path.Join("/run", strings.TrimPrefix(containerPath, "/var/run"))
		}
		// cannot use /run as most distros mount a tmpfs on top of that so mounts from lxd are not visible in the container
		if strings.HasPrefix(containerPath, "/run") {
			containerPath = path.Join("/mnt", strings.TrimPrefix(containerPath, "/run"))
		}

		c.Devices.Upsert(&device.Disk{
			Path:     containerPath,
			Source:   hostPath,
			Readonly: mnt.GetReadonly(),
			Optional: false,
		})
	}

	for _, dev := range req.GetConfig().GetDevices() {
		c.Devices.Upsert(&device.Block{
			Source: dev.GetHostPath(),
			Path:   dev.GetContainerPath(),
		})
	}

	c.Privileged = req.GetConfig().GetLinux().GetSecurityContext().GetPrivileged()

	// get metadata & cloud-init if defined
	for _, env := range req.GetConfig().GetEnvs() {
		switch {
		case env.GetKey() == "user-data":
			c.CloudInitUserData = env.GetValue()
		case env.GetKey() == "meta-data":
			c.CloudInitMetaData = env.GetValue()
		case env.GetKey() == "network-config":
			c.CloudInitNetworkConfig = env.GetValue()
		default:
			c.Environment[env.GetKey()] = env.GetValue()
		}
	}

	// append other envs below metadata
	if c.CloudInitMetaData != "" && len(c.Environment) > 0 {
		c.CloudInitMetaData += "\n"
	}

	// process limits
	resrc := req.GetConfig().GetLinux().GetResources()
	if resrc != nil {
		c.Resources = &opencontainers.LinuxResources{}
		c.Resources.CPU = &opencontainers.LinuxCPU{}
		c.Resources.Memory = &opencontainers.LinuxMemory{}
		shares := uint64(resrc.CpuShares)
		c.Resources.CPU.Shares = &shares
		c.Resources.CPU.Quota = &resrc.CpuQuota
		period := uint64(resrc.CpuPeriod)
		c.Resources.CPU.Period = &period
		c.Resources.Memory.Limit = &resrc.MemoryLimitInBytes
	}

	err = c.Apply()
	if err != nil {
		return nil, AnnErr(log, codes.Unknown, err, "unable to create container")
	}

	sb, err := c.Sandbox()
	if err != nil {
		return nil, AnnErr(log, codes.Unknown, err, "unable to find sandbox")
	}

	// create network
	if sb.NetworkConfig.Mode != lxf.NetworkHost {
		podNet, err := s.network.PodNetwork(sb.ID, sb.Annotations)
		if err != nil {
			return nil, AnnErr(log, codes.Unknown, err, "can't enter pod network context")
		}

		contNet, err := podNet.ContainerNetwork(c.ID, c.Annotations)
		if err != nil {
			return nil, AnnErr(log, codes.Unknown, err, "can't enter container network context")
		}

		res, err := contNet.WhenCreated(ctx, &network.Properties{})
		if err != nil {
			return nil, AnnErr(log, codes.Unknown, err, "can't create container network")
		}

		err = s.handleNetworkResult(sb, res)
		if err != nil {
			return nil, AnnErr(log, codes.Unknown, err, "unable to save create container network result")
		}
	}

	log.Info("create container successful")

	return &rtApi.CreateContainerResponse{ContainerId: c.ID}, nil
}

// StartContainer starts the container.
func (s RuntimeServer) StartContainer(ctx context.Context, req *rtApi.StartContainerRequest) (*rtApi.StartContainerResponse, error) {
	log := log.WithContext(ctx).WithField("containerid", req.GetContainerId())
	log.Info("start container")

	c, err := s.lxf.GetContainer(req.GetContainerId())
	if err != nil {
		return nil, AnnErr(log, codes.Unknown, err, "unable to get container")
	}

	err = c.Start()
	if err != nil {
		return nil, AnnErr(log, codes.Unknown, err, "unable to start container")
	}

	log.Info("start container successful")

	return &rtApi.StartContainerResponse{}, nil
}

// StopContainer stops a running container with a grace period (i.e., timeout). This call is idempotent, and must not
// return an error if the container has already been stopped.
func (s RuntimeServer) StopContainer(ctx context.Context, req *rtApi.StopContainerRequest) (*rtApi.StopContainerResponse, error) {
	log := log.WithContext(ctx).WithField("containerid", req.GetContainerId())
	log.Info("stop container")

	c, err := s.lxf.GetContainer(req.GetContainerId())
	if err != nil {
		if lxf.IsNotFoundError(err) {
			return &rtApi.StopContainerResponse{}, nil
		}

		return nil, AnnErr(log, codes.Unknown, err, "unable to get container")
	}

	err = s.stopContainer(c, int(req.Timeout))
	if err != nil {
		return nil, AnnErr(log, codes.Unknown, err, "unable to stop container")
	}

	log.Info("stop container successful")

	return &rtApi.StopContainerResponse{}, nil
}

// RemoveContainer removes the container. If the container is running, the container must be forcibly removed. This call
// is idempotent, and must not return an error if the container has already been removed. nolint: dupl
func (s RuntimeServer) RemoveContainer(ctx context.Context, req *rtApi.RemoveContainerRequest) (*rtApi.RemoveContainerResponse, error) {
	log := log.WithContext(ctx).WithField("containerid", req.GetContainerId())
	log.Info("remove container")

	c, err := s.lxf.GetContainer(req.GetContainerId())
	if err != nil {
		if lxf.IsNotFoundError(err) {
			return &rtApi.RemoveContainerResponse{}, nil
		}

		return nil, AnnErr(log, codes.Unknown, err, "unable to get container")
	}

	err = s.deleteContainer(ctx, c)
	if err != nil {
		return nil, AnnErr(log, codes.Unknown, err, "unable to remove container")
	}

	log.Info("remove container successful")

	return &rtApi.RemoveContainerResponse{}, nil
}

// ListContainers lists all containers by filters.
func (s RuntimeServer) ListContainers(ctx context.Context, req *rtApi.ListContainersRequest) (*rtApi.ListContainersResponse, error) { // nolint: cyclop
	log := log.WithContext(ctx).WithField("filter", req.GetFilter().String())

	response := &rtApi.ListContainersResponse{}

	cl, err := s.lxf.ListContainers()
	if err != nil {
		return nil, AnnErr(log, codes.Unknown, err, "unable to get container list")
	}

	for _, c := range cl {
		if req.GetFilter() != nil {
			filter := req.GetFilter()
			if filter.GetId() != "" && filter.GetId() != c.ID {
				continue
			}

			if filter.GetState() != nil && filter.GetState().GetState() != stateContainerAsCri(c.StateName) {
				continue
			}

			if filter.GetPodSandboxId() != "" && filter.GetPodSandboxId() != c.SandboxID() {
				continue
			}

			if !CompareFilterMap(c.Labels, filter.GetLabelSelector()) {
				continue
			}
		}

		response.Containers = append(response.Containers, toCriContainer(c))
	}

	return response, nil
}

// ContainerStatus returns status of the container. If the container is not present, returns an error.
func (s RuntimeServer) ContainerStatus(ctx context.Context, req *rtApi.ContainerStatusRequest) (*rtApi.ContainerStatusResponse, error) {
	log := log.WithContext(ctx).WithField("containerid", req.GetContainerId())

	ct, err := s.lxf.GetContainer(req.GetContainerId())
	if err != nil {
		if lxf.IsNotFoundError(err) {
			return nil, AnnErr(log, codes.NotFound, err, "container not found")
		}

		return nil, AnnErr(log, codes.Unknown, err, "unable to get container")
	}

	response := toCriStatusResponse(ct)

	return response, nil
}

// UpdateContainerResources updates ContainerConfig of the container.
func (s RuntimeServer) UpdateContainerResources(ctx context.Context, req *rtApi.UpdateContainerResourcesRequest) (*rtApi.UpdateContainerResourcesResponse, error) {
	return nil, SilErr(log, codes.Unimplemented, ErrNotImplemented, "")
}

// ReopenContainerLog asks runtime to reopen the stdout/stderr log file for the container. This is often called after
// the log file has been rotated. If the container is not running, container runtime can choose to either create a new
// log file and return nil, or return an error. Once it returns error, new container log file MUST NOT be created.
func (s RuntimeServer) ReopenContainerLog(ctx context.Context, req *rtApi.ReopenContainerLogRequest) (*rtApi.ReopenContainerLogResponse, error) {
	return nil, SilErr(log, codes.Unimplemented, ErrNotImplemented, "")
}

// ExecSync runs a command in a container synchronously.
func (s RuntimeServer) ExecSync(ctx context.Context, req *rtApi.ExecSyncRequest) (*rtApi.ExecSyncResponse, error) {
	log := log.WithContext(ctx).WithFields(logrus.Fields{
		"containerid": req.GetContainerId(),
		"cmd":         req.GetCmd(),
	})

	stdin := bytes.NewReader(nil)
	stdinR := ioutil.NopCloser(stdin)
	stdout := bytes.NewBuffer(nil)
	stdoutW := ioutils.WriteCloserWrapper(stdout)
	stderr := bytes.NewBuffer(nil)
	stderrW := ioutils.WriteCloserWrapper(stderr)

	code, err := s.lxf.Exec(req.GetContainerId(), req.GetCmd(), stdinR, stdoutW, stderrW, false, false, req.GetTimeout(), nil)
	if err != nil {
		return nil, AnnErr(log, codes.Unknown, err, "unable to exec")
	}

	log = log.WithField("exit", code)
	log.Debug("exec finished")

	return &rtApi.ExecSyncResponse{
		Stdout:   stdout.Bytes(),
		Stderr:   stderr.Bytes(),
		ExitCode: code,
	}, err
}

// Exec prepares a streaming endpoint to execute a command in the container.
func (s RuntimeServer) Exec(ctx context.Context, req *rtApi.ExecRequest) (*rtApi.ExecResponse, error) {
	log := log.WithContext(ctx).WithFields(logrus.Fields{
		"containerid": req.GetContainerId(),
		"cmd":         req.GetCmd(),
	})

	resp, err := s.stream.streamServer.GetExec(req)
	if err != nil {
		return nil, AnnErr(log, codes.Unknown, err, "unable to get exec stream")
	}

	return resp, nil
}

// Attach prepares a streaming endpoint to attach to a running container.
func (s RuntimeServer) Attach(ctx context.Context, req *rtApi.AttachRequest) (*rtApi.AttachResponse, error) {
	return nil, SilErr(log, codes.Unimplemented, ErrNotImplemented, "")
}

// PortForward prepares a streaming endpoint to forward ports from a PodSandbox.
func (s RuntimeServer) PortForward(ctx context.Context, req *rtApi.PortForwardRequest) (resp *rtApi.PortForwardResponse, err error) {
	log := log.WithContext(ctx).WithFields(logrus.Fields{
		"podid": req.GetPodSandboxId(),
		"port":  req.GetPort(),
	})

	resp, err = s.stream.streamServer.GetPortForward(req)
	if err != nil {
		return nil, AnnErr(log, codes.Unknown, err, "unable to create port forward")
	}

	return resp, nil
}

// ContainerStats returns stats of the container. If the container does not exist, the call returns an error.
func (s RuntimeServer) ContainerStats(ctx context.Context, req *rtApi.ContainerStatsRequest) (*rtApi.ContainerStatsResponse, error) {
	log := log.WithContext(ctx).WithField("containerid", req.GetContainerId())

	cntStat, err := s.lxf.GetContainer(req.GetContainerId())
	if err != nil {
		return nil, AnnErr(log, codes.Unknown, err, "unable to get container")
	}

	stats, err := toCriStats(cntStat)
	if err != nil {
		return nil, AnnErr(log, codes.Unknown, err, "unable to get stats")
	}

	return &rtApi.ContainerStatsResponse{Stats: stats}, nil
}

// ListContainerStats returns stats of all running containers.
func (s RuntimeServer) ListContainerStats(ctx context.Context, req *rtApi.ListContainerStatsRequest) (*rtApi.ListContainerStatsResponse, error) {
	log := log.WithContext(ctx).WithField("filter", req.GetFilter())

	response := &rtApi.ListContainerStatsResponse{}

	if req.GetFilter() != nil && req.GetFilter().GetId() != "" {
		log = log.WithField("containerid", req.GetFilter().GetId())

		c, err := s.lxf.GetContainer(req.Filter.Id)
		if err != nil {
			return nil, AnnErr(log, codes.Unknown, err, "unable to get container")
		}

		st, err := toCriStats(c)
		if err != nil {
			return nil, AnnErr(log, codes.Unknown, err, "unable to get stats")
		}

		response.Stats = append(response.Stats, st)

		return response, nil
	}

	cts, err := s.lxf.ListContainers()
	if err != nil {
		return nil, AnnErr(log, codes.Unknown, err, "unable to list containers")
	}

	for _, c := range cts {
		log = log.WithField("containerid", c.ID)

		st, err := toCriStats(c)
		if err != nil {
			return nil, AnnErr(log, codes.Unknown, err, "unable to get stats")
		}

		response.Stats = append(response.Stats, st)
	}

	return response, nil
}

// UpdateRuntimeConfig updates the runtime configuration based on the given request.
func (s RuntimeServer) UpdateRuntimeConfig(ctx context.Context, req *rtApi.UpdateRuntimeConfigRequest) (*rtApi.UpdateRuntimeConfigResponse, error) {
	log := log.WithContext(ctx).WithField("cidr", req.GetRuntimeConfig().GetNetworkConfig().GetPodCidr())

	err := s.network.UpdateRuntimeConfig(req.GetRuntimeConfig())
	if err != nil {
		return nil, AnnErr(log, codes.Unknown, err, "unable to update runtime config")
	}

	return &rtApi.UpdateRuntimeConfigResponse{}, nil
}

// Status returns the status of the runtime.
func (s RuntimeServer) Status(ctx context.Context, req *rtApi.StatusRequest) (*rtApi.StatusResponse, error) {
	// TODO: actually check services!
	response := &rtApi.StatusResponse{
		Status: &rtApi.RuntimeStatus{
			Conditions: []*rtApi.RuntimeCondition{
				{
					Type:   rtApi.RuntimeReady,
					Status: true,
				},
				{
					Type:   rtApi.NetworkReady,
					Status: true,
				},
			},
		},
	}

	return response, nil
}
