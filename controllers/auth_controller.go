package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/dring1/orm/models"
)

func Login(w http.ResponseWriter, r *http.Request) {
	requestUser := new(models.User)
	json.NewDecoder(r.Body).Decode(&requestUser)

  responseStatues, token := services.Login(requestUser)
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(responseStatues)
  w.Write(token)
}

func RefreshToken(w http.ResponseWriter, r *http.Request){
	requestUser := new(models.User)
	json.NewDecoder(r.Body).Decode(&requestUser)

  w.Header().Set("Content-Type", "application/json")
  w.Write(services.RefreshToken(requestUser))
}

func Logout(w http.ResponseWriter, r *http.Request){
  err := services.Logout(r)
  w.Header().Set("Content-Type", "application/json")
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
  } else {
    w.WriteHeader(http.StatusOK)
  }
}
