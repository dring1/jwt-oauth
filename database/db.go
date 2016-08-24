package database

import (
	"log"

	"github.com/dring1/jwt-oauth/models"
	"github.com/dring1/jwt-oauth/services"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type DatabaseService struct {
	*gorm.DB
}

func NewDatabaseService() *DatabaseService {
	db, err := gorm.Open("postgres", "user=postgres sslmode=disable")
	if err != nil {
		log.Fatal("Unable to reach db", err)
	}
	d := &DatabaseService{db}
	d.AutoMigrate(&models.User{})
	return d
}

func (db *DatabaseService) RegisterService(s *[]services.Service) {

}

// func (db *DataBase) InsertSubmissions(subs []*reddit.Submission) error {
// 	for _, s := range subs {
// 		db.gm.Create(s)
// 	}
// 	return nil
// }
