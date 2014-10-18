package actions

import (
	"log"

	"postman/tunnel"
)

type MailMsg struct {
	Content   string `json:"content"`
	Frome     string `json:"from"`
	To        string `json:"to"`
	WebHook   string `json:"web_hook"`
	Immediate string `json:"immediate"`
	Block     string `json:"block"`
}

func SendMail(c *tunnel.Client, args interface{}) {
	mail := args.(*MailMsg)
	log.Print(mail.Content)
}
