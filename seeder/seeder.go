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

// read the json :
//`{
//    cocktails: [
//    {
//        name,
//        ing,
//        ins
//    }
//    ]
//}`
// set some value in db to an array of cocktails ?
//func Seed(Db *databse.Service, c *Config) error {
//    file, err := ioutil.ReadFile(c.SeedDataFilePath)
//    if err != nil {
//        return err
//    }

//    var cocktails []models.CocktailDS
//    err = json.Unmarshal(file, &cocktails)
//    if err != nil {
//        return err
//    }

//    //before seeding delete everything in the table.
//    // move seeder to its own app?
//    sqlDeleteStatement := "DELETE FROM cocktails"
//    _, err = c.Db.DB.Exec(sqlDeleteStatement)
//    if err != nil {
//        return err
//    }
//    for _, v := range cocktails {

//        sqlStatement := `
//        INSERT INTO cocktails (name, ingredients, instructions)
//        VALUES ($1, $2, $3)
//        `
//        _, err := c.Db.DB.Exec(sqlStatement, v.Name, v.Ingredients, v.Instructions)
//        if err != nil {
//            log.Println("Error on cocktail: %s", v.Name)
//            return err
//        }
//    }
//    return nil
//}

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
