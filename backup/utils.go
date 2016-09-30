package backup

import "os"

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func FileKey(filepath string) (key string, err error) {
	stat, err := os.Stat(filepath)
	if err != nil {
		return
	}

	key = stat.Name()
	return
}
