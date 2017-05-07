package seeder

import (
	"encoding/json"
	"io/ioutil"

	"github.com/dring1/jwt-oauth/database"
	"github.com/dring1/jwt-oauth/models"
)

type Config struct {
	SeedDataFilePath string
	DbName           string
}

func Seed(db *database.Service, c *Config) error {
	file, err := ioutil.ReadFile(c.SeedDataFilePath)
	if err != nil {
		return err
	}

	var cocktails []models.Cocktail
	err = json.Unmarshal(file, &cocktails)
	if err != nil {
		return err
	}

	session := db.Session.Copy()
	defer session.Close()

	collection := session.DB(c.DbName).C("cocktails")

	for _, cocktail := range cocktails {
		err := collection.Insert(cocktail)
		if err != nil {
			return err
		}
	}

	return nil

}
