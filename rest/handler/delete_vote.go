package handler

import (
	"net/http"
	"time"
	"voting/rest/connector"
	"voting/protobuf/payload"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/batch_pb2"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/transaction_pb2"

	"github.com/gin-gonic/gin"
)

// DeleteVoteRequest represents the format of request
type DeleteVoteRequest struct {
	PrivateKey string `json:"private_key" binding:"required"`
	VoteID     string `json:"vote_id" binding:"required"`
}

// DeleteVote submit transaction to delete a vote
func (t *Handler) DeleteVote(context *gin.Context) {
	// parse json
	var request DeleteVoteRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content: " + err.Error()})
		return
	}

	// define payload
	payload := &payload.OrganizerPayload{
		Action:      payload.OrganizerPayload_DELETE_VOTE,
		SubmittedAt: time.Now().Unix(),
		DeleteVote: &payload.OrganizerPayload_DeleteVoteData{
			VoteId: request.VoteID,
		},
	}

	// serialize payload
	payloadBytes, err := proto.Marshal(payload)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Failed to serialize: " + err.Error()})
		return
	}

	// get related address
	address := connector.GetVoteAddress(request.VoteID)

	// signer of requester
	signer, err := connector.NewSigner(request.PrivateKey)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid private key: " + err.Error()})
		return
	}

	// submit transaction
	transaction, err := connector.NewTransaction("voting-organizer", payloadBytes, []string{address}, []string{address}, signer)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create transaction: " + err.Error()})
		return
	}

	batch, err := connector.NewBatch([]*transaction_pb2.Transaction{transaction}, signer)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create batch: " + err.Error()})
		return
	}

	batchIDs, err := connector.SubmitBatches([]*batch_pb2.Batch{batch}, signer, t.RestURL)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Failed to submit batches: " + err.Error()})
		return
	}

	// success
	context.JSON(http.StatusOK, gin.H{
		"action":    "deleteVote",
		"batch_ids": batchIDs,
		"data": gin.H{
			"vote_id": payload.GetDeleteVote().GetVoteId(),
		},
		"submitted_at": payload.GetSubmittedAt(),
	})
	return
}
