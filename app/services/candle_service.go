package services

import (
	"database/sql"
	"encoding/json"
	"gotrading/algo"
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
	candleRepository      repositories.CandleRepository
	signalEventRepository repositories.SignalEventRepository
	DataFrame             *models.DataFrameCandle
}

func NewDataFrameCandleService(candleRepository repositories.CandleRepository, signalEventRepository repositories.SignalEventRepository) *DataFrameCandleService {
	return &DataFrameCandleService{
		candleRepository:      candleRepository,
		signalEventRepository: signalEventRepository,
	}
}

func (s *DataFrameCandleService) SetDataFrame(productCode string, duration time.Duration, limit int) error {
	candles, err := s.candleRepository.FindWithLimit(productCode, duration, limit)
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

func (s *DataFrameCandleService) Closes() []float64 {
	return s.DataFrame.Closes()
}

func (s *DataFrameCandleService) AddSMA(period int) bool {
	closes := s.Closes()

	if len(closes) >= period {
		s.DataFrame.SMAs = append(s.DataFrame.SMAs, models.SMA{
			Period: period,
			Values: talib.Sma(closes, period),
		})
	}

	return false
}

func (s *DataFrameCandleService) AddEMA(period int) bool {
	closes := s.Closes()

	if len(closes) >= period {
		s.DataFrame.EMAs = append(s.DataFrame.EMAs, models.EMA{
			Period: period,
			Values: talib.Ema(closes, period),
		})
	}

	return false
}

// moving average type in talib
// const (
//
//	SMA MaType = iota
//	EMA
//	WMA
//	DEMA
//	TEMA
//	TRIMA
//	KAMA
//	MAMA
//	T3MA
//
// )
func (s *DataFrameCandleService) AddBBands(period int, k float64, maType string) bool {
	closes := s.Closes()

	typeMap := map[string]int{
		"sma":   0,
		"ema":   1,
		"wma":   2,
		"dema":  3,
		"tema":  4,
		"trima": 5,
	}

	if len(closes) >= period {
		upper, mid, lower := talib.BBands(closes, period, k, k, talib.MaType(typeMap[maType]))
		s.DataFrame.BBands = &models.BBands{
			Period: period,
			K:      k,
			Upper:  upper,
			Mid:    mid,
			Lower:  lower,
		}

		return true
	}

	return false
}

func (s *DataFrameCandleService) AddRSI(period int) bool {
	closes := s.Closes()

	if len(closes) >= period {
		s.DataFrame.RSI = &models.RSI{
			Period: period,
			Values: talib.Rsi(closes, period),
		}

		return true
	}

	return false
}

func (s *DataFrameCandleService) AddMACD(fastPeriod, slowPeriod, signalPeriod int) bool {
	closes := s.Closes()

	if len(closes) >= slowPeriod {
		values, signalValues, histgram := talib.Macd(closes, fastPeriod, slowPeriod, signalPeriod)

		s.DataFrame.MACD = &models.MACD{
			FastPeriod:   fastPeriod,
			SlowPeriod:   slowPeriod,
			SignalPeriod: signalPeriod,
			Values:       values,
			SignalValues: signalValues,
			Histgram:     histgram,
		}

		return true
	}

	return false
}

func (s *DataFrameCandleService) AddHV(period int) bool {
	closes := s.Closes()

	if len(closes) >= period {
		s.DataFrame.HVs = append(s.DataFrame.HVs, models.HV{
			Period: period,
			Values: algo.HV(closes, period),
		})

		return true
	}

	return false
}

func (s *DataFrameCandleService) AddEvents(t time.Time) bool {
	signalEvents, err := s.signalEventRepository.FindByProductCodeAndAfterTime(s.DataFrame.ProductCode, t)
	if err != nil {
		log.Println(err)
		return false
	}

	if len(signalEvents.Signals) == 0 {
		return false
	}

	s.DataFrame.Events = signalEvents
	return true
}
