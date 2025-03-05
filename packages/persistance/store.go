package persistance

import "github.com/zacharyworks/sync/packages/dir"

// Persistence is a simple implementation of FileState which records, in memory
// the state of a directory
// todo: implications of storing []byte of file, if the DIR contained large files
type Persistence map[string]dir.File

func New() dir.FileState {
	return make(Persistence)
}

// Update accepts a list of files in a directory, and detects ongoing changes by comparing
// the latest state of a DIR vs the known state
func (p Persistence) Update(latest []dir.File) []dir.FileChange {
	var changes []dir.FileChange

	latestPoll := make(map[string]dir.File)
	for _, file := range latest {
		latestPoll[file.Name] = file
	}

	// if there is something in persistence which is not in latest read,
	// we can infer it has been DELETE or moved
	for fileName, file := range p {
		if _, exists := latestPoll[fileName]; !exists {
			changes = append(changes, dir.FileChange{F: file, Modification: dir.Delete})
			delete(p, fileName)
		}
	}

	for _, latestFile := range latest {
		// if current file doesn't exist, then it is CREATED
		existingFile, exists := p[latestFile.Name]
		if !exists {
			changes = append(changes, dir.FileChange{F: latestFile, Modification: dir.Create})
			p[latestFile.Name] = latestFile
			continue
		}
		// otherwise we should check if the current file is older than the new file, then UPDATED
		if existingFile.LastUpdated.Before(latestPoll[latestFile.Name].LastUpdated) {
			changes = append(changes, dir.FileChange{F: latestFile, Modification: dir.Update})
			p[latestFile.Name] = latestFile
		}
	}

	return changes
}
