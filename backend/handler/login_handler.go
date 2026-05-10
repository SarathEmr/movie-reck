package handler

import (
	"movie-reck/controller"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginHandler struct {
	loginController *controller.LoginController
}

func NewLoginHandler(loginController *controller.LoginController) *LoginHandler {
	return &LoginHandler{loginController: loginController}
}

func (h *LoginHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authenticated, err := h.loginController.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if authenticated {
		c.JSON(http.StatusOK, gin.H{"message": "Login successful", "user": gin.H{"username": req.Username}})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	}
}