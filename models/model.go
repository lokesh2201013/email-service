package models

import "gorm.io/gorm"

type Sender struct {
	gorm.Model
	Email    string `json:"email" gorm:"unique;not null"`
	SMTPHost string `json:"smtp_host" gorm:"not null"`
	SMTPPort int    `json:"smtp_port" gorm:"not null"`
	Username string `json:"username" gorm:"not null"`
	Password string `json:"password" gorm:"not null"`
	Verified bool   `json:"verified" gorm:"default:false"`
}


type Template struct {
	ID      uint   `gorm:"primaryKey"`
	Name    string `gorm:"unique"`
	Subject string
	Body    string
	Format  string `gorm:"default:'text'"`
}

