package mail

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/jpoz/dkim"

	"postman/client"
	"postman/util"
)

var postman = client.Postman

type Mail struct {
	Id        string `json:"id"`
	From      string `json:"from"`
	To        string `json:"to"`
	Content   string `json:"content"`
	Immediate bool   `json:"immediate"`
	Retries   int    `json:"retries"`
	WebHook   string `json:"web_hook,omitempty"`
	Log       string `json:"log,omitempty"`
	Sender    string `json:"sender"`
}

const (
	MAIL_STORAGE_PREFIX = "mail:"
	DKIM_SELECTOR       = "mx"
)

func (m *Mail) Create() (err error) {
	m.From = fmt.Sprintf("bounce+%s-%s@%s", m.Id, strings.Replace(m.To, "@", "=", -1), postman.Hostname)
	mailStr, err := json.Marshal(m)
	if err != nil {
		return
	}
	err = postman.Store.Set(MAIL_STORAGE_PREFIX+m.Id, string(mailStr))
	return
}

func Get(messageId string) (m *Mail, err error) {
	m = &Mail{}
	mailStr, ok := postman.Store.Get(MAIL_STORAGE_PREFIX + messageId)
	if !ok {
		err = errors.New("no mail record found.")
		return
	}
	err = json.Unmarshal([]byte(mailStr), m)
	return
}

func (m *Mail) Addr() string {
	addr := strings.SplitN(m.To, "@", 2)[1]
	return strings.ToLower(addr)
}

func (m *Mail) Update() error {
	mailStr, _ := json.Marshal(m)
	return postman.Store.Set(MAIL_STORAGE_PREFIX+m.Id, string(mailStr))
}

func (m *Mail) Destroy() {
	postman.Store.Destroy(MAIL_STORAGE_PREFIX + m.Id)
}

func (m *Mail) CallWebHook(params map[string]string) (err error) {
	if len(m.WebHook) < 1 {
		return
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{
		Transport: tr,
	}
	v := url.Values{}
	v.Set("sender", postman.Hostname)
	v.Add("id", m.Id)
	v.Add("recipient", m.To)
	for key, value := range params {
		v.Add(key, value)
	}
	_, err = httpClient.PostForm(m.WebHook, v)
	return
}

func (m *Mail) Deliver() error {
	conf, err := dkim.NewConf(postman.Hostname, DKIM_SELECTOR)
	if err != nil {
		log.Fatalf("dkim: config %s", err.Error())
	}
	d, _ := dkim.New(conf, []byte(postman.PrivateKey))
	msg, _ := d.Sign([]byte(m.Content))
	log.Printf("mail: %s start sending", m.Id)
	err = util.SendMail(m.From, m.To, msg, postman.Hostname)
	result := map[string]interface{}{
		"id":      m.Id,
		"success": err == nil,
	}
	if err != nil {
		result["log"] = err.Error()
		m.Log = err.Error()
		m.Update()
	}
	postman.Tunnel.Request("log", result)
	return err
}
