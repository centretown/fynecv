package ui

import "fyne.io/fyne/v2/data/binding"

func IncrementFloat(val float64, field binding.Float, bounds *FloatEntryBounds) func() {
	return func() {
		f, _ := field.Get()
		f += val
		if f >= bounds.MinVal && f <= bounds.MaxVal {
			field.Set(f)
		}
	}
}

func IncrementInt(val int, field binding.Int, bounds *IntEntryBounds) func() {
	return func() {
		f, _ := field.Get()
		f += val
		if f >= bounds.MinVal && f <= bounds.MaxVal {
			field.Set(f)
		}
	}
}
