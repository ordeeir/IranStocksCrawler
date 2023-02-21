package bourse

import (
	"IranStocksCrawler/helpers/osh"
	"IranStocksCrawler/helpers/timeh"
	"errors"
	"os"
	"time"
)

var gatheringPrices bool = false
var gatheringIO bool = false
var gatheringPeriodicAverages bool = false

func gatherPrices() error {

	if gatheringPrices == true {
		return errors.New("gathering engine is busy")
	}

	// dont check first time
	if len(stockPriceList) > 0 {
		clocknumber := GetClockNumber()
		if clocknumber > 12*60*60 {
			return errors.New("Error in market time is not set correctly")
		}
	}

	gatheringPrices = true

	defer func() {
		gatheringPrices = false
	}()

	if consecutiveAtempts["FailedUpdatePrices"] > 10 {
		//log.Fatalf("Error in getting market data")
		return errors.New("Error in getting market data (over atempt)")
	}

	st := GetMarketStatus()

	if st == DEF_MARKET_STATUS_OPEN || len(stockPriceList) == 0 {

		content, errFetch := Fetch(DEF_URLS_PRICE_URL, DEF_PATHS_PRICE_PATH, time.Second*10)
		if errFetch != nil {
			consecutiveAtempts["FailedUpdatePrices"]++
			return errors.New("Error in getting market data (fetch)")
		}

		b := []byte(content)
		os.WriteFile(osh.GetRootPath()+"/files/pricedata"+timeh.TimeFormat(time.Now(), "Y-m-d-H-i-s")+".txt", b, 0777)
		currentPriceContent = content

		return nil
	}

	return errors.New("Market is closed")
}

func gatherIO() error {

	if gatheringIO == true {
		return errors.New("gathering engine is busy")
	}
	gatheringIO = true

	defer func() {
		gatheringIO = false
	}()

	if consecutiveAtempts["FailedUpdateIO"] > 10 {
		//log.Fatalf("Error in getting market data")
		return errors.New("Error in getting market io data (over atempt)")
	}

	if marketStatus == DEF_MARKET_STATUS_OPEN || len(stockIOList) == 0 {

		content, errFetch := Fetch(DEF_URLS_IO_URL, DEF_PATHS_IO_PATH, time.Second*10)
		if errFetch != nil {
			consecutiveAtempts["FailedUpdateIO"]++
			return errors.New("Error in getting market io data (fetch)")
		}

		b := []byte(content)
		os.WriteFile(osh.GetRootPath()+"/files/iodata"+timeh.TimeFormat(time.Now(), "Y-m-d-H-i-s")+".txt", b, 0777)

		currentIOContent = content

		return nil
	}

	return errors.New("Market is closed")
}

func gatherPeriodicAverages() error {

	if gatheringPeriodicAverages == true {
		return errors.New("gathering engine is busy")
	}
	gatheringPeriodicAverages = true

	defer func() {
		gatheringPeriodicAverages = false
	}()

	if consecutiveAtempts["FailedUpdatePeriodicAverages"] > 10 {
		//log.Fatalf("Error in getting market data")
		return errors.New("Error in getting market periodic avg data (over atempt)")
	}

	content, errFetch := Fetch(DEF_URLS_PERIODIC_AVERAGES_URL, DEF_PATHS_PERIODIC_AVERAGES_PATH, time.Hour*6)
	if errFetch != nil {
		consecutiveAtempts["FailedUpdatePeriodicAverages"]++
		return errors.New("Error in getting market periodic avg data (fetch)")
	}

	currentPeriodicAveragesContent = content

	return nil

}
