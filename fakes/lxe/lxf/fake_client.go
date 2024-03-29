// Code generated by counterfeiter. DO NOT EDIT.
package lxf

import (
	"io"
	"sync"

	"github.com/automaticserver/lxe/lxf"
	lxd "github.com/lxc/lxd/client"
	"k8s.io/client-go/tools/remotecommand"
)

type FakeClient struct {
	ExecStub        func(string, []string, io.ReadCloser, io.WriteCloser, io.WriteCloser, bool, bool, int64, <-chan remotecommand.TerminalSize) (int32, error)
	execMutex       sync.RWMutex
	execArgsForCall []struct {
		arg1 string
		arg2 []string
		arg3 io.ReadCloser
		arg4 io.WriteCloser
		arg5 io.WriteCloser
		arg6 bool
		arg7 bool
		arg8 int64
		arg9 <-chan remotecommand.TerminalSize
	}
	execReturns struct {
		result1 int32
		result2 error
	}
	execReturnsOnCall map[int]struct {
		result1 int32
		result2 error
	}
	GetContainerStub        func(string) (*lxf.Container, error)
	getContainerMutex       sync.RWMutex
	getContainerArgsForCall []struct {
		arg1 string
	}
	getContainerReturns struct {
		result1 *lxf.Container
		result2 error
	}
	getContainerReturnsOnCall map[int]struct {
		result1 *lxf.Container
		result2 error
	}
	GetFSPoolUsageStub        func() ([]lxf.FSPoolUsage, error)
	getFSPoolUsageMutex       sync.RWMutex
	getFSPoolUsageArgsForCall []struct {
	}
	getFSPoolUsageReturns struct {
		result1 []lxf.FSPoolUsage
		result2 error
	}
	getFSPoolUsageReturnsOnCall map[int]struct {
		result1 []lxf.FSPoolUsage
		result2 error
	}
	GetImageStub        func(string) (*lxf.Image, error)
	getImageMutex       sync.RWMutex
	getImageArgsForCall []struct {
		arg1 string
	}
	getImageReturns struct {
		result1 *lxf.Image
		result2 error
	}
	getImageReturnsOnCall map[int]struct {
		result1 *lxf.Image
		result2 error
	}
	GetRuntimeInfoStub        func() (*lxf.RuntimeInfo, error)
	getRuntimeInfoMutex       sync.RWMutex
	getRuntimeInfoArgsForCall []struct {
	}
	getRuntimeInfoReturns struct {
		result1 *lxf.RuntimeInfo
		result2 error
	}
	getRuntimeInfoReturnsOnCall map[int]struct {
		result1 *lxf.RuntimeInfo
		result2 error
	}
	GetSandboxStub        func(string) (*lxf.Sandbox, error)
	getSandboxMutex       sync.RWMutex
	getSandboxArgsForCall []struct {
		arg1 string
	}
	getSandboxReturns struct {
		result1 *lxf.Sandbox
		result2 error
	}
	getSandboxReturnsOnCall map[int]struct {
		result1 *lxf.Sandbox
		result2 error
	}
	GetServerStub        func() lxd.ContainerServer
	getServerMutex       sync.RWMutex
	getServerArgsForCall []struct {
	}
	getServerReturns struct {
		result1 lxd.ContainerServer
	}
	getServerReturnsOnCall map[int]struct {
		result1 lxd.ContainerServer
	}
	ListContainersStub        func() ([]*lxf.Container, error)
	listContainersMutex       sync.RWMutex
	listContainersArgsForCall []struct {
	}
	listContainersReturns struct {
		result1 []*lxf.Container
		result2 error
	}
	listContainersReturnsOnCall map[int]struct {
		result1 []*lxf.Container
		result2 error
	}
	ListImagesStub        func(string) ([]*lxf.Image, error)
	listImagesMutex       sync.RWMutex
	listImagesArgsForCall []struct {
		arg1 string
	}
	listImagesReturns struct {
		result1 []*lxf.Image
		result2 error
	}
	listImagesReturnsOnCall map[int]struct {
		result1 []*lxf.Image
		result2 error
	}
	ListSandboxesStub        func() ([]*lxf.Sandbox, error)
	listSandboxesMutex       sync.RWMutex
	listSandboxesArgsForCall []struct {
	}
	listSandboxesReturns struct {
		result1 []*lxf.Sandbox
		result2 error
	}
	listSandboxesReturnsOnCall map[int]struct {
		result1 []*lxf.Sandbox
		result2 error
	}
	NewContainerStub        func(string, ...string) *lxf.Container
	newContainerMutex       sync.RWMutex
	newContainerArgsForCall []struct {
		arg1 string
		arg2 []string
	}
	newContainerReturns struct {
		result1 *lxf.Container
	}
	newContainerReturnsOnCall map[int]struct {
		result1 *lxf.Container
	}
	NewSandboxStub        func() *lxf.Sandbox
	newSandboxMutex       sync.RWMutex
	newSandboxArgsForCall []struct {
	}
	newSandboxReturns struct {
		result1 *lxf.Sandbox
	}
	newSandboxReturnsOnCall map[int]struct {
		result1 *lxf.Sandbox
	}
	PullImageStub        func(string) (string, error)
	pullImageMutex       sync.RWMutex
	pullImageArgsForCall []struct {
		arg1 string
	}
	pullImageReturns struct {
		result1 string
		result2 error
	}
	pullImageReturnsOnCall map[int]struct {
		result1 string
		result2 error
	}
	RemoveImageStub        func(string) error
	removeImageMutex       sync.RWMutex
	removeImageArgsForCall []struct {
		arg1 string
	}
	removeImageReturns struct {
		result1 error
	}
	removeImageReturnsOnCall map[int]struct {
		result1 error
	}
	SetCRITestModeStub        func()
	setCRITestModeMutex       sync.RWMutex
	setCRITestModeArgsForCall []struct {
	}
	SetEventHandlerStub        func(lxf.EventHandler)
	setEventHandlerMutex       sync.RWMutex
	setEventHandlerArgsForCall []struct {
		arg1 lxf.EventHandler
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeClient) Exec(arg1 string, arg2 []string, arg3 io.ReadCloser, arg4 io.WriteCloser, arg5 io.WriteCloser, arg6 bool, arg7 bool, arg8 int64, arg9 <-chan remotecommand.TerminalSize) (int32, error) {
	var arg2Copy []string
	if arg2 != nil {
		arg2Copy = make([]string, len(arg2))
		copy(arg2Copy, arg2)
	}
	fake.execMutex.Lock()
	ret, specificReturn := fake.execReturnsOnCall[len(fake.execArgsForCall)]
	fake.execArgsForCall = append(fake.execArgsForCall, struct {
		arg1 string
		arg2 []string
		arg3 io.ReadCloser
		arg4 io.WriteCloser
		arg5 io.WriteCloser
		arg6 bool
		arg7 bool
		arg8 int64
		arg9 <-chan remotecommand.TerminalSize
	}{arg1, arg2Copy, arg3, arg4, arg5, arg6, arg7, arg8, arg9})
	stub := fake.ExecStub
	fakeReturns := fake.execReturns
	fake.recordInvocation("Exec", []interface{}{arg1, arg2Copy, arg3, arg4, arg5, arg6, arg7, arg8, arg9})
	fake.execMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8, arg9)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeClient) ExecCallCount() int {
	fake.execMutex.RLock()
	defer fake.execMutex.RUnlock()
	return len(fake.execArgsForCall)
}

func (fake *FakeClient) ExecCalls(stub func(string, []string, io.ReadCloser, io.WriteCloser, io.WriteCloser, bool, bool, int64, <-chan remotecommand.TerminalSize) (int32, error)) {
	fake.execMutex.Lock()
	defer fake.execMutex.Unlock()
	fake.ExecStub = stub
}

func (fake *FakeClient) ExecArgsForCall(i int) (string, []string, io.ReadCloser, io.WriteCloser, io.WriteCloser, bool, bool, int64, <-chan remotecommand.TerminalSize) {
	fake.execMutex.RLock()
	defer fake.execMutex.RUnlock()
	argsForCall := fake.execArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4, argsForCall.arg5, argsForCall.arg6, argsForCall.arg7, argsForCall.arg8, argsForCall.arg9
}

func (fake *FakeClient) ExecReturns(result1 int32, result2 error) {
	fake.execMutex.Lock()
	defer fake.execMutex.Unlock()
	fake.ExecStub = nil
	fake.execReturns = struct {
		result1 int32
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) ExecReturnsOnCall(i int, result1 int32, result2 error) {
	fake.execMutex.Lock()
	defer fake.execMutex.Unlock()
	fake.ExecStub = nil
	if fake.execReturnsOnCall == nil {
		fake.execReturnsOnCall = make(map[int]struct {
			result1 int32
			result2 error
		})
	}
	fake.execReturnsOnCall[i] = struct {
		result1 int32
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) GetContainer(arg1 string) (*lxf.Container, error) {
	fake.getContainerMutex.Lock()
	ret, specificReturn := fake.getContainerReturnsOnCall[len(fake.getContainerArgsForCall)]
	fake.getContainerArgsForCall = append(fake.getContainerArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.GetContainerStub
	fakeReturns := fake.getContainerReturns
	fake.recordInvocation("GetContainer", []interface{}{arg1})
	fake.getContainerMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeClient) GetContainerCallCount() int {
	fake.getContainerMutex.RLock()
	defer fake.getContainerMutex.RUnlock()
	return len(fake.getContainerArgsForCall)
}

func (fake *FakeClient) GetContainerCalls(stub func(string) (*lxf.Container, error)) {
	fake.getContainerMutex.Lock()
	defer fake.getContainerMutex.Unlock()
	fake.GetContainerStub = stub
}

func (fake *FakeClient) GetContainerArgsForCall(i int) string {
	fake.getContainerMutex.RLock()
	defer fake.getContainerMutex.RUnlock()
	argsForCall := fake.getContainerArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeClient) GetContainerReturns(result1 *lxf.Container, result2 error) {
	fake.getContainerMutex.Lock()
	defer fake.getContainerMutex.Unlock()
	fake.GetContainerStub = nil
	fake.getContainerReturns = struct {
		result1 *lxf.Container
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) GetContainerReturnsOnCall(i int, result1 *lxf.Container, result2 error) {
	fake.getContainerMutex.Lock()
	defer fake.getContainerMutex.Unlock()
	fake.GetContainerStub = nil
	if fake.getContainerReturnsOnCall == nil {
		fake.getContainerReturnsOnCall = make(map[int]struct {
			result1 *lxf.Container
			result2 error
		})
	}
	fake.getContainerReturnsOnCall[i] = struct {
		result1 *lxf.Container
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) GetFSPoolUsage() ([]lxf.FSPoolUsage, error) {
	fake.getFSPoolUsageMutex.Lock()
	ret, specificReturn := fake.getFSPoolUsageReturnsOnCall[len(fake.getFSPoolUsageArgsForCall)]
	fake.getFSPoolUsageArgsForCall = append(fake.getFSPoolUsageArgsForCall, struct {
	}{})
	stub := fake.GetFSPoolUsageStub
	fakeReturns := fake.getFSPoolUsageReturns
	fake.recordInvocation("GetFSPoolUsage", []interface{}{})
	fake.getFSPoolUsageMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeClient) GetFSPoolUsageCallCount() int {
	fake.getFSPoolUsageMutex.RLock()
	defer fake.getFSPoolUsageMutex.RUnlock()
	return len(fake.getFSPoolUsageArgsForCall)
}

func (fake *FakeClient) GetFSPoolUsageCalls(stub func() ([]lxf.FSPoolUsage, error)) {
	fake.getFSPoolUsageMutex.Lock()
	defer fake.getFSPoolUsageMutex.Unlock()
	fake.GetFSPoolUsageStub = stub
}

func (fake *FakeClient) GetFSPoolUsageReturns(result1 []lxf.FSPoolUsage, result2 error) {
	fake.getFSPoolUsageMutex.Lock()
	defer fake.getFSPoolUsageMutex.Unlock()
	fake.GetFSPoolUsageStub = nil
	fake.getFSPoolUsageReturns = struct {
		result1 []lxf.FSPoolUsage
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) GetFSPoolUsageReturnsOnCall(i int, result1 []lxf.FSPoolUsage, result2 error) {
	fake.getFSPoolUsageMutex.Lock()
	defer fake.getFSPoolUsageMutex.Unlock()
	fake.GetFSPoolUsageStub = nil
	if fake.getFSPoolUsageReturnsOnCall == nil {
		fake.getFSPoolUsageReturnsOnCall = make(map[int]struct {
			result1 []lxf.FSPoolUsage
			result2 error
		})
	}
	fake.getFSPoolUsageReturnsOnCall[i] = struct {
		result1 []lxf.FSPoolUsage
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) GetImage(arg1 string) (*lxf.Image, error) {
	fake.getImageMutex.Lock()
	ret, specificReturn := fake.getImageReturnsOnCall[len(fake.getImageArgsForCall)]
	fake.getImageArgsForCall = append(fake.getImageArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.GetImageStub
	fakeReturns := fake.getImageReturns
	fake.recordInvocation("GetImage", []interface{}{arg1})
	fake.getImageMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeClient) GetImageCallCount() int {
	fake.getImageMutex.RLock()
	defer fake.getImageMutex.RUnlock()
	return len(fake.getImageArgsForCall)
}

func (fake *FakeClient) GetImageCalls(stub func(string) (*lxf.Image, error)) {
	fake.getImageMutex.Lock()
	defer fake.getImageMutex.Unlock()
	fake.GetImageStub = stub
}

func (fake *FakeClient) GetImageArgsForCall(i int) string {
	fake.getImageMutex.RLock()
	defer fake.getImageMutex.RUnlock()
	argsForCall := fake.getImageArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeClient) GetImageReturns(result1 *lxf.Image, result2 error) {
	fake.getImageMutex.Lock()
	defer fake.getImageMutex.Unlock()
	fake.GetImageStub = nil
	fake.getImageReturns = struct {
		result1 *lxf.Image
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) GetImageReturnsOnCall(i int, result1 *lxf.Image, result2 error) {
	fake.getImageMutex.Lock()
	defer fake.getImageMutex.Unlock()
	fake.GetImageStub = nil
	if fake.getImageReturnsOnCall == nil {
		fake.getImageReturnsOnCall = make(map[int]struct {
			result1 *lxf.Image
			result2 error
		})
	}
	fake.getImageReturnsOnCall[i] = struct {
		result1 *lxf.Image
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) GetRuntimeInfo() (*lxf.RuntimeInfo, error) {
	fake.getRuntimeInfoMutex.Lock()
	ret, specificReturn := fake.getRuntimeInfoReturnsOnCall[len(fake.getRuntimeInfoArgsForCall)]
	fake.getRuntimeInfoArgsForCall = append(fake.getRuntimeInfoArgsForCall, struct {
	}{})
	stub := fake.GetRuntimeInfoStub
	fakeReturns := fake.getRuntimeInfoReturns
	fake.recordInvocation("GetRuntimeInfo", []interface{}{})
	fake.getRuntimeInfoMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeClient) GetRuntimeInfoCallCount() int {
	fake.getRuntimeInfoMutex.RLock()
	defer fake.getRuntimeInfoMutex.RUnlock()
	return len(fake.getRuntimeInfoArgsForCall)
}

func (fake *FakeClient) GetRuntimeInfoCalls(stub func() (*lxf.RuntimeInfo, error)) {
	fake.getRuntimeInfoMutex.Lock()
	defer fake.getRuntimeInfoMutex.Unlock()
	fake.GetRuntimeInfoStub = stub
}

func (fake *FakeClient) GetRuntimeInfoReturns(result1 *lxf.RuntimeInfo, result2 error) {
	fake.getRuntimeInfoMutex.Lock()
	defer fake.getRuntimeInfoMutex.Unlock()
	fake.GetRuntimeInfoStub = nil
	fake.getRuntimeInfoReturns = struct {
		result1 *lxf.RuntimeInfo
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) GetRuntimeInfoReturnsOnCall(i int, result1 *lxf.RuntimeInfo, result2 error) {
	fake.getRuntimeInfoMutex.Lock()
	defer fake.getRuntimeInfoMutex.Unlock()
	fake.GetRuntimeInfoStub = nil
	if fake.getRuntimeInfoReturnsOnCall == nil {
		fake.getRuntimeInfoReturnsOnCall = make(map[int]struct {
			result1 *lxf.RuntimeInfo
			result2 error
		})
	}
	fake.getRuntimeInfoReturnsOnCall[i] = struct {
		result1 *lxf.RuntimeInfo
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) GetSandbox(arg1 string) (*lxf.Sandbox, error) {
	fake.getSandboxMutex.Lock()
	ret, specificReturn := fake.getSandboxReturnsOnCall[len(fake.getSandboxArgsForCall)]
	fake.getSandboxArgsForCall = append(fake.getSandboxArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.GetSandboxStub
	fakeReturns := fake.getSandboxReturns
	fake.recordInvocation("GetSandbox", []interface{}{arg1})
	fake.getSandboxMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeClient) GetSandboxCallCount() int {
	fake.getSandboxMutex.RLock()
	defer fake.getSandboxMutex.RUnlock()
	return len(fake.getSandboxArgsForCall)
}

func (fake *FakeClient) GetSandboxCalls(stub func(string) (*lxf.Sandbox, error)) {
	fake.getSandboxMutex.Lock()
	defer fake.getSandboxMutex.Unlock()
	fake.GetSandboxStub = stub
}

func (fake *FakeClient) GetSandboxArgsForCall(i int) string {
	fake.getSandboxMutex.RLock()
	defer fake.getSandboxMutex.RUnlock()
	argsForCall := fake.getSandboxArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeClient) GetSandboxReturns(result1 *lxf.Sandbox, result2 error) {
	fake.getSandboxMutex.Lock()
	defer fake.getSandboxMutex.Unlock()
	fake.GetSandboxStub = nil
	fake.getSandboxReturns = struct {
		result1 *lxf.Sandbox
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) GetSandboxReturnsOnCall(i int, result1 *lxf.Sandbox, result2 error) {
	fake.getSandboxMutex.Lock()
	defer fake.getSandboxMutex.Unlock()
	fake.GetSandboxStub = nil
	if fake.getSandboxReturnsOnCall == nil {
		fake.getSandboxReturnsOnCall = make(map[int]struct {
			result1 *lxf.Sandbox
			result2 error
		})
	}
	fake.getSandboxReturnsOnCall[i] = struct {
		result1 *lxf.Sandbox
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) GetServer() lxd.ContainerServer {
	fake.getServerMutex.Lock()
	ret, specificReturn := fake.getServerReturnsOnCall[len(fake.getServerArgsForCall)]
	fake.getServerArgsForCall = append(fake.getServerArgsForCall, struct {
	}{})
	stub := fake.GetServerStub
	fakeReturns := fake.getServerReturns
	fake.recordInvocation("GetServer", []interface{}{})
	fake.getServerMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeClient) GetServerCallCount() int {
	fake.getServerMutex.RLock()
	defer fake.getServerMutex.RUnlock()
	return len(fake.getServerArgsForCall)
}

func (fake *FakeClient) GetServerCalls(stub func() lxd.ContainerServer) {
	fake.getServerMutex.Lock()
	defer fake.getServerMutex.Unlock()
	fake.GetServerStub = stub
}

func (fake *FakeClient) GetServerReturns(result1 lxd.ContainerServer) {
	fake.getServerMutex.Lock()
	defer fake.getServerMutex.Unlock()
	fake.GetServerStub = nil
	fake.getServerReturns = struct {
		result1 lxd.ContainerServer
	}{result1}
}

func (fake *FakeClient) GetServerReturnsOnCall(i int, result1 lxd.ContainerServer) {
	fake.getServerMutex.Lock()
	defer fake.getServerMutex.Unlock()
	fake.GetServerStub = nil
	if fake.getServerReturnsOnCall == nil {
		fake.getServerReturnsOnCall = make(map[int]struct {
			result1 lxd.ContainerServer
		})
	}
	fake.getServerReturnsOnCall[i] = struct {
		result1 lxd.ContainerServer
	}{result1}
}

func (fake *FakeClient) ListContainers() ([]*lxf.Container, error) {
	fake.listContainersMutex.Lock()
	ret, specificReturn := fake.listContainersReturnsOnCall[len(fake.listContainersArgsForCall)]
	fake.listContainersArgsForCall = append(fake.listContainersArgsForCall, struct {
	}{})
	stub := fake.ListContainersStub
	fakeReturns := fake.listContainersReturns
	fake.recordInvocation("ListContainers", []interface{}{})
	fake.listContainersMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeClient) ListContainersCallCount() int {
	fake.listContainersMutex.RLock()
	defer fake.listContainersMutex.RUnlock()
	return len(fake.listContainersArgsForCall)
}

func (fake *FakeClient) ListContainersCalls(stub func() ([]*lxf.Container, error)) {
	fake.listContainersMutex.Lock()
	defer fake.listContainersMutex.Unlock()
	fake.ListContainersStub = stub
}

func (fake *FakeClient) ListContainersReturns(result1 []*lxf.Container, result2 error) {
	fake.listContainersMutex.Lock()
	defer fake.listContainersMutex.Unlock()
	fake.ListContainersStub = nil
	fake.listContainersReturns = struct {
		result1 []*lxf.Container
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) ListContainersReturnsOnCall(i int, result1 []*lxf.Container, result2 error) {
	fake.listContainersMutex.Lock()
	defer fake.listContainersMutex.Unlock()
	fake.ListContainersStub = nil
	if fake.listContainersReturnsOnCall == nil {
		fake.listContainersReturnsOnCall = make(map[int]struct {
			result1 []*lxf.Container
			result2 error
		})
	}
	fake.listContainersReturnsOnCall[i] = struct {
		result1 []*lxf.Container
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) ListImages(arg1 string) ([]*lxf.Image, error) {
	fake.listImagesMutex.Lock()
	ret, specificReturn := fake.listImagesReturnsOnCall[len(fake.listImagesArgsForCall)]
	fake.listImagesArgsForCall = append(fake.listImagesArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.ListImagesStub
	fakeReturns := fake.listImagesReturns
	fake.recordInvocation("ListImages", []interface{}{arg1})
	fake.listImagesMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeClient) ListImagesCallCount() int {
	fake.listImagesMutex.RLock()
	defer fake.listImagesMutex.RUnlock()
	return len(fake.listImagesArgsForCall)
}

func (fake *FakeClient) ListImagesCalls(stub func(string) ([]*lxf.Image, error)) {
	fake.listImagesMutex.Lock()
	defer fake.listImagesMutex.Unlock()
	fake.ListImagesStub = stub
}

func (fake *FakeClient) ListImagesArgsForCall(i int) string {
	fake.listImagesMutex.RLock()
	defer fake.listImagesMutex.RUnlock()
	argsForCall := fake.listImagesArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeClient) ListImagesReturns(result1 []*lxf.Image, result2 error) {
	fake.listImagesMutex.Lock()
	defer fake.listImagesMutex.Unlock()
	fake.ListImagesStub = nil
	fake.listImagesReturns = struct {
		result1 []*lxf.Image
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) ListImagesReturnsOnCall(i int, result1 []*lxf.Image, result2 error) {
	fake.listImagesMutex.Lock()
	defer fake.listImagesMutex.Unlock()
	fake.ListImagesStub = nil
	if fake.listImagesReturnsOnCall == nil {
		fake.listImagesReturnsOnCall = make(map[int]struct {
			result1 []*lxf.Image
			result2 error
		})
	}
	fake.listImagesReturnsOnCall[i] = struct {
		result1 []*lxf.Image
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) ListSandboxes() ([]*lxf.Sandbox, error) {
	fake.listSandboxesMutex.Lock()
	ret, specificReturn := fake.listSandboxesReturnsOnCall[len(fake.listSandboxesArgsForCall)]
	fake.listSandboxesArgsForCall = append(fake.listSandboxesArgsForCall, struct {
	}{})
	stub := fake.ListSandboxesStub
	fakeReturns := fake.listSandboxesReturns
	fake.recordInvocation("ListSandboxes", []interface{}{})
	fake.listSandboxesMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeClient) ListSandboxesCallCount() int {
	fake.listSandboxesMutex.RLock()
	defer fake.listSandboxesMutex.RUnlock()
	return len(fake.listSandboxesArgsForCall)
}

func (fake *FakeClient) ListSandboxesCalls(stub func() ([]*lxf.Sandbox, error)) {
	fake.listSandboxesMutex.Lock()
	defer fake.listSandboxesMutex.Unlock()
	fake.ListSandboxesStub = stub
}

func (fake *FakeClient) ListSandboxesReturns(result1 []*lxf.Sandbox, result2 error) {
	fake.listSandboxesMutex.Lock()
	defer fake.listSandboxesMutex.Unlock()
	fake.ListSandboxesStub = nil
	fake.listSandboxesReturns = struct {
		result1 []*lxf.Sandbox
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) ListSandboxesReturnsOnCall(i int, result1 []*lxf.Sandbox, result2 error) {
	fake.listSandboxesMutex.Lock()
	defer fake.listSandboxesMutex.Unlock()
	fake.ListSandboxesStub = nil
	if fake.listSandboxesReturnsOnCall == nil {
		fake.listSandboxesReturnsOnCall = make(map[int]struct {
			result1 []*lxf.Sandbox
			result2 error
		})
	}
	fake.listSandboxesReturnsOnCall[i] = struct {
		result1 []*lxf.Sandbox
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) NewContainer(arg1 string, arg2 ...string) *lxf.Container {
	fake.newContainerMutex.Lock()
	ret, specificReturn := fake.newContainerReturnsOnCall[len(fake.newContainerArgsForCall)]
	fake.newContainerArgsForCall = append(fake.newContainerArgsForCall, struct {
		arg1 string
		arg2 []string
	}{arg1, arg2})
	stub := fake.NewContainerStub
	fakeReturns := fake.newContainerReturns
	fake.recordInvocation("NewContainer", []interface{}{arg1, arg2})
	fake.newContainerMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2...)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeClient) NewContainerCallCount() int {
	fake.newContainerMutex.RLock()
	defer fake.newContainerMutex.RUnlock()
	return len(fake.newContainerArgsForCall)
}

func (fake *FakeClient) NewContainerCalls(stub func(string, ...string) *lxf.Container) {
	fake.newContainerMutex.Lock()
	defer fake.newContainerMutex.Unlock()
	fake.NewContainerStub = stub
}

func (fake *FakeClient) NewContainerArgsForCall(i int) (string, []string) {
	fake.newContainerMutex.RLock()
	defer fake.newContainerMutex.RUnlock()
	argsForCall := fake.newContainerArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeClient) NewContainerReturns(result1 *lxf.Container) {
	fake.newContainerMutex.Lock()
	defer fake.newContainerMutex.Unlock()
	fake.NewContainerStub = nil
	fake.newContainerReturns = struct {
		result1 *lxf.Container
	}{result1}
}

func (fake *FakeClient) NewContainerReturnsOnCall(i int, result1 *lxf.Container) {
	fake.newContainerMutex.Lock()
	defer fake.newContainerMutex.Unlock()
	fake.NewContainerStub = nil
	if fake.newContainerReturnsOnCall == nil {
		fake.newContainerReturnsOnCall = make(map[int]struct {
			result1 *lxf.Container
		})
	}
	fake.newContainerReturnsOnCall[i] = struct {
		result1 *lxf.Container
	}{result1}
}

func (fake *FakeClient) NewSandbox() *lxf.Sandbox {
	fake.newSandboxMutex.Lock()
	ret, specificReturn := fake.newSandboxReturnsOnCall[len(fake.newSandboxArgsForCall)]
	fake.newSandboxArgsForCall = append(fake.newSandboxArgsForCall, struct {
	}{})
	stub := fake.NewSandboxStub
	fakeReturns := fake.newSandboxReturns
	fake.recordInvocation("NewSandbox", []interface{}{})
	fake.newSandboxMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeClient) NewSandboxCallCount() int {
	fake.newSandboxMutex.RLock()
	defer fake.newSandboxMutex.RUnlock()
	return len(fake.newSandboxArgsForCall)
}

func (fake *FakeClient) NewSandboxCalls(stub func() *lxf.Sandbox) {
	fake.newSandboxMutex.Lock()
	defer fake.newSandboxMutex.Unlock()
	fake.NewSandboxStub = stub
}

func (fake *FakeClient) NewSandboxReturns(result1 *lxf.Sandbox) {
	fake.newSandboxMutex.Lock()
	defer fake.newSandboxMutex.Unlock()
	fake.NewSandboxStub = nil
	fake.newSandboxReturns = struct {
		result1 *lxf.Sandbox
	}{result1}
}

func (fake *FakeClient) NewSandboxReturnsOnCall(i int, result1 *lxf.Sandbox) {
	fake.newSandboxMutex.Lock()
	defer fake.newSandboxMutex.Unlock()
	fake.NewSandboxStub = nil
	if fake.newSandboxReturnsOnCall == nil {
		fake.newSandboxReturnsOnCall = make(map[int]struct {
			result1 *lxf.Sandbox
		})
	}
	fake.newSandboxReturnsOnCall[i] = struct {
		result1 *lxf.Sandbox
	}{result1}
}

func (fake *FakeClient) PullImage(arg1 string) (string, error) {
	fake.pullImageMutex.Lock()
	ret, specificReturn := fake.pullImageReturnsOnCall[len(fake.pullImageArgsForCall)]
	fake.pullImageArgsForCall = append(fake.pullImageArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.PullImageStub
	fakeReturns := fake.pullImageReturns
	fake.recordInvocation("PullImage", []interface{}{arg1})
	fake.pullImageMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeClient) PullImageCallCount() int {
	fake.pullImageMutex.RLock()
	defer fake.pullImageMutex.RUnlock()
	return len(fake.pullImageArgsForCall)
}

func (fake *FakeClient) PullImageCalls(stub func(string) (string, error)) {
	fake.pullImageMutex.Lock()
	defer fake.pullImageMutex.Unlock()
	fake.PullImageStub = stub
}

func (fake *FakeClient) PullImageArgsForCall(i int) string {
	fake.pullImageMutex.RLock()
	defer fake.pullImageMutex.RUnlock()
	argsForCall := fake.pullImageArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeClient) PullImageReturns(result1 string, result2 error) {
	fake.pullImageMutex.Lock()
	defer fake.pullImageMutex.Unlock()
	fake.PullImageStub = nil
	fake.pullImageReturns = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) PullImageReturnsOnCall(i int, result1 string, result2 error) {
	fake.pullImageMutex.Lock()
	defer fake.pullImageMutex.Unlock()
	fake.PullImageStub = nil
	if fake.pullImageReturnsOnCall == nil {
		fake.pullImageReturnsOnCall = make(map[int]struct {
			result1 string
			result2 error
		})
	}
	fake.pullImageReturnsOnCall[i] = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) RemoveImage(arg1 string) error {
	fake.removeImageMutex.Lock()
	ret, specificReturn := fake.removeImageReturnsOnCall[len(fake.removeImageArgsForCall)]
	fake.removeImageArgsForCall = append(fake.removeImageArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.RemoveImageStub
	fakeReturns := fake.removeImageReturns
	fake.recordInvocation("RemoveImage", []interface{}{arg1})
	fake.removeImageMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeClient) RemoveImageCallCount() int {
	fake.removeImageMutex.RLock()
	defer fake.removeImageMutex.RUnlock()
	return len(fake.removeImageArgsForCall)
}

func (fake *FakeClient) RemoveImageCalls(stub func(string) error) {
	fake.removeImageMutex.Lock()
	defer fake.removeImageMutex.Unlock()
	fake.RemoveImageStub = stub
}

func (fake *FakeClient) RemoveImageArgsForCall(i int) string {
	fake.removeImageMutex.RLock()
	defer fake.removeImageMutex.RUnlock()
	argsForCall := fake.removeImageArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeClient) RemoveImageReturns(result1 error) {
	fake.removeImageMutex.Lock()
	defer fake.removeImageMutex.Unlock()
	fake.RemoveImageStub = nil
	fake.removeImageReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeClient) RemoveImageReturnsOnCall(i int, result1 error) {
	fake.removeImageMutex.Lock()
	defer fake.removeImageMutex.Unlock()
	fake.RemoveImageStub = nil
	if fake.removeImageReturnsOnCall == nil {
		fake.removeImageReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.removeImageReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeClient) SetCRITestMode() {
	fake.setCRITestModeMutex.Lock()
	fake.setCRITestModeArgsForCall = append(fake.setCRITestModeArgsForCall, struct {
	}{})
	stub := fake.SetCRITestModeStub
	fake.recordInvocation("SetCRITestMode", []interface{}{})
	fake.setCRITestModeMutex.Unlock()
	if stub != nil {
		fake.SetCRITestModeStub()
	}
}

func (fake *FakeClient) SetCRITestModeCallCount() int {
	fake.setCRITestModeMutex.RLock()
	defer fake.setCRITestModeMutex.RUnlock()
	return len(fake.setCRITestModeArgsForCall)
}

func (fake *FakeClient) SetCRITestModeCalls(stub func()) {
	fake.setCRITestModeMutex.Lock()
	defer fake.setCRITestModeMutex.Unlock()
	fake.SetCRITestModeStub = stub
}

func (fake *FakeClient) SetEventHandler(arg1 lxf.EventHandler) {
	fake.setEventHandlerMutex.Lock()
	fake.setEventHandlerArgsForCall = append(fake.setEventHandlerArgsForCall, struct {
		arg1 lxf.EventHandler
	}{arg1})
	stub := fake.SetEventHandlerStub
	fake.recordInvocation("SetEventHandler", []interface{}{arg1})
	fake.setEventHandlerMutex.Unlock()
	if stub != nil {
		fake.SetEventHandlerStub(arg1)
	}
}

func (fake *FakeClient) SetEventHandlerCallCount() int {
	fake.setEventHandlerMutex.RLock()
	defer fake.setEventHandlerMutex.RUnlock()
	return len(fake.setEventHandlerArgsForCall)
}

func (fake *FakeClient) SetEventHandlerCalls(stub func(lxf.EventHandler)) {
	fake.setEventHandlerMutex.Lock()
	defer fake.setEventHandlerMutex.Unlock()
	fake.SetEventHandlerStub = stub
}

func (fake *FakeClient) SetEventHandlerArgsForCall(i int) lxf.EventHandler {
	fake.setEventHandlerMutex.RLock()
	defer fake.setEventHandlerMutex.RUnlock()
	argsForCall := fake.setEventHandlerArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeClient) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.execMutex.RLock()
	defer fake.execMutex.RUnlock()
	fake.getContainerMutex.RLock()
	defer fake.getContainerMutex.RUnlock()
	fake.getFSPoolUsageMutex.RLock()
	defer fake.getFSPoolUsageMutex.RUnlock()
	fake.getImageMutex.RLock()
	defer fake.getImageMutex.RUnlock()
	fake.getRuntimeInfoMutex.RLock()
	defer fake.getRuntimeInfoMutex.RUnlock()
	fake.getSandboxMutex.RLock()
	defer fake.getSandboxMutex.RUnlock()
	fake.getServerMutex.RLock()
	defer fake.getServerMutex.RUnlock()
	fake.listContainersMutex.RLock()
	defer fake.listContainersMutex.RUnlock()
	fake.listImagesMutex.RLock()
	defer fake.listImagesMutex.RUnlock()
	fake.listSandboxesMutex.RLock()
	defer fake.listSandboxesMutex.RUnlock()
	fake.newContainerMutex.RLock()
	defer fake.newContainerMutex.RUnlock()
	fake.newSandboxMutex.RLock()
	defer fake.newSandboxMutex.RUnlock()
	fake.pullImageMutex.RLock()
	defer fake.pullImageMutex.RUnlock()
	fake.removeImageMutex.RLock()
	defer fake.removeImageMutex.RUnlock()
	fake.setCRITestModeMutex.RLock()
	defer fake.setCRITestModeMutex.RUnlock()
	fake.setEventHandlerMutex.RLock()
	defer fake.setEventHandlerMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeClient) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ lxf.Client = new(FakeClient)
