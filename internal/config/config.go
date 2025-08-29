package config

type configStruct struct {
	importPath string
	exportPath string
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
