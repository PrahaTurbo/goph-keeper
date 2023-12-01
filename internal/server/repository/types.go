package repository

import "time"

// Secret is a struct that represents a Secret created by a User.
type Secret struct {
	CreatedAt time.Time
	Type      string
	Content   []byte
	MetaData  []byte
	ID        int
	UserID    int
}
