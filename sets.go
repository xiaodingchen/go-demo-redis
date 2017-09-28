package redis

// Sadd 添加一个或多个指定的member元素到集合的 key中
func (client *Client) Sadd(key string, members ...[]byte) (int, error) {
	args := []string{key}
	for _, v := range members {
		args = append(args, string(v))
	}
	res, err := client.sendCommand("SADD", args...)
	if err != nil {
		return 0, err
	}

	return int(res.(int64)), nil
}

// Scard 返回集合存储的key的基数 (集合元素的数量).
func (client *Client) Scard(key string) (int, error) {
	res, err := client.sendCommand("SCARD", key)
	if err != nil {
		return 0, err
	}

	return int(res.(int64)), nil
}

// Sdiff 返回一个集合与给定集合的差集的元素
func (client *Client) Sdiff(keys ...string) ([][]byte, error) {
	res, err := client.sendCommand("SDIFF", keys...)

	if err != nil {
		return nil, err
	}

	return res.([][]byte), nil
}

// Sdiffstore 该命令类似于 SDIFF, 不同之处在于该命令不返回结果集，而是将结果存放在destination集合中.
func (client *Client) Sdiffstore(destination string, keys ...string) (int, error) {
	args := []string{destination}
	args = append(args, keys...)
	res, err := client.sendCommand("SDIFFSTORE", args...)
	if err != nil {
		return 0, err
	}

	return int(res.(int64)), nil
}

// Sinter 返回指定所有的集合的成员的交集.
func (client *Client) Sinter(keys ...string) ([][]byte, error) {
	res, err := client.sendCommand("SINTER", keys...)
	if err != nil {
		return nil, err
	}

	return res.([][]byte), nil
}

// Sinterstore 这个命令与SINTER命令类似, 但是它并不是直接返回结果集,而是将结果保存在 destination集合中
func (client *Client) Sinterstore(destination string, keys ...string) (int, error) {
	args := []string{destination}
	args = append(args, keys...)
	res, err := client.sendCommand("SINTERSTORE", args...)
	if err != nil {
		return 0, err
	}

	return int(res.(int64)), nil
}

// Sismember 查看key中是否存在指定的member
func (client *Client) Sismember(key, member string) (bool, error) {
	res, err := client.sendCommand("SISMEMBER", key, member)
	if err != nil {
		return false, err
	}

	return res.(int64) == 1, nil
}

// Smembers 返回key集合所有的元素
func (client *Client) Smembers(key string) ([][]byte, error) {
	res, err := client.sendCommand("SMEMBERS", key)
	if err != nil {
		return nil, err
	}

	return res.([][]byte), nil
}

// Smove 将一个元素从一个集合移动到另一个集合
func (client *Client) Smove(src, dest, member string) (bool, error) {
	res, err := client.sendCommand("SMOVE", src, dest, member)
	if err != nil {
		return false, nil
	}

	return res.(int64) == 1, nil
}

// Spop 从集合中取出一个元素（随机）
func (client *Client) Spop(key string) ([]byte, error) {
	res, err := client.sendCommand("SPOP", key)
	if err != nil {
		return nil, err
	}

	if res == nil {
		return nil, RedisError("Spop failed")
	}

	data := res.([]byte)
	return data, nil
}

// Srandmember 仅提供key参数,随机返回key集合中的一个元素
func (client *Client) Srandmember(key string) ([]byte, error) {
	res, err := client.sendCommand("SRANDMEMBER", key)
	if err != nil {
		return nil, err
	}

	return res.([]byte), nil
}

// Srem 在key集合中移除指定的元素. 如果指定的元素不是key集合中的元素则忽略 如果key集合不存在则被视为一个空的集合，该命令返回0.
func (client *Client) Srem(key string, members ...string) (bool, error) {
	args := []string{key}
	args = append(args, members...)
	res, err := client.sendCommand("SREM", args...)
	if err != nil {
		return false, err
	}

	return res.(int64) == 1, nil
}

// Sunion 返回给定的多个集合的并集中的所有成员.
func (client *Client) Sunion(keys ...string) ([][]byte, error) {
	res, err := client.sendCommand("SUNION", keys...)
	if err != nil {
		return nil, err
	}

	return res.([][]byte), nil
}

// Sunionstore 该命令作用类似于SUNION命令,不同的是它并不返回结果集,而是将结果存储在destination集合中.
func (client *Client) Sunionstore(destination string, keys ...string) (int, error) {
	args := []string{destination}
	args = append(args, keys...)
	res, err := client.sendCommand("SUNIONSTORE", args...)
	if err != nil {
		return 0, nil
	}

	return int(res.(int64)), nil
}
