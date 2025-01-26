package controllers

import (
	"bookly-api-golang/database"
	"bookly-api-golang/repository"
	"bookly-api-golang/structs"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetAllBook godoc
// @Summary Get all books
// @Description Mendapatkan semua buku
// @Tags Books
// @Produce json
// @Success 200 {object} structs.APIResponse
// @Failure 500 {object} structs.APIResponse
// @Security BearerAuth
// @Router /books [get]
func GetAllBook(c *gin.Context) {
	books, err := repository.GetAllBook(database.DbConnection)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mendapatkan data buku",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": books,
	})
}

// GetBookByID godoc
// @Summary Get book by ID
// @Description Mendapatkan detail buku berdasarkan ID
// @Tags Books
// @Produce json
// @Param id path int true "Book ID"
// @Success 200 {object} structs.Book
// @Failure 400 {object} structs.APIResponse
// @Failure 404 {object} structs.APIResponse
// @Failure 500 {object} structs.APIResponse
// @Security BearerAuth
// @Router /books/{id} [get]
func GetBookByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID tidak valid",
		})
		return
	}

	book, err := repository.GetBook(database.DbConnection, id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Buku tidak ditemukan",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Gagal mendapatkan buku",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": book,
	})
}

// CreateBook godoc
// @Summary Create a new book
// @Description Menambahkan buku baru
// @Tags Books
// @Accept json
// @Produce json
// @Param book body structs.BookInput true "Book body"
// @Success 201 {object} structs.APIResponse
// @Failure 400 {object} structs.APIResponse
// @Failure 500 {object} structs.APIResponse
// @Security BearerAuth
// @Router /books [post]
func CreateBook(c *gin.Context) {
	var input structs.BookInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Data yang dimasukkan tidak valid",
		})
		return
	}

	book := structs.Book{
		Title:       input.Title,
		Description: input.Description,
		ImageURL:    input.ImageURL,
		ReleaseYear: input.ReleaseYear,
		Price:       input.Price,
		TotalPage:   input.TotalPage,
		CategoryID:  input.CategoryID,
		CreatedBy:   input.CreatedBy,
		ModifiedBy:  input.ModifiedBy,
	}

	err := repository.CreateBook(database.DbConnection, book)
	if err != nil {
		if err.Error() == "release_year harus di antara 1980 dan 2024" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal menambahkan buku",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Buku berhasil ditambahkan",
	})
}

// UpdateBook godoc
// @Summary Update a book
// @Description Mengubah data buku
// @Tags Books
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Param book body structs.UpdateBookInput true "Book object"
// @Success 200 {object} structs.APIResponse
// @Failure 400 {object} structs.APIResponse
// @Failure 404 {object} structs.APIResponse
// @Failure 500 {object} structs.APIResponse
// @Security BearerAuth
// @Router /books/{id} [put]
func UpdateBook(c *gin.Context) {
	var input structs.UpdateBookInput
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID tidak valid",
		})
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Data yang dimasukkan tidak valid",
		})
		return
	}

	book := structs.Book{
		ID:          id,
		Title:       input.Title,
		Description: input.Description,
		ImageURL:    input.ImageURL,
		ReleaseYear: input.ReleaseYear,
		Price:       input.Price,
		TotalPage:   input.TotalPage,
		CategoryID:  input.CategoryID,
		ModifiedBy:  input.ModifiedBy,
	}

	err = repository.UpdateBook(database.DbConnection, book)
	if err != nil {
		if err.Error() == "buku tidak ditemukan" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		if err.Error() == "release_year harus di antara 1980 dan 2024" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mengubah buku",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Buku berhasil diubah",
	})
}

// DeleteBook godoc
// @Summary Delete a book
// @Description Menghapus buku
// @Tags Books
// @Produce json
// @Param id path int true "Book ID"
// @Success 200 {object} structs.APIResponse
// @Failure 400 {object} structs.APIResponse
// @Failure 404 {object} structs.APIResponse
// @Failure 500 {object} structs.APIResponse
// @Security BearerAuth
// @Router /books/{id} [delete]
func DeleteBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID tidak valid",
		})
		return
	}

	err = repository.DeleteBook(database.DbConnection, id)
	if err != nil {
		if err.Error() == "buku tidak ditemukan" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal menghapus buku",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Buku berhasil dihapus",
	})
}