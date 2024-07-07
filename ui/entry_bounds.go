package ui

type FloatEntryBounds struct {
	MinVal, MaxVal, OnVal, OffVal float64
}

func NewFloatEntryBounds(min, max, on, off float64) *FloatEntryBounds {
	if min > max {
		min, max = max, min
	}
	eb := &FloatEntryBounds{min, max, on, off}
	return eb
}

type IntEntryBounds struct {
	MinVal, MaxVal, OnVal, OffVal int
}

func NewIntEntryBounds(min, max, on, off int) *IntEntryBounds {
	if min > max {
		min, max = max, min
	}
	return &IntEntryBounds{min, max, on, off}
}

var (
	RateBounds       = &IntEntryBounds{MinVal: 16, MaxVal: 360, OnVal: 48, OffVal: 0}
	HueShiftBounds   = &IntEntryBounds{MinVal: -10, MaxVal: 10, OnVal: 1, OffVal: 0}
	ScanBounds       = &IntEntryBounds{MinVal: 1, MaxVal: 10, OnVal: 1, OffVal: 0}
	HueBounds        = &FloatEntryBounds{MinVal: 0, MaxVal: 360, OnVal: 180, OffVal: 0}
	SaturationBounds = &FloatEntryBounds{MinVal: 0, MaxVal: 100, OnVal: 50, OffVal: 0}
	ValueBounds      = &FloatEntryBounds{MinVal: 0, MaxVal: 100, OnVal: 50, OffVal: 0}
)
