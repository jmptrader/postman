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
	var postman = Postman{
		Store:  store.New(dbDir, config.StoreSecret),
		Tunnel: tunnel.New(tunnelConfig()),
	}
	// TODO register all actions
	postman.Tunnel.Register("exit", func() interface{} {
		return &actions.ExitMsg{}
	}, actions.Exit)
	postman.Tunnel.Register("helo", func() interface{} {
		return &actions.HeloMsg{}
	}, actions.Helo)
	// postman.Tunnel.Register("action", func(){return &X{}, func(c Client, args interface{}){}})
	postman.Tunnel.Serve()
}
