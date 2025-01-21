package ui

import (
	"encoding/json"
	"os"
	"path/filepath"
	"quickRadio/audio"
	"quickRadio/game"
	"quickRadio/models"
	"quickRadio/radioErrors"
	"runtime"
	"strconv"
	"strings"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

func GetProjectDir() string {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filepath.Dir(filename))
	return dir
}

func GetLogoPath(logoFilename string) string {
	dir := GetProjectDir()
	path := filepath.Join(dir, "assets", "svgs", "logos", logoFilename)
	return path
}

func GetIceRinkPixmap() *gui.QPixmap {
	dir := GetProjectDir()
	path := filepath.Join(dir, "assets", "svgs", "rink.svg")
	rinkPixmap := gui.NewQPixmap3(path, "svg", core.Qt__AutoColor)
	return rinkPixmap
}

func GetTeamPixmap(teamAbbrev string) *gui.QPixmap {
	teamStringSlice := []string{teamAbbrev, "light.svg"}
	teamLogoFilename := strings.Join(teamStringSlice, "_")
	teamLogoPath := GetLogoPath(teamLogoFilename)
	teamPixmap := gui.NewQPixmap3(teamLogoPath, "svg", core.Qt__AutoColor)
	return teamPixmap
}

func GetTeamIcon(teamAbbrev string) *gui.QIcon {
	teamPixmap := GetTeamPixmap(teamAbbrev)
	teamIcon := gui.NewQIcon2(teamPixmap)
	return teamIcon
}

func CreateTeamRadioStreamButton(teamAbbrev string, radioLink string, sweaterColors map[string][]string, gameWidget *widgets.QGroupBox) *widgets.QPushButton {
	teamIcon := GetTeamIcon(teamAbbrev)
	button := widgets.NewQPushButton(gameWidget)
	button.SetStyleSheet(CreateTeamBackgroundStylesheet(teamAbbrev, sweaterColors))
	button.SetObjectName("radio_" + teamAbbrev)
	button.SetIcon(teamIcon)
	button.SetFixedHeight(320)
	button.SetFixedWidth(320)
	button.SetIconSize(button.FrameSize())
	button.SetCheckable(true)
	button.ConnectToggled(func(onCheck bool) {
		if onCheck {
			go audio.StartFun(radioLink)
		} else {
			go audio.StopRadio()
		}
	})
	return button
}

func CreateIceRinklabel(gameWidget *widgets.QGroupBox, homeTeamAbbrev models.TeamData) *widgets.QLabel {
	iceRinkPixmap := GetIceRinkPixmap()
	homeTeamPixmap := GetTeamPixmap(homeTeamAbbrev.Abbrev)
	compositePixmap := gui.NewQPixmap2(iceRinkPixmap.Size())
	compositePixmap.Fill(gui.QColor_FromRgba(0))
	painter := gui.NewQPainter2(compositePixmap)
	painter.DrawPixmap9(0, 0, iceRinkPixmap)
	painter.DrawPixmap11((iceRinkPixmap.Size().Width()/2)-30, (iceRinkPixmap.Size().Height()/2)-35, 64, 64, homeTeamPixmap)
	painter.SetOpacity(.7)
	painter.DrawPixmap9(0, 0, iceRinkPixmap)
	iceRinkLabel := widgets.NewQLabel2("", gameWidget, core.Qt__Widget)
	iceRinkLabel.SetPixmap(compositePixmap)
	return iceRinkLabel
}

func CreateDataLabel(data string, gameWidget *widgets.QGroupBox) *widgets.QLabel {
	label := widgets.NewQLabel2(data, gameWidget, core.Qt__Widget)
	font := label.Font()
	font.SetPointSize(32)
	label.SetFont(font)
	label.SetStyleSheet(CreateLabelStylesheet())
	return label
}

func CreateTeamWidget(team models.TeamData, sweaterColors map[string][]string, gameWidget *widgets.QGroupBox) *widgets.QGroupBox {
	teamLayout := widgets.NewQVBoxLayout2(nil)
	teamWidget := widgets.NewQGroupBox(gameWidget)
	teamLayout.AddWidget(CreateTeamRadioStreamButton(team.Abbrev, team.RadioLink, sweaterColors, gameWidget), 0, core.Qt__AlignCenter)
	teamLayout.AddWidget(CreateDataLabel(team.Abbrev, gameWidget), 0, core.Qt__AlignCenter)
	teamLayout.AddWidget(CreateDataLabel(strconv.Itoa(team.Score), gameWidget), 0, core.Qt__AlignCenter)
	teamWidget.SetLayout(teamLayout)
	return teamWidget
}

func CreateIceCenterWidget(gameDataObject models.GameData, gameWidget *widgets.QGroupBox) *widgets.QGroupBox {
	layout := widgets.NewQVBoxLayout2(nil)
	centerIceWidget := widgets.NewQGroupBox(gameWidget)
	layout.AddWidget(CreateIceRinklabel(gameWidget, gameDataObject.HomeTeam), 0, core.Qt__AlignCenter)
	layout.AddWidget(CreateDataLabel(gameDataObject.GameState+" - "+"P"+strconv.Itoa(gameDataObject.PeriodDescriptor.Number)+" "+gameDataObject.Clock.TimeRemaining, gameWidget),
		0, core.Qt__AlignCenter)
	centerIceWidget.SetLayout(layout)
	return centerIceWidget
}

func CreateGameWidgetFromGameDataObject(gameDataObject models.GameData, sweaterColors map[string][]string, gameStackWidget *widgets.QStackedWidget) *widgets.QGroupBox {
	layout := widgets.NewQGridLayout(nil)
	gameWidget := widgets.NewQGroupBox(gameStackWidget)
	layout.AddWidget2(CreateTeamWidget(gameDataObject.HomeTeam, sweaterColors, gameWidget), 0, 0, core.Qt__AlignCenter)
	layout.AddWidget2(CreateIceCenterWidget(gameDataObject, gameWidget), 0, 1, core.Qt__AlignCenter)
	layout.AddWidget2(CreateTeamWidget(gameDataObject.AwayTeam, sweaterColors, gameWidget), 0, 2, core.Qt__AlignCenter)
	layout.AddWidget2(CreateGameDetailsWidgetFromGameDataObject(gameDataObject, gameWidget), 3, 1, core.Qt__AlignCenter)
	gameWidget.SetLayout(layout)
	gameWidget.SetStyleSheet(CreateGameStylesheet(gameDataObject.HomeTeam.Abbrev, gameDataObject.AwayTeam.Abbrev, sweaterColors))
	//Lastly we want to make sure our audio player 1. works and that the buttons are set up.
	return gameWidget
}

func CreateGamesWidget(gameDataObjects []models.GameData, sweaterColors map[string][]string, window *widgets.QMainWindow, gameManager *widgets.QGroupBox) *widgets.QStackedWidget {
	gameStackWidget := widgets.NewQStackedWidget(gameManager)
	for _, gameDataObject := range gameDataObjects {
		gameWidget := CreateGameWidgetFromGameDataObject(gameDataObject, sweaterColors, gameStackWidget)
		gameStackWidget.AddWidget(gameWidget)
	}
	gameStackWidget.SetCurrentIndex(0)
	return gameStackWidget
}

func CreateGameDetailsWidgetFromGameDataObject(gamedataObject models.GameData, gameWidget *widgets.QGroupBox) *widgets.QGroupBox {
	//For this use MarshalIndent just to see the data
	//This will also become part of the Update
	gameDetailsJson, err := json.MarshalIndent(gamedataObject, "", "\t")
	radioErrors.ErrorCheck(err)
	gameDetailsLayout := widgets.NewQGridLayout(nil)
	gameDetailsWidget := widgets.NewQGroupBox(gameWidget)
	scrollableArea := widgets.NewQScrollArea(gameDetailsWidget)
	jsonDumpLabel := widgets.NewQLabel2("Test", scrollableArea, core.Qt__Window)
	jsonDumpLabel.SetWordWrap(true)
	jsonDumpLabel.SetText(string(gameDetailsJson))
	scrollableArea.SetWidgetResizable(true)
	scrollableArea.SetWidget(jsonDumpLabel)
	gameDetailsLayout.AddWidget(scrollableArea)
	gameDetailsWidget.SetLayout(gameDetailsLayout)
	return gameDetailsWidget
}

func CreateGameDropdownsWidget(gameDataObjects []models.GameData, gamesStack *widgets.QStackedWidget, gameManager *widgets.QGroupBox) *widgets.QComboBox {
	var gameNames []string
	for _, gameDataObject := range gameDataObjects {
		gameNames = append(gameNames, gameDataObject.HomeTeam.Abbrev+" vs "+gameDataObject.AwayTeam.Abbrev)
	}
	dropdown := widgets.NewQComboBox(gameManager)
	dropdown.SetStyleSheet(CreateDropdownStyleSheet())
	dropdown.SetFixedWidth(600)
	dropdown.AddItems(gameNames)
	dropdown.ConnectCurrentIndexChanged(func(index int) {
		gamesStack.SetCurrentIndex(index)
	})
	return dropdown
}

func CreateGameManagerWidget(gameDataObjects []models.GameData, sweaterColors map[string][]string, window *widgets.QMainWindow) *widgets.QGroupBox {
	gameManager := widgets.NewQGroupBox(nil)
	topbarLayout := widgets.NewQVBoxLayout()
	gameStackLayout := widgets.NewQStackedLayout()
	gameManagerLayout := widgets.NewQVBoxLayout2(gameManager)

	gameManagerLayout.AddLayout(topbarLayout, 1)
	gameManagerLayout.AddLayout(gameStackLayout, 1)
	gamesStack := CreateGamesWidget(gameDataObjects, sweaterColors, window, gameManager)
	gameDropdown := CreateGameDropdownsWidget(gameDataObjects, gamesStack, gameManager)

	topbarLayout.AddWidget(gameDropdown, 1, core.Qt__AlignAbsolute)
	gameStackLayout.AddWidget(gamesStack)
	gameManager.SetStyleSheet(CreateGameManagerStyleSheet())
	return gameManager
}

func CreateLoadingScreen() *widgets.QSplashScreen {
	pixmap := GetTeamPixmap("NHL")
	splash := widgets.NewQSplashScreen(pixmap, core.Qt__Widget)
	return splash
}

func UpdateUI(gameManager *widgets.QGroupBox) {
	//Update ALl UI elements. Or figure out how to add a Slot to update them every so often.
	//Do this per widget. Its better. use concuurency and all that jazz.
	gameDataObjects := game.UIGetGameDataObjects()
	println(gameDataObjects)
}

func KillFun() {
	audio.KillFun()
}

func CreateAndRunUI() {
	app := widgets.NewQApplication(len(os.Args), os.Args)
	app.SetApplicationDisplayName("QuickRadio")
	app.ConnectAboutToQuit(func() {
		KillFun()
	})
	app.ConnectDestroyQApplication(func() {
		KillFun()
	})
	app.SetWindowIcon(GetTeamIcon("NHLF"))
	loadingScreen := CreateLoadingScreen()
	loadingScreen.Show()
	window := widgets.NewQMainWindow(nil, 0)
	gameDataObjects := game.UIGetGameDataObjects()
	sweaterColors := game.GetSweaterColors()
	gameManager := CreateGameManagerWidget(gameDataObjects, sweaterColors, window)
	window.SetCentralWidget(gameManager)
	loadingScreen.Finish(nil)
	window.Show()
	app.Exec()
}
