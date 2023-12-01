package models

import "time"

type User struct {
	Login        string
	PasswordHash string
	ID           int
}

type Secret struct {
	CreatedAt time.Time
	Type      string
	Content   string
	MetaData  string
	ID        int
	UserID    int
}
