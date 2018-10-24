package lxf

import (
	"net/http"
	"time"

	lxd "github.com/lxc/lxd/client"
	"github.com/lxc/lxd/lxc/config"
)

// LXF is a facade to thin the interface needed for our usecase and map
// the cri logic to lxd.
type LXF struct {
	sleeperHash    string
	server         lxd.ContainerServer
	config         *config.Config
	cntMonitorChan chan ContainerMonitorChan
}

// ContainerMonitorChan holds data for a go routine which periodically launches
// a function for a container depending on the given task and interval
type ContainerMonitorChan struct {
	container   *Container
	task        string
	lastCheck   time.Time
	intervalSec time.Duration
	once        bool
}

// New will set up a connection and return the LXF facade
func New(socket string, configPath string) (*LXF, error) {
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
	// is neither official supported nor encouraged and LXE is not yet opensourced filing a bugreport against
	// LXD is rather pointless.

	args := lxd.ConnectionArgs{
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
	server, err := lxd.ConnectLXDUnix(socket, &args)
	if err != nil {
		return nil, err
	}

	lxdConfig, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	// a, _, err := server.GetImageAlias("sleeper")
	// if err != nil {
	// 	log.Printf("faield to get sleeper image hash, provide a sleeper image with name sleeper")
	// }

	lxf := &LXF{
		sleeperHash:    "",
		server:         server,
		config:         lxdConfig,
		cntMonitorChan: make(chan ContainerMonitorChan),
	}

	go lxf.containerMonitor(lxf.cntMonitorChan)

	// register LXD eventhandler
	listener, err := server.GetEvents()
	if err != nil {
		return nil, err
	}
	_, err = listener.AddHandler([]string{"lifecycle"}, lxf.lifecycleEventHandler)
	if err != nil {
		return nil, err
	}

	return lxf, nil
}

// IsErrorNotFound returns true if the error is a ErrorNotFound error
func IsErrorNotFound(err error) bool {
	return err.Error() == ErrorNotFound
}
