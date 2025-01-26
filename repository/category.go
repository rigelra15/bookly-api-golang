package repository

import (
	"bookly-api-golang/structs"
	"database/sql"
	"errors"
	"strings"
)

func ValidateCategory(db *sql.DB, id int) (bool, error) {
	sql := "SELECT EXISTS (SELECT 1 FROM categories WHERE id=$1)"
	var exists bool
	err := db.QueryRow(sql, id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func GetAllCategory(db *sql.DB) (result []structs.Category, err error) {
	sql := "SELECT * FROM categories"
	rows, err := db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var category structs.Category

		err := rows.Scan(&category.ID, &category.Name, &category.CreatedAt, &category.CreatedBy, &category.ModifiedAt, &category.ModifiedBy)
		if err != nil {
			return nil, err
		}

		result = append(result, category)
	}

	return result, nil
}

func GetCategory(db *sql.DB, id int) (result structs.Category, err error) {
	sql := "SELECT * FROM categories WHERE id = $1"
	err = db.QueryRow(sql, id).Scan(&result.ID, &result.Name, &result.CreatedAt, &result.CreatedBy, &result.ModifiedAt, &result.ModifiedBy)
	if err != nil {
		return result, err
	}

	return result, nil
}

func CreateCategory(db *sql.DB, category structs.Category) (err error) {
	existingCategory, err := GetCategoryByName(db, category.Name)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if strings.EqualFold(existingCategory.Name, category.Name) {
		return errors.New("nama kategori sudah digunakan")
	}

	sql := `
		INSERT INTO categories (name, created_by, modified_by)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, modified_at
	`
	err = db.QueryRow(sql, category.Name, category.CreatedBy, category.ModifiedBy).Scan(
		&category.ID, &category.CreatedAt, &category.ModifiedAt,
	)
	if err != nil {
		return err
	}

	return nil
}


func GetCategoryByName(db *sql.DB, name string) (structs.Category, error) {
	var category structs.Category
	sql := `SELECT * FROM categories WHERE LOWER(name) = LOWER($1)`
	err := db.QueryRow(sql, name).Scan(&category.ID, &category.Name, &category.CreatedAt, &category.CreatedBy, &category.ModifiedAt, &category.ModifiedBy)
	return category, err
}

func UpdateCategory(db *sql.DB, category structs.Category) (err error) {
	exists, err := ValidateCategory(db, category.ID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("kategori tidak ditemukan")
	}

	existingCategory, err := GetCategoryByNameExcludingID(db, category.Name, category.ID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if strings.EqualFold(existingCategory.Name, category.Name) {
		return errors.New("nama kategori sudah digunakan")
	}

	sql := `
		UPDATE categories
		SET name = $1, modified_at = CURRENT_TIMESTAMP, modified_by = $2
		WHERE id = $3
	`
	_, err = db.Exec(sql, category.Name, category.ModifiedBy, category.ID)
	if err != nil {
		return err
	}

	return nil
}

func GetCategoryByNameExcludingID(db *sql.DB, name string, id int) (structs.Category, error) {
	var category structs.Category
	sql := `SELECT * FROM categories WHERE LOWER(name) = LOWER($1) AND id != $2`
	err := db.QueryRow(sql, name, id).Scan(&category.ID, &category.Name, &category.CreatedAt, &category.CreatedBy, &category.ModifiedAt, &category.ModifiedBy)
	return category, err
}

func DeleteCategory(db *sql.DB, id int) (err error) {
	hasBooks, err := CategoryHasBooks(db, id)
	if err != nil {
		return err
	}
	if hasBooks {
		return errors.New("kategori tidak dapat dihapus karena masih memiliki buku")
	}

	exists, err := ValidateCategory(db, id)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("kategori tidak ditemukan")
	}

	sql := "DELETE FROM categories WHERE id=$1"
	_, err = db.Exec(sql, id)
	if err != nil {
		return err
	}

	return nil
}

func CategoryHasBooks(db *sql.DB, categoryID int) (bool, error) {
	var exists bool
	sqlQuery := "SELECT EXISTS (SELECT 1 FROM books WHERE category_id = $1)"
	err := db.QueryRow(sqlQuery, categoryID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func GetCategoryBooks(db *sql.DB, id int) (result []structs.Book, err error) {
	exists, err := ValidateCategory(db, id)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("kategori tidak ditemukan")
	}

	sql := "SELECT * FROM books WHERE category_id=$1"
	rows, err := db.Query(sql, id)
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
