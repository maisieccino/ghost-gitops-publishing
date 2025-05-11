// internal/api/post.go

package api

type postReq struct {
	Posts []Post `json:"posts"`
}

type Post struct {
	ID            string   `json:"id,omitempty"`
	Title         string   `json:"title"`
	Slug          string   `json:"slug,omitempty"`
	Status        string   `json:"status,omitempty"`
	HTML          string   `json:"html"`
	FeatureImage  string   `json:"feature_image,omitempty"`
	Tags          []tagRef `json:"tags,omitempty"`
	CustomExcerpt string   `json:"custom_excerpt,omitempty"`
	PublishedAt   string   `json:"published_at,omitempty"`
	UpdatedAt     string   `json:"updated_at,omitempty"`
}

type tagRef struct {
	Name string `json:"name"`
	Slug string `json:"slug,omitempty"`
}

func WrapTags(tags []string) []tagRef {
	out := make([]tagRef, len(tags))
	for i, t := range tags {
		out[i] = tagRef{Name: t}
	}
	return out
}
