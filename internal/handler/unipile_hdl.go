package handler

import (
	"net/http"

	"chatsheet/config"
	"chatsheet/internal/service"
	"chatsheet/internal/unipile"

	"github.com/gin-gonic/gin"
)

type UnipileHandler struct {
	cfg        *config.AppConfig
	unipileSvc *service.UnipileService
}

func NewUnipileHandler(cfg *config.AppConfig, unipileSvc *service.UnipileService) *UnipileHandler {
	return &UnipileHandler{
		cfg:        cfg,
		unipileSvc: unipileSvc,
	}
}

// UnipileLoginRequest 處理 Username/Password 登入請求
type UnipileLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UnipileCookieRequest 處理 Cookie 登入請求
type UnipileCookieRequest struct {
	AccessToken string `json:"access_token" binding:"required"`
	UserAgent   string `json:"user_agent" binding:"required"`
}

// UnipileCheckpointRequest 處理 Checkpoint 請求
type UnipileCheckpointRequest struct {
	AccountID string `json:"account_id" binding:"required"` // Intent ID
	Code      string `json:"code" binding:"required"`       // 2FA/OTP/Phone Number
}

// handleUnipileResponse 封裝 Unipile 響應的處理邏輯
// 它負責檢查是否為 Checkpoint，並將 account_id 儲存到 session 或返回給前端。
func handleUnipileResponse(c *gin.Context, status int, response *unipile.CheckpointResponse, svc *service.UnipileService, userEmail string) {
	if status == http.StatusAccepted { // 202 Accepted, Checkpoint
		if response.Object == "Checkpoint" && response.Checkpoint != nil {
			// **TODO: 儲存 AccountID 到 Redis/Session**
			// 由於 CheckpointIntent 有 5 分鐘時限，AccountID 必須儲存並與 UserEmail 關聯。
			// 為了簡化，這裡僅返回給前端，讓前端在下一步 Checkpoint 請求中傳回。

			c.JSON(http.StatusAccepted, gin.H{
				"message":         "需要解決 Checkpoint",
				"account_id":      response.AccountID,
				"checkpoint_type": response.Checkpoint.Type,
			})
			return
		}
	}

	// 200 OK - 成功連接
	if response.AccountID != "" {
		_, err := svc.Create(c.Request.Context(), userEmail, "linkedin", response.AccountID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":    "LinkedIn 帳號成功連結",
			"account_id": response.AccountID,
		})
		return
	}

	// 其他非 Checkpoint 的成功響應，通常不應該發生
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Unipile 返回了未預期的成功響應"})
}

// @Summary LinkedInBasic
// @Description 處理 Username/Password 驗證
// @Tags unipile
// @Security BearerAuth
// @Param Authorization header string true "JWT token" default(Bearer <your_JWT_token>)
// @Produce json
// @Success 200 {object} StandardResponse{data=model.UnipileAccount.AccountID} "成功獲取帳號列表"
// @Failure 401 {object} ErrorResponse "未授權"
// @Failure 500 {object} ErrorResponse "內部伺服器錯誤"
// @Router /unipile/linkedin/basic [post]
func (h *UnipileHandler) LinkedInBasic(c *gin.Context) {
	emailAny, ok := c.Get("email") // 從 AuthMiddleware 取得
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req UnipileLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. 構建 Unipile 請求體
	unipileReq := gin.H{
		"provider": "LINKEDIN",
		"username": req.Username,
		"password": req.Password,
	}

	// 2. 呼叫 Unipile API
	var resp unipile.CheckpointResponse
	status, err := unipile.PerformRequest(h.cfg.Unipile, unipile.AccountsEndpoint, unipileReq, &resp)
	if err != nil {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	// 3. 處理響應
	handleUnipileResponse(c, status, &resp, h.unipileSvc, emailAny.(string))
}

// @Summary LinkedInCookie
// @Description 處理 Cookies 驗證
// @Tags unipile
// @Security BearerAuth
// @Param Authorization header string true "JWT token" default(Bearer <your_JWT_token>)
// @Produce json
// @Success 200 {object} StandardResponse{data=model.UnipileAccount.AccountID} "成功獲取帳號列表"
// @Failure 401 {object} ErrorResponse "未授權"
// @Failure 500 {object} ErrorResponse "內部伺服器錯誤"
// @Router /unipile/linkedin/cookie [post]
func (h *UnipileHandler) LinkedInCookie(c *gin.Context) {
	emailAny, ok := c.Get("email") // 從 AuthMiddleware 取得
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req UnipileCookieRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. 構建 Unipile 請求體
	unipileReq := gin.H{
		"provider":     "LINKEDIN",
		"access_token": req.AccessToken,
		"user_agent":   req.UserAgent,
	}

	// 2. 呼叫 Unipile API
	var resp unipile.CheckpointResponse
	status, err := unipile.PerformRequest(h.cfg.Unipile, unipile.AccountsEndpoint, unipileReq, &resp)
	if err != nil {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	// 3. 處理響應
	handleUnipileResponse(c, status, &resp, h.unipileSvc, emailAny.(string))
}

// @Summary SolveCheckpoint
// @Description 處理所有 Checkpoint 解決方案
// @Tags unipile
// @Security BearerAuth
// @Param Authorization header string true "JWT token" default(Bearer <your_JWT_token>)
// @Produce json
// @Success 200 {object} StandardResponse{data=model.UnipileAccount.AccountID} "成功獲取帳號列表"
// @Failure 401 {object} ErrorResponse "未授權"
// @Failure 500 {object} ErrorResponse "內部伺服器錯誤"
// @Router /unipile/linkedin/checkpoint [post]
func (h *UnipileHandler) Checkpoint(c *gin.Context) {
	emailAny, ok := c.Get("email") // 從 AuthMiddleware 取得
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req UnipileCheckpointRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. 構建 Unipile 請求體
	unipileReq := gin.H{
		"provider":   "LINKEDIN",
		"account_id": req.AccountID,
		"code":       req.Code, // 可能是 2FA code 或 Phone Number
	}

	// 2. 呼叫 Unipile API
	var resp unipile.CheckpointResponse
	status, err := unipile.PerformRequest(h.cfg.Unipile, unipile.CheckpointEndpoint, unipileReq, &resp)
	if err != nil {
		// 檢查是否為 408 Timeout 或 400 Bad Request (Intent 銷毀)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	// 3. 處理響應
	handleUnipileResponse(c, status, &resp, h.unipileSvc, emailAny.(string))
}

// @Summary 獲取帳號列表
// @Description 獲取使用者的帳號列表
// @Tags articles
// @Security BearerAuth
// @Param Authorization header string true "JWT token" default(Bearer <your_JWT_token>)
// @Produce json
// @Success 200 {object} StandardResponse{data=[]model.UnipileAccount} "成功獲取帳號列表"
// @Failure 401 {object} ErrorResponse "未授權"
// @Failure 500 {object} ErrorResponse "內部伺服器錯誤"
// @Router /unipile [get]
func (h *UnipileHandler) List(c *gin.Context) {
	emailAny, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	accts, err := h.unipileSvc.ListByEmail(c.Request.Context(), emailAny.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Get success",
		"accounts": accts,
	})
}
