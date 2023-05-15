package models

import (
    "time"
)

type Admin struct {
	Id       uint   `json:"id"`
	Address  string `json:"address"`
}

type Account struct {
	Id        uint   	`json:"id"`
	Address   string 	`json:"address"`
	Status    int 		`json:"status"`
	Fee 	  string    `json:"fee"`  // -: received, <address>: awaiting receipt
	Type 	  string    `json:"type"` // 0: Account, 1: Contract
	CreatedAt time.Time	`json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
