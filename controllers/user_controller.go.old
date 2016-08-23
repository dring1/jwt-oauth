package controllers

import "time"

func FindUser(email string) (*model.User, error) {
	u := new(model.User)

	if err := services.Database().Find(u, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func CreateUser(email string) (*model.User, error) {
	var user = &model.User{
		Email:        email,
		LastLoggedIn: time.Now(),
	}

	return user, services.Database().Create(user).Error
}
