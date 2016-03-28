package controllers

import "net/http"

func HelloController(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Hello, world!</h1>"))
}
