// Code generated by counterfeiter. DO NOT EDIT.
package iaas_clientsfakes

import (
	sync "sync"

	iaas_clients "github.com/cloudfoundry-incubator/stembuild/package_stemcell/iaas_clients"
)

type FakeVcenterClient struct {
	FindVMStub        func(string) error
	findVMMutex       sync.RWMutex
	findVMArgsForCall []struct {
		arg1 string
	}
	findVMReturns struct {
		result1 error
	}
	findVMReturnsOnCall map[int]struct {
		result1 error
	}
	LoginStub        func() error
	loginMutex       sync.RWMutex
	loginArgsForCall []struct {
	}
	loginReturns struct {
		result1 error
	}
	loginReturnsOnCall map[int]struct {
		result1 error
	}
	ValidateUrlStub        func() error
	validateUrlMutex       sync.RWMutex
	validateUrlArgsForCall []struct {
	}
	validateUrlReturns struct {
		result1 error
	}
	validateUrlReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeVcenterClient) FindVM(arg1 string) error {
	fake.findVMMutex.Lock()
	ret, specificReturn := fake.findVMReturnsOnCall[len(fake.findVMArgsForCall)]
	fake.findVMArgsForCall = append(fake.findVMArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("FindVM", []interface{}{arg1})
	fake.findVMMutex.Unlock()
	if fake.FindVMStub != nil {
		return fake.FindVMStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.findVMReturns
	return fakeReturns.result1
}

func (fake *FakeVcenterClient) FindVMCallCount() int {
	fake.findVMMutex.RLock()
	defer fake.findVMMutex.RUnlock()
	return len(fake.findVMArgsForCall)
}

func (fake *FakeVcenterClient) FindVMCalls(stub func(string) error) {
	fake.findVMMutex.Lock()
	defer fake.findVMMutex.Unlock()
	fake.FindVMStub = stub
}

func (fake *FakeVcenterClient) FindVMArgsForCall(i int) string {
	fake.findVMMutex.RLock()
	defer fake.findVMMutex.RUnlock()
	argsForCall := fake.findVMArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeVcenterClient) FindVMReturns(result1 error) {
	fake.findVMMutex.Lock()
	defer fake.findVMMutex.Unlock()
	fake.FindVMStub = nil
	fake.findVMReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeVcenterClient) FindVMReturnsOnCall(i int, result1 error) {
	fake.findVMMutex.Lock()
	defer fake.findVMMutex.Unlock()
	fake.FindVMStub = nil
	if fake.findVMReturnsOnCall == nil {
		fake.findVMReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.findVMReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeVcenterClient) Login() error {
	fake.loginMutex.Lock()
	ret, specificReturn := fake.loginReturnsOnCall[len(fake.loginArgsForCall)]
	fake.loginArgsForCall = append(fake.loginArgsForCall, struct {
	}{})
	fake.recordInvocation("Login", []interface{}{})
	fake.loginMutex.Unlock()
	if fake.LoginStub != nil {
		return fake.LoginStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.loginReturns
	return fakeReturns.result1
}

func (fake *FakeVcenterClient) LoginCallCount() int {
	fake.loginMutex.RLock()
	defer fake.loginMutex.RUnlock()
	return len(fake.loginArgsForCall)
}

func (fake *FakeVcenterClient) LoginCalls(stub func() error) {
	fake.loginMutex.Lock()
	defer fake.loginMutex.Unlock()
	fake.LoginStub = stub
}

func (fake *FakeVcenterClient) LoginReturns(result1 error) {
	fake.loginMutex.Lock()
	defer fake.loginMutex.Unlock()
	fake.LoginStub = nil
	fake.loginReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeVcenterClient) LoginReturnsOnCall(i int, result1 error) {
	fake.loginMutex.Lock()
	defer fake.loginMutex.Unlock()
	fake.LoginStub = nil
	if fake.loginReturnsOnCall == nil {
		fake.loginReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.loginReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeVcenterClient) ValidateUrl() error {
	fake.validateUrlMutex.Lock()
	ret, specificReturn := fake.validateUrlReturnsOnCall[len(fake.validateUrlArgsForCall)]
	fake.validateUrlArgsForCall = append(fake.validateUrlArgsForCall, struct {
	}{})
	fake.recordInvocation("ValidateUrl", []interface{}{})
	fake.validateUrlMutex.Unlock()
	if fake.ValidateUrlStub != nil {
		return fake.ValidateUrlStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.validateUrlReturns
	return fakeReturns.result1
}

func (fake *FakeVcenterClient) ValidateUrlCallCount() int {
	fake.validateUrlMutex.RLock()
	defer fake.validateUrlMutex.RUnlock()
	return len(fake.validateUrlArgsForCall)
}

func (fake *FakeVcenterClient) ValidateUrlCalls(stub func() error) {
	fake.validateUrlMutex.Lock()
	defer fake.validateUrlMutex.Unlock()
	fake.ValidateUrlStub = stub
}

func (fake *FakeVcenterClient) ValidateUrlReturns(result1 error) {
	fake.validateUrlMutex.Lock()
	defer fake.validateUrlMutex.Unlock()
	fake.ValidateUrlStub = nil
	fake.validateUrlReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeVcenterClient) ValidateUrlReturnsOnCall(i int, result1 error) {
	fake.validateUrlMutex.Lock()
	defer fake.validateUrlMutex.Unlock()
	fake.ValidateUrlStub = nil
	if fake.validateUrlReturnsOnCall == nil {
		fake.validateUrlReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.validateUrlReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeVcenterClient) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.findVMMutex.RLock()
	defer fake.findVMMutex.RUnlock()
	fake.loginMutex.RLock()
	defer fake.loginMutex.RUnlock()
	fake.validateUrlMutex.RLock()
	defer fake.validateUrlMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeVcenterClient) recordInvocation(key string, args []interface{}) {
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

var _ iaas_clients.VcenterClient = new(FakeVcenterClient)