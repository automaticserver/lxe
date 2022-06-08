package cri

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/automaticserver/lxe/lxf"
	"github.com/automaticserver/lxe/network"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	rtApi "k8s.io/cri-api/pkg/apis/runtime/v1"
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
	log        = logrus.StandardLogger().WithContext(context.TODO())
)

// Server implements the kubernetes CRI interface specification
type Server struct {
	server    *grpc.Server
	stream    *streamService
	sock      net.Listener
	criConfig *Config
}

const cniFileMode = 0660

// NewServer creates the CRI server
func NewServer(criConfig *Config) *Server { // nolint: cyclop
	err := setDefaultLXDSocketPath(criConfig)
	if err != nil {
		log.WithError(err).Fatal("Unable to find lxd socket")
	}

	err = setDefaultLXDConfigPath(criConfig)
	if err != nil {
		log.WithError(err).Fatal("Unable to find lxd remote config")
	}

	log.WithField("path", criConfig.LXDRemoteConfig).Debug("Using lxd remote config")

	client, err := lxf.NewClient(criConfig.LXDSocket, criConfig.LXDRemoteConfig)
	if err != nil {
		log.WithError(err).Fatal("Unable to initialize lxe facade")
	}

	log.WithField("lxdsocket", criConfig.LXDSocket).Info("Connected to LXD")

	if criConfig.CRITest {
		log.Warn("CRITest mode enabled")

		client.SetCRITestMode()
	}

	// Ensure profile and container schema migration
	migration := lxf.NewMigrationWorkspace(client)

	err = migration.Ensure()
	if err != nil {
		log.WithError(err).Fatal("Migration failed")
	}

	// load selected plugin
	var netPlugin network.Plugin

	switch criConfig.LXENetworkPlugin {
	case NetworkPluginCNI:
		var writer io.Writer

		switch criConfig.CNIOutputTarget {
		case "stdout":
			writer = os.Stdout
		case "stderr":
			writer = os.Stderr
		case "file":
			if criConfig.CNIOutputFile == "" {
				log.Fatal("cni output file path is required when target is set to file")
			}

			writer, err = os.OpenFile(criConfig.CNIOutputFile, os.O_CREATE|os.O_RDWR|os.O_APPEND, cniFileMode)
			if err != nil {
				log.WithError(err).Fatal("could not open cni output file")
			}
		default:
			log.WithField("target", criConfig.CNIOutputTarget).Fatal("Unknown cni output target")
		}

		netPlugin, err = network.InitPluginCNI(network.ConfCNI{
			BinPath:      criConfig.CNIBinDir,
			ConfPath:     criConfig.CNIConfDir,
			OutputWriter: writer,
		})
	case NetworkPluginBridge:
		netPlugin, err = network.InitPluginLXDBridge(client.GetServer(), network.ConfLXDBridge{
			LXDBridge:  criConfig.LXDBridgeName,
			Cidr:       criConfig.LXDBridgeDHCPRange,
			Nat:        true,
			CreateOnly: true,
		})
	default:
		err = fmt.Errorf("%w: %s", ErrUnknownNetworkPlugin, criConfig.LXENetworkPlugin)
	}

	if err != nil {
		log.WithError(err).Fatal("Unable to initialize network plugin")
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(callTracing))

	// for now we bind the http on every interface
	runtimeServer, err := NewRuntimeServer(criConfig, client, netPlugin)
	if err != nil {
		log.WithError(err).Fatal("Unable to start runtime server")
	}

	client.SetEventHandler(runtimeServer)

	err = setupStreamService(criConfig, runtimeServer)
	if err != nil {
		log.WithError(err).Fatal("unable to create streaming server")
	}

	imageServer, err := NewImageServer(runtimeServer, client)
	if err != nil {
		log.WithError(err).Fatal("Unable to start image server")
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
		log.Debugf("cleaning up stale socket")

		err = os.Remove(sock)
		if err != nil {
			log.WithError(err).Fatal("error cleaning up stale listening socket")
		}
	}

	c.sock, err = net.Listen("unix", sock)
	if err != nil {
		log.WithError(err).Fatal("error listening on socket")
	}

	defer c.sock.Close()
	defer os.Remove(sock)

	log.Infof("started %s CRI shim", Domain)

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
