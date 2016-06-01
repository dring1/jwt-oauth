package controllers

import (
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
