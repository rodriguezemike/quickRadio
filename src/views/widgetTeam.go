package views

import (
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

func (widget *TeamWidget) createDataLabel(name string, data string, genercenterLink string, gameWidget *widgets.QGroupBox) *widgets.QLabel {
	label := widgets.NewQLabel2(data, gameWidget, core.Qt__Widget)
	label.SetObjectName(name)
	label.SetAccessibleDescription(genercenterLink)
	font := label.Font()
	font.SetPointSize(32)
	label.SetFont(font)
	label.SetStyleSheet(CreateLabelStylesheet())
	label.ConnectTimerEvent(func(event *core.QTimerEvent) {
		if label.IsVisible() {
			val, ok := widget.updateMap[label.ObjectName()]
			if !ok || !val {
				if strings.Contains(label.ObjectName(), "SCORE") {
					teamAbbrev, dataLabel := widget.getTeamDataFromUIObjectName(label.ObjectName(), "_")
					label.SetText(widget.gameController.GetUIDataFromFilename(teamAbbrev, dataLabel, "0"))
					label.Repaint()
				}
				widget.updateMap[label.ObjectName()] = true
			}
		}
	})
	label.StartTimer(widget.LabelTimer, core.Qt__PreciseTimer)
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

func (widget *TeamWidget) createRadioQualityButton(gameWidget *widgets.QGroupBox) *widgets.QGroupBox {
	layout := widgets.NewQHBoxLayout2(gameWidget)
	radioQualityWidget := widgets.NewQGroupBox(gameWidget)
	radioQualityButtons := widgets.NewQButtonGroup(radioQualityWidget)

	for index, val := range []string{"124K", "192K"} {
		button := widgets.NewQPushButton(radioQualityWidget)
		button.SetText(val)
		button.SetCheckable(true)
		radioQualityButtons.AddButton(button, index)
	}
	layout.AddWidget(radioQualityWidget, 0, core.Qt__AlignCenter)
	radioQualityWidget.SetLayout(layout)
	return radioQualityWidget
}

func (widget *TeamWidget) createTeamRadioStreamButton(teamAbbrev string, radioLink string, gameWidget *widgets.QGroupBox) *widgets.QPushButton {
	radioSampleRate := "192K"
	teamIcon := widget.getTeamIcon(teamAbbrev)
	button := widgets.NewQPushButton(gameWidget)
	button.SetStyleSheet(CreateTeamButtonStylesheet(widget.gameController.Sweaters[teamAbbrev]))
	button.SetObjectName(widget.setTeamDataUIObjectName(teamAbbrev, "RADIO", "_"))
	button.SetIcon(teamIcon)
	button.SetIconSize(button.FrameSize())
	button.SetCheckable(true)
	button.ConnectToggled(func(onCheck bool) {
		teamAbbrev, _ = widget.getTeamDataFromUIObjectName(button.ObjectName(), "_")
		if onCheck {
			if radioLink != "" {
				widget.radioController = controllers.NewRadioControllerWithLock(radioLink, teamAbbrev, radioSampleRate, widget.radioLock)
				go widget.radioController.PlayRadio()
			}
		} else {
			if radioLink != "" {
				go widget.radioController.StopRadio(teamAbbrev)
				widget.radioController = nil
			}
		}
	})
	return button
}

func (widget *TeamWidget) ClearUpdateMap() {
	widget.updateMap = nil
	widget.updateMap = map[string]bool{}
}

func (widget *TeamWidget) createTeamWidget(team *models.TeamData, gamecenterLink string, gameWidget *widgets.QGroupBox) {
	if team == nil {
		log.Println("Team is Nil ")
		team = models.CreateDefaultTeam()
	}
	log.Println("Team ", team)
	teamLayout := widgets.NewQVBoxLayout2(gameWidget)
	teamWidget := widgets.NewQGroupBox(gameWidget)
	teamLayout.AddWidget(widget.createTeamRadioStreamButton(team.Abbrev, team.RadioLink, gameWidget), 0, core.Qt__AlignCenter)
	teamLayout.AddWidget(widget.createDataLabel("TeamAbbrev", team.Abbrev, gamecenterLink, gameWidget), 0, core.Qt__AlignCenter)
	teamLayout.AddWidget(widget.createDataLabel(widget.setTeamDataUIObjectName(team.Abbrev, "SCORE", "_"), strconv.Itoa(team.Score), gamecenterLink, gameWidget), 0, core.Qt__AlignCenter)
	teamLayout.AddWidget(widget.createDataLabel(widget.setTeamDataUIObjectName(team.Abbrev, "SOG", "_"), "SOG: "+strconv.Itoa(team.Sog), gamecenterLink, gameWidget), 0, core.Qt__AlignCenter)
	teamLayout.AddWidget(widget.createRadioQualityButton(gameWidget), 0, core.Qt__AlignCenter)
	teamWidget.SetLayout(teamLayout)
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
