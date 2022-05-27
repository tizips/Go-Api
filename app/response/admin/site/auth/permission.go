package auth

type TreePermission struct {
	Id        uint             `json:"id"`
	Parents   []uint           `json:"parents,omitempty"`
	Name      string           `json:"name,omitempty"`
	Slug      string           `json:"slug,omitempty"`
	Method    string           `json:"method,omitempty"`
	Path      string           `json:"path,omitempty"`
	CreatedAt string           `json:"created_at,omitempty"`
	Children  []TreePermission `json:"children,omitempty"`
}
