package lxf

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"path"
	"time"

	"github.com/automaticserver/lxe/lxf/lxo"
	lxd "github.com/lxc/lxd/client"
	"github.com/lxc/lxd/lxc/config"
	"github.com/lxc/lxd/shared/logger"
	"gopkg.in/fsnotify.v1"
	"k8s.io/client-go/tools/remotecommand"
)

var (
	ErrMissingETag = errors.New("missing ETag")
	ErrConvert     = errors.New("convert error")
	ErrParse       = errors.New("parse error")
	ErrUsage       = errors.New("usage error")
)

// Client is a facade to thin the interface to map the cri logic to lxd.
type Client interface {
	// GetServer returns the lxd ContainerServer. TODO: since it created it and others want to access lxd too (lxdbridge
	// network plugin) either return it here, or extract creation of the connection outside and pass server into
	// NewClient(), but that makes the initialisation NewClient() pretty unnecessary
	GetServer() lxd.ContainerServer
	// GetRuntimeInfo returns informations about the runtime
	GetRuntimeInfo() (*RuntimeInfo, error)
	// SetEventHandler for container's starting and stopping events
	SetEventHandler(eh EventHandler)

	// PullImage copies the given image from the remote server
	PullImage(name string) (string, error)
	// RemoveImage will remove the given image
	RemoveImage(name string) error
	// ListImages will list all local images from the lxd server
	ListImages(filter string) ([]Image, error)
	// GetImage will fetch information about the already downloaded image identified by name
	GetImage(name string) (*Image, error)
	// GetFSPoolUsage returns a list of usage information about the used storage pools
	GetFSPoolUsage() ([]FSPoolUsage, error)

	// NewSandbox creates a local representation of a sandbox
	NewSandbox() *Sandbox
	// GetSandbox will find a sandbox by id and return it.
	GetSandbox(id string) (*Sandbox, error)
	// ListSandboxes will return a list with all the available sandboxes
	ListSandboxes() ([]*Sandbox, error)

	// NewContainer creates a local representation of a container
	NewContainer(sandboxID string, additionalProfiles ...string) *Container
	// GetContainer returns the container identified by id
	GetContainer(id string) (*Container, error)
	// ListContainers returns a list of all available containers
	ListContainers() ([]*Container, error)

	// Exec will start a command on the server and attach the provided streams. It will block till the command terminated
	// AND all data was written to stdout/stdin. The caller is responsible to provide a sink which doesn't block.
	Exec(cid string, cmd []string, stdin io.ReadCloser, stdout, stderr io.WriteCloser, interactive, tty bool, timeout int64, resize <-chan remotecommand.TerminalSize) (int32, error)
}

var (
	lxdHTTPTimeout = 10 * time.Second
)

type client struct {
	server       lxd.ContainerServer
	config       *config.Config
	opwait       *lxo.LXO
	eventHandler EventHandler
	socket       string
}

// NewClient will set up a connection and return the client
func NewClient(socket string, configPath string) (Client, error) {
	config, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	cl := &client{
		config: config,
		socket: socket,
	}

	err = cl.connect()
	if err != nil {
		return nil, err
	}

	go cl.detectNeedReconnect()

	return cl, nil
}

// GetServer returns the lxd ContainerServer. TODO: since it created it and others want to access lxd too (lxdbridge
// network plugin) either return it here, or extract creation of the connection outside and pass server into
// NewClient(), but that makes the initialisation NewClient() pretty unnecessary
func (l *client) GetServer() lxd.ContainerServer {
	return l.server
}

// SetEventHandler for container's starting and stopping events
func (l *client) SetEventHandler(eh EventHandler) {
	l.eventHandler = eh
}

type RuntimeInfo struct {
	// API version of the container runtime. The string must be semver-compatible.
	Version string
}

// GetRuntimeInfo returns informations about the runtime
func (l *client) GetRuntimeInfo() (*RuntimeInfo, error) {
	server, _, err := l.server.GetServer()
	if err != nil {
		return nil, err
	}

	return &RuntimeInfo{
		// api version is only X.X, so need to add .0 for semver requirement
		Version: fmt.Sprintf("%s.0", server.APIVersion),
	}, nil
}

func (l *client) connect() error {
	args := lxd.ConnectionArgs{
		HTTPClient: &http.Client{
			// this was a byproduct of a bughunt, but i figured using TCP connections with TLS instead of unix sockets
			// might by useful some time in the future so i will leave this here commented out
			// to setup this up in LXD see: https://help.ubuntu.com/lts/serverguide/lxd.html.en#lxd-server-config

			// tlsServerCert, err := ioutil.ReadFile("/root/.config/lxc/servercerts/r1.crt")
			// if err != nil {
			// 	panic(err)
			// }
			// tlsClientCert, err := ioutil.ReadFile("/root/.config/lxc/client.crt")
			// if err != nil {
			// 	panic(err)
			// }
			// tlsClientKey, err := ioutil.ReadFile("/root/.config/lxc/client.key")
			// if err != nil {
			// 	panic(err)
			// }

			// server, err := lxd.ConnectLXD("https://127.0.0.1:8443", &lxd.ConnectionArgs{
			// 	TLSServerCert:      string(tlsServerCert),
			// 	TLSClientCert:      string(tlsClientCert),
			// 	TLSClientKey:       string(tlsClientKey),
			// 	InsecureSkipVerify: true,
			// })
			// -------------------------------------------
			// it was discovered when using a container with "hostnetwork: true" LXE
			// would leak filehandles indefinitely until the process hits the system limit and
			// LXE would stop working since no new connections could be opened.
			// This happens for unix sockets as well as for tcp/tls connections.

			// this issue could be observed by
			// a) lsof -n -p $(pidof lxe)     yielding more and more connections
			// b) pkill -SIGABRT lxe          seeing many many goroutines like this:
			// net/http.(*persistConn).readLoop(0xc000273b00)
			// 	/home/dj/src/go/src/net/http/transport.go:1761 +0x6b9
			// and
			// net/http.(*persistConn).writeLoop(0xc0002c25a0)
			// 	/home/dj/src/go/src/net/http/transport.go:1885 +0x113
			// without any stacktrace/callstack.

			// online search will lead to some golang issues at github which most of are marked as fixed
			// as well as the solution to "defer resp.Body.Close()" which is done by the LXD client api.
			// other measures like "_, err = io.Copy(ioutil.Discard, resp.Body)" were tried as well.
			// (see https://hackernoon.com/avoiding-memory-leak-in-golang-api-1843ef45fca8 e.g.)

			// the chain to track this is:
			//    call lxd.ConnectLXDUnix (lxf/client.go)
			// -> unixHttpClient (lxd/client/connection.go)
			//    >> here we force the httpClient to have a Timeout <<
			// -> unixHttpClient (lxd/client/util.go) setups a DialUnix inside of a Transport inside of the HttpClient

			// the HttpClient tries to reuse already opened connections (this is done by golangs core library
			// and is rather transparent for the caller) which does not seem to happen in this special case.
			// this commit forces a timeout on the httpClient used by the LXE (via the LXD client API) to talk to LXD.

			// this does not fix the real problem, which can be either in LXE, LXD client API or the LXD server
			// and still needs more investigation.
			// since sharing the networknamespace between host and container via "lxc.raw = lxc.net.0.type=none"
			// is neither officially supported nor encouraged, filing a bugreport against LXD is rather pointless.
			Timeout: lxdHTTPTimeout,
		},
	}

	server, err := lxd.ConnectLXDUnix(l.socket, &args)
	if err != nil {
		return err
	}

	// register LXD eventhandler
	listener, err := server.GetEvents()
	if err != nil {
		return err
	}

	_, err = listener.AddHandler([]string{"lifecycle"}, l.lifecycleEventHandler)
	if err != nil {
		return err
	}

	l.server = server
	l.opwait = lxo.NewClient(server)

	return nil
}

// detect if server needs to be connected again to. Seems to be needed if we get a lxd.RemoteOperation (e.g. in CopyImage), the op.Wait() never succeeds unless we have connected to the lxd socket again. All other lxd.Operations seem to work fine and wouldn't be needed for them.
func (l *client) detectNeedReconnect() { // nolint: gocognit
	// currently I know no way to find out when a socket is gone as all is encapsulated in lxd.ContainerServer. We can set an fsnotify to the socket file so we get an event when it was created. If we got such event, we try to connect again until it is successful.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.Crit(err.Error())
	}
	defer watcher.Close()

	err = watcher.Add(path.Dir(l.socket))
	if err != nil {
		logger.Crit(err.Error())
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			if event.Op&fsnotify.Create == fsnotify.Create && event.Name == l.socket {
				logger.Infof("socket %s got created, trying to reconnect", l.socket)

				go func() {
					for {
						err := l.connect()
						if err != nil {
							// print error and try again
							logger.Errorf("tried reconnecting to lxd socket: %v", err)
						} else {
							logger.Info("reconnected to lxd socket")

							return
						}
					}
				}()
			}

			if event.Op&fsnotify.Remove == fsnotify.Remove && event.Name == l.socket {
				logger.Warnf("socket %s got deleted, will try to reconnect once it's created again", l.socket)
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}

			logger.Crit("error: %v", err)
		}
	}
}
