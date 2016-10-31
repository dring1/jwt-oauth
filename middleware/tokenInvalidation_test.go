package middleware

import (
	"reflect"
	"testing"

	"github.com/dring1/jwt-oauth/token"
)

func TestNewTokenInvalidationMiddleware(t *testing.T) {
	type args struct {
		tokenService token.Service
	}
	tests := []struct {
		name string
		args args
		want Middleware
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if got := NewTokenInvalidationMiddleware(tt.args.tokenService); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. NewTokenInvalidationMiddleware() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
