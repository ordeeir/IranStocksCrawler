package sys

import (
	"IranStocksCrawler/helpers/stringsh"
	"IranStocksCrawler/system/cacher"
	"IranStocksCrawler/system/cacher/adapters/redis"
	"IranStocksCrawler/system/config"
	"IranStocksCrawler/system/router"
	"IranStocksCrawler/system/router/adapters/mux"

	redisv8 "github.com/go-redis/redis/v8"
)

var configData config.Config

func SetConfig(conf *config.Config) bool {

	configData = *conf
	return true
}

// cacher factory method
func CreateCacher() *cacher.Cacher {

	c := &cacher.Cacher{}

	if configData.Cache.Type == "redis" {

		redisRow := configData.Cache.Options["redis"]
		structure := redisv8.Options{}
		structure.Network = redisRow["network"]
		structure.Addr = redisRow["host"] + ":" + redisRow["port"]
		structure.Username = redisRow["username"]
		structure.Password = redisRow["password"]
		structure.DB = int(stringsh.ToInt(redisRow["db"]))
		structure.MaxRetries = int(stringsh.ToInt(redisRow["maxRetries"]))
		structure.PoolSize = int(stringsh.ToInt(redisRow["poolSize"]))
		structure.MinIdleConns = int(stringsh.ToInt(redisRow["minIdleConns"]))

		engine, _ := redis.NewCacher(&structure)
		c, _ = cacher.NewCacher(cacher.RedisCacher, engine)

	}
	if configData.Cache.Type == "filesystem" {

	}
	return c
}

// Router factory method
func CreateRouter() router.IRouter {

	if configData.Router.Type == string(router.MuxRouter) {

		//muxOptions := configData.Router.Options

		router := mux.NewMuxRouter()
		return router
	}

	return nil
}
