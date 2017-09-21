package redis

import (
	"strconv"
)

// Append 如果 key 已经存在，并且值为字符串，那么这个命令会把 value 追加到原来值（value）的结尾。 如果 key 不存在，那么它将首先创建一个空字符串的key，再执行追加操作
func (client *Client) Append(key string, value []byte) (err error) {
	_, err = client.sendCommand("APPEND", string(value))
	return
}

// Decr 对key对应的数字做减1操作
func (client *Client) Decr(key string) (int64, error) {
	data, err := client.sendCommand("DECR", key)
	if err != nil {
		return -1, err
	}

	return data.(int64), nil
}

// DecrBy 将key对应的数字减decrement
func (client *Client) DecrBy(key string, decrement int64) (int64, error) {
	data, err := client.sendCommand("DECRBY", key, strconv.FormatInt(decrement, 10))
	if err != nil {
		return -1, err
	}

	return data.(int64), nil
}

// Get 返回key的value。如果key不存在，返回特殊值nil。如果key的value不是string，就返回错误，因为GET只处理string类型的values。
func (client *Client) Get(key string) ([]byte, error) {
	data, _ := client.sendCommand("GET", key)
	if data == nil {
		return nil, RedisError("Key `" + key + "` does not exist")
	}

	return data.([]byte), nil
}

// GetSet 自动将key对应到value并且返回原来key对应的value。如果key存在但是对应的value不是字符串，就返回错误。
func (client *Client) GetSet(key string, val []byte) ([]byte, error) {
	data, err := client.sendCommand("GETSET", key, string(val))
	if err != nil {
		return nil, err
	}

	return data.([]byte), nil
}

// Incr 对存储在指定key的数值执行原子的加1操作
func (client *Client) Incr(key string) (int64, error) {
	data, err := client.sendCommand("INCR", key)
	if err != nil {
		return -1, err
	}

	return data.(int64), nil
}

// IncrBy 将key对应的数字加decrement
func (client *Client) IncrBy(key string, decrement int64) (int64, error) {
	data, err := client.sendCommand("INCRBY", key, strconv.FormatInt(decrement, 10))
	if err != nil {
		return -1, err
	}

	return data.(int64), nil
}

// Mget 返回所有指定的key的value
func (client *Client) Mget(keys ...string) ([][]byte, error) {
	data, err := client.sendCommand("MGET", keys...)
	if err != nil {
		return nil, err
	}

	return data.([][]byte), nil
}

// Mset 对应给定的keys到他们相应的values上。MSET会用新的value替换已经存在的value，就像普通的SET命令一样。
func (client *Client) Mset(mapping map[string][]byte) (err error) {
	args := make([]string, len(mapping)*2)
	i := 0
	for k, v := range mapping {
		args[i] = k
		args[i+1] = string(v)
		i += 2
	}

	_, err = client.sendCommand("MSET", args...)
	return
}

// Msetnx 对应给定的keys到他们相应的values上。只要有一个key已经存在，MSETNX一个操作都不会执行。
func (client *Client) Msetnx(mapping map[string][]byte) (flag bool, err error) {
	args := make([]string, len(mapping)*2)
	i := 0
	flag = false
	for k, v := range mapping {
		args[i] = k
		args[i+1] = string(v)
		i += 2
	}

	res, err := client.sendCommand("MSETNX", args...)
	if err != nil {
		return
	}
	if data, ok := res.(int64); ok {
		flag = (data == 1)
		return
	}

	return
}

// Set 将键key设定为指定的“字符串”值。如果	key	已经保存了一个值，那么这个操作会直接覆盖原来的值，并且忽略原始类型。
func (client *Client) Set(key string, val []byte) error {
	_, err := client.sendCommand("SET", key, string(val))

	if err != nil {
		return err
	}

	return nil
}

// Setex 设置key对应字符串value，并且设置key在给定的seconds时间之后超时过期。
func (client *Client) Setex(key string, seconds int64, val []byte) error {
	_, err := client.sendCommand("SETEX", strconv.FormatInt(seconds, 10), string(val))
	if err != nil {
		return err
	}

	return nil
}

// Setnx 将key设置值为value，如果key不存在，这种情况下等同SET命令。 当key存在时，什么也不做。
func (client *Client) Setnx(key string, val []byte) (bool, error) {
	res, err := client.sendCommand("SETNX", key, string(val))
	if err != nil {
		return false, err
	}

	if data, ok := res.(int64); ok {
		return data == 1, nil
	}

	return false, RedisError("Unexpected reply to MSETNX")
}
