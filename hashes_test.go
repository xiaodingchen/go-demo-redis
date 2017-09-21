package redis

import (
	"testing"
)

func TestHashesStart(t *testing.T) {
	t.Run("all", func(t *testing.T) {
		t.Run("SetAndGet", func(t *testing.T) {
			testSetAndGet(t)
		})

		t.Run("KeyAndVal", func(t *testing.T) {
			testKeyAndVal(t)
		})
	})

}

func testKeyAndVal(t *testing.T) {

}

func testSetAndGet(t *testing.T) {
	hkey1 := "hkey1"
	var err error
	// err = client.Hset(hkey1, "name", []byte("jerry"))
	// if err != nil {
	// 	t.Fatalf("hset error: %s", err)
	// }
	// t.Logf("hset %s %s success.", hkey1, "name")
	// fv := make(map[string][]byte)
	// fv["age"] = []byte(strconv.FormatInt(23, 10))
	// fv["sex"] = []byte("man")
	// err = client.Hmset(hkey1, fv)
	// if err != nil {
	// 	t.Fatalf("hmset error: %s", err)
	// }
	// t.Logf("hmset %s success.", hkey1)
	// var val []byte
	// val, err = client.Hget(hkey1, "name")
	// if err != nil {
	// 	t.Fatalf("hget error: %s", err)
	// }
	// t.Logf("hget %s %s value is %s", hkey1, "name", string(val))

	// var vals [][]byte
	// vals, err = client.Hmget(hkey1, "age", "name")
	// if err != nil {
	// 	t.Fatalf("hmget error: %s", err)
	// }

	// t.Logf("hmget len is %d", len(vals))
	// t.Logf(string(vals[0]))
	// t.Logf(string(vals[1]))

	var kv map[string][]byte
	kv, err = client.Hgetall(hkey1)
	if err != nil {
		t.Fatalf("Hgetall error: %s", err)
	}
	t.Logf("kv: %v", kv)

}
