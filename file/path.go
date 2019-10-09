package file

import (
	"errors"
	"os"
	"io/ioutil"
)

type LinkTreatment int

const (
	CheckFollowSymlink LinkTreatment = iota
	CheckSymlinkOnly
)

var ErrInvalidLinkTreatment = errors.New("unknown link behavior")

// Exists
func Exists(linkBehavior LinkTreatment, filename string) (bool, error) {
	var err error

	if linkBehavior == CheckFollowSymlink {
		_, err = os.Stat(filename)
	} else if linkBehavior == CheckSymlinkOnly {
		_, err = os.Lstat(filename)
	} else {
		return false, ErrInvalidLinkTreatment
	}

	if os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

// ReadDirNoStat
func ReadDirNoStat(dirname string) ([]string, error) {
	if dirname == "" {
		dirname = "."
	}

	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return f.Readdirnames(-1)
}

// ListAllFileByName
func ListAllFileByName(level int, pathSeparator string, fileDir string) []string {
	files, _ := ioutil.ReadDir(fileDir)
	fileList := make([]string, 0)
	for _, onefile := range files {
		if (onefile.IsDir()) {
			ListAllFileByName(level+1, pathSeparator, fileDir+pathSeparator+onefile.Name())
		} else {
			fileList = append(fileList, onefile.Name())
		}
	}
	return fileList
}