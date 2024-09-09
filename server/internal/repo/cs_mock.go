// Code generated by MockGen. DO NOT EDIT.
// Source: cs.go

// Package repo is a generated GoMock package.
package repo

import (
	context "context"
	reflect "reflect"

	dto "github.com/57blocks/auto-action/server/internal/dto"
	model "github.com/57blocks/auto-action/server/internal/model"
	gomock "github.com/golang/mock/gomock"
)

// MockCubeSigner is a mock of CubeSigner interface.
type MockCubeSigner struct {
	ctrl     *gomock.Controller
	recorder *MockCubeSignerMockRecorder
}

// MockCubeSignerMockRecorder is the mock recorder for MockCubeSigner.
type MockCubeSignerMockRecorder struct {
	mock *MockCubeSigner
}

// NewMockCubeSigner creates a new mock instance.
func NewMockCubeSigner(ctrl *gomock.Controller) *MockCubeSigner {
	mock := &MockCubeSigner{ctrl: ctrl}
	mock.recorder = &MockCubeSignerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCubeSigner) EXPECT() *MockCubeSignerMockRecorder {
	return m.recorder
}

// FindCSByOrgAcn mocks base method.
func (m *MockCubeSigner) FindCSByOrgAcn(c context.Context, req *dto.ReqCSRole) (*dto.RespCSRole, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindCSByOrgAcn", c, req)
	ret0, _ := ret[0].(*dto.RespCSRole)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindCSByOrgAcn indicates an expected call of FindCSByOrgAcn.
func (mr *MockCubeSignerMockRecorder) FindCSByOrgAcn(c, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindCSByOrgAcn", reflect.TypeOf((*MockCubeSigner)(nil).FindCSByOrgAcn), c, req)
}

// SyncCSKey mocks base method.
func (m *MockCubeSigner) SyncCSKey(c context.Context, key *model.CubeSignerKey) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SyncCSKey", c, key)
	ret0, _ := ret[0].(error)
	return ret0
}

// SyncCSKey indicates an expected call of SyncCSKey.
func (mr *MockCubeSignerMockRecorder) SyncCSKey(c, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SyncCSKey", reflect.TypeOf((*MockCubeSigner)(nil).SyncCSKey), c, key)
}

// ToSign mocks base method.
func (m *MockCubeSigner) ToSign(c context.Context, req *dto.ReqToSign) ([]*dto.RespToSign, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToSign", c, req)
	ret0, _ := ret[0].([]*dto.RespToSign)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ToSign indicates an expected call of ToSign.
func (mr *MockCubeSignerMockRecorder) ToSign(c, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToSign", reflect.TypeOf((*MockCubeSigner)(nil).ToSign), c, req)
}
