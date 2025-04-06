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
	CanvasWidget = 39
)

var (
	ObjCoords          map[int]g143.Rect
	CurrentWindowFrame image.Image
	ProjectName        string = "tmp_proj"
	CurrentSlide       int
	TotalSlides        int = 3
	CanvasRect         g143.Rect
	SlideFormat        map[int][]Drawn
	activeTool         int
)

type Drawn struct {
	Type      string // one of text, image, pencil
	DetailsId int
}
