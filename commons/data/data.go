package data

import (
	"encoding/json"
	"onysakura.fun/Server/commons/logrus"
)

var log = logrus.GetLogger()

type Msg struct {
	Code int
	Msg  string
	Data interface {
	}
}

func NewErrorMsg(msg string) Msg {
	return Msg{Code: -1, Msg: msg}
}

func (msg Msg) ToString() []byte {
	marshal, err := json.Marshal(msg)
	if err != nil {
		log.Warning("json marshal error: ", err)
	}
	return marshal
}
