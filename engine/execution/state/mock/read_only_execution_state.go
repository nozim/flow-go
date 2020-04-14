// Code generated by mockery v1.0.0. DO NOT EDIT.

package mock

import delta "github.com/dapperlabs/flow-go/engine/execution/state/delta"
import flow "github.com/dapperlabs/flow-go/model/flow"
import messages "github.com/dapperlabs/flow-go/model/messages"
import mock "github.com/stretchr/testify/mock"

// ReadOnlyExecutionState is an autogenerated mock type for the ReadOnlyExecutionState type
type ReadOnlyExecutionState struct {
	mock.Mock
}

// ChunkDataPackByChunkID provides a mock function with given fields: _a0
func (_m *ReadOnlyExecutionState) ChunkDataPackByChunkID(_a0 flow.Identifier) (*flow.ChunkDataPack, error) {
	ret := _m.Called(_a0)

	var r0 *flow.ChunkDataPack
	if rf, ok := ret.Get(0).(func(flow.Identifier) *flow.ChunkDataPack); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.ChunkDataPack)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(flow.Identifier) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetExecutionResultID provides a mock function with given fields: blockID
func (_m *ReadOnlyExecutionState) GetExecutionResultID(blockID flow.Identifier) (flow.Identifier, error) {
	ret := _m.Called(blockID)

	var r0 flow.Identifier
	if rf, ok := ret.Get(0).(func(flow.Identifier) flow.Identifier); ok {
		r0 = rf(blockID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(flow.Identifier)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(flow.Identifier) error); ok {
		r1 = rf(blockID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetHighestExecutedBlockID provides a mock function with given fields:
func (_m *ReadOnlyExecutionState) GetHighestExecutedBlockID() (uint64, flow.Identifier, error) {
	ret := _m.Called()

	var r0 uint64
	if rf, ok := ret.Get(0).(func() uint64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint64)
	}

	var r1 flow.Identifier
	if rf, ok := ret.Get(1).(func() flow.Identifier); ok {
		r1 = rf()
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(flow.Identifier)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func() error); ok {
		r2 = rf()
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetRegisters provides a mock function with given fields: _a0, _a1
func (_m *ReadOnlyExecutionState) GetRegisters(_a0 []byte, _a1 [][]byte) ([][]byte, error) {
	ret := _m.Called(_a0, _a1)

	var r0 [][]byte
	if rf, ok := ret.Get(0).(func([]byte, [][]byte) [][]byte); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([][]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]byte, [][]byte) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRegistersWithProofs provides a mock function with given fields: _a0, _a1
func (_m *ReadOnlyExecutionState) GetRegistersWithProofs(_a0 []byte, _a1 [][]byte) ([][]byte, [][]byte, error) {
	ret := _m.Called(_a0, _a1)

	var r0 [][]byte
	if rf, ok := ret.Get(0).(func([]byte, [][]byte) [][]byte); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([][]byte)
		}
	}

	var r1 [][]byte
	if rf, ok := ret.Get(1).(func([]byte, [][]byte) [][]byte); ok {
		r1 = rf(_a0, _a1)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([][]byte)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func([]byte, [][]byte) error); ok {
		r2 = rf(_a0, _a1)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// NewView provides a mock function with given fields: _a0
func (_m *ReadOnlyExecutionState) NewView(_a0 []byte) *delta.View {
	ret := _m.Called(_a0)

	var r0 *delta.View
	if rf, ok := ret.Get(0).(func([]byte) *delta.View); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*delta.View)
		}
	}

	return r0
}

// RetrieveStateDelta provides a mock function with given fields: blockID
func (_m *ReadOnlyExecutionState) RetrieveStateDelta(blockID flow.Identifier) (*messages.ExecutionStateDelta, error) {
	ret := _m.Called(blockID)

	var r0 *messages.ExecutionStateDelta
	if rf, ok := ret.Get(0).(func(flow.Identifier) *messages.ExecutionStateDelta); ok {
		r0 = rf(blockID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*messages.ExecutionStateDelta)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(flow.Identifier) error); ok {
		r1 = rf(blockID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// StateCommitmentByBlockID provides a mock function with given fields: _a0
func (_m *ReadOnlyExecutionState) StateCommitmentByBlockID(_a0 flow.Identifier) ([]byte, error) {
	ret := _m.Called(_a0)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(flow.Identifier) []byte); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(flow.Identifier) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
