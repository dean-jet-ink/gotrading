package services

import (
	"encoding/json"
	"gotrading/app/models"
	"gotrading/app/repositories"
	"log"
	"time"
)

type SignalEventService struct {
	signalEventRepository repositories.SignalEventRepository
	productCode           string
	signalEvents          *models.SignalEvents
}

func NewSignalEventService(repository repositories.SignalEventRepository, productCode string) *SignalEventService {
	signalEventService := &SignalEventService{
		signalEventRepository: repository,
		productCode:           productCode,
		signalEvents:          &models.SignalEvents{},
	}

	return signalEventService
}

func (s *SignalEventService) Signals() []*models.SignalEvent {
	return s.signalEvents.Signals
}

func (s *SignalEventService) SignalEvents() *models.SignalEvents {
	return s.signalEvents
}

func (s *SignalEventService) saveSignalEvent(signalEvent *models.SignalEvent, isBackTest bool) bool {
	if !isBackTest {
		if err := s.signalEventRepository.Save(signalEvent); err != nil {
			log.Println(err)
			return false
		}
	}

	s.signalEvents.Signals = append(s.Signals(), signalEvent)
	return true
}

func (s *SignalEventService) Buy(productCode string, t time.Time, price, size float64, isBacktest bool) bool {
	signals := s.Signals()
	length := len(signals)

	if length != 0 && (signals[length-1].Side == "BUY" || signals[length-1].Time.After(t)) {
		return false
	}

	signalEvent := &models.SignalEvent{
		Time:        t,
		ProductCode: productCode,
		Side:        "BUY",
		Price:       price,
		Size:        size,
	}

	return s.saveSignalEvent(signalEvent, isBacktest)
}

func (s *SignalEventService) Sell(productCode string, t time.Time, price, size float64, isBacktest bool) bool {
	length := len(s.Signals())
	if length == 0 {
		return false
	}

	lastSignal := s.Signals()[length-1]
	if lastSignal.Side == "SELL" || lastSignal.Time.After(t) {
		return false
	}

	signalEvent := &models.SignalEvent{
		Time:        t,
		ProductCode: productCode,
		Side:        "SELL",
		Price:       price,
		Size:        size,
	}

	return s.saveSignalEvent(signalEvent, isBacktest)
}

func (s *SignalEventService) MarshalSignalEvents() ([]byte, error) {
	js, err := json.Marshal(s.signalEvents)
	if err != nil {
		return nil, err
	}

	return js, nil
}

func (s *SignalEventService) GetEventWithLimit(limit int) (*models.SignalEvents, error) {
	return s.signalEventRepository.FindByProductCodeWithLimit(s.productCode, limit)
}

func (s *SignalEventService) GetEventAfterTime(t time.Time) (*models.SignalEvents, error) {
	return s.signalEventRepository.FindByProductCodeAndAfterTime(s.productCode, t)
}

func (s *SignalEventService) CollectEventAfterTime(t time.Time) *models.SignalEvents {
	return s.signalEvents.CollectAfterTime(t)
}
