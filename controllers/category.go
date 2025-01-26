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

// GetAllCategory godoc
// @Summary Get all categories
// @Description Mendapatkan semua kategori
// @Tags Categories
// @Produce json
// @Success 200 {object} structs.APIResponse
// @Failure 500 {object} structs.APIResponse
// @Security BearerAuth
// @Router /categories [get]
func GetAllCategory(c *gin.Context) {
	categories, err := repository.GetAllCategory(database.DbConnection)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mendapatkan data kategori",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": categories,
	})
}

// GetCategoryByID godoc
// @Summary Get category by ID
// @Description Mendapatkan detail kategori berdasarkan ID
// @Tags Categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} structs.Category
// @Failure 400 {object} structs.APIResponse
// @Failure 404 {object} structs.APIResponse
// @Failure 500 {object} structs.APIResponse
// @Security BearerAuth
// @Router /categories/{id} [get]
func GetCategoryByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID tidak valid",
		})
		return
	}

	category, err := repository.GetCategory(database.DbConnection, id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Kategori tidak ditemukan",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Gagal mendapatkan kategori",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": category,
	})
}

// CreateCategory godoc
// @Summary Create a category
// @Description Menambahkan kategori baru
// @Tags Categories
// @Accept json
// @Produce json
// @Param category body structs.CategoryInput true "Category object"
// @Success 201 {object} structs.APIResponse
// @Failure 400 {object} structs.APIResponse
// @Failure 500 {object} structs.APIResponse
// @Security BearerAuth
// @Router /categories [post]
func CreateCategory(c *gin.Context) {
	var input structs.CategoryInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Input tidak valid",
		})
		return
	}

	category := structs.Category{
		Name:      input.Name,
		CreatedBy: input.CreatedBy,
		ModifiedBy: input.ModifiedBy,
	}

	err := repository.CreateCategory(database.DbConnection, category)
	if err != nil {
		if err.Error() == "nama kategori sudah digunakan" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal menambahkan kategori",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Kategori berhasil ditambahkan",
	})
}

// UpdateCategory godoc
// @Summary Update a category
// @Description Memperbarui kategori
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param category body structs.UpdateCategoryInput true "Category object"
// @Success 200 {object} structs.APIResponse
// @Failure 400 {object} structs.APIResponse
// @Failure 404 {object} structs.APIResponse
// @Failure 500 {object} structs.APIResponse
// @Security BearerAuth
// @Router /categories/{id} [put]
func UpdateCategory(c *gin.Context) {
	var input structs.UpdateCategoryInput

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID tidak valid",
		})
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Input tidak valid",
		})
		return
	}

	category := structs.Category{
		ID:         id,
		Name:       input.Name,
		ModifiedBy: input.ModifiedBy,
	}

	err = repository.UpdateCategory(database.DbConnection, category)
	if err != nil {
		if err.Error() == "kategori tidak ditemukan" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		
		if err.Error() == "nama kategori sudah digunakan" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal memperbarui kategori",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Kategori berhasil diperbarui",
	})
}

// DeleteCategory godoc
// @Summary Delete a category
// @Description Menghapus kategori
// @Tags Categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} structs.APIResponse
// @Failure 400 {object} structs.APIResponse
// @Failure 404 {object} structs.APIResponse
// @Failure 500 {object} structs.APIResponse
// @Security BearerAuth
// @Router /categories/{id} [delete]
func DeleteCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID tidak valid",
		})
		return
	}

	err = repository.DeleteCategory(database.DbConnection, id)
	if err != nil {
		if err.Error() == "kategori tidak dapat dihapus karena masih memiliki buku" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		if err.Error() == "kategori tidak ditemukan" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal menghapus kategori",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Kategori berhasil dihapus",
	})
}

// GetCategoryBooks godoc
// @Summary Get books by category ID
// @Description Mendapatkan buku berdasarkan ID kategori
// @Tags Categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} structs.APIResponse
// @Failure 400 {object} structs.APIResponse
// @Failure 500 {object} structs.APIResponse
// @Security BearerAuth
// @Router /categories/{id}/books [get]
func GetCategoryBooks(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID tidak valid",
		})
		return
	}

	books, err := repository.GetCategoryBooks(database.DbConnection, id)
	if err != nil {
		if err.Error() == "kategori tidak ditemukan" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mendapatkan buku untuk kategori ini",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": books,
	})
}