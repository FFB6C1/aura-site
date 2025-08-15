package gui

import (
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
			sourceLabel.SetText("source: " + sourceFolder)
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
			destLabel.SetText("destination: " + destFolder)
		}, window)
	})
	convertButton := widget.NewButton("convert .md files to .html", func() {
		file.ConvertFilesFromFolder(sourceFolder, destFolder)
	})

	return container.NewVBox(sourceLabel, destLabel, sourceFolderButton, destFolderButton, convertButton)

}
