package rest

import (
	"fmt"
	"net/http"
	"onysakura.fun/Server/commons"
	"onysakura.fun/Server/commons/logrus"
	"onysakura.fun/Server/rest/notes"
	"onysakura.fun/Server/rest/user"
	"strconv"
	"strings"
)

var log = logrus.GetLogger()

type Router struct {
	path    string
	handler func(writer http.ResponseWriter, request *http.Request)
}

var Routers = []Router{
	{"/", Index},
	{"/user/login/", user.Login},
	{"/notes/add/", notes.Notes},
	{"/notes/", ServeFile},
}

func Run(port int) {
	mux := http.NewServeMux()
	for i := range Routers {
		mux.HandleFunc(Routers[i].path, handleInterceptor(Routers[i].handler))
	}
	_ = http.ListenAndServe(":"+strconv.Itoa(port), mux)
}

func ServeFile(writer http.ResponseWriter, request *http.Request) {
	authorization := request.Header.Get("Authorization")
	_, err := user.Auth(authorization)
	if err != nil {
		writer.WriteHeader(401)
	} else {
		path := commons.Configs.Notes.Path + request.URL.Path
		path = strings.Replace(path, "/notes/", "", 1)
		log.Debug("serve file path: ", path)
		http.ServeFile(writer, request, path)
	}
}

func handleInterceptor(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug(fmt.Sprintf("%7s:path: %s", r.Method, r.URL.Path))
		h(w, r)
	}
}
