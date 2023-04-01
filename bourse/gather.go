package bourse

import (
	"encoding/base64"
	"errors"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

var gatheringPrices bool = false
var gatheringIO bool = false
var gatheringPeriodicAverages bool = false
var gatheringIndiOrga365Days bool = false

var LastSuccessGatheringPrices time.Time

func gatherPrices() error {

	if gatheringPrices == true {
		return errors.New("gathering engine is busy")
	}
	logrus.Debugf("gatherPrices tick. stockPriceList has %v symbols", len(stockPriceList))

	gatheringPrices = true

	defer func() {
		gatheringPrices = false
	}()

	content, errFetch := Fetch(DEF_URLS_PRICE_URL, DEF_PATHS_PRICE_PATH, time.Second*10)
	if errFetch != nil {

		consecutiveAtempts["FailedUpdatePrices"]++
		return errors.New("Error in getting market data (fetch)")

	}

	//b := []byte(content)
	//os.WriteFile(osh.GetRootPath()+"/files/pricedata"+timeh.TimeFormat(time.Now(), "Y-m-d-H-i-s")+".txt", b, 0777)

	// check valid str
	contentParts := strings.Split(content, "@")

	if len(contentParts) < 2 {

		logrus.Warn("Url fetched but content of prices is invalid")

		return errors.New("Error in getting market data (fetch)")
	}

	logrus.Debug("Gathering Prices done")

	currentPriceContent = content

	return nil
}

func gatherIO() error {

	if gatheringIO == true {
		return errors.New("gathering engine is busy")
	}
	gatheringIO = true

	defer func() {
		gatheringIO = false
	}()

	logrus.Debugf("gatherIO tick. stockIOList has %v symbols", len(stockIOList))

	content, errFetch := Fetch(DEF_URLS_IO_URL, DEF_PATHS_IO_PATH, time.Second*10)
	if errFetch != nil {
		consecutiveAtempts["FailedUpdateIO"]++
		return errors.New("Error in getting market io data (fetch)")
	}

	records := strings.Split(content, ";")

	if len(records) < 10 {
		logrus.Warn("Url fetched but content of IO is invalid")
	}

	logrus.Debug("Gathering IO done")

	currentIOContent = content

	return nil

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

func gatherIndiOrga365Days(tseCode string) (string, error) {

	if gatheringIndiOrga365Days == true {
		return "", errors.New("gathering engine is busy")
	}
	gatheringIndiOrga365Days = true

	defer func() {
		gatheringIndiOrga365Days = false
	}()

	//if consecutiveAtempts["FailedUpdatePeriodicAverages"] > 10 {
	//log.Fatalf("Error in getting market data")
	//return "",errors.New("Error in getting market periodic avg data (over atempt)")
	//}

	url := strings.ReplaceAll(DEF_URLS_INDIORGA_DAYS_DATA_URL, "{TSE_CODE}", tseCode)

	agent := settings["curl-agent"]
	base64Url := base64.StdEncoding.EncodeToString([]byte(url))
	url = strings.ReplaceAll(agent, "{BASE64_URL}", base64Url)

	content, errFetch := Fetch(url, DEF_PATHS_PERIODIC_AVERAGES_PATH, 0)
	if errFetch != nil {
		consecutiveAtempts["FailedUpdateIndiOrga365Days"]++
		return "", errors.New("Error in getting indi orga 365 days data (fetch)")
	}

	return content, nil

}
