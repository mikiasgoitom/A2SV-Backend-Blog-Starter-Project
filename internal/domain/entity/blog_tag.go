package entity

// BlogTag represents the many-to-many relationship between blogs and tags
type BlogTag struct {
	BlogID string `json:"blog_id" db:"blog_id"`
	TagID  string `json:"tag_id" db:"tag_id"`
}
