package redis

import (
	"errors"
	"reflect"
	"strconv"
)

// Hset 设置 key 指定的哈希集中指定字段的值。
func (client *Client) Hset(key string, field string, value []byte) (err error) {
	_, err = client.sendCommand("HSET", key, field, string(value))
	return
}

// Hget 返回 key 指定的哈希集中该字段所关联的值
func (client *Client) Hget(key, field string) ([]byte, error) {
	res, err := client.sendCommand("HGET", key, field)
	if err != nil {
		return nil, err
	}

	return res.([]byte), nil
}

// Hexists 返回hash里面field是否存在
func (client *Client) Hexists(key, field string) (bool, error) {
	res, err := client.sendCommand("HEXISTS", key, field)
	if err != nil {
		return false, err
	}

	return res.(int64) == 1, nil
}

// Hdel 从 key 指定的哈希集中移除指定的域。在哈希集中不存在的域将被忽略。
func (client *Client) Hdel(key string, fields ...string) (bool, error) {
	args := []string{key}
	args = append(args, fields...)
	res, err := client.sendCommand("HDEL", args...)
	if err != nil {
		return false, nil
	}

	return res.(int64) > 0, nil
}

// Hgetall 返回 key 指定的哈希集中所有的字段和值
func (client *Client) Hgetall(key string) (map[string][]byte, error) {
	res, err := client.sendCommand("HGETALL", key)
	if err != nil {
		return nil, err
	}
	var data [][]byte
	var ok bool
	if data, ok = res.([][]byte); ok {
		var ret map[string][]byte
		ct := len(data)
		ret = make(map[string][]byte)
		for i := 0; i < ct; i++ {
			if i%2 == 0 {
				ret[string(data[i])] = data[i+1]
			}
		}

		return ret, nil
	}
	return nil, RedisError("Hgetall %s error: reply type error")
}

// Hincrby 增加 key 指定的哈希集中指定字段的数值。如果 key 不存在，会创建一个新的哈希集并与 key 关联。如果字段不存在，则字段的值在该操作执行前被设置为 0
func (client *Client) Hincrby(key, field string, num int64) (int64, error) {
	res, err := client.sendCommand("HINCRBY", key, field, strconv.FormatInt(num, 10))
	if err != nil {
		return 0, err
	}

	return res.(int64), nil
}

// Hkeys 返回 key 指定的哈希集中所有字段的名字。
func (client *Client) Hkeys(key string) ([]string, error) {
	res, err := client.sendCommand("HKEYS", key)
	if err != nil {
		return nil, err
	}

	data := res.([][]byte)
	if data == nil || len(data) == 0 {
		return nil, RedisError("Key `" + key + "` does not exist")
	}
	ret := make([]string, len(data))
	for i, v := range data {
		ret[i] = string(v)
	}

	return ret, nil
}

// Hvals 返回 key 指定的哈希集中所有字段的值。
func (client *Client) Hvals(key string) ([][]byte, error) {
	res, err := client.sendCommand("HVALS", key)

	if err != nil {
		return nil, err
	}
	return res.([][]byte), nil
}

// Hlen 返回 key 指定的哈希集包含的字段的数量。
func (client *Client) Hlen(key string) (int, error) {
	res, err := client.sendCommand("HLEN", key)
	if err != nil {
		return 0, err
	}

	return int(res.(int64)), nil
}

// Hmset 设置 key 指定的哈希集中指定字段的值。
func (client *Client) Hmset(key string, mapping interface{}) error {
	args := []string{key}
	err := containerToString(reflect.ValueOf(mapping), &args)
	if err != nil {
		return err
	}
	_, err = client.sendCommand("HMSET", args...)
	if err != nil {
		return err
	}
	return nil
}

// Hmget 返回 key 指定的哈希集中指定字段的值。
func (client *Client) Hmget(key string, fields ...string) ([][]byte, error) {
	args := []string{key}
	args = append(args, fields...)
	res, err := client.sendCommand("HMGET", args...)
	if err != nil {
		return nil, err
	}

	return res.([][]byte), nil
}

// Hsetnx 只在 key 指定的哈希集中不存在指定的字段时，设置字段的值。
// 如果 key 指定的哈希集不存在，会创建一个新的哈希集并与 key 关联。如果字段已存在，该操作无效果。
func (client *Client) Hsetnx(key, field string, value []byte) (bool, error) {
	res, err := client.sendCommand("HSETNX", key, field, string(value))
	if err != nil {
		return false, err
	}

	return res.(int64) == 1, nil
}

// Hstrlen 返回hash指定field的value的字符串长度，如果hash或者field不存在，返回0.
func (client *Client) Hstrlen(key, field string) (int, error) {
	res, err := client.sendCommand("HSTRLEN", key, field)
	if err != nil {
		return 0, nil
	}

	return int(res.(int64)), nil
}

func containerToString(val reflect.Value, args *[]string) error {
	switch v := val; v.Kind() {
	case reflect.Ptr:
		return containerToString(reflect.Indirect(v), args)
	case reflect.Interface:
		return containerToString(v.Elem(), args)
	case reflect.Map:
		if v.Type().Key().Kind() != reflect.String {
			return errors.New("Unsupported type - map key must be a string")
		}
		// 将map变成[]string
		for _, k := range v.MapKeys() {
			*args = append(*args, k.String())
			s, err := valueToString(v.MapIndex(k))
			if err != nil {
				return err
			}
			*args = append(*args, s)
		}
	case reflect.Struct:
		st := v.Type()
		// 将struct变成[]string
		for i := 0; i < st.NumField(); i++ {
			ft := st.FieldByIndex([]int{i})
			*args = append(*args, ft.Name)
			fv, err := valueToString(v.FieldByIndex([]int{i}))
			if err != nil {
				return err
			}
			*args = append(*args, fv)
		}
	}

	return nil
}

func valueToString(v reflect.Value) (string, error) {
	if !v.IsValid() {
		return "null", nil
	}
	switch v.Kind() {
	case reflect.Ptr:
		return valueToString(reflect.Indirect(v))
	case reflect.Interface:
		return valueToString(v.Elem())
	case reflect.Bool:
		x := v.Bool()
		if x {
			return "true", nil
		}

		return "false", nil
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8, reflect.Uintptr:
		return strconv.FormatUint(v.Uint(), 10), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10), nil
	case reflect.UnsafePointer:
		return strconv.FormatUint(uint64(v.UnsafeAddr()), 10), nil
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'g', -1, 64), nil
	case reflect.String:
		return v.String(), nil
	case reflect.Slice: // 处理切片，这里只处理整数类型的切片，所以在client调用的时候请把字符串方式的转为[]byte类型
		typ := v.Type()
		if typ.Elem().Kind() == reflect.Uint || typ.Elem().Kind() == reflect.Uint8 || typ.Elem().Kind() == reflect.Uint16 || typ.Elem().Kind() == reflect.Uint32 || typ.Elem().Kind() == reflect.Uint64 || typ.Elem().Kind() == reflect.Uintptr {
			if v.Len() > 0 {
				if v.Index(0).OverflowUint(257) {
					return string(v.Interface().([]byte)), nil
				}
			}
		}
	}

	return "", errors.New("Unsupported type")
}
