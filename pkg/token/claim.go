package token

import (
	"go-reminder-bot/pkg/enum"
)

type Claim struct {
	UserID    uint
	UserEmail string
	Role      enum.UserRole
	IssuedAt  int64
}
