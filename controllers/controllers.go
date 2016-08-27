package controllers

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/dring1/jwt-oauth/services"
)

type Controller interface {
	http.Handler

	Route() string
	Methods() []string
}

type C struct {
	Route   string
	Methods []string
}

func New(services ...services.Service) []Controller {
	ctrls := []Controller{
		&HelloController{
			C: C{
				Route:   "/hello",
				Methods: []string{"GET"},
			},
		},
	}

	for _, ctrl := range ctrls {
		val := reflect.ValueOf(ctrl).Elem()
		for index := 0; index < val.NumField(); index++ {
			for _, s := range services {
				if t := val.Type().Field(index).Type.String(); t == reflect.TypeOf(s).String() {
					if x := val.Field(index); x.CanSet() {
						x.Set(reflect.ValueOf(s))
					}
				}
			}
		}
		fmt.Printf("%+v", ctrl)
	}
	return ctrls
}
