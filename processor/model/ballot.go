package model

import (
	"fmt"
	"voting/lib"
	"voting/protobuf/voting"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
)

// SaveBallot save a Ballot to the blockchain
func SaveBallot(ballot *voting.Ballot, context *processor.Context, namespace string) error {
	// marshal data
	data, err := proto.Marshal(ballot)
	if err != nil {
		return fmt.Errorf("Failed to serialize: %v", err)
	}

	// delete repeated ballot
	address = GetBallotAddress(ballot.GetHashedCode(), ballot.GetVoteId(), !(ballot.State == voting.Ballot_IN_RESULT_CASTED), namespace)
	_, err = context.DeleteState([]string{address})
	if err != nil {
		return fmt.Errorf("Failed to delete state: %v", err)
	}

	// generate address
	address := GetBallotAddress(ballot.GetHashedCode(), ballot.GetVoteId(), (ballot.State == voting.Ballot_IN_RESULT_CASTED), namespace)

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

// LoadBallot search and return a Ballot
func LoadBallot(hashedCode string, voteID string, context *processor.Context, namespace string) (*voting.Ballot, error) {
	// get address
	addresses := []string{
		GetBallotAddress(hashedCode, voteID, false, namespace),
		GetBallotAddress(hashedCode, voteID, true, namespace),
	}

	for _, address := range addresses {
		// get data from states
		results, err := context.GetState([]string{address})
		if err != nil {
			return nil, fmt.Errorf("Failed to get state: %v", err)
		}

		// check data valid
		if len(string(results[address])) > 0 {
			ballot := &voting.Ballot{}
			err := proto.Unmarshal(results[address], ballot)
			if err != nil {
				return nil, fmt.Errorf("Failed to unmarshal ballot: %v", err)
			}
			return ballot, nil
		}
	}

	// return nil if no data
	return nil, nil
}

// GetBallotAddressPrefix returns prefix of ballot address
func GetBallotAddressPrefix(voteID string, counted bool, namespace string) string {
	countFlag := "00"
	if counted {
		countFlag = "01"
	}
	return namespace + "01" + lib.Hexdigest256(voteID)[:20] + countFlag
}

// GetBallotAddress returns full ballot address
func GetBallotAddress(hashedCode string, voteID string, counted bool, namespace string) string {
	// format: namespace(6) + ballot(2) + voteID(20) + counted(2) + hashedCode(40)
	return GetBallotAddressPrefix(voteID, counted, namespace) + hashedCode[:40]
}
