package storage

import (
	"context"
	"library-api/models"
	"library-api/storage/helpers"
	bookMSG "library-api/storage/msg"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
)

type BookQuery struct {
	tableName string
	conn      *pgx.Conn
}

func NewBookQuery(tableName string, conn *pgx.Conn) *BookQuery {
	return &BookQuery{
		tableName: tableName,
		conn:      conn,
	}
}

func (b *BookQuery) CreateBook(ctx context.Context, book models.Book) (int, error) {
	var id int
	q := `
	INSERT INTO ` + b.tableName + `
	(title, review, author, published, pages, completed, created_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;`

	err := b.conn.QueryRow(
		ctx,
		q,
		book.Title,
		book.Review,
		book.Author,
		book.PublishedYear,
		book.Pages,
		book.Completed,
		book.CreatedAt).Scan(&id)

	if err != nil {
		log.Println(bookMSG.ErrQueryRow, err)
		return -1, err
	}

	return id, nil
}

func (b *BookQuery) DeleteBook(ctx context.Context, bookID int) error {
	q := `
	DELETE FROM ` + b.tableName + `
	WHERE id=$1
	`
	tag, err := b.conn.Exec(ctx, q, bookID)
	if err != nil {
		log.Println(bookMSG.ErrExec, err)
	}

	if tag.RowsAffected() == 0 {
		return bookMSG.ErrBookNotFound
	}

	return nil
}

func (b *BookQuery) SelectBooksByID(ctx context.Context, ids []int) ([]models.Book, error) {
	q := `
	SELECT id, title, review, published, author, pages, completed, created_at, completed_at FROM ` + b.tableName + `
	WHERE id = ANY($1)
	ORDER BY id ASC`

	rows, err := b.conn.Query(ctx, q, ids)
	if err != nil {
		log.Println(bookMSG.ErrQuery, err)
		return nil, err
	}

	defer rows.Close()

	return helpers.ConvertRows(rows)
}

func (b *BookQuery) SelectAllBooks(ctx context.Context) ([]models.Book, error) {
	q := `
    SELECT id, title, review, published, author, pages, completed, created_at, completed_at
    FROM ` + b.tableName

	rows, err := b.conn.Query(ctx, q)
	if err != nil {
		log.Println(bookMSG.ErrQuery, err)
		return nil, err
	}

	defer rows.Close()
	return helpers.ConvertRows(rows)
}

func (b *BookQuery) SelectBooksByParams(ctx context.Context, field string, value string) ([]models.Book, error) {
	q := `
	SELECT id, title, review, published, author, pages, completed, created_at, completed_at FROM ` + b.tableName + `
 	WHERE ` + field + `= $1 ORDER BY id ASC`

	rows, err := b.conn.Query(ctx, q, value)
	if err != nil {
		log.Println(bookMSG.ErrQuery, err)
		return nil, err
	}

	defer rows.Close()
	return helpers.ConvertRows(rows)
}

func (b *BookQuery) UpdateBookStatus(ctx context.Context, id int, isCompleted bool) error {
	q := `
	UPDATE ` + b.tableName + `
	SET completed=$1, completed_at=$2
	WHERE id = $3`

	_, err := b.conn.Exec(
		ctx,
		q,
		isCompleted,
		time.Now(),
		id,
	)

	if err != nil {
		log.Println(bookMSG.ErrExec, err)
		return err
	}

	return nil
}
