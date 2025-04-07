package main

import (
	"fmt"
	"math"
	"strings"

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
	sTRS := theCtx.drawButtonA(SelectTool, 50, 10, "Select", fontColor, "#D9D5B0", "#D9D5B0")
	tTX := nextHorizontalX(sTRS, 20)
	tTRS := theCtx.drawButtonA(TextTool, tTX, 10, "Text", fontColor, "#D9D5B0", "#AEAC9C")
	iTX := nextHorizontalX(tTRS, 20)
	iTRS := theCtx.drawButtonA(ImageTool, iTX, 10, "Image", fontColor, "#D9D5B0", "#D9D5B0")
	pTX := nextHorizontalX(iTRS, 20)
	pTRS := theCtx.drawButtonA(PencilTool, pTX, 10, "Pencil", fontColor, "#D9D5B0", "#D9D5B0")
	mTX := nextHorizontalX(pTRS, 20)
	mTRS := theCtx.drawButtonA(MoveTool, mTX, 10, "Move", fontColor, "#D9D5B0", "#D9D5B0")
	dividerX := nextHorizontalX(mTRS, 20)
	theCtx.ggCtx.SetHexColor("#444")
	theCtx.ggCtx.DrawRectangle(float64(dividerX), 10, 2, float64(tTRS.Height))
	theCtx.ggCtx.Fill()

	// pencil extras
	mSRS := theCtx.drawButtonB(MinusSizeTool, dividerX+20, 10+5, "-", "#fff", "#aaa")
	pSX := nextHorizontalX(mSRS, 20)
	pSRS := theCtx.drawButtonB(PlusSizeTool, pSX, 10+5, "+", "#fff", "#aaa")
	tCX := nextHorizontalX(pSRS, 20)
	theCtx.drawColorBox(TextColorTool, tCX, 10+5, tTRS.Height-15, fontColor)

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

		theCtx.ggCtx.SetHexColor("#fff")
		theCtx.ggCtx.DrawRoundedRectangle(float64(10+FontSize+5)+1, float64(currentY)+1, (WorkAreaWidth*0.15 - 2),
			(WorkAreaHeight*0.15 - 2), 10)
		theCtx.ggCtx.Fill()

		currentY += 10 + int(math.Ceil(WorkAreaHeight*0.15))
	}

	// canvas
	canvasX := int(math.Ceil(WorkAreaWidth*0.15)) + 10 + FontSize + 5 + 20

	theCtx.ggCtx.SetHexColor(fontColor)
	theCtx.ggCtx.DrawRectangle(float64(canvasX), 80, WorkAreaWidth*0.8, WorkAreaHeight*0.8)
	theCtx.ggCtx.Fill()

	theCtx.ggCtx.SetHexColor("#fff")
	theCtx.ggCtx.DrawRectangle(float64(canvasX+1), 80+1, (WorkAreaWidth*0.8 - 2),
		(WorkAreaHeight*0.8 - 2))
	theCtx.ggCtx.Fill()

	CanvasRect = g143.NewRect(canvasX+1, 80+1, int(math.Ceil(WorkAreaWidth*0.8-2)),
		int(math.Ceil(WorkAreaHeight*0.8-2)))

	// draw contents of SlideFormat to canvas
	for _, obj := range SlideFormat[CurrentSlide] {
		if obj.Type == TextType {
			textDetail := TextDetails[obj.DetailsId]
			strs := strings.Split(strings.ReplaceAll(textDetail.Text, "\r", ""), "\n")
			for j, str := range strs {
				theCtx.ggCtx.SetHexColor(textDetail.Color)
				theCtx.ggCtx.DrawString(str, float64(obj.X+canvasX),
					float64(obj.Y+CanvasRect.OriginY+10+((j+1)*FontSize)))
				currentY += FontSize + 5
			}
		} else if obj.Type == ImageType {

		} else if obj.Type == PencilType {

		}
	}

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
	case SelectTool, TextTool, ImageTool, PencilTool, MoveTool:
		activeTool = widgetCode

		theCtx := Continue2dCtx(CurrentWindowFrame, &ObjCoords)
		toolNames := map[int]string{
			SelectTool: "Select", MoveTool: "Move",
			TextTool: "Text", ImageTool: "Image", PencilTool: "Pencil",
		}
		// clear all tools
		for _, toolId := range []int{SelectTool, TextTool, ImageTool, PencilTool, MoveTool} {
			toolRS := ObjCoords[toolId]
			theCtx.drawButtonA(toolId, toolRS.OriginX, toolRS.OriginY, toolNames[toolId],
				fontColor, "#D9D5B0", "#D9D5B0")
		}
		// place indicator on activeTool
		activeToolRS := ObjCoords[activeTool]
		theCtx.drawButtonA(activeTool, activeToolRS.OriginX, activeToolRS.OriginY, toolNames[activeTool],
			fontColor, "#D9D5B0", "#AEAC9C")

		// send the frame to glfw window
		g143.DrawImage(wWidth, wHeight, theCtx.ggCtx.Image(), theCtx.windowRect())
		window.SwapBuffers()

		// save the frame
		CurrentWindowFrame = theCtx.ggCtx.Image()

	case CanvasWidget:
		// theCtx := Continue2dCtx(CurrentWindowFrame, &ObjCoords)
		ctrlState := window.GetKey(glfw.KeyLeftControl)
		canvasRS := ObjCoords[CanvasWidget]

		translastedMouseX, translatedMouseY := xPos-float64(canvasRS.OriginX), yPos-float64(canvasRS.OriginY)

		if activeTool == TextTool && ctrlState == glfw.Release {

			activeX, activeY = int(translastedMouseX), int(translatedMouseY)
			window.SetMouseButtonCallback(nil)
			window.SetCursorPosCallback(nil)
			InTPickerChannel <- ""
		}

	}

}
