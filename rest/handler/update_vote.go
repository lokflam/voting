package handler

import (
	"net/http"
	"time"
	"voting/rest/connector"
	"voting/protobuf/payload"
	"voting/protobuf/voting"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/batch_pb2"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/transaction_pb2"

	"github.com/gin-gonic/gin"
)

// UpdateVoteRequest represents the format of request
type UpdateVoteRequest struct {
	PrivateKey string       `json:"private_key" binding:"required"`
	Vote       *voting.Vote `json:"vote" binding:"required"`
}

// UpdateVote submit transaction to update a vote
func (t *Handler) UpdateVote(context *gin.Context) {
	// parse json
	var request UpdateVoteRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content: " + err.Error()})
		return
	}

	// define payload
	payload := &payload.OrganizerPayload{
		Action:      payload.OrganizerPayload_UPDATE_VOTE,
		SubmittedAt: time.Now().Unix(),
		UpdateVote: &payload.OrganizerPayload_UpdateVoteData{
			Vote: request.Vote,
		},
	}

	// serialize payload
	payloadBytes, err := proto.Marshal(payload)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Failed to serialize: " + err.Error()})
		return
	}

	// get related address
	address := connector.GetVoteAddress(payload.GetUpdateVote().GetVote().GetId())

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
		"action":    "updateVote",
		"batch_ids": batchIDs,
		"data": gin.H{
			"vote_id": payload.GetUpdateVote().GetVote().GetId(),
		},
		"submitted_at": payload.GetSubmittedAt(),
	})
	return
}
