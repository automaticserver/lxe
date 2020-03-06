package cri

import (
	"net"
	"os"

	"emperror.dev/errors"
	"github.com/automaticserver/lxe/lxf"
	"github.com/automaticserver/lxe/network"
	"github.com/automaticserver/lxe/shared"
	"github.com/lxc/lxd/shared/logger"
	"google.golang.org/grpc"
	runtimeapi "k8s.io/kubernetes/pkg/kubelet/apis/cri/runtime/v1alpha2"
)

// NetworkPlugin defines how the pod network should be setup.
// NetworkPluginDefault creates and manages a lxd bridge which the containers are attached to
// NetworkPluginCNI uses the kubernetes cni tools to let it attach interfaces to containers
const (
	NetworkPluginDefault = ""
	NetworkPluginCNI     = "cni"
)

// Server implements the kubernetes CRI interface specification
type Server struct {
	server    *grpc.Server
	sock      net.Listener
	criConfig *Config
}

// NewServer creates the CRI server
func NewServer(criConfig *Config) *Server {
	configPath, err := getLXDConfigPath(criConfig)
	if err != nil {
		logger.Critf("Unable to find lxc config: %v", err)
		os.Exit(shared.ExitCodeUnspecified)
	}

	client, err := lxf.NewClient(criConfig.LXDSocket, configPath)
	if err != nil {
		logger.Critf("Unable to initialize lxe facade: %v", err)
		os.Exit(shared.ExitCodeUnspecified)
	}

	logger.Infof("Connected to LXD via %q", criConfig.LXDSocket)

	// Ensure profile and container schema migration
	migration := lxf.NewMigrationWorkspace(client)

	err = migration.Ensure()
	if err != nil {
		logger.Critf("Migration failed: %v", err)
		os.Exit(shared.ExitCodeSchemaMigrationFailure)
	}

	// load selected plugin
	var netPlugin network.Plugin

	switch criConfig.LXENetworkPlugin {
	case NetworkPluginCNI:
		netPlugin, err = network.InitPluginCNI(network.ConfCNI{
			BinPath:  criConfig.CNIBinDir,
			ConfPath: criConfig.CNIConfDir,
		})
	case NetworkPluginDefault:
		netPlugin, err = network.InitPluginLXDBridge(client.GetServer(), network.ConfLXDBridge{
			LXDBridge:  criConfig.LXEBridgeName,
			Cidr:       criConfig.LXEBridgeDHCPRange,
			Nat:        true,
			CreateOnly: true,
		})
	default:
		err = errors.Errorf("Unknown network plugin %s", criConfig.LXENetworkPlugin)
	}

	if err != nil {
		logger.Critf("Unable to initialize network plugin: %v", err)
		os.Exit(shared.ExitCodeUnspecified)
	}

	grpcServer := grpc.NewServer()

	// for now we bind the http on every interface
	runtimeServer, err := NewRuntimeServer(criConfig, client, netPlugin)
	if err != nil {
		logger.Critf("Unable to start runtime server: %v", err)
		os.Exit(shared.ExitCodeUnspecified)
	}

	client.SetEventHandler(runtimeServer)

	imageServer, err := NewImageServer(runtimeServer, client)
	if err != nil {
		logger.Critf("Unable to start image server: %v", err)
		os.Exit(shared.ExitCodeUnspecified)
	}

	runtimeapi.RegisterRuntimeServiceServer(grpcServer, *runtimeServer)
	runtimeapi.RegisterImageServiceServer(grpcServer, *imageServer)

	return &Server{
		server:    grpcServer,
		criConfig: criConfig,
	}
}

// Serve creates the cri socket and wraps for grpc.Serve
func (c *Server) Serve() error {
	var err error

	sock := c.criConfig.UnixSocket

	if _, err = os.Stat(sock); err == nil {
		logger.Debugf("Cleaning up stale socket")

		err = os.Remove(sock)
		if err != nil {
			logger.Critf("Error cleaning up stale listening socket %q: %v ", sock, err)
			os.Exit(shared.ExitCodeUnspecified)
		}
	}

	c.sock, err = net.Listen("unix", sock)
	if err != nil {
		logger.Critf("Error listening on socket %q: %v ", sock, err)
		os.Exit(shared.ExitCodeUnspecified)
	}

	logger.Infof("Started %s/%s CRI shim on UNIX socket %q", Domain, Version, sock)

	defer c.sock.Close()
	defer os.Remove(c.criConfig.UnixSocket)

	return c.server.Serve(c.sock)
}

// Stop stops the cri socket
func (c *Server) Stop() error {
	c.server.Stop()

	err := c.sock.Close()
	if err != nil {
		return err
	}

	return os.Remove(c.criConfig.UnixSocket)
}
