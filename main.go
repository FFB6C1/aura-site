package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/ffb6c1/aura-site/internal/builder"
	"github.com/ffb6c1/aura-site/internal/config"
	"github.com/ffb6c1/aura-site/internal/gui"
)

func main() {
	app := app.New()
	win := app.NewWindow("aura-site")

	win.SetContent(gui.GetMainScreen(win))
	win.Resize(fyne.NewSize(640, 480))
	win.ShowAndRun()

	config := config.GetConfig()

	config.SetImportPath("/home/chesca/workspace/aura-site/test-site/src")
	config.SetExportPath("/home/chesca/workspace/aura-site/test-site/dest")
	themes, themeNames, err := builder.GetThemes()
	if err != nil {
		log.Fatal(err)
	}
	config.SetThemes(themes)
	config.SetSelectedTheme(themeNames[0])

	if err := builder.Build(); err != nil {
		log.Fatal(err)
	}
}
