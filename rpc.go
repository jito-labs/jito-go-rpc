package jitorpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *JitoJsonRpcClient) sendRequest(endpoint, method string, params interface{}) (json.RawMessage, error) {
	url := fmt.Sprintf("%s%s", c.BaseURL, endpoint)

	request := JsonRpcRequest{
		JsonRpc: "2.0",
		ID:      1,
		Method:  method,
		Params:  params,
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	fmt.Printf("Sending request to: %s\n", url)
	fmt.Printf("Request body: %s\n", string(requestBody))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.UUID != "" {
		req.Header.Set("x-jito-auth", c.UUID)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	fmt.Printf("Response status: %s\n", resp.Status)

	var jsonResp JsonRpcResponse
	err = json.NewDecoder(resp.Body).Decode(&jsonResp)
	if err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	if jsonResp.Error != nil {
		return nil, fmt.Errorf("RPC error: %s", jsonResp.Error.Message)
	}

	fmt.Printf("Response body: %s\n", string(jsonResp.Result))

	return jsonResp.Result, nil
}
