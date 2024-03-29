// Code generated by counterfeiter. DO NOT EDIT.
package libcni

import (
	"context"
	"sync"

	"github.com/containernetworking/cni/libcni"
	"github.com/containernetworking/cni/pkg/types"
)

type FakeCNI struct {
	AddNetworkStub        func(context.Context, *libcni.NetworkConfig, *libcni.RuntimeConf) (types.Result, error)
	addNetworkMutex       sync.RWMutex
	addNetworkArgsForCall []struct {
		arg1 context.Context
		arg2 *libcni.NetworkConfig
		arg3 *libcni.RuntimeConf
	}
	addNetworkReturns struct {
		result1 types.Result
		result2 error
	}
	addNetworkReturnsOnCall map[int]struct {
		result1 types.Result
		result2 error
	}
	AddNetworkListStub        func(context.Context, *libcni.NetworkConfigList, *libcni.RuntimeConf) (types.Result, error)
	addNetworkListMutex       sync.RWMutex
	addNetworkListArgsForCall []struct {
		arg1 context.Context
		arg2 *libcni.NetworkConfigList
		arg3 *libcni.RuntimeConf
	}
	addNetworkListReturns struct {
		result1 types.Result
		result2 error
	}
	addNetworkListReturnsOnCall map[int]struct {
		result1 types.Result
		result2 error
	}
	CheckNetworkStub        func(context.Context, *libcni.NetworkConfig, *libcni.RuntimeConf) error
	checkNetworkMutex       sync.RWMutex
	checkNetworkArgsForCall []struct {
		arg1 context.Context
		arg2 *libcni.NetworkConfig
		arg3 *libcni.RuntimeConf
	}
	checkNetworkReturns struct {
		result1 error
	}
	checkNetworkReturnsOnCall map[int]struct {
		result1 error
	}
	CheckNetworkListStub        func(context.Context, *libcni.NetworkConfigList, *libcni.RuntimeConf) error
	checkNetworkListMutex       sync.RWMutex
	checkNetworkListArgsForCall []struct {
		arg1 context.Context
		arg2 *libcni.NetworkConfigList
		arg3 *libcni.RuntimeConf
	}
	checkNetworkListReturns struct {
		result1 error
	}
	checkNetworkListReturnsOnCall map[int]struct {
		result1 error
	}
	DelNetworkStub        func(context.Context, *libcni.NetworkConfig, *libcni.RuntimeConf) error
	delNetworkMutex       sync.RWMutex
	delNetworkArgsForCall []struct {
		arg1 context.Context
		arg2 *libcni.NetworkConfig
		arg3 *libcni.RuntimeConf
	}
	delNetworkReturns struct {
		result1 error
	}
	delNetworkReturnsOnCall map[int]struct {
		result1 error
	}
	DelNetworkListStub        func(context.Context, *libcni.NetworkConfigList, *libcni.RuntimeConf) error
	delNetworkListMutex       sync.RWMutex
	delNetworkListArgsForCall []struct {
		arg1 context.Context
		arg2 *libcni.NetworkConfigList
		arg3 *libcni.RuntimeConf
	}
	delNetworkListReturns struct {
		result1 error
	}
	delNetworkListReturnsOnCall map[int]struct {
		result1 error
	}
	GetNetworkCachedConfigStub        func(*libcni.NetworkConfig, *libcni.RuntimeConf) ([]byte, *libcni.RuntimeConf, error)
	getNetworkCachedConfigMutex       sync.RWMutex
	getNetworkCachedConfigArgsForCall []struct {
		arg1 *libcni.NetworkConfig
		arg2 *libcni.RuntimeConf
	}
	getNetworkCachedConfigReturns struct {
		result1 []byte
		result2 *libcni.RuntimeConf
		result3 error
	}
	getNetworkCachedConfigReturnsOnCall map[int]struct {
		result1 []byte
		result2 *libcni.RuntimeConf
		result3 error
	}
	GetNetworkCachedResultStub        func(*libcni.NetworkConfig, *libcni.RuntimeConf) (types.Result, error)
	getNetworkCachedResultMutex       sync.RWMutex
	getNetworkCachedResultArgsForCall []struct {
		arg1 *libcni.NetworkConfig
		arg2 *libcni.RuntimeConf
	}
	getNetworkCachedResultReturns struct {
		result1 types.Result
		result2 error
	}
	getNetworkCachedResultReturnsOnCall map[int]struct {
		result1 types.Result
		result2 error
	}
	GetNetworkListCachedConfigStub        func(*libcni.NetworkConfigList, *libcni.RuntimeConf) ([]byte, *libcni.RuntimeConf, error)
	getNetworkListCachedConfigMutex       sync.RWMutex
	getNetworkListCachedConfigArgsForCall []struct {
		arg1 *libcni.NetworkConfigList
		arg2 *libcni.RuntimeConf
	}
	getNetworkListCachedConfigReturns struct {
		result1 []byte
		result2 *libcni.RuntimeConf
		result3 error
	}
	getNetworkListCachedConfigReturnsOnCall map[int]struct {
		result1 []byte
		result2 *libcni.RuntimeConf
		result3 error
	}
	GetNetworkListCachedResultStub        func(*libcni.NetworkConfigList, *libcni.RuntimeConf) (types.Result, error)
	getNetworkListCachedResultMutex       sync.RWMutex
	getNetworkListCachedResultArgsForCall []struct {
		arg1 *libcni.NetworkConfigList
		arg2 *libcni.RuntimeConf
	}
	getNetworkListCachedResultReturns struct {
		result1 types.Result
		result2 error
	}
	getNetworkListCachedResultReturnsOnCall map[int]struct {
		result1 types.Result
		result2 error
	}
	ValidateNetworkStub        func(context.Context, *libcni.NetworkConfig) ([]string, error)
	validateNetworkMutex       sync.RWMutex
	validateNetworkArgsForCall []struct {
		arg1 context.Context
		arg2 *libcni.NetworkConfig
	}
	validateNetworkReturns struct {
		result1 []string
		result2 error
	}
	validateNetworkReturnsOnCall map[int]struct {
		result1 []string
		result2 error
	}
	ValidateNetworkListStub        func(context.Context, *libcni.NetworkConfigList) ([]string, error)
	validateNetworkListMutex       sync.RWMutex
	validateNetworkListArgsForCall []struct {
		arg1 context.Context
		arg2 *libcni.NetworkConfigList
	}
	validateNetworkListReturns struct {
		result1 []string
		result2 error
	}
	validateNetworkListReturnsOnCall map[int]struct {
		result1 []string
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeCNI) AddNetwork(arg1 context.Context, arg2 *libcni.NetworkConfig, arg3 *libcni.RuntimeConf) (types.Result, error) {
	fake.addNetworkMutex.Lock()
	ret, specificReturn := fake.addNetworkReturnsOnCall[len(fake.addNetworkArgsForCall)]
	fake.addNetworkArgsForCall = append(fake.addNetworkArgsForCall, struct {
		arg1 context.Context
		arg2 *libcni.NetworkConfig
		arg3 *libcni.RuntimeConf
	}{arg1, arg2, arg3})
	stub := fake.AddNetworkStub
	fakeReturns := fake.addNetworkReturns
	fake.recordInvocation("AddNetwork", []interface{}{arg1, arg2, arg3})
	fake.addNetworkMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeCNI) AddNetworkCallCount() int {
	fake.addNetworkMutex.RLock()
	defer fake.addNetworkMutex.RUnlock()
	return len(fake.addNetworkArgsForCall)
}

func (fake *FakeCNI) AddNetworkCalls(stub func(context.Context, *libcni.NetworkConfig, *libcni.RuntimeConf) (types.Result, error)) {
	fake.addNetworkMutex.Lock()
	defer fake.addNetworkMutex.Unlock()
	fake.AddNetworkStub = stub
}

func (fake *FakeCNI) AddNetworkArgsForCall(i int) (context.Context, *libcni.NetworkConfig, *libcni.RuntimeConf) {
	fake.addNetworkMutex.RLock()
	defer fake.addNetworkMutex.RUnlock()
	argsForCall := fake.addNetworkArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeCNI) AddNetworkReturns(result1 types.Result, result2 error) {
	fake.addNetworkMutex.Lock()
	defer fake.addNetworkMutex.Unlock()
	fake.AddNetworkStub = nil
	fake.addNetworkReturns = struct {
		result1 types.Result
		result2 error
	}{result1, result2}
}

func (fake *FakeCNI) AddNetworkReturnsOnCall(i int, result1 types.Result, result2 error) {
	fake.addNetworkMutex.Lock()
	defer fake.addNetworkMutex.Unlock()
	fake.AddNetworkStub = nil
	if fake.addNetworkReturnsOnCall == nil {
		fake.addNetworkReturnsOnCall = make(map[int]struct {
			result1 types.Result
			result2 error
		})
	}
	fake.addNetworkReturnsOnCall[i] = struct {
		result1 types.Result
		result2 error
	}{result1, result2}
}

func (fake *FakeCNI) AddNetworkList(arg1 context.Context, arg2 *libcni.NetworkConfigList, arg3 *libcni.RuntimeConf) (types.Result, error) {
	fake.addNetworkListMutex.Lock()
	ret, specificReturn := fake.addNetworkListReturnsOnCall[len(fake.addNetworkListArgsForCall)]
	fake.addNetworkListArgsForCall = append(fake.addNetworkListArgsForCall, struct {
		arg1 context.Context
		arg2 *libcni.NetworkConfigList
		arg3 *libcni.RuntimeConf
	}{arg1, arg2, arg3})
	stub := fake.AddNetworkListStub
	fakeReturns := fake.addNetworkListReturns
	fake.recordInvocation("AddNetworkList", []interface{}{arg1, arg2, arg3})
	fake.addNetworkListMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeCNI) AddNetworkListCallCount() int {
	fake.addNetworkListMutex.RLock()
	defer fake.addNetworkListMutex.RUnlock()
	return len(fake.addNetworkListArgsForCall)
}

func (fake *FakeCNI) AddNetworkListCalls(stub func(context.Context, *libcni.NetworkConfigList, *libcni.RuntimeConf) (types.Result, error)) {
	fake.addNetworkListMutex.Lock()
	defer fake.addNetworkListMutex.Unlock()
	fake.AddNetworkListStub = stub
}

func (fake *FakeCNI) AddNetworkListArgsForCall(i int) (context.Context, *libcni.NetworkConfigList, *libcni.RuntimeConf) {
	fake.addNetworkListMutex.RLock()
	defer fake.addNetworkListMutex.RUnlock()
	argsForCall := fake.addNetworkListArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeCNI) AddNetworkListReturns(result1 types.Result, result2 error) {
	fake.addNetworkListMutex.Lock()
	defer fake.addNetworkListMutex.Unlock()
	fake.AddNetworkListStub = nil
	fake.addNetworkListReturns = struct {
		result1 types.Result
		result2 error
	}{result1, result2}
}

func (fake *FakeCNI) AddNetworkListReturnsOnCall(i int, result1 types.Result, result2 error) {
	fake.addNetworkListMutex.Lock()
	defer fake.addNetworkListMutex.Unlock()
	fake.AddNetworkListStub = nil
	if fake.addNetworkListReturnsOnCall == nil {
		fake.addNetworkListReturnsOnCall = make(map[int]struct {
			result1 types.Result
			result2 error
		})
	}
	fake.addNetworkListReturnsOnCall[i] = struct {
		result1 types.Result
		result2 error
	}{result1, result2}
}

func (fake *FakeCNI) CheckNetwork(arg1 context.Context, arg2 *libcni.NetworkConfig, arg3 *libcni.RuntimeConf) error {
	fake.checkNetworkMutex.Lock()
	ret, specificReturn := fake.checkNetworkReturnsOnCall[len(fake.checkNetworkArgsForCall)]
	fake.checkNetworkArgsForCall = append(fake.checkNetworkArgsForCall, struct {
		arg1 context.Context
		arg2 *libcni.NetworkConfig
		arg3 *libcni.RuntimeConf
	}{arg1, arg2, arg3})
	stub := fake.CheckNetworkStub
	fakeReturns := fake.checkNetworkReturns
	fake.recordInvocation("CheckNetwork", []interface{}{arg1, arg2, arg3})
	fake.checkNetworkMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeCNI) CheckNetworkCallCount() int {
	fake.checkNetworkMutex.RLock()
	defer fake.checkNetworkMutex.RUnlock()
	return len(fake.checkNetworkArgsForCall)
}

func (fake *FakeCNI) CheckNetworkCalls(stub func(context.Context, *libcni.NetworkConfig, *libcni.RuntimeConf) error) {
	fake.checkNetworkMutex.Lock()
	defer fake.checkNetworkMutex.Unlock()
	fake.CheckNetworkStub = stub
}

func (fake *FakeCNI) CheckNetworkArgsForCall(i int) (context.Context, *libcni.NetworkConfig, *libcni.RuntimeConf) {
	fake.checkNetworkMutex.RLock()
	defer fake.checkNetworkMutex.RUnlock()
	argsForCall := fake.checkNetworkArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeCNI) CheckNetworkReturns(result1 error) {
	fake.checkNetworkMutex.Lock()
	defer fake.checkNetworkMutex.Unlock()
	fake.CheckNetworkStub = nil
	fake.checkNetworkReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeCNI) CheckNetworkReturnsOnCall(i int, result1 error) {
	fake.checkNetworkMutex.Lock()
	defer fake.checkNetworkMutex.Unlock()
	fake.CheckNetworkStub = nil
	if fake.checkNetworkReturnsOnCall == nil {
		fake.checkNetworkReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.checkNetworkReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeCNI) CheckNetworkList(arg1 context.Context, arg2 *libcni.NetworkConfigList, arg3 *libcni.RuntimeConf) error {
	fake.checkNetworkListMutex.Lock()
	ret, specificReturn := fake.checkNetworkListReturnsOnCall[len(fake.checkNetworkListArgsForCall)]
	fake.checkNetworkListArgsForCall = append(fake.checkNetworkListArgsForCall, struct {
		arg1 context.Context
		arg2 *libcni.NetworkConfigList
		arg3 *libcni.RuntimeConf
	}{arg1, arg2, arg3})
	stub := fake.CheckNetworkListStub
	fakeReturns := fake.checkNetworkListReturns
	fake.recordInvocation("CheckNetworkList", []interface{}{arg1, arg2, arg3})
	fake.checkNetworkListMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeCNI) CheckNetworkListCallCount() int {
	fake.checkNetworkListMutex.RLock()
	defer fake.checkNetworkListMutex.RUnlock()
	return len(fake.checkNetworkListArgsForCall)
}

func (fake *FakeCNI) CheckNetworkListCalls(stub func(context.Context, *libcni.NetworkConfigList, *libcni.RuntimeConf) error) {
	fake.checkNetworkListMutex.Lock()
	defer fake.checkNetworkListMutex.Unlock()
	fake.CheckNetworkListStub = stub
}

func (fake *FakeCNI) CheckNetworkListArgsForCall(i int) (context.Context, *libcni.NetworkConfigList, *libcni.RuntimeConf) {
	fake.checkNetworkListMutex.RLock()
	defer fake.checkNetworkListMutex.RUnlock()
	argsForCall := fake.checkNetworkListArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeCNI) CheckNetworkListReturns(result1 error) {
	fake.checkNetworkListMutex.Lock()
	defer fake.checkNetworkListMutex.Unlock()
	fake.CheckNetworkListStub = nil
	fake.checkNetworkListReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeCNI) CheckNetworkListReturnsOnCall(i int, result1 error) {
	fake.checkNetworkListMutex.Lock()
	defer fake.checkNetworkListMutex.Unlock()
	fake.CheckNetworkListStub = nil
	if fake.checkNetworkListReturnsOnCall == nil {
		fake.checkNetworkListReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.checkNetworkListReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeCNI) DelNetwork(arg1 context.Context, arg2 *libcni.NetworkConfig, arg3 *libcni.RuntimeConf) error {
	fake.delNetworkMutex.Lock()
	ret, specificReturn := fake.delNetworkReturnsOnCall[len(fake.delNetworkArgsForCall)]
	fake.delNetworkArgsForCall = append(fake.delNetworkArgsForCall, struct {
		arg1 context.Context
		arg2 *libcni.NetworkConfig
		arg3 *libcni.RuntimeConf
	}{arg1, arg2, arg3})
	stub := fake.DelNetworkStub
	fakeReturns := fake.delNetworkReturns
	fake.recordInvocation("DelNetwork", []interface{}{arg1, arg2, arg3})
	fake.delNetworkMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeCNI) DelNetworkCallCount() int {
	fake.delNetworkMutex.RLock()
	defer fake.delNetworkMutex.RUnlock()
	return len(fake.delNetworkArgsForCall)
}

func (fake *FakeCNI) DelNetworkCalls(stub func(context.Context, *libcni.NetworkConfig, *libcni.RuntimeConf) error) {
	fake.delNetworkMutex.Lock()
	defer fake.delNetworkMutex.Unlock()
	fake.DelNetworkStub = stub
}

func (fake *FakeCNI) DelNetworkArgsForCall(i int) (context.Context, *libcni.NetworkConfig, *libcni.RuntimeConf) {
	fake.delNetworkMutex.RLock()
	defer fake.delNetworkMutex.RUnlock()
	argsForCall := fake.delNetworkArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeCNI) DelNetworkReturns(result1 error) {
	fake.delNetworkMutex.Lock()
	defer fake.delNetworkMutex.Unlock()
	fake.DelNetworkStub = nil
	fake.delNetworkReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeCNI) DelNetworkReturnsOnCall(i int, result1 error) {
	fake.delNetworkMutex.Lock()
	defer fake.delNetworkMutex.Unlock()
	fake.DelNetworkStub = nil
	if fake.delNetworkReturnsOnCall == nil {
		fake.delNetworkReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.delNetworkReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeCNI) DelNetworkList(arg1 context.Context, arg2 *libcni.NetworkConfigList, arg3 *libcni.RuntimeConf) error {
	fake.delNetworkListMutex.Lock()
	ret, specificReturn := fake.delNetworkListReturnsOnCall[len(fake.delNetworkListArgsForCall)]
	fake.delNetworkListArgsForCall = append(fake.delNetworkListArgsForCall, struct {
		arg1 context.Context
		arg2 *libcni.NetworkConfigList
		arg3 *libcni.RuntimeConf
	}{arg1, arg2, arg3})
	stub := fake.DelNetworkListStub
	fakeReturns := fake.delNetworkListReturns
	fake.recordInvocation("DelNetworkList", []interface{}{arg1, arg2, arg3})
	fake.delNetworkListMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeCNI) DelNetworkListCallCount() int {
	fake.delNetworkListMutex.RLock()
	defer fake.delNetworkListMutex.RUnlock()
	return len(fake.delNetworkListArgsForCall)
}

func (fake *FakeCNI) DelNetworkListCalls(stub func(context.Context, *libcni.NetworkConfigList, *libcni.RuntimeConf) error) {
	fake.delNetworkListMutex.Lock()
	defer fake.delNetworkListMutex.Unlock()
	fake.DelNetworkListStub = stub
}

func (fake *FakeCNI) DelNetworkListArgsForCall(i int) (context.Context, *libcni.NetworkConfigList, *libcni.RuntimeConf) {
	fake.delNetworkListMutex.RLock()
	defer fake.delNetworkListMutex.RUnlock()
	argsForCall := fake.delNetworkListArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeCNI) DelNetworkListReturns(result1 error) {
	fake.delNetworkListMutex.Lock()
	defer fake.delNetworkListMutex.Unlock()
	fake.DelNetworkListStub = nil
	fake.delNetworkListReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeCNI) DelNetworkListReturnsOnCall(i int, result1 error) {
	fake.delNetworkListMutex.Lock()
	defer fake.delNetworkListMutex.Unlock()
	fake.DelNetworkListStub = nil
	if fake.delNetworkListReturnsOnCall == nil {
		fake.delNetworkListReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.delNetworkListReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeCNI) GetNetworkCachedConfig(arg1 *libcni.NetworkConfig, arg2 *libcni.RuntimeConf) ([]byte, *libcni.RuntimeConf, error) {
	fake.getNetworkCachedConfigMutex.Lock()
	ret, specificReturn := fake.getNetworkCachedConfigReturnsOnCall[len(fake.getNetworkCachedConfigArgsForCall)]
	fake.getNetworkCachedConfigArgsForCall = append(fake.getNetworkCachedConfigArgsForCall, struct {
		arg1 *libcni.NetworkConfig
		arg2 *libcni.RuntimeConf
	}{arg1, arg2})
	stub := fake.GetNetworkCachedConfigStub
	fakeReturns := fake.getNetworkCachedConfigReturns
	fake.recordInvocation("GetNetworkCachedConfig", []interface{}{arg1, arg2})
	fake.getNetworkCachedConfigMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	return fakeReturns.result1, fakeReturns.result2, fakeReturns.result3
}

func (fake *FakeCNI) GetNetworkCachedConfigCallCount() int {
	fake.getNetworkCachedConfigMutex.RLock()
	defer fake.getNetworkCachedConfigMutex.RUnlock()
	return len(fake.getNetworkCachedConfigArgsForCall)
}

func (fake *FakeCNI) GetNetworkCachedConfigCalls(stub func(*libcni.NetworkConfig, *libcni.RuntimeConf) ([]byte, *libcni.RuntimeConf, error)) {
	fake.getNetworkCachedConfigMutex.Lock()
	defer fake.getNetworkCachedConfigMutex.Unlock()
	fake.GetNetworkCachedConfigStub = stub
}

func (fake *FakeCNI) GetNetworkCachedConfigArgsForCall(i int) (*libcni.NetworkConfig, *libcni.RuntimeConf) {
	fake.getNetworkCachedConfigMutex.RLock()
	defer fake.getNetworkCachedConfigMutex.RUnlock()
	argsForCall := fake.getNetworkCachedConfigArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeCNI) GetNetworkCachedConfigReturns(result1 []byte, result2 *libcni.RuntimeConf, result3 error) {
	fake.getNetworkCachedConfigMutex.Lock()
	defer fake.getNetworkCachedConfigMutex.Unlock()
	fake.GetNetworkCachedConfigStub = nil
	fake.getNetworkCachedConfigReturns = struct {
		result1 []byte
		result2 *libcni.RuntimeConf
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeCNI) GetNetworkCachedConfigReturnsOnCall(i int, result1 []byte, result2 *libcni.RuntimeConf, result3 error) {
	fake.getNetworkCachedConfigMutex.Lock()
	defer fake.getNetworkCachedConfigMutex.Unlock()
	fake.GetNetworkCachedConfigStub = nil
	if fake.getNetworkCachedConfigReturnsOnCall == nil {
		fake.getNetworkCachedConfigReturnsOnCall = make(map[int]struct {
			result1 []byte
			result2 *libcni.RuntimeConf
			result3 error
		})
	}
	fake.getNetworkCachedConfigReturnsOnCall[i] = struct {
		result1 []byte
		result2 *libcni.RuntimeConf
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeCNI) GetNetworkCachedResult(arg1 *libcni.NetworkConfig, arg2 *libcni.RuntimeConf) (types.Result, error) {
	fake.getNetworkCachedResultMutex.Lock()
	ret, specificReturn := fake.getNetworkCachedResultReturnsOnCall[len(fake.getNetworkCachedResultArgsForCall)]
	fake.getNetworkCachedResultArgsForCall = append(fake.getNetworkCachedResultArgsForCall, struct {
		arg1 *libcni.NetworkConfig
		arg2 *libcni.RuntimeConf
	}{arg1, arg2})
	stub := fake.GetNetworkCachedResultStub
	fakeReturns := fake.getNetworkCachedResultReturns
	fake.recordInvocation("GetNetworkCachedResult", []interface{}{arg1, arg2})
	fake.getNetworkCachedResultMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeCNI) GetNetworkCachedResultCallCount() int {
	fake.getNetworkCachedResultMutex.RLock()
	defer fake.getNetworkCachedResultMutex.RUnlock()
	return len(fake.getNetworkCachedResultArgsForCall)
}

func (fake *FakeCNI) GetNetworkCachedResultCalls(stub func(*libcni.NetworkConfig, *libcni.RuntimeConf) (types.Result, error)) {
	fake.getNetworkCachedResultMutex.Lock()
	defer fake.getNetworkCachedResultMutex.Unlock()
	fake.GetNetworkCachedResultStub = stub
}

func (fake *FakeCNI) GetNetworkCachedResultArgsForCall(i int) (*libcni.NetworkConfig, *libcni.RuntimeConf) {
	fake.getNetworkCachedResultMutex.RLock()
	defer fake.getNetworkCachedResultMutex.RUnlock()
	argsForCall := fake.getNetworkCachedResultArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeCNI) GetNetworkCachedResultReturns(result1 types.Result, result2 error) {
	fake.getNetworkCachedResultMutex.Lock()
	defer fake.getNetworkCachedResultMutex.Unlock()
	fake.GetNetworkCachedResultStub = nil
	fake.getNetworkCachedResultReturns = struct {
		result1 types.Result
		result2 error
	}{result1, result2}
}

func (fake *FakeCNI) GetNetworkCachedResultReturnsOnCall(i int, result1 types.Result, result2 error) {
	fake.getNetworkCachedResultMutex.Lock()
	defer fake.getNetworkCachedResultMutex.Unlock()
	fake.GetNetworkCachedResultStub = nil
	if fake.getNetworkCachedResultReturnsOnCall == nil {
		fake.getNetworkCachedResultReturnsOnCall = make(map[int]struct {
			result1 types.Result
			result2 error
		})
	}
	fake.getNetworkCachedResultReturnsOnCall[i] = struct {
		result1 types.Result
		result2 error
	}{result1, result2}
}

func (fake *FakeCNI) GetNetworkListCachedConfig(arg1 *libcni.NetworkConfigList, arg2 *libcni.RuntimeConf) ([]byte, *libcni.RuntimeConf, error) {
	fake.getNetworkListCachedConfigMutex.Lock()
	ret, specificReturn := fake.getNetworkListCachedConfigReturnsOnCall[len(fake.getNetworkListCachedConfigArgsForCall)]
	fake.getNetworkListCachedConfigArgsForCall = append(fake.getNetworkListCachedConfigArgsForCall, struct {
		arg1 *libcni.NetworkConfigList
		arg2 *libcni.RuntimeConf
	}{arg1, arg2})
	stub := fake.GetNetworkListCachedConfigStub
	fakeReturns := fake.getNetworkListCachedConfigReturns
	fake.recordInvocation("GetNetworkListCachedConfig", []interface{}{arg1, arg2})
	fake.getNetworkListCachedConfigMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	return fakeReturns.result1, fakeReturns.result2, fakeReturns.result3
}

func (fake *FakeCNI) GetNetworkListCachedConfigCallCount() int {
	fake.getNetworkListCachedConfigMutex.RLock()
	defer fake.getNetworkListCachedConfigMutex.RUnlock()
	return len(fake.getNetworkListCachedConfigArgsForCall)
}

func (fake *FakeCNI) GetNetworkListCachedConfigCalls(stub func(*libcni.NetworkConfigList, *libcni.RuntimeConf) ([]byte, *libcni.RuntimeConf, error)) {
	fake.getNetworkListCachedConfigMutex.Lock()
	defer fake.getNetworkListCachedConfigMutex.Unlock()
	fake.GetNetworkListCachedConfigStub = stub
}

func (fake *FakeCNI) GetNetworkListCachedConfigArgsForCall(i int) (*libcni.NetworkConfigList, *libcni.RuntimeConf) {
	fake.getNetworkListCachedConfigMutex.RLock()
	defer fake.getNetworkListCachedConfigMutex.RUnlock()
	argsForCall := fake.getNetworkListCachedConfigArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeCNI) GetNetworkListCachedConfigReturns(result1 []byte, result2 *libcni.RuntimeConf, result3 error) {
	fake.getNetworkListCachedConfigMutex.Lock()
	defer fake.getNetworkListCachedConfigMutex.Unlock()
	fake.GetNetworkListCachedConfigStub = nil
	fake.getNetworkListCachedConfigReturns = struct {
		result1 []byte
		result2 *libcni.RuntimeConf
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeCNI) GetNetworkListCachedConfigReturnsOnCall(i int, result1 []byte, result2 *libcni.RuntimeConf, result3 error) {
	fake.getNetworkListCachedConfigMutex.Lock()
	defer fake.getNetworkListCachedConfigMutex.Unlock()
	fake.GetNetworkListCachedConfigStub = nil
	if fake.getNetworkListCachedConfigReturnsOnCall == nil {
		fake.getNetworkListCachedConfigReturnsOnCall = make(map[int]struct {
			result1 []byte
			result2 *libcni.RuntimeConf
			result3 error
		})
	}
	fake.getNetworkListCachedConfigReturnsOnCall[i] = struct {
		result1 []byte
		result2 *libcni.RuntimeConf
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeCNI) GetNetworkListCachedResult(arg1 *libcni.NetworkConfigList, arg2 *libcni.RuntimeConf) (types.Result, error) {
	fake.getNetworkListCachedResultMutex.Lock()
	ret, specificReturn := fake.getNetworkListCachedResultReturnsOnCall[len(fake.getNetworkListCachedResultArgsForCall)]
	fake.getNetworkListCachedResultArgsForCall = append(fake.getNetworkListCachedResultArgsForCall, struct {
		arg1 *libcni.NetworkConfigList
		arg2 *libcni.RuntimeConf
	}{arg1, arg2})
	stub := fake.GetNetworkListCachedResultStub
	fakeReturns := fake.getNetworkListCachedResultReturns
	fake.recordInvocation("GetNetworkListCachedResult", []interface{}{arg1, arg2})
	fake.getNetworkListCachedResultMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeCNI) GetNetworkListCachedResultCallCount() int {
	fake.getNetworkListCachedResultMutex.RLock()
	defer fake.getNetworkListCachedResultMutex.RUnlock()
	return len(fake.getNetworkListCachedResultArgsForCall)
}

func (fake *FakeCNI) GetNetworkListCachedResultCalls(stub func(*libcni.NetworkConfigList, *libcni.RuntimeConf) (types.Result, error)) {
	fake.getNetworkListCachedResultMutex.Lock()
	defer fake.getNetworkListCachedResultMutex.Unlock()
	fake.GetNetworkListCachedResultStub = stub
}

func (fake *FakeCNI) GetNetworkListCachedResultArgsForCall(i int) (*libcni.NetworkConfigList, *libcni.RuntimeConf) {
	fake.getNetworkListCachedResultMutex.RLock()
	defer fake.getNetworkListCachedResultMutex.RUnlock()
	argsForCall := fake.getNetworkListCachedResultArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeCNI) GetNetworkListCachedResultReturns(result1 types.Result, result2 error) {
	fake.getNetworkListCachedResultMutex.Lock()
	defer fake.getNetworkListCachedResultMutex.Unlock()
	fake.GetNetworkListCachedResultStub = nil
	fake.getNetworkListCachedResultReturns = struct {
		result1 types.Result
		result2 error
	}{result1, result2}
}

func (fake *FakeCNI) GetNetworkListCachedResultReturnsOnCall(i int, result1 types.Result, result2 error) {
	fake.getNetworkListCachedResultMutex.Lock()
	defer fake.getNetworkListCachedResultMutex.Unlock()
	fake.GetNetworkListCachedResultStub = nil
	if fake.getNetworkListCachedResultReturnsOnCall == nil {
		fake.getNetworkListCachedResultReturnsOnCall = make(map[int]struct {
			result1 types.Result
			result2 error
		})
	}
	fake.getNetworkListCachedResultReturnsOnCall[i] = struct {
		result1 types.Result
		result2 error
	}{result1, result2}
}

func (fake *FakeCNI) ValidateNetwork(arg1 context.Context, arg2 *libcni.NetworkConfig) ([]string, error) {
	fake.validateNetworkMutex.Lock()
	ret, specificReturn := fake.validateNetworkReturnsOnCall[len(fake.validateNetworkArgsForCall)]
	fake.validateNetworkArgsForCall = append(fake.validateNetworkArgsForCall, struct {
		arg1 context.Context
		arg2 *libcni.NetworkConfig
	}{arg1, arg2})
	stub := fake.ValidateNetworkStub
	fakeReturns := fake.validateNetworkReturns
	fake.recordInvocation("ValidateNetwork", []interface{}{arg1, arg2})
	fake.validateNetworkMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeCNI) ValidateNetworkCallCount() int {
	fake.validateNetworkMutex.RLock()
	defer fake.validateNetworkMutex.RUnlock()
	return len(fake.validateNetworkArgsForCall)
}

func (fake *FakeCNI) ValidateNetworkCalls(stub func(context.Context, *libcni.NetworkConfig) ([]string, error)) {
	fake.validateNetworkMutex.Lock()
	defer fake.validateNetworkMutex.Unlock()
	fake.ValidateNetworkStub = stub
}

func (fake *FakeCNI) ValidateNetworkArgsForCall(i int) (context.Context, *libcni.NetworkConfig) {
	fake.validateNetworkMutex.RLock()
	defer fake.validateNetworkMutex.RUnlock()
	argsForCall := fake.validateNetworkArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeCNI) ValidateNetworkReturns(result1 []string, result2 error) {
	fake.validateNetworkMutex.Lock()
	defer fake.validateNetworkMutex.Unlock()
	fake.ValidateNetworkStub = nil
	fake.validateNetworkReturns = struct {
		result1 []string
		result2 error
	}{result1, result2}
}

func (fake *FakeCNI) ValidateNetworkReturnsOnCall(i int, result1 []string, result2 error) {
	fake.validateNetworkMutex.Lock()
	defer fake.validateNetworkMutex.Unlock()
	fake.ValidateNetworkStub = nil
	if fake.validateNetworkReturnsOnCall == nil {
		fake.validateNetworkReturnsOnCall = make(map[int]struct {
			result1 []string
			result2 error
		})
	}
	fake.validateNetworkReturnsOnCall[i] = struct {
		result1 []string
		result2 error
	}{result1, result2}
}

func (fake *FakeCNI) ValidateNetworkList(arg1 context.Context, arg2 *libcni.NetworkConfigList) ([]string, error) {
	fake.validateNetworkListMutex.Lock()
	ret, specificReturn := fake.validateNetworkListReturnsOnCall[len(fake.validateNetworkListArgsForCall)]
	fake.validateNetworkListArgsForCall = append(fake.validateNetworkListArgsForCall, struct {
		arg1 context.Context
		arg2 *libcni.NetworkConfigList
	}{arg1, arg2})
	stub := fake.ValidateNetworkListStub
	fakeReturns := fake.validateNetworkListReturns
	fake.recordInvocation("ValidateNetworkList", []interface{}{arg1, arg2})
	fake.validateNetworkListMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeCNI) ValidateNetworkListCallCount() int {
	fake.validateNetworkListMutex.RLock()
	defer fake.validateNetworkListMutex.RUnlock()
	return len(fake.validateNetworkListArgsForCall)
}

func (fake *FakeCNI) ValidateNetworkListCalls(stub func(context.Context, *libcni.NetworkConfigList) ([]string, error)) {
	fake.validateNetworkListMutex.Lock()
	defer fake.validateNetworkListMutex.Unlock()
	fake.ValidateNetworkListStub = stub
}

func (fake *FakeCNI) ValidateNetworkListArgsForCall(i int) (context.Context, *libcni.NetworkConfigList) {
	fake.validateNetworkListMutex.RLock()
	defer fake.validateNetworkListMutex.RUnlock()
	argsForCall := fake.validateNetworkListArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeCNI) ValidateNetworkListReturns(result1 []string, result2 error) {
	fake.validateNetworkListMutex.Lock()
	defer fake.validateNetworkListMutex.Unlock()
	fake.ValidateNetworkListStub = nil
	fake.validateNetworkListReturns = struct {
		result1 []string
		result2 error
	}{result1, result2}
}

func (fake *FakeCNI) ValidateNetworkListReturnsOnCall(i int, result1 []string, result2 error) {
	fake.validateNetworkListMutex.Lock()
	defer fake.validateNetworkListMutex.Unlock()
	fake.ValidateNetworkListStub = nil
	if fake.validateNetworkListReturnsOnCall == nil {
		fake.validateNetworkListReturnsOnCall = make(map[int]struct {
			result1 []string
			result2 error
		})
	}
	fake.validateNetworkListReturnsOnCall[i] = struct {
		result1 []string
		result2 error
	}{result1, result2}
}

func (fake *FakeCNI) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.addNetworkMutex.RLock()
	defer fake.addNetworkMutex.RUnlock()
	fake.addNetworkListMutex.RLock()
	defer fake.addNetworkListMutex.RUnlock()
	fake.checkNetworkMutex.RLock()
	defer fake.checkNetworkMutex.RUnlock()
	fake.checkNetworkListMutex.RLock()
	defer fake.checkNetworkListMutex.RUnlock()
	fake.delNetworkMutex.RLock()
	defer fake.delNetworkMutex.RUnlock()
	fake.delNetworkListMutex.RLock()
	defer fake.delNetworkListMutex.RUnlock()
	fake.getNetworkCachedConfigMutex.RLock()
	defer fake.getNetworkCachedConfigMutex.RUnlock()
	fake.getNetworkCachedResultMutex.RLock()
	defer fake.getNetworkCachedResultMutex.RUnlock()
	fake.getNetworkListCachedConfigMutex.RLock()
	defer fake.getNetworkListCachedConfigMutex.RUnlock()
	fake.getNetworkListCachedResultMutex.RLock()
	defer fake.getNetworkListCachedResultMutex.RUnlock()
	fake.validateNetworkMutex.RLock()
	defer fake.validateNetworkMutex.RUnlock()
	fake.validateNetworkListMutex.RLock()
	defer fake.validateNetworkListMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeCNI) recordInvocation(key string, args []interface{}) {
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

var _ libcni.CNI = new(FakeCNI)
