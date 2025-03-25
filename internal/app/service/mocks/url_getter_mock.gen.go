// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/Te8va/shortURL/internal/app/service (interfaces: URLGetter)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockURLGetter is a mock of URLGetter interface.
type MockURLGetter struct {
	ctrl     *gomock.Controller
	recorder *MockURLGetterMockRecorder
}

// MockURLGetterMockRecorder is the mock recorder for MockURLGetter.
type MockURLGetterMockRecorder struct {
	mock *MockURLGetter
}

// NewMockURLGetter creates a new mock instance.
func NewMockURLGetter(ctrl *gomock.Controller) *MockURLGetter {
	mock := &MockURLGetter{ctrl: ctrl}
	mock.recorder = &MockURLGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockURLGetter) EXPECT() *MockURLGetterMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockURLGetter) Get(arg0 context.Context, arg1 string, arg2 chan error) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1, arg2)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockURLGetterMockRecorder) Get(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockURLGetter)(nil).Get), arg0, arg1, arg2)
}

// GetUserURLs mocks base method.
func (m *MockURLGetter) GetUserURLs(arg0 context.Context, arg1 int) ([]map[string]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserURLs", arg0, arg1)
	ret0, _ := ret[0].([]map[string]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserURLs indicates an expected call of GetUserURLs.
func (mr *MockURLGetterMockRecorder) GetUserURLs(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserURLs", reflect.TypeOf((*MockURLGetter)(nil).GetUserURLs), arg0, arg1)
}
