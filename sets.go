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
