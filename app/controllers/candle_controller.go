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

func NewBitflyerCandleController(db *sql.DB) CandleController {
	candleRepository := repositories.NewBitflyerCandleRepository(db)
	dataFrameCandleService := services.NewDataFrameCandleService(candleRepository)

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

	js, err := c.dataFrameCandleService.MarshalDataFrame()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
