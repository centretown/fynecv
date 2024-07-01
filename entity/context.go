package entity

type Context struct {
	ID       string `json:"id" yaml:"id"`
	ParentID string `json:"parent_id" yaml:"parent_id"`
	UserID   string `json:"user_id" yaml:"user_id"`
}
