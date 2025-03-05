package persistance

import (
	"github.com/zacharyworks/sync/packages/dir"
	"testing"
	"time"
)

func TestStore(t *testing.T) {
	type action struct {
		in          []dir.File
		expectedOut map[string]dir.Modification
	}

	initialTime := time.Now()

	tests := []struct {
		name    string
		actions []action
	}{
		{
			name: "Create {A,B}, Delete 'B'",
			actions: []action{
				{ // CREATE initial files
					in: []dir.File{
						{Name: "A", LastUpdated: initialTime.Add(-60 * time.Second)},
						{Name: "B", LastUpdated: initialTime.Add(-30 * time.Second)},
					},
					expectedOut: map[string]dir.Modification{
						"A": dir.Create,
						"B": dir.Create,
					},
				},
				{ // DELETE ONE FILE
					in: []dir.File{
						{Name: "A", LastUpdated: initialTime.Add(-60 * time.Second)},
					},
					expectedOut: map[string]dir.Modification{
						"B": dir.Delete,
					},
				},
			},
		},
		{
			name: "Create {A,B}, Update 'A' & Delete 'B'",
			actions: []action{
				{ // CREATE initial files
					in: []dir.File{
						{Name: "A", LastUpdated: initialTime.Add(-60 * time.Second)},
						{Name: "B", LastUpdated: initialTime.Add(-30 * time.Second)},
					},
					expectedOut: map[string]dir.Modification{
						"A": dir.Create,
						"B": dir.Create,
					},
				},
				{ // DELETE ONE FILE
					in: []dir.File{
						{Name: "A", LastUpdated: initialTime.Add(-30 * time.Second)},
					},
					expectedOut: map[string]dir.Modification{
						"A": dir.Update,
						"B": dir.Delete,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := New()
			// for each action
			for _, a := range tt.actions {
				// we supply some skeleton file info
				changes := store.Update(a.in)
				// AND check changes are expected
				for _, c := range changes {
					if a.expectedOut[c.F.Name] != c.Modification {
						t.Logf("for file %s expected modification %s but got %s", c.F.Name, a.expectedOut[c.F.Name], c.Modification)
						t.Fail()
					}
				}
			}
		})
	}
}
