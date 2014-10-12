package tunnel

import (
	"bufio"
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"postman/store"
)

const COMMAND_KEY_PREFIX = "cmd:"

type Action struct {
	Instance func() interface{}
	Handler  func(*Client, interface{})
}

type Config struct {
	Conf   *tls.Config
	Store  store.Store
	Remote string
	Secret string
}

type Client struct {
	config        Config
	RequestChan   chan interface{}
	authBlockChan chan bool
	actionMap     map[string]*Action
	online        bool
	buf           *bufio.ReadWriter
	conn          *tls.Conn
	name          *tls.Conn
}

func (c *Client) Serve() {
	c.online = true
	c.serve()
	if c.online {
		<-time.After(time.Second * 10)
		c.Serve()
	}
}

// send auth command to server
// exit if error meet
func (c *Client) Auth(str string) {
	hasher := md5.New()
	hasher.Write([]byte(c.config.Secret + str))
	cmd := newCommand(c, "auth", map[string]string{
		"result": hex.EncodeToString(hasher.Sum(nil)),
	})
	err := c.sendCmd(cmd.String())
	if err != nil {
		log.Fatalf("client: sende auth command %s", err.Error())
	}
}

func (c *Client) SetAuthenticated() {
	c.authBlockChan <- true
}

// send command to remote server
func (c *Client) Request(action string, args interface{}) {
	command := newCommand(c, action, args)
	c.RequestChan <- command
}

// register action for client
func (c *Client) Register(action string, instance func() interface{}, handler func(*Client, interface{})) {
	_, ok := c.actionMap[action]
	if ok {
		log.Fatalf("register action %s can not be the same", action)
	}
	c.actionMap[action] = &Action{
		Instance: instance,
		Handler:  handler,
	}
}

// handle request content
func (c *Client) handle(reply string) {
	command, err := receiveCommand(c, reply)
	if err != nil {
		if os.Getenv("POSTMAN_DEBUG_MODE") == "true" {
			log.Print(err)
		}
		return
	}
	command.Handler(c, command.Args)
}

// read buffer from server
func (c *Client) handleConn() {
	for {
		reply, err := c.buf.ReadString('\n')
		if err == io.EOF {
			log.Printf("\033[1;33;40mremote server: %s disconnect.\033[m\r\nReconnect will start after 10 seconds.", c.config.Remote)
			return
		}
		if !c.online {
			return
		}
		if err != nil {
			log.Printf("client: read buffer %s", err.Error())
		}
		if strings.HasPrefix(reply, "-") {
			continue
		}
		// parse command and send to handle
		if os.Getenv("POSTMAN_DEBUG_MODE") == "true" {
			log.Print("RECEIVE: ", reply)
		}
		go c.handle(reply)
	}
}

// send command string to server
func (c *Client) sendCmd(cmd string) error {
	if os.Getenv("POSTMAN_DEBUG_MODE") == "true" {
		log.Print("SEND: ", cmd)
	}
	c.buf.Write([]byte(cmd))
	c.buf.WriteByte('\n')
	return c.buf.Flush()
}

// send command to server
func (c *Client) handleReq() {
	// block until authenticated
	<-c.authBlockChan
	close(c.authBlockChan)
	for command := range c.RequestChan {
		var cmd, cmdId string
		// receive command via chan
		cmd, ok := command.(string)
		if ok {
			bytes := []byte(cmd)
			cmdId = string(bytes[0:5])
		} else {
			cmdSt, _ := command.(*Command)
			cmd, cmdId = cmdSt.String(), cmdSt.Id
			c.config.Store.Set(COMMAND_KEY_PREFIX+cmdId, cmd)
		}
		// then send it
		err := c.sendCmd(cmd)
		if err != nil {
			log.Printf("client: send %s: %s", command, err)
			// TODO: resent after 1 min
			return
		}
		c.config.Store.Destroy(COMMAND_KEY_PREFIX + cmdId)
	}
}

// close conn from client
func (c *Client) Close() {
	c.online = false
	close(c.RequestChan)
	c.conn.Close()
}

// start tls client and handshake
func (c *Client) serve() {
	conn, err := tls.Dial("tcp", c.config.Remote, c.config.Conf)
	if err != nil {
		log.Printf("\033[1;33;40mclient: %s.\033[m\r\nReconnect will start after 10 seconds.", err)
		return
	}
	err = conn.Handshake()
	if err != nil {
		log.Printf("\033[1;33;40mclient handshake: %s.\033[m", err)
		return
	}
	log.Println("client: connected to: ", conn.RemoteAddr())
	defer conn.Close()
	c.conn = conn
	br := bufio.NewReader(conn)
	bw := bufio.NewWriter(conn)
	c.buf = bufio.NewReadWriter(br, bw)
	go c.handleReq()
	go func() {
		// resend all fail request
		for _, key := range c.config.Store.Keys(COMMAND_KEY_PREFIX) {
			cmd, ok := c.config.Store.Get(key)
			if ok {
				c.RequestChan <- cmd
			}
		}
	}()
	c.handleConn()
}
