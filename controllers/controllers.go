package controllers

import (
	"reflect"

	"github.com/dring1/jwt-oauth/services"
)

type Controller interface {
	Create(interface{}) error
	Read(interface{}) error
	Update(interface{}) error
	Delete(interface{}) error
}

func New(services ...services.Service) []Controller {
	ctrls := []Controller{
		&HelloController{},
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
	}
	return ctrls
}
