package entity

import "github.com/google/uuid"

// Group represents a group of users (project, organization, public)
type Group struct {
	ID          uuid.UUID `json:"id"`
	Type        string    `json:"type"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}
