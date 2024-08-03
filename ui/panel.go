package ui

import (
	"encoding/json"
	"fmt"
	"fynecv/appdata"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type AddTab struct {
	btn *widget.Button
	sel *widget.Select
	ctr fyne.CanvasObject
}

func newAddTab() *AddTab {
	atab := &AddTab{
		sel: widget.NewSelect([]string{}, func(s string) {}),
		btn: widget.NewButton("Select Entity", func() {}),
	}
	atab.ctr = container.NewBorder(nil, nil,
		atab.btn, nil, atab.sel)

	return atab
}

func NewTabs(data *appdata.AppData, win fyne.Window) *container.AppTabs {
	tabs := container.NewAppTabs()

	lightIDs := []string{
		"light.led_matrix_24",
		"light.led_strip_24"}

	for _, id := range lightIDs {
		tab := container.NewTabItem(makeLabel(id), NewLightPanel(id, win, data))
		tabs.Append(tab)
	}
	atab := newAddTab()

	addTab := container.NewTabItem("add item", atab.ctr)
	addTab.Icon = theme.ContentAddIcon()
	tabs.Append(addTab)
	tabs.OnSelected = func(ti *container.TabItem) {
		atab.sel.SetOptions(data.EntityList())
	}

	atab.btn.OnTapped = func() {
		s := atab.sel.Selected
		if s == "" {
			return
		}

		src, ok := data.Entities[s]
		if !ok {
			return
		}

		ctr := NewAnyPanel(src)
		tab := container.NewTabItem(makeLabel(src.EntityID), ctr)
		tabs.Remove(addTab)
		tabs.Append(tab)
		tabs.Append(addTab)
		tabs.Select(tab)
	}
	return tabs
}

func NewAnyPanel(entity *appdata.Entity[json.RawMessage]) fyne.CanvasObject {

	var anyData appdata.AnyData
	anyData.Copy(entity)
	adj := container.NewHBox()
	adj.Add(container.NewHBox(widget.NewLabel("State:"),
		widget.NewLabel(fmt.Sprint(anyData.State))))
	adj.Add(widget.NewSeparator())
	att, _ := anyData.Attributes.(map[string]any)
	for k, v := range att {
		adj.Add(container.NewHBox(widget.NewLabel(makeLabel(k)),
			widget.NewLabel(fmt.Sprint(v))))
		adj.Add(widget.NewSeparator())
	}
	ctr := container.NewHScroll(adj)
	return ctr
}

func makeLabel(s string) string {
	slen := len(s)
	var lab string
	for i := 0; i < slen; {
		if i == 0 {
			lab += strings.ToUpper(s[:1])
			i++
			continue
		}

		if s[i] == '_' || s[i] == '.' {
			lab += " "
			i++
			if i >= slen {
				break
			}
			lab += strings.ToUpper(s[i : i+1])
			i++
			continue
		}
		lab += s[i : i+1]
		i++
	}
	lab += ":"
	return lab
}
