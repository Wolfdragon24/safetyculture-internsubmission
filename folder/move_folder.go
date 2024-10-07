package folder

import (
	"errors"
)

func getIndex(folders []Folder, name string) int {
	for i, f := range folders {
		if f.Name == name {
			return i
		}
	}
	return -1
}

func (f *driver) MoveFolder(name string, dst string) ([]Folder, error) {
	folders := f.folders

	if len(folders) == 0 {
		return nil, errors.New("Error: Source folder does not exist")
	} else if name == dst {
		return nil, errors.New("Error: Cannot move a folder to itself")
	}

	// Get and check the existence of the source and destination folders
	srcFolder := -1
	destFolder := -1

	for i, f := range folders {
		if f.Name == name {
			srcFolder = i
		}
		if f.Name == dst {
			destFolder = i
		}
	}

	if srcFolder == -1 {
		return nil, errors.New("Error: Source folder does not exist")
	} else if destFolder == -1 {
		return nil, errors.New("Error: Destination folder does not exist")
	} else if folders[srcFolder].OrgId != folders[destFolder].OrgId {
		return nil, errors.New("Error: Cannot move a folder to a different organization")
	}

	childFolders, _ := f.GetAllChildFolders(folders[srcFolder].OrgId, name)

	// Check that destination is not child of parent
	for _, f := range childFolders {
		if f.Name == dst {
			return nil, errors.New("Error: Cannot move a folder to a child of itself")
		}
	}

	// Move folder and child folders

	folders[srcFolder].Paths = folders[destFolder].Paths + "." + folders[srcFolder].Name

	for _, f := range childFolders {
		folderIndex := getIndex(folders, f.Name)

		folders[folderIndex].Paths = folders[srcFolder].Paths + "." + f.Name
	}

	updatedFolders := f.folders

	return updatedFolders, nil
}
