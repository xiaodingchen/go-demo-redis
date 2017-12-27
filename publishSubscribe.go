package redis

import (
	"bufio"
	"io"
)

// 发布/订阅

// Message 接收到的消息
type Message struct {
	ChannelMatche string // 频道模式
	Channel       string // 具体频道
	Message       []byte // 具体消息
}

// Subscribe 订阅频道
func (client *Client) Subscribe(subscribe <-chan string, unscribe <-chan string, psubscribe <-chan string, punscribe <-chan string, messages chan<- Message) error {
	cmds := make(chan []string, 0)
	data := make(chan interface{}, 0)
	// 开个协程处理分析命令模式
	go func() {
		for {
			var channel string
			var cmd string
			select {
			case channel = <-subscribe:
				cmd = "SUBSCRIBE"
			case channel = <-unscribe:
				cmd = "UNSUBSCRIBE"
			case channel = <-psubscribe:
				cmd = "PSUBSCRIBE"
			case channel = <-punscribe:
				cmd = "PUNSUBSCRIBE"
			}

			if channel == "" {
				break
			} else {
				cmds <- []string{cmd, channel}
			}
		}

		close(cmds)
		close(data)
	}()

	// 把频道广播的值放到messages里
	go func() {
		for response := range data {
			db := response.([][]byte)
			messageType := string(db[0])
			switch messageType {
			case "message":
				channel, message := string(db[1]), db[2]
				messages <- Message{channel, channel, message}
			case "subscribe":
				// Ignore
			case "unsubscribe":
				// Ignore
			case "pmessage":
				channelMatched, channel, message := string(db[1]), string(db[2]), db[3]
				messages <- Message{channelMatched, channel, message}
			case "psubscribe":
				// Ignore
			case "punsubscribe":
				// Ignore

			default:
				// log.Printf("Unknown message '%s'", messageType)
			}
		}
	}()
	err := client.sendCommands(cmds, data)

	return err
}

func (client *Client) sendCommands(cmds <-chan []string, data chan<- interface{}) (err error) {
	c, err := client.popCon()
	var reader *bufio.Reader
	var pong interface{}
	var errs chan error
	var errClosed = false

	if err != nil {
		goto End
	}

	reader = bufio.NewReader(c)
	err = writeRequest(c, "PING")

	if err == io.EOF {
		c, err = client.openConnection()
		if err != nil {
			goto End
		}
		reader = bufio.NewReader(c)
	} else {
		pong, err = readResponse(reader)
		if pong != "PONG" {
			return RedisError("Unexpected response to PING.")
		}

		if err != nil {
			goto End
		}
	}

	errs = make(chan error) // 无缓存的通道，只要有进入就会阻塞
	// 开个协程处理命令
	go func() {
		for cmd := range cmds {
			err = writeRequest(c, cmd[0], cmd[1:]...)
			if err != nil {
				if !errClosed {
					errs <- err
				}

				break
			}
		}

		if !errClosed {
			errClosed = true
			close(errs)
		}
	}()

	// 处理返回数据
	go func() {
		for {
			response, err := readResponse(reader)
			if err != nil {
				if !errClosed {
					errs <- err
				}

				break
			}

			data <- response
		}

		if !errClosed {
			errClosed = true
			close(errs)
		}
	}()

	for e := range errs {
		err = e
	}

End:
	c.Close()
	client.pushCon(nil)
	return
}

// Publish 频道发布信息
func (client *Client) Publish(channel string, val []byte) error {
	_, err := client.sendCommand("PUBLISH", channel, string(val))
	return err
}
