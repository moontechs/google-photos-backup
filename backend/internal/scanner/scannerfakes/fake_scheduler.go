// Code generated by counterfeiter. DO NOT EDIT.
package scannerfakes

import (
	"google-backup/internal/scanner"
	"sync"
)

type FakeScheduler struct {
	ScheduleRescanStub        func(string, string) error
	scheduleRescanMutex       sync.RWMutex
	scheduleRescanArgsForCall []struct {
		arg1 string
		arg2 string
	}
	scheduleRescanReturns struct {
		result1 error
	}
	scheduleRescanReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeScheduler) ScheduleRescan(arg1 string, arg2 string) error {
	fake.scheduleRescanMutex.Lock()
	ret, specificReturn := fake.scheduleRescanReturnsOnCall[len(fake.scheduleRescanArgsForCall)]
	fake.scheduleRescanArgsForCall = append(fake.scheduleRescanArgsForCall, struct {
		arg1 string
		arg2 string
	}{arg1, arg2})
	stub := fake.ScheduleRescanStub
	fakeReturns := fake.scheduleRescanReturns
	fake.recordInvocation("ScheduleRescan", []interface{}{arg1, arg2})
	fake.scheduleRescanMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeScheduler) ScheduleRescanCallCount() int {
	fake.scheduleRescanMutex.RLock()
	defer fake.scheduleRescanMutex.RUnlock()
	return len(fake.scheduleRescanArgsForCall)
}

func (fake *FakeScheduler) ScheduleRescanCalls(stub func(string, string) error) {
	fake.scheduleRescanMutex.Lock()
	defer fake.scheduleRescanMutex.Unlock()
	fake.ScheduleRescanStub = stub
}

func (fake *FakeScheduler) ScheduleRescanArgsForCall(i int) (string, string) {
	fake.scheduleRescanMutex.RLock()
	defer fake.scheduleRescanMutex.RUnlock()
	argsForCall := fake.scheduleRescanArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeScheduler) ScheduleRescanReturns(result1 error) {
	fake.scheduleRescanMutex.Lock()
	defer fake.scheduleRescanMutex.Unlock()
	fake.ScheduleRescanStub = nil
	fake.scheduleRescanReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeScheduler) ScheduleRescanReturnsOnCall(i int, result1 error) {
	fake.scheduleRescanMutex.Lock()
	defer fake.scheduleRescanMutex.Unlock()
	fake.ScheduleRescanStub = nil
	if fake.scheduleRescanReturnsOnCall == nil {
		fake.scheduleRescanReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.scheduleRescanReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeScheduler) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.scheduleRescanMutex.RLock()
	defer fake.scheduleRescanMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeScheduler) recordInvocation(key string, args []interface{}) {
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

var _ scanner.Scheduler = new(FakeScheduler)