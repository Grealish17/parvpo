package routes

import (
	"fmt"
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

	router.HandleFunc(fmt.Sprintf("/tickets/{%s:[0-9]+}", server.QueryParamId), func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			implemetation.Get(w, req)
		default:
			//log.Println("error")
			logger.Warn("Default case in HandleFunc /tickets/[0-9]+")
		}
	})

	return router
}
