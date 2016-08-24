package controllers

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/dring1/jwt-oauth/services"
)

type Controller interface {
	http.Handler
}

type C struct {
	Route   string
	Methods []string
}

func New(services ...services.Service) []Controller {
	ctrls := []Controller{
		&HellController{
			C: C{
				Route:   "/hello",
				Methods: []string{"GET"},
			},
			// DatabaseService: nil,
			// CacheService:    nil,
		},
	}

	for _, ctrl := range ctrls {
		val := reflect.ValueOf(ctrl).Elem()
		for index := 0; index < val.NumField(); index++ {
			for _, s := range services {
				// fmt.Println(val.Type().Field(index).Type.String(), reflect.TypeOf(s).String())
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
