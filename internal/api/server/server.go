package server

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/Grealish17/parvpo/infrastructure/logger"
	"github.com/Grealish17/parvpo/infrastructure/redis"
	"github.com/Grealish17/parvpo/internal/model"
	"github.com/gorilla/mux"
)

const QueryParamId = "id"

type service interface {
	Buy(context.Context, *model.BuyTicketRequest) error
	Get(context.Context, uint64) (*model.ReadTicketResponse, error)
}

type Server struct {
	service service
	cache   redis.Cache
}

func NewServer(service service, cache redis.Cache) Server {
	return Server{
		service: service,
		cache:   cache,
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

func (s *Server) Get(w http.ResponseWriter, req *http.Request) {
	id, ok := mux.Vars(req)[QueryParamId]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	idInt, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := s.cache.Get(req.Context(), req.URL.Path)
	if err == nil {
		w.WriteHeader(resp.Status)
		w.Write(resp.Data)
		return
	}

	data, status := s.ticketGet(req.Context(), idInt)

	err = s.cache.Set(req.Context(), req.URL.Path,
		&model.Response{
			Status: status,
			Data:   data,
		},
	)
	if err != nil {
		//log.Println(err)
		logger.Error(err)
	}

	w.WriteHeader(status)
	w.Write(data)
}

func (s *Server) ticketGet(ctx context.Context, id uint64) ([]byte, int) {
	point, err := s.service.Get(ctx, id)
	if err != nil {
		if errors.Is(err, model.ErrObjectNotFound) {
			return []byte(err.Error()), http.StatusNotFound
		}
		return nil, http.StatusInternalServerError
	}

	pointJson, err := json.Marshal(point)
	if err != nil {
		return nil, http.StatusInternalServerError
	}

	return pointJson, http.StatusOK
}
