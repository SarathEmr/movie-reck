package controller

import "movie-reck/repo"

type LoginController struct {
userRepo *repo.UserRepo
}

func NewLoginController(userRepo *repo.UserRepo) *LoginController {
	return &LoginController{userRepo: userRepo}
}

func (c *LoginController) Login(username, password string) (bool, error) {
	return c.userRepo.Authenticate(username, password)
}