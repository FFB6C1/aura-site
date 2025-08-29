package builder

import (
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/ffb6c1/aura-site/internal/config"
	"github.com/ffb6c1/aura-site/internal/file"
)

func GetThemes() (map[string]config.ThemeSettings, []string, error) {
	themes := map[string]config.ThemeSettings{}
	themeNames := []string{}
	files, err := os.ReadDir("components/css")
	if err != nil {
		return nil, nil, err
	}

	for _, item := range files {
		if file.CheckType(item.Name(), ".css") {
			theme, err := file.FileToString(filepath.Join("components/css", item.Name()))
			if err != nil {
				return nil, nil, err
			}
			settings := getSettings(theme)
			themes[settings.GetName()] = settings
			themeNames = append(themeNames, settings.GetName())
		}
	}

	config.GetConfig().SetThemes(themes)

	return themes, themeNames, nil
}

func getSettings(css string) config.ThemeSettings {
	return splitCommentToSettings(splitComment(getComment(css)), css)
}

func getComment(css string) string {
	return strings.Trim(strings.Split(css, "*/")[0], "/* \n")
}

func splitComment(comment string) map[string]string {
	commentMap := map[string]string{}
	commentSplit := strings.Split(comment, "\n")
	for _, field := range commentSplit {
		keyAndVal := strings.Split(field, ":")
		if len(keyAndVal) != 2 {
			continue
		}
		commentMap[strings.TrimSpace(keyAndVal[0])] = strings.TrimSpace(keyAndVal[1])
	}
	return commentMap
}

func splitCommentToSettings(commentMap map[string]string, css string) config.ThemeSettings {
	priority := omitOrSplit(commentMap["priority"], ", ")
	contentIndex := slices.Index(priority, "content")
	settings := config.NewThemeSettings(
		commentMap["name"],
		commentMap["about"],
		css,
		omitOrSplit(commentMap["required"], ", "),
		omitOrSplit(commentMap["optional"], ", "),
		priority,
		contentIndex,
	)
	return settings
}

func omitOrSplit(text, split string) []string {
	if text == "" {
		return []string{}
	}
	return strings.Split(text, split)
}
