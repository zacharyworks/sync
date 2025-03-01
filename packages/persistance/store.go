package persistance

import "github.com/zacharyworks/sync/packages/dir"

// Persistence is a
type Persistence map[string]dir.File

func New() Persistence {
	return make(Persistence)
}

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
			changes = append(changes, dir.FileChange{F: existingFile, Modification: dir.Update})
			p[latestFile.Name] = latestFile
		}
	}

	return changes
}
