package bitflyer

import (
	"log"
	"time"
)

type Balance struct {
	CurrencyCode string  `json:"currency_code"`
	Amount       float64 `json:"amount"`
	Available    float64 `json:"available"`
}

type Ticker struct {
	ProductCode     string  `json:"product_code"`
	State           string  `json:"state"`
	Timestamp       string  `json:"timestamp"`
	TickID          int     `json:"tick_id"`
	BestBid         float64 `json:"best_bid"`
	BestAsk         float64 `json:"best_ask"`
	BestBidSize     float64 `json:"best_bid_size"`
	BestAskSize     float64 `json:"best_ask_size"`
	TotalBidDepth   float64 `json:"total_bid_depth"`
	TotalAskDepth   float64 `json:"total_ask_depth"`
	MarketBidSize   float64 `json:"market_bid_size"`
	MarketAskSize   float64 `json:"market_ask_size"`
	Ltp             float64 `json:"ltp"`
	Volume          float64 `json:"volume"`
	VolumeByProduct float64 `json:"volume_by_product"`
}

func (t *Ticker) MidPrice() float64 {
	return (t.BestBid + t.BestAsk) / 2
}

func (t *Ticker) DateTime() time.Time {
	dateTime, err := time.Parse(time.RFC3339, t.Timestamp)
	if err != nil {
		log.Printf("DateTime: %s", err.Error())
	}
	return dateTime
}

func (t *Ticker) TruncateDateTime(duration time.Duration) time.Time {
	return t.DateTime().Truncate(duration)
}

type Order struct {
	ID                     int     `json:"id"`
	ChildOrderAcceptanceID string  `json:"child_order_acceptance_id"`
	ProductCode            string  `json:"product_code"`
	ChildOrderType         string  `json:"child_order_type"`
	Side                   string  `json:"side"`
	Price                  float64 `json:"price"`
	Size                   float64 `json:"size"`
	MinuteToExpires        int     `json:"minute_to_expire"`
	TimeInForce            string  `json:"time_in_force"`
	Status                 string  `json:"status"`
	ErrorMessage           string  `json:"error_message"`
	AveragePrice           float64 `json:"average_price"`
	ChildOrderState        string  `json:"child_order_state"`
	ExpireDate             string  `json:"expire_date"`
	ChildOrderDate         string  `json:"child_order_date"`
	OutstandingSize        float64 `json:"outstanding_size"`
	CancelSize             float64 `json:"cancel_size"`
	ExecutedSize           float64 `json:"executed_size"`
	TotalCommission        float64 `json:"total_commission"`
	Count                  int     `json:"count"`
	Before                 int     `json:"before"`
	After                  int     `json:"after"`
}

type ResponseSendChildOrder struct {
	ChildOrderAcceptanceID string `json:"child_order_acceptance_id"`
}
