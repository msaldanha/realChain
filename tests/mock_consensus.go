// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/msaldanha/realChain/consensus (interfaces: Consensus)

// Package tests is a generated GoMock package.
package tests

import (
	gomock "github.com/golang/mock/gomock"
	consensus "github.com/msaldanha/realChain/consensus"
	reflect "reflect"
)

// MockConsensus is a mock of Consensus interface
type MockConsensus struct {
	ctrl     *gomock.Controller
	recorder *MockConsensusMockRecorder
}

// MockConsensusMockRecorder is the mock recorder for MockConsensus
type MockConsensusMockRecorder struct {
	mock *MockConsensus
}

// NewMockConsensus creates a new mock instance
func NewMockConsensus(ctrl *gomock.Controller) *MockConsensus {
	mock := &MockConsensus{ctrl: ctrl}
	mock.recorder = &MockConsensusMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockConsensus) EXPECT() *MockConsensusMockRecorder {
	return m.recorder
}

// Accept mocks base method
func (m *MockConsensus) Accept(arg0 *consensus.AcceptRequest) (*consensus.AcceptResult, error) {
	ret := m.ctrl.Call(m, "Accept", arg0)
	ret0, _ := ret[0].(*consensus.AcceptResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Accept indicates an expected call of Accept
func (mr *MockConsensusMockRecorder) Accept(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Accept", reflect.TypeOf((*MockConsensus)(nil).Accept), arg0)
}

// Vote mocks base method
func (m *MockConsensus) Vote(arg0 *consensus.VoteRequest) (*consensus.VoteResult, error) {
	ret := m.ctrl.Call(m, "Vote", arg0)
	ret0, _ := ret[0].(*consensus.VoteResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Vote indicates an expected call of Vote
func (mr *MockConsensusMockRecorder) Vote(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Vote", reflect.TypeOf((*MockConsensus)(nil).Vote), arg0)
}
