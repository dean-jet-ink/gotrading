package repositories

import (
	"database/sql"
	"fmt"
	"gotrading/app/models"
	"log"
	"strings"
	"time"
)

type SignalEventRepository interface {
	FindByProductCodeWithLimit(productCode string, limit int) (*models.SignalEvents, error)
	FindByProductCodeAndAfterTime(productCode string, t time.Time) (*models.SignalEvents, error)
	Save(*models.SignalEvent) error
}

type SignalEventRepositoryImpl struct {
	db        *sql.DB
	tableName string
}

func NewSignalRepository(db *sql.DB) SignalEventRepository {
	return &SignalEventRepositoryImpl{
		db:        db,
		tableName: "signal_events",
	}
}

func (r *SignalEventRepositoryImpl) getSignalEvents(cmd string, args ...any) (*models.SignalEvents, error) {
	rows, err := r.db.Query(cmd, args...)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	signalEvents := &models.SignalEvents{}
	for rows.Next() {
		signalEvent := &models.SignalEvent{}

		err = rows.Scan(&signalEvent.Time, &signalEvent.ProductCode, &signalEvent.Side, &signalEvent.Price, &signalEvent.Size)
		if err != nil {
			return nil, err
		}

		signalEvents.Signals = append(signalEvents.Signals, signalEvent)
	}

	return signalEvents, nil
}

func (r *SignalEventRepositoryImpl) FindByProductCodeWithLimit(productCode string, limit int) (*models.SignalEvents, error) {
	cmd := fmt.Sprintf(`
		SELECT *
		FROM (
			SELECT time, product_code, side, price, size
			FROM %s
			WHERE product_code = ?
			ORDER BY time DESC
			LIMIT ?
		)
		ORDER BY time ASC
	`, r.tableName)

	signalEvents, err := r.getSignalEvents(cmd, productCode, limit)
	if err != nil {
		return nil, err
	}

	return signalEvents, nil
}

func (r *SignalEventRepositoryImpl) FindByProductCodeAndAfterTime(productCode string, t time.Time) (*models.SignalEvents, error) {
	cmd := fmt.Sprintf(`
		SELECT time, product_code, side, price, size
		FROM %s
		WHERE product_code = ? AND DATETIME(time) >= DATETIME(?)
	`, r.tableName)

	signalEvents, err := r.getSignalEvents(cmd, productCode, t)
	if err != nil {
		return nil, err
	}

	return signalEvents, nil
}

func (r *SignalEventRepositoryImpl) Save(event *models.SignalEvent) error {
	cmd := fmt.Sprintf(`
		INSERT INTO %s (time, product_code, side, price, size)
		VALUES (?, ?, ?, ?, ?)
	`, r.tableName)

	_, err := r.db.Exec(cmd, event.Time.Format(time.RFC3339), event.ProductCode, event.Side, event.Price, event.Size)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			log.Println(err)
			return nil
		}

		log.Println(err)
		return err
	}

	return nil
}
