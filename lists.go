package redis

import (
	"strconv"
)

// Redis列表操作

// Lpush 将所有指定的值插入到存于 key 的列表的头部。
// func (client *Client) Lpush(key string, val ...string) (int64, error) {
// 	args := []string{key}
// 	args = append(args, val...)
// 	res, err := client.sendCommand("LPUSH", args...)
// 	if err != nil {
// 		return 0, err
// 	}

// 	return res.(int64), nil
// }

// Lpush 将所有指定的值插入到存于 key 的列表的头部。
func (client *Client) Lpush(key string, val []byte) error {
	_, err := client.sendCommand("LPUSH", key, string(val))
	if err != nil {
		return err
	}

	return nil
}

// Lpop 移除并且返回 key 对应的 list 的第一个元素。
func (client *Client) Lpop(key string) ([]byte, error) {
	res, err := client.sendCommand("LPOP", key)
	if err != nil {
		return nil, err
	}

	return res.([]byte), nil
}

// Lpushx 只有当 key 已经存在并且存着一个 list 的时候，在这个 key 下面的 list 的头部插入 value。
func (client *Client) Lpushx(key string, val []byte) error {
	_, err := client.sendCommand("LPUSHX", key, string(val))
	if err != nil {
		return err
	}

	return nil
}

// Lindex 返回列表里的元素索引为index的值
func (client *Client) Lindex(key string, index int) ([]byte, error) {
	res, err := client.sendCommand("LINDEX", key, strconv.Itoa(int(index)))
	if err != nil {
		return nil, err
	}

	return res.([]byte), nil
}

// Llen 返回存储在 key 里的list的长度
func (client *Client) Llen(key string) (int64, error) {
	res, err := client.sendCommand("LLEN", key)
	if err != nil {
		return 0, err
	}

	return res.(int64), nil
}

// Lrange 返回存储在 key 的列表里指定范围内的元素。 start 和 end 偏移量都是基于0的下标，即list的第一个元素下标是0（list的表头），第二个元素下标是1，以此类推。偏移量也可以是负数，
func (client *Client) Lrange(key string, start, stop int) ([][]byte, error) {
	res, err := client.sendCommand("LRANGE", key, strconv.Itoa(start), strconv.Itoa(stop))
	if err != nil {
		return nil, err
	}
	return res.([][]byte), nil
}

// Lrem 从存于 key 的列表里移除前 count 次出现的值为 value 的元素。 这个 count 参数通过下面几种方式影响这个操作：
// count > 0: 从头往尾移除值为 value 的元素。
// count < 0: 从尾往头移除值为 value 的元素。
// count = 0: 移除所有值为 value 的元素。
func (client *Client) Lrem(key string, count int, value []byte) (int, error) {
	res, err := client.sendCommand("LREM", key, strconv.Itoa(count), string(value))
	if err != nil {
		return -1, err
	}
	return int(res.(int64)), nil
}

// Lset 设置 index 位置的list元素的值为 value。
func (client *Client) Lset(key string, index uint, val []byte) (err error) {
	_, err = client.sendCommand("LSET", key, strconv.Itoa(int(index)), string(val))
	return
}

// Ltrim 修剪(trim)一个已存在的 list，这样 list 就会只包含指定范围的指定元素
// start 和 stop 都是由0开始计数的 可以为负数
func (client *Client) Ltrim(key string, start, stop int) (err error) {
	_, err = client.sendCommand("LTRIM", key, strconv.Itoa(start), strconv.Itoa(stop))
	return
}

// Rpop 移除并返回存于 key 的 list 的最后一个元素。
func (client *Client) Rpop(key string) ([]byte, error) {
	res, err := client.sendCommand("RPOP", key)
	if err != nil {
		return nil, err
	}

	return res.([]byte), nil
}

// RpopLpush 原子性地返回并移除存储在 source 的列表的最后一个元素（列表尾部元素）， 并把该元素放入存储在 destination 的列表的第一个元素位置（列表头部）。
func (client *Client) RpopLpush(source, destination string) ([]byte, error) {
	res, err := client.sendCommand("RPOPLPUSH", source, destination)
	if err != nil {
		return nil, err
	}

	return res.([]byte), nil
}

// Rpush 向存于 key 的列表的尾部插入所有指定的值。
func (client *Client) Rpush(key string, val []byte) (err error) {
	_, err = client.sendCommand("RPUSH", key, string(val))
	return
}

// Rpushx 向存于 key 的列表的尾部插入所有指定的值。
func (client *Client) Rpushx(key string, val []byte) (err error) {
	_, err = client.sendCommand("RPUSHX", key, string(val))
	return
}

// Blpop BLPOP 是阻塞式列表的弹出原语
// @link http://redis.cn/commands/blpop.html
func (client *Client) Blpop(keys []string, timeoutSecs uint) (*string, []byte, error) {
	return client.bpop("BLPOP", keys, timeoutSecs)
}

// Brpop Brpop 是阻塞式列表的弹出原语
// @link http://redis.cn/commands/brpop.html
func (client *Client) Brpop(keys []string, timeoutSecs uint) (*string, []byte, error) {
	return client.bpop("BRPOP", keys, timeoutSecs)
}

// BrpopLpush BRPOPLPUSH 是 RPOPLPUSH 的阻塞版本
func (client *Client) BrpopLpush(source, destination string, timeoutSecs uint) ([]byte, error) {
	res, err := client.sendCommand("BRPOPLPUSH", source, destination, strconv.FormatUint(uint64(timeoutSecs), 10))
	if err != nil {
		return nil, err
	}
	// 如果达到 timeout 时限，会返回一个空的多批量回复(nil-reply)。
	if _, ok := res.([][]byte); ok {
		return nil, nil
	}

	return res.([]byte), nil
}

func (client *Client) bpop(cmd string, keys []string, timeoutSecs uint) (*string, []byte, error) {
	args := append(keys, strconv.FormatUint(uint64(timeoutSecs), 10))
	res, err := client.sendCommand(cmd, args...)
	if err != nil {
		return nil, nil, err
	}
	kv := res.([][]byte)
	// Check for timeout
	// 当没有元素的时候会弹出一个 nil 的多批量值，并且 timeout 过期。
	if len(kv) != 2 {
		return nil, nil, nil
	}

	k := string(kv[0])
	v := kv[1]

	return &k, v, nil
}
