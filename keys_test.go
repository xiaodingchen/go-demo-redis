package redis

import (
	"testing"
)

func TestKeys(t *testing.T) {
	t.Run("all", func(t *testing.T) {
		t.Run("randomkey", func(t *testing.T) {
			testRandomKey(t)
		})

		t.Run("typekey", func(t *testing.T) {
			testTypeKey(t)
		})
		testKeys(t)
	})
}
func testKeys(t *testing.T) {
	addr := "127.0.0.1:6379"
	db := 4
	client := NewRedisClient(addr, uint(db), "", 0)
	if client == nil {
		t.Fatalf("创建Redis实例错误， 地址：%s", addr)
	}
	keys, err := client.Keys("*")
	if err != nil {
		t.Fatalf("testkeys error: %s", err)
	}

	t.Logf("keys total is: %d", len(keys))
}
func testRandomKey(t *testing.T) {
	addr := "127.0.0.1:6379"
	db := 4
	client := NewRedisClient(addr, uint(db), "", 0)
	if client == nil {
		t.Fatalf("创建Redis实例错误， 地址：%s", addr)
	}
	key, err := client.RandomKey()
	if err != nil {
		t.Fatalf("random error: %s", err)
	}

	t.Logf("randomkey: %s", key)
}

func testTypeKey(t *testing.T) {
	addr := "127.0.0.1:6379"
	db := 4
	client := NewRedisClient(addr, uint(db), "", 0)
	if client == nil {
		t.Fatalf("创建Redis实例错误， 地址：%s", addr)
	}

	keytype, err := client.Type("g3")
	if err != nil {
		t.Fatalf("typekey error: %s", err)
	}

	t.Logf("g3 type is: %s", keytype)
}
