package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"taskbot/pkg/enum"
	"taskbot/pkg/task"
)

type Handler struct {
	task.Storage
}

func NewHandler(storage task.Storage) *Handler {
	return &Handler{storage}
}

func (h Handler) findTask(c *gin.Context) {
	var req GetListTaskRequest
	err := c.ShouldBindQuery(&req)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	fmt.Println("Req:", req)
	p := req.toGetListParams()
	fmt.Println("params", p)
	list, count, err := h.GetList(c, p)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	fmt.Println(list, count)
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
		Message:     req.Message,
		Webhook:     req.Webhook,
		WebhookType: enum.ParseWebhookType(req.WebhookType),
	}
	task, err = h.Storage.Create(c, task)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, transformTaskFromTaskDB(task))
}
