package bourse

import (
	"IranStocksCrawler/helpers/stringsh"
	"IranStocksCrawler/helpers/timeh"
	"IranStocksCrawler/system/config"
	"errors"
	"strings"
	"time"
	"unicode"

	"github.com/antchfx/htmlquery"
	ptime "github.com/yaa110/go-persian-calendar"
)

var marketStatus MarketStatusType = DEF_MARKET_STATUS_UNKNOWN
var marketTime string
var marketDate string
var marketStatusLastUpdate time.Time
var marketStatusLastCheck time.Time

var settings config.ConfigList

type MarketStatusType int

func updateMarketDetails() (string, error) {

	if consecutiveAtempts["FailedUpdateMarketDetails"] > 10 {
		return "", errors.New("many efforts without expectable result")
	}

	if isMarketUpdated() {
		return "", nil
	}

	conten, errFetch := Fetch(DEF_URLS_MARKET_STATUS_URL, DEF_PATHS_MARKET_STATUS_PATH, time.Second*0)
	if errFetch != nil {
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
		return "", errFetch
		// }

	}

	//doc, err := htmlquery.LoadDoc(DEF_PATHS_MARKET_STATUS_PATH)
	doc, err := htmlquery.Parse(strings.NewReader(conten))
	if err != nil {
		consecutiveAtempts["FailedUpdateMarketDetails"]++
		return "", errors.New("gathered data is incomplete")
	}

	container := htmlquery.Find(doc, "//table[@class='table1']//tr[1]//td")
	if container == nil {
		consecutiveAtempts["FailedUpdateMarketDetails"]++
		return "", errors.New("Time not found in gathered data")
	}

	containerText := htmlquery.InnerText(container[1])

	container2 := htmlquery.Find(doc, "//table[@class='table1']//tr[5]//td")
	if container2 == nil {
		consecutiveAtempts["FailedUpdateMarketDetails"]++
		return "", errors.New("Date not found in gathered data")
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
		setMarketStatus(DEF_MARKET_STATUS_OPEN)
	} else {
		setMarketStatus(DEF_MARKET_STATUS_CLOSE)
	}

	if parts[1] == "" {
		parts[1] = "00:00"
	}
	timeParts := strings.Split(parts[1], ":")

	setTime(timeParts[0], timeParts[1])

	// gathering data is done successfully
	consecutiveAtempts["FailedUpdateMarketDetails"] = 0

	marketStatusLastUpdate = time.Now()

	return "", nil
}

func setMarketStatus(st MarketStatusType) {
	marketStatus = st
}

func GetMarketStatus() MarketStatusType {
	updateMarketDetails()
	return marketStatus
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

	//dura := toInt(hour)*60 + toInt(minute)

	//marketOpenTime = time.Date(int(toInt(tStr[0])), time.Month(int(toInt(tStr[1]))), int(toInt(tStr[2])), int(toInt(tStr[3])), int(toInt(tStr[4])), int(toInt(tStr[5])), 0, time.UTC)
}

func GetTime() string {
	updateMarketDetails()
	return marketTime
}

func GetClockNumber() int64 {

	to, tc := GetMarketOpenCloseTime()

	_ = tc

	//if time.Now().Sub(to) > time.Hour*10 {

	//marketOpenTime = time.Date(2022, time.Month(10), time.Now().Day(), 9, 0, 0, 0, time.UTC)
	//marketOpenTime = t
	//}

	// //setStatus(DEF_MARKET_STATUS_OPEN) // del

	return int64(time.Now().Sub(to).Seconds())
}

func SetConfigSettings(options config.ConfigList) {
	settings = options
}

func GetMarketOpenCloseTime() (time.Time, time.Time) {

	timeZone, _ := time.LoadLocation("Asia/Tehran")

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

func __GetMarketStatus() MarketStatusType {

	tnow := time.Now()
	to, tc := GetMarketOpenCloseTime()

	if tnow.Unix() > tc.Unix() || tnow.Unix() < to.Unix() {
		return DEF_MARKET_STATUS_CLOSE
	}

	for i := 0; i < 4; i++ {

		content2, _ := Fetch(DEF_URLS_PRICE_URL, DEF_PATHS_MARKET_STATUS_PATH, time.Second*0)
		//content2, _ := Fetch(DEF_URLS_MARKET_STAT_URL, DEF_PATHS_MARKET_STATUS_PATH, time.Second*0)

		contentParts := strings.Split(content2, "@")

		if contentParts != nil && contentParts[1] != "" {
			c := strings.Split(contentParts[1], ",")
			d := strings.Split(c[0], " ")
			da := strings.Split(d[0], "/")
			ti := strings.Split(d[1], ":")

			tStr := strings.Split(timeh.JalaliToDate("14"+da[0]+"-"+da[1]+"-"+da[2]+" "+ti[0]+":"+ti[1]+":"+ti[2], "Y-m-d-H-i-s"), "-")

			t := time.Date(int(toInt(tStr[0])), time.Month(int(toInt(tStr[1]))), int(toInt(tStr[2])), int(toInt(tStr[3])), int(toInt(tStr[4])), 0, 0, time.UTC)

			if tnow.Unix() > t.Add(-5*time.Minute).Unix() && tnow.Unix() < t.Add(5*time.Minute).Unix() {
				return DEF_MARKET_STATUS_OPEN
			}

		}
	}

	return DEF_MARKET_STATUS_CLOSE
}

func isMarketUpdated() bool {

	pt := ptime.Now()

	if marketStatus == DEF_MARKET_STATUS_CLOSE {
		if pt.Format("hh:mm") >= "08:00" && pt.Format("hh:mm") < "08:02" {
			if pt.Format("hh:mm") >= "09:00" && pt.Format("hh:mm") < "09:02" {
				if pt.Format("hh:mm") >= "10:00" && pt.Format("hh:mm") < "10:02" {
					if time.Now().Sub(marketStatusLastCheck).Seconds() > 60 {
						marketStatusLastCheck = time.Now()
						return false
					}
				}
			}
		}
	}

	if marketStatus == DEF_MARKET_STATUS_OPEN {
		if pt.Format("hh:mm") >= "12:30" && pt.Format("hh:mm") < "12:32" {
			if pt.Format("hh:mm") >= "13:30" && pt.Format("hh:mm") < "13:32" {
				if pt.Format("hh:mm") >= "11:30" && pt.Format("hh:mm") < "11:32" {
					if time.Now().Sub(marketStatusLastCheck).Seconds() > 60 {
						marketStatusLastCheck = time.Now()
						return false
					}
				}
			}
		}
	}

	if marketStatus == DEF_MARKET_STATUS_UNKNOWN {
		if time.Now().Sub(marketStatusLastCheck).Seconds() > 5*60 {
			marketStatusLastCheck = time.Now()
			return false
		}
	}

	return true
}
