// internal/frontmatter/parser.go

package frontmatter

import (
	"bytes"
	"os"

	fm "github.com/adrg/frontmatter"
	"gopkg.in/yaml.v3"
)

// Meta holds every key ghostpost cares about.
// Add more tags as your workflow grows.
type Meta struct {
	Title          string   `yaml:"title"`
	Slug           string   `yaml:"slug,omitempty"`
	Status         string   `yaml:"status,omitempty"` // draft | published | scheduled
	PublishedAt    string   `yaml:"published_at,omitempty"`
	Visibility     string   `yaml:"visibility,omitempty"` // public | members | paid | specific
	Tiers          []string `yaml:"tiers,omitempty"`
	Featured       bool     `yaml:"featured,omitempty"`
	CustomExcerpt  string   `yaml:"custom_excerpt,omitempty"`
	Authors        []string `yaml:"authors,omitempty"`
	CustomTemplate string   `yaml:"custom_template,omitempty"`
	FeatureImage   string   `yaml:"feature_image,omitempty"`
	Tags           []string `yaml:"tags,omitempty"`
	PostID         string   `yaml:"post_id,omitempty"` // set after first publish
	Hash           string   `yaml:"hash,omitempty"`    // SHA256 of Markdown body
}

// ParseFile reads a Markdown file and returns its meta + body bytes.
func ParseFile(path string) (Meta, []byte, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return Meta{}, nil, err
	}
	var meta Meta
	body, err := fm.Parse(bytes.NewReader(raw), &meta)
	return meta, body, err
}

// WriteFile rewrites the Markdown file with updated front-matter.
func WriteFile(path string, meta Meta, body []byte) error {
	// marshal the front-matter
	fmBytes, _ := yaml.Marshal(meta)

	// strip out any leading "\n" so we only get one blank line
	body = bytes.TrimLeft(body, "\n")

	// build the new file
	var buf bytes.Buffer
	buf.WriteString("---\n")
	buf.Write(fmBytes)
	buf.WriteString("---\n\n") // one blank line before the content
	buf.Write(body)

	return os.WriteFile(path, buf.Bytes(), 0o644)
}
