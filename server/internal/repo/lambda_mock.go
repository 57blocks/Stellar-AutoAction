// Code generated by MockGen. DO NOT EDIT.
// Source: lambda.go

// Package repo is a generated GoMock package.
package repo

import (
	context "context"
	sql "database/sql"
	reflect "reflect"

	dto "github.com/57blocks/auto-action/server/internal/dto"
	gomock "github.com/golang/mock/gomock"
	gorm "gorm.io/gorm"
)

// MockLambda is a mock of Lambda interface.
type MockLambda struct {
	ctrl     *gomock.Controller
	recorder *MockLambdaMockRecorder
}

// MockLambdaMockRecorder is the mock recorder for MockLambda.
type MockLambdaMockRecorder struct {
	mock *MockLambda
}

// NewMockLambda creates a new mock instance.
func NewMockLambda(ctrl *gomock.Controller) *MockLambda {
	mock := &MockLambda{ctrl: ctrl}
	mock.recorder = &MockLambdaMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLambda) EXPECT() *MockLambdaMockRecorder {
	return m.recorder
}

// DeleteLambdaTX mocks base method.
func (m *MockLambda) DeleteLambdaTX(c context.Context, f func(*gorm.DB) error, opts ...*sql.TxOptions) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{c, f}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteLambdaTX", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteLambdaTX indicates an expected call of DeleteLambdaTX.
func (mr *MockLambdaMockRecorder) DeleteLambdaTX(c, f interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{c, f}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteLambdaTX", reflect.TypeOf((*MockLambda)(nil).DeleteLambdaTX), varargs...)
}

// FindByAccount mocks base method.
func (m *MockLambda) FindByAccount(c context.Context, accountId uint64) ([]*dto.RespInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByAccount", c, accountId)
	ret0, _ := ret[0].([]*dto.RespInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByAccount indicates an expected call of FindByAccount.
func (mr *MockLambdaMockRecorder) FindByAccount(c, accountId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByAccount", reflect.TypeOf((*MockLambda)(nil).FindByAccount), c, accountId)
}

// FindByNameOrARN mocks base method.
func (m *MockLambda) FindByNameOrARN(c context.Context, input string) (*dto.RespInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByNameOrARN", c, input)
	ret0, _ := ret[0].(*dto.RespInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByNameOrARN indicates an expected call of FindByNameOrARN.
func (mr *MockLambdaMockRecorder) FindByNameOrARN(c, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByNameOrARN", reflect.TypeOf((*MockLambda)(nil).FindByNameOrARN), c, input)
}

// LambdaInfo mocks base method.
func (m *MockLambda) LambdaInfo(c context.Context, req *dto.ReqURILambda) (*dto.RespInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LambdaInfo", c, req)
	ret0, _ := ret[0].(*dto.RespInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LambdaInfo indicates an expected call of LambdaInfo.
func (mr *MockLambdaMockRecorder) LambdaInfo(c, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LambdaInfo", reflect.TypeOf((*MockLambda)(nil).LambdaInfo), c, req)
}

// PersistRegResult mocks base method.
func (m *MockLambda) PersistRegResult(c context.Context, fc func(*gorm.DB) error, opts ...*sql.TxOptions) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{c, fc}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "PersistRegResult", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// PersistRegResult indicates an expected call of PersistRegResult.
func (mr *MockLambdaMockRecorder) PersistRegResult(c, fc interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{c, fc}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PersistRegResult", reflect.TypeOf((*MockLambda)(nil).PersistRegResult), varargs...)
}
