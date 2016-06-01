package services

import (
	"log"
	"sync"

	"github.com/dring1/jwt-oauth/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type database struct {
	*gorm.DB
}

var (
	onceDb sync.Once
	d      *database
)

func Database() *database {
	onceDb.Do(func() {
		db, err := gorm.Open("postgres", "user=postgres sslmode=disable")
		if err != nil {
			log.Fatal("Unable to reach db", err)
		}
		d = &database{db}
		d.AutoMigrate(&models.User{})
	})
	return d
}

// func (db *DataBase) InsertSubmissions(subs []*reddit.Submission) error {
// 	for _, s := range subs {
// 		db.gm.Create(s)
// 	}
// 	return nil
// }
