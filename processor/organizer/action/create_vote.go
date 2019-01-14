package action

import (
	"fmt"

	"voting/processor/model"
	"voting/protobuf/payload"

	"github.com/hyperledger/sawtooth-sdk-go/processor"
)

// CreateVote represents the action of creating vote
type CreateVote struct {
	Context   *processor.Context
	Namespace string
	Payload   *payload.OrganizerPayload
}

// Execute create a new vote
func (t *CreateVote) Execute() error {
	// check argument
	arg := t.Payload.GetCreateVote()
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
	if checkVote != nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Vote exists")}
	}

	// assign time
	vote.CreatedAt = t.Payload.GetSubmittedAt()
	vote.UpdatedAt = t.Payload.GetSubmittedAt()

	// save vote
	err = model.SaveVote(vote, t.Context, t.Namespace)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Error saving vote: %v", err)}
	}

	return nil
}
