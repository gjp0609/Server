package main

import (
	"onysakura.fun/Server/commons"
	"onysakura.fun/Server/commons/logrus"
	"onysakura.fun/Server/rest"
	"strconv"
)

var log = logrus.GetLogger()

type Han struct {
}

func main() {
	port := commons.Configs.Port
	log.Info("server will start at http://127.0.0.1:" + strconv.Itoa(port))
	rest.Run(port)
}
