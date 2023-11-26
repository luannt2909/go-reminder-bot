package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	UserKey             = "user"
	AuthorizationHeader = "Authorization"
)

func AuthenticateUser(c *gin.Context) {
	tokenBearer := c.GetHeader(AuthorizationHeader)
	splitToken := strings.Split(tokenBearer, "Bearer ")
	var tokenStr string
	if len(splitToken) == 2 {
		tokenStr = splitToken[1]
	}
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("JWT_SIGNING_KEY")), nil
	})
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		userClaim := claims["user"]
		var user User
		_ = ParseJson(userClaim, &user)
		c.Set(UserKey, user)
		ctx := context.WithValue(c.Request.Context(), UserKey, user)
		c.Request.WithContext(ctx)
	}
	c.Next()
}

func ParseJson(data, target interface{}) error {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = json.Unmarshal(dataBytes, &target)
	return err
}

func ExtractUserFromCtx(c *gin.Context) User {
	if u, existed := c.Get(UserKey); existed {
		return u.(User)
	}
	return User{}
}
