package library

import "time"

type Book struct {
	Title         string     `json:"title"`
	Author        string     `json:"author"`
	Review        string     `json:"review"`
	PublishedYear string     `json:"published_year"`
	Pages         int        `json:"pages"`
	Completed     bool       `json:"completed"`
	CreatedAt     time.Time  `json:"created_at"`
	CompletedAt   *time.Time `json:"completed_at"`
}

func NewBook(title, author string, pages int) *Book {
	return &Book{
		Title:     title,
		Author:    author,
		Pages:     pages,
		Completed: false,
		CreatedAt: time.Now(),
	}
}

func (b *Book) ChangeCompletedStatus(status bool) {
	b.Completed = status
	b.CompletedAt = new(time.Now())
}
