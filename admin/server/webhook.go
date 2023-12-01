package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-reminder-bot/pkg/enum"
	"go-reminder-bot/pkg/pusher"
	"net/http"
)

func (h Handler) SendMessage(c *gin.Context) {
	ctx := c.Request.Context()
	var req struct {
		WebhookType string `json:"webhook_type"`
		Webhook     string `json:"webhook"`
		Message     string `json:"message"`
	}
	err := c.ShouldBind(&req)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	user := ExtractUserFromCtx(c)
	footer := fmt.Sprintf(pusher.GGChatMsgFooter, user.Email)
	message := fmt.Sprint(req.Message, footer)
	err = h.pusher.PushMessage(ctx, enum.ParseWebhookType(req.WebhookType), req.Webhook, message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}
