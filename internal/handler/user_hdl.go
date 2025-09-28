package handler

import (
	"chatsheet/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	userService *service.UserService
	AuthService *service.AuthService
	validate    *validator.Validate
}

func NewUserHandler(userSvc *service.UserService, authSvc *service.AuthService) *UserHandler {
	return &UserHandler{
		userService: userSvc,
		AuthService: authSvc,
		validate:    validator.New(),
	}
}

// @Summary 註冊新使用者
// @Description 使用者註冊一個新帳號
// @Tags users
// @Accept json
// @Produce json
// @Param request body SignupRequest true "註冊請求"
// @Success 201 {object} StandardResponse{data=model.User}
// @Failure 400 {object} ErrorResponse "無效的請求"
// @Failure 500 {object} ErrorResponse "內部伺服器錯誤"
// @Router /signup [post]
func (h *UserHandler) Signup(c *gin.Context) {
	var req SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	user, err := h.userService.Create(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "SignUp success",
		"user":    user,
	})
}

// @Summary 使用者登入
// @Description 使用者憑 E-mail 和密碼登入
// @Tags users
// @Accept json
// @Produce json
// @Param request body LoginRequest true "登入請求"
// @Success 200 {object} StandardResponse{data=object{token=string}}
// @Failure 400 {object} ErrorResponse "無效的請求"
// @Failure 401 {object} ErrorResponse "憑證無效"
// @Failure 500 {object} ErrorResponse "內部伺服器錯誤"
// @Router /login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	user, err := h.userService.Authenticate(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		// 建議回傳通用的錯誤訊息，避免暴露使用者不存在等細節
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	token, err := h.AuthService.GenerateToken(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login success",
		"token":   token,
	})
}
