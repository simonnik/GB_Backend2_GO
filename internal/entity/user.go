package entity

import "github.com/google/uuid"

// User represents an employee
type User struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}
