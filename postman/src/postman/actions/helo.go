package actions

import (
	"postman/client"
)

type HeloMsg struct {
	Key string `json:"auth_key"`
}

func Helo(args interface{}) {
	heloMsg := args.(*HeloMsg)
	client.Postman.Tunnel.Auth(heloMsg.Key)
}
