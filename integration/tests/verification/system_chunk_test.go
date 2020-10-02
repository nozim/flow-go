package verification

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestSystemChunk(t *testing.T) {
	suite.Run(t, new(SystemChunkTestSuite))
}

type SystemChunkTestSuite struct {
	Suite
}

// TestSystemChunkIDsShouldBeDifferent evaluates that system chunk of consecutive blocks that
// do not cause state change have different chunk Ids.
func (st *SystemChunkTestSuite) TestSystemChunkIDsShouldBeDifferent() {
	// wait for first finalized block, called blockA.
	blockA := st.BlockState.WaitForFirstFinalized(st.T())
	st.T().Logf("blockA generated, height: %v ID: %v", blockA.Header.Height, blockA.Header.ID())

	// waits for the next finalized block after blockA, called blockB.
	blockB := st.BlockState.WaitForFinalizedChild(st.T(), blockA)

	// wait for execution receipt for blockA from execution node, called receiptA.
	receiptA := st.ReceiptState.WaitForReceiptFrom(st.T(), blockA.Header.ID(), st.exeID)
	resultAId := receiptA.ExecutionResult.ID()
	st.T().Logf("receipt for blockA generated: result ID: %x", resultAId)

	// wait for execution receipt for blockB from execution node, called receiptA.
	receiptB := st.ReceiptState.WaitForReceiptFrom(st.T(), blockB.Header.ID(), st.exeID)
	resultBId := receiptB.ExecutionResult.ID()
	st.T().Logf("receipt for blockA generated: result ID: %x", resultBId)

	// Todo: drop this part once system chunk changes the state
	// requires that execution state is not changed between block A and B
	stateA := receiptA.ExecutionResult.FinalStateCommit
	stateB := receiptB.ExecutionResult.FinalStateCommit
	require.Equal(st.T(), stateA, stateB)

	// computes id of system chunk in block A and B
	systemChunkAId := receiptA.ExecutionResult.Chunks[0].ID()
	systemChunkBId := receiptB.ExecutionResult.Chunks[0].ID()

	// requires that system chunk Id of execution results be different
	require.NotEqual(st.T(), systemChunkAId, systemChunkBId)
}
