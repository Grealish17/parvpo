package server

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/Grealish17/parvpo/infrastructure/logger"
	"github.com/Grealish17/parvpo/internal/model"
)

type service interface {
	Buy(context.Context, *model.BuyTicketRequest) error
}

type Server struct {
	service service
}

func NewServer(service service) Server {
	return Server{
		service: service,
	}
}

var respChans = make(map[uint64]chan model.Message)
var mu sync.RWMutex

func (s *Server) Buy(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer req.Body.Close()

	if string(body) == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var ticket model.BuyTicketRequest
	if err = json.Unmarshal(body, &ticket); err != nil {
		//fmt.Println(err)
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	mu.Lock()
	respChans[ticket.ID] = make(chan model.Message)
	mu.Unlock()

	err = s.service.Buy(req.Context(), &ticket)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
		return
	}

	var rm model.Message
	mu.RLock()
	select {
	case rm = <-respChans[ticket.ID]:
		break
	case <-time.After(5 * time.Second):
		mu.RUnlock()
		if rm.UserEmail == "" {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("Oops, something went wrong"))
			logger.Warn("The timeout has triggered")
			return
		}
	}
	mu.RUnlock()

	if rm.UserEmail == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("This ticket has already been bought"))
		return
	}

	resp := &model.BuyTicketResponse{
		ID:        rm.ID,
		UserEmail: rm.UserEmail,
		Price:     rm.Price,
		HomeTeam:  rm.HomeTeam,
		AwayTeam:  rm.AwayTeam,
		DateTime:  rm.DateTime,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
