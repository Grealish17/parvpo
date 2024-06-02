package service

import (
	"context"
	"time"

	"github.com/Grealish17/parvpo/infrastructure/logger"
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
		//fmt.Println("Api send sync message error: ", err)
		logger.Errorf("Api send sync message error: ", err)
		return model.ErrSendReq
	}

	return nil
}

func (s *Service) Get(_ context.Context, id uint64) (*model.ReadTicketResponse, error) {
	time := time.Date(2024, time.May, 30, 0, 0, 0, 0, time.UTC)
	ticket := model.Ticket{
		ID:        id,
		UserEmail: "mock@yandex.ru",
		Price:     200,
		HomeTeam:  "Barabella",
		AwayTeam:  "Mifcar",
		DateTime:  &time,
	}

	resp := &model.ReadTicketResponse{
		ID:        ticket.ID,
		UserEmail: ticket.UserEmail,
		Price:     ticket.Price,
		HomeTeam:  ticket.HomeTeam,
		AwayTeam:  ticket.AwayTeam,
		DateTime:  ticket.DateTime,
	}

	return resp, nil
}
