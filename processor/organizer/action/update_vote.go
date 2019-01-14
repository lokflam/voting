package action

import (
	"fmt"

	"voting/processor/model"
	"voting/protobuf/payload"

	"github.com/hyperledger/sawtooth-sdk-go/processor"
)

// UpdateVote represents the action of updating vote
type UpdateVote struct {
	Context   *processor.Context
	Namespace string
	Payload   *payload.OrganizerPayload
}

// Execute update a vote
func (t *UpdateVote) Execute() error {
	// check argument
	arg := t.Payload.GetUpdateVote()
	if arg == nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Invalid arguments")}
	}

	// get vote object from payload
	vote := arg.GetVote()

	// validate
	err := model.ValidateVote(vote)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Invalid vote: %v", err)}
	}

	// check exists
	checkVote, err := model.LoadVote(vote.GetId(), t.Context, t.Namespace)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Failed to load vote: %v", err)}
	}
	if checkVote == nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Vote not exists")}
	}

	// update time
	vote.UpdatedAt = t.Payload.GetSubmittedAt()
	vote.CreatedAt = checkVote.GetCreatedAt()

	// save vote
	err = model.SaveVote(vote, t.Context, t.Namespace)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Error saving vote: %v", err)}
	}

	return nil
}
