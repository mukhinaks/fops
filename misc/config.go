package misc

import (
	"encoding/json"
	"fmt"
	"os"
)

func ReadConfig(configPath string) map[string]interface{} {
	defaultConfigPath := "config.json"
	if configPath == "" {
		configPath = defaultConfigPath
		fmt.Println("Default path is used for configuration file.")
	}
	file, err := os.Open(configPath)
	if err != nil {
		fmt.Println("Unable to open configuration file.")
		return nil
	}
	var config map[string]interface{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return config
}
