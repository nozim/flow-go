// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocknetwork

import (
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	irrecoverable "github.com/onflow/flow-go/module/irrecoverable"

	mock "github.com/stretchr/testify/mock"

	network "github.com/onflow/flow-go/network"
)

// Network is an autogenerated mock type for the Network type
type Network struct {
	mock.Mock
}

// Done provides a mock function with given fields:
func (_m *Network) Done() <-chan struct{} {
	ret := _m.Called()

	var r0 <-chan struct{}
	if rf, ok := ret.Get(0).(func() <-chan struct{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan struct{})
		}
	}

	return r0
}

// Ready provides a mock function with given fields:
func (_m *Network) Ready() <-chan struct{} {
	ret := _m.Called()

	var r0 <-chan struct{}
	if rf, ok := ret.Get(0).(func() <-chan struct{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan struct{})
		}
	}

	return r0
}

// Register provides a mock function with given fields: channel, engine
func (_m *Network) Register(channel network.Channel, engine network.Engine) (network.Conduit, error) {
	ret := _m.Called(channel, engine)

	var r0 network.Conduit
	if rf, ok := ret.Get(0).(func(network.Channel, network.Engine) network.Conduit); ok {
		r0 = rf(channel, engine)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(network.Conduit)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(network.Channel, network.Engine) error); ok {
		r1 = rf(channel, engine)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RegisterBlockExchange provides a mock function with given fields: channel, store
func (_m *Network) RegisterBlockExchange(channel network.Channel, store blockstore.Blockstore) (network.BlockExchange, error) {
	ret := _m.Called(channel, store)

	var r0 network.BlockExchange
	if rf, ok := ret.Get(0).(func(network.Channel, blockstore.Blockstore) network.BlockExchange); ok {
		r0 = rf(channel, store)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(network.BlockExchange)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(network.Channel, blockstore.Blockstore) error); ok {
		r1 = rf(channel, store)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Start provides a mock function with given fields: _a0
func (_m *Network) Start(_a0 irrecoverable.SignalerContext) {
	_m.Called(_a0)
}