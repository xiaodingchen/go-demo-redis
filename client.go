package redis

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"strconv"
	"strings"
)

var defaultAddr = "127.0.0.1:6379" // 默认地址

const defaultPoolSize = 5 // 默认最大连接池

// Client Redis客户端结构体
type Client struct {
	Addr        string
	Db          uint
	Password    string
	MaxPoolSize uint
	pool        chan net.Conn // 连接通道
}

// RedisError 定义错误类型
type RedisError string

func (err RedisError) Error() string {
	return "Redis Error: " + string(err)
}

var doesNotExist = RedisError("Key or index does not exist")

// NewRedisClient 创建一个Redis客户端实例
func NewRedisClient(addr string, db uint, pwd string, maxpoolsize uint) *Client {

	if addr == "" {
		addr = defaultAddr
	}
	if maxpoolsize == 0 {
		maxpoolsize = defaultPoolSize
	}
	client := &Client{Addr: addr, Db: db, Password: pwd, MaxPoolSize: maxpoolsize}
	// 设置连接池
	client.pool = make(chan net.Conn, client.MaxPoolSize)
	for i := uint(0); i < client.MaxPoolSize; i++ {
		client.pool <- nil
	}

	return client
}

// rawSend 执行Redis原生命令
func (client *Client) rawSend(c net.Conn, cmd []byte) (interface{}, error) {
	// 发送命令
	_, err := c.Write(cmd)
	if err != nil {
		return nil, err
	}
	// 开始redis服务端的响应
	reader := bufio.NewReader(c)
	data, err := readResponse(reader)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// sendCommand 发送Redis命令并执行
func (client *Client) sendCommand(cmd string, args ...string) (data interface{}, err error) {
	var b []byte
	c, err := client.popCon()
	if err != nil {
		goto END
	}
	// 将输入的命令转为byte
	b = commandBytes(cmd, args...)
	data, err = client.rawSend(c, b)
	// 处理EOF错误，EOF
	// 如果设置了连接超时Redis server会主动断开连接。客户端这边从一个超时的连接请求就会得到EOF错误
	// 按照道理来说应该在连接池初始化的时候做一个keeplive，但这里没有，只是新打开一个连接而已
	if err == io.EOF {
		c, err = client.openConnection()
		if err != nil {
			goto END
		}
		data, err = client.rawSend(c, b)
	}
END:
	client.pushCon(c)
	return
}

// 执行Redis命令，不处理返回结果
func writeRequest(writer io.Writer, cmd string, args ...string) error {
	b := commandBytes(cmd, args...)
	_, err := writer.Write(b)
	return err
}

// commandBytes 根据Redis通信协议（请求）组合命令
// *paramsNum\r\n$param1Len\r\nparam1\r\n$param2Len\r\nparam2\r\n
func commandBytes(cmd string, args ...string) []byte {
	var buffer bytes.Buffer
	// str := fmt.Sprintf("*%d\\r\\n$%d\\r\\n%s\\r\\n", len(args)+1, len(cmd), cmd)
	// buffer.WriteString(str)
	fmt.Fprintf(&buffer, "*%d\r\n$%d\r\n%s\r\n", len(args)+1, len(cmd), cmd)
	for _, s := range args {
		// buffer.WriteString(fmt.Sprintf("$%d\\r\\n%s\\r\\n", len(s), s))
		fmt.Fprintf(&buffer, "$%d\r\n%s\r\n", len(s), s)
	}
	return buffer.Bytes()
}

// readResponse 根据Redis通信协议读取服务端响应
func readResponse(reader *bufio.Reader) (interface{}, error) {
	var line string
	var err error
	// 读取数据，直到读到为止
	for {
		line, err = reader.ReadString('\n') // ReadString 这个方法会在读取到数据之后，会设置下次读取起始点为上一次读取之后的长度
		if len(line) == 0 || err != nil {
			return nil, err
		}
		line = strings.TrimSpace(line)
		if len(line) > 0 {
			break
		}
	}
	// 根据Redis通信协议（回复）处理Redis响应
	// 单行回复
	if line[0] == '+' {
		return strings.TrimSpace(line[1:]), nil
	}
	// 错误回复
	if line[0] == '-' {
		return nil, RedisError(strings.TrimSpace(line[1:]))
	}
	// 整数回复
	if line[0] == ':' {
		n, err := strconv.ParseInt(strings.TrimSpace(line[1:]), 10, 64)
		if err != nil {
			return nil, RedisError("Integer reply is not number")
		}

		return n, nil
	}

	// 多批量回复
	// 多条批量回复的第一个字节为 "*" ， 后跟一个字符串表示的整数值， 这个值记录了多条批量回复所包含的回复数量， 再后面是一个 CRLF 。
	// 多条批量回复，可以包含任意类型的回复
	if line[0] == '*' {
		size, err := strconv.Atoi(strings.TrimSpace(line[1:]))
		if err != nil {
			return nil, RedisError("MultiBulk reply expected a number")
		}

		if size <= 0 {
			return make([][]byte, 0), nil
		}

		res := make([][]byte, size)
		for i := 0; i < size; i++ {
			res[i], err = readBulk(reader, "")
			// 如果请求的值不存在，跳出循环，继续取值
			if err == doesNotExist {
				continue
			}

			if err != nil {
				return nil, err
			}

		}

		return res, err
	}

	return readBulk(reader, line)
}

func readBulk(reader *bufio.Reader, head string) ([]byte, error) {
	var err error
	var data []byte

	if head == "" {
		head, err = reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
	}
	switch head[0] {
	// 整数回复
	case ':':
		data = []byte(strings.TrimSpace(head[1:]))
	// 批量回复
	case '$':
		// 获取批量回复的长度，回复长度
		size, err := strconv.Atoi(strings.TrimSpace(head[1:]))
		if err != nil {
			return nil, err
		}
		// 如果被请求的值不存在， 那么批量回复会将特殊值 -1 用作回复的长度值
		if size == -1 {
			return nil, doesNotExist
		}
		// 设置此次批量回复的读取长度
		lr := io.LimitReader(reader, int64(size))
		// 读取指定长度的全部
		data, err = ioutil.ReadAll(lr)
		if err == nil {
			// 从指定长度的结束处开始读取，查看是否有错误
			_, err = reader.ReadString('\n')
		}
	default:
		return nil, RedisError("Expecting Prefix '$' or ':'")
	}

	return data, err
}
