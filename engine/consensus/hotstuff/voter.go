package hotstuff

import (
	"github.com/dapperlabs/flow-go/model/hotstuff"
)

// Voter produces votes for the given block
type Voter struct {
	signer        Signer
	viewState     *ViewState
	forks         Forks
	lastVotedView uint64 // need to keep track of the last view we voted for so we don't double vote accidentally
}

// NewVoter creates a new Voter instance
func (v *Voter) NewVoter(signer Signer, viewState *ViewState, forks Forks) *Voter {
	return &Voter{
		signer:        signer,
		viewState:     viewState,
		forks:         forks,
		lastVotedView: 0,
	}
}

// ProduceVoteIfVotable will make a decision on whether it will vote for the given proposal, the returned
// boolean indicates whether to vote or not.
// In order to ensure that only a safe node will be voted, Voter will ask Forks whether a vote is a safe node or not.
// The curView is taken as input to ensure Voter will only vote for proposals at current view and prevent double voting.
// This method will only ever _once_ return a `non-nil, true` vote: the very first time it encounters a safe block of the
//  current view to vote for. Subsequently, voter does _not_ vote for any other block with the same (or lower) view.
// (including repeated calls with the initial block we voted for also return `nil, false`).
func (v *Voter) ProduceVoteIfVotable(block *hotstuff.Block, curView uint64) (*hotstuff.Vote, bool) {
	if v.forks.IsSafeBlock(block) {
		return nil, false
	}

	if curView != block.View {
		return nil, false
	}

	if curView <= v.lastVotedView {
		return nil, false
	}

	vote, err := v.signer.VoteFor(block)
	if err != nil {
		return nil, false
	}

	return vote, true
}
