package file

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ffb6c1/aura-site/internal/markdown"
)

type markdownFile struct {
	name    string
	content string
}

type htmlFile struct {
	name    string
	content string
}

func ConvertFilesFromFolder(importPath string, exportPath string) error {
	markdownFiles, err := readFiles(importPath)
	if err != nil {
		return err
	}
	htmlFiles, err := converter(markdownFiles)
	if err != nil {
		return err
	}
	if err := writeFiles(htmlFiles, exportPath); err != nil {
		return err
	}
	return nil
}

func readFiles(path string) ([]markdownFile, error) {
	dir, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open directory: %w", err)
	}
	defer dir.Close()

	files, err := dir.ReadDir(-1)
	if err != nil {
		return nil, fmt.Errorf("could not read files from directory: %w", err)
	}

	markdownFiles := make([]markdownFile, 0, len(files))

	for _, file := range files {
		if filename, markdown := isMarkdown(file.Name()); markdown {
			path := filepath.Join(path, file.Name())
			content, err := readFile(path)
			if err != nil {
				return nil, fmt.Errorf("could not read file %s: %w", path, err)
			}

			mdFile := markdownFile{
				name:    filename,
				content: content,
			}
			markdownFiles = append(markdownFiles, mdFile)
		}

	}

	return markdownFiles, nil
}

func isMarkdown(name string) (string, bool) {
	return strings.CutSuffix(name, ".md")
}

func converter(markdownFiles []markdownFile) ([]htmlFile, error) {
	htmlFiles := make([]htmlFile, 0, len(markdownFiles))
	wrapper, err := os.ReadFile("components/html/wrapper.html")
	wrapperParts := strings.Split(string(wrapper), "!!!CONTENT!!!")
	if err != nil {
		return nil, fmt.Errorf("could not open html wrapper: %w", err)
	}
	for _, file := range markdownFiles {
		content := markdown.Convert(file.content)

		html := htmlFile{
			name:    file.name + ".html",
			content: wrapperParts[0] + content + wrapperParts[1],
		}
		htmlFiles = append(htmlFiles, html)
	}
	return htmlFiles, nil
}

func writeFiles(htmlFiles []htmlFile, path string) error {
	for _, file := range htmlFiles {
		fullPath := filepath.Join(path, file.name)
		if err := os.WriteFile(fullPath, []byte(file.content), 0o777); err != nil {
			return fmt.Errorf("error writing file %s: %w", file.name, err)
		}
	}
	return nil
}

func readFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
