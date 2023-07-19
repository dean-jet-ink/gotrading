package models

import (
	"encoding/json"
	"time"
)

type SignalEvent struct {
	Time        time.Time `json:"time,omitempty"`
	ProductCode string    `json:"product_code,omitempty"`
	Side        string    `json:"side,omitempty"`
	Price       float64   `json:"price,omitempty"`
	Size        float64   `json:"size,omitempty"`
}

type SignalEvents struct {
	Signals []*SignalEvent `json:"signals,omitempty"`
}

func (s *SignalEvents) CollectAfterTime(t time.Time) *SignalEvents {
	signalEvents := &SignalEvents{}

	for i, signalEvent := range s.Signals {
		if signalEvent.Time.Before(t) {
			continue
		}

		signalEvents.Signals = s.Signals[i:]
		return signalEvents
	}

	return nil
}

func (s *SignalEvents) Profit() float64 {
	total := 0.0
	beforeSell := 0.0
	isHold := false

	for _, signalEvent := range s.Signals {
		if signalEvent.Side == "BUY" {
			total -= signalEvent.Price
			isHold = true
		} else {
			total += signalEvent.Price
			beforeSell = total
			isHold = false
		}
	}

	if isHold {
		return beforeSell
	}
	return total
}

func (s *SignalEvents) MarshalJSON() ([]byte, error) {
	js, err := json.Marshal(&struct {
		Signals []*SignalEvent `json:"signals,omitempty"`
		Profit  float64        `json:"profit,omitempty"`
	}{
		Signals: s.Signals,
		Profit:  s.Profit(),
	})
	if err != nil {
		return nil, err
	}

	return js, nil
}
