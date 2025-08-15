package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/ffb6c1/aura-site/internal/gui"
)

func main() {
	app := app.New()
	win := app.NewWindow("aura-site")

	win.SetContent(gui.MainScreen(win))
	win.Resize(fyne.NewSize(640, 480))
	win.ShowAndRun()
}
