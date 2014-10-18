package main

import (
	"postman/actions"
	"postman/store"
	"postman/tunnel"
)

type Postman struct {
	Tunnel tunnel.Tunnel
	Store  store.Store
}

func main() {
	st := store.New(dbDir, config.StoreSecret)
	postman := Postman{
		Store:  st,
		Tunnel: tunnel.New(tunnelConfig(st)),
	}
	// TODO register all actions
	postman.Tunnel.Register("exit", func() interface{} {
		return &actions.ExitMsg{}
	}, actions.Exit)

	postman.Tunnel.Register("helo", func() interface{} {
		return &actions.HeloMsg{}
	}, actions.Helo)

	postman.Tunnel.Register("authenticated", func() interface{} {
		return &actions.AuthenticatedMsg{}
	}, actions.Authenticated)

	postman.Tunnel.Register("sendMail", func() interface{} {
		return &actions.MailMsg{}
	}, actions.SendMail)

	// connect to middleware
	postman.Tunnel.Serve()
}
