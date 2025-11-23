package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Borrow struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	BookID      uint      `json:"book_id" gorm:"not null"`
	ReaderID    uint      `json:"reader_id" gorm:"not null"`
	BorrowedAt  time.Time `json:"borrowed_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	ReturnedAt  *time.Time `json:"returned_at,omitempty"` // nil if not returned yet
	IsReturned  bool      `json:"is_returned" gorm:"default:false"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	
	// Relationships
	Book   Book   `json:"book" gorm:"foreignKey:BookID"`
	Reader Reader `json:"reader" gorm:"foreignKey:ReaderID"`
}

// Validate validates the borrow data
func (b *Borrow) Validate() error {
	if b.BookID == 0 {
		return errors.New("book_id is required")
	}
	
	if b.ReaderID == 0 {
		return errors.New("reader_id is required")
	}
	
	// Can't return a book before borrowing it
	if b.ReturnedAt != nil && b.BorrowedAt.After(*b.ReturnedAt) {
		return errors.New("return date cannot be before borrow date")
	}
	
	return nil
}

// BeforeCreate is a GORM hook that runs before creating a borrow
func (b *Borrow) BeforeCreate(tx *gorm.DB) error {
	return b.Validate()
}

// BeforeUpdate is a GORM hook that runs before updating a borrow
func (b *Borrow) BeforeUpdate(tx *gorm.DB) error {
	return b.Validate()
}