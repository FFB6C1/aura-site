package file

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func MakeDirectory(path string) error {
	if err := os.Mkdir(path, 0o777); err != nil {
		if os.IsExist(err) {
			return nil
		}
		return err
	}
	return nil
}

func GetMDFiles(path string) (map[string]string, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	mdFiles := make(map[string]string)
	for _, file := range files {
		name, isMD := strings.CutSuffix(file.Name(), ".md")
		if !isMD {
			continue
		}
		body, err := FileToString(filepath.Join(path, file.Name()))
		if err != nil {
			fmt.Printf("file not read: %s. reason: %s", file.Name(), err.Error())
			continue
		}
		mdFiles[name] = body
	}
	return mdFiles, nil
}

func GetImgFiles(path string) ([][]string, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	imgFiles := [][]string{}
	for _, file := range files {
		fileWithType := strings.Split(file.Name(), ".")
		if fileWithType[1] == "jpg" ||
			fileWithType[1] == "jpeg" ||
			fileWithType[1] == "png" ||
			fileWithType[1] == "gif" ||
			fileWithType[1] == "webp" {
			imgFiles = append(imgFiles, fileWithType)
		}
	}
	return imgFiles, nil
}

func FileToString(path string) (string, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(file), nil
}

func CheckType(filename, filetype string) bool {
	return strings.HasSuffix(filetype, filetype)
}
