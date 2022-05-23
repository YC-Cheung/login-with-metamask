package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	PublicAddress string
	Username      string
	Nonce         string
}
