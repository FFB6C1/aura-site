package builder

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ffb6c1/aura-site/internal/config"
	"github.com/ffb6c1/aura-site/internal/file"
)

func Build() error {
	config := config.GetConfig()
	css := config.GetCSS()

	if err := file.WriteFileFromString(filepath.Join(config.GetExportPath(), "styles.css"), css); err != nil {
		return err
	}

	mdPages, err := file.GetMDFiles(config.GetImportPath())
	if err != nil {
		return err
	}

	galleries := []gallery{}

	htmlPages := map[string]string{}
	for name, content := range mdPages {
		html, gals := checkGallery(content)
		htmlPages[name] = html
		galleries = append(galleries, gals...)
	}

	finalPages, err := buildTemplate(htmlPages)
	if err != nil {
		return err
	}

	if err := createPages(finalPages, config.GetExportPath()); err != nil {
		return err
	}

	finishGalleries(galleries)

	return nil
}

func buildTemplate(pages map[string]string) (map[string]string, error) {
	theme, ok := config.GetConfig().GetSelectedTheme()
	if !ok {
		return nil, fmt.Errorf("no or unknown theme selected")
	}
	haveFiles, absentFiles := theme.CheckRequiredFromMap(pages)
	if !haveFiles {
		return nil, fmt.Errorf("missing required files for selected theme: %v", absentFiles)
	}

	wrapper, err := getHTMLWrapper()
	if err != nil {
		return nil, err
	}

	templateStart := wrapper[0]
	templateEnd := "\n</div>"

	priority, contentIndex := theme.GetPriorityAndContentIndex()
	if contentIndex == -1 {
		return nil, fmt.Errorf("content not listed in theme priority settings")
	}

	for _, page := range priority[:contentIndex] {
		if page == "MAINSTART" {
			templateStart += "\n<div id=\"main\">\n"
			continue
		}
		if page == "MAINEND" {
			templateStart += "\n</div>\n"
			continue
		}
		templateStart += fmt.Sprintf("<div id=\"%s\">\n%s\n</div>\n", page, pages[page])
		delete(pages, page)
	}
	templateStart += "\n<div id=\"content\">"

	for _, page := range priority[contentIndex+1:] {
		if page == "MAINSTART" {
			templateEnd += "\n<div id=\"main\">\n"
			continue
		}
		if page == "MAINEND" {
			templateEnd += "\n</div>\n"
			continue
		}
		templateEnd += fmt.Sprintf("<div id=\"%s\">\n%s\n</div>\n", page, pages[page])
		delete(pages, page)
	}
	templateEnd += wrapper[1]

	config.GetConfig().SetTemplate(templateStart, templateEnd)

	return pages, nil
}

func createPages(pages map[string]string, exportPath string) error {
	if err := file.MakeDirectory(exportPath); err != nil {
		return err
	}
	template := config.GetConfig().GetTemplate()
	for name, page := range pages {
		fullPage := template[0] + page + template[1]
		path := filepath.Join(exportPath, name+".html")
		if err := file.WriteFileFromString(path, fullPage); err != nil {
			return err
		}
	}
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
