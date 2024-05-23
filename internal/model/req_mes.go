package model

import "time"

type RequestMessage struct {
	ID        uint64
	UserEmail string
	Price     uint64
	HomeTeam  string
	AwayTeam  string
	DateTime  *time.Time
}
