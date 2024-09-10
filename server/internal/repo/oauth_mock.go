// Code generated by MockGen. DO NOT EDIT.
// Source: oauth.go

// Package repo is a generated GoMock package.
package repo

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

// FindOrg mocks base method.
func (m *MockOAuth) FindOrg(c context.Context, id uint64) (*dto.RespOrg, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindOrg", c, id)
	ret0, _ := ret[0].(*dto.RespOrg)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindOrg indicates an expected call of FindOrg.
func (mr *MockOAuthMockRecorder) FindOrg(c, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOrg", reflect.TypeOf((*MockOAuth)(nil).FindOrg), c, id)
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

// FindTokenByRefresh mocks base method.
func (m *MockOAuth) FindTokenByRefresh(c context.Context, refresh string) (*model.Token, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindTokenByRefresh", c, refresh)
	ret0, _ := ret[0].(*model.Token)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindTokenByRefresh indicates an expected call of FindTokenByRefresh.
func (mr *MockOAuthMockRecorder) FindTokenByRefresh(c, refresh interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindTokenByRefresh", reflect.TypeOf((*MockOAuth)(nil).FindTokenByRefresh), c, refresh)
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
