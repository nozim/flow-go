package synchronization

import (
	"context"
	"fmt"
	"io"
	"math"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	mockconsensus "github.com/onflow/flow-go/engine/consensus/mock"
	"github.com/onflow/flow-go/model/chainsync"
	"github.com/onflow/flow-go/model/flow"
	"github.com/onflow/flow-go/model/flow/filter"
	"github.com/onflow/flow-go/model/messages"
	synccore "github.com/onflow/flow-go/module/chainsync"
	"github.com/onflow/flow-go/module/id"
	"github.com/onflow/flow-go/module/irrecoverable"
	"github.com/onflow/flow-go/module/metrics"
	module "github.com/onflow/flow-go/module/mock"
	netint "github.com/onflow/flow-go/network"
	"github.com/onflow/flow-go/network/channels"
	"github.com/onflow/flow-go/network/mocknetwork"
	"github.com/onflow/flow-go/network/p2p/cache"
	protocolint "github.com/onflow/flow-go/state/protocol"
	protocolEvents "github.com/onflow/flow-go/state/protocol/events"
	protocol "github.com/onflow/flow-go/state/protocol/mock"
	storerr "github.com/onflow/flow-go/storage"
	storage "github.com/onflow/flow-go/storage/mock"
	"github.com/onflow/flow-go/utils/rand"
	"github.com/onflow/flow-go/utils/unittest"
)

func TestSyncEngine(t *testing.T) {
	suite.Run(t, new(SyncSuite))
}

type SyncSuite struct {
	suite.Suite
	myID         flow.Identifier
	participants flow.IdentityList
	head         *flow.Header
	heights      map[uint64]*flow.Block
	blockIDs     map[flow.Identifier]*flow.Block
	net          *mocknetwork.Network
	con          *mocknetwork.Conduit
	me           *module.Local
	state        *protocol.State
	snapshot     *protocol.Snapshot
	blocks       *storage.Blocks
	comp         *mockconsensus.Compliance
	core         *module.SyncCore
	e            *Engine
}

func (ss *SyncSuite) SetupTest() {
	// generate own ID
	ss.participants = unittest.IdentityListFixture(3, unittest.WithRole(flow.RoleConsensus))
	keys := unittest.NetworkingKeys(len(ss.participants))

	for i, p := range ss.participants {
		p.NetworkPubKey = keys[i].PublicKey()
	}
	ss.myID = ss.participants[0].NodeID

	// generate a header for the final state
	header := unittest.BlockHeaderFixture()
	ss.head = header

	// create maps to enable block returns
	ss.heights = make(map[uint64]*flow.Block)
	ss.blockIDs = make(map[flow.Identifier]*flow.Block)

	// set up the network module mock
	ss.net = &mocknetwork.Network{}
	ss.net.On("Register", mock.Anything, mock.Anything).Return(
		func(channel channels.Channel, engine netint.MessageProcessor) netint.Conduit {
			return ss.con
		},
		nil,
	)

	// set up the network conduit mock
	ss.con = &mocknetwork.Conduit{}

	// set up the local module mock
	ss.me = &module.Local{}
	ss.me.On("NodeID").Return(
		func() flow.Identifier {
			return ss.myID
		},
	)

	// set up the protocol state mock
	ss.state = &protocol.State{}
	ss.state.On("Final").Return(
		func() protocolint.Snapshot {
			return ss.snapshot
		},
	)
	ss.state.On("AtBlockID", mock.Anything).Return(
		func(blockID flow.Identifier) protocolint.Snapshot {
			if ss.head.ID() == blockID {
				return ss.snapshot
			} else {
				return unittest.StateSnapshotForUnknownBlock()
			}
		},
	).Maybe()

	// set up the snapshot mock
	ss.snapshot = &protocol.Snapshot{}
	ss.snapshot.On("Head").Return(
		func() *flow.Header {
			return ss.head
		},
		nil,
	)
	ss.snapshot.On("Identities", mock.Anything).Return(
		func(selector flow.IdentityFilter) flow.IdentityList {
			return ss.participants.Filter(selector)
		},
		nil,
	)

	// set up blocks storage mock
	ss.blocks = &storage.Blocks{}
	ss.blocks.On("ByHeight", mock.Anything).Return(
		func(height uint64) *flow.Block {
			return ss.heights[height]
		},
		func(height uint64) error {
			_, enabled := ss.heights[height]
			if !enabled {
				return storerr.ErrNotFound
			}
			return nil
		},
	)
	ss.blocks.On("ByID", mock.Anything).Return(
		func(blockID flow.Identifier) *flow.Block {
			return ss.blockIDs[blockID]
		},
		func(blockID flow.Identifier) error {
			_, enabled := ss.blockIDs[blockID]
			if !enabled {
				return storerr.ErrNotFound
			}
			return nil
		},
	)

	// set up compliance engine mock
	ss.comp = mockconsensus.NewCompliance(ss.T())
	ss.comp.On("Process", mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

	// set up sync core
	ss.core = &module.SyncCore{}

	// initialize the engine
	log := zerolog.New(io.Discard)
	metrics := metrics.NewNoopCollector()

	idCache, err := cache.NewProtocolStateIDCache(log, ss.state, protocolEvents.NewDistributor())
	require.NoError(ss.T(), err, "could not create protocol state identity cache")
	e, err := New(log, metrics, ss.net, ss.me, ss.state, ss.blocks, ss.comp, ss.core,
		id.NewIdentityFilterIdentifierProvider(
			filter.And(
				filter.HasRole(flow.RoleConsensus),
				filter.Not(filter.HasNodeID(ss.me.NodeID())),
			),
			idCache,
		),
		NewSpamDetectionConfig())
	require.NoError(ss.T(), err, "should pass engine initialization")

	ss.e = e
}

// TestOnSyncRequest_LowerThanReceiver_WithinTolerance tests that a sync request that's within tolerance of the receiver doesn't trigger
// a response, even if request height is lower than receiver.
func (ss *SyncSuite) TestOnSyncRequest_LowerThanReceiver_WithinTolerance() {
	nonce, err := rand.Uint64()
	require.NoError(ss.T(), err, "should generate nonce")
	// generate origin and request message
	originID := unittest.IdentifierFixture()
	req := &messages.SyncRequest{
		Nonce:  nonce,
		Height: 0,
	}

	// regardless of request height, if within tolerance, we should not respond
	ss.core.On("HandleHeight", ss.head, req.Height)
	ss.core.On("WithinTolerance", ss.head, req.Height).Return(true)
	ss.Assert().NoError(ss.e.requestHandler.onSyncRequest(originID, req))
	ss.con.AssertNotCalled(ss.T(), "Unicast", mock.Anything, mock.Anything)
	ss.core.AssertExpectations(ss.T())
}

// TestOnSyncRequest_HigherThanReceiver_OutsideTolerance tests that a sync request that's higher
// than the receiver's height doesn't trigger a response, even if outside tolerance.
func (ss *SyncSuite) TestOnSyncRequest_HigherThanReceiver_OutsideTolerance() {
	nonce, err := rand.Uint64()
	require.NoError(ss.T(), err, "should generate nonce")
	// generate origin and request message
	originID := unittest.IdentifierFixture()
	req := &messages.SyncRequest{
		Nonce:  nonce,
		Height: 0,
	}

	// if request height is higher than local finalized, we should not respond
	req.Height = ss.head.Height + 1

	ss.core.On("HandleHeight", ss.head, req.Height)
	ss.core.On("WithinTolerance", ss.head, req.Height).Return(false)
	ss.Assert().NoError(ss.e.requestHandler.onSyncRequest(originID, req))
	ss.con.AssertNotCalled(ss.T(), "Unicast", mock.Anything, mock.Anything)
	ss.core.AssertExpectations(ss.T())
}

// TestProcess_SyncRequest_HigherThanReceiver_OutsideTolerance_NoMisbehaviorReport tests that a sync request that's higher
// than the receiver's height doesn't trigger a response, even if outside tolerance and does not generate ALSP
// spamming misbehavior report (simulating the most likely probability).
func (ss *SyncSuite) TestProcess_SyncRequest_HigherThanReceiver_OutsideTolerance_NoMisbehaviorReport() {
	ctx, cancel := irrecoverable.NewMockSignalerContextWithCancel(ss.T(), context.Background())
	ss.e.Start(ctx)
	unittest.AssertClosesBefore(ss.T(), ss.e.Ready(), time.Second)
	defer cancel()

	// generate origin and request message
	originID := unittest.IdentifierFixture()

	nonce, err := rand.Uint64()
	require.NoError(ss.T(), err, "should generate nonce")

	req := &messages.SyncRequest{
		Nonce:  nonce,
		Height: 0,
	}

	// if request height is higher than local finalized, we should not respond
	req.Height = ss.head.Height + 1

	ss.core.On("HandleHeight", ss.head, req.Height).Once()
	ss.core.On("WithinTolerance", ss.head, req.Height).Return(false).Once()

	ss.con.AssertNotCalled(ss.T(), "Unicast", mock.Anything, mock.Anything)

	ss.e.spamDetectionConfig.syncRequestProbability = 0.0 // force not creating misbehavior report

	require.NoError(ss.T(), ss.e.Process(channels.SyncCommittee, originID, req))

	// give at least some time to process items
	time.Sleep(time.Millisecond * 100)

	ss.core.AssertExpectations(ss.T())
	ss.con.AssertExpectations(ss.T())
}

// TestLoad_Process_SyncRequest_HigherThanReceiver_OutsideTolerance_AlwaysReportSpam tests that a sync request that's higher
// than the receiver's height doesn't trigger a response, even if outside tolerance and generates ALSP
// spamming misbehavior report (simulating the unlikely probability).
// This load test ensures that a misbehavior report is generated every time when the probability factor is set to 1.0.
func (ss *SyncSuite) TestLoad_Process_SyncRequest_HigherThanReceiver_OutsideTolerance_AlwaysReportSpam() {
	ctx, cancel := irrecoverable.NewMockSignalerContextWithCancel(ss.T(), context.Background())
	ss.e.Start(ctx)
	unittest.AssertClosesBefore(ss.T(), ss.e.Ready(), time.Second)
	defer cancel()

	load := 1000

	// reset misbehavior report counter for each subtest
	misbehaviorsCounter := 0

	for i := 0; i < load; i++ {
		// generate origin and request message
		originID := unittest.IdentifierFixture()

		nonce, err := rand.Uint64()
		require.NoError(ss.T(), err, "should generate nonce")

		req := &messages.SyncRequest{
			Nonce:  nonce,
			Height: 0,
		}

		// if request height is higher than local finalized, we should not respond
		req.Height = ss.head.Height + 1

		// assert that HandleHeight, WithinTolerance are not called because misbehavior is reported
		// also, check that response is never sent
		ss.core.AssertNotCalled(ss.T(), "HandleHeight")
		ss.core.AssertNotCalled(ss.T(), "WithinTolerance")
		ss.con.AssertNotCalled(ss.T(), "Unicast", mock.Anything, mock.Anything)

		// count misbehavior reports over the course of a load test
		ss.con.On("ReportMisbehavior", mock.Anything).Return(mock.Anything).Run(
			func(args mock.Arguments) {
				misbehaviorsCounter++
			},
		)

		// force creating misbehavior report by setting syncRequestProbability to 1.0 (i.e. report misbehavior 100% of the time)
		ss.e.spamDetectionConfig.syncRequestProbability = 1.0

		require.NoError(ss.T(), ss.e.Process(channels.SyncCommittee, originID, req))
	}

	ss.core.AssertExpectations(ss.T())
	ss.con.AssertExpectations(ss.T())
	assert.Equal(ss.T(), misbehaviorsCounter, load) // should generate misbehavior report every time
}

// TestLoad_Process_SyncRequest_HigherThanReceiver_OutsideTolerance_SometimesReportSpam load tests that a sync request that's higher
// than the receiver's height doesn't trigger a response, even if outside tolerance. It checks that an ALSP
// spam misbehavior report was generated and that the number of misbehavior reports is within a reasonable range.
// This load test ensures that a misbehavior report is generated an appropriate range of times when the probability factor is set to different values.
func (ss *SyncSuite) TestLoad_Process_SyncRequest_HigherThanReceiver_OutsideTolerance_SometimesReportSpam() {
	ctx, cancel := irrecoverable.NewMockSignalerContextWithCancel(ss.T(), context.Background())
	ss.e.Start(ctx)
	unittest.AssertClosesBefore(ss.T(), ss.e.Ready(), time.Second)
	defer cancel()

	load := 1000

	type loadGroup struct {
		syncRequestProbabilityFactor float32
		expectedMisbehaviorsLower    int
		expectedMisbehaviorsUpper    int
	}

	loadGroups := []loadGroup{}

	// expect to never get misbehavior report
	loadGroups = append(loadGroups, loadGroup{0.0, 0, 0})

	// expect to get misbehavior report between 10% of the time
	loadGroups = append(loadGroups, loadGroup{0.1, 75, 140})

	// expect to get misbehavior report between 1% of the time
	loadGroups = append(loadGroups, loadGroup{0.01, 5, 15})

	// expect to get misbehavior report between 0.1% of the time (1 in 1000 requests)
	loadGroups = append(loadGroups, loadGroup{0.001, 0, 7})

	// expect to get misbehavior report between 50% of the time
	loadGroups = append(loadGroups, loadGroup{0.5, 450, 550})

	// expect to get misbehavior report between 90% of the time
	loadGroups = append(loadGroups, loadGroup{0.9, 850, 950})

	// reset misbehavior report counter for each subtest
	misbehaviorsCounter := 0

	for _, loadGroup := range loadGroups {
		ss.T().Run(fmt.Sprintf("load test; pfactor=%f lower=%d upper=%d", loadGroup.syncRequestProbabilityFactor, loadGroup.expectedMisbehaviorsLower, loadGroup.expectedMisbehaviorsUpper), func(t *testing.T) {
			for i := 0; i < load; i++ {
				ss.T().Log("load iteration", i)
				nonce, err := rand.Uint64()
				require.NoError(ss.T(), err, "should generate nonce")

				// generate origin and request message
				originID := unittest.IdentifierFixture()
				req := &messages.SyncRequest{
					Nonce:  nonce,
					Height: 0,
				}

				// if request height is higher than local finalized, we should not respond
				req.Height = ss.head.Height + 1

				ss.core.On("HandleHeight", ss.head, req.Height)
				ss.core.On("WithinTolerance", ss.head, req.Height).Return(false)
				ss.con.AssertNotCalled(ss.T(), "Unicast", mock.Anything, mock.Anything)

				// maybe function calls that might or might not occur over the course of the load test
				ss.core.On("ScanPending", ss.head).Return([]chainsync.Range{}, []chainsync.Batch{}).Maybe()
				ss.con.On("Multicast", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

				// count misbehavior reports over the course of a load test
				ss.con.On("ReportMisbehavior", mock.Anything).Return(mock.Anything).Maybe().Run(
					func(args mock.Arguments) {
						misbehaviorsCounter++
					},
				)
				ss.e.spamDetectionConfig.syncRequestProbability = loadGroup.syncRequestProbabilityFactor
				require.NoError(ss.T(), ss.e.Process(channels.SyncCommittee, originID, req))
			}

			// check function call expectations at the end of the load test; otherwise, load test would take much longer
			ss.core.AssertExpectations(ss.T())
			ss.con.AssertExpectations(ss.T())

			// check that correct range of misbehavior reports were generated (between 1-2 reports per 1000 requests)
			// since we're using a random method to generate misbehavior reports, we can't guarantee the exact number, so we
			// check that it's within a larger range, but that at least 1 misbehavior report was generated

			ss.T().Logf("misbehaviors counter after load test: %d (expected lower bound: %d expected upper bound: %d)", misbehaviorsCounter, loadGroup.expectedMisbehaviorsLower, loadGroup.expectedMisbehaviorsUpper)
			assert.GreaterOrEqual(ss.T(), misbehaviorsCounter, loadGroup.expectedMisbehaviorsLower)
			assert.LessOrEqual(ss.T(), misbehaviorsCounter, loadGroup.expectedMisbehaviorsUpper) // too many reports would indicate a bug

			misbehaviorsCounter = 0 // reset counter for next subtest
		})
	}
}

// TestOnSyncRequest_LowerThanReceiver_OutsideTolerance tests that a sync request that's outside tolerance and
// lower than the receiver's height triggers a response.
func (ss *SyncSuite) TestOnSyncRequest_LowerThanReceiver_OutsideTolerance() {
	nonce, err := rand.Uint64()
	require.NoError(ss.T(), err, "should generate nonce")

	// generate origin and request message
	originID := unittest.IdentifierFixture()
	req := &messages.SyncRequest{
		Nonce:  nonce,
		Height: 0,
	}

	// if the request height is lower than head and outside tolerance, we should expect correct response
	req.Height = ss.head.Height - 1
	ss.core.On("HandleHeight", ss.head, req.Height)
	ss.core.On("WithinTolerance", ss.head, req.Height).Return(false)
	ss.con.On("Unicast", mock.Anything, mock.Anything).Return(nil).Run(
		func(args mock.Arguments) {
			res := args.Get(0).(*messages.SyncResponse)
			assert.Equal(ss.T(), ss.head.Height, res.Height, "response should contain head height")
			assert.Equal(ss.T(), req.Nonce, res.Nonce, "response should contain request nonce")
			recipientID := args.Get(1).(flow.Identifier)
			assert.Equal(ss.T(), originID, recipientID, "should send response to original sender")
		},
	)
	err = ss.e.requestHandler.onSyncRequest(originID, req)
	require.NoError(ss.T(), err, "smaller height sync request should pass")

	ss.core.AssertExpectations(ss.T())
}

func (ss *SyncSuite) TestOnSyncResponse() {
	nonce, err := rand.Uint64()
	require.NoError(ss.T(), err, "should generate nonce")

	height, err := rand.Uint64()
	require.NoError(ss.T(), err, "should generate height")

	// generate origin ID and response message
	originID := unittest.IdentifierFixture()
	res := &messages.SyncResponse{
		Nonce:  nonce,
		Height: height,
	}

	// the height should be handled
	ss.core.On("HandleHeight", ss.head, res.Height)
	ss.e.onSyncResponse(originID, res)
	ss.core.AssertExpectations(ss.T())
}

func (ss *SyncSuite) TestOnRangeRequest() {
	nonce, err := rand.Uint64()
	require.NoError(ss.T(), err, "should generate nonce")

	// generate originID and range request
	originID := unittest.IdentifierFixture()
	req := &messages.RangeRequest{
		Nonce:      nonce,
		FromHeight: 0,
		ToHeight:   0,
	}

	// fill in blocks at heights -1 to -4 from head
	ref := ss.head.Height
	for height := ref; height >= ref-4; height-- {
		block := unittest.BlockFixture()
		block.Header.Height = height
		ss.heights[height] = &block
	}

	// empty range should be a no-op
	ss.T().Run("empty range", func(t *testing.T) {
		req.FromHeight = ref
		req.ToHeight = ref - 1
		err := ss.e.requestHandler.onRangeRequest(originID, req)
		require.NoError(ss.T(), err, "empty range request should pass")
		ss.con.AssertNumberOfCalls(ss.T(), "Unicast", 0)
	})

	// range with only unknown block should be a no-op
	ss.T().Run("range with unknown block", func(t *testing.T) {
		req.FromHeight = ref + 1
		req.ToHeight = ref + 3
		err := ss.e.requestHandler.onRangeRequest(originID, req)
		require.NoError(ss.T(), err, "unknown range request should pass")
		ss.con.AssertNumberOfCalls(ss.T(), "Unicast", 0)
	})

	// a request for same from and to should send single block
	ss.T().Run("from == to", func(t *testing.T) {
		req.FromHeight = ref - 1
		req.ToHeight = ref - 1
		ss.con.On("Unicast", mock.Anything, mock.Anything).Return(nil).Once().Run(
			func(args mock.Arguments) {
				res := args.Get(0).(*messages.BlockResponse)
				expected := ss.heights[ref-1]
				actual := res.Blocks[0].ToInternal()
				assert.Equal(ss.T(), expected, actual, "response should contain right block")
				assert.Equal(ss.T(), req.Nonce, res.Nonce, "response should contain request nonce")
				recipientID := args.Get(1).(flow.Identifier)
				assert.Equal(ss.T(), originID, recipientID, "should send response to original requester")
			},
		)
		err := ss.e.requestHandler.onRangeRequest(originID, req)
		require.NoError(ss.T(), err, "range request with higher to height should pass")
	})

	// a request for a range that we partially have should send partial response
	ss.T().Run("have partial range", func(t *testing.T) {
		req.FromHeight = ref - 2
		req.ToHeight = ref + 2
		ss.con.On("Unicast", mock.Anything, mock.Anything).Return(nil).Once().Run(
			func(args mock.Arguments) {
				res := args.Get(0).(*messages.BlockResponse)
				expected := []*flow.Block{ss.heights[ref-2], ss.heights[ref-1], ss.heights[ref]}
				assert.ElementsMatch(ss.T(), expected, res.BlocksInternal(), "response should contain right blocks")
				assert.Equal(ss.T(), req.Nonce, res.Nonce, "response should contain request nonce")
				recipientID := args.Get(1).(flow.Identifier)
				assert.Equal(ss.T(), originID, recipientID, "should send response to original requester")
			},
		)
		err := ss.e.requestHandler.onRangeRequest(originID, req)
		require.NoError(ss.T(), err, "valid range with missing blocks should fail")
	})

	// a request for a range we entirely have should send all blocks
	ss.T().Run("have entire range", func(t *testing.T) {
		req.FromHeight = ref - 2
		req.ToHeight = ref
		ss.con.On("Unicast", mock.Anything, mock.Anything).Return(nil).Once().Run(
			func(args mock.Arguments) {
				res := args.Get(0).(*messages.BlockResponse)
				expected := []*flow.Block{ss.heights[ref-2], ss.heights[ref-1], ss.heights[ref]}
				assert.ElementsMatch(ss.T(), expected, res.BlocksInternal(), "response should contain right blocks")
				assert.Equal(ss.T(), req.Nonce, res.Nonce, "response should contain request nonce")
				recipientID := args.Get(1).(flow.Identifier)
				assert.Equal(ss.T(), originID, recipientID, "should send response to original requester")
			},
		)
		err := ss.e.requestHandler.onRangeRequest(originID, req)
		require.NoError(ss.T(), err, "valid range request should pass")
	})

	// a request for a range larger than MaxSize should be clamped
	ss.T().Run("oversized range", func(t *testing.T) {
		req.FromHeight = ref - 4
		req.ToHeight = math.MaxUint64
		ss.con.On("Unicast", mock.Anything, mock.Anything).Return(nil).Once().Run(
			func(args mock.Arguments) {
				res := args.Get(0).(*messages.BlockResponse)
				expected := []*flow.Block{ss.heights[ref-4], ss.heights[ref-3], ss.heights[ref-2]}
				assert.ElementsMatch(ss.T(), expected, res.BlocksInternal(), "response should contain right blocks")
				assert.Equal(ss.T(), req.Nonce, res.Nonce, "response should contain request nonce")
				recipientID := args.Get(1).(flow.Identifier)
				assert.Equal(ss.T(), originID, recipientID, "should send response to original requester")
			},
		)

		// Rebuild sync core with a smaller max size
		var err error
		config := synccore.DefaultConfig()
		config.MaxSize = 2
		ss.e.requestHandler.core, err = synccore.New(ss.e.log, config, metrics.NewNoopCollector(), flow.Localnet)
		require.NoError(ss.T(), err)

		err = ss.e.requestHandler.onRangeRequest(originID, req)
		require.NoError(ss.T(), err, "valid range request exceeding max size should still pass")
	})
}

func (ss *SyncSuite) TestOnBatchRequest() {
	nonce, err := rand.Uint64()
	require.NoError(ss.T(), err, "should generate nonce")

	// generate origin ID and batch request
	originID := unittest.IdentifierFixture()
	req := &messages.BatchRequest{
		Nonce:    nonce,
		BlockIDs: nil,
	}

	// an empty request should not lead to response
	ss.T().Run("empty request", func(t *testing.T) {
		req.BlockIDs = []flow.Identifier{}
		err := ss.e.requestHandler.onBatchRequest(originID, req)
		require.NoError(ss.T(), err, "should pass empty request")
		ss.con.AssertNumberOfCalls(ss.T(), "Unicast", 0)
	})

	// a non-empty request for missing block ID should be a no-op
	ss.T().Run("request for missing blocks", func(t *testing.T) {
		req.BlockIDs = unittest.IdentifierListFixture(1)
		err := ss.e.requestHandler.onBatchRequest(originID, req)
		require.NoError(ss.T(), err, "should pass request for missing block")
		ss.con.AssertNumberOfCalls(ss.T(), "Unicast", 0)
	})

	// a non-empty request for existing block IDs should send right response
	ss.T().Run("request for existing blocks", func(t *testing.T) {
		block := unittest.BlockFixture()
		block.Header.Height = ss.head.Height - 1
		req.BlockIDs = []flow.Identifier{block.ID()}
		ss.blockIDs[block.ID()] = &block
		ss.con.On("Unicast", mock.Anything, mock.Anything).Return(nil).Run(
			func(args mock.Arguments) {
				res := args.Get(0).(*messages.BlockResponse)
				assert.Equal(ss.T(), &block, res.Blocks[0].ToInternal(), "response should contain right block")
				assert.Equal(ss.T(), req.Nonce, res.Nonce, "response should contain request nonce")
				recipientID := args.Get(1).(flow.Identifier)
				assert.Equal(ss.T(), originID, recipientID, "response should be send to original requester")
			},
		).Once()
		err := ss.e.requestHandler.onBatchRequest(originID, req)
		require.NoError(ss.T(), err, "should pass request with valid block")
	})

	// a request for too many blocks should be clamped
	ss.T().Run("oversized range", func(t *testing.T) {
		// setup request for 5 blocks. response should contain the first 2 (MaxSize)
		ss.blockIDs = make(map[flow.Identifier]*flow.Block)
		req.BlockIDs = make([]flow.Identifier, 5)
		for i := 0; i < len(req.BlockIDs); i++ {
			b := unittest.BlockFixture()
			b.Header.Height = ss.head.Height - uint64(i)
			req.BlockIDs[i] = b.ID()
			ss.blockIDs[b.ID()] = &b
		}
		ss.con.On("Unicast", mock.Anything, mock.Anything).Return(nil).Run(
			func(args mock.Arguments) {
				res := args.Get(0).(*messages.BlockResponse)
				assert.ElementsMatch(ss.T(), []*flow.Block{ss.blockIDs[req.BlockIDs[0]], ss.blockIDs[req.BlockIDs[1]]}, res.BlocksInternal(), "response should contain right block")
				assert.Equal(ss.T(), req.Nonce, res.Nonce, "response should contain request nonce")
				recipientID := args.Get(1).(flow.Identifier)
				assert.Equal(ss.T(), originID, recipientID, "response should be send to original requester")
			},
		)

		// Rebuild sync core with a smaller max size
		var err error
		config := synccore.DefaultConfig()
		config.MaxSize = 2
		ss.e.requestHandler.core, err = synccore.New(ss.e.log, config, metrics.NewNoopCollector(), flow.Localnet)
		require.NoError(ss.T(), err)

		err = ss.e.requestHandler.onBatchRequest(originID, req)
		require.NoError(ss.T(), err, "valid batch request exceeding max size should still pass")
	})
}

func (ss *SyncSuite) TestOnBlockResponse() {
	nonce, err := rand.Uint64()
	require.NoError(ss.T(), err, "should generate nonce")

	// generate origin and block response
	originID := unittest.IdentifierFixture()
	res := &messages.BlockResponse{
		Nonce:  nonce,
		Blocks: []messages.UntrustedBlock{},
	}

	// add one block that should be processed
	processable := unittest.BlockFixture()
	ss.core.On("HandleBlock", processable.Header).Return(true)
	res.Blocks = append(res.Blocks, messages.UntrustedBlockFromInternal(&processable))

	// add one block that should not be processed
	unprocessable := unittest.BlockFixture()
	ss.core.On("HandleBlock", unprocessable.Header).Return(false)
	res.Blocks = append(res.Blocks, messages.UntrustedBlockFromInternal(&unprocessable))

	ss.comp.On("OnSyncedBlocks", mock.Anything).Run(func(args mock.Arguments) {
		res := args.Get(0).(flow.Slashable[[]*messages.BlockProposal])
		converted := res.Message[0].Block.ToInternal()
		ss.Assert().Equal(processable.Header, converted.Header)
		ss.Assert().Equal(processable.Payload, converted.Payload)
		ss.Assert().Equal(originID, res.OriginID)
	})

	ss.e.onBlockResponse(originID, res)
	ss.core.AssertExpectations(ss.T())
}

func (ss *SyncSuite) TestPollHeight() {

	// check that we send to three nodes from our total list
	others := ss.participants.Filter(filter.HasNodeID(ss.participants[1:].NodeIDs()...))
	ss.con.On("Multicast", mock.Anything, synccore.DefaultPollNodes, others[0].NodeID, others[1].NodeID).Return(nil).Run(
		func(args mock.Arguments) {
			req := args.Get(0).(*messages.SyncRequest)
			require.Equal(ss.T(), ss.head.Height, req.Height, "request should contain finalized height")
		},
	)
	ss.e.pollHeight()
	ss.con.AssertExpectations(ss.T())
}

func (ss *SyncSuite) TestSendRequests() {

	ranges := unittest.RangeListFixture(1)
	batches := unittest.BatchListFixture(1)

	// should submit and mark requested all ranges
	ss.con.On("Multicast", mock.AnythingOfType("*messages.RangeRequest"), synccore.DefaultBlockRequestNodes, mock.Anything, mock.Anything).Return(nil).Run(
		func(args mock.Arguments) {
			req := args.Get(0).(*messages.RangeRequest)
			ss.Assert().Equal(ranges[0].From, req.FromHeight)
			ss.Assert().Equal(ranges[0].To, req.ToHeight)
		},
	)
	ss.core.On("RangeRequested", ranges[0])

	// should submit and mark requested all batches
	ss.con.On("Multicast", mock.AnythingOfType("*messages.BatchRequest"), synccore.DefaultBlockRequestNodes, mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(
		func(args mock.Arguments) {
			req := args.Get(0).(*messages.BatchRequest)
			ss.Assert().Equal(batches[0].BlockIDs, req.BlockIDs)
		},
	)
	ss.core.On("BatchRequested", batches[0])

	// exclude my node ID
	ss.e.sendRequests(ss.participants[1:].NodeIDs(), ranges, batches)
	ss.con.AssertExpectations(ss.T())
}

// test a synchronization engine can be started and stopped
func (ss *SyncSuite) TestStartStop() {
	ctx, cancel := irrecoverable.NewMockSignalerContextWithCancel(ss.T(), context.Background())
	ss.e.Start(ctx)
	unittest.AssertClosesBefore(ss.T(), ss.e.Ready(), time.Second)
	cancel()
	unittest.AssertClosesBefore(ss.T(), ss.e.Done(), time.Second)
}

// TestProcessingMultipleItems tests that items are processed in async way
func (ss *SyncSuite) TestProcessingMultipleItems() {
	ctx, cancel := irrecoverable.NewMockSignalerContextWithCancel(ss.T(), context.Background())
	ss.e.Start(ctx)
	unittest.AssertClosesBefore(ss.T(), ss.e.Ready(), time.Second)
	defer cancel()

	originID := unittest.IdentifierFixture()
	for i := 0; i < 5; i++ {
		msg := &messages.SyncResponse{
			Nonce:  uint64(i),
			Height: uint64(1000 + i),
		}
		ss.core.On("HandleHeight", mock.Anything, msg.Height).Once()
		require.NoError(ss.T(), ss.e.Process(channels.SyncCommittee, originID, msg))
	}

	finalHeight := ss.head.Height
	for i := 0; i < 5; i++ {
		msg := &messages.SyncRequest{
			Nonce:  uint64(i),
			Height: finalHeight - 100,
		}

		originID := unittest.IdentifierFixture()
		ss.core.On("WithinTolerance", mock.Anything, mock.Anything).Return(false)
		ss.core.On("HandleHeight", mock.Anything, msg.Height).Once()
		ss.con.On("Unicast", mock.Anything, mock.Anything).Return(nil)

		require.NoError(ss.T(), ss.e.Process(channels.SyncCommittee, originID, msg))
	}

	// give at least some time to process items
	time.Sleep(time.Millisecond * 100)

	ss.core.AssertExpectations(ss.T())
}

// TestProcessUnsupportedMessageType tests that Process and ProcessLocal correctly handle a case where invalid message type
// was submitted from network layer.
func (ss *SyncSuite) TestProcessUnsupportedMessageType() {
	invalidEvent := uint64(42)
	engines := []netint.MessageProcessor{ss.e, ss.e.requestHandler}
	for _, e := range engines {
		err := e.Process("ch", unittest.IdentifierFixture(), invalidEvent)
		// shouldn't result in error since byzantine inputs are expected
		require.NoError(ss.T(), err)
	}
}
