package bourse

import (
	"IranStocksCrawler/helpers/stringsh"
	"IranStocksCrawler/helpers/timeh"
	"IranStocksCrawler/system/cacher"
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
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

var stockIndiOrga365DaysList map[string]*StockIndiOrga365Days = map[string]*StockIndiOrga365Days{}
var stockTseTOSymbolList map[string]string

var lastClockNumber int64 = -1

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

type StockIndiOrga struct {
	QuantityIndiBuy  string
	QuantityOrgaBuy  string
	QuantityIndiSell string
	QuantityOrgaSell string
	VolumeIndiBuy    string
	VolumeOrgaBuy    string
	VolumeIndiSell   string
	VolumeOrgaSell   string
	AmountIndiBuy    string
	AmountOrgaBuy    string
	AmountIndiSell   string
	AmountOrgaSell   string
}

type StockIndiOrga365Days struct {
	Days       map[string]StockIndiOrga
	LastUpdate string
}

func UpdatePrices(cacher *cacher.Cacher) bool {

	var err error

	err = errors.New("")

	if IsOpenTime() == false {
		SetMarketClose()
	}

	if IsOpenTime() && MarketBeenOpenToday() {

		marketIsOpen := MarketIsOpen()

		if marketIsOpen == false {
			ResetGatheredData(cacher)
			SetMarketOpen()
		}

		err = gatherPrices()

	}

	if IsGeneralInfoEmpty() {
		err = gatherPrices()
	}

	if err == nil {

		providePriceDetails(cacher)

		provideAskBidTable(cacher)

		provideTodaySeries(cacher)

	}

	return false
}

func UpdateIO(cacher *cacher.Cacher) bool {

	err := errors.New("")

	if IsOpenTime() && MarketBeenOpenToday() {
		err = gatherIO()

	}

	if err == nil {

		provideIODetails(cacher)
	}

	return false
}

func UpdatePeriodicAverages(cacher *cacher.Cacher) bool {

	err := errors.New("")

	if IsPeriodicAveragesEmpty() {
		err = gatherPeriodicAverages()
	}

	if err == nil {

		providePeriodicAverages(cacher)
	}

	return false
}

func GetIndiOrga365DaysFromFilesystem(sym string) (string, string, error) {

	item := stockPriceList[sym]

	path := strings.ReplaceAll(DEF_PATHS_INDIORGA_DAYS_DATA_PATH, "{TSE_CODE}", item.TSE_Code)

	content, err := ioutil.ReadFile(path)

	if err != nil {
		return "", "", err
	}

	fi, err2 := os.Stat(path)

	if err2 != nil {
		return "", "", err2
	}

	mtime := fi.ModTime()

	tf := mtime.Format("01-02")

	return string(content), tf, nil
}

func UpdateIndiOrga365Days(cacher *cacher.Cacher) bool {

	if nowHM() <= "06:00" {

		logrus.Debugf("nowHM: %v , UpdateIndiOrga365Days rejected", nowHM())

		return false
	}

	logrus.Debugf("nowHM: %v , UpdateIndiOrga365Days is gathering...", nowHM())

	if IsGeneralInfoEmpty() == false {

		if len(stockIndiOrga365DaysList) == 0 {

			for sym := range stockPriceList {

				content, time, err := GetIndiOrga365DaysFromFilesystem(sym)

				if err != nil {
					continue
				}

				logrus.Debugf("GetIndiOrga365DaysFromFilesystem reterned. symbol: %v , content length: %v ", sym, len(content))

				symDays := &StockIndiOrga365Days{}

				if len(content) > 0 {
					symDays = generateIndiOrga365FromContent(sym, content)
				} else {
					symDays.Days = map[string]StockIndiOrga{}
				}

				symDays.LastUpdate = time

				stockIndiOrga365DaysList[sym] = symDays
			}
		}

		i := 0

		for sym, item := range stockPriceList {

			i++

			symDays, ok := stockIndiOrga365DaysList[sym]

			if ok {
				if symDays.LastUpdate == todayMD() {
					continue
				}
				// symDays = StockIndiOrga365Days{
				// 	LastUpdate: "00-00",
				// 	Days:       map[string]StockIndiOrga{},
				// }
			}

			content, err := gatherIndiOrga365Days(item.TSE_Code)

			if err != nil {
				continue
			}

			symDays = generateIndiOrga365FromContent(sym, content)

			stockIndiOrga365DaysList[sym] = symDays

			if i >= 1 {

				logrus.Debugf("trying to save stockIndiOrga365DaysList(size %d)...", getRealSizeOf(stockIndiOrga365DaysList))

				logrus.Debugf("stockIndiOrga365DaysList with %v symbols gathered and stored into filesystem", len(stockIndiOrga365DaysList))

				return true
			}
		}
	}

	return false
}

func generateIndiOrga365FromContent(sym string, content string) *StockIndiOrga365Days {

	symDays := &StockIndiOrga365Days{
		LastUpdate: "00-00",
		Days:       map[string]StockIndiOrga{},
	}

	rows := strings.Split(content, ";")

	row := []string{}

	j := 0

	for _, jrow := range rows {

		j++
		if j > 310 {
			break
		}

		row = strings.Split(jrow, ",")

		sio := StockIndiOrga{}
		sio.QuantityIndiBuy = row[1]
		sio.QuantityOrgaBuy = row[2]
		sio.QuantityIndiSell = row[3]
		sio.QuantityOrgaSell = row[4]
		sio.VolumeIndiBuy = row[5]
		sio.VolumeOrgaBuy = row[6]
		sio.VolumeIndiSell = row[7]
		sio.VolumeOrgaSell = row[8]
		sio.AmountIndiBuy = row[9]
		sio.AmountOrgaBuy = row[10]
		sio.AmountIndiSell = row[11]
		sio.AmountOrgaSell = row[12]

		symDays.Days[row[0]] = sio
	}

	symDays.LastUpdate = todayMD()

	return symDays
}

/*
func UpdateIndiOrga365Days(cacher *cacher.Cacher) bool {

	if nowHM() <= "06:00" {

		logrus.Debugf("nowHM: %v , UpdateIndiOrga365Days rejected", nowHM())

		return false
	}

	logrus.Debugf("nowHM: %v , UpdateIndiOrga365Days is gathering...", nowHM())

	if IsGeneralInfoEmpty() == false {

		if len(stockIndiOrga365DaysList) == 0 {

			json := cacher.Get("stockIndiOrga365DaysList")

			if json != nil {
				stockIndiOrga365DaysList = ConvertMapOfInterfaceToStockIndiOrga365Days(json)
			}
		}

		i := 0

		for sym, item := range stockPriceList {

			i++

			content, err := gatherIndiOrga365Days(item.TSE_Code)

			symDays, ok := stockIndiOrga365DaysList[sym]
			if !ok {
				symDays = StockIndiOrga365Days{
					LastUpdate: "00-00",
					Days:       map[string]StockIndiOrga{},
				}
			}

			if symDays.LastUpdate == todayMD() {
				continue
			}

			if err == nil {

				rows := strings.Split(content, ";")

				row := []string{}

				j := 0

				for _, jrow := range rows {

					j++
					if j > 200 {
						break
					}

					row = strings.Split(jrow, ",")

					sio := StockIndiOrga{}
					sio.QuantityIndiBuy = row[1]
					sio.QuantityOrgaBuy = row[2]
					sio.QuantityIndiSell = row[3]
					sio.QuantityOrgaSell = row[4]
					sio.VolumeIndiBuy = row[5]
					sio.VolumeOrgaBuy = row[6]
					sio.VolumeIndiSell = row[7]
					sio.VolumeOrgaSell = row[8]
					sio.AmountIndiBuy = row[9]
					sio.AmountOrgaBuy = row[10]
					sio.AmountIndiSell = row[11]
					sio.AmountOrgaSell = row[12]

					symDays.Days[row[0]] = sio
				}

				symDays.LastUpdate = todayMD()

			}

			stockIndiOrga365DaysList[sym] = symDays

			if i >= 1 {

				logrus.Debugf("trying to save stockIndiOrga365DaysList(size %d)...", getRealSizeOf(stockIndiOrga365DaysList))

				err := cacher.Put("stockIndiOrga365DaysList", stockIndiOrga365DaysList, 30*24*60*60)

				if err != nil {
					logrus.Error(err)
				} else {
					logrus.Debugf("stockIndiOrga365DaysList with %v symbols gathered and stored into redis", len(stockIndiOrga365DaysList))
				}

				return true
			}
		}
	}

	return false
}
*/

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

		logrus.Debugf("trying to save stockPriceList(size %d)...", getRealSizeOf(stockPriceList))

		err := cacher.Put("stockPriceList", stockPriceList, 30*24*60*60)

		if err != nil {
			logrus.Error(err)
		} else {
			logrus.Debugf("stockPriceList with %v symbols gathered and stored into redis", len(stockPriceList))
			logrus.Debugf("stockTseTOSymbolList with %v symbols stored into redis", len(stockTseTOSymbolList))
		}

		err2 := cacher.Put("lastTimeStorage", todayYMDHM(), 30*24*60*60)

		if err2 != nil {
			logrus.Error(err2)
		} else {
			logrus.Debugf("lastTimeStorage and stored into redis")
		}

	}
}

func provideAskBidTable(cacher *cacher.Cacher) {

	contentParts := strings.Split(currentPriceContent, "@")

	if len(contentParts) <= 3 {
		return
	}

	tableRowsData := strings.Split(contentParts[3], ";")

	stockAskBidTableList = map[string]*StockAskBidTable{}

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

	logrus.Debugf("trying to save stockAskBidTableList(size %d)...", getRealSizeOf(stockAskBidTableList))

	err := cacher.Put("stockAskBidTableList", stockAskBidTableList, 30*24*60*60)

	if err != nil {
		logrus.Error(err)
	} else {
		logrus.Debugf("stockAskBidTableList with %v symbols gathered and stored into redis", len(stockAskBidTableList))
	}

}

func ConvertMapOfInterfaceToStockTodaySeries(interf interface{}) map[string]*StockTodaySeries {

	result := map[string]*StockTodaySeries{}

	mainMap := interf.(map[string]interface{})

	for symbol, inter1 := range mainMap {

		ii := inter1.(map[string]interface{})

		bv := ii["BuyVolume"].(map[string]interface{})
		sv := ii["SellVolume"].(map[string]interface{})
		lp := ii["LastPrice"].(map[string]interface{})
		ibs := ii["IndiBuySaraneh"].(map[string]interface{})
		iss := ii["IndiSellSaraneh"].(map[string]interface{})

		buyVolume := LineInt{}
		sellVolume := LineInt{}
		lastPrice := LineInt{}
		indiBuySaraneh := LineFloat{}
		indiSellSaraneh := LineFloat{}

		for _key, inter2 := range bv {
			i, _ := strconv.ParseInt(_key, 10, 64)
			buyVolume[i] = int64(inter2.(float64))
		}
		for _key, inter2 := range sv {
			i, _ := strconv.ParseInt(_key, 10, 64)
			sellVolume[i] = int64(inter2.(float64))
		}
		for _key, inter2 := range lp {
			i, _ := strconv.ParseInt(_key, 10, 64)
			lastPrice[i] = int64(inter2.(float64))
		}
		for _key, inter2 := range ibs {
			i, _ := strconv.ParseInt(_key, 10, 64)
			indiBuySaraneh[i] = inter2.(float64)
		}
		for _key, inter2 := range iss {
			i, _ := strconv.ParseInt(_key, 10, 64)
			indiSellSaraneh[i] = inter2.(float64)
		}

		result[symbol] = &StockTodaySeries{
			BuyVolume:       buyVolume,
			SellVolume:      sellVolume,
			LastPrice:       lastPrice,
			IndiBuySaraneh:  indiBuySaraneh,
			IndiSellSaraneh: indiSellSaraneh,
		}

	}

	return result
}

func ConvertMapOfInterfaceToStockIndiOrga365Days(interf interface{}) map[string]StockIndiOrga365Days {

	result := map[string]StockIndiOrga365Days{}

	mainMap := interf.(map[string]interface{})

	for symbol, inter1 := range mainMap {

		ii := inter1.(map[string]interface{})

		_days := ii["Days"].(map[string]interface{})

		lastUpdate := ii["LastUpdate"].(string)

		days := map[string]StockIndiOrga{}

		for _key, inter2 := range _days {
			row := inter2.(map[string]interface{})

			sio := StockIndiOrga{}
			sio.AmountIndiBuy = row["AmountIndiBuy"].(string)
			sio.AmountIndiSell = row["AmountIndiSell"].(string)
			sio.AmountOrgaBuy = row["AmountOrgaBuy"].(string)
			sio.AmountOrgaSell = row["AmountOrgaSell"].(string)
			sio.QuantityIndiBuy = row["QuantityIndiBuy"].(string)
			sio.QuantityIndiSell = row["QuantityIndiSell"].(string)
			sio.QuantityOrgaBuy = row["QuantityOrgaBuy"].(string)
			sio.QuantityOrgaSell = row["QuantityOrgaSell"].(string)
			sio.VolumeIndiBuy = row["VolumeIndiBuy"].(string)
			sio.VolumeIndiSell = row["VolumeIndiSell"].(string)
			sio.VolumeOrgaBuy = row["VolumeOrgaBuy"].(string)
			sio.VolumeOrgaSell = row["VolumeOrgaSell"].(string)

			days[_key] = sio
		}

		result[symbol] = StockIndiOrga365Days{
			Days:       days,
			LastUpdate: lastUpdate,
		}

	}

	return result
}

func nowHM() string {

	tz := time.Now()

	timeZone, _ := time.LoadLocation("Asia/Tehran")

	t := tz.In(timeZone)

	return t.Format("15:04")
}

func todayMD() string {
	t := time.Now().Truncate(3 * time.Hour)
	today := fmt.Sprintf("%02d-%02d", t.Month(), t.Day())
	return today
}

func todayYMDHM() string {
	t := time.Now()
	ti := t.Format("2006-01-02 15:04:05")
	return ti
}

func provideTodaySeries(cacher *cacher.Cacher) {

	if len(stockIOList) == 0 {
		UpdateIO(cacher)
	}

	//today := todayMD()

	clocknumber := GetClockNumber()

	// last clock number
	if lastClockNumber == -1 {

		_interface := cacher.Get("lastClockNumber")

		if _interface != nil {
			lastClockNumber = int64(_interface.(float64))
		} else {
			lastClockNumber = 4 * 60 * 60
		}

	}

	if clocknumber < lastClockNumber {
		stockTodaySeriesList = map[string]*StockTodaySeries{}
	}

	if len(stockTodaySeriesList) == 0 {

		// today series
		_interface := cacher.Get("stockTodaySeriesList")

		if _interface != nil {

			stockTodaySeriesList = ConvertMapOfInterfaceToStockTodaySeries(_interface)

			logrus.Debugf("stockTodaySeriesList loaded from redis (last clocknumber = %v)", lastClockNumber)
		}
	}

	//
	needToStore := false

	// iterate
	for _, stock := range stockPriceList {

		// Today Stock Charts
		//
		ts, ok := stockTodaySeriesList[stock.Symbol]
		if ok == false {
			ts = &StockTodaySeries{
				LastPrice:       LineInt{},
				BuyVolume:       LineInt{},
				SellVolume:      LineInt{},
				IndiBuySaraneh:  LineFloat{},
				IndiSellSaraneh: LineFloat{},
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
				needToStore = true
			}
		}

		// series of sell queue details
		sellr, ok2 := stockAskBidTableList[stock.Symbol]
		if ok2 == true {
			ts.SellVolume[clocknumber] = 0
			sellr, ok2 := sellr.SellRows[stock.RangeBottomPrice]
			if ok2 == true {
				ts.SellVolume[clocknumber] = sellr.Volume
				needToStore = true
			}
		}

		// series of IO details
		ior, ok3 := stockIOList[stock.Symbol]
		if ok3 == true {
			ts.IndiBuySaraneh[clocknumber] = ior.IndiBuySaraneh
			ts.IndiSellSaraneh[clocknumber] = ior.IndiSellSaraneh
			needToStore = true
		}

		stockTodaySeriesList[stock.Symbol] = ts

	}

	if needToStore {

		lastClockNumber = clocknumber

		logrus.Debugf("One row added to stockTodaySeriesList (last clocknumber = %v)", lastClockNumber)

		logrus.Debugf("trying to save stockTodaySeriesList(size %d)...", getRealSizeOf(stockTodaySeriesList))

		err := cacher.Put("stockTodaySeriesList", stockTodaySeriesList, 30*24*60*60)

		if err != nil {
			logrus.Error(err)
		} else {
			logrus.Debugf("stockTodaySeriesList with %v symbols stored to stockTodaySeriesList in redis (last clocknumber = %v)", len(stockTodaySeriesList), lastClockNumber)
		}

		cacher.Put("lastClockNumber", lastClockNumber, 30*24*60*60)

	}

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
		row := strings.Split(row+",,,,,,,,,", ",")

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

	logrus.Debugf("trying to save stockIOList(size %d)...", getRealSizeOf(stockIOList))

	err := cacher.Put("stockIOList", stockIOList, 30*24*60*60)

	if err != nil {
		logrus.Error(err)
	} else {
		logrus.Debugf("stockIOList with %v symbols stored into redis", len(stockIOList))
	}

}

func providePeriodicAverages(cacher *cacher.Cacher) {

	records := strings.Split(currentPeriodicAveragesContent, ";")

	logrus.Debugf("currentPeriodicAveragesContent has %v length, records have %v rows", len(currentPeriodicAveragesContent), len(records))

	stockPeriodicAveragesList = map[string]*StockAverages{}

	// var colsNumber = map[string]string{
	// 	"Last3MonthAverageAmount": "1",
	// 	"Last3MonthAverageVolume": "5",
	// }

	//i := 0

	for _, _row := range records {
		//logrus.Debug("11111: " + _row)

		// i++
		// if i > 5 {
		// 	continue
		// }

		row := strings.Split(_row, ",")

		sar := &StockAverages{}

		if len(row) == 3 {

			sar.TSE_Code = row[0]

			if row[1] == "1" {

				sar.Last3MonthAverageAmount = toInt(row[2])

			} else if row[1] == "5" {

				sar.Last3MonthAverageVolume = toInt(row[2])
			}

		} else if len(row) == 2 {

			if row[0] == "1" {

				sar.Last3MonthAverageAmount = toInt(row[1])

			} else if row[0] == "5" {

				sar.Last3MonthAverageVolume = toInt(row[1])
			}
		}

		//logrus.Debug("11111 " + sar.TSE_Code)

		// reject unnecessary rows
		_, isExist := stockTseTOSymbolList[sar.TSE_Code]
		if isExist == false {
			continue
		}

		//logrus.Debug("22222 " + sar.Symbol)

		sar.Symbol = stockTseTOSymbolList[sar.TSE_Code]

		stockPeriodicAveragesList[sar.Symbol] = sar

	}

	logrus.Debugf("trying to save stockPeriodicAveragesList(size %d)...", getRealSizeOf(stockPeriodicAveragesList))

	err := cacher.Put("stockPeriodicAveragesList", stockPeriodicAveragesList, 30*24*60*60)

	if err != nil {
		logrus.Error(err)
	} else {
		logrus.Debugf("stockPeriodicAveragesList with %v symbols stored into redis", len(stockPeriodicAveragesList))
	}

}

func proviseIndiOrga365Days() {

}

func GetTodaySeries(symbol string) string {

	data := stockTodaySeriesList[symbol]

	return fmt.Sprint(data)
}

func GetIndiOrgaListSymbols() []string {

	var list = []string{}

	for i := range stockIndiOrga365DaysList {
		list = append(list, i)
	}

	return list
}

func GetIndiOrga(sym string) (*StockIndiOrga365Days, error) {

	var data = &StockIndiOrga365Days{}

	for i, j := range stockIndiOrga365DaysList {
		if i == sym {
			data = j
			return data, nil

		}
	}

	return data, errors.New("Not Found")
}

func ToString() string {
	var array map[string]StockTodaySeries = make(map[string]StockTodaySeries)

	i := 0
	for key, item := range stockTodaySeriesList {
		if i > 10 {
			break
		}

		array[key] = *item
	}

	return fmt.Sprint(array)
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

func IsGeneralInfoEmpty() bool {
	if len(stockPriceList) == 0 {
		return true
	}
	return false
}

func IsPeriodicAveragesEmpty() bool {

	if len(stockPeriodicAveragesList) == 0 {

		logrus.Debugf("stockPeriodicAveragesList is not empty, it has %v rows", len(stockPeriodicAveragesList))

		return true
	}

	logrus.Debugf("stockPeriodicAveragesList is empty, it has %v rows", len(stockPeriodicAveragesList))

	return false
}

func ResetGatheredData(cacher *cacher.Cacher) {

	stockTodaySeriesList = map[string]*StockTodaySeries{}
	stockAskBidTableList = map[string]*StockAskBidTable{}
	currentPriceContent = ""
	currentIOContent = ""
	currentPeriodicAveragesContent = ""

	stockPriceList = map[string]*StockPrice{}
	stockIOList = map[string]*StockIO{}
	stockPeriodicAveragesList = map[string]*StockAverages{}
	stockAskBidTableList = map[string]*StockAskBidTable{}
	stockTodaySeriesList = map[string]*StockTodaySeries{}

	stockTseTOSymbolList = map[string]string{}

	cacher.Put("stockPriceList", stockPriceList, 30*24*60*60)
	cacher.Put("stockAskBidTableList", stockAskBidTableList, 30*24*60*60)
	cacher.Put("stockTodaySeriesList", stockTodaySeriesList, 30*24*60*60)
	cacher.Put("stockIOList", stockIOList, 30*24*60*60)
	cacher.Put("stockPeriodicAveragesList", stockPeriodicAveragesList, 30*24*60*60)

	logrus.Debug("storage is reset and stored in redis")

	return
}

func getRealSizeOf(v interface{}) int {
	b := new(bytes.Buffer)
	if err := gob.NewEncoder(b).Encode(v); err != nil {
		return 0
	}
	return b.Len()
}
