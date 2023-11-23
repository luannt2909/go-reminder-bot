package db

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"taskbot/pkg/enum"
	"taskbot/pkg/task"
	"time"
)

const dbFileName = "./tmp/taskbot.db"

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
func InitSQLiteDB() (*gorm.DB, error) {
	isExistedDB := false
	if fileExists(dbFileName) {
		fmt.Println("db is existed already")
		isExistedDB = true
	} else {
		fmt.Println("db is not exists, start create and migrate db")
		os.MkdirAll("./tmp", 0755)
		os.Create(dbFileName)
	}

	db, err := gorm.Open(sqlite.Open(dbFileName), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	//db = db.Debug()
	if !isExistedDB {
		_ = db.AutoMigrate(&task.Task{})
		db.Create(&task.Task{
			Name:        "Daily reminder task report",
			Schedule:    "* * * * * *",
			Message:     "report task please",
			Webhook:     "webhook",
			WebhookType: enum.WTGoogleChat,
			IsActive:    false,
			Model: gorm.Model{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				DeletedAt: gorm.DeletedAt{},
			},
		})
	}
	return db, nil
}
