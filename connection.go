package redis

import (
	"fmt"
	"net"
	"strconv"
)

// 实现连接池
// popCon 从连接池中取出一个连接
func (client *Client) popCon() (net.Conn, error) {
	c := <-client.pool
	if c == nil {
		return client.openConnection()
	}

	return c, nil
}

// pushCon 向连接池中推送一个连接
func (client *Client) pushCon(c net.Conn) {
	client.pool <- c
}

// openConnection 实现Redis连接
func (client *Client) openConnection() (c net.Conn, err error) {
	c, err = net.Dial("tcp", client.Addr)
	if err != nil {
		return
	}
	// 开始验证密码
	if client.Password != "" {
		cmd := fmt.Sprintf("AUTH %s\r\n", client.Password)
		_, err = client.rawSend(c, []byte(cmd))
		if err != nil {
			return
		}
	}

	if client.Db != 0 {
		cmd := fmt.Sprintf("SELECT %d\r\n", client.Db)
		_, err = client.rawSend(c, []byte(cmd))
		if err != nil {
			return
		}
	}

	return
}

// Auth redis密码验证
func (client *Client) Auth(password string) error {
	_, err := client.sendCommand("AUTH", password)
	if err != nil {
		return err
	}

	return nil
}

// Echo 原样输出
func (client *Client) Echo(val []byte) ([]byte, error) {
	data, err := client.sendCommand("ECHO", string(val))
	if err != nil {
		return nil, err
	}

	return data.([]byte), nil
}

// Ping 如果后面没有参数时返回PONG，否则会返回后面带的参数
func (client *Client) Ping() error {
	res, err := client.sendCommand("PING")
	if err != nil {
		return err
	}

	if res != "PONG" {
		return RedisError("PING error.")
	}

	return nil
}

// Quit 请求服务器关闭连接
func (client *Client) Quit() error {
	_, err := client.sendCommand("QUIT")
	if err != nil {
		return err
	}

	return nil
}

// Select 选择一个数据库，下标值从0开始，一个新连接默认连接的数据库是DB0。
func (client *Client) Select(db uint) error {
	_, err := client.sendCommand("SELECT", strconv.FormatUint(uint64(db), 10))
	if err != nil {
		return err
	}

	return nil
}

// ExecCmd 执行指定的命令
func (client *Client) ExecCmd(cmd string, args ...string) (interface{}, error) {
	res, err := client.sendCommand(cmd, args...)
	if err != nil {
		return nil, err
	}

	return res, nil
}
