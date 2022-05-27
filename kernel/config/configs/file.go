package configs

type File struct {
	Driver string `default:"system"`
	Path   string `default:"upload"`
}

const (
	FileDriverSystem = "system"
	FileDriverQiniu  = "qiniu"
)
