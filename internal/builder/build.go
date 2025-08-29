package builder

import (
	"fmt"
	"strings"

	"github.com/ffb6c1/aura-site/internal/config"
	"github.com/ffb6c1/aura-site/internal/file"
	"github.com/ffb6c1/aura-site/internal/markdown"
)

func Build() error {
	config := config.GetConfig()

	mdPages, err := file.GetMDFiles(config.GetImportPath())
	if err != nil {
		return err
	}

	htmlPages := map[string]string{}
	for name, content := range mdPages {
		htmlPages[name] = markdown.Convert(content, "builder")
	}

	if err := buildTemplate(htmlPages); err != nil {
		return err
	}

	return nil
}

func buildTemplate(pages map[string]string) error {
	theme, ok := config.GetConfig().GetSelectedTheme()
	if !ok {
		return fmt.Errorf("no or unknown theme selected")
	}
	haveFiles, absentFiles := theme.CheckRequiredFromMap(pages)
	if !haveFiles {
		return fmt.Errorf("missing required files for selected theme: %v", absentFiles)
	}

	wrapper, err := getHTMLWrapper()
	if err != nil {
		return err
	}

	templateStart := wrapper[0]
	templateEnd := wrapper[1]

	priority, contentIndex := theme.GetPriorityAndContentIndex()
	if contentIndex == -1 {
		return fmt.Errorf("content not listed in theme priority settings")
	}

	for _, page := range priority[:contentIndex] {
		templateStart += pages[page]
	}
	templateStart += "\n<div id=\"content\">"
	for _, page := range priority[contentIndex+1:] {
		templateEnd = pages[page] + templateEnd
	}
	templateEnd = "\n</div>" + templateEnd

	config.GetConfig().SetTemplate(templateStart, templateEnd)

	return nil
}

func getHTMLWrapper() ([]string, error) {
	html, err := file.FileToString("components/html/wrapper.html")
	if err != nil {
		return nil, err
	}
	wrapper := strings.Split(html, "!!!CONTENT!!!")

	if len(wrapper) != 2 {
		return nil, fmt.Errorf("malformed html wrapper")
	}

	return wrapper, nil
}
