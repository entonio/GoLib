package xlsconstants

const DefaultHeightPoints = 15

type Color string

// Sort by red / green / blue
const (
	Gray      Color = "888888"
	LightGray Color = "F8F8F8"
	Red       Color = "FF0000"
	Magenta   Color = "FF00FF"
	Blue      Color = "0000FF"
	LightBlue Color = "3366FF"
	Orange    Color = "FF7F00"
	Green     Color = "008000"
	White     Color = "FFFFFF"
)

const (
	B2    Color = "eeece0"
	T2L40 Color = "538dd5"
)

type Pattern string

const (
	Solid Pattern = "solid"
)

type BorderStyle string

const (
	Medium BorderStyle = "medium"
	Thin   BorderStyle = "thin"
)

type Alignment string

const (
	Top Alignment = "top"
	// The 'middle' value does not exist in Excel, it's 'center'.
	// It's provided here only so that other APIs can distinguish V from H.
	Middle Alignment = "middle"
	Bottom Alignment = "bottom"

	Center Alignment = "center"
	Left   Alignment = "left"
	Right  Alignment = "right"
)

type Font string

const (
	Calibri Font = "Calibri"
)
