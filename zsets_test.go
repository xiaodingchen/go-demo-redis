package redis

import (
	"testing"
)

func TestZsetStart(t *testing.T) {
	t.Run("all", func(t *testing.T) {
		// t.Run("testCURD", func(t *testing.T) {
		// 	testCURD(t)
		// })

		// t.Run("testRANGE", func(t *testing.T) {
		// 	testRANGE(t)
		// })

		// t.Run("testRANK", func(t *testing.T) {
		// 	testRANK(t)
		// })

		// t.Run("testREM", func(t *testing.T) {
		// 	testREM(t)
		// })

		t.Run("testStore", func(t *testing.T) {
			testStore(t)
		})
	})
}

func testCURD(t *testing.T) {
	zkey := "myset1"
	// t.Logf("test add start")
	// bol, err := client.Zadd(zkey, 1, []byte("hello"))
	// client.Zadd(zkey, 2, []byte("world"))
	// client.Zadd(zkey, 1, []byte("hhhworld"))
	// if err != nil {
	// 	t.Fatalf("err is %s", err.Error())
	// }
	// if bol == false {
	// 	t.Logf("test add key: %s fail", zkey)
	// }
	// t.Logf("test add end")

	// t.Logf("test ZCARD")
	// res, err := client.Zcard(zkey)
	// if err != nil {
	// 	t.Fatalf("err is %s", err.Error())
	// }
	// t.Logf("num is %d", res)

	// t.Logf("test ZCOUNT")

	// res, err := client.Zcount(zkey, 0, 4)
	// if err != nil {
	// 	t.Fatalf("err is %s", err.Error())
	// }

	// t.Logf("num is %d", res)

	// t.Logf("test ZINCRBY")
	// res, err := client.Zincrby(zkey, 2, "hello")
	// if err != nil {
	// 	t.Fatalf("err is %s", err.Error())
	// }

	// t.Logf("zincrby hello: %s", res)
	// t.Logf("ZINCRBY END")

	t.Logf("test ZSCORE")
	res, err := client.Zscore(zkey, "hello")
	if err != nil {
		t.Fatalf("err is %s", err.Error())
	}
	t.Logf("key:%s,member:%s,score:%.2f", zkey, "hello", res)
	t.Logf("ZSCORE end")

}

func testRANGE(t *testing.T) {
	t.Logf("test range")
	zkey := "myset1"
	var res [][]byte
	var err error

	t.Logf("test ZRANGE")
	res, err = client.Zrange(zkey, 0, -1)
	if err != nil {
		t.Fatalf("err is %s", err.Error())
	}

	for k, v := range res {
		t.Logf("%d => %s", k, string(v))
	}

	t.Logf("ZRANGE end")

	t.Logf("test ZRANGEBYSCORE")
	res, err = client.Zrangebyscore(zkey, 1, 1)

	if err != nil {
		t.Fatalf("err is %s", err.Error())
	}

	for k, v := range res {
		t.Logf("%d => %s", k, string(v))
	}

	t.Logf("ZRANGEBYSCORE end")

	t.Logf("test ZREVRANGE")
	res, err = client.Zrevrange(zkey, 0, -1)

	if err != nil {
		t.Fatalf("err is %s", err.Error())
	}

	for k, v := range res {
		t.Logf("%d => %s", k, string(v))
	}

	t.Logf("ZREVRANGE end")

	t.Logf("test ZREVRANGEBYSCORE")
	res, err = client.Zrevrangebyscore(zkey, 1, 1)
	if err != nil {
		t.Fatalf("err is %s", err.Error())
	}

	for k, v := range res {
		t.Logf("%d => %s", k, string(v))
	}
	t.Logf("ZREVRANGEBYSCORE END")

}

func testRANK(t *testing.T) {
	t.Logf("test rank")
	t.Logf("test ZRANK")
	var res int
	var err error
	zkey := "myset1"
	res, err = client.Zrank(zkey, "world")

	if err != nil {
		t.Fatalf("err is %s", err.Error())
	}
	t.Logf("the member for world is rank: %d", res)
	t.Logf("ZRANK end")

	t.Log("test ZREVRANK")
	res, err = client.Zrevrank(zkey, "hello")
	if err != nil {
		t.Fatalf("err is %s", err.Error())
	}
	t.Logf("the member for hello is rank: %d", res)
	t.Logf("ZREVRANK end")
}

func testREM(t *testing.T) {
	t.Logf("test REM")
	var res bool
	var err error
	zkey := "myset1"

	// t.Logf("test ZREM")
	// res, err = client.Zrem(zkey, "hello")
	// if err != nil {
	// 	t.Fatalf("err is %s", err.Error())
	// }

	// t.Logf("ZREM RESULT: %v", res)
	// t.Log("ZREM end")

	// t.Logf("test ZREMRANGEBYRANK")
	// res, err = client.Zremrangebyrank(zkey, 1, 1)
	// if err != nil {
	// 	t.Fatalf("err is %s", err.Error())
	// }

	// t.Logf("ZREMRANGEBYRANK RESULT: %v", res)
	// t.Log("ZREMRANGEBYRANK end")

	t.Logf("test ZREMRANGEBYSCORE")
	res, err = client.Zremrangebyscore(zkey, 1, 1)
	if err != nil {
		t.Fatalf("err is %s", err.Error())
	}

	t.Logf("ZREMRANGEBYSCORE RESULT: %v", res)
	t.Log("ZREMRANGEBYSCORE end")
}

func testStore(t *testing.T) {
	key1 := "zkey1"
	key2 := "zkey2"
	client.Zadd(key1, 1, []byte("hello"))
	client.Zadd(key1, 1, []byte("a"))
	client.Zadd(key1, 2, []byte("b"))
	client.Zadd(key1, 0, []byte("c"))
	client.Zadd(key2, 2, []byte("b"))
	client.Zadd(key2, 0, []byte("a"))
	client.Zadd(key2, 1, []byte("c"))
	client.Zadd(key2, 1, []byte("d"))
	client.Zadd(key2, 2, []byte("hello"))
	t.Logf("test ZINTERSTORE")
	res, err := client.Zinterstore("zteststore", 2, key1, key2)
	if err != nil {
		t.Fatalf("err is %s", err.Error())
	}

	t.Logf("ZINTERSTORE result: %d", res)

	t.Log("test end")
}
