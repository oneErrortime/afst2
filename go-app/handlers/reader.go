package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"library-go/database"
	"library-go/models"
)

type ReaderHandler struct{}

func NewReaderHandler() *ReaderHandler {
	return &ReaderHandler{}
}

// GetReaders retrieves all readers
func (h *ReaderHandler) GetReaders(c *gin.Context) {
	var readers []models.Reader
	
	// Get pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	offset := (page - 1) * limit

	if err := database.DB.Offset(offset).Limit(limit).Find(&readers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch readers"})
		return
	}

	c.JSON(http.StatusOK, readers)
}

// GetReader retrieves a specific reader by ID
func (h *ReaderHandler) GetReader(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reader ID"})
		return
	}

	var reader models.Reader
	if err := database.DB.First(&reader, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reader not found"})
		return
	}

	c.JSON(http.StatusOK, reader)
}

// CreateReader creates a new reader
func (h *ReaderHandler) CreateReader(c *gin.Context) {
	var input models.Reader
	
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if email already exists (if provided)
	if input.Email != nil {
		var existingReader models.Reader
		if err := database.DB.Where("email = ?", *input.Email).First(&existingReader).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Reader with this email already exists"})
			return
		}
	}

	if err := database.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create reader"})
		return
	}

	c.JSON(http.StatusCreated, input)
}

// UpdateReader updates an existing reader
func (h *ReaderHandler) UpdateReader(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reader ID"})
		return
	}

	var reader models.Reader
	if err := database.DB.First(&reader, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reader not found"})
		return
	}

	var input struct {
		FirstName *string `json:"first_name"`
		LastName  *string `json:"last_name"`
		Email     *string `json:"email"`
		Phone     *string `json:"phone"`
		Address   *string `json:"address"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields if provided
	if input.FirstName != nil {
		reader.FirstName = *input.FirstName
	}
	if input.LastName != nil {
		reader.LastName = *input.LastName
	}
	if input.Email != nil {
		// Check if email is being updated and if it already exists
		if *input.Email != *reader.Email {
			var existingReader models.Reader
			if err := database.DB.Where("email = ?", *input.Email).First(&existingReader).Error; err == nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Reader with this email already exists"})
				return
			}
		}
		reader.Email = input.Email
	}
	if input.Phone != nil {
		reader.Phone = input.Phone
	}
	if input.Address != nil {
		reader.Address = input.Address
	}

	// Validate the updated reader
	if err := reader.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Save(&reader).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update reader"})
		return
	}

	c.JSON(http.StatusOK, reader)
}

// DeleteReader deletes a reader
func (h *ReaderHandler) DeleteReader(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reader ID"})
		return
	}

	var reader models.Reader
	if err := database.DB.First(&reader, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reader not found"})
		return
	}

	// Check if reader has active borrows
	var activeBorrows int64
	database.DB.Model(&models.Borrow{}).Where("reader_id = ? AND is_returned = ?", uint(id), false).Count(&activeBorrows)

	if activeBorrows > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete reader that has active borrows"})
		return
	}

	if err := database.DB.Delete(&reader).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete reader"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reader deleted successfully"})
}