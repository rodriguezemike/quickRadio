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

func (widget *TeamWidget) createStaticDataLabel(name string, data string, gameWidget *widgets.QGroupBox) *widgets.QLabel {
	label := widgets.NewQLabel2(data, gameWidget, core.Qt__Widget)
	label.SetObjectName(name)
	label.SetProperty("label-type", core.NewQVariant12("static"))
	label.SetStyleSheet(CreateStaticDataLabelStylesheet())
	log.Println("Static Label stylesheet", label.StyleSheet())
	return label
}

func (widget *TeamWidget) createDynamicDataLabel(name string, data string, gamecenterLink string, gameWidget *widgets.QGroupBox) *widgets.QLabel {
	label := widgets.NewQLabel2(data, gameWidget, core.Qt__Widget)
	label.SetObjectName(name)
	label.SetAccessibleDescription(gamecenterLink)
	label.SetProperty("label-type", core.NewQVariant12("dynamic"))
	label.SetStyleSheet(CreateDynamicDataLabelStylesheet())
	label.ConnectTimerEvent(func(event *core.QTimerEvent) {
		if label.IsVisible() {
			val, ok := widget.updateMap[label.ObjectName()]
			if !ok || !val {
				if strings.Contains(label.ObjectName(), "SCORE") {
					teamAbbrev, dataLabel := widget.getTeamDataFromUIObjectName(label.ObjectName(), "_")
					label.SetText(teamAbbrev + "  " + widget.gameController.GetUIDataFromFilename(teamAbbrev, dataLabel, "0"))
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

func (widget *TeamWidget) createRadioQualityButtons(gameWidget *widgets.QGroupBox) (*widgets.QPushButton, *widgets.QPushButton) {
	buttonLow := widgets.NewQPushButton(gameWidget)
	buttonLow.SetText("Bad Internet")
	buttonLow.SetToolTip("124K")
	buttonLow.SetCheckable(true)
	buttonLow.SetCheckedDefault(false)
	buttonLow.SetEnabled(false)
	buttonHigh := widgets.NewQPushButton(gameWidget)
	buttonHigh.SetText("Good Internet")
	buttonHigh.SetToolTip("192K")
	buttonHigh.SetCheckable(true)
	buttonHigh.SetCheckedDefault(true)
	buttonHigh.ConnectToggled(func(onCheck bool) {
		if onCheck {
			buttonLow.SetEnabled(false)
			buttonLow.SetChecked(false)
		} else {
			buttonLow.SetEnabled(true)
		}
	})
	buttonLow.ConnectToggled(func(onCheck bool) {
		if onCheck {
			buttonHigh.SetEnabled(false)
			buttonHigh.SetChecked(false)
		} else {
			buttonHigh.SetEnabled(true)
		}
	})
	return buttonLow, buttonHigh

}

func (widget *TeamWidget) createTeamRadioStreamButton(teamAbbrev string, radioLink string, gameWidget *widgets.QGroupBox, radioQualityButtonLow *widgets.QPushButton, radioQualityButtonHigh *widgets.QPushButton) *widgets.QPushButton {
	teamIcon := widget.getTeamIcon(teamAbbrev)
	button := widgets.NewQPushButton(gameWidget)
	button.SetToolTip(fmt.Sprintf("Play %s Radio", teamAbbrev))
	button.SetProperty("button-type", core.NewQVariant12("radio"))
	button.SetStyleSheet(CreateInactiveRadioStreamButtonStylesheet(widget.gameController.Sweaters[teamAbbrev]))
	button.SetObjectName(widget.setTeamDataUIObjectName(teamAbbrev, "RADIO", "_"))
	button.SetIcon(teamIcon)
	button.SetIconSize(button.FrameSize())
	button.SetCheckable(true)
	button.ConnectToggled(func(onCheck bool) {
		var radioSampleRate string
		teamAbbrev, _ = widget.getTeamDataFromUIObjectName(button.ObjectName(), "_")
		if radioQualityButtonHigh.IsChecked() {
			radioSampleRate = radioQualityButtonHigh.ToolTip()
		} else {
			radioSampleRate = radioQualityButtonLow.ToolTip()
		}
		if onCheck {
			if radioLink != "" {
				widget.radioController = controllers.NewRadioControllerWithLock(radioLink, teamAbbrev, radioSampleRate, widget.radioLock)
				go widget.radioController.PlayRadio()
			}
			button.SetStyleSheet(CreateActiveRadioStreamButtonStylesheet((widget.gameController.Sweaters[teamAbbrev])))
		} else {
			if radioLink != "" {
				go widget.radioController.StopRadio(teamAbbrev)
				widget.radioController = nil
			}
			button.SetStyleSheet(CreateInactiveRadioStreamButtonStylesheet((widget.gameController.Sweaters[teamAbbrev])))
		}
	})
	log.Println(button.StyleSheet())
	return button
}

func (widget *TeamWidget) ClearUpdateMap() {
	widget.updateMap = nil
	widget.updateMap = map[string]bool{}
}

func (widget *TeamWidget) createTeamWidget(team *models.TeamData, gamecenterLink string, gameWidget *widgets.QGroupBox) {
	if team == nil {
		team = models.CreateDefaultTeam()
	}
	teamLayout := widgets.NewQVBoxLayout2(gameWidget)
	radioQualityLayout := widgets.NewQHBoxLayout2(gameWidget)
	teamWidget := widgets.NewQGroupBox(gameWidget)
	teamWidget.SetProperty("widget-type", core.NewQVariant12("team"))
	//Radio Quality Buttons and Label (Static label)
	radioQualityButtonLow, radioQualityButtonHigh := widget.createRadioQualityButtons(gameWidget)
	radioQualityLayout.AddWidget(radioQualityButtonLow, 0, core.Qt__AlignCenter)
	radioQualityLayout.AddWidget(radioQualityButtonHigh, 0, core.Qt__AlignCenter)
	teamLayout.AddLayout(radioQualityLayout, 0)
	//Team Radio Stream Button and Data Labels (Dynamic Labels)
	teamLayout.AddWidget(widget.createTeamRadioStreamButton(team.Abbrev, team.RadioLink, gameWidget, radioQualityButtonLow, radioQualityButtonHigh), 0, core.Qt__AlignCenter)
	teamLayout.AddWidget(widget.createDynamicDataLabel(widget.setTeamDataUIObjectName(team.Abbrev, "SCORE", "_"), team.Abbrev+"  "+strconv.Itoa(team.Score), gamecenterLink, gameWidget), 0, core.Qt__AlignCenter)
	teamLayout.AddWidget(widget.createDynamicDataLabel(widget.setTeamDataUIObjectName(team.Abbrev, "SOG:", "_"), "SOG  "+strconv.Itoa(team.Sog), gamecenterLink, gameWidget), 0, core.Qt__AlignCenter)
	// Work off a scaling factor - base = 100 (base*1.77)*ScalingFactor and base*scalingFactor ::Scaling Factor is 2. :: 1.77 is Desired Aspect Ratio.
	teamWidget.SetMinimumSize(core.NewQSize2(200, 354))
	teamWidget.SetMaximumSize(core.NewQSize2(200, 354))
	teamWidget.SetLayout(teamLayout)
	teamWidget.SetStyleSheet(CreateTeamStylesheet())
	log.Println("Team Widget Stylesheet ", teamWidget.StyleSheet())
	widget.UI = teamWidget
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
