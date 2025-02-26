package views

import (
	"fmt"
	"log"
	"quickRadio/controllers"
	"quickRadio/quickio"
	"strconv"
	"strings"
	"sync"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

type TeamWidget struct {
	LabelTimer      int
	UIWidget        *widgets.QGroupBox
	UILayout        *widgets.QVBoxLayout
	parentWidget    *widgets.QGroupBox
	teamController  *controllers.TeamController
	radioController *controllers.RadioController
	radioLock       *sync.Mutex
	updateMap       map[string]bool
}

func GetTeamPixmap(teamAbbrev string) *gui.QPixmap {
	teamStringSlice := []string{teamAbbrev, "light.svg"}
	teamLogoFilename := strings.Join(teamStringSlice, "_")
	teamLogoPath := quickio.GetLogoPath(teamLogoFilename)
	teamPixmap := gui.NewQPixmap3(teamLogoPath, "svg", core.Qt__AutoColor)
	return teamPixmap
}

func GetTeamIcon(teamAbbrev string) *gui.QIcon {
	teamPixmap := GetTeamPixmap(teamAbbrev)
	teamIcon := gui.NewQIcon2(teamPixmap)
	return teamIcon
}

// All updated func needs to be added
func (widget *TeamWidget) ClearUpdateMap() {
	widget.updateMap = nil
	widget.updateMap = map[string]bool{}
}

func (widget *TeamWidget) RadioLockReferenceTest(lock *sync.Mutex) bool {
	return lock == widget.radioLock
}

func (widget TeamWidget) GameControllerReferenceTest(controller *controllers.TeamController) bool {
	return controller == widget.teamController
}

func (widget TeamWidget) LabelTimerTest(labelTimer int) bool {
	return widget.LabelTimer == labelTimer
}

func (widget *TeamWidget) getTeamDataFromUIObjectName(objectName string, delimiter string) (string, string) {
	objectNameSplit := strings.Split(objectName, delimiter)
	return objectNameSplit[0], objectNameSplit[1]
}

func (widget *TeamWidget) getTeamPixmap(teamAbbrev string) *gui.QPixmap {
	teamStringSlice := []string{teamAbbrev, "light.svg"}
	teamLogoFilename := strings.Join(teamStringSlice, "_")
	teamLogoPath := quickio.GetLogoPath(teamLogoFilename)
	teamPixmap := gui.NewQPixmap3(teamLogoPath, "svg", core.Qt__AutoColor)
	return teamPixmap
}

func (widget *TeamWidget) getTeamIcon(teamAbbrev string) *gui.QIcon {
	teamPixmap := widget.getTeamPixmap(teamAbbrev)
	teamIcon := gui.NewQIcon2(teamPixmap)
	return teamIcon
}

func (widget *TeamWidget) setTeamDataUIObjectName(teamAbbrev string, uiLabel string, delimiter string) string {
	return teamAbbrev + delimiter + uiLabel
}

func (widget *TeamWidget) createStaticDataLabel(name string, data string, fontSize int) *widgets.QLabel {
	label := widgets.NewQLabel2(data, widget.UIWidget, core.Qt__Widget)
	label.SetObjectName(name)
	label.SetProperty("label-type", core.NewQVariant12("static"))
	label.SetStyleSheet(CreateStaticDataLabelStylesheet(fontSize))
	log.Println("Static Label stylesheet", label.StyleSheet())
	return label
}

func (widget *TeamWidget) createDynamicDataLabel(name string, data string, gamecenterLink string, fontSize int) *widgets.QLabel {
	label := widgets.NewQLabel2(data, widget.UIWidget, core.Qt__Widget)
	label.SetObjectName(name)
	label.SetAccessibleDescription(gamecenterLink)
	label.SetProperty("label-type", core.NewQVariant12("dynamic"))
	label.SetStyleSheet(CreateDynamicDataLabelStylesheet(fontSize))
	label.ConnectTimerEvent(func(event *core.QTimerEvent) {
		if label.IsVisible() {
			val, ok := widget.updateMap[label.ObjectName()]
			if !ok || !val {
				if strings.Contains(label.ObjectName(), "SCORE") {
					_, dataLabel := widget.getTeamDataFromUIObjectName(label.ObjectName(), "_")
					label.SetText(widget.teamController.GetUIDataFromFilename(dataLabel, "-10"))
					label.Repaint()
				} else if strings.Contains(label.ObjectName(), "SOG") {
					_, dataLabel := widget.getTeamDataFromUIObjectName(label.ObjectName(), "_")
					label.SetText(widget.teamController.GetUIDataFromFilename(dataLabel, "-10"))
					label.Repaint()
				}
				widget.updateMap[label.ObjectName()] = true
			}
		}
	})
	label.StartTimer(widget.LabelTimer, core.Qt__PreciseTimer)
	log.Println("Dynamic Label stylesheet", label.StyleSheet())
	return label
}

func (widget *TeamWidget) createRadioQualityButtons() (*widgets.QPushButton, *widgets.QPushButton) {
	//Button Low Quality Audio
	sweater := widget.teamController.Sweater
	buttonLow := widgets.NewQPushButton(widget.UIWidget)
	buttonLow.SetProperty("button-type", core.NewQVariant12("glass"))
	buttonLow.SetText("Low")
	buttonLow.SetToolTip("124K - For slower networks (Either lots of people on the network or slow internet access)")
	buttonLow.SetCheckable(true)
	buttonLow.SetEnabled(true)
	//Button HIgh Quality Audio
	buttonHigh := widgets.NewQPushButton(widget.UIWidget)
	buttonHigh.SetProperty("button-type", core.NewQVariant12("glass"))
	buttonHigh.SetText("High")
	buttonHigh.SetToolTip("192K - For faster networks (Either less of people on the network or fast internet access)")
	buttonHigh.SetCheckable(true)
	buttonHigh.SetEnabled(true)
	//Connect Toggle Events
	buttonHigh.ConnectToggled(func(onCheck bool) {
		if onCheck {
			buttonLow.SetChecked(false)
		}
	})
	buttonLow.ConnectToggled(func(onCheck bool) {
		if onCheck {
			buttonHigh.SetChecked(false)
		}
	})
	//Stylesheets
	buttonLow.SetStyleSheet(CreateGlassButtonStylesheet(sweater))
	buttonHigh.SetStyleSheet(CreateGlassButtonStylesheet(sweater))
	return buttonLow, buttonHigh
}

func (widget *TeamWidget) createTeamRadioStreamButton(teamAbbrev string, radioLink string, radioQualityButtonLow *widgets.QPushButton, radioQualityButtonHigh *widgets.QPushButton, radioQualityLabel *widgets.QLabel) *widgets.QPushButton {
	teamIcon := widget.getTeamIcon(teamAbbrev)
	button := widgets.NewQPushButton(widget.UIWidget)
	button.SetToolTip(fmt.Sprintf("Play %s Radio", teamAbbrev))
	button.SetProperty("button-type", core.NewQVariant12("stream"))
	button.SetStyleSheet(CreateInactiveRadioStreamButtonStylesheet(widget.teamController.Sweater))
	button.SetObjectName(widget.setTeamDataUIObjectName(teamAbbrev, "RADIO", "_"))
	button.SetIcon(teamIcon)
	button.SetIconSize(button.FrameSize())
	button.SetCheckable(true)
	button.ConnectToggled(func(onCheck bool) {
		var radioSampleRate string
		teamAbbrev, _ = widget.getTeamDataFromUIObjectName(button.ObjectName(), "_")
		if radioQualityButtonLow.IsChecked() {
			radioSampleRate = strings.Split(radioQualityButtonLow.ToolTip(), " ")[0]
		} else {
			radioQualityButtonHigh.SetChecked(true)
			radioQualityButtonLow.SetChecked(false)
			radioSampleRate = strings.Split(radioQualityButtonHigh.ToolTip(), " ")[0]
		}
		if onCheck {
			if radioLink != "" {
				widget.radioController = controllers.NewRadioControllerWithLock(radioLink, teamAbbrev, radioSampleRate, widget.radioLock)
				go widget.radioController.PlayRadio()
			}
			button.SetStyleSheet(CreateActiveRadioStreamButtonStylesheet((widget.teamController.Sweater)))
			radioQualityButtonHigh.SetEnabled(false)
			radioQualityButtonLow.SetEnabled(false)
			button.SetToolTip(fmt.Sprintf("Stop %s Radio", teamAbbrev))
			radioQualityLabel.SetText(fmt.Sprintf("Playing - %s Radio", teamAbbrev))
			radioQualityLabel.Repaint()
		} else {
			if radioLink != "" {
				go widget.radioController.StopRadio(teamAbbrev)
				widget.radioController = nil
			}
			button.SetStyleSheet(CreateInactiveRadioStreamButtonStylesheet((widget.teamController.Sweater)))
			radioQualityButtonHigh.SetEnabled(true)
			radioQualityButtonLow.SetEnabled(true)
			button.SetToolTip(fmt.Sprintf("Play %s Radio", teamAbbrev))
			radioQualityLabel.SetText(fmt.Sprintf("Audio Quality - %s Radio", teamAbbrev))
			radioQualityLabel.Repaint()
		}
	})
	log.Println(button.StyleSheet())
	return button
}

func (widget *TeamWidget) createScoreLayout() *widgets.QHBoxLayout {
	fontSize := 32
	scoreLayout := widgets.NewQHBoxLayout()
	scoreLayout.AddWidget(widget.createStaticDataLabel("teamAbbrev", widget.teamController.Team.Abbrev, fontSize), 1, core.Qt__AlignCenter)
	scoreLayout.AddWidget(widget.createDynamicDataLabel(widget.setTeamDataUIObjectName(widget.teamController.Team.Abbrev, "SCORE", "_"), strconv.Itoa(widget.teamController.Team.Score), widget.teamController.Landinglink, fontSize), 2, core.Qt__AlignRight)
	return scoreLayout
}

func (widget *TeamWidget) createShotsOnGoalLayout() *widgets.QHBoxLayout {
	fontSize := 24
	shotsOnGoalLayout := widgets.NewQHBoxLayout()
	shotsOnGoalLayout.AddWidget(widget.createStaticDataLabel("sog", "SOG:", fontSize), 0, core.Qt__AlignCenter)
	shotsOnGoalLayout.AddWidget(widget.createDynamicDataLabel(widget.setTeamDataUIObjectName(widget.teamController.Team.Abbrev, "SOG", "_"), strconv.Itoa(widget.teamController.Team.Sog), widget.teamController.Landinglink, fontSize), 0, core.Qt__AlignCenter)
	return shotsOnGoalLayout
}

func (widget *TeamWidget) createRadioLayout() *widgets.QVBoxLayout {
	//Radio Quality Layout
	radioQualityButtonsLayout := widgets.NewQHBoxLayout()
	radioQualityButtonLow, radioQualityButtonHigh := widget.createRadioQualityButtons()
	radioQualityButtonsLayout.AddWidget(radioQualityButtonLow, 0, core.Qt__AlignCenter)
	radioQualityButtonsLayout.AddWidget(radioQualityButtonHigh, 0, core.Qt__AlignCenter)
	//Radio Layout
	radioLayout := widgets.NewQVBoxLayout()
	internetQualityLabel := widget.createStaticDataLabel("internetQuality", fmt.Sprintf("Audio Quality - %s Radio", widget.teamController.Team.Abbrev), 9)
	radioLayout.AddWidget(internetQualityLabel, 0, core.Qt__AlignLeft)
	radioLayout.AddLayout(radioQualityButtonsLayout, 0)
	radioLayout.AddWidget(widget.createTeamRadioStreamButton(widget.teamController.Team.Abbrev, widget.teamController.Team.RadioLink, radioQualityButtonLow, radioQualityButtonHigh, internetQualityLabel), 0, core.Qt__AlignCenter)
	return radioLayout
}

func (widget *TeamWidget) createSpacerLayout() *widgets.QVBoxLayout {
	verticalLayout := widgets.NewQVBoxLayout()
	verticalLayout.AddSpacerItem(widgets.NewQSpacerItem(100, 100, widgets.QSizePolicy__Minimum, widgets.QSizePolicy__Expanding))
	return verticalLayout
}

func (widget *TeamWidget) createTeamWidget() {
	//Default Team if we have No team for the games, useful for UI Testing
	//Create team layout, groupbox and set custom properties
	teamLayout := widgets.NewQVBoxLayout2(widget.parentWidget)
	teamGroupbox := widgets.NewQGroupBox(widget.parentWidget)
	//Set UI widget and Layout
	widget.UIWidget = teamGroupbox
	widget.UILayout = teamLayout
	teamGroupbox.SetProperty("widget-type", core.NewQVariant12("team"))
	//Create Child layouts
	scoreLayout := widget.createScoreLayout()
	shotsOnGoalLayout := widget.createShotsOnGoalLayout()
	radioLayout := widget.createRadioLayout()
	spacerLayout := widget.createSpacerLayout()
	//Add Child Layouts
	teamLayout.AddLayout(radioLayout, 0)
	teamLayout.AddLayout(scoreLayout, 0)
	teamLayout.AddLayout(shotsOnGoalLayout, 0)
	teamLayout.AddLayout(spacerLayout, 0)
	//Set Size and Stylesheet - Work off a scaling factor - base = 100 (base*1.77)*ScalingFactor and base*scalingFactor ::Scaling Factor is 2. :: 1.77 is Desired Aspect Ratio.
	teamGroupbox.SetMinimumSize(core.NewQSize2(200, 354))
	teamGroupbox.SetMaximumSize(core.NewQSize2(500, 885))
	teamGroupbox.SetLayout(teamLayout)
	teamGroupbox.SetStyleSheet(CreateTeamStylesheet())
	log.Println("Team Widget Stylesheet ", teamGroupbox.StyleSheet())
}

func CreateNewTeamWidget(labelTimer int, controller *controllers.TeamController, radioLock *sync.Mutex, gameWidget *widgets.QGroupBox) *TeamWidget {
	widget := TeamWidget{}
	widget.LabelTimer = labelTimer
	widget.teamController = controller
	widget.radioLock = radioLock
	widget.updateMap = map[string]bool{}
	widget.parentWidget = gameWidget
	widget.createTeamWidget()
	return &widget
}
