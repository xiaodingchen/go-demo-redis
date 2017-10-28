package redis

// Save 执行一个同步保存操作，将当前 Redis 实例的所有数据快照(snapshot)以 RDB 文件的形式保存到硬盘。
func (client *Client) Save() error {
	_, err := client.sendCommand("SAVE")

	return err
}

// Bgsave 在后台一部保存当前数据库的数据到磁盘
func (client *Client) Bgsave() error {
	_, err := client.sendCommand("BGSAVE")
	return err
}

// Bgwriteaof 异步执行一个 AOF（AppendOnly File） 文件重写操作
func (client *Client) Bgwriteaof() error {
	_, err := client.sendCommand("BGWRITEAOF")
	return err
}

// Lastsave 返回最近一次 Redis 成功将数据保存到磁盘上的时间，以 UNIX 时间戳格式表示
func (client *Client) Lastsave() error {
	_, err := client.sendCommand("LASTSAVE")

	return err
}

// Flushall 清空所有数据库
func (client *Client) Flushall() (err error) {
	_, err = client.sendCommand("FLUSHALL")
	return
}

// Flushdb 清空当前数据库的数据
func (client *Client) Flushdb() (err error) {
	_, err = client.sendCommand("FLUSHDB")
	return
}

// Dbsize 返回当前数据库key的数量
func (client *Client) Dbsize() (int, error) {
	res, err := client.sendCommand("DBSIZE")
	if err != nil {
		return 0, err
	}

	return int(res.(int64)), nil
}
