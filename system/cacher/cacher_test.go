package cacher

// import (
// 	"testing"

// 	"BourseCrawler/system/cacher/adapter"

// 	"github.com/go-macaron/cache"
// )

// func TestCacher_Methods(t *testing.T) {
// 	ce, ct, err := NewCacherEngine(RedisCacher, adapter.CEOptions{
// 		Engine: adapter.CEGoMacaron,
// 		Options: cache.Options{
// 			Adapter:       "redis",
// 			AdapterConfig: "addr:127.0.0.1:6379",
// 		},
// 	})

// 	c, err := NewCacher(ct, ce)
// 	if err != nil {
// 		t.Errorf("NewCacher: %s", err.Error())
// 		return
// 	}
// 	c.Context = "Context"
// 	c.Container = "Container"
// 	c.Component = "Component"

// 	testCacher(c, t)
// }
