package main

import (
	"image"
	"time"

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
	AddSlideBtn  = 34
	MoveTool     = 35
	BackBtn      = 36
	CanvasWidget = 39

	ColorPickerBtn = 312
	MinusSizeBtn   = 314
	PlusSizeBtn    = 315
	DrawnSizeInput = 316
	FontNameInput  = 317

	PROJ_NameInput  = 51
	PROJ_NewProject = 52
	PROJ_OpenWDBtn  = 53

	FD_CloseBtn    = 61
	FD_FontRegular = 62
	FD_FontBold    = 63
	FD_FontRough   = 64
)

var (
	ProjObjCoords       map[int]g143.Rect
	NameInputEnteredTxt string

	ObjCoords          map[int]g143.Rect
	CurrentWindowFrame image.Image
	ProjectName        string
	CurrentSlide       int
	TotalSlides        int  = 1
	SlidesOffset       int  = 0
	InWorkView         bool = false

	CanvasRect g143.Rect
	activeTool int

	SlideFormat [][]Drawn

	InputsState     map[string]string
	SlidePreview    map[int][]*image.Image
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
	scrollEventsCount = 0

	toolNames = map[int]string{
		TextTool: "Text", ImageTool: "Image", MoveTool: "Move",
	}

	ToMoveIndex int = -1

	FDObjCoords map[int]g143.Rect
)

type DrawnType int

const (
	TextType DrawnType = iota
	ImageType
)

type Drawn struct {
	WidgetCode int
	Type       DrawnType
	X, Y       int
	W, H       int
	// text things
	Text     string
	Color    string
	FontName string
	// image things
	ImagePath string
	// both text and image
	Size int
}

type ToSortProject struct {
	Name    string
	ModTime time.Time
}

type FontDef struct {
	Class    string
	Index    int
	FontData []byte
}
