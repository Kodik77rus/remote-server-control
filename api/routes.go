package api

import (
	"net/http"
	"remote-server-control/internal/handlers"
)

const _apiPrifix = "/api/v1"

func CreateRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle(_apiPrifix+"/remote-execution", http.HandlerFunc(handlers.ExecuteRemoteCommand))

	return mux
}
