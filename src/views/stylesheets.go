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
		QPushButton[button-type="stream"]:disabled{
			background-color: %s;
			color-adjust: grayscale;
			border: none;
			cursor: not-allowed;
		}
	`
	styleSheet := fmt.Sprintf(stylesheet, sweater.PrimaryColor, sweater.PrimaryColor, sweater.SecondaryColor, sweater.SecondaryColor, sweater.PrimaryColor, sweater.PrimaryColor, //Default
		sweater.SecondaryColor, sweater.SecondaryColor, sweater.PrimaryColor, sweater.PrimaryColor, //hover
		sweater.SecondaryColor, //disabled
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
			border: 2px solid %s; 
		}
		QPushButton[button-type="stream"]:disabled{
			background-color: %s;
			color-adjust: grayscale;
			border: none;
			cursor: not-allowed;

		}
	`
	styleSheet := fmt.Sprintf(stylesheet,
		sweater.PrimaryColor, sweater.SecondaryColor, sweater.PrimaryColor, //Default
		sweater.SecondaryColor, sweater.PrimaryColor, //Hover
		sweater.SecondaryColor, //disabled
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
			line-height: 2.0; /* A bit of space between lines for clarity */
		}
	`
	return fmt.Sprintf(stylesheet, strconv.Itoa(fontSize))

}

// Home on left or right? Well sort that out later. Holding the idea
func CreateSliderStylesheet(homeSweater models.Sweater, awaySweater models.Sweater, homeHandle bool) string {
	stylesheet := `
		QSlider {
			min-width: 300px;
			min-height: 10px;
		}
		
		QSlider:disabled {
			background-color: #D3D3D3; /* Light grey background for the whole slider */
		}
		
		QSlider:disabled::groove:horizontal {
			background: #A9A9A9; /* Darker grey track color */
			border-radius: 4px;
			height: 8px;
		}
		
		QSlider:disabled::handle:horizontal {
			background: %s; /* Specific color for the handle (tomato red) */
			border: 2px solid %s; /* Light grey border for the handle */
			width: 20px;
			height: 20px;
			border-radius: 10px; /* Round shape for the handle */
		}
		
		QSlider:disabled::sub-control:horizontal {
			background: %s; /* Color of the left side of the handle */
		}
		
		QSlider:disabled::sub-control:horizontal:handle {
			background: %s; /* Color of the right side of the handle */
		}
	`
	if homeHandle {
		return fmt.Sprintf(stylesheet,
			homeSweater.PrimaryColor, homeSweater.SecondaryColor,
			homeSweater.PrimaryColor,
			homeSweater.PrimaryColor,
		)
	} else {
		return fmt.Sprintf(stylesheet,
			awaySweater.PrimaryColor, awaySweater.SecondaryColor,
			awaySweater.PrimaryColor,
			awaySweater.PrimaryColor,
		)
	}

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

func CreateGameStatsAndGamestateStylesheet() string {
	stylesheet := `
		QGroupBox[widget-type="gamestatsAndGamestate"] {
			opacity:0.77;
			background-color: #002b36;
			border:none;
		}
	`
	return stylesheet
}

func CreateGameStylesheet() string {
	stylesheet := `
		QGroupBox[view-type="gameView"] {
			opacity:0.77;
			background-color: #002b36;
		}
	`
	return stylesheet
}
func CreateDropdownStyleSheet() string {
	styleSheet := `
		background-color:rgb(168, 168, 168)
	`
	return styleSheet
}
func CreateGameManagerStyleSheet() string {
	styleSheet := `
		background-color : #002b36
	`
	return styleSheet
}

func CreateVolumeSliderStylesheet() string {
	stylesheet := `
		QSlider::groove:horizontal {
			border: 1px solid #999999;
			height: 8px;
			background: qlineargradient(x1:0, y1:0, x2:0, y2:1, stop:0 #B1B1B1, stop:1 #c4c4c4);
			margin: 2px 0;
		}
		QSlider::handle:horizontal {
			background: qlineargradient(x1:0, y1:0, x2:1, y2:1, stop:0 #b4b4b4, stop:1 #8f8f8f);
			border: 1px solid #5c5c5c;
			width: 18px;
			margin: -2px 0;
			border-radius: 3px;
		}
	`
	return stylesheet
}

func CreatePowerplayLabelStylesheet() string {
	stylesheet := `
		QLabel {
			background: qlineargradient(x1:0, y1:0, x2:1, y2:1, stop:0 #FF6B35, stop:0.5 #FF4500, stop:1 #FF6B35);
			color: white;
			font-size: 18px;
			font-weight: bold;
			padding: 8px;
			border-radius: 5px;
			border: 2px solid #FF4500;
			text-align: center;
		}
	`
	return stylesheet
}

func CreateVolumeControlWidgetStylesheet() string {
	stylesheet := `
		QGroupBox[widget-type="volumeControl"] {
			background-color: transparent;
			border: none;
		}
		QLabel[objectName="volumeLabel"] {
			color: white;
			font-size: 12px;
			font-weight: bold;
		}
		QPushButton[objectName="muteButton"] {
			background-color: #2aa198;
			border: 1px solid #268bd2;
			border-radius: 3px;
			color: white;
			font-size: 16px;
			padding: 2px;
		}
		QPushButton[objectName="muteButton"]:hover {
			background-color: #268bd2;
		}
		QPushButton[objectName="muteButton"]:pressed {
			background-color: #073642;
		}
		QPushButton[objectName="muteButton"]:disabled {
			background-color: #586e75;
			color: #93a1a1;
		}
	`
	return stylesheet
}
