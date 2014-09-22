package tunnel

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

const (
	commandPrefix = "["
	commandSuffix = "]"
)

type action struct {
	ArgsSt  interface{}
	Handler *func(interface{}) (string, error)
}

type Client struct {
	Id          string
	Secret      string
	Cert        tls.Certificate
	Remote      string
	RequestChan chan string
	actionMap   map[string]*action
	Timeout     time.Duration
}

func createClient() (*Client, error) {

}

func (c *Client) Server() {

}

func (c *Client) Request(action string, args interface{}) {

}

func (c *Client) Register(action string, argsSt interface{}, handler *func(interface{}) (string, error)) {
	_, ok = c.actionMap[action]
	if ok {
		log.Fatal("register action can not be the same")
	}
	c.actionMap[action] = &action{
		ArgsSt:  argsSt,
		Handler: handler,
	}
}

func (c *Client) handle(reply string) {
	command := receiveCommand(c, reply)
	message, err := command.Handler(command.Args)
	if err != nil {
		command.Response("500", err.Error())
		return
	}
	command.Response("200", message)
}

func (c *Client) setRequestFinished(id string) {

}

func (c *Client) server() {
	config := tls.Config{
		Certificates:       []tls.Certificate{c.Cert},
		InsecureSkipVerify: true,
	}
	conn, err := tls.Dial("tcp", c.Remote, &config)
	defer conn.Close()
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
	var hasPrefix, hasSuffix = false, false
	for {
		n, err = conn.Read(reply)
		if n == 0 || err != nil {
			return
		}
		// parse command and send to handle
		_reply = string(reply[:n])
		if strings.HasPrefix(_reply, commandPrefix) {
			_reply = strings.TrimPrefix(_reply, commandPrefix)
			hasPrefix = true
		}
		if strings.HasSuffix(_reply, commandSuffix) {
			_reply = strings.TrimSuffix(_reply, commandSuffix)
			hasSuffix = true
		}
		replyStr += _reply
		if hasPrefix && hasSuffix {
			go c.handle(replyStr)
			hasPrefix, hasSuffix, replyStr = false, false, ""
		}
	}
}
