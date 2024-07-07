package ui

import "fyne.io/fyne/v2"

func Clipboard() fyne.Clipboard {
	return fyne.CurrentApp().Driver().AllWindows()[0].Clipboard()
}

func CanvasForObject(cp fyne.CanvasObject) fyne.Canvas {
	return fyne.CurrentApp().Driver().CanvasForObject(cp)
}
