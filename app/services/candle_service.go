package services

import (
	"database/sql"
	"encoding/json"
	"gotrading/api/bitflyer"
	"gotrading/app/models"
	"gotrading/app/repositories"
	"log"
	"time"

	"github.com/markcheno/go-talib"
)

type StreamCandleService struct {
	candleRepository repositories.CandleRepository
}

func NewStreamCandleService(candleRepository repositories.CandleRepository) *StreamCandleService {
	return &StreamCandleService{
		candleRepository: candleRepository,
	}
}

func (s *StreamCandleService) createCandleWithDuration(ticker *bitflyer.Ticker, productCode string, duration time.Duration) (bool, error) {
	price := ticker.MidPrice()
	currentCandle, err := s.candleRepository.FindByTime(productCode, duration, ticker.TruncateDateTime(duration))
	if err != nil {
		if err == sql.ErrNoRows {
			currentCandle = &models.Candle{
				Time:   ticker.TruncateDateTime(duration),
				Open:   price,
				Close:  price,
				High:   price,
				Low:    price,
				Volume: ticker.Volume,
			}

			if err = s.candleRepository.Create(productCode, duration, currentCandle); err != nil {
				return false, err
			}

			return true, nil
		} else {
			return false, err
		}
	}

	if currentCandle.High < price {
		currentCandle.High = price
	} else if currentCandle.Low < price {
		currentCandle.Low = price
	}
	currentCandle.Close = price
	currentCandle.Volume += ticker.Volume
	if err = s.candleRepository.Update(productCode, duration, currentCandle); err != nil {
		return false, err
	}

	return false, nil
}

func (s *StreamCandleService) StreamIngestionData(apiKey, secret, productCode string, durations map[string]time.Duration, tradeDuration time.Duration) {
	ch := make(chan *bitflyer.Ticker)
	client, err := bitflyer.NewBitflyerClient(apiKey, secret)
	if err != nil {
		log.Fatal(err)
	}

	go client.GetRealTimeTicker(productCode, ch)

	for ticker := range ch {
		// log.Printf("StreamIngestionData: %v", ticker)

		for _, duration := range durations {
			isCreated, err := s.createCandleWithDuration(ticker, productCode, duration)
			if err != nil {
				log.Fatal(err)
			}

			if isCreated && duration == tradeDuration {
			}
		}
	}
}

type DataFrameCandleService struct {
	candleRepository repositories.CandleRepository
	DataFrame        *models.DataFrameCandle
}

func NewDataFrameCandleService(candleRepository repositories.CandleRepository) *DataFrameCandleService {
	return &DataFrameCandleService{
		candleRepository: candleRepository,
	}
}

func (s *DataFrameCandleService) SetDataFrame(productCode string, duration time.Duration, limit int) error {
	candles, err := s.candleRepository.FindAllWithLimit(productCode, duration, limit)
	if err != nil {
		return err
	}

	s.DataFrame = &models.DataFrameCandle{
		ProductCode: productCode,
		Duration:    duration,
		Candles:     candles,
	}

	return nil
}

func (s *DataFrameCandleService) MarshalDataFrame() ([]byte, error) {
	js, err := json.Marshal(s.DataFrame)
	if err != nil {
		return nil, err
	}

	return js, nil
}

func (s *DataFrameCandleService) AddSMA(period int) {
	closes := s.DataFrame.Closes()

	if len(closes) >= period {
		s.DataFrame.SMAs = append(s.DataFrame.SMAs, models.SMA{
			Period: period,
			Values: talib.Sma(closes, period),
		})
	}
}
