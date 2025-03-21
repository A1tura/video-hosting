package main

import (
	"fmt"
	"net/http"
)

func handle(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./frontend/video.html")
}

func main() {

	http.Handle("/statics/", http.StripPrefix("/statics", http.FileServer(http.Dir("./statics/"))))
	http.HandleFunc("/videos/", handle)
	http.Handle("/", http.FileServer(http.Dir("./frontend/")))

	fmt.Println("xdd start")
	http.ListenAndServe("0.0.0.0:8080", nil)
}
