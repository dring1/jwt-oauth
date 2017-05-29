package graphql

import (
	"github.com/dring1/jwt-oauth/database"
	"github.com/dring1/jwt-oauth/models"
)

func GetByIngredients(db *database.Service, ingredients []string) (error, []models.Cocktail) {
	return nil, []models.Cocktail{}
}
