package actions

import (
	"postman/tunnel"
)

type HeloMsg struct {
	Key string `json:"auth_key"`
}

func Helo(c *tunnel.Client, args interface{}) {
	heloMsg := args.(*HeloMsg)
	result := c.Secret(heloMsg.Key)
	c.Request("auth", map[string]string{"result": result})
}
