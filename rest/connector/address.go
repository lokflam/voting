package connector

import (
	"voting/lib"
)

// GetBallotAddressPrefix returns prefix of ballot address
func GetBallotAddressPrefix(voteID string) string {
	return GetNamespace() + "01" + lib.Hexdigest256(voteID)[:20]
}

// GetBallotAddress returns full ballot address
func GetBallotAddress(hashedCode string, voteID string, counted bool) string {
	// format: namespace(6) + ballot(2) + voteID(20) + counted(2) + hashedCode(40)
	countFlag := "00"
	if counted {
		countFlag = "01"
	}
	return GetBallotAddressPrefix(voteID) + countFlag + hashedCode[:40]
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
