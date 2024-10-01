// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package lambda is a generated GoMock package.
package lambda

import (
	context "context"
	reflect "reflect"

	dto "github.com/57blocks/auto-action/server/internal/dto"
	gomock "github.com/golang/mock/gomock"
	websocket "github.com/gorilla/websocket"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// Info mocks base method.
func (m *MockService) Info(c context.Context, r *dto.ReqURILambda) (*dto.RespInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Info", c, r)
	ret0, _ := ret[0].(*dto.RespInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Info indicates an expected call of Info.
func (mr *MockServiceMockRecorder) Info(c, r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Info", reflect.TypeOf((*MockService)(nil).Info), c, r)
}

// Invoke mocks base method.
func (m *MockService) Invoke(c context.Context, r *dto.ReqInvoke) (*dto.RespInvoke, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Invoke", c, r)
	ret0, _ := ret[0].(*dto.RespInvoke)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Invoke indicates an expected call of Invoke.
func (mr *MockServiceMockRecorder) Invoke(c, r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Invoke", reflect.TypeOf((*MockService)(nil).Invoke), c, r)
}

// List mocks base method.
func (m *MockService) List(c context.Context, isFull bool) (interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", c, isFull)
	ret0, _ := ret[0].(interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockServiceMockRecorder) List(c, isFull interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockService)(nil).List), c, isFull)
}

// Logs mocks base method.
func (m *MockService) Logs(c context.Context, r *dto.ReqURILambda, upgrader *websocket.Upgrader) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Logs", c, r, upgrader)
	ret0, _ := ret[0].(error)
	return ret0
}

// Logs indicates an expected call of Logs.
func (mr *MockServiceMockRecorder) Logs(c, r, upgrader interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Logs", reflect.TypeOf((*MockService)(nil).Logs), c, r, upgrader)
}

// Register mocks base method.
func (m *MockService) Register(c context.Context, r *dto.ReqRegister) ([]*dto.RespRegister, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", c, r)
	ret0, _ := ret[0].([]*dto.RespRegister)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Register indicates an expected call of Register.
func (mr *MockServiceMockRecorder) Register(c, r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockService)(nil).Register), c, r)
}

// Remove mocks base method.
func (m *MockService) Remove(c context.Context, r *dto.ReqURILambda) (*dto.RespRemove, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Remove", c, r)
	ret0, _ := ret[0].(*dto.RespRemove)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Remove indicates an expected call of Remove.
func (mr *MockServiceMockRecorder) Remove(c, r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Remove", reflect.TypeOf((*MockService)(nil).Remove), c, r)
}
