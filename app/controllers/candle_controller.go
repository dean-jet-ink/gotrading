package controllers

import (
	"database/sql"
	"fmt"
	"gotrading/app/models"
	"gotrading/app/repositories"
	"gotrading/app/services"
	"gotrading/config"
	"net/http"
	"regexp"
	"strconv"
)

type CandleController interface {
	APIUrl() string
	APIHandler(w http.ResponseWriter, r *http.Request)
	validAPIHandler(w http.ResponseWriter, url, productCode string) bool
}

type bitflyerCandleController struct {
	apiUrl                 string
	validPath              *regexp.Regexp
	dataFrameCandleService *services.DataFrameCandleService
}

func (c *bitflyerCandleController) APIUrl() string {
	return c.apiUrl
}

func (c *bitflyerCandleController) DataFrame() *models.DataFrameCandle {
	return c.dataFrameCandleService.DataFrame
}

func NewBitflyerCandleController(db *sql.DB) CandleController {
	candleRepository := repositories.NewBitflyerCandleRepository(db)
	signalEventRepository := repositories.NewSignalRepository(db)
	dataFrameCandleService := services.NewDataFrameCandleService(candleRepository, signalEventRepository)

	return &bitflyerCandleController{
		apiUrl:                 "/api/candle/",
		validPath:              regexp.MustCompile("^/api/candle/$"),
		dataFrameCandleService: dataFrameCandleService,
	}
}

func (c *bitflyerCandleController) validAPIHandler(w http.ResponseWriter, url, productCode string) bool {
	match := c.validPath.FindStringSubmatch(url)
	if len(match) == 0 {
		models.APIError(w, "Not found", http.StatusNotFound)
		return false
	}

	if productCode == "" {
		models.APIError(w, "No product-code", http.StatusBadRequest)
		return false
	}

	return true
}

func (c *bitflyerCandleController) APIHandler(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	fmt.Println(values)
	fmt.Println("url", r.URL)

	productCode := values.Get("product_code")
	if !c.validAPIHandler(w, r.URL.Path, productCode) {
		return
	}

	strLimit := values.Get("limit")
	limit, err := strconv.Atoi(strLimit)
	if strLimit == "" || err != nil || limit <= 0 || limit > 1000 {
		limit = 1000
	}

	duration := values.Get("duration")
	if duration == "" {
		duration = "1m"
	}
	cfg := config.Config()
	durationTime := cfg.Durations()[duration]

	c.dataFrameCandleService.SetDataFrame(productCode, durationTime, limit)

	if values.Get("sma") != "" {
		var periodStrs []string
		periodStrs = append(periodStrs, values.Get("sma_period1"))
		periodStrs = append(periodStrs, values.Get("sma_period2"))
		periodStrs = append(periodStrs, values.Get("sma_period3"))
		defaultPeriods := []int{7, 14, 50}

		for i, periodStr := range periodStrs {
			period, err := strconv.Atoi(periodStr)
			if periodStr == "" || err != nil || period < 0 {
				c.dataFrameCandleService.AddSMA(defaultPeriods[i])
			} else {
				c.dataFrameCandleService.AddSMA(period)
			}
		}
	}

	if values.Get("ema") != "" {
		var periodStrs []string
		periodStrs = append(periodStrs, values.Get("ema_period1"))
		periodStrs = append(periodStrs, values.Get("ema_period2"))
		periodStrs = append(periodStrs, values.Get("ema_period3"))
		defaultPeriods := []int{7, 14, 50}

		for i, periodStr := range periodStrs {
			period, err := strconv.Atoi(periodStr)
			if periodStr == "" || err != nil || period < 0 {
				c.dataFrameCandleService.AddEMA(defaultPeriods[i])
			} else {
				c.dataFrameCandleService.AddEMA(period)
			}
		}
	}

	if values.Get("bbands") != "" {
		periodStr := values.Get("bbands_period")
		kStr := values.Get("bbands_k")
		maType := values.Get("bbands_maType")

		period, err := strconv.Atoi(periodStr)
		if periodStr == "" || err != nil || period < 0 {
			period = 20
		}
		k, err := strconv.ParseFloat(kStr, 64)
		if kStr == "" || err != nil || k < 0 {
			k = 2
		}

		c.dataFrameCandleService.AddBBands(period, k, maType)
	}

	if values.Get("rsi") != "" {
		periodStr := values.Get("rsi_period")
		period, err := strconv.Atoi(periodStr)
		if periodStr == "" || err != nil || period < 0 {
			period = 14
		}
		c.dataFrameCandleService.AddRSI(period)
	}

	if values.Get("macd") != "" {
		var periodStrs []string
		periodStrs = append(periodStrs, values.Get("macd_fastPeriod"))
		periodStrs = append(periodStrs, values.Get("macd_slowPeriod"))
		periodStrs = append(periodStrs, values.Get("macd_signalPeriod"))
		defaultPeriods := []int{9, 26, 12}
		var periods []int

		for i, periodStr := range periodStrs {
			period, err := strconv.Atoi(periodStr)
			if periodStr == "" || err != nil || period < 0 {
				period = defaultPeriods[i]
			}

			periods = append(periods, period)
		}

		c.dataFrameCandleService.AddMACD(periods[0], periods[1], periods[2])
	}

	if values.Get("hv") != "" {
		var periodStrs []string
		periodStrs = append(periodStrs, values.Get("hv_period1"))
		periodStrs = append(periodStrs, values.Get("hv_period2"))
		periodStrs = append(periodStrs, values.Get("hv_period3"))
		defaultPeriods := []int{21, 63, 252}

		for i, periodStr := range periodStrs {
			period, err := strconv.Atoi(periodStr)
			if periodStr == "" || err != nil || period < 0 {
				c.dataFrameCandleService.AddHV(defaultPeriods[i])
			} else {
				c.dataFrameCandleService.AddHV(period)
			}
		}
	}

	if values.Get("event") != "" {
		t := c.DataFrame().Candles[0].Time
		c.dataFrameCandleService.AddEvents(t)
	}

	js, err := c.dataFrameCandleService.MarshalDataFrame()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
