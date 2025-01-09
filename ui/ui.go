package ui

import (
	"os"
	"path/filepath"
	"quickRadio/game"
	"quickRadio/models"
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

func CreateTeamRadioStreamButton(teamAbbrev string, radioLink string) *widgets.QPushButton {
	teamIcon := GetTeamIcon(teamAbbrev)
	button := widgets.NewQPushButton(nil)
	button.SetIcon(teamIcon)
	button.SetFixedHeight(320)
	button.SetFixedWidth(320)
	button.SetIconSize(button.FrameSize())
	if radioLink == "" {
		radioLink = "Game Over/No Link"
	}
	button.ConnectClicked(func(bool) {
		widgets.QMessageBox_Information(nil, "Radio Link for Streaming", radioLink, widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
	})
	return button
}

func CreateIceRinklabel() *widgets.QLabel {
	iceRinkPixmap := GetIceRinkPixmap()
	iceRinkLabel := widgets.NewQLabel2("", nil, core.Qt__Widget)
	iceRinkLabel.SetPixmap(iceRinkPixmap)
	return iceRinkLabel
}

func CreateGameWidget() *widgets.QGroupBox {
	homeTeam := "NHL"
	awayTeam := "NHL"
	layout := widgets.NewQGridLayout(nil)
	gameWidget := widgets.NewQGroupBox(nil)
	homeTeamButton := CreateTeamRadioStreamButton(homeTeam, "TEST RADIO LINK")
	awayTeamButton := CreateTeamRadioStreamButton(awayTeam, "TEST RADIO LINK")
	iceRinkLabel := CreateIceRinklabel()
	layout.AddWidget2(homeTeamButton, 0, 0, core.Qt__AlignLeft)
	layout.AddWidget2(iceRinkLabel, 0, 1, core.Qt__AlignCenter)
	layout.AddWidget2(awayTeamButton, 0, 2, core.Qt__AlignRight)
	gameWidget.SetLayout(layout)
	return gameWidget
}

func CreateDataLabel(data string) *widgets.QLabel {
	label := widgets.NewQLabel2(data, nil, core.Qt__Widget)
	font := label.Font()
	font.SetPointSize(32)
	label.SetFont(font)
	return label
}

func CreateTeamWidgets(teamData struct{}) []*widgets.QWidget {
	return nil
}

func CreateGameWidgetFromGameDataObject(gameDataObject models.GameDataStruct) *widgets.QGroupBox {
	layout := widgets.NewQGridLayout(nil)
	gameWidget := widgets.NewQGroupBox(nil)
	homeTeamButton := CreateTeamRadioStreamButton(gameDataObject.HomeTeam.Abbrev, gameDataObject.HomeTeam.RadioLink)
	homeTeamLabel := CreateDataLabel(gameDataObject.HomeTeam.Abbrev)
	homeTeamScore := CreateDataLabel(strconv.Itoa(gameDataObject.HomeTeam.Score))
	awayTeamButton := CreateTeamRadioStreamButton(gameDataObject.AwayTeam.Abbrev, gameDataObject.AwayTeam.RadioLink)
	awayTeamLabel := CreateDataLabel(gameDataObject.AwayTeam.Abbrev)
	awayTeamScore := CreateDataLabel(strconv.Itoa(gameDataObject.AwayTeam.Score))
	gameStateLabel := CreateDataLabel(gameDataObject.GameState)
	iceRinkLabel := CreateIceRinklabel()
	layout.AddWidget2(homeTeamButton, 0, 0, core.Qt__AlignLeft)
	layout.AddWidget2(homeTeamLabel, 1, 0, core.Qt__AlignCenter)
	layout.AddWidget2(homeTeamScore, 2, 0, core.Qt__AlignCenter)
	layout.AddWidget2(iceRinkLabel, 0, 1, core.Qt__AlignCenter)
	layout.AddWidget2(gameStateLabel, 1, 1, core.Qt__AlignCenter)
	layout.AddWidget2(awayTeamButton, 0, 2, core.Qt__AlignRight)
	layout.AddWidget2(awayTeamLabel, 1, 2, core.Qt__AlignCenter)
	layout.AddWidget2(awayTeamScore, 2, 2, core.Qt__AlignCenter)

	gameWidget.SetLayout(layout)
	//Here we wnt to create the GameDetails widget which will be updated every so often
	//Lastly we want to make sure our audio player 1. works and that the buttons are set up.
	return gameWidget
}

func CreateGameWidgetFromLandinglink(landingLink string) *widgets.QGroupBox {
	gameDataObject := game.GetGameDataObjectFromResponse(landingLink)
	layout := widgets.NewQGridLayout(nil)
	gameWidget := widgets.NewQGroupBox(nil)
	homeTeamButton := CreateTeamRadioStreamButton(gameDataObject.HomeTeam.Abbrev, gameDataObject.HomeTeam.RadioLink)
	awayTeamButton := CreateTeamRadioStreamButton(gameDataObject.AwayTeam.Abbrev, gameDataObject.AwayTeam.RadioLink)
	iceRinkLabel := CreateIceRinklabel()
	layout.AddWidget2(homeTeamButton, 0, 0, core.Qt__AlignLeft)
	layout.AddWidget2(iceRinkLabel, 0, 1, core.Qt__AlignCenter)
	layout.AddWidget2(awayTeamButton, 0, 2, core.Qt__AlignRight)
	gameWidget.SetLayout(layout)
	//Here we wnt to create the GameDetails widget which will be updated every so often
	//Lastly we want to make sure our audio player 1. works and that the buttons are set up.
	return gameWidget
}

func CreateGamesWidget(gameDataObjects []models.GameDataStruct) *widgets.QStackedWidget {
	gameStackWidget := widgets.NewQStackedWidget(nil)
	for _, gameDataObject := range gameDataObjects {
		gameWidget := CreateGameWidgetFromGameDataObject(gameDataObject)
		gameStackWidget.AddWidget(gameWidget)
	}
	gameStackWidget.SetCurrentIndex(0)
	return gameStackWidget
}

func CreateGameDetailsWidget() *widgets.QGroupBox {
	return nil
}

func CreateGameDropdownsWidget(gameDataObjects []models.GameDataStruct, gamesStack *widgets.QStackedWidget) *widgets.QComboBox {
	var gameNames []string
	for _, gameDataObject := range gameDataObjects {
		gameNames = append(gameNames, gameDataObject.HomeTeam.Abbrev+" vs "+gameDataObject.AwayTeam.Abbrev)
	}
	dropdown := widgets.NewQComboBox(nil)
	dropdown.SetFixedWidth(600)
	dropdown.AddItems(gameNames)
	dropdown.ConnectCurrentIndexChanged(func(index int) {
		gamesStack.SetCurrentIndex(index)
	})
	return dropdown
}

func CreateGameManagerWidget(gameDataObjects []models.GameDataStruct) *widgets.QGroupBox {
	gameManager := widgets.NewQGroupBox(nil)
	topbarLayout := widgets.NewQVBoxLayout()
	gameStackLayout := widgets.NewQStackedLayout()
	gameManagerLayout := widgets.NewQVBoxLayout2(gameManager)

	gameManagerLayout.AddLayout(topbarLayout, 1)
	gameManagerLayout.AddLayout(gameStackLayout, 1)
	gamesStack := CreateGamesWidget(gameDataObjects)
	gameDropdown := CreateGameDropdownsWidget(gameDataObjects, gamesStack)

	topbarLayout.AddWidget(gameDropdown, 1, core.Qt__AlignAbsolute)
	gameStackLayout.AddWidget(gamesStack)

	return gameManager
}

func CreateLoadingScreen() *widgets.QSplashScreen {
	pixmap := GetTeamPixmap("NHL")
	splash := widgets.NewQSplashScreen(pixmap, core.Qt__Widget)
	return splash
}

func CreateAndRunUI() {
	app := widgets.NewQApplication(len(os.Args), os.Args)
	loadingScreen := CreateLoadingScreen()
	loadingScreen.Show()
	gameDataObjects := game.UIGetGameDataObjects()
	gameManager := CreateGameManagerWidget(gameDataObjects)
	loadingScreen.Finish(nil)
	window := widgets.NewQMainWindow(nil, 0)
	window.SetCentralWidget(gameManager)
	window.Show()
	app.Exec()
}
