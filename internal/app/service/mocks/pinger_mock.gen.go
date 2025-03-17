// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/Te8va/shortURL/internal/app/service (interfaces: Pinger)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockPinger is a mock of Pinger interface.
type MockPinger struct {
	ctrl     *gomock.Controller
	recorder *MockPingerMockRecorder
}

// MockPingerMockRecorder is the mock recorder for MockPinger.
type MockPingerMockRecorder struct {
	mock *MockPinger
}

// NewMockPinger creates a new mock instance.
func NewMockPinger(ctrl *gomock.Controller) *MockPinger {
	mock := &MockPinger{ctrl: ctrl}
	mock.recorder = &MockPingerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPinger) EXPECT() *MockPingerMockRecorder {
	return m.recorder
}

// PingPg mocks base method.
func (m *MockPinger) PingPg(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PingPg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// PingPg indicates an expected call of PingPg.
func (mr *MockPingerMockRecorder) PingPg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PingPg", reflect.TypeOf((*MockPinger)(nil).PingPg), arg0)
}
