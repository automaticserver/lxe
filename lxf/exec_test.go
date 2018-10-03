package lxf_test

import (
	"io"
	"io/ioutil"
	"strconv"
	"sync"
	"testing"

	"github.com/lxc/lxe/lxf"
)

func TestExecSyncInParallel(t *testing.T) {
	lt := newLXFTest(t)
	lt.createContainer(&lxf.Container{
		Name:    "roosevelt",
		Sandbox: setUpSandbox(lt, "roosevelt"),
		Image:   imgBusybox,
	})

	lt.startContainer("roosevelt")

	wg := sync.WaitGroup{}
	n := 10
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func(i int) {
			sti := strconv.Itoa(i)
			out := lt.execSync("roosevelt", []string{"echo", "foo" + sti})
			if string(out.StdOut) != "foo"+sti+"\n" {
				t.Errorf("stdout should be 'foo%v\\n' but is '%v'", sti, string(out.StdOut))
			}
			wg.Done()
		}(i)
	}
	wg.Wait()

}

func TestExecSync(t *testing.T) {
	lt := newLXFTest(t)
	lt.createContainer(&lxf.Container{
		Name:    "roosevelt",
		Sandbox: setUpSandbox(lt, "roosevelt"),
		Image:   imgBusybox,
	})

	lt.startContainer("roosevelt")
	n := 10
	for i := 0; i < n; i++ {
		sti := strconv.Itoa(i)
		out := lt.execSync("roosevelt", []string{"echo", "foo" + sti})
		if string(out.StdOut) != "foo"+sti+"\n" {
			t.Errorf("stdout should be 'foo%v\\n' but is '%v'", sti, string(out.StdOut))
		}
	}
}

func TestExecSyncSuccess(t *testing.T) {
	lt := newLXFTest(t)
	lt.createContainer(&lxf.Container{
		Name:    "roosevelt",
		Sandbox: setUpSandbox(lt, "roosevelt"),
		Image:   imgBusybox,
	})

	lt.startContainer("roosevelt")
	out := lt.execSync("roosevelt", []string{"echo", "foo"})
	if out.Code != 0 {
		t.Errorf("exec sync should return 0 code but is %v", out.Code)
	}
}
func TestExecSyncFailure(t *testing.T) {
	lt := newLXFTest(t)
	lt.createContainer(&lxf.Container{
		Name:    "roosevelt",
		Sandbox: setUpSandbox(lt, "roosevelt"),
		Image:   imgBusybox,
	})

	lt.startContainer("roosevelt")
	out := lt.execSync("roosevelt", []string{"ls", "--alala"})
	if out.Code == 0 {
		t.Errorf("exec sync should not return 0 code")
	}
}

func TestNonInteractiveExec(t *testing.T) {
	lt := newLXFTest(t)
	lt.createContainer(&lxf.Container{
		Name:    "roosevelt",
		Sandbox: setUpSandbox(lt, "roosevelt"),
		Image:   imgBusybox,
	})

	lt.startContainer("roosevelt")

	outin, stdout := io.Pipe()
	errin, stderr := io.Pipe()
	output := "fjdskfhk d saf sdfrsafrefewfesfsdfdsfsd"

	go func() {
		ioutil.ReadAll(errin)
	}()

	go func() {
		stdoutbytes, err := ioutil.ReadAll(outin)
		if err != nil {
			t.Errorf("failed to read from exec stdout, %v", err)
		}
		actual := string(stdoutbytes)
		expected := output + "\n"
		if expected != actual {
			t.Errorf("expected stdout to be '%v' but is '%v'", expected, actual)
		}
	}()

	lt.exec("roosevelt", []string{"echo", output}, nil, stdout, stderr)

}
