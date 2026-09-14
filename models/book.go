package models

import "time"

type BookModel struct {
	ID            int
	Title         string
	Review        string
	PublishedYear string
	Author        string
	Pages         int
	Completed     bool
	CreatedAt     time.Time
	CompletedAt   *time.Time
}
