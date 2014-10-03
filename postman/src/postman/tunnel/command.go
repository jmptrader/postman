package tunnel

import (
	"encoding/json"
	"errors"
	"log"
	"strings"

	"postman/util"
)

const (
	commandSep = "|"
)

type Command struct {
	Id      string
	Action  string
	Args    interface{}
	Handler func(*Client, interface{})
	client  *Client
}

// create command with action and args
func newCommand(client *Client, action string, args interface{}) *Command {
	return &Command{
		Id:     "-" + util.RandSeq(16),
		Action: action,
		Args:   args,
		client: client,
	}
}

// parse request string to command struct
func receiveCommand(client *Client, command string) (c *Command, err error) {
	commandArr := strings.SplitN(command, commandSep, 3)
	c = &Command{
		Id:     commandArr[0],
		Action: commandArr[1],
		client: client,
	}
	actionSt, ok := client.actionMap[c.Action]
	if !ok {
		log.Printf("action %s not found in client", c.Action)
		err = errors.New("action not found")
		return
	}
	args := actionSt.Instance()
	err = json.Unmarshal([]byte(commandArr[2]), args)
	if err != nil {
		return
	}
	c.Args = args
	c.Handler = actionSt.Handler
	return
}

// parse struct to request buffer string
func (cm *Command) String() string {
	args, err := json.Marshal(cm.Args)
	if err != nil {
		log.Fatalf("parse interface to msg %s", err.Error())
	}
	return strings.Join([]string{cm.Id, cm.Action, string(args)}, commandSep)
}
