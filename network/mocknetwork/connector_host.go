// Code generated by mockery v2.21.4. DO NOT EDIT.

package mocknetwork

import (
	network "github.com/libp2p/go-libp2p/core/network"
	mock "github.com/stretchr/testify/mock"

	peer "github.com/libp2p/go-libp2p/core/peer"
)

// ConnectorHost is an autogenerated mock type for the ConnectorHost type
type ConnectorHost struct {
	mock.Mock
}

// ClosePeer provides a mock function with given fields: id
func (_m *ConnectorHost) ClosePeer(id peer.ID) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(peer.ID) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Connections provides a mock function with given fields:
func (_m *ConnectorHost) Connections() []network.Conn {
	ret := _m.Called()

	var r0 []network.Conn
	if rf, ok := ret.Get(0).(func() []network.Conn); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]network.Conn)
		}
	}

	return r0
}

// ID provides a mock function with given fields:
func (_m *ConnectorHost) ID() peer.ID {
	ret := _m.Called()

	var r0 peer.ID
	if rf, ok := ret.Get(0).(func() peer.ID); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(peer.ID)
	}

	return r0
}

// IsProtected provides a mock function with given fields: id
func (_m *ConnectorHost) IsProtected(id peer.ID) bool {
	ret := _m.Called(id)

	var r0 bool
	if rf, ok := ret.Get(0).(func(peer.ID) bool); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// PeerInfo provides a mock function with given fields: id
func (_m *ConnectorHost) PeerInfo(id peer.ID) peer.AddrInfo {
	ret := _m.Called(id)

	var r0 peer.AddrInfo
	if rf, ok := ret.Get(0).(func(peer.ID) peer.AddrInfo); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(peer.AddrInfo)
	}

	return r0
}

type mockConstructorTestingTNewConnectorHost interface {
	mock.TestingT
	Cleanup(func())
}

// NewConnectorHost creates a new instance of ConnectorHost. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewConnectorHost(t mockConstructorTestingTNewConnectorHost) *ConnectorHost {
	mock := &ConnectorHost{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
