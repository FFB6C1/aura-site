package file

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ffb6c1/aura-site/internal/markdown"
)

func getWrapper() ([]string, error) {
	wrapper, err := readFile("components/html/wrapper.html")
	if err != nil {
		return nil, err
	}

	return strings.Split(wrapper, "!!!CONTENT!!!"), nil
}

func getFiles(path string) (map[string]string, error) {
	files, err := os.ReadDir(path)
	fileMap := map[string]string{}
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".md") {
			filePath := filepath.Join(path, file.Name())
			content, err := readFile(filePath)
			if err != nil {
				return nil, err
			}
			fileMap[strings.TrimSuffix(file.Name(), ".md")] = content
		}
	}
	return fileMap, nil
}

func readFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("could not read file at %s: %w", path, err)
	}
	return string(content), err
}

func writeFile(content, path string) error {
	if err := os.WriteFile(path, []byte(content), 0o777); err != nil {
		return fmt.Errorf("could not write file at %s: %w", path, err)
	}
	return nil
}

func converter(content, name string) string {
	converted := markdown.Convert(content, name)
	return converted
}
