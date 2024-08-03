package vision

import "image"

type Hook interface {
	UpdateHook(img image.Image)
	CloseHook(int)
}

type UiHook interface {
	Hook
	SetUi(ui interface{})
}
