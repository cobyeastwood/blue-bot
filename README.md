# BlueBot

### Quick Note
BlueBot and all mentioned services are free to use, including supported financial APIs.

### Overview

BlueBot is a self-healing stock bot. You can get started by building a custom stock screener from [Free Stock Screener - Yahoo Finance](https://finance.yahoo.com/screener/new) and setting it in the global scope. 

This project supports the following:

1. Custom Stock Screening
2. Historical & Real-Time Data Collection
3. Technical Analysis (eg. Bollinger Bands, VWAP)
4. Custom Trade Strategies (eg. Backtesting) 
5. Conditional Trade Execution (WIP)

### How To Get Started
Root file main.go contains the general structure and configuration for this project. 

1. Custom Yahoo Finance stock screener endpoint can be attached on line [line 16](https://github.com/cobyeastwood/BlueBot/blob/master/main.go#L16).

	<br/>	

		YF_SCREENER_URL = "" // Set custom yahoo finance screener endpoint


	Note: global URL will live for only a month and will have recreated if used actively. A future OAuth tool will allow endpoints to add so that the Yahoo Finance screener will come from inside a Yahoo account.
	
	<br/>

2. Trading strategies can be easily implemented.
	
	<br/>
	
		// An example trade strategy
		
		func Strategy(mc *MainConfig) (string, interface{}) {
			t.Sort(2.00)

			m1 := Short(t, mc.c)

			m2 := SiftFrom(m1)

			m2 := Mechanics(FetchLong(m1.Ticker, NewFetchConfig(HOUR, "1", 365)), NewLongMechConfig())

			Check(m1.Ticker, EMA)
		}
	
	Later a custom strategy can be added and placed inside the following code block on [line 48](https://github.com/cobyeastwood/BlueBot/blob/master/main.go#L48).
	<br/>
	
		for {

			t := <-mc.tOut

			Strategy(mc) // Custom trade strategy goes here

			time.Sleep(60 * time.Second)

		}

	<br/> 
	

### Historical Support

For historical data collection, you can choose from three different services: [Alpha Vantage](https://www.alphavantage.co/), [IEX Cloud](https://iexcloud.io/?gclid=CjwKCAjwuIWHBhBDEiwACXQYsRZK32T9FfG4LsdaTr8IvUFY9LnJG-KAQkrjIzkSzMQ1O3u90Z-QhRoCzQ0QAvD_BwE), and [Polygon.io](https://polygon.io/stocks?gclid=CjwKCAjwuIWHBhBDEiwACXQYsWGZBgzKC7eFBdpJUEYbqBgjqXkfoYtUUkwsIsBjF_n_hfQyGeJisRoCEZMQAvD_BwE). Any API keys will need to be placed inside a local env file.

### Real-Time Support

For real-time data collection, you will need to collect API keys from [Alpaca](https://alpaca.markets/docs/about-us/), an algorithmic stock trading for developers with margin availability and fractional shares. Any API keys will need to be placed inside a local env file.
