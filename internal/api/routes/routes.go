package routes

import (
	"net/http"
	"time"

	"github.com/Grealish17/parvpo/infrastructure/logger"
	"github.com/Grealish17/parvpo/internal/api/server"

	"github.com/gorilla/mux"
)

func CreateRouter(implemetation server.Server) *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/tickets", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			start := time.Now()
			implemetation.Buy(w, req)
			logger.Infof("Time to response: %v", time.Since(start))
		default:
			//log.Println("error")
			logger.Warn("Default case in HandleFunc /tickets")
		}
	})

	return router
}
