package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/alpaca"
)

type DataSet struct {
	Ticker     string
	StringData [][]string
}

func FetchLong(tick string, fc *FetchConfig) DataSet {

	var d DataSet

	d.Ticker = StringUpper(tick)

	re, err := http.Get(fmt.Sprintf("https://api.polygon.io/v2/aggs/ticker/%s/range/%s/%s/%s/%s?unadjusted=true&sort=asc&limit=50000&apiKey="+os.Getenv("KEY"), d.Ticker, fc.Value, fc.Unit, fc.End, fc.Start))

	if err != nil {
		log.Fatal(err)
	}

	defer re.Body.Close()

	switch re.StatusCode {
	case 200:
		var r RawTicks

		json.NewDecoder(re.Body).Decode(&r)

		for _, s := range r.Results {
			var tempStr []string

			tempStr = append(tempStr, StringConv(s.Timestamp), StringConv(s.Open), StringConv(s.Close), StringConv(s.High), StringConv(s.Low), StringConv(s.Volume))

			d.StringData = append(d.StringData, tempStr)
		}
		break
	case 429:
		time.Sleep(1 * time.Second)
		break

	return d

}

func FetchShort(tick string, ac *alpaca.Client) DataSet {

	var d DataSet

	d.Ticker = StringUpper(tick)

	now, end := time.Unix(time.Now().Unix()-int64(420*60), 0), time.Now()

	bars, err := ac.GetSymbolBars(tick, alpaca.ListBarParams{Timeframe: "minute", StartDt: &now, EndDt: &end})

	if err != nil {
		log.Fatal(err)
	}

	for _, b := range bars {
		var tempStr []string

		tempStr = append(tempStr, StringConv(b.Time*1000), StringConv(b.Open), StringConv(b.Close), StringConv(b.High), StringConv(b.Low), StringConv(b.Volume))

		d.StringData = append(d.StringData, tempStr)
	}

	return d

}
