package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"gotrading/config"

	_ "github.com/mattn/go-sqlite3"
)

var dbConn *sql.DB

const tableName = "signal_events"

func init() {
	cfg := config.Config()
	var err error
	dbConn, err = sql.Open(cfg.DBDriver(), cfg.DBName())
	if err != nil {
		log.Fatalf("Database connection failed: %s", err.Error())
	}

	// signal_events
	cmd := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		time DATETIME PRIMARY KEY NOT NULL,
		product_code STRING,
		side STRING,
		price FLOAT,
		size FLOAT
	)`, tableName)

	if _, err := dbConn.Exec(cmd); err != nil {
		log.Fatalf("init in database: %s", err.Error())
	}

	// BTC_XXX_1X
	for _, duration := range cfg.Durations() {
		tableName := CandleTableName(cfg.ProductCode(), duration)
		cmd = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
			time DATETIME PRIMARY KEY NOT NULL,
			open FLOAT,
			close FLOAT,
			high FLOAT,
			low FLOAT,
			volume FLOAT
		)`, tableName)

		if _, err = dbConn.Exec(cmd); err != nil {
			log.Fatalf("init in database: %s", err.Error())
		}
	}
}

func CandleTableName(productCode string, duration time.Duration) string {
	return fmt.Sprintf("%s_%s", productCode, duration)
}

func DBConn() *sql.DB {
	return dbConn
}
