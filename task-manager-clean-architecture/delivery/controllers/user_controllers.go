package controllers

import (
	"net/http"

	"task-manager-api/domains"
	"task-manager-api/infrastructure"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userUsecase domains.UserUsecase
}

func NewUserController(userUsecase domains.UserUsecase) *UserController {
	return &UserController{
		userUsecase: userUsecase,
	}
}

// Register handles user registration
func (uc *UserController) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Role     string `json:"role" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := uc.userUsecase.RegisterUser(c.Request.Context(), req.Username, req.Password, req.Role)
	if err != nil {
		if err.Error() == "username already exists" {
			c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// Login handles user login
func (uc *UserController) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, authenticated := uc.userUsecase.AuthenticateUser(c.Request.Context(), req.Username, req.Password)
	if !authenticated {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	token, err := infrastructure.GenerateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
