package handler

import (
	"net/http"
	"voting/lib"
	"voting/protobuf/voting"
	"voting/rest/connector"

	"github.com/golang/protobuf/proto"

	"github.com/gin-gonic/gin"
)

// GetVote returns the targeted vote
func (t *Handler) GetVote(context *gin.Context) {
	// parse param
	voteID := context.Param("voteID")

	address := connector.GetVoteAddress(voteID)

	// get state data
	stateResponse, err := lib.GetStatesFromRest(&lib.StateOptions{Address: address}, t.RestURL)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get state: " + err.Error()})
		return
	}
	if len(stateResponse.Data) > 0 {
		vote := &voting.Vote{}
		err = proto.Unmarshal(stateResponse.Data[0], vote)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode state: " + err.Error()})
			return
		}
		// success
		context.JSON(http.StatusOK, gin.H{"data": vote})
		return
	}

	// not data
	context.JSON(http.StatusOK, gin.H{"data": nil})
	return
}
