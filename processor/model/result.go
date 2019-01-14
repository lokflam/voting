package model

import (
	"fmt"
	"strconv"
	"voting/lib"
	"voting/protobuf/voting"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
)

// SaveResult save a Result to the blockchain
func SaveResult(result *voting.Result, context *processor.Context, namespace string) error {
	// generate address
	address := GetResultAddress(result.GetVoteId(), result.GetCreatedAt(), namespace)

	// marshal data
	data, err := proto.Marshal(result)
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

// GetResultAddressPrefix returns prefix of result address
func GetResultAddressPrefix(id string, namespace string) string {
	// format: namespace(6) + result(2) + vote_id(46) + timestamp(16)
	return namespace + "02" + lib.Hexdigest256(id)[:46]
}

// GetResultAddress returns result address
func GetResultAddress(id string, timestamp int64, namespace string) string {
	// format: namespace(6) + result(2) + vote_id(46) + timestamp(16)
	return GetResultAddressPrefix(id, namespace) + fmt.Sprintf("%016s", strconv.FormatInt(timestamp, 16))
}
