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

	// top panel
	tTRS := theCtx.drawButtonB(TextTool, 50, 10, "Text", fontColor, "#D9D5B0", "#AEAC9C")
	iTX := nextHorizontalX(tTRS, 20)
	iTRS := theCtx.drawButtonB(ImageTool, iTX, 10, "Image", fontColor, "#D9D5B0", "#D9D5B0")
	pTX := nextHorizontalX(iTRS, 20)
	theCtx.drawButtonB(PencilTool, pTX, 10, "Pencil", fontColor, "#D9D5B0", "#D9D5B0")

	activeTool = TextTool

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

	CanvasRect = g143.NewRect(workPanelX+1, 80+1, int(math.Ceil(WorkAreaWidth*0.8-2)),
		int(math.Ceil(WorkAreaHeight*0.8-2)))

	ObjCoords[CanvasWidget] = CanvasRect
	// send the frame to glfw window
	g143.DrawImage(wWidth, wHeight, theCtx.ggCtx.Image(), theCtx.windowRect())
	window.SwapBuffers()

	// save the frame
	CurrentWindowFrame = theCtx.ggCtx.Image()
}

func workViewMouseCallback(window *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if action != glfw.Release {
		return
	}

	xPos, yPos := window.GetCursorPos()
	xPosInt := int(xPos)
	yPosInt := int(yPos)

	wWidth, wHeight := window.GetSize()

	var widgetCode int

	for code, RS := range ObjCoords {
		if g143.InRect(RS, xPosInt, yPosInt) {
			// widgetRS = RS
			widgetCode = code
			break
		}
	}

	if widgetCode == 0 {
		return
	}

	// rootPath, _ := GetRootPath()

	switch widgetCode {
	case TextTool, ImageTool, PencilTool:
		activeTool = widgetCode

		theCtx := Continue2dCtx(CurrentWindowFrame, &ObjCoords)
		toolNames := map[int]string{
			TextTool: "Text", ImageTool: "Image", PencilTool: "Pencil",
		}
		// clear all tools
		for _, toolId := range []int{TextTool, ImageTool, PencilTool} {
			toolRS := ObjCoords[toolId]
			theCtx.drawButtonB(toolId, toolRS.OriginX, toolRS.OriginY, toolNames[toolId],
				fontColor, "#D9D5B0", "#D9D5B0")
		}
		// place indicator on activeTool
		activeToolRS := ObjCoords[activeTool]
		theCtx.drawButtonB(activeTool, activeToolRS.OriginX, activeToolRS.OriginY, toolNames[activeTool],
			fontColor, "#D9D5B0", "#AEAC9C")

		// send the frame to glfw window
		g143.DrawImage(wWidth, wHeight, theCtx.ggCtx.Image(), theCtx.windowRect())
		window.SwapBuffers()

		// save the frame
		CurrentWindowFrame = theCtx.ggCtx.Image()

	}

}
