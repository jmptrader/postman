package tunnel

import (
	"log"
	"strings"

	"postman/util"
)

const (
	commandSep = "|"
)

// Todo: hanlder may need to be replace with interface{}
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
		Id:     util.RandSeq(16),
		Action: action,
		Args:   args,
		client: client,
	}
}

// parse request string to command struct
func receiveCommand(client *Client, command string) (c *Command) {
	commandArr := strings.SplitN(command, commandSep, 3)
	c = &Command{
		Id:     commandArr[0],
		Action: commandArr[1],
		client: client,
	}
	actionSt, ok := client.actionMap[c.Action]
	if !ok {
		log.Printf("action %s not found in client", c.Action)
		return
	}
	args := actionSt.Instance()
	util.MsgDecode([]byte(commandArr[2]), args)
	c.Args = args
	c.Handler = actionSt.Handler
	return
}

// parse struct to request buffer string
func (cm *Command) String() string {
	args, err := util.MsgEncode(cm.Args)
	if err != nil {
		log.Fatalf("parse interface to msg %s", err.Error())
	}
	return strings.Join([]string{cm.Id, cm.Action, string(args)}, commandSep)
}
