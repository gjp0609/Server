package test

import (
	"encoding/json"
	"net/http"
	"onysakura.fun/Server/commons/logrus"
)

var log = logrus.GetLogger()

func Get(writer http.ResponseWriter, request *http.Request) {
	log.Info("params: ", request.URL.Query())
	writer.WriteHeader(200)
}

func Post(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var params map[string]interface{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Warning(err)
		writer.WriteHeader(500)
		return
	}
	log.Info("params: ", params)
	writer.WriteHeader(200)

}
