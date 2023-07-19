package main

import (
	"fmt"
	"gotrading/app/repositories"
	"gotrading/app/services"
	"gotrading/config"
	"gotrading/database"
	"gotrading/utils"
)

func main() {
	cfg := config.Config()
	productCode := cfg.ProductCode()
	tradeDuration := cfg.TradeDuration()
	utils.SetLogging(cfg.LogFile())

	db := database.DBConn()

	candleRepo := repositories.NewBitflyerCandleRepository(db)
	candleService := services.NewDataFrameCandleService(candleRepo)
	candleService.SetDataFrame(productCode, tradeDuration, 10)
	c1 := candleService.DataFrame.Candles[0]
	c2 := candleService.DataFrame.Candles[5]

	signalRepo := repositories.NewSignalRepository(db)
	signalService := services.NewSignalEventService(signalRepo, productCode)
	signalService.Buy(productCode, c1.Time.UTC(), c1.Close, 1.0, false)
	signalService.Sell(productCode, c2.Time.UTC(), c2.Close, 1.0, false)
	signalEvents1, _ := signalService.GetEventAfterTime(c1.Time)
	signalEvents2, _ := signalService.GetEventWithLimit(1)
	signalEvents3 := signalService.CollectEventAfterTime(c1.Time)

	fmt.Println(signalEvents1.Signals)
	fmt.Println(signalEvents2.Signals)
	fmt.Println(signalEvents3.Signals)
}
