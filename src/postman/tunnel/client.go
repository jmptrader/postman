package tunnel

import (
	"crypto/tls"
	// "crypto/x509"
	"io"
	"log"
	"time"

	"postman/store"
)

const (
	commandPrefix = "DATA"
	commandSuffix = "END"
)

type Action struct {
	Instance func() interface{}
	Handler  func(interface{})
}

type Client struct {
	Id          string
	Secret      string
	Cert        tls.Certificate
	Remote      string
	RequestChan chan string
	actionMap   map[string]*Action
	store       *store.Store
}

func createClient() (*Client, error) {
	return &Client{
		RequestChan: make(chan string, 10),
		actionMap:   map[string]*Action{},
	}, nil
}

func (c *Client) Server() {

}

func (c *Client) Request(action string, args interface{}) {

}

func (c *Client) Register(action string, instance func() interface{}, handler func(interface{})) {
	_, ok := c.actionMap[action]
	if ok {
		log.Fatal("register action can not be the same")
	}
	c.actionMap[action] = &Action{
		Instance: instance,
		Handler:  handler,
	}
}

func (c *Client) handle(reply string) {
	command := receiveCommand(c, reply)
	command.Handler(command.Args)
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
			// TODO: Here need to add DATA and END
			_, err := io.WriteString(conn, command)
			if err != nil {
				log.Printf("client: send %s: %s", command, err)
			}
		}
	}()
	reply := make([]byte, 300)
	var replyStr, _reply string
	var hasPrefix = false
	for {
		n, err := conn.Read(reply)
		if n == 0 || err != nil {
			return
		}
		// parse command and send to handle
		_reply = string(reply[:n])
		if _reply == commandPrefix {
			hasPrefix = true
			continue
		}
		if hasPrefix && _reply == commandSuffix {
			go c.handle(replyStr)
			hasPrefix, replyStr = false, ""
			continue
		}
		replyStr += _reply
	}
}
