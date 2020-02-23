package gui

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

type GUI struct {
	App    *interface{}
	Window *interface{}
}

// InitializeHome sets up the home gui
func InitializeHome() {
	app := app.New()

	w := app.NewWindow("CubbyBot v1")
	w.SetFullScreen(true)

	//windowSize := &canvas.Rectangle{}
	//windowSize.SetMinSize(fyne.NewSize(800, 480))
	rightPanel := widget.NewVBox(
		widget.NewLabel("CubbyBot v1"),
		widget.NewButton("Quit", func() {
			app.Quit()
		}),
	)

	fyne.NewCo

	w.SetContent(fyne.NewContainer(layout.NewBorderLayout(nil, nil, nil, rightPanel)))

	w.ShowAndRun()
}
