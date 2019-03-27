package handler

import (
	"net/http"
	"time"
	"voting/lib"
	"voting/protobuf/payload"
	"voting/rest/connector"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/batch_pb2"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/transaction_pb2"
	"github.com/hyperledger/sawtooth-sdk-go/signing"

	"github.com/gin-gonic/gin"
)

// CastBallotRequest represents the format of request
type CastBallotRequest struct {
	VoteID string `json:"vote_id" binding:"required"`
	Choice string `json:"choice" binding:"required"`
	Code   string `json:"code" binding:"required"`
}

// CastBallot submit transaction to cast a ballot
func (t *Handler) CastBallot(context *gin.Context) {
	// parse json
	var request CastBallotRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content: " + err.Error()})
		return
	}

	// define payload
	payload := &payload.VoterPayload{
		Action:      payload.VoterPayload_CAST_BALLOT,
		SubmittedAt: time.Now().Unix(),
		CastBallot: &payload.VoterPayload_CastBallotData{
			VoteId: request.VoteID,
			Choice: request.Choice,
			Code:   request.Code,
		},
	}

	// serialize payload
	payloadBytes, err := proto.Marshal(payload)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Failed to serialize: " + err.Error()})
		return
	}

	// get related address
	ballotAddresses := []string{
		connector.GetBallotAddress(lib.Hexdigest256(request.Code), request.VoteID),
		connector.GetBallotLogAddress(lib.Hexdigest256(request.Code), request.VoteID, payload.GetSubmittedAt()),
	}
	voteAddress := connector.GetVoteAddress(request.VoteID)

	// use a new random requester
	signerContext := signing.CreateContext("secp256k1")
	privateKey := signerContext.NewRandomPrivateKey()
	signer, err := connector.NewSigner(privateKey.AsHex())
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid private key: " + err.Error()})
		return
	}

	// submit transaction
	transaction, err := connector.NewTransaction("voting-voter", payloadBytes, append(ballotAddresses, voteAddress), ballotAddresses, signer)
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
		"action":    "castBallot",
		"batch_ids": batchIDs,
		"data": gin.H{
			"vote_id":     request.VoteID,
			"hashed_code": lib.Hexdigest256(request.Code),
		},
		"submitted_at": payload.GetSubmittedAt(),
	})
	return
}
