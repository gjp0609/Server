package rest

import "net/http"

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
