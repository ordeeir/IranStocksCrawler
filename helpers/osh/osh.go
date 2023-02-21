package osh

import (
	"os"
	"path/filepath"
	"runtime"
)

// IsFileExist is a function to check FileIsExist
func IsFileExist(file string) (bool, error) {
	_, err := os.Stat(file)
	if err != nil {
		return false, err
	}
	return true, nil
}

/*func GetRootPath() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return dir
}*/

func GetRootPath() string {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(filepath.Dir(filepath.Dir(b)))
	dir, _ := filepath.Abs(basepath)
	return dir
}
