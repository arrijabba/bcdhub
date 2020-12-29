// Code generated by MockGen. DO NOT EDIT.
// Source: interface.go

// Package mock is a generated GoMock package.
package mock

import (
	models "github.com/baking-bad/bcdhub/internal/models"
	contract "github.com/baking-bad/bcdhub/internal/models/contract"
	gomock "github.com/golang/mock/gomock"
	io "io"
	reflect "reflect"
)

// MockGeneralRepository is a mock of GeneralRepository interface
type MockGeneralRepository struct {
	ctrl     *gomock.Controller
	recorder *MockGeneralRepositoryMockRecorder
}

// MockGeneralRepositoryMockRecorder is the mock recorder for MockGeneralRepository
type MockGeneralRepositoryMockRecorder struct {
	mock *MockGeneralRepository
}

// NewMockGeneralRepository creates a new mock instance
func NewMockGeneralRepository(ctrl *gomock.Controller) *MockGeneralRepository {
	mock := &MockGeneralRepository{ctrl: ctrl}
	mock.recorder = &MockGeneralRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockGeneralRepository) EXPECT() *MockGeneralRepositoryMockRecorder {
	return m.recorder
}

// CreateIndexes mocks base method
func (m *MockGeneralRepository) CreateIndexes() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateIndexes")
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateIndexes indicates an expected call of CreateIndexes
func (mr *MockGeneralRepositoryMockRecorder) CreateIndexes() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateIndexes", reflect.TypeOf((*MockGeneralRepository)(nil).CreateIndexes))
}

// DeleteIndices mocks base method
func (m *MockGeneralRepository) DeleteIndices(indices []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteIndices", indices)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteIndices indicates an expected call of DeleteIndices
func (mr *MockGeneralRepositoryMockRecorder) DeleteIndices(indices interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteIndices", reflect.TypeOf((*MockGeneralRepository)(nil).DeleteIndices), indices)
}

// DeleteByLevelAndNetwork mocks base method
func (m *MockGeneralRepository) DeleteByLevelAndNetwork(arg0 []string, arg1 string, arg2 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteByLevelAndNetwork", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteByLevelAndNetwork indicates an expected call of DeleteByLevelAndNetwork
func (mr *MockGeneralRepositoryMockRecorder) DeleteByLevelAndNetwork(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteByLevelAndNetwork", reflect.TypeOf((*MockGeneralRepository)(nil).DeleteByLevelAndNetwork), arg0, arg1, arg2)
}

// DeleteByContract mocks base method
func (m *MockGeneralRepository) DeleteByContract(indices []string, network, address string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteByContract", indices, network, address)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteByContract indicates an expected call of DeleteByContract
func (mr *MockGeneralRepositoryMockRecorder) DeleteByContract(indices, network, address interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteByContract", reflect.TypeOf((*MockGeneralRepository)(nil).DeleteByContract), indices, network, address)
}

// GetAll mocks base method
func (m *MockGeneralRepository) GetAll(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetAll indicates an expected call of GetAll
func (mr *MockGeneralRepositoryMockRecorder) GetAll(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockGeneralRepository)(nil).GetAll), arg0)
}

// GetByID mocks base method
func (m *MockGeneralRepository) GetByID(arg0 models.Model) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetByID indicates an expected call of GetByID
func (mr *MockGeneralRepositoryMockRecorder) GetByID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockGeneralRepository)(nil).GetByID), arg0)
}

// GetByIDs mocks base method
func (m *MockGeneralRepository) GetByIDs(output interface{}, ids ...string) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{output}
	for _, a := range ids {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetByIDs", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetByIDs indicates an expected call of GetByIDs
func (mr *MockGeneralRepositoryMockRecorder) GetByIDs(output interface{}, ids ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{output}, ids...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByIDs", reflect.TypeOf((*MockGeneralRepository)(nil).GetByIDs), varargs...)
}

// GetByNetwork mocks base method
func (m *MockGeneralRepository) GetByNetwork(arg0 string, arg1 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByNetwork", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetByNetwork indicates an expected call of GetByNetwork
func (mr *MockGeneralRepositoryMockRecorder) GetByNetwork(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByNetwork", reflect.TypeOf((*MockGeneralRepository)(nil).GetByNetwork), arg0, arg1)
}

// GetByNetworkWithSort mocks base method
func (m *MockGeneralRepository) GetByNetworkWithSort(arg0, arg1, arg2 string, arg3 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByNetworkWithSort", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetByNetworkWithSort indicates an expected call of GetByNetworkWithSort
func (mr *MockGeneralRepositoryMockRecorder) GetByNetworkWithSort(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByNetworkWithSort", reflect.TypeOf((*MockGeneralRepository)(nil).GetByNetworkWithSort), arg0, arg1, arg2, arg3)
}

// UpdateDoc mocks base method
func (m *MockGeneralRepository) UpdateDoc(model models.Model) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateDoc", model)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateDoc indicates an expected call of UpdateDoc
func (mr *MockGeneralRepositoryMockRecorder) UpdateDoc(model interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateDoc", reflect.TypeOf((*MockGeneralRepository)(nil).UpdateDoc), model)
}

// UpdateFields mocks base method
func (m *MockGeneralRepository) UpdateFields(arg0, arg1 string, arg2 interface{}, arg3 ...string) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1, arg2}
	for _, a := range arg3 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdateFields", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateFields indicates an expected call of UpdateFields
func (mr *MockGeneralRepositoryMockRecorder) UpdateFields(arg0, arg1, arg2 interface{}, arg3 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1, arg2}, arg3...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateFields", reflect.TypeOf((*MockGeneralRepository)(nil).UpdateFields), varargs...)
}

// GetEvents mocks base method
func (m *MockGeneralRepository) GetEvents(arg0 []models.SubscriptionRequest, arg1, arg2 int64) ([]models.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEvents", arg0, arg1, arg2)
	ret0, _ := ret[0].([]models.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEvents indicates an expected call of GetEvents
func (mr *MockGeneralRepositoryMockRecorder) GetEvents(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEvents", reflect.TypeOf((*MockGeneralRepository)(nil).GetEvents), arg0, arg1, arg2)
}

// SearchByText mocks base method
func (m *MockGeneralRepository) SearchByText(arg0 string, arg1 int64, arg2 []string, arg3 map[string]interface{}, arg4 bool) (models.Result, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchByText", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(models.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchByText indicates an expected call of SearchByText
func (mr *MockGeneralRepositoryMockRecorder) SearchByText(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchByText", reflect.TypeOf((*MockGeneralRepository)(nil).SearchByText), arg0, arg1, arg2, arg3, arg4)
}

// CreateAWSRepository mocks base method
func (m *MockGeneralRepository) CreateAWSRepository(arg0, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAWSRepository", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateAWSRepository indicates an expected call of CreateAWSRepository
func (mr *MockGeneralRepositoryMockRecorder) CreateAWSRepository(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAWSRepository", reflect.TypeOf((*MockGeneralRepository)(nil).CreateAWSRepository), arg0, arg1, arg2)
}

// ListRepositories mocks base method
func (m *MockGeneralRepository) ListRepositories() ([]models.Repository, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListRepositories")
	ret0, _ := ret[0].([]models.Repository)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListRepositories indicates an expected call of ListRepositories
func (mr *MockGeneralRepositoryMockRecorder) ListRepositories() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListRepositories", reflect.TypeOf((*MockGeneralRepository)(nil).ListRepositories))
}

// CreateSnapshots mocks base method
func (m *MockGeneralRepository) CreateSnapshots(arg0, arg1 string, arg2 []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSnapshots", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateSnapshots indicates an expected call of CreateSnapshots
func (mr *MockGeneralRepositoryMockRecorder) CreateSnapshots(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSnapshots", reflect.TypeOf((*MockGeneralRepository)(nil).CreateSnapshots), arg0, arg1, arg2)
}

// RestoreSnapshots mocks base method
func (m *MockGeneralRepository) RestoreSnapshots(arg0, arg1 string, arg2 []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RestoreSnapshots", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// RestoreSnapshots indicates an expected call of RestoreSnapshots
func (mr *MockGeneralRepositoryMockRecorder) RestoreSnapshots(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RestoreSnapshots", reflect.TypeOf((*MockGeneralRepository)(nil).RestoreSnapshots), arg0, arg1, arg2)
}

// ListSnapshots mocks base method
func (m *MockGeneralRepository) ListSnapshots(arg0 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListSnapshots", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListSnapshots indicates an expected call of ListSnapshots
func (mr *MockGeneralRepositoryMockRecorder) ListSnapshots(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListSnapshots", reflect.TypeOf((*MockGeneralRepository)(nil).ListSnapshots), arg0)
}

// SetSnapshotPolicy mocks base method
func (m *MockGeneralRepository) SetSnapshotPolicy(arg0, arg1, arg2, arg3 string, arg4 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetSnapshotPolicy", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetSnapshotPolicy indicates an expected call of SetSnapshotPolicy
func (mr *MockGeneralRepositoryMockRecorder) SetSnapshotPolicy(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetSnapshotPolicy", reflect.TypeOf((*MockGeneralRepository)(nil).SetSnapshotPolicy), arg0, arg1, arg2, arg3, arg4)
}

// GetAllPolicies mocks base method
func (m *MockGeneralRepository) GetAllPolicies() ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllPolicies")
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllPolicies indicates an expected call of GetAllPolicies
func (mr *MockGeneralRepositoryMockRecorder) GetAllPolicies() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllPolicies", reflect.TypeOf((*MockGeneralRepository)(nil).GetAllPolicies))
}

// GetMappings mocks base method
func (m *MockGeneralRepository) GetMappings(arg0 []string) (map[string]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMappings", arg0)
	ret0, _ := ret[0].(map[string]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMappings indicates an expected call of GetMappings
func (mr *MockGeneralRepositoryMockRecorder) GetMappings(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMappings", reflect.TypeOf((*MockGeneralRepository)(nil).GetMappings), arg0)
}

// CreateMapping mocks base method
func (m *MockGeneralRepository) CreateMapping(arg0 string, arg1 io.Reader) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateMapping", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateMapping indicates an expected call of CreateMapping
func (mr *MockGeneralRepositoryMockRecorder) CreateMapping(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateMapping", reflect.TypeOf((*MockGeneralRepository)(nil).CreateMapping), arg0, arg1)
}

// ReloadSecureSettings mocks base method
func (m *MockGeneralRepository) ReloadSecureSettings() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReloadSecureSettings")
	ret0, _ := ret[0].(error)
	return ret0
}

// ReloadSecureSettings indicates an expected call of ReloadSecureSettings
func (mr *MockGeneralRepositoryMockRecorder) ReloadSecureSettings() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReloadSecureSettings", reflect.TypeOf((*MockGeneralRepository)(nil).ReloadSecureSettings))
}

// GetNetworkCountStats mocks base method
func (m *MockGeneralRepository) GetNetworkCountStats(arg0 string) (map[string]int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNetworkCountStats", arg0)
	ret0, _ := ret[0].(map[string]int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNetworkCountStats indicates an expected call of GetNetworkCountStats
func (mr *MockGeneralRepositoryMockRecorder) GetNetworkCountStats(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNetworkCountStats", reflect.TypeOf((*MockGeneralRepository)(nil).GetNetworkCountStats), arg0)
}

// GetDateHistogram mocks base method
func (m *MockGeneralRepository) GetDateHistogram(period string, opts ...models.HistogramOption) ([][]int64, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{period}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetDateHistogram", varargs...)
	ret0, _ := ret[0].([][]int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDateHistogram indicates an expected call of GetDateHistogram
func (mr *MockGeneralRepositoryMockRecorder) GetDateHistogram(period interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{period}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDateHistogram", reflect.TypeOf((*MockGeneralRepository)(nil).GetDateHistogram), varargs...)
}

// GetCallsCountByNetwork mocks base method
func (m *MockGeneralRepository) GetCallsCountByNetwork() (map[string]int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCallsCountByNetwork")
	ret0, _ := ret[0].(map[string]int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCallsCountByNetwork indicates an expected call of GetCallsCountByNetwork
func (mr *MockGeneralRepositoryMockRecorder) GetCallsCountByNetwork() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCallsCountByNetwork", reflect.TypeOf((*MockGeneralRepository)(nil).GetCallsCountByNetwork))
}

// GetContractStatsByNetwork mocks base method
func (m *MockGeneralRepository) GetContractStatsByNetwork() (map[string]models.ContractCountStats, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetContractStatsByNetwork")
	ret0, _ := ret[0].(map[string]models.ContractCountStats)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetContractStatsByNetwork indicates an expected call of GetContractStatsByNetwork
func (mr *MockGeneralRepositoryMockRecorder) GetContractStatsByNetwork() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetContractStatsByNetwork", reflect.TypeOf((*MockGeneralRepository)(nil).GetContractStatsByNetwork))
}

// GetFACountByNetwork mocks base method
func (m *MockGeneralRepository) GetFACountByNetwork() (map[string]int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFACountByNetwork")
	ret0, _ := ret[0].(map[string]int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFACountByNetwork indicates an expected call of GetFACountByNetwork
func (mr *MockGeneralRepositoryMockRecorder) GetFACountByNetwork() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFACountByNetwork", reflect.TypeOf((*MockGeneralRepository)(nil).GetFACountByNetwork))
}

// GetLanguagesForNetwork mocks base method
func (m *MockGeneralRepository) GetLanguagesForNetwork(network string) (map[string]int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLanguagesForNetwork", network)
	ret0, _ := ret[0].(map[string]int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLanguagesForNetwork indicates an expected call of GetLanguagesForNetwork
func (mr *MockGeneralRepositoryMockRecorder) GetLanguagesForNetwork(network interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLanguagesForNetwork", reflect.TypeOf((*MockGeneralRepository)(nil).GetLanguagesForNetwork), network)
}

// IsRecordNotFound mocks base method
func (m *MockGeneralRepository) IsRecordNotFound(err error) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsRecordNotFound", err)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsRecordNotFound indicates an expected call of IsRecordNotFound
func (mr *MockGeneralRepositoryMockRecorder) IsRecordNotFound(err interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsRecordNotFound", reflect.TypeOf((*MockGeneralRepository)(nil).IsRecordNotFound), err)
}

// MockBulkRepository is a mock of BulkRepository interface
type MockBulkRepository struct {
	ctrl     *gomock.Controller
	recorder *MockBulkRepositoryMockRecorder
}

// MockBulkRepositoryMockRecorder is the mock recorder for MockBulkRepository
type MockBulkRepositoryMockRecorder struct {
	mock *MockBulkRepository
}

// NewMockBulkRepository creates a new mock instance
func NewMockBulkRepository(ctrl *gomock.Controller) *MockBulkRepository {
	mock := &MockBulkRepository{ctrl: ctrl}
	mock.recorder = &MockBulkRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockBulkRepository) EXPECT() *MockBulkRepositoryMockRecorder {
	return m.recorder
}

// Insert mocks base method
func (m *MockBulkRepository) Insert(arg0 []models.Model) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Insert indicates an expected call of Insert
func (mr *MockBulkRepositoryMockRecorder) Insert(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockBulkRepository)(nil).Insert), arg0)
}

// Update mocks base method
func (m *MockBulkRepository) Update(arg0 []models.Model) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update
func (mr *MockBulkRepositoryMockRecorder) Update(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockBulkRepository)(nil).Update), arg0)
}

// Delete mocks base method
func (m *MockBulkRepository) Delete(arg0 []models.Model) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockBulkRepositoryMockRecorder) Delete(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockBulkRepository)(nil).Delete), arg0)
}

// RemoveField mocks base method
func (m *MockBulkRepository) RemoveField(arg0 string, arg1 []models.Model) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveField", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveField indicates an expected call of RemoveField
func (mr *MockBulkRepositoryMockRecorder) RemoveField(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveField", reflect.TypeOf((*MockBulkRepository)(nil).RemoveField), arg0, arg1)
}

// UpdateField mocks base method
func (m *MockBulkRepository) UpdateField(where []contract.Contract, fields ...string) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{where}
	for _, a := range fields {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdateField", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateField indicates an expected call of UpdateField
func (mr *MockBulkRepositoryMockRecorder) UpdateField(where interface{}, fields ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{where}, fields...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateField", reflect.TypeOf((*MockBulkRepository)(nil).UpdateField), varargs...)
}