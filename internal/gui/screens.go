package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/ffb6c1/aura-site/internal/builder"
	"github.com/ffb6c1/aura-site/internal/config"
)

func GetMainScreen(win fyne.Window) fyne.CanvasObject {

	themeSelector := GetThemeSelector(win)

	screen := container.NewVBox(themeSelector)
	return screen
}

func GetThemeSelector(win fyne.Window) fyne.CanvasObject {
	config := config.GetConfig()
	themes, themeNames, err := builder.GetThemes()
	if err != nil {
		fyne.LogError("could not retrieve themes", err)
	}
	config.SetThemes(themes)
	themeLabel := widget.NewLabel("")
	themeLabel.Wrapping = fyne.TextWrapBreak

	selector := widget.NewSelect(themeNames, func(t string) {
		theme := themes[t]
		config.SetSelectedTheme(t)
		themeLabel.SetText(fmt.Sprintf("%s: %s\n\nRequired files (md): %s\nOptional files (md): %s\n", theme.GetName(), theme.GetAbout(), theme.GetRequiredAsString(), theme.GetOptionalAsString()))
	})
	return container.NewVBox(selector, themeLabel)
}
