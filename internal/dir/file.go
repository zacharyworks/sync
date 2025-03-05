package dir

import "time"

type Modification string

type FileState interface {
	Update(file []File) []FileChange
}

type FileUpdate interface {
}

const (
	Create = Modification("create")
	Update = Modification("update")
	Delete = Modification("delete")
)

// FileChange represents an update to the state of the store from the previous
// time the caller updated it
type FileChange struct {
	F            File `json:"file"`
	Modification `json:"modification"`
}

// File represents the recentness of a particular file name
type File struct {
	Name        string    `json:"name"`
	Path        string    `json:"path"`
	Content     []byte    `json:"content"`
	LastUpdated time.Time `json:"last_updated"`
}
