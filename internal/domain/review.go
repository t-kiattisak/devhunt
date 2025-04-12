package domain

import "time"

type Review struct {
	ID        int       `json:"id"`
	ToolID    int       `json:"tool_id"`
	UserID    string    `json:"user_id"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}
