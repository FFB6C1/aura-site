package file

import (
	"fmt"
	"path/filepath"
)

type pageTemplate struct {
	pages         map[string]string
	templateStart string
	templateEnd   string
}

func buildPagesFromTemplate(pageTemplate pageTemplate, exportPath string) error {
	for name, content := range pageTemplate.pages {
		content := pageTemplate.templateStart + content + pageTemplate.templateEnd
		filePath := filepath.Join(exportPath, name+".html")
		if err := writeFile(content, filePath); err != nil {
			return err
		}
	}
	return nil
}

func buildTemplate(path string, themeSettings themeSettings) (pageTemplate, error) {
	files, err := getFiles(path)
	if err != nil {
		return pageTemplate{}, err
	}

	// find content in priority list
	contentIndex := -1
	for i, item := range themeSettings.priority {
		if item == "content" {
			contentIndex = i
		}
	}
	if contentIndex == -1 {
		return pageTemplate{}, fmt.Errorf("content missing from priority list in theme settings")
	}

	// check for presence of all required items
	if item, ok := requiredChecker(themeSettings.required, files); !ok {
		return pageTemplate{}, fmt.Errorf("missing required file for theme: %s", item)
	}

	wrapper, err := getWrapper()
	if err != nil {
		return pageTemplate{}, err
	}

	// build templateStart

	templateStart := wrapper[0]
	for _, item := range themeSettings.priority[:contentIndex] {
		templateStart += converter(files[item], item)
		delete(files, item)
	}

	//build templateEnd
	templateEnd := ""
	if contentIndex+1 <= len(themeSettings.priority) {
		for _, item := range themeSettings.priority[contentIndex+1:] {
			templateEnd += converter(files[item], item)
			delete(files, item)
		}
	}
	templateEnd += wrapper[1]

	pages := map[string]string{}

	for key, value := range files {
		pages[key] = converter(value, key)
	}

	return pageTemplate{
		pages:         pages,
		templateStart: templateStart,
		templateEnd:   templateEnd,
	}, nil
}

func requiredChecker(required []string, files map[string]string) (string, bool) {
	if len(required) == 0 {
		return "", true
	}
	for _, item := range required {
		if files[item] == "" {
			return item, false
		}
	}
	return "", true
}
