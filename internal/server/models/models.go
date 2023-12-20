// Package models contains definitions for service-layer structures.
package models

import "time"

// User is a struct that represents a User in the system.
type User struct {
	Login        string
	PasswordHash string
	ID           int
}

// Secret is a struct that represents a Secret created by a User.
type Secret struct {
	CreatedAt time.Time
	Type      string
	Content   string
	MetaData  string
	ID        int
	UserID    int
}
