package note

import (
	"errors"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	now := time.Date(2026, 10, 15, 9, 30, 0, 0, time.UTC)

	tests := []struct {
		name      string
		id        ID
		title     string
		wantErr   error
		wantTitle string
	}{
		{"valid", "n1", "Write-Ahead Log", nil, "Write-Ahead Log"},
		{"trims whitespace", "n1", "  Spaced  ", nil, "Spaced"},
		{"empty title", "n1", "", ErrEmptyTitle, ""},
		{"blank title", "n1", "   ", ErrEmptyTitle, ""},
		{"empty id", "", "Title", ErrEmptyID, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.id, tt.title, now)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("New() error = %v, want %v", err, tt.wantErr)
			}
			if err != nil {
				return
			}
			if got.Title != tt.wantTitle {
				t.Errorf("Title = %q, want %q", got.Title, tt.wantTitle)
			}
			if !got.CreatedAt.Equal(now) || !got.UpdatedAt.Equal(now) {
				t.Errorf("timestamps = %v / %v, want both %v", got.CreatedAt, got.UpdatedAt, now)
			}
		})
	}
}

func TestSetBody(t *testing.T) {
	created := time.Date(2026, 10, 15, 9, 30, 0, 0, time.UTC)
	later := created.Add(time.Hour)

	n, err := New("n1", "Title", created)
	if err != nil {
		t.Fatal(err)
	}

	n.SetBody("hello", later)

	if n.Body != "hello" {
		t.Errorf("Body = %q, want %q", n.Body, "hello")
	}
	if !n.UpdatedAt.Equal(later) {
		t.Errorf("UpdatedAt = %v, want %v", n.UpdatedAt, later)
	}
	if !n.CreatedAt.Equal(created) {
		t.Errorf("CreatedAt changed to %v, want %v", n.CreatedAt, created)
	}
}
