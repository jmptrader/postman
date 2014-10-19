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
	Handler func(interface{})
	proto   *Proto
}

// create command with action and args
func newCommand(proto *Proto, action string, args interface{}) *Command {
	return &Command{
		Id:     "-" + util.RandSeq(4),
		Action: action,
		Args:   args,
		proto:  proto,
	}
}

// parse request string to command struct
func receiveCommand(proto *Proto, command string) (c *Command, err error) {
	commandArr := strings.SplitN(command, commandSep, 3)
	c = &Command{
		Id:     commandArr[0],
		Action: commandArr[1],
		proto:  proto,
	}
	actionSt, ok := proto.actionMap[c.Action]
	if !ok {
		log.Printf("action %s not found in proto", c.Action)
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
