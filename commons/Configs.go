package commons

import (
	"encoding/json"
	"log"
	"os"
)

type Configuration struct {
	Port  int
	Notes Notes
}

type Notes struct {
	Path string
}

var Configs Configuration

func init() {
	configFile := "config.json"
	stat, err := os.Stat(configFile)
	if err != nil || stat.IsDir() {
		log.Println("Can't find config file!", err)
		Configs = Configuration{
			Port:  8080,
			Notes: Notes{Path: "./"},
		}
	} else {
		log.Println(os.Stat(configFile))
		file, _ := os.Open(configFile)
		defer func(file *os.File) {
			_ = file.Close()
		}(file)
		decoder := json.NewDecoder(file)
		configuration := Configuration{}
		err = decoder.Decode(&configuration)
		if err != nil {
			log.Println("error:", err)
		}
		Configs = configuration
		log.Println("configs: ", Configs)
	}
}

func GetConfig() Configuration {
	return Configs
}
