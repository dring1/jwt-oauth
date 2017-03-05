package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/dring1/jwt-oauth/lib/contextkeys"
)

type JSONResponse struct {
	Value interface{} `json:"value"`
	Error interface{} `json:"error"`
}

// TODO: Possible - a before and after middleware hook, with this being the last
func JsonResponseHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		jsonResponse := &JSONResponse{
			Value: ctx.Value(contextkeys.Value),
			Error: ctx.Value(contextkeys.Error),
		}
		json.NewEncoder(w).Encode(jsonResponse)
		return
	}
	return http.HandlerFunc(fn)
}
