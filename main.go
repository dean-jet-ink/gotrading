package main

import (
	"gotrading/app/controllers"
	"gotrading/app/repositories"
	"gotrading/app/services"
	"gotrading/config"
	"gotrading/database"
	"gotrading/utils"
)

func main() {
	cfg := config.Config()
	utils.SetLogging(cfg.LogFile())
	dbConn := database.DBConn()
	repository := repositories.NewBitflyerCandleRepository(dbConn)
	streamService := services.NewStreamCandleService(repository)
	go streamService.StreamIngestionData(cfg.ApiKey(), cfg.ApiSecret(), cfg.ProductCode(), cfg.Durations(), cfg.TradeDuration())
	ws := controllers.NewWebServer("", 8080, dbConn)
	ws.Start()
}
