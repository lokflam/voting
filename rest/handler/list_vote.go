package handler

import (
	"net/http"
	"strconv"
	"voting/lib"
	"voting/protobuf/voting"
	"voting/rest/connector"

	"github.com/golang/protobuf/proto"

	"github.com/gin-gonic/gin"
)

// ListVote returns the targeted vote
func (t *Handler) ListVote(context *gin.Context) {
	// parse param
	limit, err := strconv.Atoi(context.DefaultQuery("limit", "100"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit: " + err.Error()})
		return
	}
	if limit < 1 || limit > 100 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Limit must be within 1 to 100"})
		return
	}
	head := context.DefaultQuery("head", "")
	start := context.DefaultQuery("start", "")

	// get state data
	stateResponse, err := lib.GetStatesFromRest(&lib.StateOptions{
		Address: connector.GetVoteAddressPrefix(),
		Limit:   limit,
		Head:    head,
		Start:   start,
		Reverse: false,
	}, t.RestURL)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get state: " + err.Error()})
		return
	}

	var votes []*voting.Vote
	for _, data := range stateResponse.Data {
		vote := &voting.Vote{}
		err = proto.Unmarshal(data, vote)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode state: " + err.Error()})
			return
		}
		votes = append(votes, vote)
	}

	// success
	context.JSON(http.StatusOK, gin.H{
		"data":          votes,
		"head":          stateResponse.Head,
		"start":         stateResponse.Start,
		"next_position": stateResponse.NextPosition,
	})
	return
}
