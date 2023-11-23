package task

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"taskbot/pkg/enum"
)

type Task struct {
	gorm.Model
	Name        string           `json:"name" gorm:"name"`
	Schedule    string           `json:"schedule" gorm:"schedule"`
	Message     string           `json:"message" gorm:"message"`
	Webhook     string           `json:"webhook" gorm:"webhook"`
	WebhookType enum.WebhookType `json:"webhook_type" gorm:"webhook_type"`
	Type        int32            `json:"type" gorm:"type"`
	IsActive    bool             `json:"is_active" gorm:"is_active"`
}
type Storage interface {
	Create(ctx context.Context, task Task) (Task, error)
	GetList(ctx context.Context, p GetListParams) ([]Task, int64, error)
	GetOne(ctx context.Context, id int64) (Task, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, task Task) error
	GetActiveTasks(ctx context.Context) ([]Task, error)
}

type storage struct {
	db *gorm.DB
}

func (t *storage) Update(ctx context.Context, task Task) (err error) {
	err = t.db.WithContext(ctx).Select("*").Updates(&task).Error
	return
}

func (t *storage) Delete(ctx context.Context, id int64) (err error) {
	err = t.db.WithContext(ctx).Delete(&Task{}, id).Error
	return
}

func (t *storage) Create(ctx context.Context, task Task) (result Task, err error) {
	db := t.db.WithContext(ctx).Create(&task)
	return task, db.Error
}

func (t *storage) GetList(ctx context.Context, param GetListParams) (tasks []Task, count int64, err error) {
	err = t.db.WithContext(ctx).Offset(param.Offset).
		Limit(param.Limit).
		Order(fmt.Sprintf("%s %s", param.SortBy, param.SortType)).
		Find(&tasks).Count(&count).Error
	if err != nil {
		return
	}
	return
}

func (t *storage) GetActiveTasks(ctx context.Context) (tasks []Task, err error) {
	err = t.db.WithContext(ctx).
		Where("is_active", true).
		Find(&tasks).Error
	if err != nil {
		return
	}
	return
}

func (t *storage) GetOne(ctx context.Context, id int64) (task Task, err error) {
	err = t.db.WithContext(ctx).First(&task, id).Error
	fmt.Println(err)
	return
}

func NewStorage(db *gorm.DB) Storage {
	return &storage{db: db}
}
