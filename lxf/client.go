package lxf

import (
	"net/http"
	"time"

	"github.com/automaticserver/lxe/lxf/lxo"
	lxd "github.com/lxc/lxd/client"
	"github.com/lxc/lxd/lxc/config"
	"github.com/lxc/lxd/shared/api"
)

// Client is a facade to thin the interface to map the cri logic to lxd.
type Client struct {
	server lxd.ContainerServer
	config *config.Config
	opwait *lxo.LXO
}

// NewClient will set up a connection and return the client
func NewClient(socket string, configPath string) (*Client, error) {
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
			Timeout: 10 * time.Second,
		},
	}

	server, err := lxd.ConnectLXDUnix(socket, &args)
	if err != nil {
		return nil, err
	}

	config, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	client := &Client{
		server: server,
		config: config,
		opwait: lxo.NewClient(server),
	}

	// register LXD eventhandler
	listener, err := server.GetEvents()
	if err != nil {
		return nil, err
	}

	_, err = listener.AddHandler([]string{"lifecycle"}, client.lifecycleEventHandler)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (l *Client) GetRuntimeInfo() (*api.Server, error) {
	server, _, err := l.server.GetServer()
	if err != nil {
		return nil, err
	}

	return server, nil
}
