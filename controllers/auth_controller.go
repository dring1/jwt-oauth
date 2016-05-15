package controllers

import (
	"encoding/json"
	"net/http"

	"golang.org/x/net/context"

	"github.com/dghubble/ctxh"
	"github.com/dring1/orm/models"
	"github.com/dring1/orm/services"
)

func Login() ctxh.ContextHandler {
	fn := func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		// githubUser, err := github.UserFromContext(ctx)
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }
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
	return ctxh.ContextHandlerFunc(fn)
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
