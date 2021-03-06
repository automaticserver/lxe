package cri // import "github.com/automaticserver/lxe/cri"

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path"

	"github.com/automaticserver/lxe/lxf"
	"github.com/automaticserver/lxe/network"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	rtApi "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
	"k8s.io/utils/exec"
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

// NewServer creates the CRI server
func NewServer(criConfig *Config) *Server {
	configPath, err := getLXDConfigPath(criConfig)
	if err != nil {
		log.WithError(err).Fatal("Unable to find lxc config")
	}

	client, err := lxf.NewClient(criConfig.LXDSocket, configPath)
	if err != nil {
		log.WithError(err).Fatal("Unable to initialize lxe facade")
	}

	log.WithField("lxdsocket", criConfig.LXDSocket).Info("Connected to LXD")

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

			writer, err = os.OpenFile(criConfig.CNIOutputFile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0660)
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
			LXDBridge:  criConfig.LXEBridgeName,
			Cidr:       criConfig.LXEBridgeDHCPRange,
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
	defer os.Remove(c.criConfig.UnixSocket)

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

// callTracing logs requests, responses and error returned by the handler. What gets logged is influenced by what error types the handler returns and the log level. This simplifies error logging in the CRI implementation.
func callTracing(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log := log.WithContext(ctx)
	method := path.Base(info.FullMethod)

	resp, err := handler(ctx, req)
	if err != nil {
		// Depending on the error type the logging is influenced
		switch e := err.(type) {
		// The AnnotatedError uses the provided logger entry to set fields of the actual logger
		case AnnotatedError:
			log.WithError(e.Err).WithFields(e.Log.Data).Error(fmt.Sprintf("%s: %s", method, e.Msg))
		// SilentErrors are useful for not implemented functions, still return the error to the caller!
		case SilentError:
			if e.Msg != "" {
				err = fmt.Errorf("%s: %w", e.Msg, e.Err)
			} else {
				err = e.Err
			}
		// CodeExitError is a special wrapping of AnnotatedError and exec.CodeExitError
		// TODO: this can be made better
		case *exec.CodeExitError:
			a, is := e.Err.(AnnotatedError)
			if is {
				log.WithError(a.Err).WithFields(a.Log.Data).Error(fmt.Sprintf("%s: %s", method, a.Msg))
			} else {
				log.Error(fmt.Sprintf("%s: %s", method, err.Error()))
			}
		// In any other case just log the error
		default:
			log.Error(fmt.Sprintf("%s: %s", method, err.Error()))
		}
	}

	log.WithError(err).WithFields(logrus.Fields{
		"req":  req,
		"resp": resp,
	}).Trace(fmt.Sprintf("grpc %s", method))

	// It seems like CRI clients don't care about the effective grpc code. The way they interact with errors is the effective error type, so not modifying the error further
	// if err != nil {
	// 	err = status.Errorf(codes.NotFound, err.Error())
	// }

	return resp, err
}
