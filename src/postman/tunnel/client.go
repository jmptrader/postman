package tunnel

// client => [53ff21479560ce464d000001|dkim-query|{"domain": "open.jianxin.io"}]
// server => [53ff21479560ce464d000001|response|{"code": "200"}]
// TODO:
// client 发出请求后，请求内容不销毁，放在内存/文件中，等待收到 response 后再做销毁
// 否则超时后进行重新发送
// command 主要负责封装请求的解析和格式化
// 考虑使用钩子/反射 将请求与 action 绑定，并设定请求 args 类型并自定解析作为参数传递

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"log"
	"strings"
	"time"

	"postman/store"
)

type Client struct {
	Id          string
	Secret      string
	Cert        tls.Certificate
	Remote      string
	RequestChan chan string
	Timeout     time.Duration
}

func createClient() (*Client, error) {

}

func (c *Client) handle(reply string) {

}

func (c *Client) Request(command Command) {

}

func (c *Client) Server() {
	defer conn.Close()
	config := tls.Config{
		Certificates:       []tls.Certificate{c.Cert},
		InsecureSkipVerify: true,
	}
	conn, err := tls.Dial("tcp", c.Remote, &config)
	if err != nil {
		log.Fatalf("client: dial: %s", err)
	}
	log.Println("client: connected to: ", conn.RemoteAddr())
	go func() {
		var command string
		for {
			// receive command via chan
			// then send it
			command = <-c.RequestChan
			_, err := io.WriteString(conn, command)
			if err != nil {
				log.Printf("client: send %s: %s", command, err)
			}
		}
	}()
	reply := make([]byte, 300)
	var replyStr, _reply string
	for {
		n, err = conn.Read(reply)
		if n == 0 || err != nil {
			return
		}
		_reply = string(reply[:n])
		c.handle()
	}
}
