package models

import "time"

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
	// Dummy implementation for unique ID generation
	return "UN-" + time.Now().Format("20060102150405")
}
