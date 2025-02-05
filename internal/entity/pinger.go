package entity

import "time"

// Pinger -.
type Pinger struct {
	Id        int       `json:"id" db:"id" example:"1"`
	Name      string    `json:"name" db:"name" example:"linux-vm-78"`
	Password  string    `json:"password" db:"password" example:"various-pass"`
	CreatedAt time.Time `json:"created_at" db:"created_at" example:"2020-09-20T14:00:00Z"`
}
