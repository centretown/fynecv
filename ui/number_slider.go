package ui

import (
	"fmt"
	"fynecv/appdata"
	"fynecv/svc"
	"log"

	"fyne.io/fyne/v2/widget"
)

func NewNumberSlider(entity string, data *appdata.AppData) *widget.Slider {
	var (
		number           appdata.Number
		value            float64
		slider           = widget.NewSlider(0, 180)
		fromSubscription bool
	)

	slider.OnChangeEnded = func(f float64) {
		if !fromSubscription {
			data.CallService(svc.NumberCmd(entity, svc.ServiceData{
				Key:   "value",
				Value: fmt.Sprintf("%.0f", f),
			}))
		}
		fromSubscription = false

	}

	data.Subscribe(entity,
		appdata.NewSubcription(&number.Entity, func(c appdata.Consumer) {
			var f float64
			_, err := fmt.Sscanf(number.State, "%f", &f)
			if err != nil {
				log.Println(err, entity)
				return
			}
			if f != value {
				value = f
				fromSubscription = true
				slider.SetValue(value)
			}
		}))

	return slider
}
