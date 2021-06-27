package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Configuration struct {
	ControllerURL string `json:"controller_url"`
	Site          string `json:"site"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	SkipSSLVerify bool   `json:"skip_ssl_verify"`
	Timeout       int    `json:"timeout"`
}

func (c *Configuration) Load(confFile string) error {
	file, err := os.Open(confFile)

	if err != nil {
		return fmt.Errorf("can't open config file: %s", confFile)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Print("error occurred when close configuration file: ", err.Error())
		}
	}(file)

	decoder := json.NewDecoder(file)

	err = decoder.Decode(c)
	if err != nil {
		return fmt.Errorf("fatal error occurred when decoding json config file: %s", err.Error())
	}

	return nil
}
