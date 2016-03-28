package services

import (
	"github.com/dring1/orm/models"
	"github.com/dring1/orm/reddit"
	"github.com/jinzhu/gorm"
)

type DataBase struct {
	gm *gorm.DB
}

func NewDataBase() (*gorm.DB, error) {
	db, err := gorm.Open("postgres", "user=postgres sslmode=disable")
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&model.Song{})
	return db, nil
}

func (db *DataBase) InsertSubmissions(subs []*reddit.Submission) error {
	for _, s := range subs {
		db.gm.Create(s)
	}
	return nil
}
