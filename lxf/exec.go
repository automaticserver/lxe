package lxf

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

	lxd "github.com/lxc/lxd/client"
	lxdApi "github.com/lxc/lxd/shared/api"
)

// ExecResponse returns the stdout and err from an exec call
type ExecResponse struct {
	StdOut []byte
	StdErr []byte
	Code   int
}

// ExecSync runs a command on a container and blocks till it's finished
func (l *Client) ExecSync(cid string, cmd []string) (*ExecResponse, error) {
	tempStderr := NewWriteCloserBuffer()
	tempStdout := NewWriteCloserBuffer()

	dataDone := make(chan bool)

	op, err := l.opwait.ExecContainer(cid, lxdApi.ContainerExecPost{
		Command:     cmd,
		Interactive: false,
		Width:       0,
		Height:      0,
		WaitForWS:   true,
	}, &lxd.ContainerExecArgs{
		Stderr:   tempStderr,
		Stdout:   tempStdout,
		Stdin:    ioutil.NopCloser(bytes.NewReader(nil)),
		DataDone: dataDone,
	})
	if err != nil {
		return nil, err
	}

	ret, has := op.Get().Metadata["return"].(float64)
	if !has {
		return nil, fmt.Errorf("exec sync could not read the return code")
	}

	<-dataDone // wait till all data is written (stdout, stderr)

	return &ExecResponse{
		StdErr: tempStderr.Bytes(),
		StdOut: tempStdout.Bytes(),
		Code:   int(ret),
	}, nil
}

// Exec will start a command on the server and attach the provided streams. It will block till the command terminated
// AND all data was written to stdout/stdin. The caller is responsible to provide a sink which doesn't block.
func (l *Client) Exec(cid string, cmd []string, stdin io.Reader, stdout, stderr io.WriteCloser) (int, error) {
	// we get io.Reader interface from the kubelet but lxd wants ReadCloser interface
	var stdinCloser io.ReadCloser
	// kubelet might give us stdin==nil but lxd expects something there otherwise it will segfault
	if stdin == nil {
		stdinCloser = ioutil.NopCloser(bytes.NewBufferString(""))
	} else {
		stdinCloser = ioutil.NopCloser(stdin)
	}

	environment := map[string]string{
		"TERM": "xterm",
	}

	dataDone := make(chan bool)

	// TODO: Is no op.Wait() intentional?
	_, err := l.server.ExecContainer(cid,
		lxdApi.ContainerExecPost{
			Command:      cmd,
			WaitForWS:    true,
			Interactive:  (stdin != nil), // if there is no stdin, exec won't be interactive
			Environment:  environment,
			Width:        80,
			Height:       24,
			RecordOutput: false,
		}, &lxd.ContainerExecArgs{
			Stdin:    stdinCloser,
			Stdout:   stdout,
			Stderr:   stderr,
			Control:  nil,
			DataDone: dataDone,
		})
	if err != nil {
		return 1, err
	}

	<-dataDone

	// we close as soon as connections are terminated and all data got sent. it seems they won't be closed automatically
	// but i'm not sure if i miss something
	if stdout != nil {
		err = stdout.Close()
		if err != nil {
			return 0, err
		}
	}

	if stderr != nil {
		err = stderr.Close()
		if err != nil {
			return 0, err
		}
	}

	// do not wait for operation, it will wait till command finished executing
	return 0, nil
}
