package config

import (
	"os"
	"strings"
)

type ThemeSettings struct {
	name         string
	required     []string
	optional     []string
	priority     []string
	about        string
	full         string
	contentIndex int
}

func NewThemeSettings(name, about, full string, required, optional, priority []string, contentIndex int) ThemeSettings {
	return ThemeSettings{
		name:         name,
		about:        about,
		full:         full,
		required:     required,
		optional:     optional,
		priority:     priority,
		contentIndex: contentIndex,
	}
}

func (t ThemeSettings) GetRequired() []string {
	return t.required
}

func (t ThemeSettings) CheckRequired(path string) (bool, []string, error) {
	if path == "" {
		return false, t.required, nil
	}
	files, err := getMDFiles(path)
	if err != nil {
		return false, t.required, err
	}

	absent := []string{}

	for _, file := range t.required {
		_, ok := files[file]
		if !ok {
			absent = append(absent, file)
		}
	}

	if len(absent) != 0 {
		return false, absent, nil
	}

	return true, nil, nil

}

func (t ThemeSettings) CheckRequiredFromMap(files map[string]string) (bool, []string) {
	absent := []string{}

	for _, file := range t.required {
		_, ok := files[file]
		if !ok {
			absent = append(absent, file)
		}
	}

	if len(absent) != 0 {
		return false, absent
	}

	return true, nil
}

func (t ThemeSettings) GetPriorityAndContentIndex() ([]string, int) {
	return t.priority, t.contentIndex
}

func (t ThemeSettings) GetName() string {
	return t.name
}

func (t ThemeSettings) GetAbout() string {
	return t.about
}

func (t ThemeSettings) GetOptional() []string {
	return t.optional
}

func (t ThemeSettings) GetOptionalAsString() string {
	optional := ""
	for _, item := range t.optional {
		optional += ", " + item
	}
	if optional == "" {
		return "None"
	}
	return optional[2:]
}

func (t ThemeSettings) GetRequiredAsString() string {
	required := ""
	for _, item := range t.required {
		required += ", " + item
	}
	if required == "" {
		return "None"
	}
	return required[2:]
}

func getMDFiles(path string) (map[string]string, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	mdFiles := make(map[string]string)
	for _, item := range files {
		name, isMD := strings.CutSuffix(item.Name(), ".md")
		if !isMD {
			continue
		}
		mdFiles[name] = item.Name()
	}
	return mdFiles, nil
}
