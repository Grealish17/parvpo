package model

import "time"

type BuyTicketRequest struct {
	ID        uint64     `json:"ID"`
	UserEmail string     `json:"UserEmail"`
	Price     uint64     `json:"Price"`
	HomeTeam  string     `json:"HomeTeam"`
	AwayTeam  string     `json:"AwayTeam"`
	DateTime  *time.Time `json:"DateTime"`
}

type BuyTicketResponse struct {
	ID        uint64     `json:"ID"`
	UserEmail string     `json:"UserEmail"`
	Price     uint64     `json:"Price"`
	HomeTeam  string     `json:"HomeTeam"`
	AwayTeam  string     `json:"AwayTeam"`
	DateTime  *time.Time `json:"DateTime"`
}

type AddTicketRequest struct {
	ID        uint64     `json:"ID"`
	UserEmail string     `json:"UserEmail"`
	Price     uint64     `json:"Price"`
	HomeTeam  string     `json:"HomeTeam"`
	AwayTeam  string     `json:"AwayTeam"`
	DateTime  *time.Time `json:"DateTime"`
}

type AddTicketResponse struct {
	ID        uint64     `json:"ID"`
	UserEmail string     `json:"UserEmail"`
	Price     uint64     `json:"Price"`
	HomeTeam  string     `json:"HomeTeam"`
	AwayTeam  string     `json:"AwayTeam"`
	DateTime  *time.Time `json:"DateTime"`
}

type ReadTicketResponse struct {
	ID        uint64     `json:"ID"`
	UserEmail string     `json:"UserEmail"`
	Price     uint64     `json:"Price"`
	HomeTeam  string     `json:"HomeTeam"`
	AwayTeam  string     `json:"AwayTeam"`
	DateTime  *time.Time `json:"DateTime"`
}

type Response struct {
	Status int
	Data   []byte
}
