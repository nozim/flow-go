package validator

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/rs/zerolog"

	"github.com/onflow/flow-go/model/flow"
	cborcodec "github.com/onflow/flow-go/network/codec/cbor"
	"github.com/onflow/flow-go/network/message"
)

func init() {
	// initialize the authorized roles map the first time this package is imported.
	initializeAuthorizedRolesMap()
}

// authorizedRolesMap is a mapping of message type to a list of roles authorized to send them.
var authorizedRolesMap map[uint8]flow.RoleList

// initializeAuthorizedRolesMap initializes authorizedRolesMap.
func initializeAuthorizedRolesMap() {
	authorizedRolesMap = make(map[uint8]flow.RoleList)

	// consensus
	authorizedRolesMap[cborcodec.CodeBlockProposal] = flow.RoleList{flow.RoleConsensus}
	authorizedRolesMap[cborcodec.CodeBlockVote] = flow.RoleList{flow.RoleConsensus}

	// protocol state sync
	authorizedRolesMap[cborcodec.CodeSyncRequest] = flow.Roles()
	authorizedRolesMap[cborcodec.CodeSyncResponse] = flow.Roles()
	authorizedRolesMap[cborcodec.CodeRangeRequest] = flow.Roles()
	authorizedRolesMap[cborcodec.CodeBatchRequest] = flow.Roles()
	authorizedRolesMap[cborcodec.CodeBlockResponse] = flow.Roles()

	// cluster consensus
	authorizedRolesMap[cborcodec.CodeClusterBlockProposal] = flow.RoleList{flow.RoleCollection}
	authorizedRolesMap[cborcodec.CodeClusterBlockVote] = flow.RoleList{flow.RoleCollection}
	authorizedRolesMap[cborcodec.CodeClusterBlockResponse] = flow.RoleList{flow.RoleCollection}

	// collections, guarantees & transactions
	authorizedRolesMap[cborcodec.CodeCollectionGuarantee] = flow.RoleList{flow.RoleCollection}
	authorizedRolesMap[cborcodec.CodeTransactionBody] = flow.RoleList{flow.RoleCollection}
	authorizedRolesMap[cborcodec.CodeTransaction] = flow.RoleList{flow.RoleCollection}

	// core messages for execution & verification
	authorizedRolesMap[cborcodec.CodeExecutionReceipt] = flow.RoleList{flow.RoleExecution}
	authorizedRolesMap[cborcodec.CodeResultApproval] = flow.RoleList{flow.RoleVerification}

	// execution state synchronization
	// NOTE: these messages have been deprecated
	authorizedRolesMap[cborcodec.CodeExecutionStateSyncRequest] = flow.RoleList{}
	authorizedRolesMap[cborcodec.CodeExecutionStateDelta] = flow.RoleList{}

	// data exchange for execution of blocks
	authorizedRolesMap[cborcodec.CodeChunkDataRequest] = flow.RoleList{flow.RoleVerification}
	authorizedRolesMap[cborcodec.CodeChunkDataResponse] = flow.RoleList{flow.RoleExecution}

	// result approvals
	authorizedRolesMap[cborcodec.CodeApprovalRequest] = flow.RoleList{flow.RoleConsensus}
	authorizedRolesMap[cborcodec.CodeApprovalResponse] = flow.RoleList{flow.RoleVerification}

	// generic entity exchange engines
	authorizedRolesMap[cborcodec.CodeEntityRequest] = flow.RoleList{flow.RoleAccess, flow.RoleConsensus, flow.RoleCollection} // only staked access nodes
	authorizedRolesMap[cborcodec.CodeEntityResponse] = flow.RoleList{flow.RoleCollection, flow.RoleExecution}

	// testing
	authorizedRolesMap[cborcodec.CodeEcho] = flow.Roles()

	// dkg
	authorizedRolesMap[cborcodec.CodeDKGMessage] = flow.RoleList{flow.RoleConsensus} // sn nodes for next epoch
}

// AuthorizedSenderValidator using the getIdentity func will check if the role of the sender
// is part of the authorized roles list for the type of message being sent. A node is considered
// to be authorized to send a message if all of the following are true.
// 1. The message type is a known message type (initialized in the authorizedRolesMap).
// 2. The authorized roles list for the message type contains the senders role.
// 3. The node has a weight > 0 and is not ejected
func AuthorizedSenderValidator(log zerolog.Logger, getIdentity func(peer.ID) (*flow.Identity, bool)) MessageValidator {
	log = log.With().
		Str("component", "authorized_sender_validator").
		Logger()

	return func(ctx context.Context, from peer.ID, msg *message.Message) pubsub.ValidationResult {
		identity, ok := getIdentity(from)
		if !ok {
			log.Warn().Str("peer_id", from.String()).Msg("could not get identity of sender")
			return pubsub.ValidationReject
		}

		msgType := msg.Payload[0]

		if err := isAuthorizedNodeRole(identity.Role, msgType); err != nil {
			log.Warn().
				Err(err).
				Str("peer_id", from.String()).
				Str("role", identity.Role.String()).
				Uint8("message_type", msgType).
				Msg("rejecting message")

			return pubsub.ValidationReject
		}

		if err := isActiveNode(identity); err != nil {
			log.Warn().
				Err(err).
				Str("peer_id", from.String()).
				Str("role", identity.Role.String()).
				Uint8("message_type", msgType).
				Msg("rejecting message")

			return pubsub.ValidationReject
		}

		return pubsub.ValidationAccept
	}
}

// isAuthorizedNodeRole checks if a role is authorized to send message type
func isAuthorizedNodeRole(role flow.Role, msgType uint8) error {
	roleList, ok := authorizedRolesMap[msgType]
	if !ok {
		return fmt.Errorf("unknown message type does not match any code from the cbor codec")
	}

	if !roleList.Contains(role) {
		return fmt.Errorf("sender is not authorized to send this message type")
	}

	return nil
}

// isActiveNode checks that the node has a weight > 0 and is not ejected
func isActiveNode(identity *flow.Identity) error {
	if identity.Weight <= 0 {
		return fmt.Errorf("node %s has an invalid weight of %d is not an active node", identity.NodeID, identity.Weight)
	}

	if identity.Ejected {
		return fmt.Errorf("node %s is an ejected node", identity.NodeID)
	}

	return nil
}
