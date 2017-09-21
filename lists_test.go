package redis

import (
	"math/rand"
	"testing"
	"time"
)

var client *Client

func init() {
	addr := "127.0.0.1:6379"
	db := 4
	client = NewRedisClient(addr, uint(db), "", 0)
}

func TestListsStart(t *testing.T) {
	t.Run("all", func(t *testing.T) {
		// t.Run("PushAndPop", func(t *testing.T) {
		// 	testPushAndPop(t)
		// })
		// t.Run("PushX", func(t *testing.T) {
		// 	testPushX(t)
		// })
		// t.Run("Index", func(t *testing.T) {
		// 	testListsIndex(t)
		// })
		// t.Run("Index rm", func(t *testing.T) {
		// 	testListsRm(t)
		// })
		t.Run("bpopAndPush", func(t *testing.T) {
			testListsBpopAndPush(t)
		})
		t.Run("rpopLpush", func(t *testing.T) {
			testRpopLpush(t)
		})
	})
}

func testPushAndPop(t *testing.T) {
	//lkey1 := "lkey1"
	// t.Logf("lpush %s elem start.", lkey1)
	// for i := 0; i < 3; i++ {
	// 	lval1 := random(4, 3)
	// 	err := client.Lpush(lkey1, lval1)
	// 	if err != nil {
	// 		t.Fatalf("lpush %s elem error: %s", lkey1, err)
	// 	}
	// 	t.Logf("lpush %s elem success. val is: %s", lkey1, string(lval1))
	// }
	// t.Logf("lpush end.")

	// t.Logf("lpop start")
	// res, err := client.Lpop(lkey1)
	// if err != nil {
	// 	t.Fatalf("lpop %s elem error: %s", lkey1, err)
	// }
	// t.Logf("key %s index 0 val is: %s", lkey1, string(res))

	// res, err = client.Rpop(lkey1)

	// if err != nil {
	// 	t.Fatalf("rpop %s elem error: %s", lkey1, err)
	// }
	// t.Logf("key %s index 0 val is: %s", lkey1, string(res))

	rkey1 := "rkey1"
	t.Logf("rpush start")
	err := client.Rpush(rkey1, random(5, 2))
	if err != nil {
		t.Fatalf("rpush %s elem error: %s", rkey1, err)
	}

	t.Logf("rpush end.")
}

func testPushX(t *testing.T) {
	t.Logf("pushx start")
	pushxKey := "pushx"
	flag, err := client.Exists(pushxKey)
	if err != nil {
		t.Fatalf("exists key error: %s", err)
	}
	if flag {
		err = client.Lpushx(pushxKey, random(3, 3))
		if err != nil {
			t.Fatalf("pushx error: %s", err)
		}

		err = client.Rpush(pushxKey, random(6, 4))
	} else {
		err := client.Lpush(pushxKey, random(3, 2))
		if err != nil {
			t.Fatalf("lpush error: %s", err)
		}
	}

	t.Logf("pushx end.")
}

func testListsIndex(t *testing.T) {
	t.Logf("index start")
	key1 := "lkey1"
	key2 := "pushx"
	var res []byte
	var err error
	key1len, _ := client.Llen(key1)
	res, err = client.Lindex(key1, 4)
	if err != nil {
		t.Fatalf("lindex error: %s", err)
	}

	t.Logf("%s length is: %d", key1, int(key1len))
	t.Logf("lindex %s %d val: %s", key1, 2, string(res))
	t.Logf("lrange start")
	reses, err := client.Lrange(key2, 0, 4)
	if err != nil {
		t.Fatalf("lrange error: %s", err)
	}
	for i, res := range reses {
		t.Logf("key %s index %d val is: %s", key2, i, string(res))
	}

	t.Logf("lset start")
	key3 := "rkey1"
	err = client.Lset(key3, 0, random(6, 1))
	if err != nil {
		t.Fatalf("lset error: %s", err)

	}
	t.Logf("lset end")
}

func testListsRm(t *testing.T) {
	// for i := 0; i < 5; i++ {
	// 		client.Lpush("rmkey1", random(5, 2))
	// }
	// client.Lset("rmkey1", 0, []byte("hello"))
	// client.Lset("rmkey1", 3, []byte("hello"))
	res, err := client.Lrem("rmkey1", 2, []byte("hello"))
	if err != nil {
		t.Fatalf("lrem error: %s", err)
	}

	t.Logf("res: %v", res)
	t.Logf("ltrim start")
	err = client.Ltrim("pushx", 3, 9)
	if err != nil {
		t.Fatalf("ltrim error: %s", err)
	}
	t.Logf("ltrim end")

}

func testListsBpopAndPush(t *testing.T) {
	syncChan := make(chan struct{}, 1)
	bkey1 := "bkey1"
	go func() {

		k, v, err := client.Brpop([]string{bkey1}, 3)
		if err != nil {
			//syncChan <- struct{}{}
			t.Fatalf("blpop error: %s", err)
		}
		if k == nil && v == nil {
			t.Fatalf("key not exists, timeout")
		}

		t.Logf("k: %s, v: %s", *k, string(v))
		syncChan <- struct{}{}
	}()

	go func() {
		err := client.Lpush(bkey1, []byte("hello"))
		if err != nil {
			t.Fatalf("Lpush error: %s", err)
		}
	}()
	<-syncChan
}

func testRpopLpush(t *testing.T) {
	// v, err := client.RpopLpush("rkey2", "lkey1")
	v, err := client.BrpopLpush("rkey2", "lkey1", 3)
	if err != nil {
		t.Fatalf("RpopLpush error: %s", err)
	}

	t.Logf("v: %s", string(v))
}

/**
* size 随机码的位数
* kind 0    // 纯数字
       1    // 小写字母
       2    // 大写字母
       3    // 数字、大小写字母
*/

func random(size int, kind int) []byte {
	ikind, kinds, result := kind, [][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}, make([]byte, size)
	isAll := kind > 2 || kind < 0
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		if isAll { // random ikind
			ikind = rand.Intn(3)
		}
		scope, base := kinds[ikind][0], kinds[ikind][1]
		result[i] = uint8(base + rand.Intn(scope))
	}
	return result
}
