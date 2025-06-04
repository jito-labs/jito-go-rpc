package jitorpc

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type TipFloor struct {
	Time                        string  `json:"time"`
	LandedTips25thPercentile    float64 `json:"landed_tips_25th_percentile"`
	LandedTips50thPercentile    float64 `json:"landed_tips_50th_percentile"`
	LandedTips75thPercentile    float64 `json:"landed_tips_75th_percentile"`
	LandedTips95thPercentile    float64 `json:"landed_tips_95th_percentile"`
	LandedTips99thPercentile    float64 `json:"landed_tips_99th_percentile"`
	EmaLandedTips50thPercentile float64 `json:"ema_landed_tips_50th_percentile"`
}

func (c *JitoJsonRpcClient) GetTipFloors() ([]TipFloor, error) {
	url := fmt.Sprintf("%s/bundles/tip_floor", c.BaseURL)
	if c.UUID != "" {
		url = fmt.Sprintf("%s?uuid=%s", url, c.UUID)
	}

	if c.isDebugEnabled() {
		fmt.Printf("Sending request to: %s\n", url)
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
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

	var tipFloors []TipFloor
	err = json.NewDecoder(resp.Body).Decode(&tipFloors)
	if err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return tipFloors, nil
}
