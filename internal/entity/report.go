package entity

import (
	"github.com/jackc/pgtype"
	"time"
)

// Report -.
type Report struct {
	Id        int          `json:"id" db:"id" example:"2"`
	PingerId  int          `json:"pinger_id" db:"pinger_id" example:"10"`
	Content   pgtype.JSONB `json:"content" db:"content"`
	CreatedAt time.Time    `json:"created_at" db:"created_at" example:"2020-09-09T14:51:00Z"`
}
