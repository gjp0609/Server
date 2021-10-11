package commons

import (
	"encoding/json"
	"onysakura.fun/Server/commons/logrus"
	"os"
)

var log = logrus.GetLogger()

type Configuration struct {
	Port       int
	SQLitePath string
	PrivateKey string
	Notes      Notes
}

type Notes struct {
	Path string
}

var Configs Configuration

func init() {
	configFile := "./config.json"
	stat, err := os.Stat(configFile)
	if err != nil || stat.IsDir() {
		log.Warn("Can't find config file! Use default config.", err)
		Configs = Configuration{
			Port:       8080,
			SQLitePath: "./data.db",
			Notes:      Notes{Path: "./"},
		}
	} else {
		file, _ := os.Open(configFile)
		defer func(file *os.File) {
			_ = file.Close()
		}(file)
		decoder := json.NewDecoder(file)
		configuration := Configuration{}
		err = decoder.Decode(&configuration)
		if err != nil {
			log.Panic("error: ", err)
		}
		Configs = configuration
		log.Info("configs: ", Configs)
	}
}
