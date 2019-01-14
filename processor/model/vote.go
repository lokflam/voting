package model

import (
	"fmt"

	"voting/lib"
	"voting/protobuf/voting"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
)

// ValidateVote validate a vote and return the error
func ValidateVote(vote *voting.Vote) error {
	// validate ID
	if vote.GetId() == "" {
		return fmt.Errorf("Missing 'id'")
	}

	// validate name
	if vote.GetName() == "" {
		return fmt.Errorf("Missing 'name'")
	}

	// validate start_at
	if vote.GetStartAt() == 0 {
		return fmt.Errorf("Missing 'start_at'")
	}

	// validate end_at
	if vote.GetEndAt() == 0 {
		return fmt.Errorf("Missing 'end_at'")
	}
	if vote.GetEndAt() < vote.GetStartAt() {
		return fmt.Errorf("Invalid 'end_at'")
	}

	// validate candidates
	if len(vote.GetCandidates()) == 0 {
		return fmt.Errorf("Missing 'candidates'")
	}
	codes := make(map[string]bool)
	for i, candidate := range vote.GetCandidates() {
		if candidate.GetCode() == "" {
			return fmt.Errorf("Missing candidate %v 'code'", i)
		}
		if candidate.GetName() == "" {
			return fmt.Errorf("Missing candidate %v 'name'", i)
		}
		if _, ok := codes[candidate.GetCode()]; ok {
			return fmt.Errorf("Candidate %v 'code' duplicated", i)
		}
		codes[candidate.GetCode()] = true
	}

	return nil
}

// DeleteVote delete a vote
func DeleteVote(vote *voting.Vote, context *processor.Context, namespace string) error {
	// get address
	address := getVoteAddress(vote.GetId(), namespace)

	// delete vote
	responses, err := context.DeleteState([]string{address})
	if err != nil {
		return fmt.Errorf("Failed to delete state: %v", err)
	}
	if len(responses) == 0 {
		return fmt.Errorf("No state was not deleted")
	}

	return nil
}

// SaveVote save a Vote to the blockchain
func SaveVote(vote *voting.Vote, context *processor.Context, namespace string) error {
	// generate address
	address := getVoteAddress(vote.GetId(), namespace)

	// marshal data
	data, err := proto.Marshal(vote)
	if err != nil {
		return fmt.Errorf("Failed to serialize: %v", err)
	}

	// add data to state
	addresses, err := context.SetState(map[string][]byte{address: data})
	if err != nil {
		return fmt.Errorf("Failed to set state: %v", err)
	}
	if len(addresses) == 0 {
		return fmt.Errorf("No addresses in set response")
	}

	return nil
}

// LoadVote search and return a Vote
func LoadVote(id string, context *processor.Context, namespace string) (*voting.Vote, error) {
	// get address
	address := getVoteAddress(id, namespace)

	// get data from states
	results, err := context.GetState([]string{address})
	if err != nil {
		return nil, fmt.Errorf("Failed to get state: %v", err)
	}

	// check data valid
	if len(string(results[address])) > 0 {
		vote := &voting.Vote{}
		err := proto.Unmarshal(results[address], vote)
		if err != nil {
			return nil, fmt.Errorf("Failed to unmarshal Vote: %v", err)
		}

		return vote, nil
	}

	// return nil if no data
	return nil, nil
}

func getVoteAddress(id string, namespace string) string {
	// format: namespace(6) + vote(2) + vote_id(62)
	return namespace + "00" + lib.Hexdigest256(id)[:62]
}
