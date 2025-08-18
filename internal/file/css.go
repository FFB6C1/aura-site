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
		required: strings.Split(commentMap["required"], ", "),
		optional: strings.Split(commentMap["optional"], ", "),
		priority: strings.Split(commentMap["priority"], ", "),
		about:    commentMap["about"],
		full:     css,
	}
	return settings
}
