package action

import (
	"fmt"
	"voting/protobuf/voting"

	"voting/processor/model"
	"voting/protobuf/payload"

	"github.com/hyperledger/sawtooth-sdk-go/processor"
)

// AddBallot represents the action of adding new ballot
type AddBallot struct {
	Context   *processor.Context
	Namespace string
	Payload   *payload.OrganizerPayload
}

// Execute create a new user
func (t *AddBallot) Execute() error {
	// check argument
	arg := t.Payload.GetAddBallot()
	if arg == nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Invalid arguments")}
	}

	// get date from payload
	voteID := arg.GetVoteId()
	hashedCode := arg.GetHashedCode()

	// check valid code
	if len(hashedCode) != 64 {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Invalid hashed code")}
	}

	// check ballot exists
	checkBallot, err := model.LoadBallot(hashedCode, voteID, t.Context, t.Namespace)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Failed to load ballot: %v", err)}
	}
	if checkBallot != nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Ballot exists")}
	}

	// check vote exists
	checkVote, err := model.LoadVote(voteID, t.Context, t.Namespace)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Failed to load vote: %v", err)}
	}
	if checkVote == nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Vote not exists")}
	}

	// check time
	if checkVote.GetEndAt() < t.Payload.GetSubmittedAt() {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Invalid time to add ballot")}
	}

	// create object
	ballot := &voting.Ballot{
		VoteId:     voteID,
		HashedCode: hashedCode,
		Choice:     "",
		State:      voting.Ballot_NOT_IN_RESULT,
		CastedAt:   0,
		CreatedAt:  t.Payload.GetSubmittedAt(),
	}

	// save ballot
	err = model.SaveBallot(ballot, t.Context, t.Namespace)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Error saving ballot: %v", err)}
	}

	return nil
}
