package vision

import "image"

type Hook interface {
	Update(img image.Image)
	Close(int)
}

type UiHook interface {
	Hook
	SetUi(ui interface{})
}
