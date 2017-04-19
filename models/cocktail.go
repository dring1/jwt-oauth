package models

type Cocktail struct {
	Name         string   `json:"name"`
	Instructions []string `json:"instructions"`
	Ingredients  []string `json:"ingredients"`
}
