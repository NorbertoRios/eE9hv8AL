package utils

import "path/filepath"

//GetAbsFilePath returns absolute file path
func GetAbsFilePath(dir, fileDest string) string {
	dir = filepath.Dir(dir)
	absPath, _ := filepath.Abs(dir + fileDest)
	return absPath
}
