package action

import (
	"fmt"
	"voting/lib"
	"voting/protobuf/voting"

	"voting/processor/model"
	"voting/protobuf/payload"

	"github.com/hyperledger/sawtooth-sdk-go/processor"
)

// CastBallot represents the action of casting a Ballot
type CastBallot struct {
	Context   *processor.Context
	Namespace string
	Payload   *payload.VoterPayload
}

// Execute create a new user
func (t *CastBallot) Execute() error {
	// check argument
	arg := t.Payload.GetCastBallot()
	if arg == nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Invalid arguments")}
	}

	voteID := arg.GetVoteId()
	choice := arg.GetChoice()
	code := arg.GetCode()
	hashedCode := lib.Hexdigest256(code)

	// check argument
	if voteID == "" {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Missing vote ID")}
	}
	if choice == "" {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Missing choice")}
	}
	if code == "" {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Missing code")}
	}

	// get ballot and check
	ballot, err := model.LoadBallot(hashedCode, voteID, t.Context, t.Namespace)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Failed to load ballot: %v", err)}
	}
	if ballot == nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Ballot not exists")}
	}
	if ballot.GetChoice() != "" {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Ballot already casted")}
	}
	if ballot.GetCreatedAt() > t.Payload.GetSubmittedAt() {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Invalid cast time")}
	}

	// validate vote
	vote, err := model.LoadVote(voteID, t.Context, t.Namespace)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Failed to load vote: %v", err)}
	}
	if vote == nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Vote not exists")}
	}

	// check time
	if vote.GetStartAt() > t.Payload.GetSubmittedAt() || vote.GetEndAt() < t.Payload.GetSubmittedAt() {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Invalid time to vote")}
	}

	// check choice
	validChoice := false
	for _, candidate := range vote.GetCandidates() {
		if candidate.GetCode() == choice && candidate.GetStatus() != voting.CandidateStatus_DISQUALIFIED {
			validChoice = true
		}
	}
	if !validChoice {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Invalid choice")}
	}

	// modify ballot data
	ballot.Choice = choice
	ballot.CastedAt = t.Payload.GetSubmittedAt()

	// save ballot
	err = model.SaveBallot(ballot, t.Context, t.Namespace)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Error saving ballot: %v", err)}
	}

	// create log
	log := &voting.BallotLog{
		VoteId:      ballot.GetVoteId(),
		HashedCode:  ballot.GetHashedCode(),
		Choice:      ballot.GetChoice(),
		ProcessedAt: 0,
		LoggedAt:    t.Payload.GetSubmittedAt(),
	}

	// save log
	err = model.SaveBallotLog(log, t.Context, t.Namespace)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Error saving ballot log: %v", err)}
	}

	return nil
}
