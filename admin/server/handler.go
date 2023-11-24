package server

import (
	"context"
	"fmt"
	"github.com/asaskevich/EventBus"
	"github.com/gin-gonic/gin"
	"go-reminder-bot/pkg/consts"
	"go-reminder-bot/pkg/enum"
	"go-reminder-bot/pkg/reminder"
	"net/http"
	"strconv"
)

type Handler struct {
	reminder.Storage
	EventBus.Bus
}

func NewHandler(storage reminder.Storage, eventBus EventBus.Bus) *Handler {
	return &Handler{storage, eventBus}
}

func (h Handler) findReminders(c *gin.Context) {
	var req GetListReminderRequest
	err := c.ShouldBindQuery(&req)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	p := req.toGetListParams()
	list, count, err := h.GetList(c, p)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Header("Content-Range", fmt.Sprintf("%d-%d/%d", p.Offset, p.Limit, count))
	c.JSON(http.StatusOK, transformRemindersFromRemindersDB(list))
}

func (h Handler) getOneReminder(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	reminder, err := h.GetOne(c, id)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, transformReminderFromReminderDB(reminder))
}

func (h Handler) deleteReminder(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	reminder, err := h.GetOne(c, id)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	err = h.Delete(c, id)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	h.triggerReloadJob(c.Request.Context())
	c.JSON(http.StatusOK, transformReminderFromReminderDB(reminder))
}

func (h Handler) updateReminder(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	var req Reminder
	err := c.ShouldBind(&req)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	reminder, err := h.GetOne(c, id)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	reminder.Name = req.Name
	reminder.Schedule = req.Schedule
	reminder.Message = req.Message
	reminder.Webhook = req.Webhook
	reminder.IsActive = req.IsActive
	reminder.WebhookType = enum.ParseWebhookType(req.WebhookType)
	err = h.Update(c, reminder)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	reminder, err = h.GetOne(c, id)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	h.triggerReloadJob(c.Request.Context())
	c.JSON(http.StatusOK, transformReminderFromReminderDB(reminder))
}

func (h Handler) createReminder(c *gin.Context) {
	var req Reminder
	err := c.ShouldBind(&req)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	reminder := reminder.Reminder{
		Name:        req.Name,
		Schedule:    req.Schedule,
		IsActive:    req.IsActive,
		Message:     req.Message,
		Webhook:     req.Webhook,
		WebhookType: enum.ParseWebhookType(req.WebhookType),
	}
	reminder, err = h.Storage.Create(c, reminder)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	h.triggerReloadJob(c.Request.Context())
	c.JSON(http.StatusOK, transformReminderFromReminderDB(reminder))
}

func (h Handler) triggerReloadJob(ctx context.Context) {
	go h.Bus.Publish(consts.ReminderEventBusTopic, context.Background())
}
