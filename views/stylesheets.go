package views

import (
	"fmt"
	"path/filepath"
	"quickRadio/models"
	"runtime"
)

func GetPng(imageTitle string) string {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filepath.Dir(filename))
	imagePath := filepath.Join(dir, "assets", "pngs", imageTitle+".png")
	return imagePath
}

func CreateTeamButtonStylesheet(sweater models.Sweater) string {
	styleSheetTemplate := `
			background: qlineargradient(x1:0 y1:0, x2:1, y2:1, stop:0 %s, stop:.40 %s, stop:.60 %s, stop:1.0 %s);
			min-height: 100px;
			min-width: 100px;
	`
	styleSheet := fmt.Sprintf(styleSheetTemplate, sweater.PrimaryColor, sweater.PrimaryColor, sweater.SecondaryColor, sweater.SecondaryColor)
	return styleSheet
}
func CreateGameStylesheet(homeTeam string, awayTeam string) string {
	imagePath := GetPng("puck_texture")
	styleSheetTemplate := `
		background-image: url(%s);
		background-repeat: no-repeat;
		background-position: center;
	`
	styleSheet := fmt.Sprintf(styleSheetTemplate, imagePath)
	return styleSheet
}
func CreateDropdownStyleSheet() string {
	styleSheet := `
		background-color:rgb(168, 168, 168)
	`
	return styleSheet
}
func CreateGameManagerStyleSheet() string {
	styleSheet := `
		background-color :rgb(29, 29, 29)
	`
	return styleSheet
}
func CreateLabelStylesheet() string {
	styleSheet := `
		color :rgb(190, 190, 190)
	`
	return styleSheet
}
