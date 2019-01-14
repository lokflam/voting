package handler

import (
	"net/http"
	"time"
	"voting/rest/connector"
	"voting/lib"
	"voting/protobuf/payload"
	"voting/protobuf/voting"

	"github.com/hyperledger/sawtooth-sdk-go/protobuf/batch_pb2"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/transaction_pb2"

	"github.com/golang/protobuf/proto"

	"github.com/gin-gonic/gin"
)

// CreateVoteRequest represents the format of request
type CreateVoteRequest struct {
	PrivateKey string       `json:"private_key" binding:"required"`
	Vote       *voting.Vote `json:"vote" binding:"required"`
}

// CreateVote submit transaction to create a new vote
func (t *Handler) CreateVote(context *gin.Context) {
	// parse json
	var request CreateVoteRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content: " + err.Error()})
		return
	}

	// assign vote id
	if request.Vote.Id == "" {
		request.Vote.Id = lib.GenerateUUID(request.Vote.GetName())
	}

	// define payload
	payload := &payload.OrganizerPayload{
		Action:      payload.OrganizerPayload_CREATE_VOTE,
		SubmittedAt: time.Now().Unix(),
		CreateVote: &payload.OrganizerPayload_CreateVoteData{
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
	address := connector.GetVoteAddress(payload.GetCreateVote().GetVote().GetId())

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
		"action":    "createVote",
		"batch_ids": batchIDs,
		"data": gin.H{
			"vote_id": payload.GetCreateVote().GetVote().GetId(),
		},
		"submitted_at": payload.GetSubmittedAt(),
	})
	return
}
