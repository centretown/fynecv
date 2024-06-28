package cv

import (
	"fmt"

	"gocv.io/x/gocv"
)

type WinHook struct {
	window   *gocv.Window
	trackers []*Tracker
}

var _ Hook = (*WinHook)(nil)

func NewWinHook(device any, trackers ...*Tracker) *WinHook {
	wh := &WinHook{}
	wh.window = gocv.NewWindow(fmt.Sprint("Camera:", device))
	wh.trackers = make([]*Tracker, len(trackers))
	for i, tr := range trackers {
		tr.trackBar = wh.window.CreateTrackbar(tr.Title, tr.Max)
		tr.trackBar.SetPos(tr.Pos)
		wh.trackers[i] = tr
	}
	return wh
}

func (wh *WinHook) Update(img *gocv.Mat) {
	wh.window.IMShow(*img)
	for _, tr := range wh.trackers {
		newPos := tr.trackBar.GetPos()
		if tr.Pos != newPos {
			tr.Pos = newPos
			tr.Handler(tr.Pos)
		}
	}
	wh.window.WaitKey(100)
}

func (wh *WinHook) Close(int) {
	wh.window.Close()
}
