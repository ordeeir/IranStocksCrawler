package gomacaron

/*
import (
	"fmt"
	"testing"
	"time"

	_ "github.com/go-macaron/cache/redis"
	"github.com/muharihar/d3ta-go/system/config"
)

func newConfig(t *testing.T) (*config.Config, error) {
	c, _, err := config.NewConfig("../../../../conf")
	if err != nil {
		return nil, err
	}
	return c, nil
}

func TestGoMacaron_Methods(t *testing.T) {
	cfg, err := newConfig(t)
	if err != nil {
		t.Errorf("newConfig: %v", err.Error())
		return
	}

	c, err := NewCacher(
		"redis",
		/*
			cache.Options{
				Adapter:       "redis",
				AdapterConfig: "addr:127.0.0.1:6379",
				...
			})
*/ /*
		ConfigParser(cfg.Caches.TmpDataCache.Configurations),
	)

	if err != nil {
		t.Errorf("NewRedisCacher: %v", err.Error())
		return
	}

	t.Log("Put: key1")
	value := fmt.Sprintf("value-%s", time.Now().String())
	err = c.Put("key1", value, 0)
	if err != nil {
		t.Errorf("Put: %v", err.Error())
		return
	}

	t.Log("Get: key1")
	g := c.Get("key1")
	t.Logf("Value: %#v", g)
	if g == nil {
		t.Errorf("Get: %#v", g)
		return
	}

	t.Log("IsExist: key1")
	exist := c.IsExist("key1")
	t.Logf("Exist: %#v", exist)

	t.Log("Delete: key1")
	err = c.Delete("key1")
	if err != nil {
		t.Errorf("Delete: %s", err.Error())
		return
	}

	t.Log("IsExist: key1")
	exist = c.IsExist("key1")
	t.Logf("Exist: %#v", exist)

	t.Log("Delete: key-not-found")
	err = c.Delete("key-not-found")
	if err != nil {
		t.Errorf("Delete: %s", err.Error())
		return
	}

	err = c.Put("key-increment", 134, 0)
	t.Log("Incr: key-increment")
	err = c.Incr("key-increment")
	if err != nil {
		t.Errorf("Incr: %s", err.Error())
		return
	}

	t.Log("Get: key-increment")
	g = c.Get("key-increment")
	t.Logf("Value: %#v", g)
	if g == nil {
		t.Errorf("Get: %#v", g)
		return
	}

	t.Log("Decr: key-increment")
	err = c.Decr("key-increment")
	if err != nil {
		t.Errorf("Decr: %s", err.Error())
		return
	}

	t.Log("Get: key-increment")
	g = c.Get("key-increment")
	t.Logf("Value: %#v", g)
	if g == nil {
		t.Errorf("Get: %#v", g)
		return
	}

	/*
		t.Log("Flush")
		err = c.Flush()
		if err != nil {
			t.Errorf("Flush: %s", err.Error())
			return
		}
*/ /*
}
*/
