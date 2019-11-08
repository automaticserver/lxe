package lxf

import (
	"testing"

	"github.com/automaticserver/lxe/lxf/lxdfakes"
	lxd "github.com/lxc/lxd/client"
	lxdApi "github.com/lxc/lxd/shared/api"
	"github.com/stretchr/testify/assert"
)

func TestClient_Exec_BasicOk(t *testing.T) {
	t.Parallel()

	client, fake := testClient()
	fakeOp := &lxdfakes.FakeOperation{}
	// fakeControl := &websocket.Conn{}
	// fakeSes := &session{}
	// fakeDataDone := make(chan bool)
	fake.ExecContainerCalls(func(arg1 string, arg2 lxdApi.ContainerExecPost, arg3 *lxd.ContainerExecArgs) (lxd.Operation, error) {
		// 	arg3.Control = fakeSes.controlHandler
		// 	arg3.Control(fakeControl)

		// An independent routine within fake sends to this channel after ExecContainer, not Wait() related, so we'll send it here
		go func() {
			arg3.DataDone <- true
		}()

		return fakeOp, nil
	})
	fakeOp.WaitReturns(nil)

	fakeOp.GetReturns(lxdApi.Operation{
		Metadata: map[string]interface{}{
			"return": float64(8),
		},
	})

	exitCode, err := client.Exec("", nil, nil, nil, nil, false, false, 0, nil)
	assert.NoError(t, err)
	assert.Equal(t, int32(8), exitCode)
}

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
// 		ioutil.ReadAll(errin)
// 	}()

// 	go func() {
// 		stdoutbytes, err := ioutil.ReadAll(outin)
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
