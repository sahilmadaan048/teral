// Package notefile converts notes to and from their on-disk format:
// a Markdown file with a YAML front matter header.
package notefile

import (
	"bytes"
	"fmt"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/sahilmadaan048/teral/internal/core/note"
)

// fence marks the start and end of the front matter block.
const fence = "---"

// frontMatter is the YAML header of a note file. It is a separate type
// from note.Note so that the file format can evolve without changing the
// domain type, and so the domain stays free of YAML details.
type frontMatter struct {
	ID      string    `yaml:"id"`
	Title   string    `yaml:"title"`
	Created time.Time `yaml:"created"`
	Updated time.Time `yaml:"updated"`
	Tags    []string  `yaml:"tags,omitempty,flow"`
}

// Encode renders a note as front matter followed by its body.
func Encode(n note.Note) ([]byte, error) {
	header, err := yaml.Marshal(frontMatter{
		ID:      string(n.ID),
		Title:   n.Title,
		Created: n.CreatedAt,
		Updated: n.UpdatedAt,
		Tags:    n.Tags,
	})
	if err != nil {
		return nil, fmt.Errorf("encode front matter: %w", err)
	}

	var buf bytes.Buffer
	buf.WriteString(fence + "\n")
	buf.Write(header)
	buf.WriteString(fence + "\n")
	buf.WriteString(n.Body)

	return buf.Bytes(), nil
}
