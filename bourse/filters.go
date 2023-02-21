package bourse

import (
	"IranStocksCrawler/helpers/arraysh"
	"IranStocksCrawler/system/cacher"
)

type filterFunc func(symbol string, spr *StockPrice, ts *StockTodaySeries) bool

var filters map[string]*filterFunc = make(map[string]*filterFunc, 0)

var filterResult map[string][]string = make(map[string][]string, 0)

func ApplyFilters(cacher cacher.ICacherEngine) bool {

	for sym := range stockPriceList {
		spr := stockPriceList[sym]
		//sir := stockIORowList[sym]
		series := stockTodaySeriesList[sym]

		doFilters(sym, spr, series)
	}

	// store result of each filter
	for name, res := range filterResult {
		cacher.Put(name, res, 24*60*60)
	}

	return true
}

func doFilters(symbol string, spr *StockPrice, ts *StockTodaySeries) {

	for name, f := range filters {

		ff := *f

		ok := ff(symbol, spr, ts)

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
	addFilter("shakhes-saz-ha", func(symbol string, spr *StockPrice, series *StockTodaySeries) bool {

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
	addFilter("doubtful-volume", func(symbol string, spr *StockPrice, series *StockTodaySeries) bool {

		return false
	})

	// saraneh kharid so'udi
	addFilter("saraneh-kharid-so'udi", func(symbol string, spr *StockPrice, series *StockTodaySeries) bool {

		return false
	})

}
