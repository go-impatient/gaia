package file

import "os"

// DirExists checks if a dir is existed
func DirExists(configPath string) bool {
	fi, err := os.Stat(configPath)
	if err != nil {
		return false
	}
	return fi.IsDir()
}
