package connector

import (
	"fmt"
	"strconv"
	"voting/lib"
)

// GetBallotAddress returns full ballot address
func GetBallotAddress(hashedCode string, voteID string) string {
	// format: namespace(6) + ballot(2) + voteID(16) + hashedCode(44)
	return GetNamespace() + "01" + lib.Hexdigest256(voteID)[:16] + hashedCode[:44]
}

// GetBallotLogAddressPrefix returns prefix of ballot log address
func GetBallotLogAddressPrefix(voteID string) string {
	return GetNamespace() + "11" + lib.Hexdigest256(voteID)[:16]
}

// GetBallotLogAddress returns full ballot log address
func GetBallotLogAddress(hashedCode string, voteID string, timestamp int64) string {
	// format: namespace(6) + ballotLog(2) + voteID(16) + timestamp(16) + hashedCode(30)
	return GetBallotLogAddressPrefix(voteID) + fmt.Sprintf("%016s", strconv.FormatInt(timestamp, 16)) + hashedCode[:30]
}

// GetVoteAddress returns address of vote
func GetVoteAddress(voteID string) string {
	return GetVoteAddressPrefix() + lib.Hexdigest256(voteID)[:62]
}

// GetVoteAddressPrefix returns address prefix of vote
func GetVoteAddressPrefix() string {
	return GetNamespace() + "00"
}

// GetResultAddressPrefix returns prefix of result address
func GetResultAddressPrefix(id string) string {
	// format: namespace(6) + result(2) + vote_id(46) + timestamp(16)
	return GetNamespace() + "02" + lib.Hexdigest256(id)[:46]
}
