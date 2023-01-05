// Code generated by MockGen. DO NOT EDIT.
// Source: .\pkg\facade\go.mongodb.org\mongo-driver\mongo\database.go

// Package mock_mongo is a generated GoMock package.
package mock_mongo

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	mongo "go.mongodb.org/mongo-driver/mongo"
	options "go.mongodb.org/mongo-driver/mongo/options"
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

// Collection mocks base method.
func (m *MockDatabase) Collection(db *mongo.Database, name string, opts ...*options.CollectionOptions) *mongo.Collection {
	m.ctrl.T.Helper()
	varargs := []interface{}{db, name}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Collection", varargs...)
	ret0, _ := ret[0].(*mongo.Collection)
	return ret0
}

// Collection indicates an expected call of Collection.
func (mr *MockDatabaseMockRecorder) Collection(db, name interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{db, name}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Collection", reflect.TypeOf((*MockDatabase)(nil).Collection), varargs...)
}
