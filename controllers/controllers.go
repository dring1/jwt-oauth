package controllers

import (
	"net/http"
	"reflect"

	"github.com/dring1/jwt-oauth/services"
)

type Controller interface {
	http.Handler
	Route() string
	Methods() []string
}

func New(services ...services.Service) []Controller {
	ctrls := []Controller{
		&HellController{},
	}

	for _, ctrl := range ctrls {
		val := reflect.ValueOf(ctrl).Elem()
		for index := 0; index < val.NumField(); index++ {
			for s := range services {
				if val.Type().Field(index).Name == reflect.TypeOf(s).String() {
					// set the value of the field to the service
				}
			}
		}
	}
	return ctrls
}
