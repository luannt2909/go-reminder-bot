package db

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"go-reminder-bot/pkg/config"
	"go-reminder-bot/pkg/enum"
	"go-reminder-bot/pkg/reminder"
	"go-reminder-bot/pkg/user"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

const dbFileName = "./tmp/go-reminder-bot.db"

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
		_ = db.AutoMigrate(&reminder.Reminder{}, &user.User{})
		db.Create(&reminder.DefaultReminder)
		db.Create(&user.User{
			Model: gorm.Model{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Email:    "admin@reminderbot.com",
			Password: "reminderbot",
			Role:     enum.RoleAdmin,
			IsActive: true,
		})
		db.Create(&user.User{
			Model: gorm.Model{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Email:    "guest@reminderbot.com",
			Password: "reminderbot",
			Role:     enum.RoleGuest,
			IsActive: true,
		})
	}
	return db, nil
}

func InitDatabase(cfg config.DBConfig) (db *gorm.DB, err error) {
	log.Println("config", cfg)
	db, err = initDB(cfg)
	if err != nil {
		return
	}
	err = db.AutoMigrate(&user.User{}, &reminder.Reminder{})
	return
}

func initDB(cfg config.DBConfig) (db *gorm.DB, err error) {
	switch cfg.DBClient {
	case "mysql":
		db, err = gorm.Open(mysql.Open(cfg.DBConnectionURI), &gorm.Config{})
		if err != nil {
			return
		}
		return
	default:
		return InitSQLiteDB()
	}
}
