package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/context"

	"github.com/dghubble/ctxh"
	"github.com/dghubble/gologin/github"
	"github.com/dring1/orm/lib/errors"
	"github.com/dring1/orm/models"
	"github.com/dring1/orm/services"
)

func Login() ctxh.ContextHandler {
	fn := func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		githubUser, err := github.UserFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		u, err := FindUser(*githubUser.Email)
		if err != nil {

			if err.Error() == errors.RecordNotFound {
				log.Println("User does not exists", err)
				log.Printf("Creating new user... %s", *githubUser.Email)
				w.WriteHeader(500)
				w.Write([]byte("could not find user!"))
				return
			}
		}
		if u == nil {
			log.Println("User does not exists")
			w.WriteHeader(500)
			w.Write([]byte("shit!"))
			return
		}

		// requestUser := new(models.User)
		token, err := services.Login(u)
		// If user already exists - update last logged in
		// If User does not exist, create it
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		// val := services.Database().NewRecord(&models.User{ID: uuid.NewV4(), Email: *githubUser.Email})
		// u := &models.User{Email: *githubUser.Email}
		if inserted := services.Database().Create(u).Error; inserted != nil {

		}
		users := []models.User{}
		services.Database().Find(&users)
		for _, u := range users {
			log.Println(u)
		}
		// log.Println("Val", val)
		// if val {
		// }
		// insert token into redis
		services.Cache().Set(u.Email, token, 5*time.Minute)
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
