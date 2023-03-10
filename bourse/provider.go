package bourse

import (
	"IranStocksCrawler/helpers/stringsh"
	"IranStocksCrawler/helpers/timeh"
	"IranStocksCrawler/system/cacher"
	"fmt"
	"math"
	"regexp"
	"strings"
	"time"
)

var consecutiveAtempts map[string]int = make(map[string]int)

var currentPriceContent string
var currentIOContent string
var currentPeriodicAveragesContent string

var stockPriceList map[string]*StockPrice
var stockIOList map[string]*StockIO
var stockPeriodicAveragesList map[string]*StockAverages
var stockAskBidTableList map[string]*StockAskBidTable
var stockTodaySeriesList map[string]*StockTodaySeries = map[string]*StockTodaySeries{}

var stockTseTOSymbolList map[string]string

//

//var todayStocks = map[string]TodayStock{"12": {LastPrice: LineInt{1: 4, 2: 44}}}

type StockTodaySeries struct {
	LastPrice       LineInt
	BuyVolume       LineInt
	SellVolume      LineInt
	IndiBuySaraneh  LineFloat
	IndiSellSaraneh LineFloat
}

type LineInt map[int64]int64
type LineFloat map[int64]float64

type StockPrice struct {
	Symbol            string
	CompanyName       string
	TSE_Code          string
	IR_Code           string
	Time              string
	ClosePrice        int64
	ClosePriceRate    float64
	LastPrice         int64
	LastPriceRate     float64
	MaxPrice          int64
	MinPrice          int64
	RangeTopPrice     int64
	RangeBottomPrice  int64
	YesterdayPrice    int64
	TodayFirstPrice   int64
	BuyPriceAtFirst   int64
	SellPriceAtFirst  int64
	Volume            int64
	BaseVolume        int64
	BuyVolumeAtFirst  int64
	SellVolumeAtFirst int64
	FullAmount        int64
	Quantity          int64
	Shares            int64
}

type StockIO struct {
	Symbol           string
	CompanyName      string
	TSE_Code         string
	Time             string
	LastPrice        int64
	LastPriceRate    float64
	IndiBuyQuantity  int64
	OrgaBuyQuantity  int64
	IndiBuyVolume    int64
	OrgaBuyVolume    int64
	IndiSellQuantity int64
	OrgaSellQuantity int64
	IndiSellVolume   int64
	OrgaSellVolume   int64
	IndiBuySaraneh   float64
	IndiSellSaraneh  float64
	IndiBuyPower     float64
	IndiBuyPercent   float64
	OrgaSellPercent  float64
}
type StockAverages struct {
	Symbol                  string
	TSE_Code                string
	Last3MonthAverageAmount int64
	Last3MonthAverageVolume int64
}

type StockAskBidTable struct {
	BuyRows  map[int64]AskBidTableRow
	SellRows map[int64]AskBidTableRow
}

type AskBidTableRow struct {
	Volume   int64
	Quantity int
}

func UpdatePrices(cacher *cacher.Cacher) bool {

	err := gatherPrices()

	if err == nil {

		providePriceDetails(cacher)

		provideAskBidTable(cacher)

		provideTodaySeries(cacher)

	}

	return false
}

func UpdateIO(cacher *cacher.Cacher) bool {

	err := gatherIO()

	if err == nil {

		provideIODetails(cacher)
	}

	return false
}

func UpdatePeriodicAverages(cacher *cacher.Cacher) bool {

	err := gatherPeriodicAverages()

	if err == nil {

		providePeriodicAverages(cacher)
	}

	return false
}

func providePriceDetails(cacher *cacher.Cacher) {

	contentParts := strings.Split(currentPriceContent, "@")

	records := strings.Split(contentParts[2], ";")

	stockPriceList = make(map[string]*StockPrice)
	stockTseTOSymbolList = make(map[string]string)

	var colsNumber = map[string]int{
		"Symbol":           2,
		"CompanyName":      3,
		"TSE_Code":         0,
		"IR_Code":          0,
		"ClosePrice":       6,
		"LastPrice":        7,
		"MaxPrice":         12,
		"MinPrice":         11,
		"RangeTopPrice":    19,
		"RangeBottomPrice": 20,
		"YesterdayPrice":   13,
		"TodayFirstPrice":  5,
		//"BuyPriceAtFirst":   999,
		//"SellPriceAtFirst":  999,
		"Volume":     9,
		"BaseVolume": 15,
		//"BuyVolumeAtFirst":  999,
		//"SellVolumeAtFirst": 999,
		"FullAmount": 10,
		"Quantity":   8,
		"Shares":     21,
	}

	timeStr := timeh.TimeFormat(time.Now(), "Y-m-d H:i:s")
	containDegit := regexp.MustCompile(`\d`)

	for _, row := range records {
		row := strings.Split(row+",,,,,,,,,,,,,,,,,,,,", ",")

		// reject all symbols that contain number
		if containDegit.MatchString(row[colsNumber["Symbol"]]) {
			continue
		}

		// reject incomplete data
		if len(row) < 22 {
			continue
		}

		row[colsNumber["Symbol"]] = stringsh.TextPersianize(row[colsNumber["Symbol"]])

		sr := &StockPrice{
			Symbol:           row[colsNumber["Symbol"]],
			CompanyName:      row[colsNumber["CompanyName"]],
			TSE_Code:         row[colsNumber["TSE_Code"]],
			IR_Code:          row[colsNumber["IR_Code"]],
			Time:             timeStr,
			ClosePrice:       toInt(row[colsNumber["ClosePrice"]]),
			LastPrice:        toInt(row[colsNumber["LastPrice"]]),
			MaxPrice:         toInt(row[colsNumber["MaxPrice"]]),
			MinPrice:         toInt(row[colsNumber["MinPrice"]]),
			RangeTopPrice:    toInt(row[colsNumber["RangeTopPrice"]]),
			RangeBottomPrice: toInt(row[colsNumber["RangeBottomPrice"]]),
			YesterdayPrice:   toInt(row[colsNumber["YesterdayPrice"]]),
			TodayFirstPrice:  toInt(row[colsNumber["TodayFirstPrice"]]),
			//BuyPriceAtFirst:   toInt(row[colsNumber["BuyPriceAtFirst"]]),
			//SellPriceAtFirst:  toInt(row[colsNumber["SellPriceAtFirst"]]),
			Volume:     toInt(row[colsNumber["Volume"]]),
			BaseVolume: toInt(row[colsNumber["BaseVolume"]]),
			//BuyVolumeAtFirst:  toInt(row[colsNumber["BuyVolumeAtFirst"]]),
			//SellVolumeAtFirst: toInt(row[colsNumber["SellVolumeAtFirst"]]),
			FullAmount: toInt(row[colsNumber["FullAmount"]]),
			Quantity:   toInt(row[colsNumber["Quantity"]]),
			Shares:     toInt(row[colsNumber["Shares"]]),
		}

		// calcul price rates
		sr.LastPriceRate = float64(sr.LastPrice-sr.YesterdayPrice) * 100 / float64(sr.YesterdayPrice)
		sr.ClosePriceRate = float64(sr.ClosePrice-sr.YesterdayPrice) * 100 / float64(sr.YesterdayPrice)

		stockPriceList[sr.Symbol] = sr

		stockTseTOSymbolList[sr.TSE_Code] = sr.Symbol

	}

	if len(stockPriceList) > 200 {
		cacher.Put("stockPriceList", stockPriceList, 24*60*60)

	}
}

func provideAskBidTable(cacher *cacher.Cacher) {

	contentParts := strings.Split(currentPriceContent, "@")

	if len(contentParts) <= 3 {
		return
	}

	tableRowsData := strings.Split(contentParts[3], ";")

	stockAskBidTableList = make(map[string]*StockAskBidTable)

	for _, rowData := range tableRowsData {
		rowData := strings.Split(rowData, ",")

		// reject unnecessary rows
		tseCode := rowData[0]
		symbol, isExist := stockTseTOSymbolList[tseCode]
		if isExist == false {
			continue
		}

		if len(rowData) == 8 {

			priceAskBidTableRow := &StockAskBidTable{
				BuyRows:  make(map[int64]AskBidTableRow),
				SellRows: make(map[int64]AskBidTableRow),
			}

			//
			if stockAskBidTableList[symbol] != nil {
				priceAskBidTableRow = stockAskBidTableList[symbol]
			}

			// buy row
			if toInt(rowData[4]) > 0 {
				buyAskBidTableRow := AskBidTableRow{
					Volume:   toInt(rowData[6]),
					Quantity: int(toInt(rowData[3])),
				}

				priceAskBidTableRow.BuyRows[(toInt(rowData[4]))] = buyAskBidTableRow
			}

			// sell row
			if toInt(rowData[5]) > 0 {
				sellAskBidTableRow := AskBidTableRow{
					Volume:   toInt(rowData[7]),
					Quantity: int(toInt(rowData[2])),
				}

				priceAskBidTableRow.SellRows[toInt(rowData[5])] = sellAskBidTableRow
			}

			stockAskBidTableList[symbol] = priceAskBidTableRow

		}
	}

	cacher.Put("stockAskBidTableList", stockAskBidTableList, 24*60*60)
}

func provideTodaySeries(cacher *cacher.Cacher) {

	if len(stockIOList) == 0 {
		UpdateIO(cacher)
	}

	clocknumber := GetClockNumber()

	for _, stock := range stockPriceList {

		// Today Stock Charts
		//
		ts, ok := stockTodaySeriesList[stock.Symbol]
		if ok == false {
			ts = &StockTodaySeries{
				LastPrice:       make(LineInt, 0),
				BuyVolume:       make(LineInt, 0),
				SellVolume:      make(LineInt, 0),
				IndiBuySaraneh:  make(LineFloat, 0),
				IndiSellSaraneh: make(LineFloat, 0),
			}
		}

		ts.LastPrice[clocknumber] = stock.LastPrice

		// series of buy queue details
		buyr, ok1 := stockAskBidTableList[stock.Symbol]
		if ok1 == true {
			ts.BuyVolume[clocknumber] = 0
			buyr, ok1 := buyr.BuyRows[stock.RangeTopPrice]
			if ok1 == true {
				ts.BuyVolume[clocknumber] = buyr.Volume
			}
		}

		// series of sell queue details
		sellr, ok2 := stockAskBidTableList[stock.Symbol]
		if ok2 == true {
			ts.SellVolume[clocknumber] = 0
			sellr, ok2 := sellr.SellRows[stock.RangeBottomPrice]
			if ok2 == true {
				ts.SellVolume[clocknumber] = sellr.Volume
			}
		}

		// series of IO details
		ior, ok3 := stockIOList[stock.Symbol]
		if ok3 == true {
			ts.IndiBuySaraneh[clocknumber] = ior.IndiBuySaraneh
			ts.IndiSellSaraneh[clocknumber] = ior.IndiSellSaraneh
		}

		stockTodaySeriesList[stock.Symbol] = ts

	}

	cacher.Put("stockTodaySeriesList", stockTodaySeriesList, 24*60*60)
}

func provideIODetails(cacher *cacher.Cacher) {

	//contentParts := strings.Split(currentIOContent, "@")

	records := strings.Split(currentIOContent, ";")

	stockIOList = make(map[string]*StockIO)

	var colsNumber = map[string]int{
		"Symbol":           0,
		"TSE_Code":         0,
		"IndiBuyQuantity":  1,
		"OrgaBuyQuantity":  2,
		"IndiBuyVolume":    3,
		"OrgaBuyVolume":    4,
		"IndiSellQuantity": 5,
		"OrgaSellQuantity": 6,
		"IndiSellVolume":   7,
		"OrgaSellVolume":   8,
	}

	for _, row := range records {
		row := strings.Split(row+",,,,,,,", ",")

		sir := &StockIO{
			TSE_Code:         row[colsNumber["TSE_Code"]],
			IndiBuyQuantity:  toInt(row[colsNumber["IndiBuyQuantity"]]),
			OrgaBuyQuantity:  toInt(row[colsNumber["OrgaBuyQuantity"]]),
			IndiBuyVolume:    toInt(row[colsNumber["IndiBuyVolume"]]),
			OrgaBuyVolume:    toInt(row[colsNumber["OrgaBuyVolume"]]),
			IndiSellQuantity: toInt(row[colsNumber["IndiSellQuantity"]]),
			OrgaSellQuantity: toInt(row[colsNumber["OrgaSellQuantity"]]),
			IndiSellVolume:   toInt(row[colsNumber["IndiSellVolume"]]),
			OrgaSellVolume:   toInt(row[colsNumber["OrgaSellVolume"]]),
		}

		// reject unnecessary rows
		sym, isExist := stockTseTOSymbolList[sir.TSE_Code]
		if isExist == false {
			continue
		}

		sir.Symbol = sym

		spr := stockPriceList[sir.Symbol]

		sir.Time = spr.Time
		sir.LastPrice = spr.LastPrice
		sir.LastPriceRate = spr.LastPriceRate

		if sir.IndiBuySaraneh = 0; sir.IndiBuyQuantity > 0 {
			sir.IndiBuySaraneh = math.Round(float64(sir.IndiBuyVolume)*float64(sir.LastPrice)/(float64(sir.IndiBuyQuantity)*10000000)*100) / 100
		}
		if sir.IndiSellSaraneh = 0; sir.IndiSellQuantity > 0 {
			sir.IndiSellSaraneh = math.Round(float64(sir.IndiSellVolume)*float64(sir.LastPrice)/(float64(sir.IndiSellQuantity)*10000000)*100) / 100
		}
		if sir.IndiBuyPower = 0; sir.IndiSellSaraneh > 0 {
			sir.IndiBuyPower = math.Round(float64(sir.IndiBuySaraneh)/(float64(sir.IndiSellSaraneh))*10) / 10
		}
		if sir.IndiBuyPercent = 0; sir.IndiBuyVolume+sir.OrgaBuyVolume > 0 {
			sir.IndiBuyPercent = math.Round(float64(sir.IndiBuyVolume*100) / (float64(sir.IndiBuyVolume + sir.OrgaBuyVolume)))
		}
		if sir.OrgaSellPercent = 0; sir.OrgaSellVolume+sir.IndiSellVolume > 0 {
			sir.OrgaSellPercent = math.Round(float64(sir.OrgaSellVolume*100) / (float64(sir.OrgaSellVolume + sir.IndiSellVolume)))
		}

		stockIOList[sir.Symbol] = sir

	}

	cacher.Put("stockIOList", stockIOList, 24*60*60)
}

func providePeriodicAverages(cacher *cacher.Cacher) {

	//contentParts := strings.Split(currentIOContent, "@")

	records := strings.Split(currentPeriodicAveragesContent, ";")

	stockPeriodicAveragesList = make(map[string]*StockAverages)

	var colsNumber = map[string]int{
		"TSE_Code":                0,
		"Last3MonthAverageAmount": 1,
		"Last3MonthAverageVolume": 5,
	}

	for _, row := range records {
		row := strings.Split(row, ",")

		sar := &StockAverages{
			TSE_Code:                row[colsNumber["TSE_Code"]],
			Last3MonthAverageAmount: toInt(row[colsNumber["Last3MonthAverageAmount"]]),
			Last3MonthAverageVolume: toInt(row[colsNumber["Last3MonthAverageVolume"]]),
		}

		// reject unnecessary rows
		_, isExist := stockTseTOSymbolList[sar.TSE_Code]
		if isExist == false {
			continue
		}

		sar.Symbol = stockTseTOSymbolList[sar.TSE_Code]

		stockPeriodicAveragesList[sar.Symbol] = sar

	}

	cacher.Put("stockPeriodicAveragesList", stockPeriodicAveragesList, 24*60*60)
}

func storeTodayLines() {

	// todayStock, ok := todayStocks[sr.Symbol]
	// if ok != true {
	// 	todayStock = TodayStock{}
	// }
	// todayStock.LastPrice[marketTime] = sr.LastPrice
	// todayStock.BuyVolume[marketTime] = sr.BuyVolumeAtFirst
	// todayStock.SellVolume[marketTime] = sr.BaseVolume
	// todayStock.IndiBySaraneh[marketTime] = sr.BaseVolume
}

func ToString() string {
	//s := "length of stockPriceList: " + strconv.Itoa(len(stockPriceList)) + "\r\n"

	//stock := stockTodaySeriesList["????????"]
	//x := fmt.Sprint(stock)
	//s = s + "Series1: " + arrayToString(stock.LastPrice, ",") + "\r\n"

	return fmt.Sprint(stockTodaySeriesList)
}

func toInt(str string) int64 {
	if str == "" {
		return 0
	}
	return stringsh.ToInt(str)
}

func arrayToString(a LineInt, delimiter string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(a), " ", delimiter, -1), "[]")
}
