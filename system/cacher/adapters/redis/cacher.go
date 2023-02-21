package redis

import (
	"IranStocksCrawler/system/cacher"
	"context"
	"encoding/json"
	"fmt"
	"time"

	redisv8 "github.com/go-redis/redis/v8"
)

// NewCacher new Redis Cacher
func NewCacher(options *redisv8.Options) (cacher.ICacherEngine, error) {
	c := new(Cacher)
	c.ctx = context.Background()
	c.engine = redisv8.NewClient(options)
	return c, nil
}

// Cacher represent Redis Cacher
type Cacher struct {
	engine *redisv8.Client
	ctx    context.Context
}

// GetEngine get default Engine
func (c *Cacher) GetEngine() interface{} {
	return c.engine
}

// IsExist check key is exist
func (c *Cacher) IsExist(key string) bool {
	r, err := c.engine.Exists(c.ctx, key).Result()
	if err != nil {
		return false
	}
	return r > 0
}

// Put put key value with (timeout) expiration
func (c *Cacher) Put(key string, val interface{}, expirationSeconds int64) error {

	v, _ := json.Marshal(val)

	expiration, err := time.ParseDuration(fmt.Sprintf("%ds", expirationSeconds))
	if err != nil {
		return err
	}

	return c.engine.Set(c.ctx, key, string(v), expiration).Err()
}

// Get key value
func (c *Cacher) Get(key string) interface{} {

	r, err := c.engine.Get(c.ctx, key).Result()

	if err != nil {
		return nil
	}
	byt := []byte(r)

	var v interface{}
	_ = json.Unmarshal(byt, &v)

	return v
}

// Incr increment key value
func (c *Cacher) Incr(key string) error {
	return c.engine.Incr(c.ctx, key).Err()
}

// Decr decrement key value
func (c *Cacher) Decr(key string) error {
	return c.engine.Decr(c.ctx, key).Err()
}

// Delete delete key
func (c *Cacher) Delete(key string) error {
	return c.engine.Del(c.ctx, key).Err()
}

// Flush flush all data
func (c *Cacher) Flush() error {
	return c.engine.FlushAll(c.ctx).Err()
}
