package bourse

import (
	"IranStocksCrawler/helpers/arraysh"
	"IranStocksCrawler/system/cacher"
	"strconv"
	"time"
)

type FilterParams struct {
	GeneralInfo         *StockPrice
	StockIOInfo         *StockIO
	StockTodaySeries    *StockTodaySeries
	IndiOrga365DaysList *StockIndiOrga365Days
}

type filterFunc func(symbol string, params FilterParams, lastOpenDays []string) bool

var filters map[string]*filterFunc = make(map[string]*filterFunc, 0)

var filterResult map[string][]string = make(map[string][]string, 0)

func ApplyFilters(cacher cacher.ICacherEngine) bool {

	days := createDays()

	for sym := range stockPriceList {

		params := FilterParams{
			GeneralInfo:         stockPriceList[sym],
			StockIOInfo:         stockIOList[sym],
			StockTodaySeries:    stockTodaySeriesList[sym],
			IndiOrga365DaysList: stockIndiOrga365DaysList[sym],
		}

		doFilters(sym, params, days)
	}

	// store result of each filter
	for name, res := range filterResult {
		cacher.Put(name, res, 24*60*60)
	}

	return true
}

func doFilters(symbol string, params FilterParams, days []string) {

	for name, f := range filters {

		ff := *f

		ok := ff(symbol, params, days)

		if ok == true {
			filterResult[name] = append(filterResult[name], symbol)
		}

	}
}

func addFilter(name string, filter filterFunc) {

	filters[name] = &filter
}

func InitFilters() {

	if len(filters) > 0 {
		return
	}

	// shakhes-saz-ha
	addFilter("shakhes-saz-ha", func(symbol string, params FilterParams, days []string) bool {

		shakhesSymbols := []string{
			"شپنا", "شبندر", "شتران", "شستا", "فارس", "فولاد",
			"فملی", "وبصادر", "خپارس", "تاپیکو", "وپارس",
			"وتجارت", "ذوب", "کچاد", "وغدیر", "اخابر", "پارسان",
			"خساپا", "خگستر", "خودرو", "وبملت", "فخوز", "رمپنا",
			"وامید", "پارس", "همراه", "وصندوق", "شپدیس", "کگل", "ومعادن"}

		if arraysh.Contains(shakhesSymbols, symbol) {
			return true
		}

		return false
	})

	// doubtful volume
	addFilter("doubtful-volume", func(symbol string, params FilterParams, days []string) bool {

		if params.IndiOrga365DaysList == nil {
			return false
		}

		todayV := params.GeneralInfo.Volume

		sumVolume := int64(0)

		i := 0
		for _, j := range days {
			i++

			indiorga, exist := params.IndiOrga365DaysList.Days[j]
			if !exist {
				continue
			}

			vib, _ := strconv.ParseInt(indiorga.VolumeIndiBuy, 10, 64)
			vob, _ := strconv.ParseInt(indiorga.VolumeOrgaBuy, 10, 64)

			sumVolume += (vib + vob)

			if i >= 7 {
				break
			}
		}

		if todayV > 2*sumVolume/int64(len(days)) {
			return true
		}

		return false
	})

	// saraneh kharid ascending
	addFilter("saraneh-kharid-ascending", func(symbol string, params FilterParams, days []string) bool {

		if params.IndiOrga365DaysList == nil {
			return false
		}

		//todayV := params.GeneralInfo.Volume
		todaySaranehBuy := params.StockIOInfo.IndiBuySaraneh
		//todayV := params.GeneralInfo.Volume

		//sumVolume := int64(0)

		maxSaranehBuy := int64(0)

		i := 0
		for _, j := range days {
			i++

			indiorga, exist := params.IndiOrga365DaysList.Days[j]
			if !exist {
				continue
			}
			//vib, _ := strconv.ParseInt(indiorga.VolumeIndiBuy, 10, 64)
			//vob, _ := strconv.ParseInt(indiorga.VolumeOrgaBuy, 10, 64)

			aib, _ := strconv.ParseInt(indiorga.AmountIndiBuy, 10, 64)
			qib, _ := strconv.ParseInt(indiorga.QuantityIndiBuy, 10, 64)

			//ais, _ := strconv.ParseInt(indiorga.AmountIndiSell, 10, 64)
			//qis, _ := strconv.ParseInt(indiorga.QuantityIndiSell, 10, 64)

			saranehBuy := (aib / (qib * 10000000))
			//saranehSell := (ais / qis)

			//buyPower := saranehBuy / saranehSell
			if saranehBuy > maxSaranehBuy {
				maxSaranehBuy = saranehBuy
			}

			if i >= 7 {
				break
			}
		}

		if todaySaranehBuy > (1.5)*float64(maxSaranehBuy) {
			return true
		}

		return false
	})

}

func createDays() []string {

	days := []string{}

	days = append(days, time.Now().Format(DEF_DATE_PATTERN))
	days = append(days, time.Now().Truncate(1*24*60*60).Format(DEF_DATE_PATTERN))
	days = append(days, time.Now().Truncate(2*24*60*60).Format(DEF_DATE_PATTERN))
	days = append(days, time.Now().Truncate(3*24*60*60).Format(DEF_DATE_PATTERN))
	days = append(days, time.Now().Truncate(4*24*60*60).Format(DEF_DATE_PATTERN))
	days = append(days, time.Now().Truncate(5*24*60*60).Format(DEF_DATE_PATTERN))
	days = append(days, time.Now().Truncate(6*24*60*60).Format(DEF_DATE_PATTERN))
	days = append(days, time.Now().Truncate(7*24*60*60).Format(DEF_DATE_PATTERN))
	days = append(days, time.Now().Truncate(8*24*60*60).Format(DEF_DATE_PATTERN))
	days = append(days, time.Now().Truncate(9*24*60*60).Format(DEF_DATE_PATTERN))
	days = append(days, time.Now().Truncate(10*24*60*60).Format(DEF_DATE_PATTERN))

	return days
}

func GetFilterResult(filter string) []string {
	result := filterResult[filter]
	return result
}
