package handler

import (
	"chatsheet/config"
	"chatsheet/internal/middleware"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
)

// @title Chatsheet API
// @version 1.0
// @description An app for connecting user's LinkedIn account by Unipile's native authentication。
// @host localhost:8080
// @BasePath /
func SetupRouter(cfg *config.AppConfig, userHdl *UserHandler, unipileHdl *UnipileHandler) *gin.Engine {
	r := gin.Default()

	// CORS 設定
	r.Use(middleware.CORSMiddleware(cfg.App.FrontendURL))

	// **Swagger 文件路由** (完成後取消註釋)
	// r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	authApi := r.Group("/auth")
	{
		authApi.POST("/signup", userHdl.Signup)
		authApi.POST("/login", userHdl.Login)
	}

	// 路由群組
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware(userHdl.AuthService))
	{
		unipileApi := api.Group("/unipile")
		{
			unipileApi.GET("/", unipileHdl.List)
			unipileApi.POST("/linkedin/basic", unipileHdl.LinkedInBasic)
			unipileApi.POST("/linkedin/cookie", unipileHdl.LinkedInCookie)
			unipileApi.POST("/linkedin/checkpoint", unipileHdl.Checkpoint)
		}
	}

	const staticPath = "web/myapp/dist"
	r.Static("/assets", path.Join(staticPath, "assets"))

	// 3. 處理 SPA 路由回退 (Fallback)
	// 這個萬用路由將處理所有不屬於 /api 的 GET 請求
	r.NoRoute(func(c *gin.Context) {
		// 確保不是發往 /api/ 的請求
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			// 如果是 API 請求，但找不到對應路由，返回 404
			c.JSON(http.StatusNotFound, gin.H{"error": "API 路由未找到"})
			return
		}

		// 嘗試服務根目錄的 index.html
		indexPath := path.Join(staticPath, "index.html")

		// 檢查 index.html 是否存在
		if _, err := os.Stat(indexPath); os.IsNotExist(err) {
			c.String(http.StatusNotFound, "找不到前端 index.html 檔案，請先執行 npm run build。")
			return
		}

		// 處理所有前端路由 (如 /accounts, /login)
		// Gin 將 index.html 的內容發送給瀏覽器，由 Svelte 路由器處理前端路徑。
		c.File(indexPath)
	})

	return r
}
