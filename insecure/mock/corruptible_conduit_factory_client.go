// Code generated by mockery v1.0.0. DO NOT EDIT.

package mockinsecure

import (
	context "context"

	grpc "google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"

	insecure "github.com/onflow/flow-go/insecure"

	mock "github.com/stretchr/testify/mock"
)

// CorruptibleConduitFactoryClient is an autogenerated mock type for the CorruptibleConduitFactoryClient type
type CorruptibleConduitFactoryClient struct {
	mock.Mock
}

// ProcessAttackerMessage provides a mock function with given fields: ctx, opts
func (_m *CorruptibleConduitFactoryClient) ProcessAttackerMessage(ctx context.Context, opts ...grpc.CallOption) (insecure.CorruptibleConduitFactory_ProcessAttackerMessageClient, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 insecure.CorruptibleConduitFactory_ProcessAttackerMessageClient
	if rf, ok := ret.Get(0).(func(context.Context, ...grpc.CallOption) insecure.CorruptibleConduitFactory_ProcessAttackerMessageClient); ok {
		r0 = rf(ctx, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(insecure.CorruptibleConduitFactory_ProcessAttackerMessageClient)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RegisterAttacker provides a mock function with given fields: ctx, in, opts
func (_m *CorruptibleConduitFactoryClient) RegisterAttacker(ctx context.Context, in *insecure.AttackerRegisterMessage, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *insecure.AttackerRegisterMessage, ...grpc.CallOption) *emptypb.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *insecure.AttackerRegisterMessage, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}