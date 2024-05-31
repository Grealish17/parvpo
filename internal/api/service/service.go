package service

import (
	"context"
	"fmt"

	"github.com/Grealish17/parvpo/internal/model"
)

type sender interface {
	SendMessage(message model.Message) error
}

type Service struct {
	sender sender
}

func NewService(sender sender) *Service {
	return &Service{sender: sender}
}

func (s *Service) Buy(_ context.Context, ticket *model.BuyTicketRequest) error {
	err := s.sender.SendMessage(
		model.Message{
			ID:        ticket.ID,
			UserEmail: ticket.UserEmail,
			Price:     ticket.Price,
			HomeTeam:  ticket.HomeTeam,
			AwayTeam:  ticket.AwayTeam,
			DateTime:  ticket.DateTime,
		},
	)

	if err != nil {
		fmt.Println("Api send sync message error: ", err)
		return model.ErrSendReq
	}

	return nil
}
