/*
	@name: BlueBot
	@desc: Self-Trading Stock Bot
	@date: 2021-03-08
	@vers: 1.0.0
*/

package main

import (
	"fmt"
	"time"
)

const (
	PORT            = "" // Set oauth port
	YF_SCREENER_URL = "" // Set custom yahoo finance screener endpoint
)

// An example trade strategy
func Strategy(mc *MainConfig) (string, interface{}) {
	t.Sort(2.00)

	m1 := Short(t, mc.c)

	m2 := Sift(m1)

	m2 := Mechanics(FetchLong(m1.Ticker, NewFetchConfig(HOUR, "1", 365)), NewLongMechConfig())

	Check(m1.Ticker, EMA)
}

func main() {

	mc := NewMainConfig()

	// Start oauth service
	// OAuth(mc, PORT)

	go Status(mc)

	for info := range mc.mOut {
		if info.IsOpen {
			fmt.Println("Bot is Ready. Starting...")

			go Scan(mc)

			for {

				t := <-mc.tOut

				Strategy(mc) // Custom trade strategy goes here

				time.Sleep(60 * time.Second)

			}

		} else {
			fmt.Println("Bot is Waiting...", mc.m.Wait)

			go Status(mc)
		}

	}

}

// ########## GoDoc References ##########
// https://github.com/sdcoffey/techan
// https://alpaca.Statuss/docs/api-documentation/api-v2/
