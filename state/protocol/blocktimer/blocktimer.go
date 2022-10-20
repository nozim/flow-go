package blocktimer

import (
	"fmt"
	"time"

	"github.com/onflow/flow-go/state/protocol"
)

// Is a functor that generates a current timestamp, usually it's just time.Now().
// Used to make testing easier
type timestampGenerator = func() time.Time

// BlockTimestamp is a helper structure that performs building and validation of valid
// timestamp for blocks that are generated by block builder and checked by hotstuff event loop.
// Let τ be the time stamp of the parent block and t be the current clock time of the proposer that is building the child block
// An honest proposer sets the Timestamp of its proposal according to the following rule:
// if t is within the interval [τ + minInterval, τ + maxInterval], then the proposer sets Timestamp := t
// otherwise, the proposer chooses the time stamp from the interval that is closest to its current time t, i.e.
// if t < τ + minInterval, the proposer sets Timestamp := τ + minInterval
// if τ + maxInterval < t, the proposer sets Timestamp := τ + maxInterval
type BlockTimestamp struct {
	minInterval time.Duration
	maxInterval time.Duration
	generator   timestampGenerator
}

var DefaultBlockTimer = NewNoopBlockTimer()

// NewBlockTimer creates new block timer with specific intervals and time.Now as generator
func NewBlockTimer(minInterval, maxInterval time.Duration) (*BlockTimestamp, error) {
	if minInterval >= maxInterval {
		return nil, fmt.Errorf("invariant minInterval < maxInterval is not satisfied, %d >= %d", minInterval, maxInterval)
	}
	if minInterval <= 0 {
		return nil, fmt.Errorf("invariant minInterval > 0 it not satisifed")
	}

	return &BlockTimestamp{
		minInterval: minInterval,
		maxInterval: maxInterval,
		generator:   func() time.Time { return time.Now().UTC() },
	}, nil
}

// Build generates a timestamp based on definition of valid timestamp.
func (b BlockTimestamp) Build(parentTimestamp time.Time) time.Time {
	// calculate the timestamp and cutoffs
	timestamp := b.generator()
	from := parentTimestamp.Add(b.minInterval)
	to := parentTimestamp.Add(b.maxInterval)

	// adjust timestamp if outside of cutoffs
	if timestamp.Before(from) {
		timestamp = from
	}
	if timestamp.After(to) {
		timestamp = to
	}

	return timestamp
}

// Validate accepts parent and current timestamps and checks if current timestamp satisfies
// definition of valid timestamp.
// Timestamp is valid if: Timestamp ∈ [τ + minInterval, τ + maxInterval]
// Returns:
//   - model.ErrInvalidBlockTimestamp - timestamp is invalid
//   - nil - success
func (b BlockTimestamp) Validate(parentTimestamp, currentTimestamp time.Time) error {
	from := parentTimestamp.Add(b.minInterval)
	to := parentTimestamp.Add(b.maxInterval)
	if currentTimestamp.Before(from) || currentTimestamp.After(to) {
		return protocol.NewInvalidBlockTimestamp("timestamp %v is not within interval [%v; %v]", currentTimestamp, from, to)
	}
	return nil
}
