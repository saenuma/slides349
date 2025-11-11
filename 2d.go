package main

import (
	"image"

	g143 "github.com/bankole7782/graphics143"
	"github.com/fogleman/gg"
)

type Ctx struct {
	WindowWidth     int
	WindowHeight    int
	ggCtx           *gg.Context
	ObjCoords       *map[int]g143.Rect
	CurrentFontSize int
}

func New2dCtx(wWidth, wHeight int, objCoords *map[int]g143.Rect) Ctx {
	// frame buffer
	ggCtx := gg.NewContext(wWidth, wHeight)

	// background rectangle
	ggCtx.DrawRectangle(0, 0, float64(wWidth), float64(wHeight))
	ggCtx.SetHexColor("#ffffff")
	ggCtx.Fill()

	// load font
	fontPath := getFontPath(1)
	err := ggCtx.LoadFontFace(fontPath, 20)
	if err != nil {
		panic(err)
	}

	ctx := Ctx{WindowWidth: wWidth, WindowHeight: wHeight, ggCtx: ggCtx,
		ObjCoords: objCoords, CurrentFontSize: 20}
	return ctx
}

func Continue2dCtx(img image.Image, objCoords *map[int]g143.Rect) Ctx {
	ggCtx := gg.NewContextForImage(img)

	// load font
	fontPath := getFontPath(1)
	err := ggCtx.LoadFontFace(fontPath, 20)
	if err != nil {
		panic(err)
	}

	ctx := Ctx{WindowWidth: img.Bounds().Dx(), WindowHeight: img.Bounds().Dy(), ggCtx: ggCtx,
		ObjCoords: objCoords, CurrentFontSize: 20}
	return ctx
}

func (ctx *Ctx) setFontSize(fontSize int) {
	// load font
	fontPath := getFontPath(1)
	err := ctx.ggCtx.LoadFontFace(fontPath, float64(fontSize))
	if err != nil {
		panic(err)
	}

	ctx.CurrentFontSize = fontSize
}

func (ctx *Ctx) drawButtonB(btnId, originX, originY int, text, textColor, bgColor string) g143.Rect {
	// draw bounding rect
	textW, textH := ctx.ggCtx.MeasureString(text)
	width, height := textW+float64(ctx.CurrentFontSize), textH+float64(ctx.CurrentFontSize)
	ctx.ggCtx.SetHexColor(bgColor)
	ctx.ggCtx.DrawRectangle(float64(originX), float64(originY), float64(width), float64(height))
	ctx.ggCtx.Fill()

	textOffsetY := float64(ctx.CurrentFontSize) / 5
	textOffsetX := float64(ctx.CurrentFontSize) / 2
	// draw text
	ctx.ggCtx.SetHexColor(textColor)
	ctx.ggCtx.DrawString(text, float64(originX+int(textOffsetX)), float64(originY+ctx.CurrentFontSize)+textOffsetY)

	// save dimensions
	btnARect := g143.NewRect(originX, originY, int(width), int(height))
	(*ctx.ObjCoords)[btnId] = btnARect
	return btnARect
}

func (ctx *Ctx) drawButtonA(btnId, originX, originY int, text, textColor, bgColor string, active bool) g143.Rect {
	// draw bounding rect
	textW, textH := ctx.ggCtx.MeasureString(text)
	width, height := textW+float64(ctx.CurrentFontSize), textH+float64(ctx.CurrentFontSize)
	ctx.ggCtx.SetHexColor(bgColor)
	ctx.ggCtx.DrawRectangle(float64(originX), float64(originY), float64(width), float64(height))
	ctx.ggCtx.Fill()

	textOffsetY := float64(ctx.CurrentFontSize) / 5
	textOffsetX := float64(ctx.CurrentFontSize) / 2
	// draw text
	ctx.ggCtx.SetHexColor(textColor)
	ctx.ggCtx.DrawString(text, float64(originX+int(textOffsetX)), float64(originY+ctx.CurrentFontSize)+textOffsetY)

	if active {
		ctx.ggCtx.SetHexColor(bgColor)
	} else {
		ctx.ggCtx.SetHexColor("#fff")
	}

	// draw rect
	ctx.ggCtx.DrawRectangle(float64(originX), float64(originY)+height+5, width, 3)
	ctx.ggCtx.Fill()

	// save dimensions
	btnARect := g143.NewRect(originX, originY, int(width), int(height))
	(*ctx.ObjCoords)[btnId] = btnARect
	return btnARect
}

func (ctx *Ctx) drawColorBox(inputId, originX, originY, width int, pickedColor string) g143.Rect {
	ctx.ggCtx.SetHexColor(pickedColor)
	ctx.ggCtx.DrawRectangle(float64(originX), float64(originY), float64(width), float64(width))
	ctx.ggCtx.Fill()

	colorBoxRect := g143.NewRect(originX, originY, width, width)
	(*ctx.ObjCoords)[inputId] = colorBoxRect

	return colorBoxRect
}

func (ctx *Ctx) drawInput(inputId, originX, originY int, writtenStr string) g143.Rect {
	ctx.ggCtx.SetHexColor("#fff")
	ctx.ggCtx.DrawRectangle(float64(originX), float64(originY), 40, float64(originY)+FontSize+5)
	ctx.ggCtx.Fill()

	wSW, _ := ctx.ggCtx.MeasureString(writtenStr)
	ctx.ggCtx.SetHexColor("#909BD0")
	ctx.ggCtx.DrawRectangle(float64(originX), 50, wSW+10, 3)
	ctx.ggCtx.Fill()

	entryRect := g143.Rect{Width: int(wSW) + 10, Height: 50, OriginX: originX, OriginY: 10}
	(*ctx.ObjCoords)[inputId] = entryRect

	ctx.ggCtx.SetHexColor("#444")
	ctx.ggCtx.DrawString(writtenStr, float64(originX+5), float64(originY)+FontSize+5)
	return entryRect
}
func (ctx *Ctx) drawInputB(inputId, originX, originY, inputWidth int, placeholder string, isDefault bool) g143.Rect {
	height := float64(ctx.CurrentFontSize) * 1.8
	ctx.ggCtx.SetHexColor(fontColor)
	ctx.ggCtx.DrawRectangle(float64(originX), float64(originY), float64(inputWidth), float64(height))
	ctx.ggCtx.Fill()

	ctx.ggCtx.SetHexColor("#fff")
	ctx.ggCtx.DrawRectangle(float64(originX)+2, float64(originY)+2, float64(inputWidth)-4, float64(height)-4)
	ctx.ggCtx.Fill()

	entryRect := g143.Rect{Width: inputWidth, Height: int(height), OriginX: originX, OriginY: originY}
	(*ctx.ObjCoords)[inputId] = entryRect

	if isDefault {
		ctx.ggCtx.SetHexColor("#444")
		ctx.ggCtx.DrawString(placeholder, float64(originX+ctx.CurrentFontSize/3),
			float64(originY+ctx.CurrentFontSize/5)+float64(ctx.CurrentFontSize))
	} else {
		ctx.ggCtx.SetHexColor("#aaa")
		ctx.ggCtx.DrawString(placeholder, float64(originX+ctx.CurrentFontSize/3),
			float64(originY+ctx.CurrentFontSize/5)+float64(ctx.CurrentFontSize))
	}
	return entryRect
}

func (ctx *Ctx) drawAFont(inputId, originX, originY int, fontClass string) g143.Rect {
	sampleText1 := "The quick brown fox jumps over the lazy dog"
	sampleText2 := ".,?!12345678"

	ctx.ggCtx.SetHexColor(fontColor)
	fontDef := getFontDef(fontClass)
	// load font
	fontPath := getFontPath(1)
	ctx.ggCtx.LoadFontFace(fontPath, FontSize)
	ctx.ggCtx.DrawString(fontClass, float64(originX), float64(originY)+FontSize)
	previewFontPath := getFontPath(fontDef.Index)
	previewFontSize := 30.0
	ctx.ggCtx.LoadFontFace(previewFontPath, previewFontSize)
	ctx.ggCtx.DrawString(sampleText1, float64(originX)+10, float64(originY)+FontSize+previewFontSize)
	maxH := FontSize + previewFontSize + previewFontSize
	ctx.ggCtx.DrawString(sampleText2, float64(originX)+10, float64(originY)+maxH)
	maxW, _ := ctx.ggCtx.MeasureString(sampleText1)
	fontRect := g143.NewRect(originX, originY, int(maxW)+10, int(maxH)+10)
	(*ctx.ObjCoords)[inputId] = fontRect

	return fontRect
}

func (ctx *Ctx) windowRect() g143.Rect {
	return g143.NewRect(0, 0, ctx.WindowWidth, ctx.WindowHeight)
}

func nextX(aRect g143.Rect, margin int) int {
	return aRect.OriginX + aRect.Width + margin
}

func nextY(aRect g143.Rect, margin int) int {
	return aRect.OriginY + aRect.Height + margin
}
