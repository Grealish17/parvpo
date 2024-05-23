package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Grealish17/parvpo/internal/model"
)

type service interface {
	Buy(context.Context, *model.BuyTicketRequest) (*model.BuyTicketResponse, error)
}

type Server struct {
	service service
}

func NewServer(service service) Server {
	return Server{
		service: service,
	}
}

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
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, err := s.service.Buy(req.Context(), &ticket)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
	}

	data, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
