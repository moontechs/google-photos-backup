// Code generated by counterfeiter. DO NOT EDIT.
package settingsfakes

import (
	"google-backup/internal/settings"
	"sync"
)

type FakeRepository struct {
	FindStub        func() ([]byte, error)
	findMutex       sync.RWMutex
	findArgsForCall []struct {
	}
	findReturns struct {
		result1 []byte
		result2 error
	}
	findReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	SaveStub        func([]byte) error
	saveMutex       sync.RWMutex
	saveArgsForCall []struct {
		arg1 []byte
	}
	saveReturns struct {
		result1 error
	}
	saveReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeRepository) Find() ([]byte, error) {
	fake.findMutex.Lock()
	ret, specificReturn := fake.findReturnsOnCall[len(fake.findArgsForCall)]
	fake.findArgsForCall = append(fake.findArgsForCall, struct {
	}{})
	stub := fake.FindStub
	fakeReturns := fake.findReturns
	fake.recordInvocation("Find", []interface{}{})
	fake.findMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeRepository) FindCallCount() int {
	fake.findMutex.RLock()
	defer fake.findMutex.RUnlock()
	return len(fake.findArgsForCall)
}

func (fake *FakeRepository) FindCalls(stub func() ([]byte, error)) {
	fake.findMutex.Lock()
	defer fake.findMutex.Unlock()
	fake.FindStub = stub
}

func (fake *FakeRepository) FindReturns(result1 []byte, result2 error) {
	fake.findMutex.Lock()
	defer fake.findMutex.Unlock()
	fake.FindStub = nil
	fake.findReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *FakeRepository) FindReturnsOnCall(i int, result1 []byte, result2 error) {
	fake.findMutex.Lock()
	defer fake.findMutex.Unlock()
	fake.FindStub = nil
	if fake.findReturnsOnCall == nil {
		fake.findReturnsOnCall = make(map[int]struct {
			result1 []byte
			result2 error
		})
	}
	fake.findReturnsOnCall[i] = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *FakeRepository) Save(arg1 []byte) error {
	var arg1Copy []byte
	if arg1 != nil {
		arg1Copy = make([]byte, len(arg1))
		copy(arg1Copy, arg1)
	}
	fake.saveMutex.Lock()
	ret, specificReturn := fake.saveReturnsOnCall[len(fake.saveArgsForCall)]
	fake.saveArgsForCall = append(fake.saveArgsForCall, struct {
		arg1 []byte
	}{arg1Copy})
	stub := fake.SaveStub
	fakeReturns := fake.saveReturns
	fake.recordInvocation("Save", []interface{}{arg1Copy})
	fake.saveMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeRepository) SaveCallCount() int {
	fake.saveMutex.RLock()
	defer fake.saveMutex.RUnlock()
	return len(fake.saveArgsForCall)
}

func (fake *FakeRepository) SaveCalls(stub func([]byte) error) {
	fake.saveMutex.Lock()
	defer fake.saveMutex.Unlock()
	fake.SaveStub = stub
}

func (fake *FakeRepository) SaveArgsForCall(i int) []byte {
	fake.saveMutex.RLock()
	defer fake.saveMutex.RUnlock()
	argsForCall := fake.saveArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeRepository) SaveReturns(result1 error) {
	fake.saveMutex.Lock()
	defer fake.saveMutex.Unlock()
	fake.SaveStub = nil
	fake.saveReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeRepository) SaveReturnsOnCall(i int, result1 error) {
	fake.saveMutex.Lock()
	defer fake.saveMutex.Unlock()
	fake.SaveStub = nil
	if fake.saveReturnsOnCall == nil {
		fake.saveReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.saveReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeRepository) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.findMutex.RLock()
	defer fake.findMutex.RUnlock()
	fake.saveMutex.RLock()
	defer fake.saveMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeRepository) recordInvocation(key string, args []interface{}) {
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

var _ settings.Repository = new(FakeRepository)
