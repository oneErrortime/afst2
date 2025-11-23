package models

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Book struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title" gorm:"not null"`
	Author      string    `json:"author" gorm:"not null"`
	Year        *int      `json:"year,omitempty"` // Publication year, optional
	ISBN        *string   `json:"isbn,omitempty" gorm:"unique"` // Unique ISBN, optional
	Copies      int       `json:"copies" gorm:"default:1"` // Number of copies available
	Description *string   `json:"description,omitempty"` // Added for the second migration
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	
	// Relationship with borrows
	Borrows []Borrow `json:"-" gorm:"foreignKey:BookID"`
}

// Validate validates the book data
func (b *Book) Validate() error {
	// Validate title
	if b.Title == "" || strings.TrimSpace(b.Title) == "" {
		return errors.New("title cannot be empty")
	}
	if len(b.Title) > 500 {
		return errors.New("title must be less than 500 characters")
	}
	b.Title = strings.TrimSpace(b.Title)

	// Validate author
	if b.Author == "" || strings.TrimSpace(b.Author) == "" {
		return errors.New("author cannot be empty")
	}
	if len(b.Author) > 200 {
		return errors.New("author name must be less than 200 characters")
	}
	b.Author = strings.TrimSpace(b.Author)

	// Validate year if provided
	if b.Year != nil {
		currentYear := time.Now().Year()
		if *b.Year < 1000 || *b.Year > currentYear+10 { // Allow up to 10 years in the future
			return errors.New("year must be between 1000 and " + fmt.Sprintf("%d", currentYear+10))
		}
	}

	// Validate ISBN if provided
	if b.ISBN != nil {
		if err := validateISBN(*b.ISBN); err != nil {
			return err
		}
	}

	// Validate copies
	if b.Copies < 0 {
		return errors.New("copies cannot be negative")
	}

	// Validate description if provided
	if b.Description != nil && len(*b.Description) > 2000 {
		return errors.New("description must be less than 2000 characters")
	}

	return nil
}

// validateISBN validates the ISBN format
func validateISBN(isbn string) error {
	// Remove hyphens and spaces
	cleanISBN := regexp.MustCompile(`[-\s]`).ReplaceAllString(isbn, "")
	
	if len(cleanISBN) != 10 && len(cleanISBN) != 13 {
		return errors.New("ISBN must be either 10 or 13 characters long")
	}

	// Check if all characters are digits (for basic validation)
	if !isNumeric(cleanISBN) && !(len(cleanISBN) == 10 && isNumeric(cleanISBN[:9]) && (cleanISBN[9] == 'X' || cleanISBN[9] == 'x' || isNumeric(cleanISBN[9:10]))) {
		return errors.New("invalid ISBN format")
	}

	return nil
}

// isNumeric checks if a string contains only numeric characters
func isNumeric(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// BeforeCreate is a GORM hook that runs before creating a book
func (b *Book) BeforeCreate(tx *gorm.DB) error {
	return b.Validate()
}

// BeforeUpdate is a GORM hook that runs before updating a book
func (b *Book) BeforeUpdate(tx *gorm.DB) error {
	return b.Validate()
}