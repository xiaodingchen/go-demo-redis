package redis

import (
	"testing"
	"time"
)

func TestSubscribe(t *testing.T) {
	subscribe := make(chan string, 0)
	unsubscribe := make(chan string, 0)
	psubscribe := make(chan string, 0)
	punsubscribe := make(chan string, 0)
	messages := make(chan Message, 0)
	defer func() {
		close(subscribe)
		close(unsubscribe)
		close(psubscribe)
		close(punsubscribe)
		close(messages)
	}()

	go func() {
		t.Logf("start 1 %v", time.Now().UnixNano())
		err := client.Subscribe(subscribe, unsubscribe, psubscribe, punsubscribe, messages)
		if err != nil {
			t.Fatal("Subscribed failed", err)
		}
		t.Log("end 1")
	}()
	t.Log("start test")
	subscribe <- "news"
	data := []byte("hh")
	quit := make(chan bool, 0)
	defer close(quit)

	go func() {
		t.Log("start 2")
		tick := time.Tick(10 * 1000 * 1000)     // 10ms
		timeout := time.Tick(100 * 1000 * 1000) // 100ms
		for {
			select {
			case <-quit:
				return
			case <-timeout:
				t.Fatal("test timeout")
			case <-tick:
				if err := client.Publish("news", data); err != nil {
					t.Fatal("Publish faild", err)
				}

			}
		}
		t.Log("end 2")
	}()
	t.Logf("end test %v", time.Now().UnixNano())
	msg := <-messages

	quit <- true

	if msg.Channel != "news" {
		t.Fatal("Unexpected channel name")
	}
	if string(msg.Message) != string(data) {
		t.Fatalf("Expected %s but got %s", string(data), string(msg.Message))
	}
}
