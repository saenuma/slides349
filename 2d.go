package main

import (
	"image"

	g143 "github.com/bankole7782/graphics143"
	"github.com/fogleman/gg"
)

type Ctx struct {
	WindowWidth  int
	WindowHeight int
	ggCtx        *gg.Context
	ObjCoords    *map[int]g143.Rect
}

func New2dCtx(wWidth, wHeight int, objCoords *map[int]g143.Rect) Ctx {
	// frame buffer
	ggCtx := gg.NewContext(wWidth, wHeight)

	// background rectangle
	ggCtx.DrawRectangle(0, 0, float64(wWidth), float64(wHeight))
	ggCtx.SetHexColor("#ffffff")
	ggCtx.Fill()

	// load font
	fontPath := GetDefaultFontPath()
	err := ggCtx.LoadFontFace(fontPath, FontSize)
	if err != nil {
		panic(err)
	}

	ctx := Ctx{WindowWidth: wWidth, WindowHeight: wHeight, ggCtx: ggCtx,
		ObjCoords: objCoords}
	return ctx
}

func Continue2dCtx(img image.Image, objCoords *map[int]g143.Rect) Ctx {
	ggCtx := gg.NewContextForImage(img)

	// load font
	fontPath := GetDefaultFontPath()
	err := ggCtx.LoadFontFace(fontPath, FontSize)
	if err != nil {
		panic(err)
	}

	ctx := Ctx{WindowWidth: img.Bounds().Dx(), WindowHeight: img.Bounds().Dy(), ggCtx: ggCtx,
		ObjCoords: objCoords}
	return ctx
}

func (ctx *Ctx) drawButtonA(btnId, originX, originY int, text, textColor, bgColor, circleColor string) g143.Rect {
	// draw bounding rect
	textW, textH := ctx.ggCtx.MeasureString(text)
	width, height := textW+80, textH+30
	ctx.ggCtx.SetHexColor(bgColor)
	ctx.ggCtx.DrawRectangle(float64(originX), float64(originY), float64(width), float64(height))
	ctx.ggCtx.Fill()

	// draw text
	ctx.ggCtx.SetHexColor(textColor)
	ctx.ggCtx.DrawString(text, float64(originX)+20, float64(originY)+FontSize+10)

	// draw circle
	ctx.ggCtx.SetHexColor(circleColor)
	ctx.ggCtx.DrawCircle(float64(originX)+width-30, float64(originY)+(height/2), 10)
	ctx.ggCtx.Fill()

	// save dimensions
	btnARect := g143.NewRect(originX, originY, int(width), int(height))
	(*ctx.ObjCoords)[btnId] = btnARect
	return btnARect
}

func (ctx *Ctx) drawButtonB(btnId, originX, originY int, text, textColor, bgColor string) g143.Rect {
	// draw bounding rect
	textW, textH := ctx.ggCtx.MeasureString(text)
	width, height := textW+20, textH+15
	ctx.ggCtx.SetHexColor(bgColor)
	ctx.ggCtx.DrawRectangle(float64(originX), float64(originY), float64(width), float64(height))
	ctx.ggCtx.Fill()

	// draw text
	ctx.ggCtx.SetHexColor(textColor)
	ctx.ggCtx.DrawString(text, float64(originX)+10, float64(originY)+FontSize)

	// save dimensions
	btnARect := g143.NewRect(originX, originY, int(width), int(height))
	if btnId != 0 {
		(*ctx.ObjCoords)[btnId] = btnARect
	}
	return btnARect
}

func (ctx *Ctx) drawColorBox(inputId, originX, originY, width int, pickedColor string) g143.Rect {
	ctx.ggCtx.SetHexColor("#000")
	ctx.ggCtx.DrawRectangle(float64(originX), float64(originY), float64(width), float64(width))
	ctx.ggCtx.Fill()

	ctx.ggCtx.SetHexColor(pickedColor)
	ctx.ggCtx.DrawRectangle(float64(originX)+2, float64(originY)+2, float64(width)-4, float64(width)-4)
	ctx.ggCtx.Fill()

	colorBoxRect := g143.NewRect(originX, originY, width, width)
	(*ctx.ObjCoords)[inputId] = colorBoxRect

	return colorBoxRect
}

func (ctx *Ctx) drawInput(inputId, originX, originY int, writtenStr string) g143.Rect {

	ctx.ggCtx.SetHexColor("#fff")
	ctx.ggCtx.DrawRectangle(float64(originX), float64(originY), 40, float64(originY)+FontSize+5)
	ctx.ggCtx.Fill()

	ctx.ggCtx.SetHexColor("#909BD0")
	ctx.ggCtx.DrawRectangle(float64(originX), 50, 40, 3)
	ctx.ggCtx.Fill()

	entryRect := g143.Rect{Width: 40, Height: 50, OriginX: originX, OriginY: 10}
	(*ctx.ObjCoords)[inputId] = entryRect

	ctx.ggCtx.SetHexColor("#444")
	ctx.ggCtx.DrawString(writtenStr, float64(originX+5), float64(originY)+FontSize+5)
	return entryRect
}

func (ctx *Ctx) drawInputB(inputId, originX, originY, inputWidth int, placeholder string, isDefault bool) g143.Rect {
	height := 30
	ctx.ggCtx.SetHexColor(fontColor)
	ctx.ggCtx.DrawRectangle(float64(originX), float64(originY), float64(inputWidth), float64(height))
	ctx.ggCtx.Fill()

	ctx.ggCtx.SetHexColor("#fff")
	ctx.ggCtx.DrawRectangle(float64(originX)+2, float64(originY)+2, float64(inputWidth)-4, float64(height)-4)
	ctx.ggCtx.Fill()

	entryRect := g143.Rect{Width: inputWidth, Height: height, OriginX: originX, OriginY: originY}
	(*ctx.ObjCoords)[inputId] = entryRect

	if isDefault {
		ctx.ggCtx.SetHexColor("#444")
		ctx.ggCtx.DrawString(placeholder, float64(originX+15), float64(originY)+FontSize)
	} else {
		ctx.ggCtx.SetHexColor("#aaa")
		ctx.ggCtx.DrawString(placeholder, float64(originX+15), float64(originY)+FontSize)
	}
	return entryRect
}

func (ctx *Ctx) windowRect() g143.Rect {
	return g143.NewRect(0, 0, ctx.WindowWidth, ctx.WindowHeight)
}

func nextHorizontalX(aRect g143.Rect, margin int) int {
	return aRect.OriginX + aRect.Width + margin
}

func nextVerticalY(aRect g143.Rect, margin int) int {
	return aRect.OriginY + aRect.Height + margin
}
