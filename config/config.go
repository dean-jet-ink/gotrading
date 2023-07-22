package config

import (
	"log"
	"time"

	"gopkg.in/ini.v1"
)

type configList struct {
	apiTimeout    int
	logFile       string
	tradeDuration time.Duration
	durations     map[string]time.Duration
	apiKey        string
	apiSecret     string
	productCode   string
	dbName        string
	dbDriver      string
	port          int
	backtest      bool
	usePercent    float64
	stopPercent   float64
	dateLimit     int
}

var config *configList

func init() {
	cfg, err := ini.Load("config/config.ini")
	if err != nil {
		log.Println(err)
		return
	}

	durations := map[string]time.Duration{
		"1h": time.Hour,
		"1m": time.Minute,
		"1s": time.Second,
	}

	config = &configList{
		apiTimeout:    cfg.Section("api").Key("timeout").MustInt(),
		logFile:       cfg.Section("gotrading").Key("log_file").String(),
		productCode:   cfg.Section("gotrading").Key("product_code").String(),
		tradeDuration: durations[cfg.Section("gotrading").Key("trade_duration").String()],
		durations:     durations,
		backtest:      cfg.Section("gotrading").Key("backtest").MustBool(),
		usePercent:    cfg.Section("gotrading").Key("use_percent").MustFloat64(),
		stopPercent:   cfg.Section("gotrading").Key("stop_percent").MustFloat64(),
		dateLimit:     cfg.Section("gotrading").Key("date_limit").MustInt(),
		apiKey:        cfg.Section("bitflyer").Key("api_key").String(),
		apiSecret:     cfg.Section("bitflyer").Key("api_secret").String(),
		dbName:        cfg.Section("db").Key("name").String(),
		dbDriver:      cfg.Section("db").Key("driver").String(),
		port:          cfg.Section("web").Key("port").MustInt(),
	}
}

func Config() *configList {
	return config
}

func (c *configList) APITimeout() int {
	return c.apiTimeout
}

func (c *configList) LogFile() string {
	return c.logFile
}

func (c *configList) TradeDuration() time.Duration {
	return c.tradeDuration
}

func (c *configList) Durations() map[string]time.Duration {
	return c.durations
}

func (c *configList) ApiKey() string {
	return c.apiKey
}

func (c *configList) ApiSecret() string {
	return c.apiSecret
}

func (c *configList) ProductCode() string {
	return c.productCode
}

func (c *configList) DBName() string {
	return c.dbName
}

func (c *configList) DBDriver() string {
	return c.dbDriver
}

func (c *configList) Port() int {
	return c.port
}

func (c *configList) Backtest() bool {
	return c.backtest
}

func (c *configList) UsePercent() float64 {
	return c.usePercent
}

func (c *configList) StopPercent() float64 {
	return c.stopPercent
}

func (c *configList) DateLimit() int {
	return c.dateLimit
}
