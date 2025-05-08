package models

import (
	"time"
)

// OrganigramTemplate represents a time slot in the organigram
type OrganigramTemplate struct {
	TeacherID uint   `json:"teacher_id"`
	Type      string `json:"type"` // "TP", "TD", or "Cours"
}

// Assignment represents a teacher's assignment to a module
type Assignment struct {
	TeacherID     uint   `json:"teacher_id"`
	ModuleID      uint   `json:"module_id"`
	Type          string `json:"type"` // "TP", "TD", or "Cours"
	HoursAssigned int    `json:"hours_assigned"`
}

// OrganigramInput represents the input structure for organigramme processing
type OrganigramInput struct {
	TimeSlots        []OrganigramTemplate `json:"time_slots"`
	TeacherWishes    []Voeux              `json:"teacher_wishes"`
	AvailableModules []Module             `json:"available_modules"`
	Teachers         []Teacher            `json:"teachers"`
}

// OrganigramOutput represents the output structure for assignments
type OrganigramOutput struct {
	Assignments []Assignment `json:"assignments"`
}

// SavedOrganigram represents a saved version of an organigram
type SavedOrganigram struct {
	ID          uint      `gorm:"primaryKey"`
	Name        string    `gorm:"not null"`
	Description string    `gorm:"type:text"`
	Year        string    `gorm:"not null"`
	IsTemplate  bool      `gorm:"not null;default:false"`
	CreatedAt   time.Time `gorm:"not null"`
	UpdatedAt   time.Time `gorm:"not null"`
	CreatedBy   uint      `gorm:"not null"`           // User ID of the creator
	Data        string    `gorm:"type:json;not null"` // JSON string of OrganigramOutput
}
