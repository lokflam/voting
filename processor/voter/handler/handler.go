package handler

import (
	"fmt"
	"time"

	"voting/lib"
	cmd "voting/processor/command"
	act "voting/processor/voter/action"
	pl "voting/protobuf/payload"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/sawtooth-sdk-go/logging"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/processor_pb2"
)

var logger = logging.Get()
var namespace = lib.Hexdigest512("voting")[:6]
var acceptedDelay int64 = 60 * 5 // 5 mins

// Handler represents handler for application
type Handler struct {
}

// FamilyName returns name
func (t *Handler) FamilyName() string {
	return "voting-voter"
}

// FamilyVersions returns version
func (t *Handler) FamilyVersions() []string {
	return []string{"1.0"}
}

// Namespaces returns namespace
func (t *Handler) Namespaces() []string {
	return []string{namespace}
}

// Apply handles request
func (t *Handler) Apply(request *processor_pb2.TpProcessRequest, context *processor.Context) error {
	// decode payload
	payload := &pl.VoterPayload{}
	err := proto.Unmarshal(request.GetPayload(), payload)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Invalid payload: %v", err)}
	}

	// check submitted time
	if payload.GetSubmittedAt() > time.Now().Unix() || payload.GetSubmittedAt()+acceptedDelay < time.Now().Unix() {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Accepted session ended")}
	}

	// list of valid actions
	actions := map[pl.VoterPayload_Action]cmd.Command{
		pl.VoterPayload_CAST_BALLOT: &act.CastBallot{
			Context:   context,
			Namespace: namespace,
			Payload:   payload,
		},
	}

	// check action exists
	action, ok := actions[payload.GetAction()]
	if !ok || action == nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Invalid action")}
	}

	// run action
	return action.Execute()
}
