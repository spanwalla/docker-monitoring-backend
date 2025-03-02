// Code generated by MockGen. DO NOT EDIT.
// Source: internal/repository/repository.go

// Package repository_mocks is a generated GoMock package.
package repository_mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	entity "github.com/spanwalla/docker-monitoring-backend/internal/entity"
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

// CreatePinger mocks base method.
func (m *MockPinger) CreatePinger(ctx context.Context, pinger entity.Pinger) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePinger", ctx, pinger)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreatePinger indicates an expected call of CreatePinger.
func (mr *MockPingerMockRecorder) CreatePinger(ctx, pinger interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePinger", reflect.TypeOf((*MockPinger)(nil).CreatePinger), ctx, pinger)
}

// GetPingerById mocks base method.
func (m *MockPinger) GetPingerById(ctx context.Context, id int) (entity.Pinger, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPingerById", ctx, id)
	ret0, _ := ret[0].(entity.Pinger)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPingerById indicates an expected call of GetPingerById.
func (mr *MockPingerMockRecorder) GetPingerById(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPingerById", reflect.TypeOf((*MockPinger)(nil).GetPingerById), ctx, id)
}

// GetPingerByName mocks base method.
func (m *MockPinger) GetPingerByName(ctx context.Context, name string) (entity.Pinger, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPingerByName", ctx, name)
	ret0, _ := ret[0].(entity.Pinger)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPingerByName indicates an expected call of GetPingerByName.
func (mr *MockPingerMockRecorder) GetPingerByName(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPingerByName", reflect.TypeOf((*MockPinger)(nil).GetPingerByName), ctx, name)
}

// GetPingerByNameAndPassword mocks base method.
func (m *MockPinger) GetPingerByNameAndPassword(ctx context.Context, name, password string) (entity.Pinger, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPingerByNameAndPassword", ctx, name, password)
	ret0, _ := ret[0].(entity.Pinger)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPingerByNameAndPassword indicates an expected call of GetPingerByNameAndPassword.
func (mr *MockPingerMockRecorder) GetPingerByNameAndPassword(ctx, name, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPingerByNameAndPassword", reflect.TypeOf((*MockPinger)(nil).GetPingerByNameAndPassword), ctx, name, password)
}

// MockReport is a mock of Report interface.
type MockReport struct {
	ctrl     *gomock.Controller
	recorder *MockReportMockRecorder
}

// MockReportMockRecorder is the mock recorder for MockReport.
type MockReportMockRecorder struct {
	mock *MockReport
}

// NewMockReport creates a new mock instance.
func NewMockReport(ctrl *gomock.Controller) *MockReport {
	mock := &MockReport{ctrl: ctrl}
	mock.recorder = &MockReportMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockReport) EXPECT() *MockReportMockRecorder {
	return m.recorder
}

// CreateReport mocks base method.
func (m *MockReport) CreateReport(ctx context.Context, report entity.Report) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateReport", ctx, report)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateReport indicates an expected call of CreateReport.
func (mr *MockReportMockRecorder) CreateReport(ctx, report interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateReport", reflect.TypeOf((*MockReport)(nil).CreateReport), ctx, report)
}

// GetLatestReportByEveryPinger mocks base method.
func (m *MockReport) GetLatestReportByEveryPinger(ctx context.Context) ([]entity.Report, []string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLatestReportByEveryPinger", ctx)
	ret0, _ := ret[0].([]entity.Report)
	ret1, _ := ret[1].([]string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetLatestReportByEveryPinger indicates an expected call of GetLatestReportByEveryPinger.
func (mr *MockReportMockRecorder) GetLatestReportByEveryPinger(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLatestReportByEveryPinger", reflect.TypeOf((*MockReport)(nil).GetLatestReportByEveryPinger), ctx)
}
