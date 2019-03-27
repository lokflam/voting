package model

import (
	"fmt"
	"strconv"
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

	// generate address
	address := GetBallotAddress(ballot.GetHashedCode(), ballot.GetVoteId(), namespace)

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

// SaveBallotLog save a BallotLog to the blockchain
func SaveBallotLog(log *voting.BallotLog, context *processor.Context, namespace string) error {
	// marshal data
	data, err := proto.Marshal(log)
	if err != nil {
		return fmt.Errorf("Failed to serialize: %v", err)
	}

	// generate address
	address := GetBallotLogAddress(log.GetHashedCode(), log.GetVoteId(), log.GetLoggedAt(), namespace)

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
	address := GetBallotAddress(hashedCode, voteID, namespace)

	// get data from state
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

	// return nil if no data
	return nil, nil
}

// GetBallotAddress returns full ballot address
func GetBallotAddress(hashedCode string, voteID string, namespace string) string {
	// format: namespace(6) + ballot(2) + voteID(16) + hashedCode(46)
	return namespace + "01" + lib.Hexdigest256(voteID)[:16] + hashedCode[:46]
}

// GetBallotLogAddressPrefix returns prefix of ballot log address
func GetBallotLogAddressPrefix(voteID string, namespace string) string {
	return namespace + "11" + lib.Hexdigest256(voteID)[:16]
}

// GetBallotLogAddress returns full ballot log address
func GetBallotLogAddress(hashedCode string, voteID string, timestamp int64, namespace string) string {
	// format: namespace(6) + ballotLog(2) + voteID(16) + timestamp(16) + hashedCode(30)
	return GetBallotLogAddressPrefix(voteID, namespace) + fmt.Sprintf("%016s", strconv.FormatInt(timestamp, 16)) + hashedCode[:30]
}
