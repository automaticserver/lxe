package cri // import "github.com/automaticserver/lxe/cri"

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/url"
	"os/exec"
	"strings"

	"github.com/docker/docker/pkg/pools"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	utilNet "k8s.io/apimachinery/pkg/util/net"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/kubernetes/pkg/kubelet/server/streaming"
	utilExec "k8s.io/utils/exec"
)

// streamService implements streaming.Runtime.
type streamService struct {
	streaming.Runtime
	conf          streaming.Config
	runtimeServer *RuntimeServer // needed by Exec() endpoint
	streamServer  streaming.Server
}

func setupStreamService(criConfig *Config, runtime *RuntimeServer) error {
	sHost, sPort, err := net.SplitHostPort(criConfig.LXEStreamingBindAddr)
	if err != nil {
		return err
	}

	var bHost string
	var bPort string

	if criConfig.LXEStreamingBaseURL != "" {
		bHost, bPort, err = net.SplitHostPort(criConfig.LXEStreamingBaseURL)
		if err != nil {
			var aerr *net.AddrError
			if errors.As(err, &aerr) && aerr.Err == "missing port in address" {
				// we allow port to be missing here
			} else {
				return err
			}
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
			return fmt.Errorf("could not find suitable host interface: %w", err)
		}

		bHost = outboundIP.String()
	}

	sService := &streamService{
		runtimeServer: runtime,
	}

	// Prepare streaming server
	sService.conf = streaming.DefaultConfig
	sService.conf.Addr = criConfig.LXEStreamingBindAddr
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
	log.WithFields(logrus.Fields{"endpoint": ss.conf.Addr, "baseurl": ss.conf.BaseURL}).Info("started streaming server")

	err := ss.streamServer.Start(true)
	if err != nil {
		return err
	}

	return nil
}

func (ss streamService) Exec(containerID string, cmd []string, stdinR io.Reader, stdout, stderr io.WriteCloser, tty bool, resize <-chan remotecommand.TerminalSize) error {
	log := log.WithField("container", containerID).WithField("cmd", cmd)

	var stdin io.ReadCloser
	if stdinR == nil {
		stdin = ioutil.NopCloser(bytes.NewReader(nil))
	} else {
		stdin = ioutil.NopCloser(stdinR)
	}

	interactive := (stdinR != nil)

	code, err := ss.runtimeServer.lxf.Exec(containerID, cmd, stdin, stdout, stderr, interactive, tty, 0, resize)

	log.Debugf("received exit code %v", code)
	log = log.WithField("exit", code)

	if err != nil || code != 0 {
		return &utilExec.CodeExitError{
			Err:  AnnErr(log, err, "error executing command"),
			Code: int(code),
		}
	}

	return nil
}

func (ss streamService) PortForward(podSandboxID string, port int32, stream io.ReadWriteCloser) error {
	log := log.WithField("podsandbox", podSandboxID).WithField("port", port)

	sb, err := ss.runtimeServer.lxf.GetSandbox(podSandboxID)
	if err != nil {
		return AnnErr(log, err, "unable to find pod")
	}

	podIP := ss.runtimeServer.getInetAddress(context.TODO(), sb)

	_, err = exec.LookPath("socat")
	if err != nil {
		return AnnErr(log, err, "unable to do port forwarding")
	}

	args := []string{"-", fmt.Sprintf("TCP4:%s:%d,keepalive", podIP, port)}

	commandString := fmt.Sprintf("socat %s", strings.Join(args, " "))
	log.WithField("cmd", commandString).Debug("executing port forwarding command")

	command := exec.Command("socat", args...)
	command.Stdout = stream

	stderr := new(bytes.Buffer)
	command.Stderr = stderr

	// If we use Stdin, command.Run() won't return until the goroutine that's copying from stream finishes. Unfortunately,
	// if you have a client like telnet connected via port forwarding, as long as the user's telnet client is connected to
	// the user's local listener that port forwarding sets up, the telnet session never exits. This means that even if
	// socat has finished running, command.Run() won't ever return (because the client still has the connection and stream
	// open). The work around is to use StdinPipe(), as Wait() (called by Run()) closes the pipe when the command (socat)
	// exits.
	inPipe, err := command.StdinPipe()
	if err != nil {
		return AnnErr(log, err, "unable to do port forwarding")
	}

	go func() {
		_, err = pools.Copy(inPipe, stream)
		if err != nil {
			log.WithError(err).Error("pipe copy errored")
		}

		err = inPipe.Close()
		if err != nil {
			log.WithError(err).Error("pipe close errored")
		}
	}()

	err = command.Run()
	if err != nil {
		return AnnErr(log, err, stderr.String())
	}

	return nil
}
