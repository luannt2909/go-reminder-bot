package server

import (
	"encoding/json"
	"taskbot/pkg/task"
	"taskbot/pkg/util"
)

type GetListTaskRequest struct {
	Filter map[string]interface{} `form:"filter"`
	Range  string                 `form:"range"`
	Sort   string                 `form:"sort"`
}

func (r GetListTaskRequest) toGetListParams() task.GetListParams {
	p := task.GetListParams{
		Filter:   r.Filter,
		Limit:    10,
		Offset:   0,
		SortBy:   "id",
		SortType: "ASC",
	}
	if r.Range != "" {
		var queryRange []int
		_ = json.Unmarshal([]byte(r.Range), &queryRange)
		if len(queryRange) == 2 {
			p.Offset, p.Limit = queryRange[0], queryRange[1]
		}
	}
	if r.Sort != "" {
		var querySort []string
		_ = json.Unmarshal([]byte(r.Sort), &querySort)
		if len(querySort) == 2 {
			p.SortBy, p.SortType = querySort[0], querySort[1]
		}
	}
	return p
}

type Task struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	IsActive    bool   `json:"is_active"`
	Type        string `json:"type"`
	Schedule    string `json:"schedule"`
	NextTime    string `json:"next_time"`
	Message     string `json:"message"`
	Webhook     string `json:"webhook"`
	WebhookType string `json:"webhook_type"`
}

func transformTasksFromTasksDB(tasks []task.Task) []Task {
	result := make([]Task, 0, len(tasks))
	for _, task := range tasks {
		result = append(result, transformTaskFromTaskDB(task))
	}
	return result
}

func transformTaskFromTaskDB(t task.Task) Task {
	var nextTime string
	if nTime, err := util.Parse(t.Schedule); err == nil {
		nextTime = nTime.Format("15:04:05 02-01-2006")
	} else {
		nextTime = "invalid"
	}
	return Task{
		ID:          int64(t.Model.ID),
		Name:        t.Name,
		IsActive:    t.IsActive,
		Schedule:    t.Schedule,
		NextTime:    nextTime,
		Message:     t.Message,
		Webhook:     t.Webhook,
		WebhookType: t.WebhookType.String(),
	}
}
