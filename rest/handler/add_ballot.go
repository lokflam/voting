package handler

import (
	"net/http"
	"strconv"
	"time"
	"voting/lib"
	"voting/protobuf/payload"
	"voting/rest/connector"

	"github.com/hyperledger/sawtooth-sdk-go/protobuf/batch_pb2"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/transaction_pb2"

	"github.com/golang/protobuf/proto"

	"github.com/gin-gonic/gin"
)

// AddBallotRequest represents the format of request
type AddBallotRequest struct {
	PrivateKey string   `json:"private_key" binding:"required"`
	VoteID     string   `json:"vote_id" binding:"required"`
	Quantity   int      `json:"quantity" binding:"-"`
	Codes      []string `json:"codes" binding:"-"`
}

// AddBallot submit transaction to add new empty ballots
func (t *Handler) AddBallot(context *gin.Context) {
	// parse json
	var request AddBallotRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content: " + err.Error()})
		return
	}

	var codes []string // storing generated codes
	if len(request.Codes) > 1 {
		if len(request.Codes) > 1000 {
			context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content: 'codes' length must be smaller than 1000"})
			return
		}
		codes = request.Codes
	} else {
		if request.Quantity < 0 || request.Quantity > 1000 {
			context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content: 'quantity' must be 1 to 1000"})
			return
		}

		for i := 0; i < request.Quantity; i++ {
			codes = append(codes, lib.GenerateUUID(strconv.Itoa(i)))
		}
	}

	// signer of requester
	signer, err := connector.NewSigner(request.PrivateKey)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid private key: " + err.Error()})
		return
	}

	var transactions []*transaction_pb2.Transaction         // storing generated transactions
	submittedAt := time.Now().Unix()                        // define submitted time for all ballots
	voteAddress := connector.GetVoteAddress(request.VoteID) // define vote address for all ballots

	// create transaction for each ballot
	for _, code := range codes {
		hashedCode := lib.Hexdigest256(code)

		// define payload
		payload := &payload.OrganizerPayload{
			Action:      payload.OrganizerPayload_ADD_BALLOT,
			SubmittedAt: submittedAt,
			AddBallot: &payload.OrganizerPayload_AddBallotData{
				VoteId:     request.VoteID,
				HashedCode: hashedCode,
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
			connector.GetBallotAddress(hashedCode, request.VoteID),
			connector.GetBallotLogAddress(hashedCode, request.VoteID, submittedAt),
		}

		// create transaction
		transaction, err := connector.NewTransaction("voting-organizer", payloadBytes, append(ballotAddresses, voteAddress), ballotAddresses, signer)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create transaction: " + err.Error()})
			return
		}

		transactions = append(transactions, transaction)
	}

	// submit transactions
	batch, err := connector.NewBatch(transactions, signer)
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
		"action":    "addBallot",
		"batch_ids": batchIDs,
		"data": gin.H{
			"vote_id": request.VoteID,
		},
		"submitted_at": submittedAt,
		"codes":        codes,
	})
	return
}
