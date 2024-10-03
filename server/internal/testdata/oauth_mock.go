// Code generated by MockGen. DO NOT EDIT.
// Source: oauth.go

// Package testdata is a generated GoMock package.
package testdata

import (
	context "context"
	reflect "reflect"

	dto "github.com/57blocks/auto-action/server/internal/dto"
	model "github.com/57blocks/auto-action/server/internal/model"
	gomock "github.com/golang/mock/gomock"
)

// MockOAuth is a mock of OAuth interface.
type MockOAuth struct {
	ctrl     *gomock.Controller
	recorder *MockOAuthMockRecorder
}

// MockOAuthMockRecorder is the mock recorder for MockOAuth.
type MockOAuthMockRecorder struct {
	mock *MockOAuth
}

// NewMockOAuth creates a new mock instance.
func NewMockOAuth(ctrl *gomock.Controller) *MockOAuth {
	mock := &MockOAuth{ctrl: ctrl}
	mock.recorder = &MockOAuthMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOAuth) EXPECT() *MockOAuthMockRecorder {
	return m.recorder
}

// CreateUser mocks base method.
func (m *MockOAuth) CreateUser(c context.Context, user *model.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", c, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockOAuthMockRecorder) CreateUser(c, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockOAuth)(nil).CreateUser), c, user)
}

// DeleteTokenByAccess mocks base method.
func (m *MockOAuth) DeleteTokenByAccess(c context.Context, access string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTokenByAccess", c, access)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteTokenByAccess indicates an expected call of DeleteTokenByAccess.
func (mr *MockOAuthMockRecorder) DeleteTokenByAccess(c, access interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTokenByAccess", reflect.TypeOf((*MockOAuth)(nil).DeleteTokenByAccess), c, access)
}

// FindOrgByName mocks base method.
func (m *MockOAuth) FindOrgByName(c context.Context, name string) (*dto.RespOrg, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindOrgByName", c, name)
	ret0, _ := ret[0].(*dto.RespOrg)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindOrgByName indicates an expected call of FindOrgByName.
func (mr *MockOAuthMockRecorder) FindOrgByName(c, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOrgByName", reflect.TypeOf((*MockOAuth)(nil).FindOrgByName), c, name)
}

// FindTokenByRefreshID mocks base method.
func (m *MockOAuth) FindTokenByRefreshID(c context.Context, refresh string) (*model.Token, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindTokenByRefreshID", c, refresh)
	ret0, _ := ret[0].(*model.Token)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindTokenByRefreshID indicates an expected call of FindTokenByRefreshID.
func (mr *MockOAuthMockRecorder) FindTokenByRefreshID(c, refresh interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindTokenByRefreshID", reflect.TypeOf((*MockOAuth)(nil).FindTokenByRefreshID), c, refresh)
}

// FindUserByAcn mocks base method.
func (m *MockOAuth) FindUserByAcn(c context.Context, acn string) (*dto.RespUser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindUserByAcn", c, acn)
	ret0, _ := ret[0].(*dto.RespUser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindUserByAcn indicates an expected call of FindUserByAcn.
func (mr *MockOAuthMockRecorder) FindUserByAcn(c, acn interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindUserByAcn", reflect.TypeOf((*MockOAuth)(nil).FindUserByAcn), c, acn)
}

// FindUserByOrgAcn mocks base method.
func (m *MockOAuth) FindUserByOrgAcn(c context.Context, req *dto.ReqOrgAcn) (*dto.RespUser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindUserByOrgAcn", c, req)
	ret0, _ := ret[0].(*dto.RespUser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindUserByOrgAcn indicates an expected call of FindUserByOrgAcn.
func (mr *MockOAuthMockRecorder) FindUserByOrgAcn(c, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindUserByOrgAcn", reflect.TypeOf((*MockOAuth)(nil).FindUserByOrgAcn), c, req)
}

// SyncToken mocks base method.
func (m *MockOAuth) SyncToken(c context.Context, token *model.Token) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SyncToken", c, token)
	ret0, _ := ret[0].(error)
	return ret0
}

// SyncToken indicates an expected call of SyncToken.
func (mr *MockOAuthMockRecorder) SyncToken(c, token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SyncToken", reflect.TypeOf((*MockOAuth)(nil).SyncToken), c, token)
}
