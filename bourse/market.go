package bourse

import (
	"IranStocksCrawler/helpers/stringsh"
	"IranStocksCrawler/system/config"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/antchfx/htmlquery"
	"github.com/sirupsen/logrus"
)

const (
	DateTime = "2006-01-02 15:04:05"
)

var marketStatus MarketStatusType = DEF_MARKET_STATUS_UNKNOWN

var calibrationPointDiffTimeH int
var calibrationPointDiffTimeM int

var marketTime string
var marketDate string
var marketStatusLastUpdate time.Time
var marketStatusLastCheck time.Time

var marketIsOpen bool = false
var lastTimeMarketOpen string
var lastTimeMarketUpdate string

var settings config.ConfigList

type MarketStatusType int

func updateMarketDetails() error {

	//if isMarketTime() == false {
	//	setMarketStatus(DEF_MARKET_STATUS_CLOSE)

	//	return nil
	//}

	// if consecutiveAtempts["FailedUpdateMarketDetails"] > 10 {
	// 	return "", errors.New("many efforts without expectable result")
	// }

	//if isMarketUpdated() {
	//	return nil
	//}

	agent := settings["curl-agent1"]
	base64Url := base64.StdEncoding.EncodeToString([]byte(DEF_URLS_MARKET_STATUS_URL))
	url := strings.ReplaceAll(agent, "{BASE64_URL}", base64Url)

	content, errFetch := Fetch(url, DEF_PATHS_MARKET_STATUS_PATH, 0)

	if content == "" && errFetch != nil {
		// byteData, errFetch2 := os.ReadFile("C:/Go_Projects/IranStocksCrawler/files/pricedata2022-10-01-11-36-57.txt")
		// content2 := string(byteData)
		// //content2, errFetch2 := Fetch(DEF_URLS_PRICE_URL, DEF_PATHS_MARKET_STATUS_PATH, time.Second*5)
		// contentParts := strings.Split(content2, "@")
		// //records := strings.Split(contentParts[2], ";")
		// if contentParts != nil {
		// 	//c := strings.Split(contentParts[1], ",")
		// 	//d := strings.Split(c[0], " ")
		// 	//da := strings.Split(d[0], "/")
		// 	//ti := strings.Split(d[1], ":")
		// }

		// if errFetch2 != nil {
		consecutiveAtempts["FailedUpdateMarketDetails"]++
		return errFetch
		// }

	}

	//doc, err := htmlquery.LoadDoc(DEF_PATHS_MARKET_STATUS_PATH)
	doc, err := htmlquery.Parse(strings.NewReader(content))
	if err != nil {
		consecutiveAtempts["FailedUpdateMarketDetails"]++
		return errors.New("gathered data is incomplete")
	}

	container := htmlquery.Find(doc, "//table[@class='table1']//tr[1]//td")
	if container == nil {
		consecutiveAtempts["FailedUpdateMarketDetails"]++
		return errors.New("Time not found in gathered data")
	}

	containerText := htmlquery.InnerText(container[1])

	container2 := htmlquery.Find(doc, "//table[@class='table1']//tr[5]//td")
	if container2 == nil {
		consecutiveAtempts["FailedUpdateMarketDetails"]++
		return errors.New("Date not found in gathered data")
	}
	containerText2 := htmlquery.InnerText(container2[1])
	var dateParts []string = make([]string, 0)
	if containerText2 != "" {
		dateParts = strings.Split(containerText2, " ")
		dateParts[0] = strings.ReplaceAll(dateParts[0], "/", "-")
	}

	setDate(dateParts[0])

	// clean
	containerText = strings.Map(func(r rune) rune {
		if unicode.IsPrint(r) {
			return r
		}
		s := ' '
		return s
	}, containerText)

	parts := strings.Split(containerText, " ")

	if parts[0] == "باز" {
		setLastTimeMarketOpen()
		setMarketStatus(DEF_MARKET_STATUS_OPEN)
	} else {
		setMarketStatus(DEF_MARKET_STATUS_CLOSE)
	}

	if parts[1] == "" {
		parts[1] = "00:00"
	}
	//timeParts := strings.Split(parts[1], ":")

	/////////////////////////////////////////////////////

	content, errFetch = Fetch(DEF_URLS_PRECISE_TIME_URL, "", 0)

	if content == "" && errFetch != nil {
		return errFetch
	}

	var mtime map[string]interface{}
	json.Unmarshal([]byte(content), &mtime)

	h := fmt.Sprintf("%02d", int(mtime["hour"].(float64)))
	m := fmt.Sprintf("%02d", int(mtime["minute"].(float64)))

	setTime(h, m)

	// gathering data is done successfully
	consecutiveAtempts["FailedUpdateMarketDetails"] = 0

	marketStatusLastUpdate = time.Now()
	setLastTimeMarketUpdate()

	return nil
}

func setMarketStatus(st MarketStatusType) {
	marketStatus = st
}

func GetMarketStatus() MarketStatusType {
	return marketStatus
}

func setLastTimeMarketOpen() {
	lastTimeMarketOpen = time.Now().Format(DateTime) // Y-m-d H:i:s
}

func getLastTimeMarketOpen() string {
	if lastTimeMarketOpen == "" {
		lastTimeMarketOpen = "2001-01-01 00:00:00"
	}
	return lastTimeMarketOpen
}

func setLastTimeMarketUpdate() {
	lastTimeMarketUpdate = time.Now().Format(DateTime) // Y-m-d H:i:s
}

func getLastTimeMarketUpdate() string {
	if lastTimeMarketUpdate == "" {
		lastTimeMarketUpdate = "2001-01-01 00:00:00"
	}
	return lastTimeMarketUpdate
}

func MarketBeenOpenToday() bool {

	to, err := time.Parse(DateTime, getLastTimeMarketOpen())
	if err != nil {
		return false
	}

	now := time.Now()

	if to.Month() == now.Month() {
		if to.Day() == now.Day() {
			return true
		}
	}

	tu, err := time.Parse(DateTime, getLastTimeMarketUpdate())
	if err != nil {
		return false
	}

	if now.Sub(tu).Minutes() > 5 {
		updateMarketDetails()
	}

	return false
}

func setDate(date string) {
	marketDate = date
	if date[0:2] != "14" {
		date = "14" + date
	}

	marketDate = date
}

func GetDate() string {
	return marketDate
}

func setTime(hour string, minute string) {
	marketTime = hour + ":" + minute

	calibrationPointMarketTimeH := int(toInt(hour))
	calibrationPointMarketTimeM := int(toInt(minute))

	tz := time.Now()
	timeZone, _ := time.LoadLocation("Asia/Tehran")
	t := tz.In(timeZone)

	calibrationPointServerTimeH := int(toInt(t.Format("15")))
	calibrationPointServerTimeM := int(toInt(t.Format("04")))

	calibrationPointDiffTimeH = calibrationPointMarketTimeH - calibrationPointServerTimeH
	calibrationPointDiffTimeM = calibrationPointMarketTimeM - calibrationPointServerTimeM

}

func GetTime() string {
	updateMarketDetails()

	tz := time.Now()
	timeZone, _ := time.LoadLocation("Asia/Tehran")
	t := tz.In(timeZone)

	serverTimeH := int(toInt(t.Format("15"))) + calibrationPointDiffTimeH
	serverTimeM := int(toInt(t.Format("04"))) + calibrationPointDiffTimeM

	return fmt.Sprintf("%02d:%02d", serverTimeH, serverTimeM)
}

func IsOpenTime() bool {

	to, tc := GetMarketOpenCloseTime()

	a := time.Now().Sub(to).Seconds()
	b := time.Now().Sub(tc).Seconds()

	if a > 0 && b < 0 {
		//st := GetMarketStatus()
		//if st == DEF_MARKET_STATUS_OPEN {
		return true
		//}
	}

	return false
}

func MarketIsOpen() bool {
	return marketIsOpen
}

func SetMarketOpen() {
	marketIsOpen = true
}

func SetMarketClose() {
	marketIsOpen = false
}

func GetClockNumber() int64 {

	to, tc := GetMarketOpenCloseTime()

	_ = tc

	return int64(time.Now().Sub(to).Seconds())
}

func SetConfigSettings(options config.ConfigList) {
	settings = options
}

func GetMarketOpenCloseTime() (time.Time, time.Time) {

	timeZone, err := time.LoadLocation("Asia/Tehran")

	if err != nil {
		logrus.Error(err)
	}

	start := strings.Split(settings["start-time"], ":")
	end := strings.Split(settings["end-time"], ":")

	startH := stringsh.ToInt(start[0])
	startM := stringsh.ToInt(start[1])

	endH := stringsh.ToInt(end[0])
	endM := stringsh.ToInt(end[1])

	z := time.Now()
	o := time.Date(z.Year(), z.Month(), z.Day(), int(startH), int(startM), 0, 0, timeZone).UTC()
	c := time.Date(z.Year(), z.Month(), z.Day(), int(endH), int(endM), 0, 0, timeZone).UTC()

	//t := pt.Time().UTC()

	return o, c
}

func DepartMarketStatusPureContent(content string) (pureYear string, pureMonth string, pureDay string, pureHour string, pureMin string, pureSec string) {

	contentParts := strings.Split(content, "@")

	pureYear = ""
	pureMonth = ""
	pureDay = ""

	pureHour = ""
	pureMin = ""
	pureSec = ""

	if contentParts != nil && contentParts[1] != "" {
		c := strings.Split(contentParts[1], ",")
		d := strings.Split(c[0], " ")
		da := strings.Split(d[0], "/")
		ti := strings.Split(d[1], ":")

		pureYear = da[0]
		pureMonth = da[1]
		pureDay = da[2]

		pureHour = ti[0]
		pureMin = ti[1]
		pureSec = ti[2]

	}

	return
}
