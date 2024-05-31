package service

import (
	"context"

	"github.com/Grealish17/parvpo/internal/model"
)

type repository interface {
	Add(ctx context.Context, ticket *model.Ticket) error
}

type TicketsService struct {
	repo repository
}

func NewTicketsService(repository repository) *TicketsService {
	return &TicketsService{repo: repository}
}

func (s *TicketsService) Add(ctx context.Context, ticket *model.AddTicketRequest) (*model.AddTicketResponse, error) {
	ticketRepo := &model.Ticket{
		ID:        ticket.ID,
		UserEmail: ticket.UserEmail,
		Price:     ticket.Price,
		HomeTeam:  ticket.HomeTeam,
		AwayTeam:  ticket.AwayTeam,
		DateTime:  ticket.DateTime,
	}

	err := s.repo.Add(ctx, ticketRepo)
	if err != nil {
		return nil, err
	}

	resp := &model.AddTicketResponse{
		ID:        ticket.ID,
		UserEmail: ticket.UserEmail,
		Price:     ticket.Price,
		HomeTeam:  ticket.HomeTeam,
		AwayTeam:  ticket.AwayTeam,
		DateTime:  ticket.DateTime,
	}

	return resp, nil
}
