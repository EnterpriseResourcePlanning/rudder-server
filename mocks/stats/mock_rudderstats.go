// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/rudderlabs/rudder-server/services/stats (interfaces: RudderStats)

// Package mocks_stats is a generated GoMock package.
package mocks_stats

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockRudderStats is a mock of RudderStats interface
type MockRudderStats struct {
	ctrl     *gomock.Controller
	recorder *MockRudderStatsMockRecorder
}

// MockRudderStatsMockRecorder is the mock recorder for MockRudderStats
type MockRudderStatsMockRecorder struct {
	mock *MockRudderStats
}

// NewMockRudderStats creates a new mock instance
func NewMockRudderStats(ctrl *gomock.Controller) *MockRudderStats {
	mock := &MockRudderStats{ctrl: ctrl}
	mock.recorder = &MockRudderStatsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRudderStats) EXPECT() *MockRudderStatsMockRecorder {
	return m.recorder
}

// Count mocks base method
func (m *MockRudderStats) Count(arg0 int) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Count", arg0)
}

// Count indicates an expected call of Count
func (mr *MockRudderStatsMockRecorder) Count(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockRudderStats)(nil).Count), arg0)
}

// DeferredTimer mocks base method
func (m *MockRudderStats) DeferredTimer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "DeferredTimer")
}

// DeferredTimer indicates an expected call of DeferredTimer
func (mr *MockRudderStatsMockRecorder) DeferredTimer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeferredTimer", reflect.TypeOf((*MockRudderStats)(nil).DeferredTimer))
}

// End mocks base method
func (m *MockRudderStats) End() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "End")
}

// End indicates an expected call of End
func (mr *MockRudderStatsMockRecorder) End() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "End", reflect.TypeOf((*MockRudderStats)(nil).End))
}

// Gauge mocks base method
func (m *MockRudderStats) Gauge(arg0 interface{}) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Gauge", arg0)
}

// Gauge indicates an expected call of Gauge
func (mr *MockRudderStatsMockRecorder) Gauge(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Gauge", reflect.TypeOf((*MockRudderStats)(nil).Gauge), arg0)
}

// Increment mocks base method
func (m *MockRudderStats) Increment() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Increment")
}

// Increment indicates an expected call of Increment
func (mr *MockRudderStatsMockRecorder) Increment() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Increment", reflect.TypeOf((*MockRudderStats)(nil).Increment))
}

// Start mocks base method
func (m *MockRudderStats) Start() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Start")
}

// Start indicates an expected call of Start
func (mr *MockRudderStatsMockRecorder) Start() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockRudderStats)(nil).Start))
}