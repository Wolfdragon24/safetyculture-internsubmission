package folder

import (
	"errors"
	"regexp"

	"github.com/gofrs/uuid"
)

func GetAllFolders() []Folder {
	return GetSampleData()
}

func (f *driver) GetFoldersByOrgID(orgID uuid.UUID) []Folder {
	folders := f.folders

	res := []Folder{}
	for _, f := range folders {
		if f.OrgId == orgID {
			res = append(res, f)
		}
	}

	return res

}

func (f *driver) GetAllChildFolders(orgID uuid.UUID, name string) ([]Folder, error) {
	folders := f.folders
	maxSample := []Folder{}
	foundFile := false
	sameOrg := false

	// Early error if no folders
	if len(folders) == 0 {
		return nil, errors.New("Error: Folder does not exist")
	}

	// Use regex to filter all paths which contain the name as a parent
	r, _ := regexp.Compile("^(?:[\\w.]+\\.|\\.)?" + name + "\\.")

	for _, f := range folders {
		if r.FindStringIndex(f.Paths) != nil {
			maxSample = append(maxSample, f)
		}
		if f.Name == name {
			foundFile = true
			if f.OrgId == orgID {
				sameOrg = true
			}
		}
	}

	if !foundFile {
		return nil, errors.New("Error: Folder does not exist")
	} else if !sameOrg {
		return nil, errors.New("Error: Folder does not exist in the specified organization")
	}

	return maxSample, nil
}
