package models

import "time"

type User struct {
	ID          string
	Email       string
	Password    string
	Name        string
	AccessToken string
	UpdatedAt   time.Time
}
