package redis

import (
	"fmt"
	"strconv"

	"github.com/fatih/structs"
	"github.com/go-redis/redis/v8"
)

// ConfigParser config parser for redis
func ConfigParser(conf interface{}) redis.Options {
	config := structs.Map(conf)

	opt := redis.Options{
		Network: "tcp",
		Addr:    "127.0.0.1:6379",
	}

	v, exist := config["Network"]
	if exist {
		opt.Network = fmt.Sprintf("%s", v)
	}
	v, exist = config["Host"]
	if exist {
		opt.Addr = fmt.Sprintf("%s:%s", config["Host"], config["Port"])
	}
	v, exist = config["DB"]
	if exist {
		a, _ := strconv.ParseInt(v.(string), 10, 64)
		opt.DB = int(a)
	}
	v, exist = config["MaxRetries"]
	if exist {
		a, _ := strconv.ParseInt(v.(string), 10, 64)
		opt.MaxRetries = int(a)
	}
	v, exist = config["PoolSize"]
	if exist {
		a, _ := strconv.ParseInt(v.(string), 10, 64)
		opt.PoolSize = int(a)
	}
	v, exist = config["MinIdleConns"]
	if exist {
		a, _ := strconv.ParseInt(v.(string), 10, 64)
		opt.MinIdleConns = int(a)
	}

	return opt
}

func first(n interface{}, _ interface{}) interface{} {
	return n
}
