package repositories

import (
	"database/sql"
	"fmt"
	"gotrading/app/models"
	"log"
	"time"
)

type CandleRepository interface {
	FindByTime(productCode string, duration time.Duration, t time.Time) (*models.Candle, error)
	FindWithLimit(productCode string, duration time.Duration, limit int) ([]*models.Candle, error)
	Create(productCode string, duration time.Duration, candle *models.Candle) error
	Update(productCode string, duration time.Duration, candle *models.Candle) error
}

type BitflyerCandleRepository struct {
	db *sql.DB
}

func NewBitflyerCandleRepository(db *sql.DB) CandleRepository {
	return &BitflyerCandleRepository{
		db: db,
	}
}

func (r *BitflyerCandleRepository) FindByTime(productCode string, duration time.Duration, t time.Time) (*models.Candle, error) {
	cmd := fmt.Sprintf(`SELECT time, open, close, high, low, volume FROM %s WHERE time = ?`, r.TableName(productCode, duration))
	row := r.db.QueryRow(cmd, t.Format(time.RFC3339))

	var candle = &models.Candle{}
	err := row.Scan(&candle.Time, &candle.Open, &candle.Close, &candle.High, &candle.Low, &candle.Volume)
	if err != nil {
		return nil, err
	}

	return candle, nil
}

func (r *BitflyerCandleRepository) FindWithLimit(productCode string, duration time.Duration, limit int) ([]*models.Candle, error) {
	cmd := fmt.Sprintf(`
		SELECT * 
		FROM (
			SELECT time, open, close, high, low, volume
			FROM %s
			ORDER BY time DESC
			LIMIT ?
		)
		ORDER BY time ASC
	`, r.TableName(productCode, duration))

	rows, err := r.db.Query(cmd, limit)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	var candles []*models.Candle
	for rows.Next() {
		var candle = &models.Candle{}
		err = rows.Scan(&candle.Time, &candle.Open, &candle.Close, &candle.High, &candle.Low, &candle.Volume)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		candles = append(candles, candle)
	}

	return candles, nil
}

func (r *BitflyerCandleRepository) Create(productCode string, duration time.Duration, candle *models.Candle) error {
	cmd := fmt.Sprintf(`INSERT INTO %s (time, open, close, high, low, volume)
	VALUES (?, ?, ?, ?, ?, ?)`, r.TableName(productCode, duration))

	if _, err := r.db.Exec(cmd, candle.Time.Format(time.RFC3339), candle.Open, candle.Close, candle.High, candle.Low, candle.Volume); err != nil {
		log.Println("BitflyerCandleRepository Create: ")
		return err
	}

	return nil
}

func (r *BitflyerCandleRepository) Update(productCode string, duration time.Duration, candle *models.Candle) error {
	cmd := fmt.Sprintf(`UPDATE %s SET open = ?, close = ?, high = ?, low = ?, volume = ? WHERE time = ?`, r.TableName(productCode, duration))

	if _, err := r.db.Exec(cmd, candle.Open, candle.Close, candle.High, candle.Low, candle.Volume, candle.Time.Format(time.RFC3339)); err != nil {
		log.Println("BitflyerCandleRepository Update: ")
		return err
	}

	return nil
}

func (r *BitflyerCandleRepository) TableName(productCode string, duration time.Duration) string {
	return fmt.Sprintf("%s_%s", productCode, duration)
}
