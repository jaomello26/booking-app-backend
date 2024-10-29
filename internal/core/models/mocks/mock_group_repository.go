//go:build test

// Code generated by MockGen. DO NOT EDIT.
// Source: internal/core/models/group.go
//
// Generated by this command:
//
//	mockgen -source=internal/core/models/group.go -destination=internal/core/models/mocks/mock_group_repository.go -package=mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	models "github.com/jaomello26/booking-app-backend/internal/core/models"
	gomock "go.uber.org/mock/gomock"
	gorm "gorm.io/gorm"
)

// MockGroupRepository is a mock of GroupRepository interface.
type MockGroupRepository struct {
	ctrl     *gomock.Controller
	recorder *MockGroupRepositoryMockRecorder
	isgomock struct{}
}

// MockGroupRepositoryMockRecorder is the mock recorder for MockGroupRepository.
type MockGroupRepositoryMockRecorder struct {
	mock *MockGroupRepository
}

// NewMockGroupRepository creates a new mock instance.
func NewMockGroupRepository(ctrl *gomock.Controller) *MockGroupRepository {
	mock := &MockGroupRepository{ctrl: ctrl}
	mock.recorder = &MockGroupRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGroupRepository) EXPECT() *MockGroupRepositoryMockRecorder {
	return m.recorder
}

// CreateOne mocks base method.
func (m *MockGroupRepository) CreateOne(ctx context.Context, group *models.Group) (*models.Group, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateOne", ctx, group)
	ret0, _ := ret[0].(*models.Group)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateOne indicates an expected call of CreateOne.
func (mr *MockGroupRepositoryMockRecorder) CreateOne(ctx, group any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOne", reflect.TypeOf((*MockGroupRepository)(nil).CreateOne), ctx, group)
}

// CreateOneTx mocks base method.
func (m *MockGroupRepository) CreateOneTx(ctx context.Context, tx *gorm.DB, group *models.Group) (*models.Group, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateOneTx", ctx, tx, group)
	ret0, _ := ret[0].(*models.Group)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateOneTx indicates an expected call of CreateOneTx.
func (mr *MockGroupRepositoryMockRecorder) CreateOneTx(ctx, tx, group any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOneTx", reflect.TypeOf((*MockGroupRepository)(nil).CreateOneTx), ctx, tx, group)
}

// DeleteOne mocks base method.
func (m *MockGroupRepository) DeleteOne(ctx context.Context, groupId uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteOne", ctx, groupId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteOne indicates an expected call of DeleteOne.
func (mr *MockGroupRepositoryMockRecorder) DeleteOne(ctx, groupId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteOne", reflect.TypeOf((*MockGroupRepository)(nil).DeleteOne), ctx, groupId)
}

// GetMany mocks base method.
func (m *MockGroupRepository) GetMany(ctx context.Context, userId uint) ([]*models.Group, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMany", ctx, userId)
	ret0, _ := ret[0].([]*models.Group)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMany indicates an expected call of GetMany.
func (mr *MockGroupRepositoryMockRecorder) GetMany(ctx, userId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMany", reflect.TypeOf((*MockGroupRepository)(nil).GetMany), ctx, userId)
}

// GetOne mocks base method.
func (m *MockGroupRepository) GetOne(ctx context.Context, groupId uint) (*models.Group, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOne", ctx, groupId)
	ret0, _ := ret[0].(*models.Group)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOne indicates an expected call of GetOne.
func (mr *MockGroupRepositoryMockRecorder) GetOne(ctx, groupId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOne", reflect.TypeOf((*MockGroupRepository)(nil).GetOne), ctx, groupId)
}

// UpdateOne mocks base method.
func (m *MockGroupRepository) UpdateOne(ctx context.Context, groupId uint, updateData map[string]any) (*models.Group, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateOne", ctx, groupId, updateData)
	ret0, _ := ret[0].(*models.Group)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateOne indicates an expected call of UpdateOne.
func (mr *MockGroupRepositoryMockRecorder) UpdateOne(ctx, groupId, updateData any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateOne", reflect.TypeOf((*MockGroupRepository)(nil).UpdateOne), ctx, groupId, updateData)
}