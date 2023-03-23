// Code generated by mockery v2.21.4. DO NOT EDIT.

package mockp2p

import (
	p2p "github.com/onflow/flow-go/network/p2p"
	mock "github.com/stretchr/testify/mock"
)

// PeerScore is an autogenerated mock type for the PeerScore type
type PeerScore struct {
	mock.Mock
}

// PeerScoreExposer provides a mock function with given fields:
func (_m *PeerScore) PeerScoreExposer() (p2p.PeerScoreExposer, bool) {
	ret := _m.Called()

	var r0 p2p.PeerScoreExposer
	var r1 bool
	if rf, ok := ret.Get(0).(func() (p2p.PeerScoreExposer, bool)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() p2p.PeerScoreExposer); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(p2p.PeerScoreExposer)
		}
	}

	if rf, ok := ret.Get(1).(func() bool); ok {
		r1 = rf()
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// SetPeerScoreExposer provides a mock function with given fields: e
func (_m *PeerScore) SetPeerScoreExposer(e p2p.PeerScoreExposer) {
	_m.Called(e)
}

type mockConstructorTestingTNewPeerScore interface {
	mock.TestingT
	Cleanup(func())
}

// NewPeerScore creates a new instance of PeerScore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewPeerScore(t mockConstructorTestingTNewPeerScore) *PeerScore {
	mock := &PeerScore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}