// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/msaldanha/realChain/consensus (interfaces: ConsensusClient)

// Package tests is a generated GoMock package.
package tests

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	consensus "github.com/msaldanha/realChain/consensus"
	grpc "google.golang.org/grpc"
	reflect "reflect"
)

// MockConsensusClient is a mock of ConsensusClient interface
type MockConsensusClient struct {
	ctrl     *gomock.Controller
	recorder *MockConsensusClientMockRecorder
}

// MockConsensusClientMockRecorder is the mock recorder for MockConsensusClient
type MockConsensusClientMockRecorder struct {
	mock *MockConsensusClient
}

// NewMockConsensusClient creates a new mock instance
func NewMockConsensusClient(ctrl *gomock.Controller) *MockConsensusClient {
	mock := &MockConsensusClient{ctrl: ctrl}
	mock.recorder = &MockConsensusClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockConsensusClient) EXPECT() *MockConsensusClientMockRecorder {
	return m.recorder
}

// Accept mocks base method
func (m *MockConsensusClient) Accept(arg0 context.Context, arg1 *consensus.AcceptRequest, arg2 ...grpc.CallOption) (*consensus.AcceptResult, error) {
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Accept", varargs...)
	ret0, _ := ret[0].(*consensus.AcceptResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Accept indicates an expected call of Accept
func (mr *MockConsensusClientMockRecorder) Accept(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Accept", reflect.TypeOf((*MockConsensusClient)(nil).Accept), varargs...)
}

// Vote mocks base method
func (m *MockConsensusClient) Vote(arg0 context.Context, arg1 *consensus.VoteRequest, arg2 ...grpc.CallOption) (*consensus.VoteResult, error) {
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Vote", varargs...)
	ret0, _ := ret[0].(*consensus.VoteResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Vote indicates an expected call of Vote
func (mr *MockConsensusClientMockRecorder) Vote(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Vote", reflect.TypeOf((*MockConsensusClient)(nil).Vote), varargs...)
}
