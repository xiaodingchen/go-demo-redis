package redis

import (
	"bytes"
	"strconv"
)

// Del 删除指定的一批keys，如果删除中的某些key不存在，则直接忽略。
func (client *Client) Del(keys ...string) error {
	_, err := client.sendCommand("DEL", keys...)
	if err != nil {
		return err
	}

	return nil
}

// Exists 返回key是否存在。
func (client *Client) Exists(key string) (bool, error) {
	res, err := client.sendCommand("EXISTS", key)
	if err != nil {
		return false, err
	}

	return res.(int64) == 1, nil
}

// Expire 设置key的过期时间,超过时间后，将会自动删除该key。
func (client *Client) Expire(key string, s uint64) (bool, error) {
	res, err := client.sendCommand("EXPIRE", key, strconv.FormatUint(s, 10))
	if err != nil {
		return false, err
	}

	return res.(int64) == 1, nil
}

// Expireat EXPIREAT 的作用和 EXPIRE类似，都用于为 key 设置生存时间。不同在于 EXPIREAT 命令接受的时间参数是 UNIX 时间戳 Unix timestamp 。
func (client *Client) Expireat(key string, timestamp uint64) (bool, error) {
	res, err := client.sendCommand("EXPIREAT", key, strconv.FormatUint(timestamp, 10))
	if err != nil {
		return false, err
	}

	return res.(int64) == 1, nil
}

// Keys 查找所有符合给定模式pattern（正则表达式）的 key 。
func (client *Client) Keys(pattern string) ([]string, error) {
	res, err := client.sendCommand("KEYS", pattern)
	if err != nil {
		return nil, err
	}
	var ok bool
	var keydata [][]byte
	if keydata, ok = res.([][]byte); ok {
		// 只是用来做个判断，keydata已赋值
	} else {
		keydata = bytes.Fields(res.([]byte))
	}
	ret := make([]string, len(keydata))
	for i, key := range keydata {
		ret[i] = string(key)
	}

	return ret, nil
}

// Move 将当前数据库的 key 移动到给定的数据库 db 当中。
func (client *Client) Move(key string, db uint) (bool, error) {
	res, err := client.sendCommand("MOVE", key, strconv.Itoa(int(db)))
	if err != nil {
		return false, err
	}

	return res.(int64) == 1, nil
}

// RandomKey 从当前数据库返回一个随机的key。
func (client *Client) RandomKey() (string, error) {
	res, err := client.sendCommand("RANDOMKEY")
	if err != nil {
		return "", err
	}

	return string(res.([]byte)), nil
}

// Rename 将key重命名为newkey，如果key与newkey相同，将返回一个错误。如果newkey已经存在，则值将被覆盖。
func (client *Client) Rename(key, newkey string) error {
	_, err := client.sendCommand("RENAME", key, newkey)
	if err != nil {
		return err
	}

	return nil
}

// Renamenx 当且仅当 newkey 不存在时，将 key 改名为 newkey 。当 key 不存在时，返回一个错误。
func (client *Client) Renamenx(key, newkey string) error {
	_, err := client.sendCommand("RENAMENX", key, newkey)
	if err != nil {
		return err
	}

	return nil
}

// TTL 返回key剩余的过期时间。 这种反射能力允许Redis客户端检查指定key在数据集里面剩余的有效期。
// 在Redis 2.6和之前版本，如果key不存在或者已过期时server返回-1。
// 从Redis2.8开始 如果key不存在或者已过期，server返回 -2 如果key存在并且没有设置过期时间（永久有效），server返回 -1 。
// 这里如果server返回错误，此方法就返回-9，-9代表server执行命令出错
// 使用者可根据实际返回值和server版本进行业务操作
// 同时也可以判断error值是否为空来判断server端返回的情况
func (client *Client) TTL(key string) (int64, error) {
	res, err := client.sendCommand("TTL", key)
	if err != nil {
		return -9, err
	}

	return res.(int64), nil
}

// Type 返回key所存储的value的数据结构类型，它可以返回string, list, set, zset 和 hash等不同的类型。
func (client *Client) Type(key string) (string, error) {
	res, err := client.sendCommand("TYPE", key)
	if err != nil {
		return "", err
	}

	return res.(string), nil
}
