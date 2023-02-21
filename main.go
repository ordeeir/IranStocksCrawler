package main

import (
	"IranStocksCrawler/bourse"
	"IranStocksCrawler/helpers/osh"
	"IranStocksCrawler/helpers/timeh"
	sys "IranStocksCrawler/system"
	"IranStocksCrawler/system/config"

	"fmt"
	"net/http"
	"time"
)

func main() {

	timeh.Init()

	conf, _, _ := config.NewConfig(osh.GetRootPath() + "/config.json")

	sys.SetConfig(conf)

	cacher := sys.CreateCacher()

	router := sys.CreateRouter()

	bourse.SetConfigSettings(conf.Settings)

	//bourse.GetTime()

	router.HttpGet("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Welcome lucky man")
	})

	router.HttpGet("/tostring", func(w http.ResponseWriter, r *http.Request) {

		s := bourse.ToString()

		w.Write([]byte(s))
	})

	router.HttpGet("/fetchUrl", func(w http.ResponseWriter, r *http.Request) {

		url := r.URL.Query().Get("url")
		path := r.URL.Query().Get("path")

		bourse.Fetch(url, path, time.Minute*5)

	})

	go func() {

		const tickPriceInterval = time.Second * 20
		const tickIOInterval = time.Second * 40
		const tickPeriodicAvgsInterval = time.Hour * 6

		tickPrices := time.NewTicker(time.Second * 1)
		tickIO := time.NewTicker(time.Second * 5)
		tickPeriodicAvgs := time.NewTicker(time.Second * 10)
		tickFilters := time.NewTicker(time.Second * 30)

		fmt.Println("\n---")

		for {
			select {

			case <-tickPrices.C:
				err := bourse.UpdatePrices(cacher)
				if err {
					tickPrices = time.NewTicker(time.Second * 5)
				} else {
					tickPrices = time.NewTicker(tickPriceInterval)
				}
				continue

			case <-tickIO.C:
				err := bourse.UpdateIO(cacher)
				if err {
					tickIO = time.NewTicker(time.Second * 5)
				} else {
					tickIO = time.NewTicker(tickIOInterval)
				}
				continue

			case <-tickPeriodicAvgs.C:
				// err := bourse.UpdatePeriodicAverages(cacher)
				// if err {
				// 	tickPeriodicAvgs = time.NewTicker(time.Second * 5)
				// } else {
				// 	tickPeriodicAvgs = time.NewTicker(tickPeriodicAvgsInterval)
				// }
				continue

			case <-tickFilters.C:
				tickFilters = time.NewTicker(time.Second * 20)

				bourse.InitFilters()

				bourse.ApplyFilters(cacher)

				continue

			}
		}

	}()

	router.HttpServe("1212")

}
