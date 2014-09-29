package tunnel

import (
	"crypto/tls"
	"io"
	"log"
	"time"
)

const (
	commandPrefix = "DATA"
	commandSuffix = "END"
)

type Action struct {
	Instance func() interface{}
	Handler  func(*Client, interface{})
}

type Config struct {
	Id     string
	Conf   *tls.Config
	Remote string
	Secret string
}

type Client struct {
	config      Config
	RequestChan chan string
	actionMap   map[string]*Action
}

func (c *Client) Serve() {
	c.serve()
	<-time.After(time.Second * 10)
	c.Serve()
}

func (c *Client) Request(action string, args interface{}) {
	command := newCommand(c, action, args)
	for _, cmd := range []string{commandPrefix, command.String(), commandSuffix} {
		c.RequestChan <- cmd
	}
}

func (c *Client) Register(action string, instance func() interface{}, handler func(*Client, interface{})) {
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
	command.Handler(c, command.Args)
}

func (c *Client) serve() {
	conn, err := tls.Dial("tcp", c.config.Remote, c.config.Conf)
	if err != nil {
		log.Printf("client: dial: %s.\nReconnect will start after 10 seconds.", err)
		return
	}
	defer conn.Close()
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
			log.Printf("remote server: %s disconnect.\nReconnect will start after 10 seconds.", c.config.Remote)
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
