// Code generated by mockery v2.21.4. DO NOT EDIT.

package backend

import (
	flow "github.com/onflow/flow-go/model/flow"
	mock "github.com/stretchr/testify/mock"
)

// NodeCommunicatorMock is an autogenerated mock type for the Communicator type
type NodeCommunicatorMock struct {
	mock.Mock
}

// CallAvailableNode provides a mock function with given fields: nodes, call, shouldTerminateOnError
func (_m *NodeCommunicatorMock) CallAvailableNode(nodes flow.IdentityList, call NodeAction, shouldTerminateOnError ErrorTerminator) error {
	ret := _m.Called(nodes, call, shouldTerminateOnError)

	var r0 error
	if rf, ok := ret.Get(0).(func(flow.IdentityList, NodeAction, ErrorTerminator) error); ok {
		r0 = rf(nodes, call, shouldTerminateOnError)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewNodeCommunicatorMock interface {
	mock.TestingT
	Cleanup(func())
}

// NewNodeCommunicatorMock creates a new instance of NodeCommunicatorMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewNodeCommunicatorMock(t mockConstructorTestingTNewNodeCommunicatorMock) *NodeCommunicatorMock {
	mock := &NodeCommunicatorMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
