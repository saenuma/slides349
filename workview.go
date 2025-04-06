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

	// slides panel
	currentY := 80
	for i := range TotalSlides {
		iInUse := i + 1
		theCtx.ggCtx.SetHexColor(fontColor)
		theCtx.ggCtx.DrawString(fmt.Sprintf("%d", iInUse), 10, float64(currentY+FontSize))

		theCtx.ggCtx.SetHexColor(fontColor)
		theCtx.ggCtx.DrawRoundedRectangle(float64(10+FontSize+5), float64(currentY), WorkAreaWidth*0.15, WorkAreaHeight*0.15, 10)
		theCtx.ggCtx.Fill()

		if i == 0 {
		}
		theCtx.ggCtx.SetHexColor("#fff")
		theCtx.ggCtx.DrawRoundedRectangle(float64(10+FontSize+5)+1, float64(currentY)+1, (WorkAreaWidth*0.15 - 2),
			(WorkAreaHeight*0.15 - 2), 10)
		theCtx.ggCtx.Fill()

		currentY += 10 + int(math.Ceil(WorkAreaHeight*0.15))
	}

	// work panel
	workPanelX := int(math.Ceil(WorkAreaWidth*0.15)) + 10 + FontSize + 5 + 20

	theCtx.ggCtx.SetHexColor(fontColor)
	theCtx.ggCtx.DrawRectangle(float64(workPanelX), 80, WorkAreaWidth*0.8, WorkAreaHeight*0.8)
	theCtx.ggCtx.Fill()

	theCtx.ggCtx.SetHexColor("#fff")
	theCtx.ggCtx.DrawRectangle(float64(workPanelX+1), 80+1, (WorkAreaWidth*0.8 - 2),
		(WorkAreaHeight*0.8 - 2))
	theCtx.ggCtx.Fill()

	// send the frame to glfw window
	windowRS := g143.Rect{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
	g143.DrawImage(wWidth, wHeight, theCtx.ggCtx.Image(), windowRS)
	window.SwapBuffers()

	// save the frame
	CurrentWindowFrame = theCtx.ggCtx.Image()
}
