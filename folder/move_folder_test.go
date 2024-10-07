package folder_test

import (
	"errors"
	"testing"

	"github.com/georgechieng-sc/interns-2022/folder"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_folder_MoveFolder(t *testing.T) {
	org1 := uuid.Must(uuid.NewV4())
	org2 := uuid.Must(uuid.NewV4())

	t.Parallel()
	tests := [...]struct {
		name        string
		destination string
		folders     []folder.Folder
		want        []folder.Folder
	}{
		{
			"alpha",
			"beta",
			[]folder.Folder{
				{"alpha", org1, "alpha"},
				{"beta", org1, "beta"},
			},
			[]folder.Folder{
				{"alpha", org1, "beta.alpha"},
				{"beta", org1, "beta"},
			},
		},
		{
			"alpha",
			"beta",
			[]folder.Folder{
				{"alpha", org1, "alpha"},
				{"beta", org1, "beta"},
				{"delta", org1, "alpha.delta"},
			},
			[]folder.Folder{
				{"alpha", org1, "beta.alpha"},
				{"beta", org1, "beta"},
				{"delta", org1, "beta.alpha.delta"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := folder.NewDriver(tt.folders)

			endFolders, err := f.MoveFolder(tt.name, tt.destination)
			assert.Equal(t, checkSliceEqual(endFolders, tt.want), true, "got %v, want %v", endFolders, tt.want)
			assert.Nil(t, err)
		})
	}

	errorTests := [...]struct {
		name        string
		destination string
		folders     []folder.Folder
		err         error
	}{
		{
			"alpha",
			"beta",
			[]folder.Folder{
				{"alpha", org1, "alpha"},
				{"beta", org1, "alpha.beta"},
			},
			errors.New("Error: Cannot move a folder to a child of itself"),
		},
		{
			"alpha",
			"alpha",
			[]folder.Folder{
				{"alpha", org1, "alpha"},
				{"beta", org1, "alpha.beta"},
			},
			errors.New("Error: Cannot move a folder to itself"),
		},
		{
			"alpha",
			"beta",
			[]folder.Folder{
				{"alpha", org1, "alpha"},
				{"beta", org2, "beta"},
			},
			errors.New("Error: Cannot move a folder to a different organization"),
		},
		{
			"alpha",
			"beta",
			[]folder.Folder{
				{"alpha", org2, "alpha"},
				{"beta", org1, "beta"},
			},
			errors.New("Error: Cannot move a folder to a different organization"),
		},
		{
			"delta",
			"beta",
			[]folder.Folder{
				{"alpha", org1, "alpha"},
				{"beta", org1, "beta"},
			},
			errors.New("Error: Source folder does not exist"),
		},
		{
			"alpha",
			"delta",
			[]folder.Folder{
				{"alpha", org1, "alpha"},
				{"beta", org1, "beta"},
			},
			errors.New("Error: Destination folder does not exist"),
		},
	}
	for _, tt := range errorTests {
		t.Run(tt.name, func(t *testing.T) {
			f := folder.NewDriver(tt.folders)

			endFolders, err := f.MoveFolder(tt.name, tt.destination)

			assert.Error(t, err)
			assert.Equal(t, tt.err, err)
			assert.Nil(t, endFolders)
		})
	}
}
