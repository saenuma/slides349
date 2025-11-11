package main

import (
	"image"

	g143 "github.com/bankole7782/graphics143"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/kovidgoyal/imaging"
)

func drawFontPicker(window *glfw.Window, currentFrame image.Image) {
	FDObjCoords = make(map[int]g143.Rect)
	InWorkView = false

	wWidth, wHeight := window.GetSize()
	// background image
	img := imaging.AdjustBrightness(currentFrame, -40)
	theCtx := Continue2dCtx(img, &FDObjCoords)

	// dialog rectangle
	dialogWidth := 950
	dialogHeight := 600

	dialogOriginX := (wWidth - dialogWidth) / 2
	dialogOriginY := (wHeight - dialogHeight) / 2

	theCtx.ggCtx.SetHexColor("#fff")
	theCtx.ggCtx.DrawRectangle(float64(dialogOriginX), float64(dialogOriginY), float64(dialogWidth),
		float64(dialogHeight))
	theCtx.ggCtx.Fill()

	// Add Form
	aFLX, aFLY := dialogOriginX+20, dialogOriginY+20
	theCtx.ggCtx.SetHexColor("#444")
	theCtx.ggCtx.DrawString("Select Font", float64(aFLX), float64(aFLY)+FontSize)

	closeBtnOriginX := dialogWidth + dialogOriginX - 80
	closeBtnStr := "close"
	theCtx.drawButtonB(FD_CloseBtn, closeBtnOriginX, dialogOriginY+20, closeBtnStr, "#fff", "#B75F5F")

	// Font Previews
	font1RS := theCtx.drawAFont(FD_FontRegular, aFLX, aFLY+40, "regular")
	f2Y := nextY(font1RS, 20)
	font2RS := theCtx.drawAFont(FD_FontBold, aFLX, f2Y, "bold")
	f3Y := nextY(font2RS, 20)
	theCtx.drawAFont(FD_FontRough, aFLX, f3Y, "rough")

	// send the frame to glfw window
	g143.DrawImage(wWidth, wHeight, theCtx.ggCtx.Image(), theCtx.windowRect())
	window.SwapBuffers()

	// save the frame
	CurrentWindowFrame = theCtx.ggCtx.Image()
}

func fontPickerMouseCallback(window *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if action != glfw.Release {
		return
	}

	xPos, yPos := window.GetCursorPos()
	xPosInt := int(xPos)
	yPosInt := int(yPos)

	// wWidth, wHeight := window.GetSize()

	var widgetCode int

	for code, RS := range FDObjCoords {
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
	case FD_CloseBtn:
		DrawWorkView(window, CurrentSlide)
		window.SetMouseButtonCallback(workViewMouseCallback)
		window.SetKeyCallback(nil)
		window.SetScrollCallback(workViewScrollCB)
		// quick hover effect
		window.SetCursorPosCallback(getHoverCB(&ObjCoords))

	case FD_FontRegular, FD_FontBold, FD_FontRough:
		names := map[int]string{FD_FontRegular: "regular", FD_FontBold: "bold", FD_FontRough: "rough"}
		InputsState["font"] = names[widgetCode]

		DrawWorkView(window, CurrentSlide)
		window.SetMouseButtonCallback(workViewMouseCallback)
		window.SetKeyCallback(nil)
		window.SetScrollCallback(workViewScrollCB)
		// quick hover effect
		window.SetCursorPosCallback(getHoverCB(&ObjCoords))
	}
}
