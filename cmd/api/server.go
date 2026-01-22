package server

import (
	"fmt"
	config "github.com/divyanshu-parihar/goFlux/config"
	"log/slog"
	"net/http"
	"os"
)

func ListQueueSize(w http.ResponseWriter, res *http.Request) {

	slog.Info("redis configuration", "host", os.Getenv("REDIS_HOST"))
	err := config.CreateRedisClient()
	if err != nil {
		slog.Error("Redis Server Failed")
		return
	}

	fmt.Fprintf(w, "hello\n")
}

func CreateServer(address, port string) {
	http.HandleFunc("/list", ListQueueSize)
	http.ListenAndServe(":"+port, nil)
}
