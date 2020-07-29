package cri // import "github.com/automaticserver/lxe/cri"

import (
	"fmt"
	"net"
	"os"

	"github.com/automaticserver/lxe/lxf"
	"github.com/automaticserver/lxe/network"
	"github.com/automaticserver/lxe/shared"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	rtApi "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

// NetworkPlugin defines how the pod network should be setup.
// NetworkPluginBridge creates and manages a lxd bridge which the containers are attached to
// NetworkPluginCNI uses the kubernetes cni tools to let it attach interfaces to containers
const (
	NetworkPluginBridge = "bridge"
	NetworkPluginCNI    = "cni"
)

var (
	ErrTimeout = errors.New("timeout error")
	log        = logrus.StandardLogger()
)

// Server implements the kubernetes CRI interface specification
type Server struct {
	server    *grpc.Server
	stream    *streamService
	sock      net.Listener
	criConfig *Config
}

// NewServer creates the CRI server
func NewServer(criConfig *Config) *Server {
	configPath, err := getLXDConfigPath(criConfig)
	if err != nil {
		log.Fatalf("Unable to find lxc config: %v", err)
		os.Exit(shared.ExitCodeUnspecified)
	}

	client, err := lxf.NewClient(criConfig.LXDSocket, configPath)
	if err != nil {
		log.Fatalf("Unable to initialize lxe facade: %v", err)
		os.Exit(shared.ExitCodeUnspecified)
	}

	log.WithField("lxd-socket", criConfig.LXDSocket).Infof("Connected to LXD")

	// Ensure profile and container schema migration
	migration := lxf.NewMigrationWorkspace(client)

	err = migration.Ensure()
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
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
	case NetworkPluginBridge:
		netPlugin, err = network.InitPluginLXDBridge(client.GetServer(), network.ConfLXDBridge{
			LXDBridge:  criConfig.LXEBridgeName,
			Cidr:       criConfig.LXEBridgeDHCPRange,
			Nat:        true,
			CreateOnly: true,
		})
	default:
		err = fmt.Errorf("%w: %s", ErrUnknownNetworkPlugin, criConfig.LXENetworkPlugin)
	}

	if err != nil {
		log.Fatalf("Unable to initialize network plugin: %v", err)
		os.Exit(shared.ExitCodeUnspecified)
	}

	grpcServer := grpc.NewServer()

	// for now we bind the http on every interface
	runtimeServer, err := NewRuntimeServer(criConfig, client, netPlugin)
	if err != nil {
		log.Fatalf("Unable to start runtime server: %v", err)
		os.Exit(shared.ExitCodeUnspecified)
	}

	client.SetEventHandler(runtimeServer)

	err = setupStreamService(criConfig, runtimeServer)
	if err != nil {
		log.Fatalf("unable to create streaming server: %v", err)
		os.Exit(shared.ExitCodeUnspecified)
	}

	imageServer, err := NewImageServer(runtimeServer, client)
	if err != nil {
		log.Fatalf("Unable to start image server: %v", err)
		os.Exit(shared.ExitCodeUnspecified)
	}

	rtApi.RegisterRuntimeServiceServer(grpcServer, *runtimeServer)
	rtApi.RegisterImageServiceServer(grpcServer, *imageServer)

	return &Server{
		server:    grpcServer,
		stream:    runtimeServer.stream,
		criConfig: criConfig,
	}
}

// Serve creates the cri socket and wraps for grpc.Serve
func (c *Server) Serve() error {
	var err error

	sock := c.criConfig.UnixSocket
	log := log.WithField("socket", sock)

	if _, err = os.Stat(sock); err == nil {
		log.Debugf("Cleaning up stale socket")

		err = os.Remove(sock)
		if err != nil {
			log.Fatalf("Error cleaning up stale listening socket: %v ", err)
			os.Exit(shared.ExitCodeUnspecified)
		}
	}

	c.sock, err = net.Listen("unix", sock)
	if err != nil {
		log.Fatalf("Error listening on socket: %v ", err)
		os.Exit(shared.ExitCodeUnspecified)
	}

	defer c.sock.Close()
	defer os.Remove(c.criConfig.UnixSocket)

	log.Infof("Started %s CRI shim", Domain)

	go func() {
		err := c.stream.serve()
		if err != nil {
			panic(fmt.Errorf("error serving stream service: %w", err))
		}
	}()

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
