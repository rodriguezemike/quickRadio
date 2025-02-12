package views

import (
	"fmt"
	"path"
	"quickRadio/models"
	"quickRadio/quickio"
)

func GetPng(imageTitle string) string {
	dir := quickio.GetProjectDirWithForwadSlash()
	imagePath := path.Join(dir, "assets", "pngs", imageTitle+".png")
	return imagePath
}

//Consider creating a Single Large stylesheet functionally, by combining these func ina func
//To generate a large stylesheet, but one that can broken up like this in funcs for ease of editing.

func CreateInactiveRadioStreamButtonStylesheet(sweater models.Sweater) string {
	stylesheet := `
		QPushButton[button-type="radio"]{
			background: qlineargradient(x1:0 y1:0, x2:1, y2:1, stop:0 %s, stop:.40 %s, stop:.60 %s, stop:1.0 %s);
			border-radius: 9px; /* Rounded borders */
			border: 2px solid %s; /* Subtle border color to match the button */
			box-shadow: 0 4px 8px %s; /* Soft shadow for elegance rgba(0, 0, 0, 0.2)* #000000C0 */
			min-height: 150px;
			min-width: 150px;
			max-height: 150px;
			max-height: 150px;
		}
	`
	styleSheet := fmt.Sprintf(stylesheet, sweater.PrimaryColor, sweater.PrimaryColor, sweater.SecondaryColor, sweater.SecondaryColor, sweater.PrimaryColor, sweater.PrimaryColor)
	return styleSheet
}

func CreateActiveRadioStreamButtonStylesheet(sweater models.Sweater) string {
	stylesheet := `
		QPushButton[button-type="radio"]{
			background-color:%s;
			border:none;
			border-radius: 9px; /* Rounded borders */
    		border: 2px solid %s; /* Subtle border color to match the button */
    		box-shadow: 0 4px 8px %s; /* Soft shadow for elegance rgba(0, 0, 0, 0.2) #000000FF*/
			min-height: 150px;
			min-width: 150px;
			max-height: 150px;
			max-height: 150px;
		}
	`
	styleSheet := fmt.Sprintf(stylesheet, sweater.PrimaryColor, sweater.SecondaryColor, sweater.PrimaryColor)
	return styleSheet
}

func CreateDynamicDataLabelStylesheet() string {
	stylesheet := `
		QLabel[label-type="dynamic"] {
			font-family: "Segoe UI", "Georgia", "Arial", sans-serif; /* Elegant font with a mix of modern and classic */
			font-size: 32px; /* Size large enough for readability, but elegant */
			color: #839496; /* Dark grey color for a soft yet sophisticated text color */
			background-color: transparent; /* No background, letting the label sit naturally */
			padding: 5px 15px; /* Padding around the text to give it some space */
			border-radius: 10px; /* Slightly rounded corners for the label */
			text-align: center; /* Centered text */
			letter-spacing: 0.5px; /* Slight spacing between characters for an elegant look */
			word-wrap: break-word; /* Ensure the text wraps nicely if too long */
			line-height: 1.4; /* A bit of space between lines for clarity */
		}
	`
	return stylesheet
}

func CreateStaticDataLabelStylesheet() string {
	stylesheet := `
		QLabel[label-type="static"] {
			font-family: "Segoe UI", "Georgia", "Arial", sans-serif; /* Elegant font with a mix of modern and classic */
			font-size: 9px; /* Size large enough for readability, but elegant */
			color: #839496; /* Dark grey color for a soft yet sophisticated text color */
			background-color: transparent; /* No background, letting the label sit naturally */
			padding: 5px 15px; /* Padding around the text to give it some space */
			border-radius: 10px; /* Slightly rounded corners for the label */
			text-align: center; /* Centered text */
			letter-spacing: 0.5px; /* Slight spacing between characters for an elegant look */
			word-wrap: break-word; /* Ensure the text wraps nicely if too long */
			line-height: 1.4; /* A bit of space between lines for clarity */
		}
	`
	return stylesheet

}

func CreateTeamStylesheet() string {
	stylesheet := `
		QGroupBox[widget-type="team"] {
			opacity:0.77;
			background-color: #002b36;
		}
	`
	return stylesheet
}

func CreateGameStylesheet(homeTeam string, awayTeam string) string {
	imagePath := GetPng("puck_texture")
	stylesheet := `
		background-image: url(%s);
		background-repeat: no-repeat;
		background-position: center;
	`
	styleSheet := fmt.Sprintf(stylesheet, imagePath)
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
