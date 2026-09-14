package storage

import (
	"context"
	"library-api/models"

	"github.com/jackc/pgx/v5"
)

type BookQueries struct {
	tableName string
}

func NewBookQueries(tableName string) *BookQueries {
	return &BookQueries{
		tableName: tableName,
	}
}

func (b *BookQueries) CreateBook(conn *pgx.Conn, ctx context.Context, book models.BookModel) error {
	q := `
	INSERT INTO ` + b.tableName + `
	(title, review, author, published, pages, completed, created_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := conn.Exec(
		ctx,
		q,
		book.Title,
		book.Review,
		book.Author,
		book.PublishedYear,
		book.Pages,
		book.Completed,
		book.CreatedAt)

	return err
}

func (b *BookQueries) DeleteBook(conn *pgx.Conn, ctx context.Context, booksIDs int) error {
	q := `
	DELETE FROM ` + b.tableName + `
	WHERE id ANY($1)
	`
	_, err := conn.Exec(ctx, q, booksIDs)
	return err
}

func (b *BookQueries) SelectBook(conn *pgx.Conn, ctx context.Context, condition string) ([]models.BookModel, error) {
	q := `
	SELECT * FROM ` + b.tableName + `
	WHERE ` + condition + `
	ORDER BY id ASC`

	rows, err := conn.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	books := make([]models.BookModel, 0)

	for rows.Next() {
		var book models.BookModel

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
			return nil, err
		}

		books = append(books, book)
	}

	return books, nil
}

func (b *BookQueries) UpdateBook(conn *pgx.Conn, ctx context.Context, book models.BookModel) error {
	q := `
	UPDATE ` + b.tableName + `
	SET title=$1, review=$2, author=$3, pages=$4, published=$5, completed=$6, created_at=$7, completed_at=$8
	WHERE id = $9`

	_, err := conn.Exec(
		ctx,
		q,
		book.Title,
		book.Review,
		book.Author,
		book.Pages,
		book.PublishedYear,
		book.Completed,
		book.CreatedAt,
		book.CompletedAt,
	)

	return err
}
