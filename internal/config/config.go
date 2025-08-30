package config

type configStruct struct {
	importPath    string
	exportPath    string
	themes        map[string]ThemeSettings
	selectedTheme string
	template      []string
}

var config configStruct = configStruct{
	importPath: "",
	exportPath: "",
}

func GetConfig() *configStruct {
	return &config
}

func (c *configStruct) SetImportPath(path string) {
	c.importPath = path
}

func (c configStruct) GetImportPath() string {
	return c.importPath
}

func (c *configStruct) SetExportPath(path string) {
	c.exportPath = path
}

func (c configStruct) GetExportPath() string {
	return c.exportPath
}

func (c *configStruct) SetThemes(themes map[string]ThemeSettings) {
	c.themes = themes
}

func (c configStruct) GetThemes() map[string]ThemeSettings {
	return c.themes
}

func (c configStruct) GetTheme(name string) (ThemeSettings, bool) {
	theme, ok := c.themes[name]
	return theme, ok
}

func (c *configStruct) SetSelectedTheme(theme string) {
	c.selectedTheme = theme
}

func (c configStruct) GetSelectedTheme() (ThemeSettings, bool) {
	if c.selectedTheme == "" {
		return ThemeSettings{}, false
	}
	theme, ok := c.themes[c.selectedTheme]
	return theme, ok
}

func (c configStruct) GetCSS() string {
	return c.themes[c.selectedTheme].full
}

func (c *configStruct) SetTemplate(templateStart, templateEnd string) {
	c.template = append(c.template, templateStart, templateEnd)
}

func (c configStruct) GetTemplate() []string {
	return c.template
}
