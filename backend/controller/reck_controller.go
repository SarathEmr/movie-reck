package controller

import (
	"movie-reck/intelligence"
	"movie-reck/model"
	"movie-reck/repo"
)

type ReckController struct {
	movieRepo *repo.MovieRepo
}

func NewReckController(movieRepo *repo.MovieRepo) *ReckController {
	return &ReckController{movieRepo: movieRepo}
}

func (c *ReckController) GetRecommendations(req model.RecommendationRequest) ([]model.Movie, error) {

	// return c.movieRepo.GetAllMovies()

	return intelligence.RecommendMoviesUsingAI(req)
}
