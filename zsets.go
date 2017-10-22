package redis

import (
	"strconv"
)

// Zadd 将一个或多个成员元素及其分数值加入到有序集当中,或者更新已存在成员的分数
// 分数值可以是整数值或双精度浮点数。
func (client *Client) Zadd(key string, score float64, value []byte) (bool, error) {
	res, err := client.sendCommand("ZADD", key, strconv.FormatFloat(score, 'f', -1, 64), string(value))
	if err != nil {
		return false, err
	}

	return res.(int64) == 1, nil
}

// Zcard 获取有序集合的成员数
func (client *Client) Zcard(key string) (int, error) {
	res, err := client.sendCommand("ZCARD", key)
	if err != nil {
		return 0, err
	}

	return int(res.(int64)), nil
}

// Zcount 计算在有序集合中指定区间分数的成员数
func (client *Client) Zcount(key string, min, max float64) (int, error) {
	res, err := client.sendCommand("ZCOUNT", key, strconv.FormatFloat(min, 'f', -1, 64), strconv.FormatFloat(max, 'f', -1, 64))
	if err != nil {
		return 0, err
	}

	return int(res.(int64)), nil
}

// Zincrby 有序集合中对指定成员的分数加上增量 increment
func (client *Client) Zincrby(key string, increment float64, member string) (bool, error) {
	res, err := client.sendCommand("ZINCRBY", key, strconv.FormatFloat(increment, 'f', -1, 64), member)
	if err != nil {
		return false, err
	}

	return res.(int64) == 1, nil
}

// Zinterstore 计算给定的一个或多个有序集的交集并将结果集存储在新的有序集合 key 中
func (client *Client) Zinterstore(des string, num int, key ...string) (int, error) {
	args := []string{des, strconv.Itoa(num)}
	args = append(args, key...)
	res, err := client.sendCommand("ZINTERSTORE", args...)
	if err != nil {
		return 0, err
	}

	return int(res.(int64)), nil
}

// Zrange 通过索引区间返回有序集合成指定区间内的成员
func (client *Client) Zrange(key string, start, stop int) ([][]byte, error) {
	res, err := client.sendCommand("ZRANGE", key, strconv.Itoa(start), strconv.Itoa(stop))
	if err != nil {
		return nil, err
	}

	return res.([][]byte), nil
}

// Zrangebyscore 通过分数返回有序集合指定区间内的成员
func (client *Client) Zrangebyscore(key string, start float64, end float64) ([][]byte, error) {
	res, err := client.sendCommand("ZRANGEBYSCORE", key, strconv.FormatFloat(start, 'f', -1, 64), strconv.FormatFloat(end, 'f', -1, 64))
	if err != nil {
		return nil, err
	}

	return res.([][]byte), nil
}

// Zrank 有序集合中指定成员的索引
func (client *Client) Zrank(key, member string) (int, error) {
	res, err := client.sendCommand("ZRANK", key, member)

	if err != nil {
		return 0, nil
	}

	return int(res.(int64)), nil
}

// Zrem 移除有序集合中的一个或多个成员
func (client *Client) Zrem(key string, members ...string) (bool, error) {
	args := []string{key}
	args = append(args, members...)
	res, err := client.sendCommand("ZREM", args...)
	if err != nil {
		return false, err
	}

	return res.(int64) > 0, nil
}

// Zremrangebyrank
