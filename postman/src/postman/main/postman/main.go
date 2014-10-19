package main

import (
	"postman/actions"
	"postman/client"
	"postman/mail"
)

func main() {
	tunnel := client.Postman.Tunnel

	// TODO register all actions
	tunnel.Register("exit", func() interface{} {
		return &actions.ExitMsg{}
	}, actions.Exit)

	tunnel.Register("helo", func() interface{} {
		return &actions.HeloMsg{}
	}, actions.Helo)

	tunnel.Register("authenticated", func() interface{} {
		return &actions.AuthenticatedMsg{}
	}, actions.Authenticated)

	tunnel.Register("sendMail", func() interface{} {
		return &mail.Mail{}
	}, actions.SendMail)

	// connect to middleware
	tunnel.Serve()
}
