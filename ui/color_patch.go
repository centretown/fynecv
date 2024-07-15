package ui

import (
	"encoding/json"
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var _ fyne.Draggable = (*ColorPatch)(nil)
var _ fyne.Tappable = (*ColorPatch)(nil)
var _ fyne.Focusable = (*ColorPatch)(nil)
var _ fyne.Widget = (*ColorPatch)(nil)
var _ desktop.Mouseable = (*ColorPatch)(nil)
var _ desktop.Hoverable = (*ColorPatch)(nil)
var _ fyne.Shortcutable = (*ColorPatch)(nil)

// var _ desktop.Keyable = (*ColorPatch)(nil)

type ColorPatch struct {
	widget.BaseWidget
	onTapped   func() `json:"-"`
	onChanged  func() `json:"-"`
	rectangle  *canvas.Rectangle
	background *canvas.Rectangle
	FillColor  color.NRGBA

	hovered, focused bool
	unused           bool

	Editing bool
}

func NewColorPatch() (patch *ColorPatch) {
	patch = NewColorPatchWithColor(color.NRGBA{R: 127, G: 127, B: 127, A: 255}, nil, nil)
	patch.unused = true
	return
}

func NewColorPatchWithColor(fill color.NRGBA, onTapped func(), onChanged func()) *ColorPatch {
	cp := &ColorPatch{
		FillColor: fill,
		rectangle: canvas.NewRectangle(fill),
		background: canvas.NewRectangle(color.NRGBA{
			R: fill.R / 2, G: fill.G / 2, B: fill.B / 2, A: 36}),
		onChanged: onChanged,
		onTapped:  onTapped,
	}
	cp.ExtendBaseWidget(cp)
	return cp
}

func (cp *ColorPatch) applyPatchTheme() {
	cp.background.FillColor = cp.backgroundColor()
	cp.Refresh()
}

func (cp *ColorPatch) backgroundColor() (c color.Color) {
	r, g, b, a := cp.GetColor().RGBA()
	c = color.NRGBA{R: uint8(r / 2), G: uint8(g / 2), B: uint8(b / 2), A: uint8(a)}
	// switch {
	// case cp.focused:
	// 	c = theme.FocusColor()
	// case cp.hovered:
	// 	c = theme.HoverColor()
	// default:
	// 	c = theme.ButtonColor()
	// }
	return
}

func (cp *ColorPatch) copy() string {
	buf, err := json.Marshal(cp.FillColor)
	if err != nil {
		return ""
	}
	return string(buf)
}

func (cp *ColorPatch) paste(s string) {
	if len(s) < 1 {
		return
	}

	b := []byte(s)
	var rgb color.NRGBA
	err := json.Unmarshal(b, &rgb)
	if err != nil {
		return
	}
	cp.FillColor = rgb
	cp.setChanged()
}

func (cp *ColorPatch) setChanged() {
	if cp.onChanged != nil {
		cp.onChanged()
	}
}

func (cp *ColorPatch) SetOnChanged(onChanged func()) {
	cp.onChanged = onChanged
}

// Shortcutable
func (cp *ColorPatch) TypedShortcut(sc fyne.Shortcut) {
	switch p := sc.(type) {
	case *fyne.ShortcutCopy:
		p.Clipboard.SetContent(cp.copy())
	case *fyne.ShortcutPaste:
		cp.paste(p.Clipboard.Content())
	case *fyne.ShortcutCut:
		p.Clipboard.SetContent(cp.copy())
		cp.SetUnused(true)
	default:
	}
}

// MouseIn is called when a desktop pointer enters the widget
func (cp *ColorPatch) MouseIn(*desktop.MouseEvent) {
	cp.hovered = true
	cp.applyPatchTheme()
}

// MouseMoved is called when a desktop pointer hovers over the widget
func (cp *ColorPatch) MouseMoved(*desktop.MouseEvent) {
}

// MouseOut is called when a desktop pointer exits the widget
func (cp *ColorPatch) MouseOut() {
	cp.hovered = false
	cp.applyPatchTheme()
}

// Mouseable
func (cp *ColorPatch) MouseUp(*desktop.MouseEvent) {
}
func (cp *ColorPatch) MouseDown(*desktop.MouseEvent) {
} // Mouseable

// Draggable
func (cp *ColorPatch) Dragged(d *fyne.DragEvent) {
	// pos := d.Position
}

func (cp *ColorPatch) DragEnd() {
	fmt.Println("DragEnd")
} // Draggable

// fyne.Focusable
func (cp *ColorPatch) TypedRune(rune) {
}

func (cp *ColorPatch) TypedKey(ev *fyne.KeyEvent) {
	switch ev.Name {
	case fyne.KeySpace:
		cp.Tapped(nil)
	}
}

func (cp *ColorPatch) FocusGained() {
	cp.focused = true
	cp.applyPatchTheme()
}

func (cp *ColorPatch) FocusLost() {
	cp.focused = false
	cp.applyPatchTheme()
} // fyne.Focusable

func (cp *ColorPatch) SetOnTapped(tapped func()) {
	cp.onTapped = tapped
}

func (cp *ColorPatch) SetUnused(b bool) {
	cp.unused = b
	cp.setFill(theme.DisabledColor())
	cp.setChanged()
}

func (cp *ColorPatch) Unused() bool {
	return cp.unused
}

func (cp *ColorPatch) GetColor() color.Color {
	return cp.FillColor
}

func (cp *ColorPatch) SetColor(c color.NRGBA) {
	cp.FillColor = c
	cp.setFill(c)
}

func (cp *ColorPatch) CopyPatch(source *ColorPatch) {
	cp.unused = source.unused
	cp.FillColor = source.FillColor
	if cp.unused {
		cp.SetUnused(true)
	} else {
		cp.setFill(cp.FillColor)
	}
}

func (cp *ColorPatch) setFill(color color.Color) {
	cp.rectangle.FillColor = color
	cp.rectangle.Refresh()
}

func (cp *ColorPatch) Tapped(_ *fyne.PointEvent) {
	if cp.onTapped != nil {
		cp.onTapped()
	}
}

func (cp *ColorPatch) EditCut() {
	cp.TypedShortcut(&fyne.ShortcutCut{Clipboard: Clipboard()})
}
func (cp *ColorPatch) EditCopy() {
	cp.TypedShortcut(&fyne.ShortcutCopy{Clipboard: Clipboard()})
}
func (cp *ColorPatch) EditPaste() {
	cp.TypedShortcut(&fyne.ShortcutPaste{Clipboard: Clipboard()})
}

func (cp *ColorPatch) TappedSecondary(pointEvent *fyne.PointEvent) {
	cp.requestFocus()
	cutItem := fyne.NewMenuItem("Cut", func() {
		cp.EditCut()
	})
	copyItem := fyne.NewMenuItem("Copy", func() {
		cp.EditCopy()
	})

	pasteItem := fyne.NewMenuItem("Paste", func() {
		cp.EditPaste()
	})

	menu := &fyne.Menu{}
	switch {
	case cp.Editing:
		menu.Items = []*fyne.MenuItem{cutItem, copyItem, pasteItem}
	default:
		menu.Items = []*fyne.MenuItem{cutItem, copyItem, pasteItem,
			fyne.NewMenuItem("Edit", func() {
				cp.Tapped(nil)
			})}
	}

	popUp := widget.NewPopUpMenu(menu, CanvasForObject(cp))
	var popUpPosition fyne.Position
	if pointEvent != nil {
		// popUpPosition = pointEvent.Position.AddXY(0, theme.Padding())
		popUpPosition = pointEvent.Position
	} else {
		popUpPosition = fyne.Position{X: cp.Size().Width / 2, Y: cp.Size().Height}
	}
	popUp.ShowAtRelativePosition(popUpPosition, cp)

}

type patchRenderer struct {
	objects []fyne.CanvasObject
	patch   *ColorPatch
}

func (cp *ColorPatch) requestFocus() {
	if c := fyne.CurrentApp().Driver().CanvasForObject(cp); c != nil {
		c.Focus(cp)
	}
}

func (cp *ColorPatch) CreateRenderer() fyne.WidgetRenderer {
	pr := &patchRenderer{
		objects: []fyne.CanvasObject{cp.rectangle},
		patch:   cp,
	}
	cp.ExtendBaseWidget(cp)
	return pr
}

func (pr *patchRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{pr.patch.background, pr.patch.rectangle}
}

func (pr *patchRenderer) Destroy() {}

func (pr *patchRenderer) Refresh() {
	pr.Layout(pr.patch.Size())
}

func (pr *patchRenderer) Layout(size fyne.Size) {
	pr.patch.background.Resize(size)
	if pr.patch.hovered || pr.patch.focused {
		diff := theme.Padding() * 2
		vec := fyne.Delta{DX: diff / 2, DY: diff / 2}
		rectPos := pr.patch.background.Position().Add(vec)
		pr.patch.rectangle.Move(rectPos)
		size = size.SubtractWidthHeight(diff, diff)
	} else {
		pr.patch.rectangle.Move(pr.patch.background.Position())
	}

	pr.patch.rectangle.Resize(size)
}

func (pr *patchRenderer) MinSize() fyne.Size {
	return fyne.NewSquareSize(theme.IconInlineSize() + theme.Padding())
}
