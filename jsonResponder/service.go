package jsonresponder

import (
	"encoding/json"
	"net/http"

	"github.com/dring1/jwt-oauth/lib/contextkeys"
)

type Service interface {
	Respond(w http.ResponseWriter, r *http.Request)
}

type JsonResponder struct{}

type JSONResponse struct {
	Value interface{} `json:"value"`
	Error interface{} `json:"error"`
}

func NewJsonResponder() Service {
	return &JsonResponder{}
}

func (j *JsonResponder) Respond(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	jsonResponse := &JSONResponse{
		Value: ctx.Value(contextkeys.Value),
		Error: ctx.Value(contextkeys.Error),
	}
	w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(jsonResponse)
	_ = err
}
