package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/sdcoffey/big"
	"github.com/sdcoffey/techan"
)

func Mechanics(d DataSet, mc *MechConfig) Mechs {

	var m Mechs

	if len(d.StringData)-1 < 10 { // Need At Least 10 Data Points
		return m
	}

	series := techan.NewTimeSeries()

	var sPeriod int64

	for i, datum := range d.StringData {

		if i == 0 {
			sPeriod, _ = strconv.ParseInt(datum[0], 10, 64)
		}

		start, _ := strconv.ParseInt(datum[0], 10, 64)
		period := techan.NewTimePeriod(time.Unix(start, 0), time.Hour*24)

		candle := techan.NewCandle(period)
		candle.OpenPrice = big.NewFromString(datum[1])
		candle.ClosePrice = big.NewFromString(datum[2])
		candle.MaxPrice = big.NewFromString(datum[3])
		candle.MinPrice = big.NewFromString(datum[4])
		candle.Volume = big.NewFromString(datum[5])

		series.AddCandle(candle)

	}

	closePrices := techan.NewClosePriceIndicator(series)

	c := series.LastCandle()
	i := series.LastIndex()

	sigma := (c.MaxPrice.Float() + c.MinPrice.Float() + c.ClosePrice.Float()) / 3

	sma1 := techan.NewSimpleMovingAverage(closePrices, mc.SmaWindow.Short)
	sma2 := techan.NewSimpleMovingAverage(closePrices, mc.SmaWindow.Medium)
	sma3 := techan.NewSimpleMovingAverage(closePrices, mc.SmaWindow.Long)
	ema1 := techan.NewEMAIndicator(closePrices, mc.EmaWindow.Short)
	ema2 := techan.NewEMAIndicator(closePrices, mc.EmaWindow.Medium)
	ema3 := techan.NewEMAIndicator(closePrices, mc.EmaWindow.Long)

	rsi := techan.NewRelativeStrengthIndexIndicator(closePrices, mc.RsiWindow)
	macd := techan.NewMACDIndicator(closePrices, mc.MacdWindow.Short, mc.MacdWindow.Long)
	macdh := techan.NewMACDHistogramIndicator(macd, mc.MacdHistogramWindow)
	bub := techan.NewBollingerUpperBandIndicator(closePrices, mc.BollingerBandWindow, sigma)
	blb := techan.NewBollingerLowerBandIndicator(closePrices, mc.BollingerBandWindow, sigma)

	m = Mechs{
		Ticker:        d.Ticker,
		StartPeriod:   sPeriod,
		EndPeriod:     c.Period.End.Unix(),
		Rsi:           rsi.Calculate(i).Float(),
		Macd:          macd.Calculate(i).Float(),
		MacdHistogram: macdh.Calculate(i).Float(),
		Ema: EmaMech{
			Short:  ema1.Calculate(i).Float(),
			Medium: ema2.Calculate(i).Float(),
			Long:   ema3.Calculate(i).Float(),
		},
		Sma: SmaMech{
			Short:  sma1.Calculate(i).Float(),
			Medium: sma2.Calculate(i).Float(),
			Long:   sma3.Calculate(i).Float(),
		},
		Bollinger: BollingerMech{
			Upper: bub.Calculate(i).Float(),
			Lower: blb.Calculate(i).Float(),
		},
		LastPrice: series.LastCandle().ClosePrice.Float(),
		Length:    len(series.Candles),
		Config:    *mc,
		Success:   true,
	}

	return m
}

func Check(t string, m Metrics) (string, interface{}) {

	var c CheckData

	if t == "" {
		return c.Ticker, c.Metric
	}

	re, err := http.Get(fmt.Sprintf("https://www.alphavantage.co/query?function=%s&symbol=%s&interval=weekly&time_period=15&series_type=close&apikey=%s", m.String(), t, os.Getenv("KEY")))

	if err != nil {
		log.Fatal(err)
	}

	var d EMAData

	json.NewDecoder(re.Body).Decode(&d)

	defer re.Body.Close()

	var raw map[string]interface{}

	for key, v := range raw {
		if key == "2021-03-12" {
			return t, v
		}
	}

	return c.Ticker, c.Metric

}
