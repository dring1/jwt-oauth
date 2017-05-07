package database

import mgo "gopkg.in/mgo.v2"

type Service struct {
	Session *mgo.Session
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
	session, err := mgo.Dial("localhost")
	if err != nil {
		return nil, err
	}
	session.SetMode(mgo.Monotonic, true)
	mg := &Service{
		Session: session,
	}

	err = Create(mg.Session, c)
	if err != nil {
		return nil, err
	}

	return mg, nil

}

func Create(s *mgo.Session, config *Config) error {
	session := s.Copy()
	defer session.Close()
	c := session.DB(config.DbName).C("cocktails")

	index := mgo.Index{
		Key:        []string{"name"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	return c.EnsureIndex(index)
}
