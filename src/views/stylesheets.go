package views

import (
	"fmt"
	"path"
	"quickRadio/models"
	"quickRadio/quickio"
	"strconv"
)

func GetPng(imageTitle string) string {
	dir := quickio.GetProjectDirWithForwadSlash()
	imagePath := path.Join(dir, "assets", "pngs", imageTitle+".png")
	return imagePath
}

//Consider creating a Single Large stylesheet functionally, by combining these func ina func
//To generate a large stylesheet, but one that can broken up like this in funcs for ease of editing.

func CreateInactiveRadioStreamButtonStylesheet(sweater *models.Sweater) string {
	stylesheet := `
		QPushButton[button-type="stream"]{
			background: qlineargradient(x1:0 y1:0, x2:1, y2:1, stop:0 %s, stop:.40 %s, stop:.60 %s, stop:1.0 %s);
			border-radius: 9px; /* Rounded borders */
			border: 2px solid %s; /* Subtle border color to match the button */
			box-shadow: 0 4px 8px %s; /* Soft shadow for elegance rgba(0, 0, 0, 0.2)* #000000C0 */
			min-height: 150px;
			min-width: 150px;
			max-height: 150px;
			max-height: 150px;
		}
		QPushButton[button-type="stream"]:hover{
			background: qlineargradient(x1:0 y1:0, x2:1, y2:1, stop:0 %s, stop:.40 %s, stop:.60 %s, stop:1.0 %s);
		}
	`
	styleSheet := fmt.Sprintf(stylesheet, sweater.PrimaryColor, sweater.PrimaryColor, sweater.SecondaryColor, sweater.SecondaryColor, sweater.PrimaryColor, sweater.PrimaryColor, //Default
		sweater.SecondaryColor, sweater.SecondaryColor, sweater.PrimaryColor, sweater.PrimaryColor, //hover
	)
	return styleSheet
}

func CreateActiveRadioStreamButtonStylesheet(sweater *models.Sweater) string {
	stylesheet := `
		QPushButton[button-type="stream"]{
			background-color:%s;
			border-radius: 9px; /* Rounded borders */
    		border: 2px solid %s; /* Subtle border color to match the button */
    		box-shadow: 0 4px 8px %s; /* Soft shadow for elegance rgba(0, 0, 0, 0.2) #000000FF*/
			min-height: 150px;
			min-width: 150px;
			max-height: 150px;
			max-height: 150px;
		}
		QPushButton[button-type="stream"]:hover{
			background-color:%s;
			border: 2px solid %s; /* Subtle border color to match the button */

		}
	`
	styleSheet := fmt.Sprintf(stylesheet,
		sweater.PrimaryColor, sweater.SecondaryColor, sweater.PrimaryColor, //Default
		sweater.SecondaryColor, sweater.PrimaryColor, //Hover
	)
	return styleSheet
}

func CreateGlassButtonStylesheet(sweater *models.Sweater) string {
	stylesheet := `
		QPushButton[button-type="glass"]{
			background: %s;
			border: 1px solid %s;
			border-radius: 3px;
			padding: 5px;
			color: %s;
			font-size: 14px;
			font-weight: bold;
			box-shadow: 0px 0px 15px rgba(0, 0, 0, 0.5);
		}
		QPushButton[button-type="glass"]:hover {
			background: %s;
			color: %s;
		}
		QPushButton[button-type="glass"]:pressed {
			background: %s;
		}
		QPushButton[button-type="glass"]:checked {
			background: %s;
			border: 1px solid %s;
			color: %s;
		}
		QPushButton[button-type="glass"]:disabled {
			background: %s;
			color: %s;
			opacity: 0.6;
			border: none;
			cursor: not-allowed;
		}
		QPushButton[button-type="glass"]:checked:disabled {
			background: %s;
			color: %s;
			opacity: 0.6;
			border:none;
			cursor: not-allowed;
		}
	`
	return fmt.Sprintf(stylesheet,
		sweater.SecondaryColor, sweater.PrimaryColor, sweater.PrimaryColor, //default
		sweater.PrimaryColor, sweater.SecondaryColor, //hover
		sweater.SecondaryColor,                                               //pressed
		sweater.PrimaryColor, sweater.SecondaryColor, sweater.SecondaryColor, //checked
		sweater.SecondaryColor, sweater.SecondaryColor, //disabled
		sweater.PrimaryColor, sweater.SecondaryColor, //checked:disabled
	)
}

func CreateDynamicDataLabelStylesheet(fontSize int) string {
	stylesheet := `
		QLabel[label-type="dynamic"] {
			font-family: Segoe UI, Georgia, Arial, sans-serif; /* Elegant font with a mix of modern and classic */
			font-weight: bold;
			font-size: %spx; /* Size large enough for readability, but elegant */
			color: #839496; /* Dark grey color for a soft yet sophisticated text color */
			background-color: transparent; /* No background, letting the label sit naturally */
			padding: 5px 15px; /* Padding around the text to give it some space */
			border-radius: 10px; /* Slightly rounded corners for the label */
			text-align: center; /* Centered text */
			letter-spacing: 0.5px; /* Slight spacing between characters for an elegant look */
			line-height: 1.3; /* A bit of space between lines for clarity */
		}
	`
	return fmt.Sprintf(stylesheet, strconv.Itoa(fontSize))
}

func CreateStaticDataLabelStylesheet(fontSize int) string {
	stylesheet := `
		QLabel[label-type="static"] {
			font-family: Segoe UI, Georgia, Arial, sans-serif; /* Elegant font with a mix of modern and classic */
			font-weight: bold;
			font-size: %spx; /* Size large enough for readability, but elegant */
			color: #839496; /* Dark grey color for a soft yet sophisticated text color */
			background-color: transparent; /* No background, letting the label sit naturally */
			padding: 3px 3px; /* Padding around the text to give it some space */
			border-radius: 10px; /* Slightly rounded corners for the label */
			text-align: center; /* Centered text */
			letter-spacing: 0.5px; /* Slight spacing between characters for an elegant look */
			line-height: 1.3; /* A bit of space between lines for clarity */
		}
	`
	return fmt.Sprintf(stylesheet, strconv.Itoa(fontSize))

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
