package handler

import (
	"fmt"
	"log"
	"movie-reck/controller"
	"movie-reck/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ReckHandler struct {
	reckController *controller.ReckController
}

func NewReckHandler(reckController *controller.ReckController) *ReckHandler {
	return &ReckHandler{reckController: reckController}
}

func (h *ReckHandler) GetRecommendations(c *gin.Context) {

	var req model.RecommendationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters"})
		return
	}
	log.Printf("--- handler req %s", req)

	movies, err := h.reckController.GetRecommendations(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": fmt.Sprintf("failed to fetch recommendations: %s", err.Error())})
		return
	}
	c.JSON(http.StatusOK, gin.H{"movies": movies})
}
