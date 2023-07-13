package bitflyer

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"gotrading/api"
	"gotrading/socket/jsonrpc"
)

type BitflyerClient struct {
	*api.APIClient
	realTimeURL *url.URL
}

func NewBitflyerClient(key, secret string) (*BitflyerClient, error) {
	baseURL, err := url.Parse("https://api.bitflyer.com/v1/")
	if err != nil {
		return nil, err
	}

	if len(key) == 0 {
		return nil, errors.New("missing key")
	}
	if len(secret) == 0 {
		return nil, errors.New("missing secret")
	}

	realTimeURL := &url.URL{Scheme: "wss", Host: "ws.lightstream.bitflyer.com", Path: "/json-rpc"}

	return &BitflyerClient{
		APIClient: &api.APIClient{
			URL: baseURL,
			HTTPClient: &http.Client{
				Timeout: 60 * time.Second,
			},
			Key:    key,
			Secret: secret,
		},
		realTimeURL: realTimeURL,
	}, nil
}

func (c *BitflyerClient) NewRequest(ctx context.Context, urlPath, method string, queryParams map[string]string, body []byte, isPrivate bool) (*http.Request, error) {
	var header map[string]string
	if isPrivate {
		header = c.PrivateAPIHeader(method, urlPath, body)
	} else {
		header = c.PublicAPIHeader()
	}

	return c.APIClient.NewRequest(ctx, urlPath, method, queryParams, header, body)
}

func (c *BitflyerClient) PrivateAPIHeader(method, url string, body []byte) map[string]string {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	text := timestamp + method + url + string(body)

	h := hmac.New(sha256.New, []byte(c.Secret))
	h.Write([]byte(text))
	sign := hex.EncodeToString(h.Sum(nil))

	return map[string]string{
		"ACCESS-KEY":       c.Key,
		"ACCESS-TIMESTAMP": timestamp,
		"ACCESS-SIGN":      sign,
		"Content-Type":     "application/json",
	}
}

func (c *BitflyerClient) PublicAPIHeader() map[string]string {
	return map[string]string{
		"Content-Type": "application/json",
	}
}

func (c *BitflyerClient) GetBalance(ctx context.Context) ([]Balance, error) {
	req, err := c.NewRequest(ctx, "me/getbalance", "GET", nil, nil, true)
	if err != nil {
		return nil, fmt.Errorf("GetBalance: %s", err.Error())
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GetBalance: %s", err.Error())
	}

	var balance []Balance
	err = c.DecodeBody(resp, balance)
	if err != nil {
		return nil, fmt.Errorf("GetBalance: %s", err.Error())
	}

	return balance, nil
}

func (c *BitflyerClient) GetTicker(ctx context.Context, productCode string) (*Ticker, error) {
	queries := map[string]string{"product_code": productCode}
	req, err := c.NewRequest(ctx, "ticker", "GET", queries, nil, false)
	if err != nil {
		return nil, fmt.Errorf("GetTicker: %s", err.Error())
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GetTicker: %s", err.Error())
	}

	var ticker = &Ticker{}
	err = c.DecodeBody(resp, ticker)
	if err != nil {
		return nil, fmt.Errorf("GetTicker: %s", err.Error())
	}

	return ticker, nil
}

func (c *BitflyerClient) GetRealTimeTicker(productCode string, ch chan<- *Ticker) {
	rpc, err := jsonrpc.NewJsonRPC2Client(c.realTimeURL.String())
	if err != nil {
		log.Fatal(fmt.Errorf("GetRealTimeTicker: %s", err.Error()))
	}
	defer rpc.Close()

	type Subscribe struct {
		Channel string `json:"channel"`
	}

	channel := fmt.Sprintf("lightning_ticker_%s", productCode)
	message := &jsonrpc.JsonRPC2{
		Version: "2.0",
		Method:  "subscribe",
		Params:  Subscribe{Channel: channel},
	}
	rpc.Send(message)

OUTER:
	for {
		rpc.Recieve(message)
		// message format(JsonRPC2)
		// &{2.0 channelMessage map[channel:lightning_ticker_BTC_JPY message:map[best_ask:4.419753e+06 best_ask_size:0.02 best_bid:4.418377e+06 best_bid_size:0.04 circuit_break_end:<nil> ltp:4.419535e+06 market_ask_size:0 market_bid_size:0 preopen_end:<nil> product_code:BTC_JPY state:RUNNING tick_id:1.6414218e+07 timestamp:2023-07-06T06:48:12.7049048Z total_ask_depth:378.14667391 total_bid_depth:572.0329027 volume:1435.03618731 volume_by_product:1435.03618731]] <nil> <nil>}

		switch v := message.Params.(type) {

		case map[string]interface{}:
			for key, value := range v {
				if key == "message" {
					marshalTic, err := json.Marshal(value)
					if err != nil {
						log.Println("GetRealTimeTicker:", err.Error())
						continue OUTER
					}

					ticker := &Ticker{}
					if err = json.Unmarshal(marshalTic, ticker); err != nil {
						log.Println("GetRealTimeTicker:", err.Error())
						continue OUTER
					}

					ch <- ticker
				}
			}
		}
	}
}

func (c *BitflyerClient) SendOrder(ctx context.Context, order *Order) (*ResponseSendChildOrder, error) {
	body, err := json.Marshal(order)
	if err != nil {
		return nil, err
	}
	req, err := c.NewRequest(ctx, "me/sendchildorder", "POST", nil, body, true)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	responseSendChildOrder := &ResponseSendChildOrder{}
	if err = c.DecodeBody(resp, responseSendChildOrder); err != nil {
		return nil, err
	}

	return responseSendChildOrder, nil
}

func (c *BitflyerClient) OrderList(ctx context.Context, queryParams map[string]string) ([]Order, error) {
	req, err := c.NewRequest(ctx, "me/getchildorders", "GET", queryParams, nil, true)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	var orderList []Order
	if err = c.DecodeBody(resp, orderList); err != nil {
		return nil, err
	}

	return orderList, nil
}
