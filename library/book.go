package library

import "time"

type Book struct {
	Title        string     `json:"title"`
	Author       string     `json:"author"`
	Pages        int        `json:"pages"`
	Completed    bool       `json:"completed"`
	Created_At   time.Time  `json:"created_at"`
	Completed_At *time.Time `json:"completed_at"`
}

func NewBook(title, author string, pages int) *Book {
	return &Book{
		Title:      title,
		Author:     author,
		Pages:      pages,
		Completed:  false,
		Created_At: time.Now(),
	}
}

func (b *Book) ChangeCompletedStatus(status bool) {
	completedAt := time.Now()
	b.Completed = status
	b.Completed_At = &completedAt
}
