package connector

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"voting/lib"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/batch_pb2"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/transaction_pb2"
	"github.com/hyperledger/sawtooth-sdk-go/signing"
)

var namespace = "voting"
var families = map[string]string{
	"voting-organizer": "1.0",
	"voting-voter":     "1.0",
}

// GetNamespace returns namespace
func GetNamespace() string {
	return lib.Hexdigest512(namespace)[:6]
}

// NewTransaction creates a new transaction
func NewTransaction(family string, payload []byte, inputs []string, outputs []string, signer *signing.Signer) (*transaction_pb2.Transaction, error) {
	if _, ok := families[family]; !ok {
		return nil, fmt.Errorf("Transaction family not exists")
	}

	txnHeader := &transaction_pb2.TransactionHeader{
		FamilyName:       family,
		FamilyVersion:    families[family],
		Dependencies:     []string{},
		Inputs:           inputs,
		Outputs:          outputs,
		SignerPublicKey:  signer.GetPublicKey().AsHex(),
		BatcherPublicKey: signer.GetPublicKey().AsHex(),
		PayloadSha512:    lib.Hexdigest512(string(payload)),
	}

	txnHeaderBytes, err := proto.Marshal(txnHeader)
	if err != nil {
		return nil, fmt.Errorf("Failed to serialize transaction header: %v", err)
	}

	signatureBytes := signer.Sign(txnHeaderBytes)
	signature := strings.ToLower(hex.EncodeToString(signatureBytes))

	return &transaction_pb2.Transaction{
		Header:          txnHeaderBytes,
		HeaderSignature: signature,
		Payload:         payload,
	}, nil
}

// NewBatch creates a new batch containing multiple transactions
func NewBatch(txns []*transaction_pb2.Transaction, signer *signing.Signer) (*batch_pb2.Batch, error) {
	txnIDs := []string{}
	for _, txn := range txns {
		txnIDs = append(txnIDs, txn.GetHeaderSignature())
	}

	batchHeader := &batch_pb2.BatchHeader{
		SignerPublicKey: signer.GetPublicKey().AsHex(),
		TransactionIds:  txnIDs,
	}

	batchHeaderBytes, err := proto.Marshal(batchHeader)
	if err != nil {
		return nil, fmt.Errorf("Failed to serialize batch header: %v", err)
	}

	batchSignatureBytes := signer.Sign(batchHeaderBytes)
	batchSignature := strings.ToLower(hex.EncodeToString(batchSignatureBytes))

	return &batch_pb2.Batch{
		Header:          batchHeaderBytes,
		HeaderSignature: batchSignature,
		Transactions:    txns,
	}, nil
}

// SubmitBatches submits array of batches and return comma seperated batch ids string for tracing
func SubmitBatches(batches []*batch_pb2.Batch, signer *signing.Signer, restURL string) (string, error) {
	batchList := &batch_pb2.BatchList{
		Batches: batches,
	}

	batchListBytes, err := proto.Marshal(batchList)
	if err != nil {
		return "", fmt.Errorf("Failed to serialize batch list: %v", err)
	}

	response, err := http.Post(restURL+"/batches", "application/octet-stream", bytes.NewBuffer(batchListBytes))
	if err != nil {
		return "", fmt.Errorf("Fatal error: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		bodyBytes, _ := ioutil.ReadAll(response.Body)
		return "", fmt.Errorf("Failed to submit: %v", string(bodyBytes))
	}

	var data map[string]interface{}
	if err = json.NewDecoder(response.Body).Decode(&data); err != nil {
		return "", fmt.Errorf("Fatal error: %v", err)
	}

	if _, ok := data["link"]; !ok {
		return "", fmt.Errorf("Missing 'link' in response")
	}

	link, ok := data["link"].(string)
	if !ok {
		return "", fmt.Errorf("Invalid 'link' in response")
	}

	u, err := url.Parse(link)
	if err != nil {
		return "", fmt.Errorf("Invalid 'link' in response")
	}

	query, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return "", fmt.Errorf("Invalid 'link' in response")
	}

	if ids, ok := query["id"]; !ok || len(ids) < 1 {
		return "", fmt.Errorf("Invalid 'link' in response")
	}

	return query["id"][0], nil
}

// NewSigner returns a signer of the private key
func NewSigner(privateKeyString string) (*signing.Signer, error) {
	privateKeyBytes, err := hex.DecodeString(privateKeyString)
	if err != nil {
		return nil, err
	}
	privateKey := signing.NewSecp256k1PrivateKey(privateKeyBytes)
	context := signing.CreateContext("secp256k1")
	signer := signing.NewCryptoFactory(context).NewSigner(privateKey)
	return signer, nil
}
