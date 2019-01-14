package action

import (
	"fmt"

	"voting/processor/model"
	"voting/protobuf/payload"

	"github.com/hyperledger/sawtooth-sdk-go/processor"
)

// DeleteVote represents the action of deleting vote
type DeleteVote struct {
	Context   *processor.Context
	Namespace string
	Payload   *payload.OrganizerPayload
}

// Execute delete a vote
func (t *DeleteVote) Execute() error {
	// check argument
	arg := t.Payload.GetDeleteVote()
	if arg == nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Invalid arguments")}
	}

	// get target
	voteID := arg.GetVoteId()
	vote, err := model.LoadVote(voteID, t.Context, t.Namespace)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Failed to load vote: %v", err)}
	}
	if vote == nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Vote not exists")}
	}

	// delete vote
	err = model.DeleteVote(vote, t.Context, t.Namespace)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Failed to delete vote: %v", err)}
	}

	return nil
}
