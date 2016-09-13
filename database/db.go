package database

import (
	"github.com/dring1/jwt-oauth/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type DatabaseService struct {
	*gorm.DB
}

func New() (*DatabaseService, error) {
	db, err := gorm.Open("postgres", "user=postgres sslmode=disable")
	if err != nil {
		return nil, err
	}
	d := &DatabaseService{db}
	d.AutoMigrate(&models.User{})
	return d, nil
}

// func (db *DatabaseService) RegisterService(s *[]services.Service) {
//
// }

// func (db *DataBase) InsertSubmissions(subs []*reddit.Submission) error {
// 	for _, s := range subs {
// 		db.gm.Create(s)
// 	}
// 	return nil
// }
