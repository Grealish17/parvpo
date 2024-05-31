package server

import (
	"context"
	"fmt"
	"log"

	"github.com/Grealish17/parvpo/internal/model"
)

type sender interface {
	SendMessage(message model.Message) error
}

type service interface {
	Add(ctx context.Context, ticket *model.AddTicketRequest) (*model.AddTicketResponse, error)
}

type Server struct {
	service service
	sender  sender
	msgChan <-chan model.Message
}

func NewServer(service service, sender sender, msgChan chan model.Message) Server {
	return Server{
		service: service,
		msgChan: msgChan,
		sender:  sender,
	}
}

func (s *Server) Listen(ctx context.Context) {
	for msg := range s.msgChan {
		ticket := &model.AddTicketRequest{
			ID:        msg.ID,
			UserEmail: msg.UserEmail,
			Price:     msg.Price,
			HomeTeam:  msg.HomeTeam,
			AwayTeam:  msg.AwayTeam,
			DateTime:  msg.DateTime,
		}

		var rm model.Message
		resp, err := s.service.Add(ctx, ticket)
		if err != nil {
			rm = model.Message{
				ID: ticket.ID,
			}
		} else {
			rm = model.Message{
				ID:        resp.ID,
				UserEmail: resp.UserEmail,
				Price:     resp.Price,
				HomeTeam:  resp.HomeTeam,
				AwayTeam:  resp.AwayTeam,
				DateTime:  resp.DateTime,
			}
		}

		err = s.sender.SendMessage(rm)
		if err != nil {
			fmt.Println("App send sync message error: ", err)
		}
	}
	log.Println("App channel closed, exiting listen goroutine")
}
