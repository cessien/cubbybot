package gui

import (
	"fmt"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
)

type GUI struct {
	App    *interface{}
	Window *interface{}
}

var (
	keyLeftEvent  = &fyne.KeyEvent{Name: fyne.KeyLeft}
	keyRightEvent = &fyne.KeyEvent{Name: fyne.KeyRight}
	keyUpEvent    = &fyne.KeyEvent{Name: fyne.KeyUp}
	keyDownEvent  = &fyne.KeyEvent{Name: fyne.KeyTab}
	keyEnterEvent = &fyne.KeyEvent{Name: fyne.KeyEnter}
	keySpaceEvent = &fyne.KeyEvent{Name: fyne.KeySpace}

	w fyne.Window
)

// InitializeHome sets up the home gui
func InitializeHome() {
	app := app.New()

	w = app.NewWindow("CubbyBot v1")
	w.Resize(fyne.NewSize(800, 600))

	button := widget.NewButton("shutdown", func() {
		fmt.Println("shutdown initiated")
	})

	enableCamera := widget.NewCheck("enable camera", func(enabled bool) {
		if enabled {
			fmt.Println("camera enabled")
		} else {
			fmt.Println("camera disabled")
		}
	})

	w.SetContent(widget.NewVBox(
		widget.NewLabel("CubbyBot v1 Options"),
		enableCamera,
		button,
	))

	w.ShowAndRun()

	time.Sleep(5 * time.Second)
	button.SetText("Better")

	w.ShowAndRun()
	w.SetFullScreen(true)
	w.CenterOnScreen()
	w.SetFixedSize(true)
	w.Show()

	time.Sleep(5 * time.Second)
	w.RequestFocus()
	w.Canvas().Focused().TypedKey(keyDownEvent)
	time.Sleep(time.Second)
	w.Canvas().Focused().TypedKey(keyEnterEvent)
	w.Show()
}

func Down() {
	w.Canvas().Focused().TypedKey(keyDownEvent)
}

func Up() {
	w.Canvas().Focused().TypedKey(keyUpEvent)
}
