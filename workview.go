package main

import (
	"fmt"
	"image"
	"math"
	"slices"
	"strconv"
	"strings"

	g143 "github.com/bankole7782/graphics143"
	"github.com/fogleman/gg"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/kovidgoyal/imaging"
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
	mSRS := theCtx.drawButtonB(MinusSizeTool, dividerX+20, 10+5, "--", "#fff", "#aaa")
	tSIX := nextHorizontalX(mSRS, 5)

	size := SlideMemory[CurrentSlide]["size"]
	tSIRS := theCtx.drawInput(DrawnSizeInput, tSIX, 10, size)
	pSX := nextHorizontalX(tSIRS, 5)
	pSRS := theCtx.drawButtonB(PlusSizeTool, pSX, 10+5, "+", "#fff", "#aaa")
	tCX := nextHorizontalX(pSRS, 20)
	selectedColor := SlideMemory[CurrentSlide]["color"]
	theCtx.drawColorBox(TextColorTool, tCX, 10+5, tTRS.Height-15, selectedColor)

	activeTool = TextTool

	// slides panel
	currentY := 80
	for i := range TotalSlides {
		iInUse := i + 1
		theCtx.ggCtx.SetHexColor(fontColor)
		theCtx.ggCtx.DrawString(fmt.Sprintf("%d", iInUse), 10, float64(currentY+FontSize))

		theCtx.ggCtx.SetHexColor(fontColor)
		theCtx.ggCtx.DrawRectangle(float64(10+FontSize+5), float64(currentY), WorkAreaWidth*0.15, WorkAreaHeight*0.15)
		theCtx.ggCtx.Fill()

		theCtx.ggCtx.SetHexColor("#fff")
		theCtx.ggCtx.DrawRectangle(float64(10+FontSize+5)+1, float64(currentY)+1, (WorkAreaWidth*0.15 - 2),
			(WorkAreaHeight*0.15 - 2))
		theCtx.ggCtx.Fill()

		currentY += 10 + int(math.Ceil(WorkAreaHeight*0.15))
	}

	// canvas
	canvasX := int(math.Ceil(WorkAreaWidth*0.15)) + 10 + FontSize + 5 + 20

	theCtx.ggCtx.SetHexColor(fontColor)
	theCtx.ggCtx.DrawRectangle(float64(canvasX), 80, WorkAreaWidth*0.8, WorkAreaHeight*0.8)
	theCtx.ggCtx.Fill()

	nWAW, nWAH := int(math.Ceil(WorkAreaWidth*0.8-2)), int(math.Ceil(WorkAreaHeight*0.8-2))
	canvasImg := drawOnCanvas()
	canvasImg = imaging.Fit(canvasImg, nWAW, nWAH, imaging.Lanczos)
	theCtx.ggCtx.DrawImage(canvasImg, canvasX+1, 80+1)
	CanvasRect = g143.NewRect(canvasX+1, 80+1, nWAW, nWAH)

	ObjCoords[CanvasWidget] = CanvasRect
	// send the frame to glfw window
	g143.DrawImage(wWidth, wHeight, theCtx.ggCtx.Image(), theCtx.windowRect())
	window.SwapBuffers()

	// save the frame
	CurrentWindowFrame = theCtx.ggCtx.Image()
}

func drawOnCanvas() image.Image {
	// frame buffer
	ggCtx := gg.NewContext(WorkAreaWidth, WorkAreaHeight)

	// background rectangle
	ggCtx.DrawRectangle(0, 0, float64(WorkAreaWidth), float64(WorkAreaHeight))
	ggCtx.SetHexColor("#ffffff")
	ggCtx.Fill()

	// load font
	fontPath := GetDefaultFontPath()
	ggCtx.LoadFontFace(fontPath, 20)

	for i, obj := range SlideFormat[CurrentSlide] {
		if obj.Type == TextType {
			strs := strings.Split(strings.ReplaceAll(obj.Text, "\r", ""), "\n")
			textFontSize := float64(obj.Size) * 15
			textFontSizeInt := int(math.Ceil(textFontSize))
			ggCtx.LoadFontFace(fontPath, textFontSize)

			maxX := 0
			currentY := obj.Y

			textHeight := 0
			for _, str := range strs {
				strW, strH := ggCtx.MeasureString(str)
				if int(strW) > maxX {
					maxX = int(strW)
				}
				textHeight += int(strH)
				ggCtx.SetHexColor(obj.Color)
				ggCtx.DrawString(str, float64(obj.X), float64(currentY+	textFontSizeInt))
				currentY += 10
			}

			obj.W = maxX
			obj.H = textHeight
			SlideFormat[CurrentSlide][i] = obj

		} else if obj.Type == ImageType {

		} else if obj.Type == PencilType {

		}
	}

	return ggCtx.Image()
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
		activeX, activeY = int(translastedMouseX), int(translatedMouseY)

		if activeTool == TextTool && ctrlState == glfw.Release {
			// stop interaction till returning from tpicker
			window.SetMouseButtonCallback(nil)
			window.SetCursorPosCallback(nil)

			foundIndex := -1
			for i, obj := range SlideFormat[CurrentSlide] {
				if obj.Type != TextType {
					continue
				}
				objRect := g143.NewRect(obj.X, obj.Y, obj.W, obj.H)
				if g143.InRect(objRect, activeX, activeY) {
					foundIndex = i
					break
				}
			}

			if foundIndex != -1 {
				DrawnEditIndex = foundIndex
				PickerChan <- []string{"text", SlideFormat[CurrentSlide][foundIndex].Text}
			} else {
				PickerChan <- []string{"text", ""}
			}

		} else if activeTool == TextTool && ctrlState == glfw.Press {
			foundIndex := -1
			for i, obj := range SlideFormat[CurrentSlide] {
				if obj.Type != TextType {
					continue
				}
				objRect := g143.NewRect(obj.X, obj.Y, obj.W, obj.H)
				if g143.InRect(objRect, activeX, activeY) {
					foundIndex = i
					break
				}
			}

			if foundIndex != -1 {
				objs := SlideFormat[CurrentSlide]
				SlideFormat[CurrentSlide] = slices.Delete(objs, foundIndex, foundIndex+1)
			}

			DrawWorkView(window, CurrentSlide)
		}

	case PlusSizeTool:
		size := SlideMemory[CurrentSlide]["size"]
		sizeInt, _ := strconv.Atoi(size)
		if sizeInt != 10 {
			sizeInt += 1
		}

		theCtx := Continue2dCtx(CurrentWindowFrame, &ObjCoords)
		widgetRS := ObjCoords[DrawnSizeInput]
		theCtx.drawInput(DrawnSizeInput, widgetRS.OriginX, widgetRS.OriginY, strconv.Itoa(sizeInt))

		SlideMemory[CurrentSlide]["size"] = strconv.Itoa(sizeInt)
		// send the frame to glfw window
		g143.DrawImage(wWidth, wHeight, theCtx.ggCtx.Image(), theCtx.windowRect())
		window.SwapBuffers()

		// save the frame
		CurrentWindowFrame = theCtx.ggCtx.Image()

	case MinusSizeTool:
		size := SlideMemory[CurrentSlide]["size"]
		sizeInt, _ := strconv.Atoi(size)
		if sizeInt != 1 {
			sizeInt -= 1
		}
		theCtx := Continue2dCtx(CurrentWindowFrame, &ObjCoords)
		widgetRS := ObjCoords[DrawnSizeInput]
		theCtx.drawInput(DrawnSizeInput, widgetRS.OriginX, widgetRS.OriginY, strconv.Itoa(sizeInt))

		SlideMemory[CurrentSlide]["size"] = strconv.Itoa(sizeInt)
		// send the frame to glfw window
		g143.DrawImage(wWidth, wHeight, theCtx.ggCtx.Image(), theCtx.windowRect())
		window.SwapBuffers()

		// save the frame
		CurrentWindowFrame = theCtx.ggCtx.Image()

	case TextColorTool:
		window.SetMouseButtonCallback(nil)
		window.SetCursorPosCallback(nil)
		PickerChan <- []string{"color", ""}
	}

}
