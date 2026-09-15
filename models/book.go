package models

import "time"

type Book struct {
	ID            int        `json:"id"`
	Title         string     `json:"title"`
	Review        string     `json:"review"`
	Author        string     `json:"author"`
	Pages         int        `json:"pages"`
	PublishedYear string     `json:"published_year"`
	Completed     bool       `json:"completed"`
	CreatedAt     time.Time  `json:"created_at"`
	CompletedAt   *time.Time `json:"completed_at"`
}
