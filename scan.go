package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func Scan(mc *MainConfig) {

	var t Ticks

	re, err := http.Get(YF_SCREENER_URL)

	if err != nil {
		log.Fatal(err)
	}

	defer re.Body.Close()

	switch re.StatusCode {
	case 200:
		var ticks []string

		doc, err := goquery.NewDocumentFromReader(re.Body)

		if err != nil {
			log.Fatal(err)
		}

		doc.Find("table").Each(func(i int, tablehtml *goquery.Selection) {
			tablehtml.Find("td").Each(func(i int, tablerow *goquery.Selection) {
				tablerow.Find("a").Each(func(i int, td *goquery.Selection) {
					ticks = append(ticks, td.Text())
				})
			})
		})

		rOut := make(chan RawStat, len(ticks))

		for _, t := range ticks {
			go ScanTicks(t, rOut)
		}

		t.Set()

		for range ticks {
			select {
			case r := <-rOut:
				t.Stats[r.Symbol] = Stat{Tick: r.Symbol, VolumePercent: r.RelativeVolume(), PreVolume: r.PreviousVolume, AvgVolume: r.AverageVolume, IexPercent: r.IexPercent}
			}
		}
		break

	case 429:
		time.Sleep(1 * time.Second)
		break


	mc.tOut <- t

}

func ScanTicks(tick string, rOut chan<- RawStat) {

	var r RawStat

	re, err := http.Get(fmt.Sprintf("https://cloud.iexapis.com/stable/stock/%s/quote?token="+os.Getenv("KEY4"), tick))

	if err != nil {
		log.Fatal(err)
	}

	defer re.Body.Close()

	json.NewDecoder(re.Body).Decode(&r)

	rOut <- r

}

// https://api.login.yahoo.com/oauth2/request_auth?client_id=dj0yJmk9WGx0QlE0UWdCa0hKJmQ9WVdrOWNrNUhXVnBhTkhFbWNHbzlNQS0tJnM9Y29uc3VtZXJzZWNyZXQmeD01OA--&response_type=code&redirect_uri=https://yahoo.com&scope=openid%20mail-r&nonce=YihsFwGKgt3KJUh6tPs2
// https://api.login.yahoo.com/oauth2/request_auth?client_id=dj0yJmk9YWFGS3RhQjhpQXpOJmQ9WVdrOVRFWTFXRzVrTURnbWNHbzlNQT09JnM9Y29uc3VtZXJzZWNyZXQmc3Y9MCZ4PTVk&redirect_uri=https%3A%2F%2Ffinance.yahoo.com%2Fscreener%2F2e87646d-0831-4e68-9e53-b5689a4f2609&response_type=code
