package views

import (
	"context"
	"os"
	"quickRadio/quickio"
	"sync"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

type QuickRadioUI struct {
	CentralWidget   *widgets.QGroupBox
	CentralLayout   *widgets.QVBoxLayout
	LandingLinks    []string
	LabelTimer      int
	DropdownWidth   int
	app             *widgets.QApplication
	window          *widgets.QMainWindow
	gameManagerView *GameManagerView
	radioLock       *sync.Mutex
	ctx             context.Context
	cancelCallback  context.CancelFunc
}

func (ui *QuickRadioUI) CreateLoadingScreen() *widgets.QSplashScreen {
	pixmap := GetTeamPixmap("NHL")
	splash := widgets.NewQSplashScreen(pixmap, core.Qt__Widget)
	return splash
}

func (ui *QuickRadioUI) KillAllTheFun() {
	//Here we want to ensure the producer gorountine is canceled
	ui.cancelCallback()
	//Call all views deconstructors - For now just set to nil.
	//Empty all tmp directories
	quickio.EmptyTmpFolder()
	//Exit and set all pointers to nil
	ui.app.CloseAllWindows()
	ui.app.Exit(0)
	//Set to nil and deconstruct any qt objects due to bindings
	ui.app = nil
	ui.window = nil
	ui.CentralLayout = nil
	ui.CentralWidget = nil
	ui.gameManagerView = nil
	ui.radioLock = nil
}

func (ui *QuickRadioUI) CreateAndRunApp() {
	var gameViews []*GameView
	ui.radioLock = &sync.Mutex{}
	quickio.EmptyTmpFolder()
	ui.app = widgets.NewQApplication(len(os.Args), os.Args)
	ui.CentralLayout = widgets.NewQVBoxLayout()
	ui.CentralWidget = widgets.NewQGroupBox(nil)
	ui.app.SetWindowIcon(GetTeamIcon("NHLF"))
	ui.window = widgets.NewQMainWindow(nil, 0)
	loadingScreen := ui.CreateLoadingScreen()
	loadingScreen.Show()
	for _, link := range ui.LandingLinks {
		gameViews = append(gameViews, CreateNewGameView(link, ui.CentralWidget, ui.radioLock, ui.LabelTimer))
	}
	ui.ctx, ui.cancelCallback = context.WithCancel(context.Background())
	ui.gameManagerView = CreateNewGameManagerView(ui.DropdownWidth, gameViews, ui.CentralWidget, ui.LabelTimer)
	ui.app.SetApplicationDisplayName("QuickRadio")
	ui.app.ConnectAboutToQuit(func() {
		ui.KillAllTheFun()
	})
	ui.app.ConnectDestroyQApplication(func() {
		ui.KillAllTheFun()
	})
	ui.CentralLayout.AddWidget(ui.gameManagerView.UIWidget, 0, core.Qt__AlignTop)
	ui.CentralWidget.SetLayout(ui.CentralLayout)
	ui.window.SetCentralWidget(ui.CentralWidget)
	loadingScreen.Finish(nil)
	ui.window.Show()
	go ui.gameManagerView.GoUpdateGames(ui.ctx)
	ui.app.Exec()
}

func CreateNewQuckRadioUI() *QuickRadioUI {
	var ui QuickRadioUI
	ui.LabelTimer = 5000
	ui.DropdownWidth = 600
	ui.LandingLinks = quickio.GetGameLandingLinks()
	return &ui
}
