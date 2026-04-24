package jitorpc

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
)

type TipFloor struct {
	Time                        time.Time `json:"time"`
	LandedTips25thPercentile    float64   `json:"landed_tips_25th_percentile"`
	LandedTips50thPercentile    float64   `json:"landed_tips_50th_percentile"`
	LandedTips75thPercentile    float64   `json:"landed_tips_75th_percentile"`
	LandedTips95thPercentile    float64   `json:"landed_tips_95th_percentile"`
	LandedTips99thPercentile    float64   `json:"landed_tips_99th_percentile"`
	EmaLandedTips50thPercentile float64   `json:"ema_landed_tips_50th_percentile"`
}

// GetTipFloor
//
// curl 'curl https://bundles.jito.wtf/api/v1/bundles/tip_floor' \
// --header 'Content-Type: application/json' \
// --data '
//
//	{
//	    "jsonrpc": "2.0",
//	    "id": 1,
//	    "method": "getTipFloor",
//	    "params": []
//	}
//
// '
// Response
//
//	{
//	  "jsonrpc": "2.0",
//	  "id": 1,
//	  "result": [
//	    {
//	      "time": "2025-04-24T13:47:32+00:00",
//	      "landed_tips_25th_percentile": 0.000005,
//	      "landed_tips_50th_percentile": 0.00001,
//	      "landed_tips_75th_percentile": 0.001,
//	      "landed_tips_95th_percentile": 0.00533476325,
//	      "landed_tips_99th_percentile": 0.00826274555000001,
//	      "ema_landed_tips_50th_percentile": 0.000009952174575705424
//	    }
//	  ]
//	}
func GetTipFloors(ctx context.Context) (tipFloors []*TipFloor, err error) {

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://bundles.jito.wtf/api/v1/bundles/tip_floor", nil)
	if err != nil {
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("status code %d", resp.StatusCode)
		return
	}

	jsonData, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	err = jsoniter.Unmarshal(jsonData, &tipFloors)

	return
}

func GetTipStream(ctx context.Context, chTipFloor chan *TipFloor) {
	ctx1, cancel1 := context.WithCancel(ctx)
	defer cancel1()

	uri := "wss://bundles.jito.wtf/api/v1/bundles/tip_stream"

	dialer := &websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 45 * time.Second,
		// EnableCompression: true,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	chClient := make(chan *websocket.Conn, 1)
	defer close(chClient)

	var (
		conn *websocket.Conn

		err error
		wg  sync.WaitGroup
	)

	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()

loop:
	for {
		wg.Go(func() {
			if conn != nil {
				conn.Close()
			}

			if conn, _, err = dialer.DialContext(ctx1, uri, nil); err != nil {
				slog.Error("dialer.DialContext(ctx1, uri, nil) fail", "uri", uri, "error", err)
				// sleep and try again
				time.Sleep(time.Millisecond * 300)

				chClient <- nil
				return
			}

			chClient <- conn
		})

		select {
		case <-ctx1.Done():
			break loop
		case conn := <-chClient:
			if conn == nil {
				break
			}

			ch := make(chan *TipFloor)

			wg.Go(func() {

				defer close(ch)

				for {
					messageType, message, err := conn.ReadMessage()
					if err != nil {
						slog.Error("conn.ReadMessage() fail", "uri", uri, "error", err)
						break
					}

					if messageType != 1 {
						slog.Error("messageType != 1 fail", "uri", uri, "error", err)
						break
					}

					var tipFloors []*TipFloor
					if err := jsoniter.Unmarshal(message, &tipFloors); err != nil {
						slog.Error("jsoniter.Unmarshal(message, &tipFloors) fail", "uri", uri, "message", string(message), "error", err)
						break
					}

					for _, tipFloor := range tipFloors {
						ch <- tipFloor
					}
				}
			})

		loop2:
			for {
				select {
				case <-ctx1.Done():
					break loop
				case d, ok := <-ch:
					if !ok {
						break loop2
					}
					chTipFloor <- d
				}
			}
		}
	}

	wg.Wait()

}
