package main

import (
	"fmt"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v2/stream"
)

func ShortStream(ticker string, sOut chan<- []string) {

	if err := stream.SubscribeBars(func(bar stream.Bar) {
		var str []string

		str = append(str, StringConv(bar.Timestamp.UnixNano()/int64(time.Millisecond)), StringConv(bar.Open), StringConv(bar.Close), StringConv(bar.High), StringConv(bar.Low), StringConv(bar.Volume))

		sOut <- str
	}, ticker); err != nil {
		panic(err)
	}

	if err := stream.SubscribeQuotes(func(quote stream.Quote) {}, ticker); err != nil {
		panic(err)
	}

	select {}
}
