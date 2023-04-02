package main

import (
	"IranStocksCrawler/bourse"
	"IranStocksCrawler/helpers/osh"
	"IranStocksCrawler/helpers/timeh"
	driver "IranStocksCrawler/system"
	"IranStocksCrawler/system/config"
	"os"
	"strings"

	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {

	timeh.Init()

	loglevel := os.Getenv("LOG_LEVEL")
	if loglevel != "" {
		ll, err := logrus.ParseLevel(loglevel)
		if err == nil {
			logrus.SetLevel(ll)
		}
	}

	conf, _, _ := config.NewConfig(osh.GetRootPath() + "/config.json")

	driver.SetConfig(conf)

	storage := driver.CreateCacher()

	router := driver.CreateRouter()

	bourse.SetConfigSettings(conf.Settings)

	_ = storage.Put("z2", "22", 24*60*60)
	z := storage.Get("z2")
	if z != nil && z.(string) == "22" {
		logrus.Debug("Connected to Redis successfully")

	} else {
		logrus.Debug("Redis Connection failed")

	}

	//
	router.HttpGet("/", func(w http.ResponseWriter, r *http.Request) {

		resp := "Welcome lucky man"
		w.Write([]byte(resp))
	})

	router.HttpGet("/GetTodaySeries/{symbol}", func(w http.ResponseWriter, r *http.Request) {
		value := router.Var(r, "symbol")

		resp := bourse.GetTodaySeries(value)
		w.Write([]byte(resp))
	})

	router.HttpGet("/GetIndiOrga/{symbol}", func(w http.ResponseWriter, r *http.Request) {
		sym := router.Var(r, "symbol")

		list := bourse.GetIndiOrgaListSymbols()

		str := "Gathered List : " + strings.Join(list, ",") + " <br><br><br>"

		data, err := bourse.GetIndiOrga(sym)

		if err == nil {
			str += "Symbol: " + sym + "<br>"
			str += "LastUpdate: " + data.LastUpdate + "<br>"
			for i, j := range data.Days {
				str += i + ": " + "<br>"
				str += "---- AmountIndiBuy: " + j.AmountIndiBuy + "<br>"
				str += "---- AmountIndiSell: " + j.AmountIndiSell + "<br>"
				str += "---- AmountOrgaBuy: " + j.AmountOrgaBuy + "<br>"
				str += "---- AmountOrgaSell: " + j.AmountOrgaSell + "<br>"
				str += "---- QuantityIndiBuy: " + j.QuantityIndiBuy + "<br>"
				str += "---- QuantityIndiSell: " + j.QuantityIndiSell + "<br>"
				str += "---- QuantityOrgaBuy: " + j.QuantityOrgaBuy + "<br>"
				str += "---- QuantityOrgaSell: " + j.QuantityOrgaSell + "<br>"
				str += "---- VolumeIndiBuy: " + j.VolumeIndiBuy + "<br>"
				str += "---- VolumeIndiSell: " + j.VolumeIndiSell + "<br>"
				str += "---- VolumeOrgaBuy: " + j.VolumeOrgaBuy + "<br>"
				str += "---- VolumeOrgaSell: " + j.VolumeOrgaSell + "<br>"

			}
		} else {
			str += "Symbol not found "

		}

		w.Write([]byte(str))
	})

	router.HttpGet("/CrawlerStatus", func(w http.ResponseWriter, r *http.Request) {

		s := "Root Path: " + osh.GetRootPath() + "\r\n"
		s += "Prices Updated At : " + timeh.TimeFormat(bourse.LastSuccessGatheringPrices, "Y-m-d H:i:s") + "\r\n"

		w.Write([]byte(s))
	})

	router.HttpGet("/Delete", func(w http.ResponseWriter, r *http.Request) {

		bourse.ResetGatheredData(storage)

		s := "Deleted.\r\n"

		w.Write([]byte(s))
	})

	router.HttpGet("/MarketString", func(w http.ResponseWriter, r *http.Request) {

		content2, _ := bourse.Fetch(bourse.DEF_URLS_PRICE_URL, bourse.DEF_PATHS_MARKET_STATUS_PATH, time.Second*0)

		pureYear, pureMonth, pureDay, pureHour, pureMin, pureSec := bourse.DepartMarketStatusPureContent(content2)

		s := "y: " + pureYear + "\r\n"
		s += "m: " + pureMonth + "\r\n"
		s += "d: " + pureDay + "\r\n"
		s += "h: " + pureHour + "\r\n"
		s += "i: " + pureMin + "\r\n"
		s += "s: " + pureSec + "\r\n"

		w.Write([]byte(s))
	})

	router.HttpGet("/Tostring", func(w http.ResponseWriter, r *http.Request) {

		s := bourse.ToString()

		w.Write([]byte(s))
	})

	router.HttpGet("/FetchUrl", func(w http.ResponseWriter, r *http.Request) {

		url := r.URL.Query().Get("url")
		path := r.URL.Query().Get("path")

		bourse.Fetch(url, path, time.Minute*5)

	})

	go func() {

		const tickPriceInterval = time.Second * 20
		const tickIOInterval = time.Second * 40
		const tickPeriodicAvgsInterval = time.Second * 60

		tickPrices := time.NewTicker(time.Second * 1)
		tickIO := time.NewTicker(time.Second * 5)
		tickPeriodicAvgs := time.NewTicker(time.Second * 10)
		tickIndiOrga := time.NewTicker(time.Second * 5)
		tickFilters := time.NewTicker(time.Second * 30)

		fmt.Println("\n---")

		for {
			select {

			case <-tickPrices.C:

				err := bourse.UpdatePrices(storage)
				if err {
					tickPrices = time.NewTicker(time.Second * 5)
				} else {
					tickPrices = time.NewTicker(tickPriceInterval)
				}
				continue

			case <-tickIO.C:

				err := bourse.UpdateIO(storage)
				if err {
					tickIO = time.NewTicker(time.Second * 5)
				} else {
					tickIO = time.NewTicker(tickIOInterval)
				}
				continue

			case <-tickPeriodicAvgs.C:

				err := bourse.UpdatePeriodicAverages(storage)
				if err {
					tickPeriodicAvgs = time.NewTicker(time.Second * 20)
				} else {
					tickPeriodicAvgs = time.NewTicker(tickPeriodicAvgsInterval)
				}
				continue

			case <-tickIndiOrga.C:

				_ = bourse.UpdateIndiOrga365Days(storage)

				continue

			case <-tickFilters.C:

				tickFilters = time.NewTicker(time.Second * 20)

				bourse.InitFilters()

				bourse.ApplyFilters(storage)

				continue

			}
		}

	}()

	router.HttpServe("1212")

}

func convert(str string) string {
	return (str)
}
