package controllers

import (
	"time"

	"github.com/dring1/jwt-oauth/models"
	"github.com/dring1/jwt-oauth/services"
)

func FindUser(email string) (*models.User, error) {
	u := new(models.User)

	if err := services.Database().Find(u, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func CreateUser(email string) (*models.User, error) {
	var user = &models.User{
		Email:        email,
		LastLoggedIn: time.Now(),
	}

	return user, services.Database().Create(user).Error
}
