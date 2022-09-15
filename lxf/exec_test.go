package lxf

import (
	"strconv"
	"sync"
	"testing"
	"time"

	lxdfakes "github.com/automaticserver/lxe/fakes/lxd/client"
	"github.com/gorilla/websocket"
	lxd "github.com/lxc/lxd/client"
	"github.com/lxc/lxd/shared/api"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/tools/remotecommand"
)

func TestClient_Exec_BasicOk(t *testing.T) {
	t.Parallel()

	client, fake := testClient()
	fakeOp := &lxdfakes.FakeOperation{}

	fake.ExecContainerCalls(func(arg1 string, arg2 api.ContainerExecPost, arg3 *lxd.ContainerExecArgs) (lxd.Operation, error) {
		go sendDataDone(arg3, 0)

		return fakeOp, nil
	})
	fakeOp.WaitReturns(nil)

	fakeOp.GetReturns(api.Operation{
		Metadata: map[string]interface{}{
			"return": float64(CodeExecError),
		},
	})

	exitCode, err := client.Exec("", nil, nil, nil, nil, false, false, 0, nil)
	assert.NoError(t, err)
	assert.Equal(t, CodeExecError, exitCode)
}

func TestClient_Exec_Timeout(t *testing.T) {
	t.Parallel()

	client, fake := testClient()
	fakeOp := &lxdfakes.FakeOperation{}
	fakeSes := &session{}

	var fakeControl *websocket.Conn

	fake.ExecContainerCalls(func(arg1 string, arg2 api.ContainerExecPost, arg3 *lxd.ContainerExecArgs) (lxd.Operation, error) {
		arg3.Control = fakeSes.controlHandler
		arg3.Control(fakeControl)
		go sendDataDone(arg3, 1200*time.Millisecond)

		return fakeOp, nil
	})
	fakeOp.WaitReturns(nil)

	fakeOp.GetReturns(api.Operation{
		Metadata: map[string]interface{}{
			"return": float64(8),
		},
	})

	exitCode, err := client.Exec("", nil, nil, nil, nil, false, false, 1, nil)
	assert.Error(t, err)
	assert.Exactly(t, ErrExecTimeout, err)
	assert.Equal(t, CodeExecTimeout, exitCode)
}

// TODO: Test timeout correctly including control websocket

func TestClient_Exec_Resize(t *testing.T) {
	t.Parallel()

	client, fake := testClient()
	fakeOp := &lxdfakes.FakeOperation{}
	resize := make(chan remotecommand.TerminalSize)
	fakeSes := &session{resize: resize}

	var fakeControl *websocket.Conn

	fake.ExecContainerCalls(func(arg1 string, arg2 api.ContainerExecPost, arg3 *lxd.ContainerExecArgs) (lxd.Operation, error) {
		arg3.Control = fakeSes.controlHandler
		arg3.Control(fakeControl)
		go sendDataDone(arg3, 0)

		return fakeOp, nil
	})
	fakeOp.WaitReturns(nil)

	fakeOp.GetReturns(api.Operation{
		Metadata: map[string]interface{}{
			"return": float64(0),
		},
	})

	exitCode, err := client.Exec("", nil, nil, nil, nil, false, false, 0, fakeSes.resize)
	assert.NoError(t, err)
	assert.Equal(t, CodeExecOk, exitCode)

	// resize happens at any time
	resizeConsumed := make(chan bool)

	go func() {
		resize <- remotecommand.TerminalSize{Width: 60, Height: 40}
		resizeConsumed <- true
	}()
	<-resizeConsumed
}

func TestClient_Exec_Parallel(t *testing.T) {
	t.Parallel()

	client, fake := testClient()

	// take the first command argument as a number and use that to return to ensure the a specific command gets the
	// matching return
	fake.ExecContainerCalls(func(arg1 string, arg2 api.ContainerExecPost, arg3 *lxd.ContainerExecArgs) (lxd.Operation, error) {
		go sendDataDone(arg3, 0)

		fakeOp := &lxdfakes.FakeOperation{}
		fakeOp.WaitReturns(nil)

		returnCode, err := strconv.ParseFloat(arg2.Command[0], 64)
		assert.NoError(t, err)

		fakeOp.GetReturns(api.Operation{
			Metadata: map[string]interface{}{
				"return": returnCode,
			},
		})

		return fakeOp, nil
	})

	wg := sync.WaitGroup{}
	n := 10
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func(i int) {
			exitCode, err := client.Exec("", []string{strconv.Itoa(i)}, nil, nil, nil, false, false, 0, nil)
			assert.NoError(t, err)
			assert.Equal(t, int32(i), exitCode)
			wg.Done()
		}(i)
	}

	wg.Wait()
}

// TODO: Test resize correctly including control websocket

// func TestExecSyncInParallel(t *testing.T) {
// 	lt := newLXFTest(t)
// 	lt.createContainer(&lxf.Container{
// 		CRIObject: lxf.CRIObject{
// 			LXDObject: lxf.LXDObject{
// 				ID: "roosevelt",
// 			},
// 		},
// 		Sandbox: setUpSandbox(lt, "roosevelt"),
// 		Image:   imgBusybox,
// 	})

// 	lt.startContainer("roosevelt")

// 	wg := sync.WaitGroup{}
// 	n := 10
// 	wg.Add(n)

// 	for i := 0; i < n; i++ {
// 		go func(i int) {
// 			sti := strconv.Itoa(i)
// 			out := lt.execSync("roosevelt", []string{"echo", "foo" + sti})
// 			if string(out.StdOut) != "foo"+sti+"\n" {
// 				t.Errorf("stdout should be 'foo%v\\n' but is '%v'", sti, string(out.StdOut))
// 			}
// 			wg.Done()
// 		}(i)
// 	}
// 	wg.Wait()

// }

// func TestExecSync(t *testing.T) {
// 	lt := newLXFTest(t)
// 	lt.createContainer(&lxf.Container{
// 		CRIObject: lxf.CRIObject{
// 			LXDObject: lxf.LXDObject{
// 				ID: "roosevelt",
// 			},
// 		},
// 		Sandbox: setUpSandbox(lt, "roosevelt"),
// 		Image:   imgBusybox,
// 	})

// 	lt.startContainer("roosevelt")
// 	n := 10
// 	for i := 0; i < n; i++ {
// 		sti := strconv.Itoa(i)
// 		out := lt.execSync("roosevelt", []string{"echo", "foo" + sti})
// 		if string(out.StdOut) != "foo"+sti+"\n" {
// 			t.Errorf("stdout should be 'foo%v\\n' but is '%v'", sti, string(out.StdOut))
// 		}
// 	}
// }

// func TestExecSyncSuccess(t *testing.T) {
// 	lt := newLXFTest(t)
// 	lt.createContainer(&lxf.Container{
// 		CRIObject: lxf.CRIObject{
// 			LXDObject: lxf.LXDObject{
// 				ID: "roosevelt",
// 			},
// 		},
// 		Sandbox: setUpSandbox(lt, "roosevelt"),
// 		Image:   imgBusybox,
// 	})

// 	lt.startContainer("roosevelt")
// 	out := lt.execSync("roosevelt", []string{"echo", "foo"})
// 	if out.Code != 0 {
// 		t.Errorf("exec sync should return 0 code but is %v", out.Code)
// 	}
// }
// func TestExecSyncFailure(t *testing.T) {
// 	lt := newLXFTest(t)
// 	lt.createContainer(&lxf.Container{
// 		CRIObject: lxf.CRIObject{
// 			LXDObject: lxf.LXDObject{
// 				ID: "roosevelt",
// 			},
// 		},
// 		Sandbox: setUpSandbox(lt, "roosevelt"),
// 		Image:   imgBusybox,
// 	})

// 	lt.startContainer("roosevelt")
// 	out := lt.execSync("roosevelt", []string{"ls", "--alala"})
// 	if out.Code == 0 {
// 		t.Errorf("exec sync should not return 0 code")
// 	}
// }

// func TestNonInteractiveExec(t *testing.T) {
// 	lt := newLXFTest(t)
// 	lt.createContainer(&lxf.Container{
// 		CRIObject: lxf.CRIObject{
// 			LXDObject: lxf.LXDObject{
// 				ID: "roosevelt",
// 			},
// 		},
// 		Sandbox: setUpSandbox(lt, "roosevelt"),
// 		Image:   imgBusybox,
// 	})

// 	lt.startContainer("roosevelt")

// 	outin, stdout := io.Pipe()
// 	errin, stderr := io.Pipe()
// 	output := "fjdskfhk d saf sdfrsafrefewfesfsdfdsfsd"

// 	go func() {
// 		io.ReadAll(errin)
// 	}()

// 	go func() {
// 		stdoutbytes, err := io.ReadAll(outin)
// 		if err != nil {
// 			t.Errorf("failed to read from exec stdout, %v", err)
// 		}
// 		actual := string(stdoutbytes)
// 		expected := output + "\n"
// 		if expected != actual {
// 			t.Errorf("expected stdout to be '%v' but is '%v'", expected, actual)
// 		}
// 	}()

// 	lt.exec("roosevelt", []string{"echo", output}, nil, stdout, stderr)

// }

// An independent routine within fake sends to this channel after ExecContainer, not Wait() related, so we'll send it here
func sendDataDone(args *lxd.ContainerExecArgs, sleep time.Duration) {
	if sleep > 0 {
		time.Sleep(1200 * time.Millisecond)
	}
	args.DataDone <- true
}
