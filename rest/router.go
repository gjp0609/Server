package rest

import (
	"fmt"
	"net/http"
	"onysakura.fun/Server/commons/logrus"
	"onysakura.fun/Server/rest/notes"
	"onysakura.fun/Server/rest/test"
	"onysakura.fun/Server/rest/user"
	"strconv"
)

var log = logrus.GetLogger()

type Router struct {
	path    string
	handler func(writer http.ResponseWriter, request *http.Request)
}

var Routers = []Router{
	{"/", Index},
	{"/test/get/", test.Get},
	{"/test/post/", test.Post},
	{"/user/login/", user.Login},
	{"/notes/list/", notes.List},
	{"/notes/get/", notes.GetNote},
	{"/notes/add/", notes.AddNote},
	{"/notes/update/", notes.UpdateNotes},
	{"/notes/", notes.ServeNotes},
}

func Run(port int) {
	mux := http.NewServeMux()
	for i := range Routers {
		mux.HandleFunc(Routers[i].path, handleInterceptor(Routers[i].handler))
	}
	_ = http.ListenAndServe(":"+strconv.Itoa(port), mux)
}

func Index(writer http.ResponseWriter, request *http.Request) {
	var name string
	if request.URL.Query()["name"] != nil {
		name = request.URL.Query()["name"][0]
	} else {
		name = "anonymous"
	}
	_, err := writer.Write([]byte("Hello " + name))
	if err != nil {
		return
	}
}

func handleInterceptor(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug(fmt.Sprintf("%7s:path: %s", r.Method, r.RequestURI))
		h(w, r)
	}
}
