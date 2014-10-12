package actions

import (
	"postman/tunnel"
)

type HeloMsg struct {
	Key string `json:"auth_key"`
}

func Helo(c *tunnel.Client, args interface{}) {
	heloMsg := args.(*HeloMsg)
	c.Auth(heloMsg.Key)
}
