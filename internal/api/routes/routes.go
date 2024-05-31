package routes

import (
	"log"
	"net/http"

	"github.com/Grealish17/parvpo/internal/api/server"

	"github.com/gorilla/mux"
)

func CreateRouter(implemetation server.Server) *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/tickets", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			implemetation.Buy(w, req)
		default:
			log.Println("error")
		}
	})

	return router
}
