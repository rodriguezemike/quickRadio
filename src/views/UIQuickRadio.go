//go:build ignore
// +build ignore

package views

import (
	"os"
	"quickRadio/controllers"
	"quickRadio/quickio"
	"strings"

	"github.com/therecipe/qt/core"
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

func (view *QuickRadioView) GetTeamDataFromUIObjectName(objectName string, delimiter string) (string, string) {
	objectNameSplit := strings.Split(objectName, delimiter)
	return objectNameSplit[0], objectNameSplit[1]
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
					label.SetText(view.gameController.GetUIDataFromFilename(teamAbbrev, dataLabel, "0"))
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

func (view *QuickRadioView) CreateGamesWidget() *widgets.QStackedWidget {
	gameStackWidget := widgets.NewQStackedWidget(view.gameManagerWidget)
	for i, gameDataObject := range view.gameController.GetGameDataObjects() {
		gameWidget := view.CreateGameWidgetFromGameDataObject(gameDataObject, view.gameController.Landinglinks[i])
		gameStackWidget.AddWidget(gameWidget)
	}
	gameStackWidget.SetCurrentIndex(view.activeGameIndex)
	return gameStackWidget
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
	view.gameController.KillFun()
	view.radioController.KillFun()
	view.gameController = nil
	view.radioController = nil
	quickio.EmptyTmpFolder()
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
	view.activeGameDataUpdateMap = map[string]bool{}
	return &view
}
