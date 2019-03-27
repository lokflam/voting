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
	Context       *processor.Context
	Namespace     string
	Rest          string
	AcceptedDelay int64
	Payload       *payload.OrganizerPayload
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
		result.Counts = []*voting.Result_Count{}
	}
	result.CreatedAt = t.Payload.GetSubmittedAt()

	// init non existing count
	for _, candidate := range vote.GetCandidates() {
		for _, count := range result.GetCounts() {
			if count.GetCandidate() == candidate.GetCode() {
				continue
			}
		}
		result.Counts = append(result.Counts, &voting.Result_Count{
			Candidate: candidate.GetCode(),
			Count:     0,
		})
	}

	// count ballot
	// get recent ballots
	options := &lib.StateOptions{
		Address: model.GetBallotLogAddressPrefix(voteID, t.Namespace),
		Limit:   100,
		Start:   "",
		Head:    "",
		Reverse: true,
	}
	for countedAll := false; !countedAll; {
		// get ballot log states
		stateResponse, err = lib.GetStatesFromRest(options, t.Rest)
		if err != nil {
			return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Failed to get ballot log states: %v", err)}
		}

		// check progress
		if stateResponse.NextPosition != "" {
			options.Start = stateResponse.NextPosition
			options.Head = stateResponse.Head
		} else {
			countedAll = true
		}

		// add ballot log to count
		for _, data := range stateResponse.Data {
			// decode ballot log
			log := &voting.BallotLog{}
			err = proto.Unmarshal(data, log)
			if err != nil {
				return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Failed to unmarshal ballot log: %v", err)}
			}

			// ignore ballots that are casted after this transaction is submitted
			if log.GetLoggedAt() > t.Payload.GetSubmittedAt() {
				continue
			}

			// only process ballots that are not processed in last count
			if log.GetLoggedAt() < prevResult.GetCreatedAt()-t.AcceptedDelay {
				break
			}

			// skip counted ballot
			if log.GetProcessedAt() > 0 {
				continue
			}

			// increment counts
			if log.GetChoice() != "" {
				result.Casted = result.Casted + 1
				for _, count := range result.GetCounts() {
					if count.GetCandidate() == log.GetChoice() {
						count.Count = count.GetCount() + 1
					}
				}
			} else {
				result.Total = result.Total + 1
			}

			log.ProcessedAt = t.Payload.GetSubmittedAt()
			err = model.SaveBallotLog(log, t.Context, t.Namespace)
			if err != nil {
				return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Failed to save ballot log: %v", err)}
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
