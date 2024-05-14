package model

import "github.com/google/uuid"

type Product struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"userId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       uint32    `json:"price"`
}
