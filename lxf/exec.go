package lxf

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	lxd "github.com/lxc/lxd/client"
	lxdApi "github.com/lxc/lxd/shared/api"
	"github.com/lxc/lxd/shared/logger"
	"golang.org/x/sys/unix"
	"k8s.io/client-go/tools/remotecommand"
)

const (
	WindowHeightDefault = 24
	WindowWidthDefault  = 80
)

var (
	ErrExecTimeout     = errors.New("timeout reached")
	ErrNoControlSocket = errors.New("no control socket found")

	cancelSignal = unix.SIGTERM

	CodeExecOk      int32 = 0
	CodeExecError   int32 = 128
	CodeExecTimeout int32 = CodeExecError + int32(cancelSignal) // 128+15=143
)

// Exec will start a command on the server and attach the provided streams. It will block till the command terminated
// AND all data was written to stdout/stdin. The caller is responsible to provide a sink which doesn't block.
func (l *client) Exec(cid string, cmd []string, stdin io.ReadCloser, stdout, stderr io.WriteCloser, interactive, tty bool, timeout int64, resize <-chan remotecommand.TerminalSize) (int32, error) {
	ses := &session{resize: resize}

	req := lxdApi.ContainerExecPost{
		Command:      cmd,
		WaitForWS:    true,
		Interactive:  interactive,
		Environment:  map[string]string{"TERM": "xterm"},
		Width:        WindowWidthDefault,
		Height:       WindowHeightDefault,
		RecordOutput: false,
	}
	args := &lxd.ContainerExecArgs{
		Stdin:    stdin,
		Stdout:   stdout,
		Stderr:   stderr,
		Control:  ses.controlHandler,
		DataDone: make(chan bool),
	}

	op, err := l.server.ExecContainer(cid, req, args)
	if err != nil {
		return CodeExecError, err
	}

	var deadline <-chan time.Time
	if timeout > 0 {
		deadline = time.After(time.Duration(timeout) * time.Second)
	}

	select {
	// Exit early if timeout is reached
	case <-deadline:
		err := ses.sendCancel()
		if err != nil {
			logger.Errorf("session control failed: %v", err)
		}

		return CodeExecTimeout, ErrExecTimeout

	// Wait for any remaining I/O to be flushed
	case <-args.DataDone:
	}

	// Wait for the operation to complete so we can get the return code
	err = op.Wait()
	if err != nil {
		return CodeExecError, err
	}

	opAPI := op.Get()

	exitCode, ok := opAPI.Metadata["return"].(float64)
	if !ok {
		return CodeExecError, fmt.Errorf("can't parse error code: %#v", opAPI.Metadata["return"])
	}

	return int32(exitCode), nil
}

type session struct {
	resize  <-chan remotecommand.TerminalSize
	control *websocket.Conn
}

func (s *session) controlHandler(control *websocket.Conn) {
	s.control = control

	go s.listenResize()
}

func (s *session) listenResize() {
	for r := range s.resize {
		err := s.sendResize(r)
		if err != nil {
			logger.Errorf("session control failed: %v", err)
		}
	}
}

func (s *session) sendResize(r remotecommand.TerminalSize) error {
	width := strconv.FormatUint(uint64(r.Width), 10)
	height := strconv.FormatUint(uint64(r.Height), 10)

	logger.Debugf("session control window size is now: %vx%v", width, height)

	if s.control == nil {
		return ErrNoControlSocket
	}

	w, err := s.control.NextWriter(websocket.TextMessage)
	if err != nil {
		return err
	}

	msg := lxdApi.ContainerExecControl{}
	msg.Command = "window-resize"
	msg.Args = make(map[string]string)
	msg.Args["width"] = width
	msg.Args["height"] = height

	buf, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	_, err = w.Write(buf)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return nil
}

func (s *session) sendCancel() error {
	logger.Debugf("timeout reached, force closing of connection")

	if s.control == nil {
		return ErrNoControlSocket
	}

	// TODO: In noninteractive mode this doesn't stop the command from executing, whatever signal it is set to
	sig := cancelSignal

	logger.Debugf("forwarding signal: %s", sig)

	w, err := s.control.NextWriter(websocket.TextMessage)
	if err != nil {
		return err
	}
	defer w.Close()

	msg := lxdApi.ContainerExecControl{}
	msg.Command = "signal"
	msg.Signal = int(sig)

	buf, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	_, err = w.Write(buf)
	if err != nil {
		return err
	}

	closeMsg := websocket.FormatCloseMessage(websocket.CloseGoingAway, "timeout reached")
	err = s.control.WriteControl(websocket.CloseMessage, closeMsg, time.Now().Add(1*time.Second)) // nolint: gomnd

	return err
}
