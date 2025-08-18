package file

import (
	"path/filepath"
)

func BuildSite(importPath, exportPath, cssPath string) error {
	css, err := readFile(cssPath)
	if err != nil {
		return err
	}

	themeSettings := splitCommentToSettings(splitComment(getComment(css)), css)

	template, err := buildTemplate(importPath, themeSettings)
	if err != nil {
		return err
	}

	if err := buildPagesFromTemplate(template, exportPath); err != nil {
		return err
	}

	if err := writeFile(css, filepath.Join(exportPath, "styles.css")); err != nil {
		return err
	}

	return nil
}
