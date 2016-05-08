package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/dring1/orm/models"
	"github.com/dring1/orm/services"
)

func Login(w http.ResponseWriter, r *http.Request) {
	requestUser := new(models.User)
	json.NewDecoder(r.Body).Decode(&requestUser)

	token, err := services.Login(requestUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(token)
}

func RefreshToken(w http.ResponseWriter, r *http.Request) {
	requestUser := new(models.User)
	json.NewDecoder(r.Body).Decode(&requestUser)

	w.Header().Set("Content-Type", "application/json")
	token, err := services.RefreshToken(requestUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(token)
}

func Logout(w http.ResponseWriter, r *http.Request) {

	requestUser := new(models.User)
	json.NewDecoder(r.Body).Decode(&requestUser)
	err := services.Logout(r)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
