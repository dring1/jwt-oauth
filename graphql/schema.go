package graphql

import (
	"github.com/dring1/jwt-oauth/database"
	"github.com/dring1/jwt-oauth/models"
	"github.com/graphql-go/graphql"
)

type Schema map[string]*graphql.Object

func NewSchema(db *database.Service) graphql.Fields {
	// {cocktail { ingredients : ['rum', 'simple syrup'], tags: ['sweet', 'sour', 'juice', 'summer', 'tiki']}}
	var cocktailType = graphql.NewObject(graphql.ObjectConfig{
		Name: "cocktail",
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"ingredients": &graphql.Field{
				Type: graphql.NewList(graphql.String),
			},
			"instructions": &graphql.Field{
				Type: graphql.NewList(graphql.String),
			},
		},
	})

	var ingredientType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Ingredient",
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type: graphql.NewList(cocktailType),
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return []models.Cocktail{
						models.Cocktail{
							Name:         "Testing",
							Instructions: []string{"Shake baby"},
							Ingredients:  []string{"Whiskey"},
						},
					}, nil
				},
			},
		},
	})

	// var queryFields = Schema{
	// 	"cocktail": cocktailType,
	// 	"ingredient": ingredientType,
	// }

	// var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	// 	Name:   "RootQuery",
	// 	Fields: queryFields,
	// })

	return graphql.Fields{
		// "cocktail":   cocktailType,
		"ingredient": &graphql.Field{
			Type:        ingredientType,
			Description: "Get cocktails contaning ingredients",
			// Args: graphql.FieldConfigArgument{
			// 	"ingredient": &graphql.ArgumentConfig{
			// 		Type: graphql.Interface,
			// 	},
			// },

			// Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			// return []models.Cocktail{
			// models.Cocktail{
			// Name:         "Testing",
			// Instructions: []string{"Shake baby"},
			// Ingredients:  []string{"Whiskey"},
			// },
			// }, nil
			// },
		},
		"cocktail": &graphql.Field{
			Type:        graphql.NewList(cocktailType),
			Description: "Get cocktails ",
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				return []models.Cocktail{
					models.Cocktail{
						Name:         "Testing",
						Instructions: []string{"Shake baby"},
						Ingredients:  []string{"Whiskey"},
					},
				}, nil
			},
		},
	}
}
