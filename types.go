package main

import (
	"math"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/alpaca"
	"golang.org/x/oauth2"
)

const (
	MINUTE Units = iota + 1
	HOUR
	DAY
)

const (
	SMA Metrics = iota + 1
	EMA
	RSI
	VWAP
	MACD
	BBANDS
)

type CheckData struct {
	Ticker string
	Metric interface{}

	EMA    EMAData
	SMA    SMAData
	RSI    RSIData
	MACD   MACDData
	BBANDS BBANDSData
	VWAP   VWAPData
}

type EMAData struct {
	Raw map[string]interface{} `json:"Technical Analysis: EMA"`
}
type RSIData struct {
	Raw map[string]interface{} `json:"Technical Analysis: RSI"`
}
type SMAData struct {
	Raw map[string]interface{} `json:"Technical Analysis: SMA"`
}
type MACDData struct {
	Raw map[string]interface{} `json:"Technical Analysis: MACD"`
}
type BBANDSData struct {
	Raw map[string]interface{} `json:"Technical Analysis: BBANDS"`
}
type VWAPData struct {
	Raw map[string]interface{} `json:"Technical Analysis: VWAP"`
}

type BollingerMech struct {
	Upper float64
	Lower float64
}

type EmaMech struct {
	Short  float64
	Medium float64
	Long   float64
}

type FetchConfig struct {
	Unit  Units
	Value string
	Start string
	End   string
}

type MainConfig struct {
	mu sync.Mutex

	c    *alpaca.Client
	m    *MarketInfo
	y    *oauth2.Config
	yc   *http.Client
	yd   time.Time
	mOut chan *MarketInfo
	tOut chan Ticks
}

type MarketInfo struct {
	IsOpen bool
	Wait   time.Duration
}

type Mechs struct {
	Ticker        string
	StartPeriod   int64
	EndPeriod     int64
	Rsi           float64
	Macd          float64
	MacdHistogram float64
	Ema           EmaMech
	Sma           SmaMech
	Bollinger     BollingerMech

	LastPrice float64
	Length    int
	Config    MechConfig
	Success   bool
}

type Metrics int

func (m Metrics) String() string {
	return [...]string{"SMA", "EMA", "RSI", "VWAP", "MACD", "BBANDS"}[m]
}

type MechConfig struct {
	RsiWindow           int
	MacdHistogramWindow int
	BollingerBandWindow int
	EmaWindow           Windows
	SmaWindow           Windows
	MacdWindow          Windows
}

type Units int

func (u Units) String() string {
	return [...]string{"minute", "hour", "day"}[u]
}

type Techs struct {
	Short map[string]Mechs
	Long  map[string]Mechs
}

func (t *Techs) Set() {
	t.Short = make(map[string]Mechs)
	t.Long = make(map[string]Mechs)
}

type Ticks struct {
	Stats       map[string]Stat
	HighVolumes map[string]Stat
	HighGains   map[string]Stat
	MaxVolume   string
	MinVolume   string
}

func (t *Ticks) Max() {
	var max float64 = math.Inf(-1)
	var str string

	for key, v := range t.HighVolumes {
		if v.IexPercent >= max {
			max = v.IexPercent
			str = key
		}
	}

	t.MaxVolume = str

}

func (t *Ticks) Min() {
	var min float64 = math.Inf(0)
	var str string

	for key, v := range t.HighVolumes {
		if v.IexPercent <= min {
			min = v.IexPercent
			str = key
		}
	}

	t.MinVolume = str

}

func (t *Ticks) Set() {
	t.Stats = make(map[string]Stat)
	t.HighGains = make(map[string]Stat)
	t.HighVolumes = make(map[string]Stat)
}

func (t *Ticks) Sort(index float64) {

	var adds int
	var todo []float64

	for key, v := range t.Stats {
		if v.VolumePercent > index {
			t.HighVolumes[key] = v
			adds++
		} else {
			todo = append(todo, v.VolumePercent)
		}
	}

	if len(todo) != 0 && adds < 3 {
		sort.Float64s(todo)

		for key, v := range t.Stats {
			switch v.VolumePercent {
			case todo[len(todo)-1]:
				t.HighVolumes[key] = v
			case todo[len(todo)-2]:
				t.HighVolumes[key] = v
			case todo[len(todo)-3]:
				t.HighVolumes[key] = v
			}
		}
	}

}

type rawStats []RawStat

type RawTicks struct {
	Ticker  string `json:"ticker"`
	Results []struct {
		Open         float32 `json:"o"`
		High         float32 `json:"h"`
		Low          float32 `json:"l"`
		Close        float32 `json:"c"`
		Volume       float32 `json:"v"`
		VolumeWeight float32 `json:"vw"`
		Timestamp    int     `json:"t"`
	} `json:"results"`
}

type RawStat struct {
	Symbol           string  `json:"symbol"`
	Company          string  `json:"companyName"`
	Exchange         string  `json:"primaryExchange"`
	CalculationPrice string  `json:"calculationPrice"`
	Open             int     `json:"open"`
	OpenTime         int     `json:"openTime"`
	Close            int     `json:"close"`
	CloseTime        int     `json:"closeTime"`
	High             int     `json:"high"`
	HighTime         int     `json:"highTime"`
	Low              int     `json:"low"`
	LowTime          int     `json:"lowTime"`
	PreviousVolume   int     `json:"previousVolume"`
	Volume           int     `json:"volume"`
	IexVolume        int     `json:"iexVolume"`
	AverageVolume    int     `json:"avgTotalVolume"`
	PE               int     `json:"peRatio"`
	Change           int     `json:"change"`
	Percent          int     `json:"changePercent"`
	IexPercent       float64 `json:"iexStatusPercent"`
	Week52High       int     `json:"week52High"`
	Week52Low        int     `json:"week52Low"`
	YtdChange        int     `json:"ytdChange"`
}

func (s *RawStat) RelativeVolume() float64 {
	return float64(s.PreviousVolume) / float64(s.AverageVolume)
}

type SmaMech struct {
	Short  float64
	Medium float64
	Long   float64
}

type Stat struct {
	Tick          string
	VolumePercent float64
	AvgVolume     int
	PreVolume     int
	IexPercent    float64
}

type Windows struct {
	Short  int
	Medium int
	Long   int
}
