package service

import (
	"context"
	"fmt"

	"github.com/Grealish17/parvpo/internal/model"
)

type sender interface {
	SendMessage(message model.RequestMessage) error
}

type Service struct {
	sender sender
}

func NewService(sender sender) *Service {
	return &Service{sender: sender}
}

func (s *Service) Buy(_ context.Context, ticket *model.BuyTicketRequest) (*model.BuyTicketResponse, error) {
	err := s.sender.SendMessage(
		model.RequestMessage{
			ID:        ticket.ID,
			UserEmail: ticket.UserEmail,
			Price:     ticket.Price,
			HomeTeam:  ticket.HomeTeam,
			AwayTeam:  ticket.AwayTeam,
			DateTime:  ticket.DateTime,
		},
	)

	if err != nil {
		fmt.Println("Send sync message error: ", err)
		return nil, model.ErrSendReq
	}

	resp := &model.BuyTicketResponse{
		ID:        ticket.ID,
		UserEmail: ticket.UserEmail,
		Price:     ticket.Price,
		HomeTeam:  ticket.HomeTeam,
		AwayTeam:  ticket.AwayTeam,
		DateTime:  ticket.DateTime,
	}

	return resp, nil
}
