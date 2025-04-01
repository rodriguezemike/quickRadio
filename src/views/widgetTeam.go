package views

import (
	"fmt"
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
	LabelTimer             int
	QuickRadioGameView     *GameView
	UIWidget               *widgets.QGroupBox
	UILayout               *widgets.QVBoxLayout
	highQualityAudioButton *widgets.QPushButton
	lowQualityAudioButton  *widgets.QPushButton
	radioStreamButton      *widgets.QPushButton
	onIceStaticLabel       *widgets.QLabel
	sinBinStaticLabel      *widgets.QLabel
	parentWidget           *widgets.QGroupBox
	teamController         *controllers.TeamController
	radioController        *controllers.RadioController
	radioLock              *sync.Mutex
	isPlaying              bool
	updateMap              map[string]bool //Update to Sync map, use locks for now
	updateMapLock          *sync.RWMutex
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

func (widget *TeamWidget) AreStaticLabelsVisible() bool {
	return widget.onIceStaticLabel.IsVisible() && widget.sinBinStaticLabel.IsVisible()
}

func (widget *TeamWidget) SetStaticLabelVisibility(vis bool) {
	widget.onIceStaticLabel.SetVisible(vis)
	widget.sinBinStaticLabel.SetVisible(vis)
}

func (widget *TeamWidget) ClearUpdateMap() {
	for key := range widget.updateMap {
		widget.updateMapLock.Lock()
		widget.updateMap[key] = false
		widget.updateMapLock.Unlock()
	}
}

func (widget *TeamWidget) IsUpdated() bool {
	widget.updateMapLock.RLock()
	for _, v := range widget.updateMap {
		if !v {
			widget.updateMapLock.RUnlock()
			return false
		}
	}
	widget.updateMapLock.RUnlock()
	return true
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
	return objectNameSplit[0], strings.Join(objectNameSplit[1:], delimiter)
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

func (widget *TeamWidget) setTeamDataUIObjectName(uiLabel string, delimiter string) string {
	return widget.teamController.Team.Abbrev + delimiter + uiLabel
}

func (widget *TeamWidget) createStaticDataLabel(name string, data string, fontSize int) *widgets.QLabel {
	label := widgets.NewQLabel2(data, widget.UIWidget, core.Qt__Widget)
	label.SetObjectName(name)
	label.SetAlignment(core.Qt__AlignTop)
	label.SetProperty("label-type", core.NewQVariant12("static"))
	label.SetStyleSheet(CreateStaticDataLabelStylesheet(fontSize))
	return label
}

func (widget *TeamWidget) createDynamicDataLabel(name string, data string, gamecenterLink string, fontSize int, defaultString *string) *widgets.QLabel {
	label := widgets.NewQLabel2(data, widget.UIWidget, core.Qt__Widget)
	label.SetObjectName(name)
	label.SetAccessibleDescription(gamecenterLink)
	label.SetProperty("label-type", core.NewQVariant12("dynamic"))
	label.SetStyleSheet(CreateDynamicDataLabelStylesheet(fontSize))
	label.SetWordWrap(true)
	label.ConnectTimerEvent(func(event *core.QTimerEvent) {
		widget.updateMapLock.RLock()
		val, ok := widget.updateMap[label.ObjectName()]
		widget.updateMapLock.RUnlock()
		if !ok || !val {
			_, dataLabel := widget.getTeamDataFromUIObjectName(label.ObjectName(), "_")
			if defaultString == nil {
				label.SetText(widget.teamController.GetUIDataFromFilename(dataLabel, label.Text()))
			} else {
				label.SetText(widget.teamController.GetUIDataFromFilename(dataLabel, *defaultString))
			}
			label.Repaint()
			widget.updateMapLock.Lock()
			widget.updateMap[label.ObjectName()] = true
			widget.updateMapLock.Unlock()
		}
	})
	widget.updateMapLock.Lock()
	widget.updateMap[label.ObjectName()] = false
	widget.updateMapLock.Unlock()
	label.StartTimer(widget.LabelTimer, core.Qt__PreciseTimer)
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
	button.SetObjectName(widget.setTeamDataUIObjectName("RADIO", "_"))
	button.SetIcon(teamIcon)
	button.SetIconSize(button.FrameSize())
	button.SetCheckable(true)
	button.ConnectToggled(func(onCheck bool) {
		if widget.radioLock.TryLock() {
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
					button.SetStyleSheet(CreateActiveRadioStreamButtonStylesheet((widget.teamController.Sweater)))
					radioQualityButtonHigh.SetEnabled(false)
					radioQualityButtonLow.SetEnabled(false)
					button.SetToolTip(fmt.Sprintf("Stop %s Radio", teamAbbrev))
					radioQualityLabel.SetText(fmt.Sprintf("Playing - %s Radio", teamAbbrev))
					radioQualityLabel.Repaint()
					go widget.radioController.PlayRadio(radioSampleRate)
					widget.isPlaying = true
					if widget.QuickRadioGameView != nil {
						widget.QuickRadioGameView.PlayRadioMode(widget)
					}
				}
			}
		} else {
			if widget.isPlaying {
				if radioLink != "" {
					widget.radioController.StopRadio()
					widget.isPlaying = false
					widget.radioLock.Unlock()
					button.SetStyleSheet(CreateInactiveRadioStreamButtonStylesheet((widget.teamController.Sweater)))
					radioQualityButtonHigh.SetEnabled(true)
					radioQualityButtonLow.SetEnabled(true)
					button.SetToolTip(fmt.Sprintf("Play %s Radio", teamAbbrev))
					radioQualityLabel.SetText(fmt.Sprintf("Audio Quality - %s Radio", teamAbbrev))
					radioQualityLabel.Repaint()
					if widget.QuickRadioGameView != nil {
						widget.QuickRadioGameView.StopRadioMode(widget)
					}
				}
			}
		}
	})
	return button
}

func (widget *TeamWidget) createScoreLayout() *widgets.QHBoxLayout {
	fontSize := 32
	scoreLayout := widgets.NewQHBoxLayout()
	scoreLayout.AddWidget(widget.createStaticDataLabel("teamAbbrev", widget.teamController.Team.Abbrev, fontSize), 1, core.Qt__AlignCenter)
	scoreLayout.AddWidget(widget.createDynamicDataLabel(widget.setTeamDataUIObjectName("SCORE", "_"), strconv.Itoa(widget.teamController.Team.Score), widget.teamController.Landinglink, fontSize, nil), 2, core.Qt__AlignRight)
	scoreLayout.SetSizeConstraint(widgets.QLayout__SetFixedSize)
	return scoreLayout
}

func (widget *TeamWidget) createShotsOnGoalLayout() *widgets.QHBoxLayout {
	fontSize := 24
	shotsOnGoalLayout := widgets.NewQHBoxLayout()
	shotsOnGoalLayout.AddWidget(widget.createStaticDataLabel("sog", "SOG:", fontSize), 0, core.Qt__AlignCenter)
	shotsOnGoalLayout.AddWidget(widget.createDynamicDataLabel(widget.setTeamDataUIObjectName("SOG", "_"), strconv.Itoa(widget.teamController.Team.Sog), widget.teamController.Landinglink, fontSize, nil), 0, core.Qt__AlignCenter)
	shotsOnGoalLayout.SetSizeConstraint(widgets.QLayout__SetFixedSize)
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
	//create radiostrteam button
	radioStreamButton := widget.createTeamRadioStreamButton(widget.teamController.Team.Abbrev, widget.teamController.Team.RadioLink, radioQualityButtonLow, radioQualityButtonHigh, internetQualityLabel)
	radioLayout.AddWidget(radioStreamButton, 0, core.Qt__AlignCenter)
	widget.highQualityAudioButton = radioQualityButtonHigh
	widget.lowQualityAudioButton = radioQualityButtonLow
	widget.radioStreamButton = radioStreamButton
	return radioLayout
}

func (widget *TeamWidget) createSpacerLayout(spacerUnits int) *widgets.QVBoxLayout {
	verticalLayout := widgets.NewQVBoxLayout()
	verticalLayout.AddSpacerItem(widgets.NewQSpacerItem(spacerUnits, spacerUnits, widgets.QSizePolicy__Minimum, widgets.QSizePolicy__Expanding))
	return verticalLayout
}

func (widget *TeamWidget) CreatePlayerOnIceLayout(player *models.PlayerOnIce, labelIndex int, fontSize int) *widgets.QHBoxLayout {
	//Number (May need a stylesheet for this)
	playerLayout := widgets.NewQHBoxLayout()
	defaultString := " "
	sweaterNumberWidget := widget.createDynamicDataLabel(widget.setTeamDataUIObjectName("SWEATERNUMBER_"+strconv.Itoa(labelIndex), "_"), " ", widget.teamController.Landinglink, fontSize, &defaultString)
	playerLayout.AddWidget(sweaterNumberWidget, 0, core.Qt__AlignLeft)
	//FullName
	nameWidget := widget.createDynamicDataLabel(widget.setTeamDataUIObjectName("PLAYERNAME_"+strconv.Itoa(labelIndex), "_"), player.Name.Default, widget.teamController.Landinglink, fontSize, &defaultString)
	playerLayout.AddWidget(nameWidget, 0, core.Qt__AlignLeft)
	//Position
	positionWidget := widget.createDynamicDataLabel(widget.setTeamDataUIObjectName("POSITIONCODE_"+strconv.Itoa(labelIndex), "_"), player.PositionCode, widget.teamController.Landinglink, fontSize, &defaultString)
	playerLayout.AddWidget(positionWidget, 0, core.Qt__AlignRight)
	return playerLayout
}

func (widget *TeamWidget) createTeamOnIceLayout() *widgets.QVBoxLayout {
	verticalLayout := widgets.NewQVBoxLayout()
	labelIndex := 0
	POIFontSize := 20
	playerFontSize := 14
	widget.onIceStaticLabel = widget.createStaticDataLabel("playersOnIce", "On Ice", POIFontSize)
	verticalLayout.AddWidget(widget.onIceStaticLabel, 10, core.Qt__AlignLeft)
	defaultTeam := models.CreateDefaultTeamOnIce()
	for _, player := range defaultTeam.Forwards {
		verticalLayout.AddLayout(widget.CreatePlayerOnIceLayout(&player, labelIndex, playerFontSize), 0)
		labelIndex += 1
	}
	for _, player := range defaultTeam.Defensemen {
		verticalLayout.AddLayout(widget.CreatePlayerOnIceLayout(&player, labelIndex, playerFontSize), 0)
		labelIndex += 1
	}
	for _, player := range defaultTeam.Goalies {
		verticalLayout.AddLayout(widget.CreatePlayerOnIceLayout(&player, labelIndex, playerFontSize), 0)
		labelIndex += 1
	}
	verticalLayout.SetSizeConstraint(widgets.QLayout__SetFixedSize)
	return verticalLayout

}

func (widget *TeamWidget) createSinBinLayout() *widgets.QVBoxLayout {
	verticalLayout := widgets.NewQVBoxLayout()
	labelIndex := 10
	peanlityBoxFontSize := 20
	playerFontSize := 14
	widget.sinBinStaticLabel = widget.createStaticDataLabel("penalityBox", "Penality Box", peanlityBoxFontSize)
	verticalLayout.AddWidget(widget.sinBinStaticLabel, 0, core.Qt__AlignLeft)
	defaultTeam := models.CreateDefaultTeamOnIce()
	for _, player := range defaultTeam.PenaltyBox {
		verticalLayout.AddLayout(widget.CreatePlayerOnIceLayout(&player, labelIndex, playerFontSize), 0)
		labelIndex += 1
	}
	verticalLayout.SetSizeConstraint(widgets.QLayout__SetFixedSize)
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
	teamOnIceLayout := widget.createTeamOnIceLayout()
	sinBinLayout := widget.createSinBinLayout()
	spacerLayout := widget.createSpacerLayout(20)
	//Add Child Layouts
	teamLayout.AddLayout(radioLayout, 0)
	teamLayout.AddLayout(scoreLayout, 0)
	teamLayout.AddLayout(shotsOnGoalLayout, 0)
	teamLayout.AddLayout(teamOnIceLayout, 0)
	teamLayout.AddLayout(sinBinLayout, 0)
	teamLayout.AddLayout(spacerLayout, 0)
	//Set Size and Stylesheet - Work off a scaling factor - base = 100 (base*1.77)*ScalingFactor and base*scalingFactor ::Scaling Factor is 2. :: 1.77 is Desired Aspect Ratio.
	teamGroupbox.SetMinimumSize(core.NewQSize2(300, 770))
	teamGroupbox.SetMaximumSize(core.NewQSize2(500, 1080))
	teamGroupbox.SetLayout(teamLayout)
	teamGroupbox.SetStyleSheet(CreateTeamStylesheet())
}

func (widget *TeamWidget) DisableButtons() {
	widget.highQualityAudioButton.SetEnabled(false)
	widget.highQualityAudioButton.SetChecked(false)
	widget.lowQualityAudioButton.SetEnabled(false)
	widget.lowQualityAudioButton.SetChecked(false)
	widget.radioStreamButton.SetEnabled(false)
}

func (widget *TeamWidget) EnableButtons() {
	widget.highQualityAudioButton.SetEnabled(true)
	widget.lowQualityAudioButton.SetEnabled(true)
	widget.radioStreamButton.SetEnabled(true)
}

func CreateNewTeamWidget(labelTimer int, teamController *controllers.TeamController, radioLock *sync.Mutex, parentWidget *widgets.QGroupBox, quickRadioGameView *GameView) *TeamWidget {
	widget := TeamWidget{}
	widget.QuickRadioGameView = quickRadioGameView
	widget.LabelTimer = labelTimer
	widget.teamController = teamController
	widget.radioLock = radioLock
	widget.updateMap = map[string]bool{}
	widget.parentWidget = parentWidget
	widget.updateMapLock = &sync.RWMutex{}
	widget.radioController = controllers.NewRadioController(widget.teamController.Team.RadioLink, widget.teamController.Team.Abbrev)
	widget.createTeamWidget()
	return &widget
}
