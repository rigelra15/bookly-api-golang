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

func CreateCategory(c *gin.Context) {
	var category structs.Category

	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Input tidak valid",
		})
		return
	}

	err := repository.CreateCategory(database.DbConnection, category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal menambahkan kategori",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Kategori berhasil ditambahkan",
	})
}

func UpdateCategory(c *gin.Context) {
	var category structs.Category
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID tidak valid",
		})
		return
	}

	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Input tidak valid",
		})
		return
	}

	category.ID = id

	err = repository.UpdateCategory(database.DbConnection, category)
	if err != nil {
		if err.Error() == "kategori tidak ditemukan" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Gagal memperbarui kategori",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Kategori berhasil diperbarui",
	})
}

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
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Kategori tidak ditemukan",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Gagal menghapus kategori",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Kategori berhasil dihapus",
	})
}

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
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mendapatkan buku untuk kategori ini",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": books,
	})
}