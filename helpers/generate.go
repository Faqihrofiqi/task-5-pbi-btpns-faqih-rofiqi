package helpers

import (
	"io"
	"os"
	"strconv"
)

func CreateDirIfNotExist(dirPath string) (bool, error) {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err := os.Mkdir(dirPath, 0755)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

func CreateFile(filePath string) (*os.File, error) {
	return os.Create(filePath)
}

func CopyFile(destination io.Writer, source io.Reader) (int64, error) {
	return io.Copy(destination, source)
}

func GenerateUniqueString(userID uint, photoID uint) string {
	return strconv.FormatUint(uint64(userID), 10) + "_" + strconv.FormatUint(uint64(photoID), 10)
}
