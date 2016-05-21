package models

// User -
type User struct {
	// ID    uuid.UUID `gorm:"primary_key;type:uuid"`
	Email string `gorm:"type:varchar(100);primary_key;unique"`
}
