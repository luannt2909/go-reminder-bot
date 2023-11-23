package server

import (
	"context"
	"fmt"
	"github.com/asaskevich/EventBus"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"taskbot/pkg/consts"
	"taskbot/pkg/enum"
	"taskbot/pkg/task"
)

type Handler struct {
	task.Storage
	EventBus.Bus
}

func NewHandler(storage task.Storage, eventBus EventBus.Bus) *Handler {
	return &Handler{storage, eventBus}
}

func (h Handler) findTask(c *gin.Context) {
	var req GetListTaskRequest
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
	c.JSON(http.StatusOK, transformTasksFromTasksDB(list))
}

func (h Handler) getOneTask(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	task, err := h.GetOne(c, id)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, transformTaskFromTaskDB(task))
}

func (h Handler) deleteTask(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	task, err := h.GetOne(c, id)
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
	c.JSON(http.StatusOK, transformTaskFromTaskDB(task))
}

func (h Handler) updateTask(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	var req Task
	err := c.ShouldBind(&req)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	task, err := h.GetOne(c, id)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	task.Name = req.Name
	task.Schedule = req.Schedule
	task.Message = req.Message
	task.Webhook = req.Webhook
	task.IsActive = req.IsActive
	task.WebhookType = enum.ParseWebhookType(req.WebhookType)
	err = h.Update(c, task)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	task, err = h.GetOne(c, id)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	h.triggerReloadJob(c.Request.Context())
	c.JSON(http.StatusOK, transformTaskFromTaskDB(task))
}

func (h Handler) createTask(c *gin.Context) {
	var req Task
	err := c.ShouldBind(&req)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	task := task.Task{
		Name:        req.Name,
		Schedule:    req.Schedule,
		IsActive:    req.IsActive,
		Message:     req.Message,
		Webhook:     req.Webhook,
		WebhookType: enum.ParseWebhookType(req.WebhookType),
	}
	task, err = h.Storage.Create(c, task)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	h.triggerReloadJob(c.Request.Context())
	c.JSON(http.StatusOK, transformTaskFromTaskDB(task))
}

func (h Handler) triggerReloadJob(ctx context.Context) {
	go h.Bus.Publish(consts.TaskEventBusTopic, context.Background())
}
