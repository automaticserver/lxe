package cri

import (
	"net"
	"os"

	"github.com/automaticserver/lxe/lxf"
	"github.com/automaticserver/lxe/shared"
	"github.com/lxc/lxd/shared/logger"
	"google.golang.org/grpc"
	runtimeapi "k8s.io/kubernetes/pkg/kubelet/apis/cri/runtime/v1alpha2"
)

// NetworkPlugin defines how the pod network should be setup.
// NetworkPluginDefault creates and manages a `LXEBridge` bridge which the containers are attached to
// NetworkPluginCNI uses the kubernetes cni tools to let it attach interfaces to containers
const (
	NetworkPluginDefault = ""
	NetworkPluginCNI     = "cni"
	LXEBridge            = "lxebr0"
)

// Server implements the kubernetes CRI interface specification
type Server struct {
	server    *grpc.Server
	sock      net.Listener
	criConfig *Config
}

// NewServer creates the CRI server
func NewServer(criConfig *Config) *Server {
	grpcServer := grpc.NewServer()

	configPath, err := getLXDConfigPath(criConfig)
	if err != nil {
		logger.Critf("Unable to find lxc config: %v", err)
		os.Exit(shared.ExitCodeUnspecified)
	}

	lxf, err := lxf.NewClient(criConfig.LXDSocket, configPath)
	if err != nil {
		logger.Critf("Unable to initialize lxe facade: %v", err)
		os.Exit(shared.ExitCodeUnspecified)
	}

	logger.Infof("Connected to LXD via %q", criConfig.LXDSocket)

	// Ensure profile and container schema migration
	migration := lxf.Migration()

	err = migration.Ensure()
	if err != nil {
		logger.Critf("Migration failed: %v", err)
		os.Exit(shared.ExitCodeSchemaMigrationFailure)
	}

	if criConfig.LXENetworkPlugin == NetworkPluginDefault {
		// Initialize lxd bridge for lxe is created with new generated cidr if missing
		err = lxf.EnsureBridge(LXEBridge, criConfig.LXEBridgeDHCPRange, true, true)
		if err != nil {
			logger.Critf("Unable to setup bridge %v: %v", LXEBridge, err)
			os.Exit(shared.ExitCodeUnspecified)
		}
	}

	// for now we bind the http on every interface
	runtimeServer, err := NewRuntimeServer(criConfig, lxf)
	if err != nil {
		logger.Critf("Unable to start runtime server: %v", err)
		os.Exit(shared.ExitCodeUnspecified)
	}

	imageServer, err := NewImageServer(runtimeServer, lxf)
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
