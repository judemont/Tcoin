package main

import (
	"time"
)

type Crypto struct {
	ID                       string
	Name                     string
	Symbol                   string
	Price                    float64
	LogoURL                  string
	PriceChangePercentageDay float64
	Description              string
	Categories               []string
	Website                  string
	ATH                      float64
	ATHDate                  time.Time
	MarketCap                float64
	MarketCapRank            int
	DayHigh                  float64
	DayLow                   float64
	TotalSupply              float64
	CirculatingSupply        float64
	Volume                   float64
}
