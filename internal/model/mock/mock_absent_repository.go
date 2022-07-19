// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/Himatro2021/API/internal/model (interfaces: AbsentRepository)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"
	time "time"

	model "github.com/Himatro2021/API/internal/model"
	gomock "github.com/golang/mock/gomock"
)

// MockAbsentRepository is a mock of AbsentRepository interface.
type MockAbsentRepository struct {
	ctrl     *gomock.Controller
	recorder *MockAbsentRepositoryMockRecorder
}

// MockAbsentRepositoryMockRecorder is the mock recorder for MockAbsentRepository.
type MockAbsentRepositoryMockRecorder struct {
	mock *MockAbsentRepository
}

// NewMockAbsentRepository creates a new mock instance.
func NewMockAbsentRepository(ctrl *gomock.Controller) *MockAbsentRepository {
	mock := &MockAbsentRepository{ctrl: ctrl}
	mock.recorder = &MockAbsentRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAbsentRepository) EXPECT() *MockAbsentRepositoryMockRecorder {
	return m.recorder
}

// CreateAbsentForm mocks base method.
func (m *MockAbsentRepository) CreateAbsentForm(arg0 context.Context, arg1 string, arg2, arg3 time.Time, arg4 int64) (*model.AbsentForm, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAbsentForm", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(*model.AbsentForm)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateAbsentForm indicates an expected call of CreateAbsentForm.
func (mr *MockAbsentRepositoryMockRecorder) CreateAbsentForm(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAbsentForm", reflect.TypeOf((*MockAbsentRepository)(nil).CreateAbsentForm), arg0, arg1, arg2, arg3, arg4)
}

// FillAbsentFormByAttendee mocks base method.
func (m *MockAbsentRepository) FillAbsentFormByAttendee(arg0 context.Context, arg1, arg2 int64, arg3 model.Status, arg4 string) (*model.AbsentList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FillAbsentFormByAttendee", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(*model.AbsentList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FillAbsentFormByAttendee indicates an expected call of FillAbsentFormByAttendee.
func (mr *MockAbsentRepositoryMockRecorder) FillAbsentFormByAttendee(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FillAbsentFormByAttendee", reflect.TypeOf((*MockAbsentRepository)(nil).FillAbsentFormByAttendee), arg0, arg1, arg2, arg3, arg4)
}

// GetAbsentFormByID mocks base method.
func (m *MockAbsentRepository) GetAbsentFormByID(arg0 context.Context, arg1 int64) (*model.AbsentForm, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAbsentFormByID", arg0, arg1)
	ret0, _ := ret[0].(*model.AbsentForm)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAbsentFormByID indicates an expected call of GetAbsentFormByID.
func (mr *MockAbsentRepositoryMockRecorder) GetAbsentFormByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAbsentFormByID", reflect.TypeOf((*MockAbsentRepository)(nil).GetAbsentFormByID), arg0, arg1)
}

// GetAbsentListByCreatorID mocks base method.
func (m *MockAbsentRepository) GetAbsentListByCreatorID(arg0 context.Context, arg1, arg2 int64) (*model.AbsentList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAbsentListByCreatorID", arg0, arg1, arg2)
	ret0, _ := ret[0].(*model.AbsentList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAbsentListByCreatorID indicates an expected call of GetAbsentListByCreatorID.
func (mr *MockAbsentRepositoryMockRecorder) GetAbsentListByCreatorID(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAbsentListByCreatorID", reflect.TypeOf((*MockAbsentRepository)(nil).GetAbsentListByCreatorID), arg0, arg1, arg2)
}

// GetAbsentListByID mocks base method.
func (m *MockAbsentRepository) GetAbsentListByID(arg0 context.Context, arg1, arg2 int64) (*model.AbsentList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAbsentListByID", arg0, arg1, arg2)
	ret0, _ := ret[0].(*model.AbsentList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAbsentListByID indicates an expected call of GetAbsentListByID.
func (mr *MockAbsentRepositoryMockRecorder) GetAbsentListByID(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAbsentListByID", reflect.TypeOf((*MockAbsentRepository)(nil).GetAbsentListByID), arg0, arg1, arg2)
}

// GetAllAbsentForm mocks base method.
func (m *MockAbsentRepository) GetAllAbsentForm(arg0 context.Context, arg1, arg2 int) ([]model.AbsentForm, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllAbsentForm", arg0, arg1, arg2)
	ret0, _ := ret[0].([]model.AbsentForm)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllAbsentForm indicates an expected call of GetAllAbsentForm.
func (mr *MockAbsentRepositoryMockRecorder) GetAllAbsentForm(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllAbsentForm", reflect.TypeOf((*MockAbsentRepository)(nil).GetAllAbsentForm), arg0, arg1, arg2)
}

// GetParticipantsByFormID mocks base method.
func (m *MockAbsentRepository) GetParticipantsByFormID(arg0 context.Context, arg1 int64) ([]model.Participant, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetParticipantsByFormID", arg0, arg1)
	ret0, _ := ret[0].([]model.Participant)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetParticipantsByFormID indicates an expected call of GetParticipantsByFormID.
func (mr *MockAbsentRepositoryMockRecorder) GetParticipantsByFormID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetParticipantsByFormID", reflect.TypeOf((*MockAbsentRepository)(nil).GetParticipantsByFormID), arg0, arg1)
}

// UpdateAbsentForm mocks base method.
func (m *MockAbsentRepository) UpdateAbsentForm(arg0 context.Context, arg1 int64, arg2 string, arg3, arg4 time.Time, arg5 int64) (*model.AbsentForm, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateAbsentForm", arg0, arg1, arg2, arg3, arg4, arg5)
	ret0, _ := ret[0].(*model.AbsentForm)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateAbsentForm indicates an expected call of UpdateAbsentForm.
func (mr *MockAbsentRepositoryMockRecorder) UpdateAbsentForm(arg0, arg1, arg2, arg3, arg4, arg5 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAbsentForm", reflect.TypeOf((*MockAbsentRepository)(nil).UpdateAbsentForm), arg0, arg1, arg2, arg3, arg4, arg5)
}

// UpdateAbsentListByAttendee mocks base method.
func (m *MockAbsentRepository) UpdateAbsentListByAttendee(arg0 context.Context, arg1 *model.AbsentList) (*model.AbsentList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateAbsentListByAttendee", arg0, arg1)
	ret0, _ := ret[0].(*model.AbsentList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateAbsentListByAttendee indicates an expected call of UpdateAbsentListByAttendee.
func (mr *MockAbsentRepositoryMockRecorder) UpdateAbsentListByAttendee(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAbsentListByAttendee", reflect.TypeOf((*MockAbsentRepository)(nil).UpdateAbsentListByAttendee), arg0, arg1)
}
