package main

func IndicatorShort(t Techs) (map[string]Mechs, bool) {

	var s map[string]Mechs
	var status bool

	s = make(map[string]Mechs, 1)

	for key, v := range t.Short {
		if v.Success && v.Rsi != 0 && v.Rsi < 80.00 {
			s[key] = v
			status = true
		}
	}

	return s, status
}

func IndicatorLong(t Techs) (map[string]Mechs, bool) {

	var s map[string]Mechs
	var status bool

	s = make(map[string]Mechs, 1)

	for key, v := range t.Long {
		if v.Success && v.Rsi != 0 && v.Rsi < 80.00 {
			s[key] = v
			status = true
		}
	}

	return s, status
}
