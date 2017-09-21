package redis

import (
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"testing"
)

func TestStart(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	t.Run("all", func(t *testing.T) {
		t.Run("newclient", func(t *testing.T) {
			testNew(t)
		})
		t.Run("ping", func(t *testing.T) {
			testPing(t)
		})
		t.Run("stringSetAndGet", func(t *testing.T) {
			testStringSetAndGet(t)
		})
	})
}

func testNew(t *testing.T) {
	addr := "127.0.0.1:6379"
	db := 4
	client := NewRedisClient(addr, uint(db), "", 0)
	if client == nil {
		t.Fatalf("创建Redis实例错误， 地址：%s", addr)
	}

	t.Logf("创建Redis实例成功，地址：%s", addr)

}

func testPing(t *testing.T) {
	addr := "127.0.0.1:6379"
	db := 4
	client := NewRedisClient(addr, uint(db), "", 0)
	if client == nil {
		t.Fatalf("创建Redis实例错误， 地址：%s", addr)
	}

	err := client.Ping()
	if err != nil {
		t.Fatalf("ping出错：%s", err)
	}

	t.Logf("ping成功，地址：%s", addr)
}

func testStringSetAndGet(t *testing.T) {
	addr := "127.0.0.1:6379"
	db := 4
	client := NewRedisClient(addr, uint(db), "", 0)
	if client == nil {
		t.Fatalf("创建Redis实例错误， 地址：%s", addr)
	}

	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 5; i++ {
		go func(n int) {
			defer wg.Done()
			key := fmt.Sprintf("g%d", n)
			val := strconv.FormatInt(int64(n*10), 10)
			err := client.Set(key, []byte(val))
			if err != nil {
				t.Logf("Set error:%s", err)
			}
			t.Logf("setg%d, key:%s,val:%s", n, key, val)

		}(i)
	}

	for i := 0; i < 5; i++ {
		go func(n int) {
			defer wg.Done()
			//t.Logf("getg%d", n)
			key := fmt.Sprintf("g%d", n)
			getval, err := client.Get(key)
			if err != nil {
				t.Logf("get error:%s", err)
			}
			t.Logf("val is:%s", string(getval))
		}(i)
	}

	wg.Wait()
}
