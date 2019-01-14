package handler

import (
	"net/http"
	"voting/lib"
	"voting/protobuf/voting"
	"voting/rest/connector"

	"github.com/golang/protobuf/proto"

	"github.com/gin-gonic/gin"
)

// GetBallot returns the targeted vote
func (t *Handler) GetBallot(context *gin.Context) {
	// parse param
	voteID := context.Param("voteID")
	code := context.Param("code")

	addresses := []string{
		connector.GetBallotAddress(lib.Hexdigest256(code), voteID, false),
		connector.GetBallotAddress(lib.Hexdigest256(code), voteID, true),
	}

	for _, address := range addresses {
		// get state data
		stateResponse, err := lib.GetStatesFromRest(&lib.StateOptions{Address: address}, t.RestURL)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get state: " + err.Error()})
			return
		}
		if len(stateResponse.Data) > 0 {
			ballot := &voting.Ballot{}
			err = proto.Unmarshal(stateResponse.Data[0], ballot)
			if err != nil {
				context.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode state: " + err.Error()})
				return
			}
			// success
			context.JSON(http.StatusOK, gin.H{"data": ballot})
			return
		}
	}

	// fail
	context.JSON(http.StatusOK, gin.H{"data": nil})
	return
}
