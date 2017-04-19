package seeder

import "github.com/dring1/jwt-oauth/database"

type Config struct {
	db               *database.Service
	seedDataFilePath string
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
func Seed(c *Config) error {

	return nil
}
