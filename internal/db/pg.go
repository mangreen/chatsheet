package db

import (
	"fmt"
	"log/slog"

	"chatsheet/config"
	"chatsheet/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB 初始化資料庫連線並執行遷移
func InitDB(cfg *config.AppConfig) (*gorm.DB, error) {
	dbCfg := cfg.Database

	// 使用配置中的值來構建 DSN
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=Asia/Taipei",
		dbCfg.Host, dbCfg.User, dbCfg.Password, dbCfg.Name, dbCfg.Port, dbCfg.SSLMode)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		slog.Error("Failed to connect to database", "err", err)
		return nil, err
	}

	// 自動遷移模型
	err = DB.AutoMigrate(&model.User{}, &model.UnipileAccount{})
	if err != nil {
		slog.Error("Failed to database auto migrate", "err", err)
		return nil, err
	}
	slog.Info("Connecting to database", "on", dsn)

	return DB, nil
}
