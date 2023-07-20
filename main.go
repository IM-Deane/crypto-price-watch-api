package main

import "net/http"

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.ListenAndServe("localhost:8080", nil)
}
