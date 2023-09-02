package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type RequestBody struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  []any  `json:"params"`

	// Id of the json-rpc call. Doesn't really matter since we're not batching multiple json-rpc calls.
	Id int `json:"id"`
}

type ResponseBody struct {
	Id      int             `json:"id"`
	Jsonrpc string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result"`
}

// EthClient is a client for interfacing with the Ethereum blockchain through the JSON-RPC API.
// See https://ethereum.org/en/developers/docs/apis/json-rpc for the respective methods in EthClient.
type EthClient struct {
	url string
	id  int
}

func NewEthClient(url string) EthClient {
	return EthClient{
		id:  1,
		url: url,
	}
}

func (ec *EthClient) BlockNumber() (int, error) {
	result, err := ec.post("eth_blockNumber", nil)
	if err != nil {
		return 0, fmt.Errorf("failed to get block number: %w", err)
	}

	var blockNumberInString string
	if err := json.Unmarshal(result, &blockNumberInString); err != nil {
		return 0, fmt.Errorf("failed to unmarshal block number result: %w", err)
	}

	blockNumber, err := strconv.ParseInt(blockNumberInString[2:], 16, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse block number: %w", err)
	}

	return int(blockNumber), nil
}

func (ec *EthClient) GetBlockByNumber(blockNumber int) (Block, error) {
	var block Block
	params := []any{
		"0x" + strconv.FormatInt(int64(blockNumber), 16),
		true,
	}
	result, err := ec.post("eth_getBlockByNumber", params)
	if err != nil {
		return block, fmt.Errorf("failed to get block by number: %w", err)
	}

	// Block isn't mined yet.
	if len(result) == 0 {
		return block, nil
	}

	if err := json.Unmarshal(result, &block); err != nil {
		return block, fmt.Errorf("failed to unmarshal get block by number result: %w", err)
	}

	return block, nil
}

// post to the blockchain through json-rpc.
func (ec *EthClient) post(method string, params []any) (json.RawMessage, error) {
	reqBody := &RequestBody{
		Jsonrpc: "2.0",
		Method:  method,
		Params:  params,
		Id:      ec.id,
	}
	buf, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Increase id so we don't reuse it.
	ec.id += 1

	resp, err := http.Post(ec.url, "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return nil, fmt.Errorf("failed to post: %w", err)
	}
	defer resp.Body.Close()

	var responseBody ResponseBody
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	return responseBody.Result, nil
}
