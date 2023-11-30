package server

import (
	"context"
	"github.com/gin-gonic/gin"
	"go-reminder-bot/pkg/token"
	"go-reminder-bot/pkg/user"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	UserKey             = "user"
	AuthorizationHeader = "Authorization"
	ExpireDuration      = 24 * time.Hour
)

func AuthenticateUserHandler(tokenizer token.Tokenizer, userStorage user.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenBearer := c.GetHeader(AuthorizationHeader)
		splitToken := strings.Split(tokenBearer, "Bearer ")
		var tokenStr string
		if len(splitToken) == 2 {
			tokenStr = splitToken[1]
		}
		claim, err := tokenizer.Parse(tokenStr)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		log.Println("Claim: ", claim)
		if time.Unix(claim.IssuedAt, 0).Add(ExpireDuration).Before(time.Now()) {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		u, err := userStorage.GetOne(c, int64(claim.UserID))
		if err != nil {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		if !u.IsActive {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		c.Set(UserKey, u)
		ctx := context.WithValue(c.Request.Context(), UserKey, u)
		c.Request.WithContext(ctx)
		c.Next()
	}
}

func ExtractUserFromCtx(c *gin.Context) user.User {
	if u, existed := c.Get(UserKey); existed {
		return u.(user.User)
	}
	return user.User{}
}
