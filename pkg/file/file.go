package file

import "os"

// DirExists checks if a dir is existed
func DirExists(configPath string) bool {
	fi, err := os.Stat(configPath)
	if os.IsNotExist(err) {
		return false
	}
	return !fi.IsDir()
}
