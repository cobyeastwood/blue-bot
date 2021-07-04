# BlueBot


### Description

BlueBot is a self-healing stock bot that is still under construction. You can get started by building a custom stock screener from [Free Stock Screener - Yahoo Finance](https://finance.yahoo.com/screener/new) and setting it in the global scope. 

This project supports the following:

1. Custom Stock Screening
2. Historical & Real-Time Data Collection
3. Technical Analysis (eg. Bollinger Bands, VWAP)
4. Custom Trade Strategies (eg. Backtesting) 
5. Conditional Trade Execution (WIP)

For real-time data collection and trading, you will need to collect API keys from Alpaca, a popular trading API. If so desired, you can choose to additionally add other supported services such as Alpha Vantage, IEX Cloud, and Polygon. Any keys will need to be placed inside a local env file.

### Overview
The main.go file contains a basic configuration for this project. 

1. Custom Yahoo Finance screener url can be attached on line [line 16](https://github.com/cobyeastwood/BlueBot/blob/master/main.go#L16).

	<br/>	

		YF_SCREENER_URL = "" // Set custom yahoo finance screener endpoint


	Note: global URL will live for only a month and will have recreated if used actively. A future OAuth tool will allow endpoints to add so that the Yahoo Finance screener will come from inside a Yahoo account.
	
	<br/>

2. Trading strategies can be placed inside the following code block on [line 43](https://github.com/cobyeastwood/BlueBot/blob/master/main.go#L43).
	
	<br/>
	
		// An example trade strategy
		
		func Strategy(mc *MainConfig) (string, interface{}) {
			t.Sort(2.00)

			m1 := Short(t, mc.c)

			m2 := SiftFrom(m1)

			m2 := Mechanics(FetchLong(m1.Ticker, NewFetchConfig(HOUR, "1", 365)), NewLongMechConfig())

			Check(m1.Ticker, EMA)
		}
	
	Later a custom strategy can be added.
	<br/>
	
		for {

			t := <-mc.tOut

			Strategy(mc) // Custom trade strategy goes here

			time.Sleep(60 * time.Second)

		}

	<br/> 
	
