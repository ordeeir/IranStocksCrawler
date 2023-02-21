package bourse

import "IranStocksCrawler/helpers/osh"

var (

	//DEF_URLS_MARKET_STATUS_URL       = "http://www.tsetmc.com/Loader.aspx?ParTree=15"
	DEF_URLS_MARKET_STATUS_URL = "https://pay.destinia.ir/includes/curl.php?u=aHR0cDovL3d3dy50c2V0bWMuY29tL0xvYWRlci5hc3B4P1BhclRyZWU9MTUNCg=="
	//DEF_URLS_PRICE_URL               = "http://www.tsetmc.com/tsev2/data/MarketWatchPlus.aspx?h=0&r=2222"
	DEF_URLS_PRICE_URL = "https://pay.destinia.ir/includes/curl.php?u=aHR0cDovL3d3dy50c2V0bWMuY29tL3RzZXYyL2RhdGEvTWFya2V0V2F0Y2hQbHVzLmFzcHg/aD0wJnI9MjIyMg=="
	//DEF_URLS_IO_URL                  = "http://www.tsetmc.com/tsev2/data/ClientTypeAll.aspx"
	DEF_URLS_IO_URL = "https://pay.destinia.ir/includes/curl.php?u=aHR0cDovL3d3dy50c2V0bWMuY29tL3RzZXYyL2RhdGEvQ2xpZW50VHlwZUFsbC5hc3B4"
	//DEF_URLS_PERIODIC_AVERAGES_URL   = "http://tsetmc.com/tsev2/data/InstValue.aspx?t=a"
	DEF_URLS_PERIODIC_AVERAGES_URL = "https://pay.destinia.ir/includes/curl.php?u=aHR0cDovL3RzZXRtYy5jb20vdHNldjIvZGF0YS9JbnN0VmFsdWUuYXNweD90PWE="

	DEF_PATHS_PRICE_PATH             = osh.GetRootPath() + "/files/bourse/realtime/prices.txt"
	DEF_PATHS_IO_PATH                = osh.GetRootPath() + "/files/bourse/realtime/indi_orga_details.txt"
	DEF_PATHS_MARKET_STATUS_PATH     = osh.GetRootPath() + "/files/bourse/realtime/market_status.txt"
	DEF_PATHS_PERIODIC_AVERAGES_PATH = osh.GetRootPath() + "/files/bourse/realtime/periodic_avgs.txt"

	DEF_MARKET_STATUS_OPEN    MarketStatusType = 1
	DEF_MARKET_STATUS_CLOSE   MarketStatusType = 0
	DEF_MARKET_STATUS_UNKNOWN MarketStatusType = -1
)
