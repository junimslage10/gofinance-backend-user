package dto

import (
	"time"
)

type UserRequest struct {
	ID          int32     `json:"user_id"`
	Username    string    `json:"username"`
	Password    string    `json:"password"`
	Email       string    `json:"email"`
	RequestedAt time.Time `json:"requested_at"`
}
