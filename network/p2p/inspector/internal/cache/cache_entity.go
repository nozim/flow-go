package cache

import (
	"time"

	"github.com/onflow/flow-go/model/flow"
)

// RecordEntity is an entity that represents a tracking record that keeps track
// of the amount of cluster prefixed topics received from a peer. This struct
// implements the flow.Entity interface and uses a flow.Identifier created from
// the records peer field for deduplication.
type RecordEntity struct {
	ClusterPrefixTopicsReceivedRecord
	lastUpdated time.Time
}

var _ flow.Entity = (*RecordEntity)(nil)

// NewRecordEntity returns a new *RecordEntity creating the Identifier from the ClusterPrefixTopicsReceivedRecord
// peer field.
func NewRecordEntity(identifier flow.Identifier) RecordEntity {
	return RecordEntity{
		ClusterPrefixTopicsReceivedRecord: NewClusterPrefixTopicsReceivedRecord(identifier),
		lastUpdated:                       time.Now(),
	}
}

// ID returns the origin id of the spam record, which is used as the unique identifier of the entity for maintenance and
// deduplication purposes in the cache.
func (r RecordEntity) ID() flow.Identifier {
	return r.Identifier
}

// Checksum returns the origin id of the spam record, it does not have any purpose in the cache.
// It is implemented to satisfy the flow.Entity interface.
func (r RecordEntity) Checksum() flow.Identifier {
	return r.Identifier
}
