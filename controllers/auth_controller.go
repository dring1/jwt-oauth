package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/context"

	"github.com/dghubble/ctxh"
	"github.com/dghubble/gologin/github"
	"github.com/dring1/jwt-oauth/cache"
	"github.com/dring1/jwt-oauth/lib/errors"
)

type AuthController struct {
	CacheService *cache.CacheService
}

func Login(callback func(*model.User)) ctxh.ContextHandler {
	fn := func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		githubUser, err := github.UserFromContext(ctx)
		log.Println("Github user", *githubUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		u, err := FindUser(*githubUser.Email)
		if err != nil {
			if err.Error() == errors.RecordNotFound {
				log.Println("User does not exists", err)
				log.Printf("Creating new user... %s", *githubUser.Email)

				user, err := CreateUser(*githubUser.Email)
				if err != nil {
					w.WriteHeader(500)
					w.Write([]byte("shit!"))
					return
				}
				u = user
				w.WriteHeader(201)
			}
		}
		token, err := services.Login(*githubUser.Email)
		log.Println("I AM HERE")
		// If user already exists - update last logged in
		// If User does not exist, create it
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error occurred during login"))
			return
		}
		// insert token into redis
		// TODO: `u`
		fmt.Println(u.Email, token)
		err = services.Cache().Set(u.Email, token, 5*time.Minute).Err()
		if err != nil {
			log.Fatal(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(token)
		callback(u)
	}
	return ctxh.ContextHandlerFunc(fn)
}

func RefreshToken(w http.ResponseWriter, r *http.Request) {
	requestUser := new(model.User)
	json.NewDecoder(r.Body).Decode(&requestUser)

	w.Header().Set("Content-Type", "application/json")
	token, err := services.RefreshToken(requestUser.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(token)
}

func Logout(w http.ResponseWriter, r *http.Request) {

	requestUser := new(model.User)
	json.NewDecoder(r.Body).Decode(&requestUser)
	err := services.Logout(r)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
