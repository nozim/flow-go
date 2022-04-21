// Code generated by mockery v1.0.0. DO NOT EDIT.

package mock

import (
	context "context"

	execution "github.com/onflow/flow/protobuf/go/flow/execution"
	mock "github.com/stretchr/testify/mock"
)

// ExecutionAPIServer is an autogenerated mock type for the ExecutionAPIServer type
type ExecutionAPIServer struct {
	mock.Mock
}

// ExecuteScriptAtBlockID provides a mock function with given fields: _a0, _a1
func (_m *ExecutionAPIServer) ExecuteScriptAtBlockID(_a0 context.Context, _a1 *execution.ExecuteScriptAtBlockIDRequest) (*execution.ExecuteScriptAtBlockIDResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *execution.ExecuteScriptAtBlockIDResponse
	if rf, ok := ret.Get(0).(func(context.Context, *execution.ExecuteScriptAtBlockIDRequest) *execution.ExecuteScriptAtBlockIDResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*execution.ExecuteScriptAtBlockIDResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *execution.ExecuteScriptAtBlockIDRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAccountAtBlockID provides a mock function with given fields: _a0, _a1
func (_m *ExecutionAPIServer) GetAccountAtBlockID(_a0 context.Context, _a1 *execution.GetAccountAtBlockIDRequest) (*execution.GetAccountAtBlockIDResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *execution.GetAccountAtBlockIDResponse
	if rf, ok := ret.Get(0).(func(context.Context, *execution.GetAccountAtBlockIDRequest) *execution.GetAccountAtBlockIDResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*execution.GetAccountAtBlockIDResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *execution.GetAccountAtBlockIDRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetBlockHeaderByID provides a mock function with given fields: _a0, _a1
func (_m *ExecutionAPIServer) GetBlockHeaderByID(_a0 context.Context, _a1 *execution.GetBlockHeaderByIDRequest) (*execution.BlockHeaderResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *execution.BlockHeaderResponse
	if rf, ok := ret.Get(0).(func(context.Context, *execution.GetBlockHeaderByIDRequest) *execution.BlockHeaderResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*execution.BlockHeaderResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *execution.GetBlockHeaderByIDRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetEventsForBlockIDs provides a mock function with given fields: _a0, _a1
func (_m *ExecutionAPIServer) GetEventsForBlockIDs(_a0 context.Context, _a1 *execution.GetEventsForBlockIDsRequest) (*execution.GetEventsForBlockIDsResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *execution.GetEventsForBlockIDsResponse
	if rf, ok := ret.Get(0).(func(context.Context, *execution.GetEventsForBlockIDsRequest) *execution.GetEventsForBlockIDsResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*execution.GetEventsForBlockIDsResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *execution.GetEventsForBlockIDsRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLatestBlockHeader provides a mock function with given fields: _a0, _a1
func (_m *ExecutionAPIServer) GetLatestBlockHeader(_a0 context.Context, _a1 *execution.GetLatestBlockHeaderRequest) (*execution.BlockHeaderResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *execution.BlockHeaderResponse
	if rf, ok := ret.Get(0).(func(context.Context, *execution.GetLatestBlockHeaderRequest) *execution.BlockHeaderResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*execution.BlockHeaderResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *execution.GetLatestBlockHeaderRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRegisterAtBlockID provides a mock function with given fields: _a0, _a1
func (_m *ExecutionAPIServer) GetRegisterAtBlockID(_a0 context.Context, _a1 *execution.GetRegisterAtBlockIDRequest) (*execution.GetRegisterAtBlockIDResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *execution.GetRegisterAtBlockIDResponse
	if rf, ok := ret.Get(0).(func(context.Context, *execution.GetRegisterAtBlockIDRequest) *execution.GetRegisterAtBlockIDResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*execution.GetRegisterAtBlockIDResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *execution.GetRegisterAtBlockIDRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTransactionResult provides a mock function with given fields: _a0, _a1
func (_m *ExecutionAPIServer) GetTransactionResult(_a0 context.Context, _a1 *execution.GetTransactionResultRequest) (*execution.GetTransactionResultResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *execution.GetTransactionResultResponse
	if rf, ok := ret.Get(0).(func(context.Context, *execution.GetTransactionResultRequest) *execution.GetTransactionResultResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*execution.GetTransactionResultResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *execution.GetTransactionResultRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTransactionResultByIndex provides a mock function with given fields: _a0, _a1
func (_m *ExecutionAPIServer) GetTransactionResultByIndex(_a0 context.Context, _a1 *execution.GetTransactionByIndexRequest) (*execution.GetTransactionResultResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *execution.GetTransactionResultResponse
	if rf, ok := ret.Get(0).(func(context.Context, *execution.GetTransactionByIndexRequest) *execution.GetTransactionResultResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*execution.GetTransactionResultResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *execution.GetTransactionByIndexRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTransactionResultsByBlockID provides a mock function with given fields: _a0, _a1
func (_m *ExecutionAPIServer) GetTransactionResultsByBlockID(_a0 context.Context, _a1 *execution.GetTransactionsByBlockIDRequest) (*execution.GetTransactionResultsResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *execution.GetTransactionResultsResponse
	if rf, ok := ret.Get(0).(func(context.Context, *execution.GetTransactionsByBlockIDRequest) *execution.GetTransactionResultsResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*execution.GetTransactionResultsResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *execution.GetTransactionsByBlockIDRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Ping provides a mock function with given fields: _a0, _a1
func (_m *ExecutionAPIServer) Ping(_a0 context.Context, _a1 *execution.PingRequest) (*execution.PingResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *execution.PingResponse
	if rf, ok := ret.Get(0).(func(context.Context, *execution.PingRequest) *execution.PingResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*execution.PingResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *execution.PingRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
