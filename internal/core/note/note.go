package note

import (
	"errors"
	"strings"
	"time"
)

// ID uniquely identifies a note for its whole life. It does not change
// when the note is renamed or moved, which is what lets links survive.
type ID string

// Validation errors returned by New. Callers compare against these with
// errors.Is instead of matching on error text.
var (
	ErrEmptyID    = errors.New("note: id must not be empty")
	ErrEmptyTitle = errors.New("note: title must not be empty")
)

// Note is the core unit of knowledge in Teral.
type Note struct {
	ID        ID
	Title     string
	Body      string
	Tags      []string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// New builds a valid Note. The caller supplies the id and the current
// time so that this function stays deterministic and easy to test.
func New(id ID, title string, now time.Time) (Note, error) {
	title = strings.TrimSpace(title)

	if id == "" {
		return Note{}, ErrEmptyID
	}
	if title == "" {
		return Note{}, ErrEmptyTitle
	}

	now = now.UTC()
	return Note{
		ID:        id,
		Title:     title,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// SetBody replaces the note's content and records when that happened.
func (n *Note) SetBody(body string, now time.Time) {
	n.Body = body
	n.UpdatedAt = now.UTC()
}
