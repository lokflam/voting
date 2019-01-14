package action

import (
	"fmt"
	"voting/lib"
	"voting/protobuf/voting"

	"voting/processor/model"
	"voting/protobuf/payload"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
)

// CountBallot represents the action of counting ballots and generating result
type CountBallot struct {
	Context   *processor.Context
	Namespace string
	Rest      string
	Payload   *payload.OrganizerPayload
}

// Execute create a new user
func (t *CountBallot) Execute() error {
	// check argument
	arg := t.Payload.GetCountBallot()
	if arg == nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Invalid arguments")}
	}

	// get date from payload
	voteID := arg.GetVoteId()

	// check vote
	vote, err := model.LoadVote(voteID, t.Context, t.Namespace)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Failed to load vote: %v", err)}
	}
	if vote == nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Vote not exists: %v", err)}
	}

	// get previous count
	stateResponse, err := lib.GetStatesFromRest(&lib.StateOptions{
		Address: model.GetResultAddressPrefix(voteID, t.Namespace),
		Limit:   1,
		Reverse: true,
	}, t.Rest)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Failed to get result states: %v", err)}
	}

	// use previous result if exists
	result := &voting.Result{VoteId: voteID}
	prevResult := &voting.Result{}
	if len(stateResponse.Data) > 0 {
		err = proto.Unmarshal(stateResponse.Data[0], prevResult)
		if err != nil {
			return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Failed to unmarshal result: %v", err)}
		}
		result.Total = prevResult.Total
		result.Casted = prevResult.Casted
		result.Counts = prevResult.Counts
	} else {
		result.Total = 0
		result.Casted = 0
		result.Counts = make(map[string]uint32)
	}
	result.CreatedAt = t.Payload.GetSubmittedAt()

	// init non existing count
	for _, candidate := range vote.GetCandidates() {
		if _, ok := result.Counts[candidate.GetCode()]; !ok {
			result.Counts[candidate.GetCode()] = 0
		}
	}

	// count ballot
	options := &lib.StateOptions{
		Address: model.GetBallotAddressPrefix(voteID, false, t.Namespace), // only get ballots that are not counted
		Limit:   100,
		Start:   "",
		Head:    "",
		Reverse: false,
	}
	for countedAll := false; !countedAll; {
		// get ballot states
		stateResponse, err = lib.GetStatesFromRest(options, t.Rest)
		if err != nil {
			return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Failed to get ballot states: %v", err)}
		}

		// check progress
		if stateResponse.NextPosition != "" {
			options.Start = stateResponse.NextPosition
			options.Head = stateResponse.Head
		} else {
			countedAll = true
		}

		// add ballot to count
		for _, data := range stateResponse.Data {
			// decode ballot
			ballot := &voting.Ballot{}
			err = proto.Unmarshal(data, ballot)
			if err != nil {
				return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Failed to unmarshal ballot: %v", err)}
			}

			modified := false

			// increment total ballot, ignore ballots that are created after this transaction is submitted
			if ballot.State == voting.Ballot_NOT_IN_RESULT && ballot.GetCreatedAt() < t.Payload.GetSubmittedAt() {
				result.Total = result.Total + 1
				ballot.State = voting.Ballot_IN_RESULT_TOTAL
				modified = true
			}

			// increment casted and counts, ignore ballots that are casted after this transaction is submitted
			if ballot.GetChoice() != "" && ballot.GetCastedAt() < t.Payload.GetSubmittedAt() {
				result.Casted = result.Casted + 1
				ballot.State = voting.Ballot_IN_RESULT_CASTED
				if _, ok := result.Counts[ballot.GetChoice()]; ok {
					result.Counts[ballot.GetChoice()] = result.Counts[ballot.GetChoice()] + 1
				}
				modified = true
			}

			// save ballot
			if modified {
				err = model.SaveBallot(ballot, t.Context, t.Namespace)
				if err != nil {
					return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Failed to save ballot: %v", err)}
				}
			}
		}
	}

	// save result
	err = model.SaveResult(result, t.Context, t.Namespace)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Failed to save result: %v", err)}
	}

	return nil
}
