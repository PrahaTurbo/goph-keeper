package repository

import "time"

type Secret struct {
	CreatedAt time.Time
	Type      string
	Content   []byte
	MetaData  []byte
	ID        int
	UserID    int
}
