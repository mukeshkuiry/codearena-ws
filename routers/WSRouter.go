package routers

import (
	"github.com/gorilla/mux"

	"github.com/CodeArena-Org/codearena-ws/handlers"
)

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/ws", handlers.WSHandler)
	return router
}
