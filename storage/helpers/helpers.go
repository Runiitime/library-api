package helpers

import (
	"library-api/models"
	bookErr "library-api/storage/msg"
	"log"

	"github.com/jackc/pgx/v5"
)

func ConvertRows(rows pgx.Rows) ([]models.Book, error) {
	books := make([]models.Book, 0)

	if err := rows.Err(); err != nil {
		log.Fatalln(bookErr.ErrRows, err)
		return books, bookErr.ErrRows
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
			log.Fatalln(bookErr.ErrRowScan, err)
			return books, bookErr.ErrRowScan
		}

		books = append(books, book)
	}

	return books, nil
}
