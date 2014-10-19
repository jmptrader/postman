package mail

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"postman/store"
	"postman/tunnel"
)

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

	store store.Store
}

const mailStoragePrefix = "mail:"

func (m *Mail) Create(c *tunnel.Client) (err error) {
	m.Sender = c.Config.Hostname
	m.From = fmt.Sprintf("bounce+%s-%s@%s", m.Id, strings.Replace(m.To, "@", "=", -1), m.Sender)
	mailStr, err := json.Marshal(m)
	if err != nil {
		return
	}
	err = c.Config.Store.Set(mailStoragePrefix+m.Id, string(mailStr))
	return
}

func GetMail(c *tunnel.Client, messageId string) (m Mail, err error) {
	store := c.Config.Store
	mailStr, ok := store.Get(mailStoragePrefix + messageId)
	if !ok {
		err = errors.New("no mail record found.")
		return
	}
	err = json.Unmarshal([]byte(mailStr), &m)
	m.store = store
	return
}

func (m *Mail) addr() string {
	addr := strings.SplitN(m.To, "@", 2)[1]
	return strings.ToLower(addr)
}

func (m *Mail) Update() error {
	mailStr, _ := json.Marshal(m)
	return m.store.Set(mailStoragePrefix+m.Id, string(mailStr))
}

func (m *Mail) CallWebHook(params map[string]string) (err error) {
	if len(m.WebHook) < 1 {
		return
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{
		Timeout:   time.Second * 20,
		Transport: tr,
	}
	v := url.Values{}
	v.Set("sender", m.Sender)
	v.Add("id", m.Id)
	v.Add("recipient", m.To)
	for key, value := range params {
		v.Add(key, value)
	}
	_, err = httpClient.PostForm(m.WebHook, v)
	return
}
