package dir

import "time"

type Modification string

const (
	Create = Modification("create")
	Update = Modification("update")
	Delete = Modification("delete")
)

// FileChange represents an update to the state of the store from the previous
// time the caller updated it
type FileChange struct {
	F File
	Modification
}

// File represents the recentness of a particular file name
type File struct {
	Name        string
	LastUpdated time.Time
}
