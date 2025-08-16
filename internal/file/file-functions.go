package file

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ffb6c1/aura-site/internal/markdown"
)

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

func readFiles(path string) (map[string]string, error) {
	dir, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open directory: %w", err)
	}
	defer dir.Close()

	files, err := dir.ReadDir(-1)
	if err != nil {
		return nil, fmt.Errorf("could not read files from directory: %w", err)
	}

	markdownFiles := make(map[string]string, len(files))

	for _, file := range files {
		if filename, markdown := isMarkdown(file.Name()); markdown {
			path := filepath.Join(path, file.Name())
			content, err := readFile(path)
			if err != nil {
				return nil, fmt.Errorf("could not read file %s: %w", path, err)
			}

			markdownFiles[filename] = content
		}

	}

	return markdownFiles, nil
}

func isMarkdown(name string) (string, bool) {
	return strings.CutSuffix(name, ".md")
}

func converter(markdownFiles map[string]string) (map[string]string, error) {
	htmlFiles := make(map[string]string)
	wrapper, err := os.ReadFile("components/html/wrapper.html")
	wrapperParts := strings.Split(string(wrapper), "!!!CONTENT!!!")
	if err != nil {
		return nil, fmt.Errorf("could not open html wrapper: %w", err)
	}
	nav := markdownFiles["nav"]
	delete(markdownFiles, "nav")
	if nav != "" {
		nav = markdown.Convert(nav) + "\n\n"
	}

	for name, content := range markdownFiles {
		convertedContent := markdown.Convert(content)

		pageParts := []string{
			nav, convertedContent,
		}

		htmlFiles[name+".html"] = pageBuilder(pageParts, wrapperParts)

	}
	return htmlFiles, nil
}

func pageBuilder(pageParts []string, wrapperParts []string) string {
	fullPage := ""
	for _, part := range pageParts {
		fullPage += part
	}
	return wrapperParts[0] + fullPage + wrapperParts[1]
}

func writeFiles(htmlFiles map[string]string, path string) error {
	for name, content := range htmlFiles {
		fullPath := filepath.Join(path, name)
		if err := os.WriteFile(fullPath, []byte(content), 0o777); err != nil {
			return fmt.Errorf("error writing file %s: %w", name, err)
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
