package server

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"go-reminder-bot/pkg/reminder"
	"go-reminder-bot/pkg/token"
	"go-reminder-bot/pkg/user"
	"go-reminder-bot/pkg/util"
	"net/http"
	"time"
)

func (h Handler) Login(c *gin.Context) {
	ctx := c.Request.Context()
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	err := c.ShouldBind(&req)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	passwordHash := util.Hash256(req.Password)
	u, err := h.userStorage.GetByEmail(ctx, req.Username)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			u, err = h.registerUser(c, req.Username, passwordHash)
			if err != nil {
				_ = c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
		} else {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}
	if u.Password != passwordHash {
		// TODO: migrate new hash password
		if req.Password != u.Password {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "password is incorrect",
			})
			return
		}
		go h.userStorage.UpdateNewPassword(context.Background(), u.ID, passwordHash)
	}
	if !u.IsActive {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "account is inactive",
		})
		return
	}
	claim := token.Claim{
		UserID:    u.ID,
		UserEmail: u.Email,
		Role:      u.Role,
		IssuedAt:  time.Now().Unix(),
	}
	tokenStr, err := h.tokenizer.Generate(claim)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	uRsp := transformUserFromUserDB(u)
	uRsp.Token = tokenStr
	c.JSON(http.StatusOK, uRsp)
}

func (h Handler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

func (h Handler) registerUser(ctx context.Context, username, password string) (u user.User, err error) {
	newUser := user.NewUser(username, password)
	u, err = h.userStorage.Create(ctx, newUser)
	if err != nil {
		return
	}
	reminderSample := reminder.DefaultReminder
	reminderSample.CreatedBy = u.Email
	_, _ = h.reminderStorage.Create(ctx, reminderSample)
	return
}
