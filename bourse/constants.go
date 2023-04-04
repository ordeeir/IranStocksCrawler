package bourse

import "IranStocksCrawler/helpers/osh"

var (
	DEF_URLS_MARKET_STATUS_URL      = "http://www.tsetmc.com/Loader.aspx?ParTree=15"
	DEF_URLS_PRICE_URL              = "http://www.tsetmc.com/tsev2/data/MarketWatchPlus.aspx?h=0&r=2222"
	DEF_URLS_IO_URL                 = "http://www.tsetmc.com/tsev2/data/ClientTypeAll.aspx"
	DEF_URLS_PERIODIC_AVERAGES_URL  = "http://tsetmc.com/tsev2/data/InstValue.aspx?t=a"
	DEF_URLS_INDIORGA_DAYS_DATA_URL = "https://members.tsetmc.com/tsev2/data/clienttype.aspx?i={TSE_CODE}"
	DEF_URLS_PRECISE_TIME_URL       = "https://timeapi.io/api/Time/current/zone?timeZone=Asia/Tehran"

	DEF_PATHS_PRICE_PATH              = osh.GetRootPath() + "/files/bourse/realtime/prices.txt"
	DEF_PATHS_IO_PATH                 = osh.GetRootPath() + "/files/bourse/realtime/indi_orga_details.txt"
	DEF_PATHS_MARKET_STATUS_PATH      = osh.GetRootPath() + "/files/bourse/realtime/market_status.txt"
	DEF_PATHS_PERIODIC_AVERAGES_PATH  = osh.GetRootPath() + "/files/bourse/realtime/periodic_avgs.txt"
	DEF_PATHS_INDIORGA_DAYS_DATA_PATH = osh.GetRootPath() + "/files/bourse/indi_orga_days/{TSE_CODE}.txt"

	DEF_MARKET_STATUS_OPEN    MarketStatusType = 1
	DEF_MARKET_STATUS_CLOSE   MarketStatusType = 0
	DEF_MARKET_STATUS_UNKNOWN MarketStatusType = -1
)
