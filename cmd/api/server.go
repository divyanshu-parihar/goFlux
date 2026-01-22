package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	config "github.com/divyanshu-parihar/goFlux/config"
	"github.com/redis/go-redis/v9"
)

type QueueHandler struct {
	rclient *redis.Client
}

type ListViewReqBody struct {
	Key string `json:"key"`
}

type ListViewRes struct {
	Data map[string]string `json:"data"`
}

func (handler *QueueHandler) ListQueueSize(w http.ResponseWriter, req *http.Request) {

	slog.Info("redis configuration", "host", os.Getenv("REDIS_HOST"))
	rclient, err := config.CreateRedisClient()

	if err != nil {
		slog.Error("Redis Server Failed")
		return
	}

	var body ListViewReqBody
	err = json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Invalid Body", http.StatusBadRequest)
		return
	}

	if body.Key == "" {
		http.Error(w, "Invalid Body Parameter", http.StatusBadRequest)
		return
	}
	value, err := config.GetRedisAllHValue(req.Context(), rclient, body.Key)
	if err != nil {
		http.Error(w, "Internal Error getting values", http.StatusInternalServerError)
		return
	}

	var response ListViewRes
	response.Data = value
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Internal Error getting values", http.StatusInternalServerError)
		return
	}
}

type QueueAddReqBody struct {
	Key   string `json:"data"`
	Value string `json:"value"`
}

func (handler *QueueHandler) QueueAdd(w http.ResponseWriter, req *http.Request) {

	var body ListViewReqBody
	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Invalid Body", http.StatusBadRequest)
		return
	}

	if body.Key == "" {
		http.Error(w, "Invalid Body Parameter", http.StatusBadRequest)
		return
	}
	var response interface{}

	json.NewDecoder(req.Body).Decode(&response)
	config.CreateRedisHValue(req.Context(), handler.rclient, body.Key, response)

}
func CreateServer(address, port string) {

	slog.Info("redis configuration", "host", os.Getenv("REDIS_HOST"))
	client, err := config.CreateRedisClient()

	if err != nil {
		slog.Error("Redis Server Failed")
		return
	}
	queueHandler := QueueHandler{
		rclient: client,
	}
	http.HandleFunc("/list", queueHandler.ListQueueSize)
	http.ListenAndServe(":"+port, nil)

}
