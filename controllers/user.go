package controllers

import (
	"bookly-api-golang/database"
	"bookly-api-golang/middlewares"
	"bookly-api-golang/repository"
	"bookly-api-golang/structs"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Login godoc
// @Summary Login user
// @Description Autentikasi user untuk mendapatkan token JWT
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body structs.Login true "User login credentials"
// @Success 200 {object} structs.APIResponse
// @Failure 400 {object} structs.APIResponse
// @Failure 401 {object} structs.APIResponse
// @Failure 500 {object} structs.APIResponse
// @Router /users/login [post]
func Login(c *gin.Context) {
	var loginRequest structs.Login

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Input tidak valid",
		})
		return
	}

	user, err := repository.GetUserByUsername(database.DbConnection, loginRequest.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Username atau password salah",
		})
		return
	}

	if !repository.CheckPasswordHash(loginRequest.Password, user.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Username atau password salah",
			})
			return
	}


	token, err := middlewares.GenerateJWT(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal membuat token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}


// GetAllUsers godoc
// @Summary Get all users
// @Description Mendapatkan semua user
// @Tags Users
// @Produce json
// @Success 200 {object} []structs.User
// @Failure 500 {object} structs.APIResponse
// @Security JWT
// @Router /users [get]
func GetAllUsers(c *gin.Context) {
	users, err := repository.GetAllUsers(database.DbConnection)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mendapatkan data user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": users,
	})
}

// GetUserByID godoc
// @Summary Get user by ID
// @Description Mendapatkan user berdasarkan ID
// @Tags Users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} structs.User
// @Failure 404 {object} structs.APIResponse
// @Failure 500 {object} structs.APIResponse
// @Security JWT
// @Router /users/{id} [get]
func GetUserByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID tidak valid",
		})
		return
	}

	user, err := repository.GetUserByID(database.DbConnection, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": user,
	})
}

// CreateUser godoc
// @Summary Create new user
// @Description Membuat user baru
// @Tags Users
// @Accept json
// @Produce json
// @Param user body structs.UserInput true "User object"
// @Success 201 {object} structs.APIResponse
// @Failure 400 {object} structs.APIResponse
// @Failure 500 {object} structs.APIResponse
// @Security JWT
// @Router /users [post]
func CreateUser(c *gin.Context) {
	var input structs.UserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Input tidak valid",
		})
		return
	}

	hashedPassword, err := repository.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal meng-hash password",
		})
		return
	}

	user := structs.User{
		Username:   input.Username,
		Password:   hashedPassword,
		CreatedBy:  input.CreatedBy,
		ModifiedBy: input.ModifiedBy,
	}

	err = repository.CreateUser(database.DbConnection, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal membuat user",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User berhasil dibuat",
	})
}

// UpdateUser godoc
// @Summary Update user
// @Description Mengubah data user
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body structs.UserInput true "User object"
// @Success 200 {object} structs.APIResponse
// @Failure 400 {object} structs.APIResponse
// @Failure 404 {object} structs.APIResponse
// @Failure 500 {object} structs.APIResponse
// @Security JWT
// @Router /users/{id} [put]
func UpdateUser(c *gin.Context) {
	var input structs.UserInput
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

	hashedPassword, err := repository.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal meng-hash password",
		})
		return
	}

	user := structs.User{
		ID:         id,
		Username:   input.Username,
		Password:   hashedPassword,
		ModifiedBy: input.ModifiedBy,
	}

	err = repository.UpdateUser(database.DbConnection, user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User berhasil diperbarui",
	})
}

// DeleteUser godoc
// @Summary Delete user
// @Description Menghapus user berdasarkan ID
// @Tags Users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} structs.APIResponse
// @Failure 404 {object} structs.APIResponse
// @Failure 500 {object} structs.APIResponse
// @Security JWT
// @Router /users/{id} [delete]
func DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID tidak valid",
		})
		return
	}

	err = repository.DeleteUser(database.DbConnection, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User berhasil dihapus",
	})
}
