package server

import (
	"errors"
	"fmt"
	"github.com/asaskevich/EventBus"
	"github.com/gin-gonic/gin"
	"go-reminder-bot/pkg/config"
	"go-reminder-bot/pkg/consts"
	"go-reminder-bot/pkg/enum"
	"go-reminder-bot/pkg/pusher"
	"go-reminder-bot/pkg/reminder"
	"go-reminder-bot/pkg/token"
	"go-reminder-bot/pkg/user"
	"net/http"
	"strconv"
)

type Handler struct {
	globalCfg       config.Config
	reminderStorage reminder.Storage
	userStorage     user.Storage
	pusher          pusher.Pusher
	EventBus.Bus
	tokenizer token.Tokenizer
}

func NewHandler(cfg config.Config, reminderStorage reminder.Storage, userStorage user.Storage, eventBus EventBus.Bus, pusher pusher.Pusher, tokenizer token.Tokenizer) *Handler {
	return &Handler{globalCfg: cfg, reminderStorage: reminderStorage, Bus: eventBus, userStorage: userStorage, pusher: pusher, tokenizer: tokenizer}
}

func (h Handler) findReminders(c *gin.Context) {
	var req GetListRequest
	err := c.ShouldBindQuery(&req)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	user := ExtractUserFromCtx(c)
	p := req.toGetListParams()
	reminderParams := reminder.GetListParams{GetListParams: p}
	if user.Role != enum.RoleAdmin {
		reminderParams.CreatedBy = &user.Email
	}
	list, count, err := h.reminderStorage.GetList(c, reminderParams)
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
	reminder, err := h.reminderStorage.GetOne(c, id)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, transformReminderFromReminderDB(reminder))
}

func (h Handler) deleteReminder(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	reminder, err := h.reminderStorage.GetOne(c, id)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	err = h.reminderStorage.Delete(c, id)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	h.triggerReloadJob(reminder.ID, enum.REDelete)
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
	reminder, err := h.reminderStorage.GetOne(c, id)
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
	err = h.reminderStorage.Update(c, reminder)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	reminder, err = h.reminderStorage.GetOne(c, id)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	h.triggerReloadJob(reminder.ID, enum.REUpdate)
	c.JSON(http.StatusOK, transformReminderFromReminderDB(reminder))
}

func (h Handler) createReminder(c *gin.Context) {
	var req Reminder
	err := c.ShouldBind(&req)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	user := ExtractUserFromCtx(c)
	reminderTotal, err := h.reminderStorage.CountRemindersByUser(c, user.Email)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if reminderTotal >= int64(h.globalCfg.MaximumReminder) {
		_ = c.AbortWithError(http.StatusInternalServerError, errors.New("you can only create a maximum of 5 reminders"))
		return
	}
	reminder := reminder.Reminder{
		Name:        req.Name,
		Schedule:    req.Schedule,
		IsActive:    req.IsActive,
		Message:     req.Message,
		Webhook:     req.Webhook,
		WebhookType: enum.ParseWebhookType(req.WebhookType),
		CreatedBy:   user.Email,
	}
	reminder, err = h.reminderStorage.Create(c, reminder)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	h.triggerReloadJob(reminder.ID, enum.RECreate)
	c.JSON(http.StatusOK, transformReminderFromReminderDB(reminder))
}

func (h Handler) triggerReloadJob(reminderID uint, event enum.ReminderEvent) {
	go h.Bus.Publish(consts.ReminderEventBusTopic, reminderID, event)
}
