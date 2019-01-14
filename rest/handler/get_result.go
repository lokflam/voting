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

// GetResult returns the result of targeted vote
func (t *Handler) GetResult(context *gin.Context) {
	// parse param
	voteID := context.Param("voteID")
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
	reverse := context.DefaultQuery("reverse", "true")

	// get state data
	stateResponse, err := lib.GetStatesFromRest(&lib.StateOptions{
		Address: connector.GetResultAddressPrefix(voteID),
		Limit:   limit,
		Head:    head,
		Start:   start,
		Reverse: reverse != "false",
	}, t.RestURL)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get state: " + err.Error()})
		return
	}

	var results []*voting.Result
	for _, data := range stateResponse.Data {
		result := &voting.Result{}
		err = proto.Unmarshal(data, result)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode state: " + err.Error()})
			return
		}
		results = append(results, result)
	}

	// success
	context.JSON(http.StatusOK, gin.H{
		"data":          results,
		"head":          stateResponse.Head,
		"start":         stateResponse.Start,
		"next_position": stateResponse.NextPosition,
	})
	return
}
