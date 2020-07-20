package cri // import "github.com/automaticserver/lxe/cri"

import (
	"net"
	"net/url"

	"github.com/sirupsen/logrus"
	utilNet "k8s.io/apimachinery/pkg/util/net"
	"k8s.io/kubernetes/pkg/kubelet/server/streaming"
)

// streamService implements streaming.Runtime.
type streamService struct {
	streaming.Runtime
	conf          streaming.Config
	runtimeServer *RuntimeServer // needed by Exec() endpoint
	streamServer  streaming.Server
}

func setupStreamService(criConfig *Config, runtime *RuntimeServer) error {
	sHost, sPort, err := net.SplitHostPort(criConfig.LXEStreamingEndpoint)
	if err != nil {
		return err
	}

	if sPort == "" {
		return &net.ParseError{Type: "Missing port", Text: criConfig.LXEStreamingEndpoint}
	}

	var bHost string
	var bPort string

	if criConfig.LXEStreamingAddress != "" {
		bHost, bPort, err = net.SplitHostPort(criConfig.LXEStreamingAddress)
		if err != nil {
			return err
		}
	}

	// If a part is empty, use it from endpoint
	if bHost == "" {
		bHost = sHost
	}

	if bPort == "" {
		bPort = sPort
	}

	// If base host is still empty, use the address of the interface of the default gateway
	if bHost == "" {
		outboundIP, err := utilNet.ChooseHostInterface()
		if err != nil {
			log.Errorf("could not find suitable host interface: %v", err)
			return err
		}

		bHost = outboundIP.String()
	}

	sService := &streamService{
		runtimeServer: runtime,
	}

	// Prepare streaming server
	sService.conf = streaming.DefaultConfig
	sService.conf.Addr = criConfig.LXEStreamingEndpoint
	sService.conf.BaseURL = &url.URL{
		Scheme: "http",
		Host:   net.JoinHostPort(bHost, bPort),
	}

	runtime.stream = sService

	sService.streamServer, err = streaming.NewServer(sService.conf, runtime.stream)
	if err != nil {
		return err
	}

	return nil
}

func (ss *streamService) serve() error {
	log.WithFields(logrus.Fields{"endpoint": ss.conf.Addr, "baseurl": ss.conf.BaseURL}).Infof("Started streaming server")

	err := ss.streamServer.Start(true)
	if err != nil {
		return err
	}

	return nil
}
