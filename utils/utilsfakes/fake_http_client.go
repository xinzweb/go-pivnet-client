// Code generated by counterfeiter. DO NOT EDIT.
package utilsfakes

import (
	"io"
	"net/http"
	"sync"

	"github.com/baotingfang/go-pivnet-client/utils"
)

type FakeHttpClient struct {
	DeleteStub        func(string) (*http.Response, error)
	deleteMutex       sync.RWMutex
	deleteArgsForCall []struct {
		arg1 string
	}
	deleteReturns struct {
		result1 *http.Response
		result2 error
	}
	deleteReturnsOnCall map[int]struct {
		result1 *http.Response
		result2 error
	}
	DoStub        func(*http.Request) (*http.Response, error)
	doMutex       sync.RWMutex
	doArgsForCall []struct {
		arg1 *http.Request
	}
	doReturns struct {
		result1 *http.Response
		result2 error
	}
	doReturnsOnCall map[int]struct {
		result1 *http.Response
		result2 error
	}
	GetStub        func(string) (*http.Response, error)
	getMutex       sync.RWMutex
	getArgsForCall []struct {
		arg1 string
	}
	getReturns struct {
		result1 *http.Response
		result2 error
	}
	getReturnsOnCall map[int]struct {
		result1 *http.Response
		result2 error
	}
	PatchStub        func(string, io.Reader) (*http.Response, error)
	patchMutex       sync.RWMutex
	patchArgsForCall []struct {
		arg1 string
		arg2 io.Reader
	}
	patchReturns struct {
		result1 *http.Response
		result2 error
	}
	patchReturnsOnCall map[int]struct {
		result1 *http.Response
		result2 error
	}
	PostStub        func(string, io.Reader) (*http.Response, error)
	postMutex       sync.RWMutex
	postArgsForCall []struct {
		arg1 string
		arg2 io.Reader
	}
	postReturns struct {
		result1 *http.Response
		result2 error
	}
	postReturnsOnCall map[int]struct {
		result1 *http.Response
		result2 error
	}
	RefreshAccessTokenStub        func(bool)
	refreshAccessTokenMutex       sync.RWMutex
	refreshAccessTokenArgsForCall []struct {
		arg1 bool
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeHttpClient) Delete(arg1 string) (*http.Response, error) {
	fake.deleteMutex.Lock()
	ret, specificReturn := fake.deleteReturnsOnCall[len(fake.deleteArgsForCall)]
	fake.deleteArgsForCall = append(fake.deleteArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("Delete", []interface{}{arg1})
	fake.deleteMutex.Unlock()
	if fake.DeleteStub != nil {
		return fake.DeleteStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.deleteReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeHttpClient) DeleteCallCount() int {
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	return len(fake.deleteArgsForCall)
}

func (fake *FakeHttpClient) DeleteCalls(stub func(string) (*http.Response, error)) {
	fake.deleteMutex.Lock()
	defer fake.deleteMutex.Unlock()
	fake.DeleteStub = stub
}

func (fake *FakeHttpClient) DeleteArgsForCall(i int) string {
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	argsForCall := fake.deleteArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeHttpClient) DeleteReturns(result1 *http.Response, result2 error) {
	fake.deleteMutex.Lock()
	defer fake.deleteMutex.Unlock()
	fake.DeleteStub = nil
	fake.deleteReturns = struct {
		result1 *http.Response
		result2 error
	}{result1, result2}
}

func (fake *FakeHttpClient) DeleteReturnsOnCall(i int, result1 *http.Response, result2 error) {
	fake.deleteMutex.Lock()
	defer fake.deleteMutex.Unlock()
	fake.DeleteStub = nil
	if fake.deleteReturnsOnCall == nil {
		fake.deleteReturnsOnCall = make(map[int]struct {
			result1 *http.Response
			result2 error
		})
	}
	fake.deleteReturnsOnCall[i] = struct {
		result1 *http.Response
		result2 error
	}{result1, result2}
}

func (fake *FakeHttpClient) Do(arg1 *http.Request) (*http.Response, error) {
	fake.doMutex.Lock()
	ret, specificReturn := fake.doReturnsOnCall[len(fake.doArgsForCall)]
	fake.doArgsForCall = append(fake.doArgsForCall, struct {
		arg1 *http.Request
	}{arg1})
	fake.recordInvocation("Do", []interface{}{arg1})
	fake.doMutex.Unlock()
	if fake.DoStub != nil {
		return fake.DoStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.doReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeHttpClient) DoCallCount() int {
	fake.doMutex.RLock()
	defer fake.doMutex.RUnlock()
	return len(fake.doArgsForCall)
}

func (fake *FakeHttpClient) DoCalls(stub func(*http.Request) (*http.Response, error)) {
	fake.doMutex.Lock()
	defer fake.doMutex.Unlock()
	fake.DoStub = stub
}

func (fake *FakeHttpClient) DoArgsForCall(i int) *http.Request {
	fake.doMutex.RLock()
	defer fake.doMutex.RUnlock()
	argsForCall := fake.doArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeHttpClient) DoReturns(result1 *http.Response, result2 error) {
	fake.doMutex.Lock()
	defer fake.doMutex.Unlock()
	fake.DoStub = nil
	fake.doReturns = struct {
		result1 *http.Response
		result2 error
	}{result1, result2}
}

func (fake *FakeHttpClient) DoReturnsOnCall(i int, result1 *http.Response, result2 error) {
	fake.doMutex.Lock()
	defer fake.doMutex.Unlock()
	fake.DoStub = nil
	if fake.doReturnsOnCall == nil {
		fake.doReturnsOnCall = make(map[int]struct {
			result1 *http.Response
			result2 error
		})
	}
	fake.doReturnsOnCall[i] = struct {
		result1 *http.Response
		result2 error
	}{result1, result2}
}

func (fake *FakeHttpClient) Get(arg1 string) (*http.Response, error) {
	fake.getMutex.Lock()
	ret, specificReturn := fake.getReturnsOnCall[len(fake.getArgsForCall)]
	fake.getArgsForCall = append(fake.getArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("Get", []interface{}{arg1})
	fake.getMutex.Unlock()
	if fake.GetStub != nil {
		return fake.GetStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeHttpClient) GetCallCount() int {
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	return len(fake.getArgsForCall)
}

func (fake *FakeHttpClient) GetCalls(stub func(string) (*http.Response, error)) {
	fake.getMutex.Lock()
	defer fake.getMutex.Unlock()
	fake.GetStub = stub
}

func (fake *FakeHttpClient) GetArgsForCall(i int) string {
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	argsForCall := fake.getArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeHttpClient) GetReturns(result1 *http.Response, result2 error) {
	fake.getMutex.Lock()
	defer fake.getMutex.Unlock()
	fake.GetStub = nil
	fake.getReturns = struct {
		result1 *http.Response
		result2 error
	}{result1, result2}
}

func (fake *FakeHttpClient) GetReturnsOnCall(i int, result1 *http.Response, result2 error) {
	fake.getMutex.Lock()
	defer fake.getMutex.Unlock()
	fake.GetStub = nil
	if fake.getReturnsOnCall == nil {
		fake.getReturnsOnCall = make(map[int]struct {
			result1 *http.Response
			result2 error
		})
	}
	fake.getReturnsOnCall[i] = struct {
		result1 *http.Response
		result2 error
	}{result1, result2}
}

func (fake *FakeHttpClient) Patch(arg1 string, arg2 io.Reader) (*http.Response, error) {
	fake.patchMutex.Lock()
	ret, specificReturn := fake.patchReturnsOnCall[len(fake.patchArgsForCall)]
	fake.patchArgsForCall = append(fake.patchArgsForCall, struct {
		arg1 string
		arg2 io.Reader
	}{arg1, arg2})
	fake.recordInvocation("Patch", []interface{}{arg1, arg2})
	fake.patchMutex.Unlock()
	if fake.PatchStub != nil {
		return fake.PatchStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.patchReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeHttpClient) PatchCallCount() int {
	fake.patchMutex.RLock()
	defer fake.patchMutex.RUnlock()
	return len(fake.patchArgsForCall)
}

func (fake *FakeHttpClient) PatchCalls(stub func(string, io.Reader) (*http.Response, error)) {
	fake.patchMutex.Lock()
	defer fake.patchMutex.Unlock()
	fake.PatchStub = stub
}

func (fake *FakeHttpClient) PatchArgsForCall(i int) (string, io.Reader) {
	fake.patchMutex.RLock()
	defer fake.patchMutex.RUnlock()
	argsForCall := fake.patchArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeHttpClient) PatchReturns(result1 *http.Response, result2 error) {
	fake.patchMutex.Lock()
	defer fake.patchMutex.Unlock()
	fake.PatchStub = nil
	fake.patchReturns = struct {
		result1 *http.Response
		result2 error
	}{result1, result2}
}

func (fake *FakeHttpClient) PatchReturnsOnCall(i int, result1 *http.Response, result2 error) {
	fake.patchMutex.Lock()
	defer fake.patchMutex.Unlock()
	fake.PatchStub = nil
	if fake.patchReturnsOnCall == nil {
		fake.patchReturnsOnCall = make(map[int]struct {
			result1 *http.Response
			result2 error
		})
	}
	fake.patchReturnsOnCall[i] = struct {
		result1 *http.Response
		result2 error
	}{result1, result2}
}

func (fake *FakeHttpClient) Post(arg1 string, arg2 io.Reader) (*http.Response, error) {
	fake.postMutex.Lock()
	ret, specificReturn := fake.postReturnsOnCall[len(fake.postArgsForCall)]
	fake.postArgsForCall = append(fake.postArgsForCall, struct {
		arg1 string
		arg2 io.Reader
	}{arg1, arg2})
	fake.recordInvocation("Post", []interface{}{arg1, arg2})
	fake.postMutex.Unlock()
	if fake.PostStub != nil {
		return fake.PostStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.postReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeHttpClient) PostCallCount() int {
	fake.postMutex.RLock()
	defer fake.postMutex.RUnlock()
	return len(fake.postArgsForCall)
}

func (fake *FakeHttpClient) PostCalls(stub func(string, io.Reader) (*http.Response, error)) {
	fake.postMutex.Lock()
	defer fake.postMutex.Unlock()
	fake.PostStub = stub
}

func (fake *FakeHttpClient) PostArgsForCall(i int) (string, io.Reader) {
	fake.postMutex.RLock()
	defer fake.postMutex.RUnlock()
	argsForCall := fake.postArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeHttpClient) PostReturns(result1 *http.Response, result2 error) {
	fake.postMutex.Lock()
	defer fake.postMutex.Unlock()
	fake.PostStub = nil
	fake.postReturns = struct {
		result1 *http.Response
		result2 error
	}{result1, result2}
}

func (fake *FakeHttpClient) PostReturnsOnCall(i int, result1 *http.Response, result2 error) {
	fake.postMutex.Lock()
	defer fake.postMutex.Unlock()
	fake.PostStub = nil
	if fake.postReturnsOnCall == nil {
		fake.postReturnsOnCall = make(map[int]struct {
			result1 *http.Response
			result2 error
		})
	}
	fake.postReturnsOnCall[i] = struct {
		result1 *http.Response
		result2 error
	}{result1, result2}
}

func (fake *FakeHttpClient) RefreshAccessToken(arg1 bool) {
	fake.refreshAccessTokenMutex.Lock()
	fake.refreshAccessTokenArgsForCall = append(fake.refreshAccessTokenArgsForCall, struct {
		arg1 bool
	}{arg1})
	fake.recordInvocation("RefreshAccessToken", []interface{}{arg1})
	fake.refreshAccessTokenMutex.Unlock()
	if fake.RefreshAccessTokenStub != nil {
		fake.RefreshAccessTokenStub(arg1)
	}
}

func (fake *FakeHttpClient) RefreshAccessTokenCallCount() int {
	fake.refreshAccessTokenMutex.RLock()
	defer fake.refreshAccessTokenMutex.RUnlock()
	return len(fake.refreshAccessTokenArgsForCall)
}

func (fake *FakeHttpClient) RefreshAccessTokenCalls(stub func(bool)) {
	fake.refreshAccessTokenMutex.Lock()
	defer fake.refreshAccessTokenMutex.Unlock()
	fake.RefreshAccessTokenStub = stub
}

func (fake *FakeHttpClient) RefreshAccessTokenArgsForCall(i int) bool {
	fake.refreshAccessTokenMutex.RLock()
	defer fake.refreshAccessTokenMutex.RUnlock()
	argsForCall := fake.refreshAccessTokenArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeHttpClient) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	fake.doMutex.RLock()
	defer fake.doMutex.RUnlock()
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	fake.patchMutex.RLock()
	defer fake.patchMutex.RUnlock()
	fake.postMutex.RLock()
	defer fake.postMutex.RUnlock()
	fake.refreshAccessTokenMutex.RLock()
	defer fake.refreshAccessTokenMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeHttpClient) recordInvocation(key string, args []interface{}) {
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

var _ utils.HttpClient = new(FakeHttpClient)
