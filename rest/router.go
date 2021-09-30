package rest

import (
	"net/http"
	"onysakura.fun/Server/rest/notes"
	"strconv"
)

type Router struct {
	path    string
	handler func(writer http.ResponseWriter, request *http.Request)
}

var Routers = []Router{
	{"/", Index},
	{"/notes/", notes.Notes},
}

func Run(port int) {
	mux := http.NewServeMux()
	for i := range Routers {
		mux.HandleFunc(Routers[i].path, Routers[i].handler)
	}
	_ = http.ListenAndServe(":"+strconv.Itoa(port), mux)
}
