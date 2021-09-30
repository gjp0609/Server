package notes

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"onysakura.fun/Server/commons/logrus"
	"onysakura.fun/Server/rest/user"
	"strconv"
)

var log = logrus.GetLogger()

func Notes(writer http.ResponseWriter, request *http.Request) {
	var err error
	defer func() {
		if err != nil {
			log.Warning(err)
		}
	}()
	var username string
	if request.URL.Query()["username"] != nil {
		username = request.URL.Query()["username"][0]
	} else {
		username = "anonymous"
	}
	userEntry := user.GetUser(username)
	if userEntry != nil {
		log.Info("user find")
	} else {
		log.Warning("user not found")
	}
	log.Info(fmt.Sprintf("username: %s", username))
	_, err = writer.Write([]byte("Hello " + username + ", id: " + strconv.Itoa(userEntry.Id)))
}
