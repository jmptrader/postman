package main

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
	DEFAULT_CONFIG_DIR = "./config"
)

type Config struct {
	AuthSecret  string    `json:"authSecret"`
	StoreSecret string    `json:"storeSecret"`
	RemoteAddr  string    `json:"remoteAddr"`
	CreatedAt   time.Time `json:"createdAt"`
}

var config Config
var dbDir, configDir string

// get postman dbDir and configDir from env path
// set as default if not found
func init() {
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
	configDir = os.Getenv("POSTMAN_CONFIG_DIR")
	if len(configDir) < 1 {
		dbDir = DEFAULT_CONFIG_DIR
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
	configFile, err := ioutil.ReadFile(configDir + "/config.json")
	if err != nil {
		return c, err
	}
	err = json.Unmarshal(configFile, &(c))
	return c, err
}

// get tunnel config
func tunnelConfig(st store.Store) tunnel.Config {
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
