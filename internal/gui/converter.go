package gui

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/ffb6c1/aura-site/internal/file"
)

func MainScreen(window fyne.Window) fyne.CanvasObject {
	sourceFolder := ""
	destFolder := ""

	themes, themeNames, err := file.GetThemes()
	if err != nil {
		log.Fatal(err.Error())
	}
	requiredLabel := widget.NewLabel("Required files for chosen theme:")

	themeSelector := widget.NewSelect(themeNames, func(theme string) {
		reqFiles := ""
		if len(themes[theme].GetRequired()) == 0 {
			reqFiles = "none"
		} else {
			for _, file := range themes[theme].GetRequired() {
				reqFiles += file + " "
			}
		}
		requiredLabel.SetText(fmt.Sprintf("Required files for chosen theme: %s", reqFiles))
	})

	sourceLabel := widget.NewLabel("Choose a source folder")
	destLabel := widget.NewLabel("Choose a destination folder")

	sourceFolderButton := widget.NewButton("select source folder", func() {
		dialog.ShowFolderOpen(func(list fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, window)
				return
			}
			if list == nil {
				log.Println("Cancelled")
				return
			}

			sourceFolder = list.Path()
			sourceLabel.SetText(sourceFolder)
		}, window)
	})
	destFolderButton := widget.NewButton("select destination folder", func() {
		dialog.ShowFolderOpen(func(list fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, window)
				return
			}
			if list == nil {
				log.Println("Cancelled")
				return
			}

			destFolder = list.Path()
			destLabel.SetText(destFolder)
		}, window)
	})
	convertButton := widget.NewButton("convert .md files to .html", func() {
		if err := file.BuildSite(sourceFolder, destFolder, "components/css/default.css"); err != nil {
			log.Fatal(err.Error())
		}
	})

	label := container.NewCenter(widget.NewLabelWithStyle("aura-site website builder", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
	source := container.NewHBox(sourceFolderButton, sourceLabel)
	destination := container.NewHBox(destFolderButton, destLabel)
	theme := container.NewVBox(themeSelector, requiredLabel)

	screen := container.NewVBox(label, theme, source, destination, convertButton)

	return screen

}
