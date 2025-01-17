package ui

import (
	"log"
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
	if radioLink == "" {
		radioLink = "Game Over/No Link"
	}
	button.ConnectClicked(func(bool) {
		widgets.QMessageBox_Information(nil, "Radio Link for Streaming", radioLink, widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
	})
	/*
				button.ConnectToggled(func(checked bool) {
		 			if checked == true {
						audio.KillFun()
						//Find A Way to get all buttons of qpush button type in order to disable em all
							for _, someObject := range window.CentralWidget().Children(){
								if strings.Contains(someObject.ObjectName(), "radio") && !strings.Contains(someObject.ObjectName(), teamAbbrev){

								}
							}
					} else {
						audio.StartFun(radioLink)
							for _, someButton := range window.CentralWidget().FindChildren(widgets.QPushButton) {
								if someButton != button && strings.Contains(someButton.AccessibleName(), "radio") {
									someButton.SetEnabled(false)
								}
							}
					}
				})
	*/
	return button
}

func CreateIceRinklabel(gameWidget *widgets.QGroupBox) *widgets.QLabel {
	iceRinkPixmap := GetIceRinkPixmap()
	iceRinkLabel := widgets.NewQLabel2("", gameWidget, core.Qt__Widget)
	iceRinkLabel.SetPixmap(iceRinkPixmap)
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

func CreateTeamObjects(team models.TeamData, sweaterColors map[string][]string, gameWidget *widgets.QGroupBox) []widgets.QWidget_ITF {
	return []widgets.QWidget_ITF{
		CreateTeamRadioStreamButton(team.Abbrev, team.RadioLink, sweaterColors, gameWidget),
		CreateDataLabel(team.Abbrev, gameWidget),
		CreateDataLabel(strconv.Itoa(team.Score), gameWidget),
	}
}

func CreateGameWidgetFromGameDataObject(gameDataObject models.GameData, sweaterColors map[string][]string, gameStackWidget *widgets.QStackedWidget) *widgets.QGroupBox {
	layout := widgets.NewQGridLayout(nil)
	gameWidget := widgets.NewQGroupBox(gameStackWidget)
	homeTeamObjects := CreateTeamObjects(gameDataObject.HomeTeam, sweaterColors, gameWidget)
	awayTeamObjects := CreateTeamObjects(gameDataObject.AwayTeam, sweaterColors, gameWidget)
	for position, obj := range homeTeamObjects {
		layout.AddWidget2(obj, position, 0, core.Qt__AlignCenter)
	}
	layout.AddWidget2(CreateIceRinklabel(gameWidget), 0, 1, core.Qt__AlignCenter)
	layout.AddWidget2(CreateDataLabel(gameDataObject.GameState, gameWidget), 1, 1, core.Qt__AlignCenter)
	layout.AddWidget2(CreateDataLabel(strconv.Itoa(gameDataObject.PeriodDescriptor.Number)+" "+gameDataObject.Clock.TimeRemaining, gameWidget),
		2, 1, core.Qt__AlignCenter)

	for position, obj := range awayTeamObjects {
		layout.AddWidget2(obj, position, 2, core.Qt__AlignCenter)
	}
	gameWidget.SetLayout(layout)
	gameWidget.SetStyleSheet(CreateGameStylesheet(gameDataObject.HomeTeam.Abbrev, gameDataObject.AwayTeam.Abbrev, sweaterColors))
	log.Println(gameWidget.StyleSheet())
	//Here we wnt to create the GameDetails widget which will be updated every so often
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

func CreateGameDetailsWidget() *widgets.QGroupBox {
	return nil
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

func CreateAndRunUI() {
	app := widgets.NewQApplication(len(os.Args), os.Args)
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
