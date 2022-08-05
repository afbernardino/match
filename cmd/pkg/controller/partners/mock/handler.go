// Code generated by MockGen. DO NOT EDIT.
// Source: ../handler.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	models "match/cmd/pkg/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockDatabase is a mock of Database interface.
type MockDatabase struct {
	ctrl     *gomock.Controller
	recorder *MockDatabaseMockRecorder
}

// MockDatabaseMockRecorder is the mock recorder for MockDatabase.
type MockDatabaseMockRecorder struct {
	mock *MockDatabase
}

// NewMockDatabase creates a new mock instance.
func NewMockDatabase(ctrl *gomock.Controller) *MockDatabase {
	mock := &MockDatabase{ctrl: ctrl}
	mock.recorder = &MockDatabaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDatabase) EXPECT() *MockDatabaseMockRecorder {
	return m.recorder
}

// GetMatches mocks base method.
func (m *MockDatabase) GetMatches(ctx context.Context, materials []uint, lat, long float32) ([]models.Partner, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMatches", ctx, materials, lat, long)
	ret0, _ := ret[0].([]models.Partner)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMatches indicates an expected call of GetMatches.
func (mr *MockDatabaseMockRecorder) GetMatches(ctx, materials, lat, long interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMatches", reflect.TypeOf((*MockDatabase)(nil).GetMatches), ctx, materials, lat, long)
}

// GetPartnerById mocks base method.
func (m *MockDatabase) GetPartnerById(ctx context.Context, id uint) (models.Partner, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPartnerById", ctx, id)
	ret0, _ := ret[0].(models.Partner)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPartnerById indicates an expected call of GetPartnerById.
func (mr *MockDatabaseMockRecorder) GetPartnerById(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPartnerById", reflect.TypeOf((*MockDatabase)(nil).GetPartnerById), ctx, id)
}
