package entity

import "github.com/google/uuid"

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" `
	Username     string    `gorm:"type:varchar(64);uniqueIndex;not null" `
	Email        string    `gorm:"type:varchar(128);uniqueIndex"`
	PasswordHash string    `gorm:"type:text" `
}
