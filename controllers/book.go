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

func CreateBook(c *gin.Context) {
	var book structs.Book

	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Data yang dimasukkan tidak valid",
		})
		return
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

func UpdateBook(c *gin.Context) {
	var book structs.Book
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID tidak valid",
		})
		return
	}

	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Data yang dimasukkan tidak valid",
		})
		return
	}

	book.ID = id

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