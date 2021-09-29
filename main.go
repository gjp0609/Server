package main

import (
	"log"
	"net/http"
	"strconv"

	"onysakura.fun/Server/commons"
	"onysakura.fun/Server/rest"
)

func main() {
	http.HandleFunc("/", rest.Index)
	port := commons.Configs.Port
	log.Println("Server will start at http://127.0.0.1:" + strconv.Itoa(port))
	_ = http.ListenAndServe(":"+strconv.Itoa(port), nil)
}
