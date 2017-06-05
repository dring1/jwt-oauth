package graphql

import (
	"github.com/dring1/jwt-oauth/database"
	"github.com/graphql-go/graphql"
)

type Service interface {
	ExecuteQuery(string) *graphql.Result
	// Schema() (*graphql.Schema, error)
}

type Gql struct {
	Schema graphql.Schema
	db     database.Service
}

func NewService(db *database.Service) (Service, error) {
	// store the schema in here, hide it from rest of app

	querySchema, err := GenerateSchema(db)
	if err != nil {
		return nil, err
	}
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    querySchema,
		Mutation: querySchema,
	})
	if err != nil {
		return nil, err
	}

	return &Gql{
		Schema: schema,
	}, nil
}

func (gql *Gql) ExecuteQuery(query string) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        gql.Schema,
		RequestString: query,
	})
	// if len(result.Errors) > 0 {
	// 	return
	// }
	return result
}

// schema config
func GenerateSchema(db *database.Service) (*graphql.Object, error) {
	// initialize all types (return map)
	queryFields := NewSchema(db)
	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name:   "RootQuery",
		Fields: queryFields,
	})
	return rootQuery, nil

}
