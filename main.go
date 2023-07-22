package main

import (
	"gotrading/app/controllers"
	"gotrading/config"
	"gotrading/database"
	"gotrading/utils"
)

func main() {
	cfg := config.Config()
	utils.SetLogging(cfg.LogFile())

	db := database.DBConn()

	webserver := controllers.NewWebServer("", 8080, db)
	webserver.Start()

	// productCode := cfg.ProductCode()
	// tradeDuration := cfg.TradeDuration()
	// candleRepo := repositories.NewBitflyerCandleRepository(db)
	// signalRepo := repositories.NewSignalRepository(db)
	// candleService.SetDataFrame(productCode, tradeDuration, 10)
	// c1 := candleService.DataFrame.Candles[0]
	// c2 := candleService.DataFrame.Candles[5]

	// signalService := services.NewSignalEventService(signalRepo, productCode)
	// signalService.Buy(productCode, c1.Time.UTC(), c1.Close, 1.0, false)
	// signalService.Sell(productCode, c2.Time.UTC(), c2.Close, 1.0, false)
}
