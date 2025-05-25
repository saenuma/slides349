package main

import (
	"image"

	g143 "github.com/bankole7782/graphics143"
)

const (
	FPS            = 24
	FontSize       = 20
	fontColor      = "#444"
	ProgTitle      = "a slides tool for videos349"
	WorkAreaWidth  = 1366
	WorkAreaHeight = 768

	TextTool     = 31
	ImageTool    = 32
	PencilTool   = 33
	SelectTool   = 34
	MoveTool     = 35
	CanvasWidget = 39

	TextColorTool  = 312
	MinusSizeTool  = 314
	PlusSizeTool   = 315
	DrawnSizeInput = 316
)

var (
	ObjCoords          map[int]g143.Rect
	CurrentWindowFrame image.Image
	ProjectName        string = "tmp_proj"
	CurrentSlide       int
	TotalSlides        int = 3
	CanvasRect         g143.Rect
	activeTool         int

	SlideFormat map[int][]Drawn

	SlideMemory     map[int]map[string]string
	SlidePreview    map[int]map[string]*image.Image
	CurrentSlideImg image.Image

	DrawnEditIndex int = -1

	activeX int
	activeY int

	PickerChan        = make(chan []string) // 0:instr, 1:data
	TextFromTPicker   string
	ClearAfterTPicker bool

	ClearAfterFPicker bool
	PathFromFPicker   string

	TextFromACPicker   string
	ClearAFterACPicker bool

	cursorEventsCount = 0

	toolNames = map[int]string{
		SelectTool: "Select", MoveTool: "Move",
		TextTool: "Text", ImageTool: "Image", PencilTool: "Pencil",
	}
)

type DrawnType int

const (
	TextType DrawnType = iota
	ImageType
	PencilType
)

type Drawn struct {
	WidgetCode int
	Type       DrawnType
	X, Y       int
	W, H       int
	// text things
	Text  string
	Color string
	Size  int
	// image things
	ImagePath string
}
