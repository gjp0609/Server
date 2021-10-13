package notes

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"onysakura.fun/Server/commons"
	"onysakura.fun/Server/commons/data"
	"onysakura.fun/Server/commons/logrus"
	"onysakura.fun/Server/rest/user"
	"os/exec"
)

var log = logrus.GetLogger()

func Notes(writer http.ResponseWriter, request *http.Request) {
	var returnMsg = data.Msg{Code: -1}
	var err error
	defer func() {
		if err != nil {
			log.Warning(err)
			returnMsg.Msg = returnMsg.Msg + ", Reason: " + err.Error()
		}
		_, _ = writer.Write(returnMsg.ToString())
	}()
	authorization := request.Header.Get("Authorization")
	log.Info("Authorization: ", authorization)
	var username *string
	username, err = user.Auth(authorization)
	if err != nil {
		returnMsg.Msg = "auth fail"
		return
	}
	log.Info(fmt.Sprintf("username: %s", *username))
	returnMsg = data.Msg{Code: 1, Msg: "Hello " + *username}
}

func UpdateNotes(writer http.ResponseWriter, request *http.Request) {
	defer func() {
		writer.WriteHeader(200)
	}()
	command := exec.Command("git", "-C", commons.Configs.Notes.Path, "pull")
	output, err := command.Output()
	log.Info("update notes ", string(output), err)
}
