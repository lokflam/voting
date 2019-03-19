package handler

import (
	"net/http"
	"voting/lib"
	"voting/protobuf/voting"
	"voting/rest/connector"

	"github.com/golang/protobuf/proto"

	"github.com/gin-gonic/gin"
)

// GetBallotRequest represents the format of request
type GetBallotRequest struct {
	VoteID string `json:"vote_id" binding:"required"`
	Code   string `json:"codes" binding:"required"`
}

// GetBallot returns the targeted vote
func (t *Handler) GetBallot(context *gin.Context) {
	// parse json
	var request GetBallotRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content: " + err.Error()})
		return
	}

	addresses := []string{
		connector.GetBallotAddress(lib.Hexdigest256(request.Code), request.VoteID, false),
		connector.GetBallotAddress(lib.Hexdigest256(request.Code), request.VoteID, true),
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
