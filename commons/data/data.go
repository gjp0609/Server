package data

import (
	"encoding/json"
	"onysakura.fun/Server/commons/logrus"
)

var log = logrus.GetLogger()

type Msg struct {
	Code MsgCode
	Msg  string
	Data interface {
	}
}

type MsgCode int

const (
	MsgFail MsgCode = -1
	MsgOk   MsgCode = 1
)

func NewErrorMsg() Msg {
	return Msg{Code: MsgFail}
}

func (msg Msg) ToString() []byte {
	marshal, err := json.Marshal(msg)
	if err != nil {
		log.Warning("json marshal error: ", err)
	}
	return marshal
}
