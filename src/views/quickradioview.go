package views

import (
	"os"
	"path/filepath"
	"quickRadio/controllers"
	"quickRadio/models"
	"quickRadio/quickio"
	"strconv"
	"strings"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

type QuickRadioView struct {
	LabelTimer              int
	activeGameDataUpdateMap map[string]bool
	gameController          *controllers.GameController
	radioController         *controllers.RadioController
	app                     *widgets.QApplication
	window                  *widgets.QMainWindow
	gameManagerWidget       *widgets.QGroupBox
	gamesStackWidget        *widgets.QStackedWidget
	activeGameWidget        *widgets.QWidget
	activeGameIndex         int
}

func (view *QuickRadioView) getIceRinkPixmap() *gui.QPixmap {
	dir := quickio.GetProjectDir()
	path := filepath.Join(dir, "assets", "svgs", "rink.svg")
	rinkPixmap := gui.NewQPixmap3(path, "svg", core.Qt__AutoColor)
	return rinkPixmap
}

func (view *QuickRadioView) getTeamPixmap(teamAbbrev string) *gui.QPixmap {
	teamStringSlice := []string{teamAbbrev, "light.svg"}
	teamLogoFilename := strings.Join(teamStringSlice, "_")
	teamLogoPath := quickio.GetLogoPath(teamLogoFilename)
	teamPixmap := gui.NewQPixmap3(teamLogoPath, "svg", core.Qt__AutoColor)
	return teamPixmap
}

func (view *QuickRadioView) GetTeamIcon(teamAbbrev string) *gui.QIcon {
	teamPixmap := view.getTeamPixmap(teamAbbrev)
	teamIcon := gui.NewQIcon2(teamPixmap)
	return teamIcon
}

func (view *QuickRadioView) SetTeamDataUIObjectName(teamAbbrev string, uiLabel string, delimiter string) string {
	return teamAbbrev + delimiter + uiLabel
}

func (view *QuickRadioView) GetTeamDataFromUIObjectName(objectName string, delimiter string) (string, string) {
	objectNameSplit := strings.Split(objectName, delimiter)
	return objectNameSplit[0], objectNameSplit[1]
}

func (view *QuickRadioView) CreateTeamRadioStreamButton(teamAbbrev string, radioLink string, gameWidget *widgets.QGroupBox) *widgets.QPushButton {
	teamIcon := view.GetTeamIcon(teamAbbrev)
	button := widgets.NewQPushButton(gameWidget)
	button.SetStyleSheet(CreateTeamButtonStylesheet(view.gameController.Sweaters[teamAbbrev]))
	button.SetObjectName(view.SetTeamDataUIObjectName(teamAbbrev, "RADIO", "_"))
	button.SetIcon(teamIcon)
	button.SetIconSize(button.FrameSize())
	button.SetCheckable(true)
	button.ConnectToggled(func(onCheck bool) {
		teamAbbrev, _ = view.GetTeamDataFromUIObjectName(button.ObjectName(), "_")
		if onCheck {
			if radioLink != "" && !quickio.IsRadioLocked() {
				view.radioController = controllers.NewRadioController(radioLink, teamAbbrev)
				go view.radioController.PlayRadio()
			}
		} else {
			if radioLink != "" && quickio.IsOurRadioLocked(teamAbbrev) {
				go view.radioController.StopRadio()
				view.radioController = nil
			}
		}
	})
	return button
}

func (view *QuickRadioView) CreateDataLabel(name string, data string, genercenterLink string, gameWidget *widgets.QGroupBox) *widgets.QLabel {
	label := widgets.NewQLabel2(data, gameWidget, core.Qt__Widget)
	label.SetObjectName(name)
	label.SetAccessibleDescription(genercenterLink)
	font := label.Font()
	font.SetPointSize(32)
	label.SetFont(font)
	label.SetStyleSheet(CreateLabelStylesheet())
	label.ConnectTimerEvent(func(event *core.QTimerEvent) {
		if label.IsVisible() {
			val, ok := view.activeGameDataUpdateMap[label.ObjectName()]
			if !ok || !val {
				if strings.Contains(label.ObjectName(), "SCORE") {
					teamAbbrev, dataLabel := view.GetTeamDataFromUIObjectName(label.ObjectName(), "_")
					label.SetText(view.gameController.GetUIDataFromFilename(teamAbbrev, dataLabel, "-1"))
				} else if strings.Contains(label.ObjectName(), "GAMESTATE") {
					label.SetText(view.gameController.GetActiveGamestateFromFile())
				}
				label.Repaint()
				view.activeGameDataUpdateMap[label.ObjectName()] = true
			}
		}
	})
	label.StartTimer(view.LabelTimer, core.Qt__PreciseTimer)
	return label
}
func (view *QuickRadioView) CreateTeamWidget(team models.TeamData, gamecenterLink string, gameWidget *widgets.QGroupBox) *widgets.QGroupBox {
	teamLayout := widgets.NewQVBoxLayout2(gameWidget)
	teamWidget := widgets.NewQGroupBox(gameWidget)
	teamLayout.AddWidget(view.CreateTeamRadioStreamButton(team.Abbrev, team.RadioLink, gameWidget), 0, core.Qt__AlignCenter)
	teamLayout.AddWidget(view.CreateDataLabel("TeamAbbrev", team.Abbrev, gamecenterLink, gameWidget), 0, core.Qt__AlignCenter)
	teamLayout.AddWidget(view.CreateDataLabel(view.SetTeamDataUIObjectName(team.Abbrev, "SOG", "_"), strconv.Itoa(team.Sog), gamecenterLink, gameWidget), 0, core.Qt__AlignCenter)
	teamLayout.AddWidget(view.CreateDataLabel(view.SetTeamDataUIObjectName(team.Abbrev, "SCORE", "_"), strconv.Itoa(team.Score), gamecenterLink, gameWidget), 0, core.Qt__AlignCenter)
	teamWidget.SetLayout(teamLayout)
	return teamWidget
}

func (view *QuickRadioView) CreateIceCenterWidget(gameDataObject models.GameData, gamecenterLink string, gameWidget *widgets.QGroupBox) *widgets.QGroupBox {
	layout := widgets.NewQVBoxLayout2(gameWidget)
	centerIceWidget := widgets.NewQGroupBox(gameWidget)
	gameStateLabel := view.CreateDataLabel("GAMESTATE", view.gameController.GetGamestateString(&gameDataObject), gamecenterLink, gameWidget)
	layout.AddWidget(CreateIceRinklabel(gameWidget, gameDataObject, view), 0, core.Qt__AlignCenter)
	layout.AddWidget(gameStateLabel, 0, core.Qt__AlignCenter)
	centerIceWidget.SetLayout(layout)
	centerIceWidget.ConnectTimerEvent(func(event *core.QTimerEvent) {
		for _, value := range view.activeGameDataUpdateMap {
			if !value {
				return
			}
		}
		for key := range view.activeGameDataUpdateMap {
			view.activeGameDataUpdateMap[key] = false
			view.gameController.ConsumeActiveGameData()
		}
	})
	centerIceWidget.StartTimer(view.LabelTimer*2, core.Qt__CoarseTimer)
	return centerIceWidget
}

func (view *QuickRadioView) CreateGameWidgetFromGameDataObject(gameDataObject models.GameData, gamecenterLink string) *widgets.QGroupBox {
	layout := widgets.NewQGridLayout(view.gamesStackWidget)
	gameWidget := widgets.NewQGroupBox(view.gamesStackWidget)
	layout.AddWidget2(view.CreateTeamWidget(gameDataObject.HomeTeam, gamecenterLink, gameWidget), 0, 0, core.Qt__AlignCenter)
	layout.AddWidget2(view.CreateIceCenterWidget(gameDataObject, gamecenterLink, gameWidget), 0, 1, core.Qt__AlignCenter)
	layout.AddWidget2(view.CreateTeamWidget(gameDataObject.AwayTeam, gamecenterLink, gameWidget), 0, 2, core.Qt__AlignCenter)
	gameWidget.SetLayout(layout)
	gameWidget.SetStyleSheet(CreateGameStylesheet(gameDataObject.HomeTeam.Abbrev, gameDataObject.AwayTeam.Abbrev))
	return gameWidget
}

func (view *QuickRadioView) CreateGamesWidget() *widgets.QStackedWidget {
	gameStackWidget := widgets.NewQStackedWidget(view.gameManagerWidget)
	for i, gameDataObject := range view.gameController.GetGameDataObjects() {
		gameWidget := view.CreateGameWidgetFromGameDataObject(gameDataObject, view.gameController.Landinglinks[i])
		gameStackWidget.AddWidget(gameWidget)
	}
	gameStackWidget.SetCurrentIndex(view.activeGameIndex)
	return gameStackWidget
}

func (view *QuickRadioView) CreateGameDropdownsWidget() *widgets.QComboBox {
	var gameNames []string
	for _, gameDataObject := range view.gameController.GetGameDataObjects() {
		gameNames = append(gameNames, gameDataObject.HomeTeam.Abbrev+" vs "+gameDataObject.AwayTeam.Abbrev)
	}
	dropdown := widgets.NewQComboBox(view.gameManagerWidget)
	dropdown.SetStyleSheet(CreateDropdownStyleSheet())
	dropdown.SetFixedWidth(600)
	dropdown.AddItems(gameNames)
	dropdown.ConnectCurrentIndexChanged(func(index int) {
		view.activeGameDataUpdateMap = nil
		view.activeGameDataUpdateMap = map[string]bool{}
		view.gameController.KillActiveGame()
		view.gameController.SwitchActiveObjects(index)
		view.gameController.RunActiveGame()
		view.gamesStackWidget.SetCurrentIndex(index)
		view.activeGameWidget = view.gamesStackWidget.CurrentWidget()
	})
	return dropdown
}

func (view *QuickRadioView) CreateGameManagerWidget() *widgets.QGroupBox {
	gameManager := widgets.NewQGroupBox(nil)
	topbarLayout := widgets.NewQVBoxLayout()
	gameStackLayout := widgets.NewQStackedLayout()
	gameManagerLayout := widgets.NewQVBoxLayout2(gameManager)

	gameManagerLayout.AddLayout(topbarLayout, 1)
	gameManagerLayout.AddLayout(gameStackLayout, 1)
	view.gamesStackWidget = view.CreateGamesWidget()
	gameDropdown := view.CreateGameDropdownsWidget()

	topbarLayout.AddWidget(gameDropdown, 1, core.Qt__AlignAbsolute)
	gameStackLayout.AddWidget(view.gamesStackWidget)
	gameManager.SetStyleSheet(CreateGameManagerStyleSheet())
	return gameManager
}

func (view *QuickRadioView) CreateLoadingScreen() *widgets.QSplashScreen {
	pixmap := view.getTeamPixmap("NHL")
	splash := widgets.NewQSplashScreen(pixmap, core.Qt__Widget)
	return splash
}

func (view *QuickRadioView) RunLoadingScream() {
	loadingScreen := view.CreateLoadingScreen()
	loadingScreen.Show()
	view.gameController = controllers.NewGameController()
	view.gameManagerWidget = view.CreateGameManagerWidget()
	view.window.SetCentralWidget(view.gameManagerWidget)
	loadingScreen.Finish(nil)
}

func (view *QuickRadioView) KillAllTheFun() {
	//Kill ALL THE FUN.
	//Stop all go rountines
	//Delete All the things
	//Close.
}

func (view *QuickRadioView) CreateAndRunApp() {
	quickio.EmptyTmpFolder()
	view.app = widgets.NewQApplication(len(os.Args), os.Args)
	view.app.SetApplicationDisplayName("QuickRadio")
	view.app.ConnectAboutToQuit(func() {
		view.KillAllTheFun()
	})
	view.app.ConnectDestroyQApplication(func() {
		view.KillAllTheFun()
	})
	view.app.SetWindowIcon(view.GetTeamIcon("NHLF"))
	view.window = widgets.NewQMainWindow(nil, 0)
	view.RunLoadingScream()

	view.window.Show()
	view.app.Exec()
}

func NewQuickRadioView() *QuickRadioView {
	var view QuickRadioView
	view.LabelTimer = 3000
	return &view
}
