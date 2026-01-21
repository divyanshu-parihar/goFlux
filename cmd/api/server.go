package server

import (
	"fmt"
	"net/http"
)

func ListQueueSize(w http.ResponseWriter, res *http.Request) {
	fmt.Fprintf(w, "hello\n")
}

func CreateServer(address, port string) {
	http.HandleFunc("/list", ListQueueSize)
	http.ListenAndServe(":"+port, nil)
}
