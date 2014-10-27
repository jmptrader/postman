package client

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"

	"postman/store"
	"postman/tunnel"
)

const (
	DEFAULT_DB_DIR     = "./data"
	DEFAULT_CONFIG_DIR = "."
	CLIENT_PRIVATE_KEY = "sys:private"
)

type Config struct {
	AuthSecret  string    `json:"authSecret"`
	StoreSecret string    `json:"storeSecret"`
	RemoteAddr  string    `json:"remoteAddr"`
	Hostname    string    `json:"hostname"`
	CreatedAt   time.Time `json:"createdAt"`
}

type Client struct {
	Tunnel     tunnel.Tunnel
	Store      store.Store
	Hostname   string
	PrivateKey string
}

var Postman *Client
var config Config
var dbDir, configDir string

func init() {
	configInit()
	st := store.New(dbDir, config.StoreSecret)
	Postman = &Client{
		Store:    st,
		Tunnel:   tunnel.New(getTunnelConfig(st)),
		Hostname: config.Hostname,
	}
	go Postman.setPrivateKey()
}

// todo: here be block for some reason
func (p *Client) setPrivateKey() {
	pk, ok := p.Store.Get(CLIENT_PRIVATE_KEY)
	if ok {
		p.PrivateKey = pk
		return
	}
	pk, err := p.Tunnel.RequestBlock("privateKey", map[string]string{})
	if err != nil {
		log.Fatalf("client: get private key %s", err.Error())
	}
	p.Store.Set(CLIENT_PRIVATE_KEY, pk)
	p.PrivateKey = pk
}

func configInit() {
	var err error
	dbDir = os.Getenv("POSTMAN_DB_DIR")
	if len(dbDir) < 1 {
		dbDir = DEFAULT_DB_DIR
	}
	// create dir if not exist
	if _, err = os.Stat(dbDir); err != nil {
		if os.IsNotExist(err) {
			// file does not exist
			if os.Mkdir(dbDir, 0755) != nil {
				log.Fatalf("can not create db dir %s", dbDir)
			}
		}
	}
	config, err = loadConfig()
	if err != nil {
		log.Fatalf("load config file %s", err.Error())
	}
}

// load config file and parse to struct
// exit with error if meet any error
func loadConfig() (Config, error) {
	c := Config{}
	configDir = os.Getenv("POSTMAN_CONFIG_DIR")
	if len(configDir) < 1 {
		configDir = DEFAULT_CONFIG_DIR
	}
	configFile, err := ioutil.ReadFile(configDir + "/config.json")
	if err != nil {
		return c, err
	}
	err = json.Unmarshal(configFile, &(c))
	return c, err
}

// get tunnel config
func getTunnelConfig(st store.Store) tunnel.Config {
	conf := &tls.Config{
		InsecureSkipVerify: true,
	}
	return tunnel.Config{
		Conf:   conf,
		Remote: config.RemoteAddr,
		Secret: config.AuthSecret,
		Store:  st,
	}
}
