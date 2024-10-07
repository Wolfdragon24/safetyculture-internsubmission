package folder_test

import (
	"errors"
	"testing"

	"github.com/georgechieng-sc/interns-2022/folder"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

func checkSliceEqual(a []folder.Folder, b []folder.Folder) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func Test_folder_GetFoldersByOrgID(t *testing.T) {
	org1 := uuid.Must(uuid.NewV4())
	org2 := uuid.Must(uuid.NewV4())

	t.Parallel()
	tests := [...]struct {
		name    string
		orgID   uuid.UUID
		folders []folder.Folder
		want    []folder.Folder
	}{
		{
			"alpha",
			org1,
			[]folder.Folder{
				{"alpha", org1, "alpha"},
			},
			[]folder.Folder{},
		},
		{
			"alpha",
			org1,
			[]folder.Folder{
				{"alpha", org1, "alpha"},
				{"beta", org1, "alpha.beta"},
			},
			[]folder.Folder{
				{"beta", org1, "alpha.beta"},
			},
		},
		{
			"alpha",
			org1,
			[]folder.Folder{
				{"beta", org1, "beta"},
				{"alpha", org1, "beta.alpha"},
				{"delta", org1, "beta.alpha.delta"},
			},
			[]folder.Folder{
				{"delta", org1, "beta.alpha.delta"},
			},
		},
		{
			"alphaa",
			org1,
			[]folder.Folder{
				{"alphaa", org1, "alphaa"},
				{"alphab", org2, "alphab"},
				{"beta", org1, "alphaa.beta"},
			},
			[]folder.Folder{
				{"beta", org1, "alphaa.beta"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := folder.NewDriver(tt.folders)

			getChild, err := f.GetAllChildFolders(tt.orgID, tt.name)
			assert.Equal(t, checkSliceEqual(getChild, tt.want), true, "got %v, want %v", getChild, tt.want)
			assert.Nil(t, err)
		})
	}

	errorTests := [...]struct {
		name    string
		orgID   uuid.UUID
		folders []folder.Folder
		err     error
	}{
		{
			"alpha",
			org1,
			[]folder.Folder{
				{"beta", org1, "beta"},
			},
			errors.New("Error: Folder does not exist"),
		},
		{
			"alpha",
			org1,
			[]folder.Folder{
				{"alalpha", org1, "alalpha"},
			},
			errors.New("Error: Folder does not exist"),
		},
		{
			"alpha",
			org1,
			[]folder.Folder{
				{"alalpha", org1, "alalpha"},
				{"beta", org1, "alalpha.beta"},
			},
			errors.New("Error: Folder does not exist"),
		},
		{
			"alpha",
			org1,
			[]folder.Folder{
				{"beta", org1, "beta"},
				{"alalpha", org1, "beta.alalpha"},
				{"delta", org1, "beta.alalpha.delta"},
			},
			errors.New("Error: Folder does not exist"),
		},
		{
			"alpha",
			org1,
			[]folder.Folder{
				{"delta", org1, "delta"},
				{"alpha", org2, "alpha"},
				{"beta", org2, "alpha.beta"},
			},
			errors.New("Error: Folder does not exist in the specified organization"),
		},
	}
	for _, tt := range errorTests {
		t.Run(tt.name, func(t *testing.T) {
			f := folder.NewDriver(tt.folders)

			getChild, err := f.GetAllChildFolders(tt.orgID, tt.name)

			assert.Error(t, err)
			assert.Equal(t, tt.err, err)
			assert.Nil(t, getChild)
		})
	}
}
