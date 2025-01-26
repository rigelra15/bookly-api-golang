package repository

import (
	"bookly-api-golang/structs"
	"database/sql"
	"errors"
	"log"
	"fmt"
	"strings"
)

func ValidateCategoryID(db *sql.DB, categoryID int) (bool, error) {
	sql := "SELECT EXISTS (SELECT 1 FROM categories WHERE id=$1)"
	var exists bool
	err := db.QueryRow(sql, categoryID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func ValidateReleaseYear(year int) error {
	if year < 1980 || year > 2024 {
		return fmt.Errorf("release_year harus di antara 1980 dan 2024")
	}
	return nil
}

func ValidateExistBook(db *sql.DB, id int) (bool, error) {
	sql := "SELECT EXISTS (SELECT 1 FROM books WHERE id=$1)"
	var exists bool
	err := db.QueryRow(sql, id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func CalculateThickness(totalPage int) string {
	if totalPage > 100 {
		return "tebal"
	}
	return "tipis"
}

func GetAllBook(db *sql.DB) (result []structs.Book, err error) {
	sql := "SELECT * FROM books"
	rows, err := db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var book structs.Book

		err := rows.Scan(&book.ID, &book.Title, &book.Description, &book.ImageURL, &book.ReleaseYear, &book.Price, &book.TotalPage, &book.Thickness, &book.CategoryID, &book.CreatedAt, &book.CreatedBy, &book.ModifiedAt, &book.ModifiedBy)
		if err != nil {
			return nil, err
		}

		result = append(result, book)
	}

	return result, nil
}

func GetBook(db *sql.DB, id int) (result structs.Book, err error) {
	sql := "SELECT * FROM books WHERE id = $1"
	err = db.QueryRow(sql, id).Scan(&result.ID, &result.Title, &result.Description, &result.ImageURL, &result.ReleaseYear, &result.Price, &result.TotalPage, &result.Thickness, &result.CategoryID, &result.CreatedAt, &result.CreatedBy, &result.ModifiedAt, &result.ModifiedBy)
	if err != nil {
		return result, err
	}

	return result, nil
}

func CreateBook(db *sql.DB, book structs.Book) (err error) {
	categoryExists, err := ValidateCategoryID(db, book.CategoryID)
	if err != nil {
		return err
	}

	if !categoryExists {
		log.Printf("Invalid category_id: %d", book.CategoryID)
		return errors.New("category_id tidak valid")
	}

	existingBook, err := GetBookByTitle(db, book.Title)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if strings.EqualFold(existingBook.Title, book.Title) {
		return errors.New("judul buku sudah digunakan oleh buku lain")
	}	

	if err := ValidateReleaseYear(book.ReleaseYear); err != nil {
		return err
	}

	book.Thickness = CalculateThickness(book.TotalPage)

	sql := `
        INSERT INTO books (title, description, image_url, release_year, price, total_page, thickness, category_id, created_by, modified_by)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        RETURNING id, created_at, modified_at
    `
	err = db.QueryRow(sql, book.Title, book.Description, book.ImageURL, book.ReleaseYear, book.Price, book.TotalPage, book.Thickness, book.CategoryID, book.CreatedBy, book.ModifiedBy).Scan(
		&book.ID, &book.CreatedAt, &book.ModifiedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func GetBookByTitle(db *sql.DB, title string) (structs.Book, error) {
	var book structs.Book
	sql := "SELECT * FROM books WHERE LOWER(title) = LOWER($1)"
	err := db.QueryRow(sql, title).Scan(
		&book.ID, &book.Title, &book.Description, &book.ImageURL,
		&book.ReleaseYear, &book.Price, &book.TotalPage, &book.Thickness,
		&book.CategoryID, &book.CreatedAt, &book.CreatedBy, &book.ModifiedAt, &book.ModifiedBy,
	)
	return book, err
}

func UpdateBook(db *sql.DB, book structs.Book) (err error) {
	exists, err := ValidateExistBook(db, book.ID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("buku tidak ditemukan")
	}

	categoryExists, err := ValidateCategoryID(db, book.CategoryID)
	if err != nil {
		return err
	}
	if !categoryExists {
		return errors.New("category_id tidak valid")
	}

	if err := ValidateReleaseYear(book.ReleaseYear); err != nil {
		return err
	}

	titleUsed, err := IsTitleUsedByOtherBook(db, book.Title, book.ID)
	if err != nil {
		return err
	}
	if titleUsed {
		return errors.New("judul buku sudah digunakan oleh buku lain")
	}

	book.Thickness = CalculateThickness(book.TotalPage)

	sql := `
        UPDATE books 
        SET title=$1, description=$2, image_url=$3, release_year=$4, price=$5, 
            total_page=$6, thickness=$7, category_id=$8, modified_at=CURRENT_TIMESTAMP, modified_by=$9 
        WHERE id=$10
    `
	result, err := db.Exec(sql, book.Title, book.Description, book.ImageURL, book.ReleaseYear, book.Price, book.TotalPage, book.Thickness, book.CategoryID, book.ModifiedBy, book.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("buku tidak ditemukan")
	}

	return nil
}

func IsTitleUsedByOtherBook(db *sql.DB, title string, bookID int) (bool, error) {
	var exists bool
	sql := `
        SELECT EXISTS (
            SELECT 1 FROM books WHERE LOWER(title) = LOWER($1) AND id != $2
        )
    `
	err := db.QueryRow(sql, title, bookID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func DeleteBook(db *sql.DB, id int) (err error) {
	var exists bool
	sqlCheck := "SELECT EXISTS (SELECT 1 FROM books WHERE id=$1)"
	err = db.QueryRow(sqlCheck, id).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("buku tidak ditemukan")
	}

	sql := "DELETE FROM books WHERE id=$1"
	_, err = db.Exec(sql, id)
	if err != nil {
		return err
	}

	return nil
}