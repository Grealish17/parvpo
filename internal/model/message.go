package model

import "time"

type Message struct {
	ID        uint64
	UserEmail string
	Price     uint64
	HomeTeam  string
	AwayTeam  string
	DateTime  *time.Time
}
