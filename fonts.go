package main

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
)

//go:embed Roboto-Light.ttf
var Roboto []byte

//go:embed Chewy-Regular.ttf
var Chewy []byte

//go:embed RubikMarkerHatch-Regular.ttf
var RubikHatch []byte

var fonts []FontDef = []FontDef{
	{Class: "regular", Index: 1, FontData: Roboto},
	{Class: "bold", Index: 2, FontData: Chewy},
	{Class: "rough", Index: 3, FontData: RubikHatch},
}

func getFontPath(index int) string {
	var fontData []byte
	for _, fontDef := range fonts {
		if index == fontDef.Index {
			fontData = fontDef.FontData
			break
		}
	}

	fontPath := filepath.Join(os.TempDir(), fmt.Sprintf("v349font_%d", index))
	if !DoesPathExists(fontPath) {
		os.WriteFile(fontPath, fontData, 0777)
	}
	return fontPath
}

func getFontDef(class string) FontDef {
	for _, fontDef := range fonts {
		if class == fontDef.Class {
			return fontDef
		}
	}

	return FontDef{}
}
