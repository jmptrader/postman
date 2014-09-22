package tunnel

import (
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
	Handler *func(interface{}) (string, error)
	client  *Client
}

// create command with action and args
func newCommand(client *Command, action string, args interface{}) *Command {
	return &Command{
		Id:     util.RandSeq(16),
		Action: action,
		Args:   args,
		client: &client,
	}
}

// parse request string to command struct
func receiveCommand(client *Command, command string) (c *Command) {
	commandArr := strings.SplitN(reply, commandSep, 3)
	c = &Command{
		Id:     commandArr[0],
		Action: commandArr[1],
		client: &client,
	}
	if c.Action == "response" {
		c.setRequestFinished(c.Id)
		return
	}
	actionSt, ok := client.actionMap[c.Action]
	if !ok {
		c.Response("404", "action "+c.Action+" not found in client")
		return
	}
	args := new(actionSt.ArgsSt)
	util.MsgDecode([]byte(commandArr[2]), args)
	c.Args = args
	c.Handler = actionSt.Handler
	return
}

// parse struct to request buffer string
func (cm *Command) String() string {
	args := util.MsgEncode(cm.Args)
	return strings.Join([3]string{cm.Id, cm.Action, args}, commandSep)
}

// client => [53ff21479560ce464d000001|dkim-query|{"domain": "open.jianxin.io"}]
// server => [53ff21479560ce464d000001|response|{"code": "200"}]
func (cm *Command) Response(code string, message string) {
	args := map[string]string{"code": code, "message": message}
	response := Command{
		Id:     cm.Id,
		Action: "response",
		Args:   args,
	}
	cm.client.Request(response.String())
}
