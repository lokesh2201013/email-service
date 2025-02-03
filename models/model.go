package models

import ("gorm.io/gorm"
		"time")

type Sender struct {
	gorm.Model
	AdminName string  `json:"admin_name" gorm:"not null"`
	Email    string `json:"email" gorm:"unique;not null"`
	SMTPHost string `json:"smtp_host" gorm:"not null"`
	SMTPPort int    `json:"smtp_port" gorm:"not null"`
	Username string `json:"username" gorm:"not null"`
	AppPassword string `json:"password" gorm:"not null"`
	Verified bool   `json:"verified" gorm:"default:false"`
}


type Template struct {
	ID      uint   `gorm:"primaryKey"`
	Name    string `gorm:"unique"`
	Subject string
	Body    string
	Format  string `gorm:"default:'text'"`
}

type User struct {
    ID        uint   `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
    Username string `json:"username" gorm:"unique;not null"`
    Password string `json:"password" gorm:"not null"`
}