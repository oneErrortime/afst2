package models

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Reader struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	FirstName   string    `json:"first_name" gorm:"not null"`
	LastName    string    `json:"last_name" gorm:"not null"`
	Email       *string   `json:"email,omitempty" gorm:"unique"` // Optional email
	Phone       *string   `json:"phone,omitempty"`              // Optional phone
	Address     *string   `json:"address,omitempty"`            // Optional address
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	
	// Relationship with borrows
	Borrows []Borrow `json:"-" gorm:"foreignKey:ReaderID"`
}

// Validate validates the reader data
func (r *Reader) Validate() error {
	// Validate first name
	if r.FirstName == "" || strings.TrimSpace(r.FirstName) == "" {
		return errors.New("first name cannot be empty")
	}
	if len(r.FirstName) > 100 {
		return errors.New("first name must be less than 100 characters")
	}
	r.FirstName = strings.TrimSpace(r.FirstName)

	// Validate last name
	if r.LastName == "" || strings.TrimSpace(r.LastName) == "" {
		return errors.New("last name cannot be empty")
	}
	if len(r.LastName) > 100 {
		return errors.New("last name must be less than 100 characters")
	}
	r.LastName = strings.TrimSpace(r.LastName)

	// Validate email if provided
	if r.Email != nil {
		if len(*r.Email) > 255 {
			return errors.New("email must be less than 255 characters")
		}
		// Basic email validation
		if !isValidEmail(*r.Email) {
			return errors.New("invalid email format")
		}
	}

	// Validate phone if provided
	if r.Phone != nil {
		if len(*r.Phone) > 20 {
			return errors.New("phone number must be less than 20 characters")
		}
	}

	// Validate address if provided
	if r.Address != nil {
		if len(*r.Address) > 500 {
			return errors.New("address must be less than 500 characters")
		}
	}

	return nil
}

// isValidEmail performs basic email validation
func isValidEmail(email string) bool {
	// Basic check: contains @ and has valid format
	if len(email) < 5 || len(email) > 255 {
		return false
	}
	
	atIndex := -1
	for i, char := range email {
		if char == '@' {
			if atIndex != -1 { // More than one @
				return false
			}
			atIndex = i
		}
	}
	
	if atIndex == -1 || atIndex == 0 || atIndex == len(email)-1 { // No @, or @ at start/end
		return false
	}
	
	// Check if there's a dot after @
	dotFound := false
	for i := atIndex + 1; i < len(email); i++ {
		if email[i] == '.' {
			dotFound = true
			break
		}
	}
	
	return dotFound
}

// BeforeCreate is a GORM hook that runs before creating a reader
func (r *Reader) BeforeCreate(tx *gorm.DB) error {
	return r.Validate()
}

// BeforeUpdate is a GORM hook that runs before updating a reader
func (r *Reader) BeforeUpdate(tx *gorm.DB) error {
	return r.Validate()
}