package file

import (
	"fmt"
	"strings"
)

type themeSettings struct {
	name     string
	required []string
	optional []string
	priority []string
	about    string
	full     string
}

func (t themeSettings) GetRequired() []string {
	return t.required
}

func getSettings(css string) themeSettings {
	return splitCommentToSettings(splitComment(getComment(css)), css)
}

func getComment(css string) string {
	return strings.Trim(strings.Split(css, "*/")[0], "/* \n")
}

func splitComment(comment string) map[string]string {
	commentMap := map[string]string{}
	commentSplit := strings.Split(comment, "\n")
	fmt.Println(commentSplit)
	for _, field := range commentSplit {
		keyAndVal := strings.Split(field, ":")
		if len(keyAndVal) != 2 {
			continue
		}
		commentMap[strings.TrimSpace(keyAndVal[0])] = strings.TrimSpace(keyAndVal[1])
	}
	return commentMap
}

func splitCommentToSettings(commentMap map[string]string, css string) themeSettings {
	settings := themeSettings{
		name:     commentMap["name"],
		required: omitOrSplit(commentMap["required"], ", "),
		optional: omitOrSplit(commentMap["optional"], ", "),
		priority: omitOrSplit(commentMap["priority"], ", "),
		about:    commentMap["about"],
		full:     css,
	}
	return settings
}

func omitOrSplit(text, split string) []string {
	if text == "" {
		return []string{}
	}
	return strings.Split(text, split)
}
