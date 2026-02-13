package models

import (
	"time"

	"github.com/google/uuid"
)

type BaseModel struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

func (b *BaseModel) SetID() {
	b.ID = generateUniqueID()
}
func generateUniqueID() string {
	return uuid.NewString()
}
