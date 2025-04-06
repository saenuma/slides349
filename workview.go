package main

import (
	"fmt"
	"math"

	g143 "github.com/bankole7782/graphics143"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func DrawWorkView(window *glfw.Window, slide int) {

	CurrentSlide = slide

	window.SetTitle(fmt.Sprintf("Project: %s ---- %s", ProjectName, ProgTitle))

	ObjCoords = make(map[int]g143.Rect)

	wWidth, wHeight := window.GetSize()
	theCtx := New2dCtx(wWidth, wHeight, &ObjCoords)

	currentY := 80
	for i := range TotalSlides {
		iInUse := i + 1
		theCtx.ggCtx.SetHexColor(fontColor)
		theCtx.ggCtx.DrawString(fmt.Sprintf("%d", iInUse), 10, float64(currentY+FontSize))

		theCtx.ggCtx.SetHexColor(fontColor)
		theCtx.ggCtx.DrawRoundedRectangle(float64(10+FontSize+5), float64(currentY), WorkAreaWidth*0.15, WorkAreaHeight*0.15, 10)
		theCtx.ggCtx.Fill()

		theCtx.ggCtx.SetHexColor("#fff")
		theCtx.ggCtx.DrawRoundedRectangle(float64(10+FontSize+5)+2, float64(currentY)+2, (WorkAreaWidth*0.15 - 4),
			(WorkAreaHeight*0.15 - 4), 10)
		theCtx.ggCtx.Fill()

		currentY += 10 + int(math.Ceil(WorkAreaHeight*0.15))
	}

	// send the frame to glfw window
	windowRS := g143.Rect{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
	g143.DrawImage(wWidth, wHeight, theCtx.ggCtx.Image(), windowRS)
	window.SwapBuffers()

	// save the frame
	CurrentWindowFrame = theCtx.ggCtx.Image()
}
