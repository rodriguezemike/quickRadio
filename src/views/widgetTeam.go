package views

import (
	"fmt"
	"log"
	"quickRadio/controllers"
	"quickRadio/models"
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
	UI              *widgets.QGroupBox
	gameController  *controllers.GameController
	radioController *controllers.RadioController
	radioLock       *sync.Mutex
	updateMap       map[string]bool
}

func (widget *TeamWidget) RadioLockReferenceTest(lock *sync.Mutex) bool {
	return lock == widget.radioLock
}

func (widget TeamWidget) GameControllerReferenceTest(controller *controllers.GameController) bool {
	return controller == widget.gameController
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

func (widget *TeamWidget) ClearUpdateMap() {
	widget.updateMap = nil
	widget.updateMap = map[string]bool{}
}

func (widget *TeamWidget) createStaticDataLabel(name string, data string, gameWidget *widgets.QGroupBox, fontSize int) *widgets.QLabel {
	label := widgets.NewQLabel2(data, gameWidget, core.Qt__Widget)
	label.SetObjectName(name)
	label.SetProperty("label-type", core.NewQVariant12("static"))
	label.SetStyleSheet(CreateStaticDataLabelStylesheet(fontSize))
	log.Println("Static Label stylesheet", label.StyleSheet())
	return label
}

func (widget *TeamWidget) createDynamicDataLabel(name string, data string, gamecenterLink string, gameWidget *widgets.QGroupBox, fontSize int) *widgets.QLabel {
	label := widgets.NewQLabel2(data, gameWidget, core.Qt__Widget)
	label.SetObjectName(name)
	label.SetAccessibleDescription(gamecenterLink)
	label.SetProperty("label-type", core.NewQVariant12("dynamic"))
	label.SetStyleSheet(CreateDynamicDataLabelStylesheet(fontSize))
	label.ConnectTimerEvent(func(event *core.QTimerEvent) {
		if label.IsVisible() {
			val, ok := widget.updateMap[label.ObjectName()]
			if !ok || !val {
				if strings.Contains(label.ObjectName(), "SCORE") {
					teamAbbrev, dataLabel := widget.getTeamDataFromUIObjectName(label.ObjectName(), "_")
					label.SetText(widget.gameController.GetUIDataFromFilename(teamAbbrev, dataLabel, "-10"))
					label.Repaint()
				} else if strings.Contains(label.ObjectName(), "SOG") {
					teamAbbrev, dataLabel := widget.getTeamDataFromUIObjectName(label.ObjectName(), "_")
					label.SetText(widget.gameController.GetUIDataFromFilename(teamAbbrev, dataLabel, "-10"))
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

func (widget *TeamWidget) createRadioQualityButtons(teamAbbrev string, teamWidget *widgets.QGroupBox) (*widgets.QPushButton, *widgets.QPushButton) {
	//Button Low Quality Audio
	sweater := widget.gameController.Sweaters[teamAbbrev]
	buttonLow := widgets.NewQPushButton(teamWidget)
	buttonLow.SetProperty("button-type", core.NewQVariant12("glass"))
	buttonLow.SetText("Low")
	buttonLow.SetToolTip("124K - For slower networks (Either lots of people on the network or slow internet access)")
	buttonLow.SetCheckable(true)
	buttonLow.SetEnabled(true)
	//Button HIgh Quality Audio
	buttonHigh := widgets.NewQPushButton(teamWidget)
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
	buttonLow.SetStyleSheet(CreateGlassButtonStylesheet(&sweater))
	buttonHigh.SetStyleSheet(CreateGlassButtonStylesheet(&sweater))
	return buttonLow, buttonHigh
}

func (widget *TeamWidget) createTeamRadioStreamButton(teamAbbrev string, radioLink string, teamWidget *widgets.QGroupBox, radioQualityButtonLow *widgets.QPushButton, radioQualityButtonHigh *widgets.QPushButton, radioQualityLabel *widgets.QLabel) *widgets.QPushButton {
	teamIcon := widget.getTeamIcon(teamAbbrev)
	button := widgets.NewQPushButton(teamWidget)
	button.SetToolTip(fmt.Sprintf("Play %s Radio", teamAbbrev))
	button.SetProperty("button-type", core.NewQVariant12("stream"))
	button.SetStyleSheet(CreateInactiveRadioStreamButtonStylesheet(widget.gameController.Sweaters[teamAbbrev]))
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
			button.SetStyleSheet(CreateActiveRadioStreamButtonStylesheet((widget.gameController.Sweaters[teamAbbrev])))
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
			button.SetStyleSheet(CreateInactiveRadioStreamButtonStylesheet((widget.gameController.Sweaters[teamAbbrev])))
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

func (widget *TeamWidget) createScoreLayout(team *models.TeamData, teamGroupbox *widgets.QGroupBox, gamecenterLink string) *widgets.QHBoxLayout {
	fontSize := 32
	scoreLayout := widgets.NewQHBoxLayout()
	scoreLayout.AddWidget(widget.createStaticDataLabel("teamAbbrev", team.Abbrev, teamGroupbox, fontSize), 1, core.Qt__AlignCenter)
	scoreLayout.AddWidget(widget.createDynamicDataLabel(widget.setTeamDataUIObjectName(team.Abbrev, "SCORE", "_"), strconv.Itoa(team.Score), gamecenterLink, teamGroupbox, fontSize), 2, core.Qt__AlignRight)
	return scoreLayout
}

func (widget *TeamWidget) createShotsOnGoalLayout(team *models.TeamData, teamGroupbox *widgets.QGroupBox, gamecenterLink string) *widgets.QHBoxLayout {
	fontSize := 24
	shotsOnGoalLayout := widgets.NewQHBoxLayout()
	shotsOnGoalLayout.AddWidget(widget.createStaticDataLabel("sog", "SOG:", teamGroupbox, fontSize), 0, core.Qt__AlignCenter)
	shotsOnGoalLayout.AddWidget(widget.createDynamicDataLabel(widget.setTeamDataUIObjectName(team.Abbrev, "SOG:", "_"), strconv.Itoa(team.Sog), gamecenterLink, teamGroupbox, fontSize), 0, core.Qt__AlignCenter)
	return shotsOnGoalLayout
}

func (widget *TeamWidget) createRadioLayout(team *models.TeamData, teamGroupbox *widgets.QGroupBox) *widgets.QVBoxLayout {
	//Radio Quality Layout
	radioQualityButtonsLayout := widgets.NewQHBoxLayout()
	radioQualityButtonLow, radioQualityButtonHigh := widget.createRadioQualityButtons(team.Abbrev, teamGroupbox)
	radioQualityButtonsLayout.AddWidget(radioQualityButtonLow, 0, core.Qt__AlignCenter)
	radioQualityButtonsLayout.AddWidget(radioQualityButtonHigh, 0, core.Qt__AlignCenter)
	//Radio Layout
	radioLayout := widgets.NewQVBoxLayout()
	internetQualityLabel := widget.createStaticDataLabel("internetQuality", fmt.Sprintf("Audio Quality - %s Radio", team.Abbrev), teamGroupbox, 9)
	radioLayout.AddWidget(internetQualityLabel, 0, core.Qt__AlignLeft)
	radioLayout.AddLayout(radioQualityButtonsLayout, 0)
	radioLayout.AddWidget(widget.createTeamRadioStreamButton(team.Abbrev, team.RadioLink, teamGroupbox, radioQualityButtonLow, radioQualityButtonHigh, internetQualityLabel), 0, core.Qt__AlignCenter)
	return radioLayout
}

func (widget *TeamWidget) createSpacerLayout() *widgets.QVBoxLayout {
	verticalLayout := widgets.NewQVBoxLayout()
	verticalLayout.AddSpacerItem(widgets.NewQSpacerItem(100, 100, widgets.QSizePolicy__Minimum, widgets.QSizePolicy__Expanding))
	return verticalLayout
}

func (widget *TeamWidget) createTeamWidget(team *models.TeamData, gamecenterLink string, gameWidget *widgets.QGroupBox) {
	//Default Team if we have No team for the games, useful for UI Testing
	if team == nil {
		team = models.CreateDefaultTeam()
	}
	//Create team layout, groupbox and set custom properties
	teamLayout := widgets.NewQVBoxLayout2(gameWidget)
	teamGroupbox := widgets.NewQGroupBox(gameWidget)
	teamGroupbox.SetProperty("widget-type", core.NewQVariant12("team"))
	//Create Child layouts
	scoreLayout := widget.createScoreLayout(team, teamGroupbox, gamecenterLink)
	shotsOnGoalLayout := widget.createShotsOnGoalLayout(team, teamGroupbox, gamecenterLink)
	radioLayout := widget.createRadioLayout(team, teamGroupbox)
	spacerLayout := widget.createSpacerLayout()
	//Add Child Layouts
	teamLayout.AddLayout(radioLayout, 0)
	teamLayout.AddLayout(scoreLayout, 0)
	teamLayout.AddLayout(shotsOnGoalLayout, 0)
	teamLayout.AddLayout(spacerLayout, 0)
	//Set Size and Stylesheet - Work off a scaling factor - base = 100 (base*1.77)*ScalingFactor and base*scalingFactor ::Scaling Factor is 2. :: 1.77 is Desired Aspect Ratio.
	teamGroupbox.SetMinimumSize(core.NewQSize2(200, 354))
	teamGroupbox.SetMaximumSize(core.NewQSize2(200, 354))
	teamGroupbox.SetLayout(teamLayout)
	teamGroupbox.SetStyleSheet(CreateTeamStylesheet())
	log.Println("Team Widget Stylesheet ", teamGroupbox.StyleSheet())
	//Set widget UI
	widget.UI = teamGroupbox
}

func CreateNewTeamWidget(labelTimer int, gameIndex int, homeTeam bool, controller *controllers.GameController, radioLock *sync.Mutex, gameWidget *widgets.QGroupBox) *TeamWidget {
	var team *models.TeamData
	var gamecenterLink string
	if gameIndex == -1 || len(controller.Landinglinks) == 0 {
		team = models.CreateDefaultTeam()
		gamecenterLink = ""
	} else {
		if homeTeam {
			team = &controller.GetGameDataObjects()[gameIndex].HomeTeam
		} else {
			team = &controller.GetGameDataObjects()[gameIndex].AwayTeam
		}
		gamecenterLink = controller.Landinglinks[gameIndex]
	}
	widget := TeamWidget{}
	widget.LabelTimer = labelTimer
	widget.gameController = controller
	widget.radioLock = radioLock
	widget.updateMap = map[string]bool{}
	widget.createTeamWidget(team, gamecenterLink, gameWidget)
	return &widget
}
