package ui

import (
	"fynecv/appdata"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type Panel struct {
	data *appdata.AppData
	Tabs *container.DocTabs
}

func NewPanel(data *appdata.AppData, win fyne.Window) *Panel {
	lp := &Panel{
		data: data,
		Tabs: container.NewDocTabs(),
	}

	lightIDs := []string{
		"light.led_matrix_24",
		"light.led_strip_24"}

	for _, id := range lightIDs {
		tab := container.NewTabItem(id,
			NewLightPanel(id, win, data))
		lp.Tabs.Append(tab)
	}
	return lp
}
