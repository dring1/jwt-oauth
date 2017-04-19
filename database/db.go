package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Service struct {
	DB *sql.DB
}

type Config struct {
	User     string
	Password string
	Host     string
	Port     int
	DbName   string
	SSL      string
}

func NewService(c *Config) (*Service, error) {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DbName, c.SSL)
	log.Println(connectionString)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &Service{DB: db}, nil
}
