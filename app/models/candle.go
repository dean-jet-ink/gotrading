package models

import (
	"time"
)

type Candle struct {
	Time   time.Time `json:"time"`
	Open   float64   `json:"open"`
	Close  float64   `json:"close"`
	High   float64   `json:"high"`
	Low    float64   `json:"low"`
	Volume float64   `json:"volume"`
}

type DataFrameCandle struct {
	ProductCode string        `json:"product_code"`
	Duration    time.Duration `json:"duration"`
	Candles     []*Candle     `json:"candles"`
	SMAs        []SMA         `json:"smas,omitempty"`
}

type SMA struct {
	Period int       `json:"period,omitempty"`
	Values []float64 `json:"values,omitempty"`
}

func (df *DataFrameCandle) Opens() []float64 {
	opens := make([]float64, len(df.Candles))
	for i, candle := range df.Candles {
		opens[i] = candle.Open
	}

	return opens
}

func (df *DataFrameCandle) Closes() []float64 {
	closes := make([]float64, len(df.Candles))
	for i, candle := range df.Candles {
		closes[i] = candle.Close
	}

	return closes
}

func (df *DataFrameCandle) Highs() []float64 {
	highs := make([]float64, len(df.Candles))
	for i, candle := range df.Candles {
		highs[i] = candle.High
	}

	return highs
}

func (df *DataFrameCandle) Lows() []float64 {
	lows := make([]float64, len(df.Candles))
	for i, candle := range df.Candles {
		lows[i] = candle.Low
	}

	return lows
}

func (df *DataFrameCandle) Volumes() []float64 {
	volumes := make([]float64, len(df.Candles))
	for i, candle := range df.Candles {
		volumes[i] = candle.Volume
	}

	return volumes
}
