package ui

import (
	"fmt"
	"fynecv/appdata"
	"fynecv/vision"
	"image"
	"log"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type CameraList struct {
	data      *appdata.AppData
	view      *CameraView
	List      *widget.List
	Container *fyne.Container
	bound     binding.UntypedList

	tools *widget.Toolbar
	// dlgAdd    *dialog.CustomDialog
	// dlgRemove *dialog.CustomDialog
}

func NewCameraList(data *appdata.AppData, win fyne.Window, view *CameraView) *CameraList {

	cl := &CameraList{
		data:  data,
		bound: binding.NewUntypedList(),
		tools: widget.NewToolbar(),
		view:  view,
	}

	dlgAdd := cl.NewAddDialog(win)
	dlgRemove := dialog.NewCustomConfirm("Remove Camera", "Sure", "Oops", widget.NewLabel("Remove"), func(bool) {}, win)

	cl.tools.Append(
		widget.NewToolbarAction(theme.ContentAddIcon(), func() {
			dlgAdd.EntitySelect.Options = data.EntityList("number.")
			dlgAdd.Show()
		}))
	cl.tools.Append(
		widget.NewToolbarAction(theme.ContentRemoveIcon(), func() { dlgRemove.Show() }))

	const (
		thumbWidth  = 320
		thumbHeight = 200
	)

	for _, camera := range data.Cameras {
		camera.ThumbHook = NewCameraHook(nil)
		cl.bound.Append(camera)
	}

	cl.List = widget.NewListWithData(
		cl.bound,
		func() fyne.CanvasObject {
			img := image.NewNRGBA(image.Rect(0, 0, thumbWidth, thumbHeight))
			imageBox := canvas.NewImageFromImage(img)
			imageBox.FillMode = canvas.ImageFillOriginal
			imageBox.SetMinSize(fyne.NewSize(thumbWidth, thumbHeight))
			return imageBox
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			d, _ := i.(binding.Untyped).Get()
			camera, _ := d.(*vision.Camera)
			camera.ThumbHook.SetUi(o)
		},
	)

	cl.Container = container.NewBorder(cl.tools, nil, nil, nil, cl.List)

	return cl
}

func (cl *CameraList) Add(s string) {
	cam := vision.NewCamera(s, "", "")
	cl.data.Cameras = append(cl.data.Cameras, cam)
	cam.MainHook = cl.view.MainHook
	cam.ThumbHook = NewCameraHook(nil)
	cl.bound.Append(cam)
}

type CameraAddDialog struct {
	*dialog.FormDialog
	EntitySelect *widget.Select
}

func (cl *CameraList) NewAddDialog(win fyne.Window) *CameraAddDialog {
	cad := &CameraAddDialog{}
	urlValue := binding.NewString()
	urlEntry := widget.NewEntryWithData(urlValue)
	urlItem := widget.NewFormItem("Url", urlEntry)

	cad.EntitySelect = widget.NewSelect(cl.data.EntityList(), func(s string) {})
	selItem := widget.NewFormItem("Entities", cad.EntitySelect)

	moveGroup := widget.NewCheckGroup([]string{"Pan", "Tilt"}, func(s []string) {
		fmt.Print(s)
	})
	moveGroup.Horizontal = true
	moveItem := widget.NewFormItem("Movement", moveGroup)

	netGroup := widget.NewRadioGroup([]string{"Local", "Network"}, func(s string) {})
	netGroup.Horizontal = true
	netItem := widget.NewFormItem("Visibility", netGroup)

	// port := binding.NewSprintf("%d", binding.NewInt())
	// portEntry := widget.NewEntryWithData(port)
	// portItem := widget.NewFormItem("Port", portEntry)

	items := make([]*widget.FormItem, 0)
	items = append(items, urlItem, selItem, moveItem, netItem)

	dlg := dialog.NewForm("Add Camera", "Add", "Cancel", items,
		func(state bool) {
			if !state {
				return
			}

			s, _ := urlValue.Get()
			if len(s) == 0 {
				log.Println("add camera", "zero length")
				return
			}

			_, err := url.Parse(s)
			if err != nil {
				log.Println("AddCamera", err)
				return
			}

			log.Println("add camera to list")
			cl.Add(s)
			log.Println("add camera", state)
		}, win)

	dlg.Resize(fyne.NewSize(600, 0))
	cad.FormDialog = dlg
	return cad
}
