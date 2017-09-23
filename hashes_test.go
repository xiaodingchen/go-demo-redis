package redis

import (
	"testing"
)

func TestHashesStart(t *testing.T) {
	t.Run("all", func(t *testing.T) {
		// t.Run("SetAndGet", func(t *testing.T) {
		// 	testSetAndGet(t)
		// })

		// t.Run("KeyAndVal", func(t *testing.T) {
		// 	testKeyAndVal(t)
		// })
		t.Run("Incr", func(t *testing.T) {
			testIncr(t)
		})
	})

}

func testIncr(t *testing.T) {
	hkeynum := "hkeynum"
	// vals := make(map[string][]byte)
	// vals["int"] = []byte(strconv.FormatInt(2, 10))
	// vals["float"] = []byte(strconv.FormatFloat(2.333, 'e', 4, 64))
	// client.Hmset(hkeynum, vals)
	client.Hincrby(hkeynum, "int", -1)
}

func testKeyAndVal(t *testing.T) {
	var err error
	var flag bool
	// flag, err = client.Hsetnx("hkey31", "day", []byte(time.Now().Format("2006-01-02 15:04:05")))
	// if err != nil {
	// 	t.Fatalf("hsetnx error: %s", err)
	// }

	// t.Logf("hsetnx result: %v", flag)
	flag, err = client.Hexists("hkey44", "name")

	if err != nil {
		t.Fatalf("hexists error: %s", err)
	}

	t.Logf("hexists result: %v", flag)

	flag, err = client.Hdel("hkey1", "name", "age")
	if err != nil {
		t.Fatalf("hdel error: %s", err)
	}
	t.Logf("hdel result: %v", flag)
	count, err := client.Hlen("hkey1")
	if err != nil {
		t.Fatalf("hlen error: %s", err)
	}
	t.Logf("hlen result: %v", count)

	keys, err := client.Hkeys("hkey1")
	if err != nil {
		t.Fatalf("hkeys error: %s", err)
	}
	t.Logf("hkeys result: %v", keys)

	vals, err := client.Hvals("hkey1")
	if err != nil {
		t.Fatalf("hvals error: %s", err)
	}

	t.Logf("hvals result: %v", vals)
	t.Logf("vals[1]: %s", string(vals[1]))
	fc, err := client.Hstrlen("hkey1", "day")
	if err != nil {
		t.Fatalf("hstrlen error: %s", err)
	}

	t.Logf("hstrlen result: %v", fc)
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
