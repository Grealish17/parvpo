package repository

import (
	"context"

	"github.com/Grealish17/parvpo/internal/app/db"
	"github.com/Grealish17/parvpo/internal/model"
)

type TicketsRepo struct {
	db db.DBops
}

func NewTickets(database db.DBops) *TicketsRepo {
	return &TicketsRepo{db: database}
}

func (r *TicketsRepo) Add(ctx context.Context, ticket *model.Ticket) error {
	var id uint64
	err := r.db.ExecQueryRow(ctx, "INSERT INTO tickets(id, user_email, price, home_team, away_team, date_time) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id;",
		ticket.ID, ticket.UserEmail, ticket.Price, ticket.HomeTeam, ticket.AwayTeam, ticket.DateTime).Scan(&id)
	if err != nil {
		return err
	}
	return nil
}
