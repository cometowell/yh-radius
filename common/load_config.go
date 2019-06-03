package common

import (
	"encoding/json"
	"io/ioutil"
)

var Config map[string]interface{}

func loadConfig() map[string]interface{} {
	configBytes, err := ioutil.ReadFile("./config/radius.json")
	if err != nil {
		panic(err)
	}

	Config = make(map[string]interface{})
	json.Unmarshal(configBytes, &Config)
	return Config
}

func GetConfig() map[string]interface{} {
	return loadConfig()
}
