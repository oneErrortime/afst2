package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"library-go/database"
	"library-go/models"
)

type BookHandler struct{}

func NewBookHandler() *BookHandler {
	return &BookHandler{}
}

// GetBooks retrieves all books
func (h *BookHandler) GetBooks(c *gin.Context) {
	var books []models.Book
	
	// Get pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	offset := (page - 1) * limit

	if err := database.DB.Offset(offset).Limit(limit).Find(&books).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch books"})
		return
	}

	c.JSON(http.StatusOK, books)
}

// GetBook retrieves a specific book by ID
func (h *BookHandler) GetBook(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	var book models.Book
	if err := database.DB.First(&book, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	c.JSON(http.StatusOK, book)
}

// CreateBook creates a new book
func (h *BookHandler) CreateBook(c *gin.Context) {
	var input models.Book
	
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if ISBN already exists (if provided)
	if input.ISBN != nil {
		var existingBook models.Book
		if err := database.DB.Where("isbn = ?", *input.ISBN).First(&existingBook).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Book with this ISBN already exists"})
			return
		}
	}

	if err := database.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create book"})
		return
	}

	c.JSON(http.StatusCreated, input)
}

// UpdateBook updates an existing book
func (h *BookHandler) UpdateBook(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	var book models.Book
	if err := database.DB.First(&book, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	var input struct {
		Title       *string  `json:"title"`
		Author      *string  `json:"author"`
		Year        *int     `json:"year"`
		ISBN        *string  `json:"isbn"`
		Copies      *int     `json:"copies"`
		Description *string  `json:"description"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields if provided
	if input.Title != nil {
		book.Title = *input.Title
	}
	if input.Author != nil {
		book.Author = *input.Author
	}
	if input.Year != nil {
		book.Year = input.Year
	}
	if input.ISBN != nil {
		// Check if ISBN is being updated and if it already exists
		if *input.ISBN != *book.ISBN {
			var existingBook models.Book
			if err := database.DB.Where("isbn = ?", *input.ISBN).First(&existingBook).Error; err == nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Book with this ISBN already exists"})
				return
			}
		}
		book.ISBN = input.ISBN
	}
	if input.Copies != nil {
		book.Copies = *input.Copies
	}
	if input.Description != nil {
		book.Description = input.Description
	}

	// Validate the updated book
	if err := book.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Save(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update book"})
		return
	}

	c.JSON(http.StatusOK, book)
}

// DeleteBook deletes a book
func (h *BookHandler) DeleteBook(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	var book models.Book
	if err := database.DB.First(&book, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	// Check if book is currently borrowed
	var activeBorrows int64
	database.DB.Model(&models.Borrow{}).Where("book_id = ? AND is_returned = ?", uint(id), false).Count(&activeBorrows)

	if activeBorrows > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete book that is currently borrowed"})
		return
	}

	if err := database.DB.Delete(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete book"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
}