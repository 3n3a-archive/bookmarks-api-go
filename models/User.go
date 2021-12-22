package models

import "github.com/jinzhu/gorm"

// Our User Struct
type User struct {
	gorm.Model
	Name  string
	Email string
}
