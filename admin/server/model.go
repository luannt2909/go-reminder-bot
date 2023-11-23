package server

import (
	"encoding/json"
	task2 "taskbot/pkg/task"
)

type GetListTaskRequest struct {
	Filter map[string]interface{} `form:"filter"`
	Range  string                 `form:"range"`
	Sort   string                 `form:"sort"`
}

func (r GetListTaskRequest) toGetListParams() task2.GetListParams {
	p := task2.GetListParams{
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
	Type        string `json:"type"`
	Schedule    string `json:"schedule"`
	Message     string `json:"message"`
	Webhook     string `json:"webhook"`
	WebhookType string `json:"webhook_type"`
}

func transformTasksFromTasksDB(tasks []task2.Task) []Task {
	result := make([]Task, 0, len(tasks))
	for _, task := range tasks {
		result = append(result, transformTaskFromTaskDB(task))
	}
	return result
}

func transformTaskFromTaskDB(t task2.Task) Task {
	return Task{
		ID:          int64(t.Model.ID),
		Name:        t.Name,
		Schedule:    t.Schedule,
		Message:     t.Message,
		Webhook:     t.Webhook,
		WebhookType: t.WebhookType.String(),
	}
}
