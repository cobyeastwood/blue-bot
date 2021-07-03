package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/alpacahq/alpaca-trade-api-go/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/common"
	"golang.org/x/oauth2"
)

// Break down min funcs to an interface
interface Min(){
	func GetMin([]map[string]Mechs, func(key string, value Mechs) bool{}) Mechs
}

func NewFetchConfig(u Units, v string, d int) *FetchConfig {
	return &FetchConfig{
		Unit:  u,
		Value: v,
		Start: time.Now().Format("2006-01-02"),
		End:   time.Now().Add(-time.Hour * 24 * time.Duration(d)).Format("2006-01-02"),
	}
}

func NewLongMechConfig() *MechConfig {
	return &MechConfig{
		EmaWindow: Windows{
			50,
			100,
			200,
		},
		SmaWindow: Windows{
			50,
			100,
			200,
		},
		RsiWindow: 14,
		MacdWindow: Windows{
			Short: 12,
			Long:  26,
		},
		MacdHistogramWindow: 9,
		BollingerBandWindow: 20,
	}
}

func NewMainConfig() *MainConfig {

	if common.Credentials().ID == "" {
		os.Setenv(common.EnvApiKeyID, os.Getenv("KEY3ID"))
	}

	if common.Credentials().Secret == "" {
		os.Setenv(common.EnvApiSecretKey, os.Getenv("KEY3"))
	}

	return &MainConfig{
		c:    alpaca.NewClient(common.Credentials()),
		m:    NewMarketInfo(),
		y:    NewYahooOAuth(),
		mOut: make(chan *MarketInfo, 1),
		tOut: make(chan Ticks, 1),
	}

}

func NewYahooOAuth() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("YAHOO_ID"),
		ClientSecret: os.Getenv("YAHOO_SECRET"),
		Scopes:       []string{"yfin-w"},
		RedirectURL:  "https://localhost:8090/oauth2/yahoo/receive",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://api.login.yahoo.com/oauth2/request_auth",
			TokenURL: "https://api.login.yahoo.com/oauth2/get_token",
		},
	}
}

func NewMarketInfo() *MarketInfo {
	return &MarketInfo{}
}

func NewShortMechConfig() *MechConfig {
	return &MechConfig{
		EmaWindow: Windows{
			5,
			8,
			13,
		},
		SmaWindow: Windows{
			5,
			8,
			13,
		},
		RsiWindow: 14,
		MacdWindow: Windows{
			Short: 5,
			Long:  35,
		},
		MacdHistogramWindow: 5,
		BollingerBandWindow: 20,
	}
}

func Long(t Ticks, fc *FetchConfig) []map[string]Mechs {

	var m []map[string]Mechs

	mc := NewLongMechConfig()

	for key := range t.Stats {
		var ts Techs

		ts.Set()

		ts.Long[key] = Mechanics(FetchLong(key, fc), mc)

		if ms, b := IndicatorLong(ts); b {
			m = append(m, ms)
		}
	}

	return m
}

func MaxEma(m []map[string]Mechs) Mechs {

	var max float64 = math.Inf(-1)

	var str string
	var idx int

	for i, e := range m {
		for key, v := range e {
			if v.Ema.Short > v.Ema.Medium && v.Ema.Medium > v.Ema.Long {
				if c := (v.Ema.Short - v.Ema.Medium - v.Ema.Long) / v.LastPrice; c >= max {
					max = c
					idx = i
					str = key
				}
			}
		}
	}

	return m[idx][str]

}

/* Substitute Better Sorting Methods */
func Sift(m []map[string]Mechs) Mechs {

	var min float64 = math.Inf(0)

	var str string
	var idx int

	for i, e := range m {
		for key, v := range e {
			if v.Ema.Short > v.Ema.Medium && v.Ema.Medium > v.Ema.Long && v.Rsi < 70.00 && v.LastPrice <= v.LastPrice+((v.Bollinger.Upper-v.Bollinger.Lower)/3) {
				if c := (v.Ema.Short - v.Ema.Medium - v.Ema.Long) / v.LastPrice; c <= min {
					min = c
					idx = i
					str = key
				}
			}
		}
	}

	return m[idx][str]

}

func MinBBand(m []map[string]Mechs) Mechs {
	var min float64 = math.Inf(0)
	var str string

	var idx int

	for i, e := range m {
		for key, v := range e {
			if (v.LastPrice-v.Bollinger.Lower)/v.LastPrice <= min {
				idx = i
				min = v.Rsi
				str = key
			}
		}
	}

	return m[idx][str]

}

func MinRsi(m []map[string]Mechs) Mechs {
	var min float64 = math.Inf(0)
	var str string

	var idx int

	for i, e := range m {
		for key, v := range e {
			if v.Rsi <= min {
				idx = i
				min = v.Rsi
				str = key
			}
		}
	}

	return m[idx][str]

}

func MinEma(m []map[string]Mechs) Mechs {

	var str string
	var idx int

	for i, e := range m {
		for key, v := range e {
			if v.Ema.Short < v.Ema.Medium || v.Ema.Medium < v.Ema.Long {
				idx = i
				str = key
			}
		}
	}

	return m[idx][str]

}

func StringConv(i interface{}) string {

	var s string

	switch v := i.(type) {
	case int:
		s = strconv.Itoa(v)
	case int32:
		s = strconv.Itoa(int(v))
	case int64:
		s = strconv.FormatInt(v, 10)
	case float32:
		s = strconv.FormatFloat(float64(v), 'E', -1, 32)
	case float64:
		s = strconv.FormatFloat(v, 'E', -1, 64)
	}

	return s
}

func StringUpper(s string) string {

	if len(s) == 0 {
		return s
	}

	var sb strings.Builder

	for _, r := range s {
		if ok := unicode.IsLower(r); ok {
			sb.WriteRune(unicode.ToUpper(r))
		} else {
			sb.WriteRune(r)
		}
	}

	return sb.String()

}

func Short(t Ticks, ac *alpaca.Client) []map[string]Mechs {

	var m []map[string]Mechs

	mc := NewShortMechConfig()

	for key := range t.Stats {

		if key != "" {
			var ts Techs

			ts.Set()

			ts.Short[key] = Mechanics(FetchShort(key, ac), mc)

			if ms, b := IndicatorShort(ts); b {
				m = append(m, ms)
			}
		}

	}

	return m
}

func Status(mc *MainConfig) {

	for range time.After(mc.m.Wait * time.Second) {

		mc.mu.Lock()

		if c, e := mc.c.GetClock(); e == nil {
			mc.m.IsOpen = c.IsOpen
			mc.m.Wait = time.Duration(c.NextOpen.Sub(c.Timestamp).Seconds())
			mc.mOut <- mc.m
		} else {
			mc.m.IsOpen = false
			mc.m.Wait = 3600
			mc.mOut <- mc.m
		}

		mc.mu.Unlock()

	}

}
