package repository

import (
	"bookly-api-golang/structs"
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

func GetAllUsers(db *sql.DB) ([]structs.User, error) {
	var users []structs.User
	sqlQuery := "SELECT id, username, created_at, created_by, modified_at, modified_by FROM users"
	rows, err := db.Query(sqlQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user structs.User
		err := rows.Scan(&user.ID, &user.Username, &user.CreatedAt, &user.CreatedBy, &user.ModifiedAt, &user.ModifiedBy)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func GetUserByID(db *sql.DB, id int) (structs.User, error) {
	var user structs.User
	sqlQuery := "SELECT id, username, created_at, created_by, modified_at, modified_by FROM users WHERE id = $1"
	err := db.QueryRow(sqlQuery, id).Scan(&user.ID, &user.Username, &user.CreatedAt, &user.CreatedBy, &user.ModifiedAt, &user.ModifiedBy)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, errors.New("user tidak ditemukan")
		}
		return user, err
	}
	return user, nil
}

func CreateUser(db *sql.DB, user structs.User) error {
	usernameExists, err := IsUsernameExists(db, user.Username)
	if err != nil {
		return err
	}
	if usernameExists {
		return errors.New("username sudah digunakan")
	}

	sqlQuery := "INSERT INTO users (username, password, created_by, modified_by) VALUES ($1, $2, $3, $4)"
	_, err = db.Exec(sqlQuery, user.Username, user.Password, user.CreatedBy, user.ModifiedBy)
	if err != nil {
		return err
	}

	return nil
}

func IsUsernameExists(db *sql.DB, username string) (bool, error) {
	var exists bool
	sqlQuery := "SELECT EXISTS (SELECT 1 FROM users WHERE username = $1)"
	err := db.QueryRow(sqlQuery, username).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func UpdateUser(db *sql.DB, user structs.User) error {
	usernameExists, err := IsUsernameExistsByOtherUser(db, user.Username, user.ID)
	if err != nil {
		return err
	}
	if usernameExists {
		return errors.New("username sudah digunakan")
	}

	sqlQuery := "UPDATE users SET username=$1, password=$2, modified_by=$3, modified_at=CURRENT_TIMESTAMP WHERE id=$4"
	result, err := db.Exec(sqlQuery, user.Username, user.Password, user.ModifiedBy, user.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("user tidak ditemukan")
	}

	return nil
}

func IsUsernameExistsByOtherUser(db *sql.DB, username string, userID int) (bool, error) {
	var exists bool
	sqlQuery := "SELECT EXISTS (SELECT 1 FROM users WHERE username = $1 AND id != $2)"
	err := db.QueryRow(sqlQuery, username, userID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}


func DeleteUser(db *sql.DB, id int) error {
	sqlQuery := "DELETE FROM users WHERE id=$1"
	result, err := db.Exec(sqlQuery, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("user tidak ditemukan")
	}
	return nil
}

func GetUserByUsername(db *sql.DB, username string) (structs.User, error) {
	var user structs.User
	sqlQuery := "SELECT id, username, password, created_at, created_by, modified_at, modified_by FROM users WHERE username = $1"
	err := db.QueryRow(sqlQuery, username).Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt, &user.CreatedBy, &user.ModifiedAt, &user.ModifiedBy)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, errors.New("user tidak ditemukan")
		}
		return user, err
	}
	return user, nil
}

func CheckPasswordHash(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
