package views

import (
	"quickRadio/controllers"
	"quickRadio/models"
	"sync"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

type GameView struct {
	LabelTimer              int
	GamecenterLink          string
	UIWidget                *widgets.QGroupBox
	UILayout                *widgets.QHBoxLayout
	AwayTeamWidget          *TeamWidget
	HomeTeamWidget          *TeamWidget
	GamestateAndStatsWidget *GamestateAndStatsWidget
	parentWidget            *widgets.QGroupBox
	radioLock               *sync.Mutex
	gameController          *controllers.GameController
}

func (view *GameView) PlayRadioMode(widget *TeamWidget) {
	if widget != view.AwayTeamWidget {
		view.AwayTeamWidget.DisableButtons()
	} else {
		view.HomeTeamWidget.DisableButtons()
	}
}

func (view *GameView) StopRadioMode(widget *TeamWidget) {
	if widget != view.AwayTeamWidget {
		view.AwayTeamWidget.EnableButtons()
	} else {
		view.HomeTeamWidget.EnableButtons()
	}
}

func (view *GameView) ClearUpdateMaps() {
	view.AwayTeamWidget.ClearUpdateMap()
	view.HomeTeamWidget.ClearUpdateMap()
	view.GamestateAndStatsWidget.ClearUpdateMap()
}

func (view *GameView) gameIsUpdated() bool {
	return view.AwayTeamWidget.IsUpdated() && view.HomeTeamWidget.IsUpdated() && view.GamestateAndStatsWidget.IsUpdated()
}

func (view *GameView) GetGameName() string {
	return view.gameController.HomeTeamController.Team.Abbrev + models.DEFAULT_VERSES_STRING + view.gameController.AwayTeamController.Team.Abbrev
}

func (view *GameView) createGameView() {
	view.gameController = controllers.CreateNewGameController(view.GamecenterLink)
	//Set UI widget and Layout
	viewLayout := widgets.NewQHBoxLayout()
	viewGroupBox := widgets.NewQGroupBox(view.parentWidget)
	viewGroupBox.SetProperty("view-type", core.NewQVariant12("gameView"))
	view.UILayout = viewLayout
	view.UIWidget = viewGroupBox
	view.UIWidget.ConnectTimerEvent(func(event *core.QTimerEvent) {
		if view.gameIsUpdated() {
			view.ClearUpdateMaps()
			view.gameController.ConsumeGameData()
		}
	})
	//Create Child layouts
	view.AwayTeamWidget = CreateNewTeamWidget(view.LabelTimer, view.gameController.AwayTeamController, view.radioLock, view.UIWidget, view)
	view.HomeTeamWidget = CreateNewTeamWidget(view.LabelTimer, view.gameController.HomeTeamController, view.radioLock, view.UIWidget, view)
	view.GamestateAndStatsWidget = CreateNewGamestateAndStatsWidget(view.LabelTimer, view.gameController, view.UIWidget)
	//Add Child Layouts
	view.UILayout.AddWidget(view.HomeTeamWidget.UIWidget, 0, core.Qt__AlignTop)
	view.UILayout.AddWidget(view.GamestateAndStatsWidget.UIWidget, 0, core.Qt__AlignTop)
	view.UILayout.AddWidget(view.AwayTeamWidget.UIWidget, 0, core.Qt__AlignTop)
	//Set Size and Stylesheet - Work off a scaling factor - base = 100 (base*1.77)*ScalingFactor and base*scalingFactor ::Scaling Factor is 2. :: 1.77 is Desired Aspect Ratio.
	view.UIWidget.SetMaximumSize(core.NewQSize2(1920, 1080))
	view.UIWidget.SetLayout(view.UILayout)
	view.UIWidget.SetStyleSheet(CreateGameStylesheet())
	view.UIWidget.StartTimer(view.LabelTimer, core.Qt__CoarseTimer)

}

// Mainly for testing purposes
func (view *GameView) createDefaultGameView() {
	view.gameController = controllers.CreateNewDefaultGameController()
	//Set UI widget and Layout
	viewLayout := widgets.NewQHBoxLayout()
	viewGroupBox := widgets.NewQGroupBox(view.parentWidget)
	viewGroupBox.SetProperty("view-type", core.NewQVariant12("gameView"))
	view.UILayout = viewLayout
	view.UIWidget = viewGroupBox
	view.UIWidget.ConnectTimerEvent(func(event *core.QTimerEvent) {
		//Run consume data here.

		if view.gameIsUpdated() {
			view.ClearUpdateMaps()
			view.gameController.ConsumeGameData()
		}
	})
	//Create Child layouts
	view.AwayTeamWidget = CreateNewTeamWidget(view.LabelTimer, view.gameController.AwayTeamController, view.radioLock, view.UIWidget, view)
	view.HomeTeamWidget = CreateNewTeamWidget(view.LabelTimer, view.gameController.HomeTeamController, view.radioLock, view.UIWidget, view)
	view.GamestateAndStatsWidget = CreateNewGamestateAndStatsWidget(view.LabelTimer, view.gameController, view.UIWidget)
	//Add Widgets
	view.UILayout.AddWidget(view.HomeTeamWidget.UIWidget, 0, core.Qt__AlignTop)
	view.UILayout.AddWidget(view.GamestateAndStatsWidget.UIWidget, 0, core.Qt__AlignTop)
	view.UILayout.AddWidget(view.AwayTeamWidget.UIWidget, 0, core.Qt__AlignTop)
	//Set Size and Stylesheet - Work off a scaling factor - base = 100 (base*1.77)*ScalingFactor and base*scalingFactor ::Scaling Factor is 2. :: 1.77 is Desired Aspect Ratio.
	view.UIWidget.SetMinimumSize(core.NewQSize2(1920, 1080))
	view.UIWidget.SetMaximumSize(core.NewQSize2(1920, 1080))
	view.UIWidget.SetLayout(view.UILayout)
	view.UIWidget.SetStyleSheet(CreateGameStylesheet())
	view.UIWidget.StartTimer(view.LabelTimer, core.Qt__CoarseTimer)
}

func CreateNewGameView(gamecenterLink string, parentWidget *widgets.QGroupBox, radioLock *sync.Mutex, labelTimer int) *GameView {
	//Create team layout, groupbox and set custom properties
	gameView := GameView{}
	gameView.LabelTimer = labelTimer
	gameView.parentWidget = parentWidget
	gameView.GamecenterLink = gamecenterLink
	gameView.radioLock = radioLock
	gameView.createGameView()
	return &gameView
}

// Mainly for testing purposes
func CreateNewDefaultGameView() *GameView {
	gameView := GameView{}
	gameView.LabelTimer = 1000
	gameView.parentWidget = widgets.NewQGroupBox(nil)
	gameView.GamecenterLink = ""
	gameView.radioLock = &sync.Mutex{}
	gameView.createDefaultGameView()
	return &gameView
}
