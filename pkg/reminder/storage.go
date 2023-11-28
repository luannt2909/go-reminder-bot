package reminder

import (
	"context"
	"fmt"
	"gorm.io/gorm"
)

type Storage interface {
	Create(ctx context.Context, reminder Reminder) (Reminder, error)
	GetList(ctx context.Context, p GetListParams) ([]Reminder, int64, error)
	GetOne(ctx context.Context, id int64) (Reminder, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, reminder Reminder) error
	GetActiveReminder(ctx context.Context) ([]Reminder, error)
	GetActiveRemindersByUsers(ctx context.Context, emails []string) ([]Reminder, error)
}

type storage struct {
	db *gorm.DB
}

func (t *storage) GetActiveRemindersByUsers(ctx context.Context, emails []string) (reminders []Reminder, err error) {
	err = t.db.WithContext(ctx).
		Where("is_active = ? AND created_by IN ?", true, emails).
		Find(&reminders).Error
	if err != nil {
		return
	}
	return
}

func (t *storage) Update(ctx context.Context, reminder Reminder) (err error) {
	err = t.db.WithContext(ctx).Select("*").Updates(&reminder).Error
	return
}

func (t *storage) Delete(ctx context.Context, id int64) (err error) {
	err = t.db.WithContext(ctx).Delete(&Reminder{}, id).Error
	return
}

func (t *storage) Create(ctx context.Context, reminder Reminder) (result Reminder, err error) {
	db := t.db.WithContext(ctx).Create(&reminder)
	return reminder, db.Error
}

func (t *storage) GetList(ctx context.Context, param GetListParams) (reminders []Reminder, count int64, err error) {
	db := t.db.WithContext(ctx).Offset(param.Offset).
		Limit(param.Limit).
		Order(fmt.Sprintf("%s %s", param.SortBy, param.SortType))
	if param.CreatedBy != nil {
		db.Where("created_by", param.CreatedBy)
	}
	err = db.Find(&reminders).Count(&count).Error
	if err != nil {
		return
	}
	return
}

func (t *storage) GetActiveReminder(ctx context.Context) (reminders []Reminder, err error) {
	err = t.db.WithContext(ctx).
		Where("is_active", true).
		Find(&reminders).Error
	if err != nil {
		return
	}
	return
}

func (t *storage) GetOne(ctx context.Context, id int64) (reminder Reminder, err error) {
	err = t.db.WithContext(ctx).First(&reminder, id).Error
	return
}

func NewStorage(db *gorm.DB) Storage {
	return &storage{db: db}
}
