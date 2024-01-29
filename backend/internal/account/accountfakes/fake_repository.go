// Code generated by counterfeiter. DO NOT EDIT.
package accountfakes

import (
	"google-backup/internal/account"
	"sync"
)

type FakeRepository struct {
	AccountExistStub        func(string) (bool, error)
	accountExistMutex       sync.RWMutex
	accountExistArgsForCall []struct {
		arg1 string
	}
	accountExistReturns struct {
		result1 bool
		result2 error
	}
	accountExistReturnsOnCall map[int]struct {
		result1 bool
		result2 error
	}
	CreateUpdateLimitsStub        func(string, []byte) error
	createUpdateLimitsMutex       sync.RWMutex
	createUpdateLimitsArgsForCall []struct {
		arg1 string
		arg2 []byte
	}
	createUpdateLimitsReturns struct {
		result1 error
	}
	createUpdateLimitsReturnsOnCall map[int]struct {
		result1 error
	}
	FindAccountStub        func(string) ([]byte, error)
	findAccountMutex       sync.RWMutex
	findAccountArgsForCall []struct {
		arg1 string
	}
	findAccountReturns struct {
		result1 []byte
		result2 error
	}
	findAccountReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	FindTokenByEmailStub        func(string) ([]byte, error)
	findTokenByEmailMutex       sync.RWMutex
	findTokenByEmailArgsForCall []struct {
		arg1 string
	}
	findTokenByEmailReturns struct {
		result1 []byte
		result2 error
	}
	findTokenByEmailReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	GetAccountOauthClientNameStub        func(string) ([]byte, error)
	getAccountOauthClientNameMutex       sync.RWMutex
	getAccountOauthClientNameArgsForCall []struct {
		arg1 string
	}
	getAccountOauthClientNameReturns struct {
		result1 []byte
		result2 error
	}
	getAccountOauthClientNameReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	GetAccountsStub        func() ([][]byte, error)
	getAccountsMutex       sync.RWMutex
	getAccountsArgsForCall []struct {
	}
	getAccountsReturns struct {
		result1 [][]byte
		result2 error
	}
	getAccountsReturnsOnCall map[int]struct {
		result1 [][]byte
		result2 error
	}
	GetLimitsStub        func(string) ([]byte, error)
	getLimitsMutex       sync.RWMutex
	getLimitsArgsForCall []struct {
		arg1 string
	}
	getLimitsReturns struct {
		result1 []byte
		result2 error
	}
	getLimitsReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	SaveAccountStub        func(string, []byte) error
	saveAccountMutex       sync.RWMutex
	saveAccountArgsForCall []struct {
		arg1 string
		arg2 []byte
	}
	saveAccountReturns struct {
		result1 error
	}
	saveAccountReturnsOnCall map[int]struct {
		result1 error
	}
	SaveTokenStub        func(string, []byte) error
	saveTokenMutex       sync.RWMutex
	saveTokenArgsForCall []struct {
		arg1 string
		arg2 []byte
	}
	saveTokenReturns struct {
		result1 error
	}
	saveTokenReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeRepository) AccountExist(arg1 string) (bool, error) {
	fake.accountExistMutex.Lock()
	ret, specificReturn := fake.accountExistReturnsOnCall[len(fake.accountExistArgsForCall)]
	fake.accountExistArgsForCall = append(fake.accountExistArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.AccountExistStub
	fakeReturns := fake.accountExistReturns
	fake.recordInvocation("AccountExist", []interface{}{arg1})
	fake.accountExistMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeRepository) AccountExistCallCount() int {
	fake.accountExistMutex.RLock()
	defer fake.accountExistMutex.RUnlock()
	return len(fake.accountExistArgsForCall)
}

func (fake *FakeRepository) AccountExistCalls(stub func(string) (bool, error)) {
	fake.accountExistMutex.Lock()
	defer fake.accountExistMutex.Unlock()
	fake.AccountExistStub = stub
}

func (fake *FakeRepository) AccountExistArgsForCall(i int) string {
	fake.accountExistMutex.RLock()
	defer fake.accountExistMutex.RUnlock()
	argsForCall := fake.accountExistArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeRepository) AccountExistReturns(result1 bool, result2 error) {
	fake.accountExistMutex.Lock()
	defer fake.accountExistMutex.Unlock()
	fake.AccountExistStub = nil
	fake.accountExistReturns = struct {
		result1 bool
		result2 error
	}{result1, result2}
}

func (fake *FakeRepository) AccountExistReturnsOnCall(i int, result1 bool, result2 error) {
	fake.accountExistMutex.Lock()
	defer fake.accountExistMutex.Unlock()
	fake.AccountExistStub = nil
	if fake.accountExistReturnsOnCall == nil {
		fake.accountExistReturnsOnCall = make(map[int]struct {
			result1 bool
			result2 error
		})
	}
	fake.accountExistReturnsOnCall[i] = struct {
		result1 bool
		result2 error
	}{result1, result2}
}

func (fake *FakeRepository) CreateUpdateLimits(arg1 string, arg2 []byte) error {
	var arg2Copy []byte
	if arg2 != nil {
		arg2Copy = make([]byte, len(arg2))
		copy(arg2Copy, arg2)
	}
	fake.createUpdateLimitsMutex.Lock()
	ret, specificReturn := fake.createUpdateLimitsReturnsOnCall[len(fake.createUpdateLimitsArgsForCall)]
	fake.createUpdateLimitsArgsForCall = append(fake.createUpdateLimitsArgsForCall, struct {
		arg1 string
		arg2 []byte
	}{arg1, arg2Copy})
	stub := fake.CreateUpdateLimitsStub
	fakeReturns := fake.createUpdateLimitsReturns
	fake.recordInvocation("CreateUpdateLimits", []interface{}{arg1, arg2Copy})
	fake.createUpdateLimitsMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeRepository) CreateUpdateLimitsCallCount() int {
	fake.createUpdateLimitsMutex.RLock()
	defer fake.createUpdateLimitsMutex.RUnlock()
	return len(fake.createUpdateLimitsArgsForCall)
}

func (fake *FakeRepository) CreateUpdateLimitsCalls(stub func(string, []byte) error) {
	fake.createUpdateLimitsMutex.Lock()
	defer fake.createUpdateLimitsMutex.Unlock()
	fake.CreateUpdateLimitsStub = stub
}

func (fake *FakeRepository) CreateUpdateLimitsArgsForCall(i int) (string, []byte) {
	fake.createUpdateLimitsMutex.RLock()
	defer fake.createUpdateLimitsMutex.RUnlock()
	argsForCall := fake.createUpdateLimitsArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeRepository) CreateUpdateLimitsReturns(result1 error) {
	fake.createUpdateLimitsMutex.Lock()
	defer fake.createUpdateLimitsMutex.Unlock()
	fake.CreateUpdateLimitsStub = nil
	fake.createUpdateLimitsReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeRepository) CreateUpdateLimitsReturnsOnCall(i int, result1 error) {
	fake.createUpdateLimitsMutex.Lock()
	defer fake.createUpdateLimitsMutex.Unlock()
	fake.CreateUpdateLimitsStub = nil
	if fake.createUpdateLimitsReturnsOnCall == nil {
		fake.createUpdateLimitsReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.createUpdateLimitsReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeRepository) FindAccount(arg1 string) ([]byte, error) {
	fake.findAccountMutex.Lock()
	ret, specificReturn := fake.findAccountReturnsOnCall[len(fake.findAccountArgsForCall)]
	fake.findAccountArgsForCall = append(fake.findAccountArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.FindAccountStub
	fakeReturns := fake.findAccountReturns
	fake.recordInvocation("FindAccount", []interface{}{arg1})
	fake.findAccountMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeRepository) FindAccountCallCount() int {
	fake.findAccountMutex.RLock()
	defer fake.findAccountMutex.RUnlock()
	return len(fake.findAccountArgsForCall)
}

func (fake *FakeRepository) FindAccountCalls(stub func(string) ([]byte, error)) {
	fake.findAccountMutex.Lock()
	defer fake.findAccountMutex.Unlock()
	fake.FindAccountStub = stub
}

func (fake *FakeRepository) FindAccountArgsForCall(i int) string {
	fake.findAccountMutex.RLock()
	defer fake.findAccountMutex.RUnlock()
	argsForCall := fake.findAccountArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeRepository) FindAccountReturns(result1 []byte, result2 error) {
	fake.findAccountMutex.Lock()
	defer fake.findAccountMutex.Unlock()
	fake.FindAccountStub = nil
	fake.findAccountReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *FakeRepository) FindAccountReturnsOnCall(i int, result1 []byte, result2 error) {
	fake.findAccountMutex.Lock()
	defer fake.findAccountMutex.Unlock()
	fake.FindAccountStub = nil
	if fake.findAccountReturnsOnCall == nil {
		fake.findAccountReturnsOnCall = make(map[int]struct {
			result1 []byte
			result2 error
		})
	}
	fake.findAccountReturnsOnCall[i] = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *FakeRepository) FindTokenByEmail(arg1 string) ([]byte, error) {
	fake.findTokenByEmailMutex.Lock()
	ret, specificReturn := fake.findTokenByEmailReturnsOnCall[len(fake.findTokenByEmailArgsForCall)]
	fake.findTokenByEmailArgsForCall = append(fake.findTokenByEmailArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.FindTokenByEmailStub
	fakeReturns := fake.findTokenByEmailReturns
	fake.recordInvocation("FindTokenByEmail", []interface{}{arg1})
	fake.findTokenByEmailMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeRepository) FindTokenByEmailCallCount() int {
	fake.findTokenByEmailMutex.RLock()
	defer fake.findTokenByEmailMutex.RUnlock()
	return len(fake.findTokenByEmailArgsForCall)
}

func (fake *FakeRepository) FindTokenByEmailCalls(stub func(string) ([]byte, error)) {
	fake.findTokenByEmailMutex.Lock()
	defer fake.findTokenByEmailMutex.Unlock()
	fake.FindTokenByEmailStub = stub
}

func (fake *FakeRepository) FindTokenByEmailArgsForCall(i int) string {
	fake.findTokenByEmailMutex.RLock()
	defer fake.findTokenByEmailMutex.RUnlock()
	argsForCall := fake.findTokenByEmailArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeRepository) FindTokenByEmailReturns(result1 []byte, result2 error) {
	fake.findTokenByEmailMutex.Lock()
	defer fake.findTokenByEmailMutex.Unlock()
	fake.FindTokenByEmailStub = nil
	fake.findTokenByEmailReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *FakeRepository) FindTokenByEmailReturnsOnCall(i int, result1 []byte, result2 error) {
	fake.findTokenByEmailMutex.Lock()
	defer fake.findTokenByEmailMutex.Unlock()
	fake.FindTokenByEmailStub = nil
	if fake.findTokenByEmailReturnsOnCall == nil {
		fake.findTokenByEmailReturnsOnCall = make(map[int]struct {
			result1 []byte
			result2 error
		})
	}
	fake.findTokenByEmailReturnsOnCall[i] = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *FakeRepository) GetAccountOauthClientName(arg1 string) ([]byte, error) {
	fake.getAccountOauthClientNameMutex.Lock()
	ret, specificReturn := fake.getAccountOauthClientNameReturnsOnCall[len(fake.getAccountOauthClientNameArgsForCall)]
	fake.getAccountOauthClientNameArgsForCall = append(fake.getAccountOauthClientNameArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.GetAccountOauthClientNameStub
	fakeReturns := fake.getAccountOauthClientNameReturns
	fake.recordInvocation("GetAccountOauthClientName", []interface{}{arg1})
	fake.getAccountOauthClientNameMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeRepository) GetAccountOauthClientNameCallCount() int {
	fake.getAccountOauthClientNameMutex.RLock()
	defer fake.getAccountOauthClientNameMutex.RUnlock()
	return len(fake.getAccountOauthClientNameArgsForCall)
}

func (fake *FakeRepository) GetAccountOauthClientNameCalls(stub func(string) ([]byte, error)) {
	fake.getAccountOauthClientNameMutex.Lock()
	defer fake.getAccountOauthClientNameMutex.Unlock()
	fake.GetAccountOauthClientNameStub = stub
}

func (fake *FakeRepository) GetAccountOauthClientNameArgsForCall(i int) string {
	fake.getAccountOauthClientNameMutex.RLock()
	defer fake.getAccountOauthClientNameMutex.RUnlock()
	argsForCall := fake.getAccountOauthClientNameArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeRepository) GetAccountOauthClientNameReturns(result1 []byte, result2 error) {
	fake.getAccountOauthClientNameMutex.Lock()
	defer fake.getAccountOauthClientNameMutex.Unlock()
	fake.GetAccountOauthClientNameStub = nil
	fake.getAccountOauthClientNameReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *FakeRepository) GetAccountOauthClientNameReturnsOnCall(i int, result1 []byte, result2 error) {
	fake.getAccountOauthClientNameMutex.Lock()
	defer fake.getAccountOauthClientNameMutex.Unlock()
	fake.GetAccountOauthClientNameStub = nil
	if fake.getAccountOauthClientNameReturnsOnCall == nil {
		fake.getAccountOauthClientNameReturnsOnCall = make(map[int]struct {
			result1 []byte
			result2 error
		})
	}
	fake.getAccountOauthClientNameReturnsOnCall[i] = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *FakeRepository) GetAccounts() ([][]byte, error) {
	fake.getAccountsMutex.Lock()
	ret, specificReturn := fake.getAccountsReturnsOnCall[len(fake.getAccountsArgsForCall)]
	fake.getAccountsArgsForCall = append(fake.getAccountsArgsForCall, struct {
	}{})
	stub := fake.GetAccountsStub
	fakeReturns := fake.getAccountsReturns
	fake.recordInvocation("GetAccounts", []interface{}{})
	fake.getAccountsMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeRepository) GetAccountsCallCount() int {
	fake.getAccountsMutex.RLock()
	defer fake.getAccountsMutex.RUnlock()
	return len(fake.getAccountsArgsForCall)
}

func (fake *FakeRepository) GetAccountsCalls(stub func() ([][]byte, error)) {
	fake.getAccountsMutex.Lock()
	defer fake.getAccountsMutex.Unlock()
	fake.GetAccountsStub = stub
}

func (fake *FakeRepository) GetAccountsReturns(result1 [][]byte, result2 error) {
	fake.getAccountsMutex.Lock()
	defer fake.getAccountsMutex.Unlock()
	fake.GetAccountsStub = nil
	fake.getAccountsReturns = struct {
		result1 [][]byte
		result2 error
	}{result1, result2}
}

func (fake *FakeRepository) GetAccountsReturnsOnCall(i int, result1 [][]byte, result2 error) {
	fake.getAccountsMutex.Lock()
	defer fake.getAccountsMutex.Unlock()
	fake.GetAccountsStub = nil
	if fake.getAccountsReturnsOnCall == nil {
		fake.getAccountsReturnsOnCall = make(map[int]struct {
			result1 [][]byte
			result2 error
		})
	}
	fake.getAccountsReturnsOnCall[i] = struct {
		result1 [][]byte
		result2 error
	}{result1, result2}
}

func (fake *FakeRepository) GetLimits(arg1 string) ([]byte, error) {
	fake.getLimitsMutex.Lock()
	ret, specificReturn := fake.getLimitsReturnsOnCall[len(fake.getLimitsArgsForCall)]
	fake.getLimitsArgsForCall = append(fake.getLimitsArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.GetLimitsStub
	fakeReturns := fake.getLimitsReturns
	fake.recordInvocation("GetLimits", []interface{}{arg1})
	fake.getLimitsMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeRepository) GetLimitsCallCount() int {
	fake.getLimitsMutex.RLock()
	defer fake.getLimitsMutex.RUnlock()
	return len(fake.getLimitsArgsForCall)
}

func (fake *FakeRepository) GetLimitsCalls(stub func(string) ([]byte, error)) {
	fake.getLimitsMutex.Lock()
	defer fake.getLimitsMutex.Unlock()
	fake.GetLimitsStub = stub
}

func (fake *FakeRepository) GetLimitsArgsForCall(i int) string {
	fake.getLimitsMutex.RLock()
	defer fake.getLimitsMutex.RUnlock()
	argsForCall := fake.getLimitsArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeRepository) GetLimitsReturns(result1 []byte, result2 error) {
	fake.getLimitsMutex.Lock()
	defer fake.getLimitsMutex.Unlock()
	fake.GetLimitsStub = nil
	fake.getLimitsReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *FakeRepository) GetLimitsReturnsOnCall(i int, result1 []byte, result2 error) {
	fake.getLimitsMutex.Lock()
	defer fake.getLimitsMutex.Unlock()
	fake.GetLimitsStub = nil
	if fake.getLimitsReturnsOnCall == nil {
		fake.getLimitsReturnsOnCall = make(map[int]struct {
			result1 []byte
			result2 error
		})
	}
	fake.getLimitsReturnsOnCall[i] = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *FakeRepository) SaveAccount(arg1 string, arg2 []byte) error {
	var arg2Copy []byte
	if arg2 != nil {
		arg2Copy = make([]byte, len(arg2))
		copy(arg2Copy, arg2)
	}
	fake.saveAccountMutex.Lock()
	ret, specificReturn := fake.saveAccountReturnsOnCall[len(fake.saveAccountArgsForCall)]
	fake.saveAccountArgsForCall = append(fake.saveAccountArgsForCall, struct {
		arg1 string
		arg2 []byte
	}{arg1, arg2Copy})
	stub := fake.SaveAccountStub
	fakeReturns := fake.saveAccountReturns
	fake.recordInvocation("SaveAccount", []interface{}{arg1, arg2Copy})
	fake.saveAccountMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeRepository) SaveAccountCallCount() int {
	fake.saveAccountMutex.RLock()
	defer fake.saveAccountMutex.RUnlock()
	return len(fake.saveAccountArgsForCall)
}

func (fake *FakeRepository) SaveAccountCalls(stub func(string, []byte) error) {
	fake.saveAccountMutex.Lock()
	defer fake.saveAccountMutex.Unlock()
	fake.SaveAccountStub = stub
}

func (fake *FakeRepository) SaveAccountArgsForCall(i int) (string, []byte) {
	fake.saveAccountMutex.RLock()
	defer fake.saveAccountMutex.RUnlock()
	argsForCall := fake.saveAccountArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeRepository) SaveAccountReturns(result1 error) {
	fake.saveAccountMutex.Lock()
	defer fake.saveAccountMutex.Unlock()
	fake.SaveAccountStub = nil
	fake.saveAccountReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeRepository) SaveAccountReturnsOnCall(i int, result1 error) {
	fake.saveAccountMutex.Lock()
	defer fake.saveAccountMutex.Unlock()
	fake.SaveAccountStub = nil
	if fake.saveAccountReturnsOnCall == nil {
		fake.saveAccountReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.saveAccountReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeRepository) SaveToken(arg1 string, arg2 []byte) error {
	var arg2Copy []byte
	if arg2 != nil {
		arg2Copy = make([]byte, len(arg2))
		copy(arg2Copy, arg2)
	}
	fake.saveTokenMutex.Lock()
	ret, specificReturn := fake.saveTokenReturnsOnCall[len(fake.saveTokenArgsForCall)]
	fake.saveTokenArgsForCall = append(fake.saveTokenArgsForCall, struct {
		arg1 string
		arg2 []byte
	}{arg1, arg2Copy})
	stub := fake.SaveTokenStub
	fakeReturns := fake.saveTokenReturns
	fake.recordInvocation("SaveToken", []interface{}{arg1, arg2Copy})
	fake.saveTokenMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeRepository) SaveTokenCallCount() int {
	fake.saveTokenMutex.RLock()
	defer fake.saveTokenMutex.RUnlock()
	return len(fake.saveTokenArgsForCall)
}

func (fake *FakeRepository) SaveTokenCalls(stub func(string, []byte) error) {
	fake.saveTokenMutex.Lock()
	defer fake.saveTokenMutex.Unlock()
	fake.SaveTokenStub = stub
}

func (fake *FakeRepository) SaveTokenArgsForCall(i int) (string, []byte) {
	fake.saveTokenMutex.RLock()
	defer fake.saveTokenMutex.RUnlock()
	argsForCall := fake.saveTokenArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeRepository) SaveTokenReturns(result1 error) {
	fake.saveTokenMutex.Lock()
	defer fake.saveTokenMutex.Unlock()
	fake.SaveTokenStub = nil
	fake.saveTokenReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeRepository) SaveTokenReturnsOnCall(i int, result1 error) {
	fake.saveTokenMutex.Lock()
	defer fake.saveTokenMutex.Unlock()
	fake.SaveTokenStub = nil
	if fake.saveTokenReturnsOnCall == nil {
		fake.saveTokenReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.saveTokenReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeRepository) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.accountExistMutex.RLock()
	defer fake.accountExistMutex.RUnlock()
	fake.createUpdateLimitsMutex.RLock()
	defer fake.createUpdateLimitsMutex.RUnlock()
	fake.findAccountMutex.RLock()
	defer fake.findAccountMutex.RUnlock()
	fake.findTokenByEmailMutex.RLock()
	defer fake.findTokenByEmailMutex.RUnlock()
	fake.getAccountOauthClientNameMutex.RLock()
	defer fake.getAccountOauthClientNameMutex.RUnlock()
	fake.getAccountsMutex.RLock()
	defer fake.getAccountsMutex.RUnlock()
	fake.getLimitsMutex.RLock()
	defer fake.getLimitsMutex.RUnlock()
	fake.saveAccountMutex.RLock()
	defer fake.saveAccountMutex.RUnlock()
	fake.saveTokenMutex.RLock()
	defer fake.saveTokenMutex.RUnlock()
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

var _ account.Repository = new(FakeRepository)
