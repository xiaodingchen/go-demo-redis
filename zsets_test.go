package redis

import (
	"testing"
)

func TestZsetStart(t *testing.T) {
	t.Run("all", func(t *testing.T) {
		t.Run("testCURD", func(t *testing.T) {
			testCURD(t)
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

	t.Logf("test ZCOUNT")

	res, err := client.Zcount(zkey, 0, 4)
	if err != nil {
		t.Fatalf("err is %s", err.Error())
	}

	t.Logf("num is %d", res)

}
