// Code generated by MockGen. DO NOT EDIT.
// Source: resty.go

// Package restyx is a generated GoMock package.
package restyx

import (
	reflect "reflect"

	dto "github.com/57blocks/auto-action/server/internal/dto"
	gomock "github.com/golang/mock/gomock"
)

// MockResty is a mock of Resty interface.
type MockResty struct {
	ctrl     *gomock.Controller
	recorder *MockRestyMockRecorder
}

// MockRestyMockRecorder is the mock recorder for MockResty.
type MockRestyMockRecorder struct {
	mock *MockResty
}

// NewMockResty creates a new mock instance.
func NewMockResty(ctrl *gomock.Controller) *MockResty {
	mock := &MockResty{ctrl: ctrl}
	mock.recorder = &MockRestyMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockResty) EXPECT() *MockRestyMockRecorder {
	return m.recorder
}

// AddCSKey mocks base method.
func (m *MockResty) AddCSKey(csToken string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddCSKey", csToken)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddCSKey indicates an expected call of AddCSKey.
func (mr *MockRestyMockRecorder) AddCSKey(csToken interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCSKey", reflect.TypeOf((*MockResty)(nil).AddCSKey), csToken)
}

// AddCSKeyToRole mocks base method.
func (m *MockResty) AddCSKeyToRole(csToken, keyId, role string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddCSKeyToRole", csToken, keyId, role)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddCSKeyToRole indicates an expected call of AddCSKeyToRole.
func (mr *MockRestyMockRecorder) AddCSKeyToRole(csToken, keyId, role interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCSKeyToRole", reflect.TypeOf((*MockResty)(nil).AddCSKeyToRole), csToken, keyId, role)
}

// AddCSRole mocks base method.
func (m *MockResty) AddCSRole(csToken, orgName, account string) (*dto.RespAddCsRole, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddCSRole", csToken, orgName, account)
	ret0, _ := ret[0].(*dto.RespAddCsRole)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddCSRole indicates an expected call of AddCSRole.
func (mr *MockRestyMockRecorder) AddCSRole(csToken, orgName, account interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCSRole", reflect.TypeOf((*MockResty)(nil).AddCSRole), csToken, orgName, account)
}

// DeleteCSKey mocks base method.
func (m *MockResty) DeleteCSKey(csToken, keyId string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCSKey", csToken, keyId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCSKey indicates an expected call of DeleteCSKey.
func (mr *MockRestyMockRecorder) DeleteCSKey(csToken, keyId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCSKey", reflect.TypeOf((*MockResty)(nil).DeleteCSKey), csToken, keyId)
}

// DeleteCSKeyFromRole mocks base method.
func (m *MockResty) DeleteCSKeyFromRole(csToken, keyId, role string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCSKeyFromRole", csToken, keyId, role)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCSKeyFromRole indicates an expected call of DeleteCSKeyFromRole.
func (mr *MockRestyMockRecorder) DeleteCSKeyFromRole(csToken, keyId, role interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCSKeyFromRole", reflect.TypeOf((*MockResty)(nil).DeleteCSKeyFromRole), csToken, keyId, role)
}
