package routes

import (
	"context"
	"net/http"

	"github.com/dring1/jwt-oauth/graphql"
	jsonresponder "github.com/dring1/jwt-oauth/jsonResponder"
	"github.com/dring1/jwt-oauth/lib/contextkeys"
)

type GraphqlRoute struct {
	Route
	GraphqlService graphql.Service       `service:"GraphqlService"`
	JsonResponder  jsonresponder.Service `service:"JsonResponder"`
}

func (rql *GraphqlRoute) CompileRoute() (*Route, error) {
	fn := func(w http.ResponseWriter, r *http.Request) {

		// if its a mutation lol
		// query, err := ioutil.ReadAll(r.Body)
		query := r.URL.Query()["query"][0]

		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusBadRequest)
		// 	return
		// }
		result := rql.GraphqlService.ExecuteQuery(string(query))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		ctx := context.WithValue(r.Context(), contextkeys.Value, result)
		r = r.WithContext(ctx)
		rql.JsonResponder.Respond(w, r)
		return

	}
	rql.Handler = http.HandlerFunc(fn)
	return &rql.Route, nil
}
