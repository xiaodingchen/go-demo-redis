package redis

import (
	"testing"
)

func TestSetsStart(t *testing.T) {
	t.Run("all", func(t *testing.T) {
		setkey1 := "setkey1"
		setkey2 := "setkey2"
		var err error

		// t.Logf("sadd start")
		// num, err := client.Sadd(setkey1, []byte(random(5, 3)), []byte(random(4, 1)))
		// if err != nil {
		// 	t.Fatalf("sadd error: %s", err)
		// }

		// t.Logf("sadd success: %d", num)

		// strs := []string{"test", "444"}
		// setSliceString(strs)
		// t.Logf("slice: %v", strs)

		// count, err := client.Scard(setkey1)
		// if err != nil {
		// 	t.Fatalf("scard: %s", err)
		// }
		// t.Logf("count: %d", count)

		// for i := 0; i < 7; i++ {
		// 	client.Sadd(setkey2, []byte(random(5, 3)), []byte(random(4, 1)))
		// }
		// data, err := client.Sdiff(setkey2, setkey1)
		// if err != nil {
		// 	t.Fatalf("Sdiff: %s", err)
		// }
		// for k, v := range data {
		// 	t.Logf("%d => %s", k, string(v))
		// }
		// nums, err := client.Sdiffstore("dest1", setkey1, setkey2)
		// if err != nil {
		// 	t.Fatalf("sdiffstore error: %s", err)
		// }
		// t.Logf("sdiffstore result: %v", nums)
		// data, err := client.Sinter(setkey1, setkey2)
		// if err != nil {
		// 	t.Fatalf("sinter error: %s", err)
		// }
		// for k, v := range data {
		// 	t.Logf("%d => %s", k, string(v))
		// }
		// nums, err := client.Sinterstore("dest2", setkey1, setkey2)
		// if err != nil {
		// 	t.Fatalf("sinterstore error: %s", err)
		// }

		// t.Logf("sinterstore result: %v", nums)
		// flag, err := client.Sismember(setkey1, "demo")
		// if err != nil {
		// 	t.Fatalf("sismember error: %s", err)
		// }

		// t.Logf("%s result: %v", setkey1, flag)
		// flag, _ = client.Sismember(setkey2, "demo")
		// t.Logf("%s result: %v", setkey2, flag)
		// data, err := client.Smembers(setkey1)
		// if err != nil {
		// 	t.Fatalf("smembers error: %s", err)
		// }
		// count, _ := client.Scard(setkey2)
		// for k, v := range data {
		// 	t.Logf("%d=>%s", k, string(v))
		// }
		// t.Logf("%s count: %d", setkey2, count)
		// flag, err := client.Smove(setkey1, setkey2, "demo")
		// if err != nil {
		// 	t.Fatalf("smove error: %s", err)
		// }

		// t.Logf("smove result: %v", flag)
		// data, err := client.Spop(setkey1)
		// if err != nil {
		// 	t.Fatalf("spop error: %s", err)
		// }
		// t.Logf("spop result: %s", string(data))
		// client.Spop(setkey2)
		// data, err := client.Srandmember(setkey1)
		// if err != nil {
		// 	t.Fatalf("srandmember error: %s", err)
		// }

		// t.Logf("data: %s", string(data))
		// client.Srem(setkey2, "demo")
		// data, err := client.Sunion(setkey1, setkey2)
		// if err != nil {
		// 	t.Fatalf("sunion error: %s", err)
		// }

		// for k, v := range data {
		// 	t.Logf("%d => %s", k, string(v))
		// }

		_, err = client.Sunionstore("dest3", setkey1, setkey2)
		if err != nil {
			t.Fatalf("sunionstore error: %s", err)
		}
		data, _ := client.Smembers("dest3")

		for k, v := range data {
			t.Logf("%d=>%s", k, string(v))
		}
	})
}

func setSliceString(strs []string) {
	str := "demo"
	for k := range strs {
		strs[k] = str
	}
}
