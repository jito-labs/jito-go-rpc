package jitorpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
)

type TipAccount struct {
	Address string `json:"address"`
}

type BundleStatusResponse struct {
	Context struct {
		Slot int64 `json:"slot"`
	} `json:"context"`
	Value []struct {
		BundleID           string   `json:"bundle_id"`
		Transactions       []string `json:"transactions"`
		Slot               int64    `json:"slot"`
		ConfirmationStatus string   `json:"confirmation_status"`
		Err                struct {
			Ok interface{} `json:"Ok"`
		} `json:"err"`
	} `json:"value"`
}

func (c *JitoJsonRpcClient) GetTipAccounts() (json.RawMessage, error) {
	endpoint := "/bundles"
	if c.UUID != "" {
		endpoint = fmt.Sprintf("%s?uuid=%s", endpoint, c.UUID)
	}
	return c.sendRequest(endpoint, "getTipAccounts", nil)
}

func (c *JitoJsonRpcClient) GetRandomTipAccount() (*TipAccount, error) {
	rawResponse, err := c.GetTipAccounts()
	if err != nil {
		return nil, err
	}

	var tipAddresses []string
	err = json.Unmarshal(rawResponse, &tipAddresses)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal tip accounts: %w", err)
	}

	if len(tipAddresses) == 0 {
		return nil, fmt.Errorf("no tip accounts available")
	}

	randomIndex := rand.Intn(len(tipAddresses))
	return &TipAccount{Address: tipAddresses[randomIndex]}, nil
}

func (c *JitoJsonRpcClient) GetBundleStatuses(bundleIds []string) (*BundleStatusResponse, error) {
	endpoint := "/getBundleStatuses"
	if c.UUID != "" {
		endpoint = fmt.Sprintf("%s?uuid=%s", endpoint, c.UUID)
	}
	params := [][]string{bundleIds}
	responseBody, err := c.sendRequest(endpoint, "getBundleStatuses", params)
	if err != nil {
		return nil, err
	}

	var response BundleStatusResponse
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal bundle statuses: %w", err)
	}

	return &response, nil
}

func (c *JitoJsonRpcClient) SendBundle(bundleTransactions [][]string) (json.RawMessage, error) {
    url := fmt.Sprintf("%s/bundles", c.BaseURL)
    if c.UUID != "" {
        url = fmt.Sprintf("%s?uuid=%s", url, c.UUID)
    }

    var transactions []string
    for _, txGroup := range bundleTransactions {
        transactions = append(transactions, txGroup...)
    }

    request := struct {
        JsonRpc string        `json:"jsonrpc"`
        ID      int           `json:"id"`
        Method  string        `json:"method"`
        Params  []interface{} `json:"params"`
    }{
        JsonRpc: "2.0",
        ID:      1,
        Method:  "sendBundle",
        Params: []interface{}{
            transactions,
            map[string]string{"encoding": "base64"},
        },
    }

	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	if c.isDebugEnabled() {
		fmt.Printf("Sending request to: %s\n", url)
		fmt.Printf("Request body: %s\n", string(requestBody))
	}

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

	if c.isDebugEnabled() {
		fmt.Printf("Response status: %s\n", resp.Status)
	}

	var jsonResp JsonRpcResponse
	err = json.NewDecoder(resp.Body).Decode(&jsonResp)
	if err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	if jsonResp.Error != nil {
		return nil, fmt.Errorf("RPC error: %s", jsonResp.Error.Message)
	}

	if c.isDebugEnabled() {
		fmt.Printf("Response body: %s\n", string(jsonResp.Result))
	}

	return jsonResp.Result, nil
}

func (c *JitoJsonRpcClient) GetInflightBundleStatuses(params interface{}) (json.RawMessage, error) {
	endpoint := "/getInflightBundleStatuses"
	if c.UUID != "" {
		endpoint = fmt.Sprintf("%s?uuid=%s", endpoint, c.UUID)
	}
	return c.sendRequest(endpoint, "getInflightBundleStatuses", params)
}