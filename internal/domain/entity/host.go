package entity

import "time"

type Host struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Vimid int    `json:"vimid"`
}

type OrderHost struct {
	Id       string    `json:"id"`
	RentDate time.Time `json:"rentDate"`
}
