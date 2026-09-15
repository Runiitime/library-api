package helpers

import (
	"library-api/models"
	"log"

	"github.com/jackc/pgx/v5"
)

func ConvertRows(rows pgx.Rows) ([]models.Book, error) {
	books := make([]models.Book, 0)

	if err := rows.Err(); err != nil {
		log.Printf("rows error: %v", err)
		return books, err
	}
	
	for rows.Next() {
		var book models.Book

		if err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.Review,
			&book.PublishedYear,
			&book.Author,
			&book.Pages,
			&book.Completed,
			&book.CreatedAt,
			&book.CompletedAt,
		); err != nil {
			log.Printf("scan error: %v", err)
			return books, err
		}

		books = append(books, book)
	}

	return books, nil
}
