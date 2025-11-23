package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"library-go/database"
	"library-go/models"
)

type BorrowHandler struct{}

func NewBorrowHandler() *BorrowHandler {
	return &BorrowHandler{}
}

// GetBorrows retrieves all borrows
func (h *BorrowHandler) GetBorrows(c *gin.Context) {
	var borrows []models.Borrow
	
	// Get pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	offset := (page - 1) * limit

	if err := database.DB.Offset(offset).Limit(limit).Preload("Book").Preload("Reader").Find(&borrows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch borrows"})
		return
	}

	c.JSON(http.StatusOK, borrows)
}

// GetBorrow retrieves a specific borrow by ID
func (h *BorrowHandler) GetBorrow(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid borrow ID"})
		return
	}

	var borrow models.Borrow
	if err := database.DB.Preload("Book").Preload("Reader").First(&borrow, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Borrow not found"})
		return
	}

	c.JSON(http.StatusOK, borrow)
}

// CreateBorrow creates a new borrow record
func (h *BorrowHandler) CreateBorrow(c *gin.Context) {
	var input struct {
		BookID   uint `json:"book_id" binding:"required"`
		ReaderID uint `json:"reader_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if book exists
	var book models.Book
	if err := database.DB.First(&book, input.BookID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	// Check if reader exists
	var reader models.Reader
	if err := database.DB.First(&reader, input.ReaderID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reader not found"})
		return
	}

	// Check if book is available (has copies)
	var activeBorrows int64
	database.DB.Model(&models.Borrow{}).Where("book_id = ? AND is_returned = ?", input.BookID, false).Count(&activeBorrows)

	if int(activeBorrows) >= book.Copies {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No available copies of this book"})
		return
	}

	// Create borrow record
	borrow := models.Borrow{
		BookID:   input.BookID,
		ReaderID: input.ReaderID,
		IsReturned: false,
		BorrowedAt: time.Now(),
	}

	if err := database.DB.Create(&borrow).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create borrow record"})
		return
	}

	// Preload relationships for response
	database.DB.Preload("Book").Preload("Reader").First(&borrow, borrow.ID)

	c.JSON(http.StatusCreated, borrow)
}

// UpdateBorrow updates an existing borrow (typically to mark as returned)
func (h *BorrowHandler) UpdateBorrow(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid borrow ID"})
		return
	}

	var borrow models.Borrow
	if err := database.DB.First(&borrow, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Borrow not found"})
		return
	}

	var input struct {
		IsReturned *bool      `json:"is_returned"`
		ReturnedAt *time.Time `json:"returned_at"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields if provided
	if input.IsReturned != nil {
		borrow.IsReturned = *input.IsReturned
		// If marking as returned and no returned_at time provided, set to now
		if *input.IsReturned && borrow.ReturnedAt == nil {
			now := time.Now()
			borrow.ReturnedAt = &now
		}
	}
	if input.ReturnedAt != nil {
		borrow.ReturnedAt = input.ReturnedAt
		// If returned_at is set, mark as returned
		if borrow.ReturnedAt != nil {
			borrow.IsReturned = true
		}
	}

	if err := database.DB.Save(&borrow).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update borrow"})
		return
	}

	// Preload relationships for response
	database.DB.Preload("Book").Preload("Reader").First(&borrow, borrow.ID)

	c.JSON(http.StatusOK, borrow)
}

// DeleteBorrow deletes a borrow record
func (h *BorrowHandler) DeleteBorrow(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid borrow ID"})
		return
	}

	var borrow models.Borrow
	if err := database.DB.First(&borrow, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Borrow not found"})
		return
	}

	// Cannot delete a borrow that hasn't been returned yet
	if !borrow.IsReturned {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete an active borrow, return the book first"})
		return
	}

	if err := database.DB.Delete(&borrow).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete borrow"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Borrow deleted successfully"})
}