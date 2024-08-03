package ui

import (
	"fynecv/vision"
	"image"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

const (
	HighlightRadius = 120
)

var (
	HighlightDropped       = color.NRGBA{R: 192, G: 192, B: 192, A: 31}
	HighlightFill          = color.NRGBA{R: 192, G: 192, B: 192, A: 31}
	HighlightStroke        = color.NRGBA{R: 0, G: 255, B: 255, A: 127}
	HighlightStrokeDropped = color.NRGBA{R: 255, G: 0, B: 0, A: 127}
)

var _ fyne.Widget = (*CameraWidget)(nil)
var _ fyne.Scrollable = (*CameraWidget)(nil)

// var _ fyne.Tappable = (*CameraWidget)(nil)
// var _ fyne.SecondaryTappable = (*CameraWidget)(nil)
var _ desktop.Mouseable = (*CameraWidget)(nil)

var _ desktop.Hoverable = (*CameraWidget)(nil)

type CameraWidgetwState struct {
	zoom     float32
	viewPort image.Rectangle
	dropPos  fyne.Position
	dropped  bool
}

type CameraWidget struct {
	widget.BaseWidget
	canv *canvas.Image
	// canvasRect *canvas.Rectangle
	// zoom           float32
	// viewPort       image.Rectangle
	cursorPosition fyne.Position
	scrolled       bool
	circle         *canvas.Circle
	// dropped        bool
	state *CameraWidgetwState
}

func NewCameraWidget(img *canvas.Image) *CameraWidget {
	ci := &CameraWidget{
		canv:  img,
		state: &CameraWidgetwState{zoom: 100},
		// canvasRect: canvas.NewRectangle(color.NRGBA{R: 255, G: 0, B: 0, A: 255}),
		circle: canvas.NewCircle(HighlightFill),
	}
	ci.circle.StrokeColor = HighlightStroke
	ci.circle.StrokeWidth = 2
	// ci.circle.FillColor = HighlightFill
	ci.circle.Resize(fyne.NewSize(HighlightRadius, HighlightRadius))
	ci.ExtendBaseWidget(ci)
	return ci
}

func (ci *CameraWidget) Scrolled(e *fyne.ScrollEvent) {
	ci.cursorPosition = e.Position
	zoom := ci.state.zoom
	zoom += e.Scrolled.DY
	if zoom < 100 {
		zoom = 100
	} else if zoom > 1000 {
		zoom = 1000
	}
	ci.state.zoom = zoom
	ci.scrolled = true
}

var _ vision.UiHook = (*CameraWidget)(nil)

func (ci *CameraWidget) CloseHook(int) {}

func (ci *CameraWidget) UpdateHook(img image.Image) {
	if ci.canv != nil {
		if ci.state.zoom > 100.0 {
			ci.canv.Image = ci.Zoom(img)
		} else {
			ci.scrolled = false
			ci.fullView(img)
			ci.canv.Image = img
		}
		ci.canv.Refresh()
	}
}

func (ci *CameraWidget) SetUi(imageBox interface{}) {
	ci.canv = imageBox.(*canvas.Image)
	ci.state.zoom = 100.0
}

func (ci *CameraWidget) SetState(state *CameraWidgetwState) {
	ci.state.dropPos = ci.circle.Position()
	ci.state = state
	ci.circle.Move(state.dropPos)
	// ci.Refresh()
}

func (ci *CameraWidget) fullView(img image.Image) {
	ci.state.viewPort.Min.X = img.Bounds().Min.X
	ci.state.viewPort.Min.Y = img.Bounds().Min.Y
	ci.state.viewPort.Max.X = img.Bounds().Max.X
	ci.state.viewPort.Max.Y = img.Bounds().Max.Y
}

func (ci *CameraWidget) canvasToImage(pos fyne.Position) image.Point {
	viewPort := ci.state.viewPort
	imgX := float32(viewPort.Min.X) + pos.X*float32(viewPort.Dx())/ci.canv.Size().Width
	imgY := float32(viewPort.Min.Y) + pos.Y*float32(viewPort.Dy())/ci.canv.Size().Height
	return image.Pt(int(imgX), int(imgY))
}

func (ci *CameraWidget) imageToCanvas(pt image.Point) fyne.Position {
	viewPort := ci.state.viewPort
	posX := float32(pt.X - viewPort.Min.X)
	wr := ci.canv.Size().Width / float32(viewPort.Dx())
	posX *= wr

	posY := float32(pt.Y - viewPort.Min.Y)
	hr := ci.canv.Size().Height / float32(viewPort.Dy())
	posY *= hr
	return fyne.NewPos(posX, posY)
}

func (ci *CameraWidget) Zoom(img image.Image) image.Image {
	ycbr, ok := img.(*image.YCbCr)
	if !ok {
		return img
	}
	if !ci.scrolled {
		return ycbr.SubImage(ci.state.viewPort)
	}
	ci.scrolled = false

	if ci.state.viewPort.Empty() {
		ci.fullView(img)
	}

	imagePoint := ci.canvasToImage(ci.circle.Position())
	// fmt.Println("Min.X", ci.viewPort.Min.X, "Max.X", ci.viewPort.Max.X)
	// fmt.Println("Min.Y", ci.viewPort.Min.Y, "Max.Y", ci.viewPort.Max.Y)
	izoom := int(ci.state.zoom)
	zoomWidth := img.Bounds().Dx() * 100 / izoom
	zoomX := imagePoint.X - zoomWidth/2
	if zoomX < 0 {
		zoomX = 0
	}
	if zoomX+zoomWidth > img.Bounds().Max.X {
		zoomX = img.Bounds().Max.X - zoomWidth
	}

	zoomHeight := img.Bounds().Dy() * 100 / izoom
	zoomY := imagePoint.Y - zoomHeight/2
	if zoomY < 0 {
		zoomY = 0
	}
	if zoomY+zoomHeight > img.Bounds().Max.Y {
		zoomY = img.Bounds().Max.Y - zoomHeight
	}

	ci.state.viewPort = image.Rect(zoomX, zoomY, zoomX+zoomWidth, zoomY+zoomHeight)
	canvPoint := ci.imageToCanvas(imagePoint)
	// fmt.Println("circle position", ci.circle.Position(), "imagePoint", imagePoint, "canvPoint", canvPoint)
	ci.circle.Move(canvPoint)
	return ycbr.SubImage(ci.state.viewPort)
}

func (ci *CameraWidget) CreateRenderer() fyne.WidgetRenderer {
	cr := &cameraRenderer{
		objects:      []fyne.CanvasObject{ci.canv},
		cameraWidget: ci,
	}
	ci.ExtendBaseWidget(ci)
	return cr
}

type cameraRenderer struct {
	objects      []fyne.CanvasObject
	cameraWidget *CameraWidget
}

func (cr *cameraRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{cr.cameraWidget.canv, cr.cameraWidget.circle}
}

func (cr *cameraRenderer) Destroy() {}

func (cr *cameraRenderer) Refresh() {
	cr.Layout(cr.cameraWidget.canv.Size())
}

func (cr *cameraRenderer) Layout(size fyne.Size) {
	canv := cr.cameraWidget.canv
	canv.Resize(size)
	canv.Refresh()
	circle := cr.cameraWidget.circle
	if cr.cameraWidget.state.dropped {
		circle.FillColor = HighlightDropped
		circle.StrokeColor = HighlightStrokeDropped
	} else {
		circle.FillColor = HighlightFill
		circle.StrokeColor = HighlightStroke
	}
	cr.cameraWidget.circle.Refresh()
}

func (cr *cameraRenderer) MinSize() fyne.Size {
	return fyne.Size{Width: 300, Height: 300}
}

// func (ai *CameraImage) Tapped(_ *fyne.PointEvent) {
// 	fmt.Println("Hello Tapped")
// }

// func (ai *CameraImage) TappedSecondary(_ *fyne.PointEvent) {
// 	fmt.Println("Hello TappedSecondary")
// }

func (ci *CameraWidget) Cursor() desktop.Cursor {
	if !ci.state.dropped {
		return desktop.CrosshairCursor
	}
	return desktop.DefaultCursor
}

// Mouseable
func (ci *CameraWidget) MouseUp(event *desktop.MouseEvent) {
}
func (ci *CameraWidget) MouseDown(event *desktop.MouseEvent) {
	ci.state.dropped = !ci.state.dropped
	ci.Refresh()
} // Mouseable

func (ci *CameraWidget) MouseIn(event *desktop.MouseEvent) {
}

func (ci *CameraWidget) MouseMoved(event *desktop.MouseEvent) {
	if !ci.state.dropped {
		pos := fyne.Position{X: event.Position.X - HighlightRadius/2,
			Y: event.Position.Y - HighlightRadius/2}
		ci.circle.Move(pos)
	}
}

func (ci *CameraWidget) MouseOut() {
}
