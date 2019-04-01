// Code generated by counterfeiter. DO NOT EDIT.
package constructfakes

import (
	"sync"
)

type FakeZipUnarchiver struct {
	UnzipStub        func([]byte, string) ([]byte, error)
	unzipMutex       sync.RWMutex
	unzipArgsForCall []struct {
		arg1 []byte
		arg2 string
	}
	unzipReturns struct {
		result1 []byte
		result2 error
	}
	unzipReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeZipUnarchiver) Unzip(arg1 []byte, arg2 string) ([]byte, error) {
	var arg1Copy []byte
	if arg1 != nil {
		arg1Copy = make([]byte, len(arg1))
		copy(arg1Copy, arg1)
	}
	fake.unzipMutex.Lock()
	ret, specificReturn := fake.unzipReturnsOnCall[len(fake.unzipArgsForCall)]
	fake.unzipArgsForCall = append(fake.unzipArgsForCall, struct {
		arg1 []byte
		arg2 string
	}{arg1Copy, arg2})
	fake.recordInvocation("Unzip", []interface{}{arg1Copy, arg2})
	fake.unzipMutex.Unlock()
	if fake.UnzipStub != nil {
		return fake.UnzipStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.unzipReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeZipUnarchiver) UnzipCallCount() int {
	fake.unzipMutex.RLock()
	defer fake.unzipMutex.RUnlock()
	return len(fake.unzipArgsForCall)
}

func (fake *FakeZipUnarchiver) UnzipCalls(stub func([]byte, string) ([]byte, error)) {
	fake.unzipMutex.Lock()
	defer fake.unzipMutex.Unlock()
	fake.UnzipStub = stub
}

func (fake *FakeZipUnarchiver) UnzipArgsForCall(i int) ([]byte, string) {
	fake.unzipMutex.RLock()
	defer fake.unzipMutex.RUnlock()
	argsForCall := fake.unzipArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeZipUnarchiver) UnzipReturns(result1 []byte, result2 error) {
	fake.unzipMutex.Lock()
	defer fake.unzipMutex.Unlock()
	fake.UnzipStub = nil
	fake.unzipReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *FakeZipUnarchiver) UnzipReturnsOnCall(i int, result1 []byte, result2 error) {
	fake.unzipMutex.Lock()
	defer fake.unzipMutex.Unlock()
	fake.UnzipStub = nil
	if fake.unzipReturnsOnCall == nil {
		fake.unzipReturnsOnCall = make(map[int]struct {
			result1 []byte
			result2 error
		})
	}
	fake.unzipReturnsOnCall[i] = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *FakeZipUnarchiver) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.unzipMutex.RLock()
	defer fake.unzipMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeZipUnarchiver) recordInvocation(key string, args []interface{}) {
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
