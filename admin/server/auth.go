package server

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go-reminder-bot/pkg/reminder"
	"go-reminder-bot/pkg/user"
	"net/http"
	"os"
	"time"
)

func (h Handler) Login(c *gin.Context) {
	ctx := c.Request.Context()
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := c.ShouldBind(&req)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	u, err := h.userStorage.Authenticate(ctx, req.Username, req.Password)
	if errors.Is(err, user.ErrNotFound) {
		newUser := user.NewUser(req.Username, req.Password)
		user, err := h.userStorage.Create(ctx, newUser)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		reminderSample := reminder.DefaultReminder
		reminderSample.CreatedBy = user.Email
		_, _ = h.reminderStorage.Create(ctx, reminderSample)
		u = user
		goto FINISH
	}
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
FINISH:
	uRsp := transformUserFromUserDB(u)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user":      uRsp,
		"issued_at": time.Now().Unix(),
	})
	key := os.Getenv("JWT_SIGNING_KEY")
	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	uRsp.Token = tokenString

	c.JSON(http.StatusOK, uRsp)
}
