package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"chatsheet/config"
	"chatsheet/internal/db"
	"chatsheet/internal/handler"
	"chatsheet/internal/repository/gormimpl"
	"chatsheet/internal/service"

	"github.com/MatusOllah/slogcolor"
	"github.com/gin-gonic/gin"
	//_ "chatsheet/docs" // Swagger 產生的文件
)

func main() {
	gin.ForceConsoleColor()

	// 初始化結構化日誌
	slog.SetDefault(slog.New(slogcolor.NewHandler(os.Stderr, slogcolor.DefaultOptions)))

	// 載入配置
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("Failed to load configuration", "err", err)
		os.Exit(1)
	}
	slog.Info("Configuration loaded successfully")

	// 等待 DB 啟動
	db, err := db.InitDB(cfg)
	if err != nil {
		slog.Error("Failed to connect to database", "err", err)
		os.Exit(1)
	}

	sqlDB, err := db.DB()
	if err != nil {
		slog.Error("failed to get underlying sql.DB:", "err", err)
		os.Exit(1)
	}
	defer sqlDB.Close()

	// 依賴注入：組裝 Repository, Service, Handler
	userRepo := gormimpl.NewUserRepository(db)
	unipileRepo := gormimpl.NewUnipileRepository(db)

	userSvc := service.NewUserService(userRepo)
	authSvc := service.NewAuthService(cfg.Server.JWTSecret)
	unipileSvc := service.NewUnipileService(unipileRepo)

	userHdl := handler.NewUserHandler(userSvc, authSvc)
	unipileHdl := handler.NewUnipileHandler(cfg, unipileSvc)

	// 設定路由
	r := handler.SetupRouter(cfg, userHdl, unipileHdl)
	slog.Info("Router setup complete")

	// 5. 將 Gin 路由器包裝在標準的 http.Server 中
	serverAddr := fmt.Sprintf(":%d", cfg.Server.Port)
	srv := &http.Server{
		Addr:    serverAddr,
		Handler: r,
	}

	// 6. 啟動伺服器 (在一個 goroutine 中非同步啟動)
	go func() {
		slog.Info(fmt.Sprintf("Server starting on port %s", serverAddr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server failed to start", "err", err)
			os.Exit(1)
		}
	}()

	// 7. Graceful Shutdown 邏輯
	// 建立一個 channel 來接收作業系統訊號
	quit := make(chan os.Signal, 1)
	// 監聽 SIGINT (Ctrl+C) 和 SIGTERM (Docker, Kubernetes 等終止訊號)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 阻塞直到接收到訊號
	<-quit
	slog.Warn("Shutting down server...")

	// 8. 執行伺服器關閉
	// 建立一個具有 5 秒超時的 Context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // 確保 context 資源被釋放

	// 嘗試關閉伺服器
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown:", "err", err)
		os.Exit(1)
	}

	slog.Info("Server exiting gracefully.")
}
