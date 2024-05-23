package model

import "time"

type Ticket struct {
	ID        uint64     `db:"id"`
	UserEmail string     `db:"user_email"`
	Price     uint64     `db:"price"`
	HomeTeam  string     `db:"home_team"`
	AwayTeam  string     `db:"away_team"`
	DateTime  *time.Time `db:"date_time"`
	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}
