package lib

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// StateOptions represents options for getting states
type StateOptions struct {
	Address string
	Limit   int
	Head    string
	Start   string
	Reverse bool
}

// StateResponse responsents response from get state function
type StateResponse struct {
	Data         [][]byte
	Head         string
	Start        string
	NextPosition string
}

// GetStatesFromRest return result of states from rest api
func GetStatesFromRest(options *StateOptions, restURL string) (*StateResponse, error) {
	// construct url
	url := restURL + "/state?address=" + options.Address
	if options.Limit != 0 {
		url = url + "&limit=" + strconv.Itoa(options.Limit)
	}
	if options.Head != "" {
		url = url + "&head=" + options.Head
	}
	if options.Start != "" {
		url = url + "&start=" + options.Start
	}
	if options.Reverse != false {
		url = url + "&reverse"
	}

	// get from rest-api
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Failed to get state: %v", err)
	}
	defer response.Body.Close()

	// decode json
	var responseJSON map[string]interface{}
	if err = json.NewDecoder(response.Body).Decode(&responseJSON); err != nil {
		return nil, fmt.Errorf("Failed to parse response: %v", err)
	}

	// check 'data'
	records, ok := responseJSON["data"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("Failed to get response 'data': %v", responseJSON["data"])
	}

	// define response object
	stateResponse := &StateResponse{}

	// update 'head'
	checkHead, ok := responseJSON["head"].(string)
	if ok {
		stateResponse.Head = checkHead
	}

	// check 'paging'
	paging, ok := responseJSON["paging"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Failed to decode 'paging': %v", err)
	}
	// update 'start'
	checkStart, ok := paging["start"].(string)
	if ok {
		stateResponse.Start = checkStart
	}
	// add 'next_position'
	checkNext, ok := paging["next_position"].(string)
	if ok {
		stateResponse.NextPosition = checkNext
	}

	// extract data
	for _, record := range records {
		data, ok := record.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("Failed to parse 'data': %v", record)
		}
		if _, ok := data["data"]; !ok {
			continue
		}
		payloadBase64, ok := data["data"].(string)
		if !ok {
			return nil, fmt.Errorf("Failed to parse 'data': %v", record)
		}
		payload, err := base64.StdEncoding.DecodeString(payloadBase64)
		if err != nil {
			return nil, fmt.Errorf("Failed to decode 'data': %v", err)
		}
		stateResponse.Data = append(stateResponse.Data, payload)
	}

	return stateResponse, nil
}
