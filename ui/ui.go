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
	//button.SetCheckable(true)
	if radioLink == "" {
		radioLink = "Game Over/No Link"
	}
	/*
		button.ConnectToggled(func(checked bool) {
			//Here When it is not checked, check it, change the color
			//Then loop through all the buttons disabling em
			//Finally start the radio

			//If it is checked
			//Kill radio, uncheck, change color,
			//Then loop through all the buttons enabling em

		})
	*/
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

func CreateTeamObjects(team models.TeamData) []widgets.QWidget_ITF {
	return []widgets.QWidget_ITF{
		CreateTeamRadioStreamButton(team.Abbrev, team.RadioLink),
		CreateDataLabel(team.Abbrev),
		CreateDataLabel(strconv.Itoa(team.Score)),
	}
}

func CreateGamePalette(homeTeam string, awayTeam string, sweaterColors map[string][]string) *gui.QPalette {
	//Aspecr Ratio and need to figure out window heights
	homeTeamColors := sweaterColors[homeTeam]
	awayTeamColors := sweaterColors[awayTeam]
	//gradient := gui.NewQLinearGradient3(0.0, 0.0, float64(window.Width())*1.77, float64(window.Width())*1.77)
	gradient := gui.NewQLinearGradient3(0.0, 0.0, float64(1920)*1.77, float64(1080)*1.77) //Holds the idea
	gradient.SetColorAt(0.0, gui.NewQColor6(strings.TrimSpace(homeTeamColors[0])))
	gradient.SetColorAt(0.13, gui.NewQColor6(strings.TrimSpace(homeTeamColors[0])))
	gradient.SetColorAt(0.24, gui.NewQColor6(strings.TrimSpace(homeTeamColors[1])))
	gradient.SetColorAt(0.39, gui.NewQColor6(strings.TrimSpace(homeTeamColors[1])))
	gradient.SetColorAt(0.40, gui.NewQColor2(core.Qt__white))
	gradient.SetColorAt(0.5, gui.NewQColor2(core.Qt__white))
	gradient.SetColorAt(0.55, gui.NewQColor2(core.Qt__white))
	gradient.SetColorAt(0.56, gui.NewQColor6(strings.TrimSpace(awayTeamColors[1])))
	gradient.SetColorAt(0.60, gui.NewQColor6(strings.TrimSpace(awayTeamColors[1])))
	gradient.SetColorAt(0.71, gui.NewQColor6(strings.TrimSpace(awayTeamColors[0])))
	gradient.SetColorAt(1.0, gui.NewQColor6(strings.TrimSpace(awayTeamColors[0])))
	gamePalette := gui.NewQPalette()
	gamePalette.SetBrush(gui.QPalette__Window, gui.NewQBrush10(gradient))
	return gamePalette
}

func CreateGameWidgetFromGameDataObject(gameDataObject models.GameData, sweaterColors map[string][]string) *widgets.QGroupBox {
	//Need to Set parent and child heirarchy to attempt to get proper palette working for this widget
	//Might need to figure it out for stacked widget
	layout := widgets.NewQGridLayout(nil)
	gameWidget := widgets.NewQGroupBox(nil)
	gamePalette := CreateGamePalette(gameDataObject.HomeTeam.Abbrev, gameDataObject.AwayTeam.Abbrev, sweaterColors)
	gameWidget.SetPalette(gamePalette)
	homeTeamObjects := CreateTeamObjects(gameDataObject.HomeTeam)
	awayTeamObjects := CreateTeamObjects(gameDataObject.AwayTeam)
	for position, obj := range homeTeamObjects {
		layout.AddWidget2(obj, position, 0, core.Qt__AlignCenter)
	}
	layout.AddWidget2(CreateIceRinklabel(), 0, 1, core.Qt__AlignCenter)
	layout.AddWidget2(CreateDataLabel(gameDataObject.GameState), 1, 1, core.Qt__AlignCenter)
	layout.AddWidget2(CreateDataLabel(strconv.Itoa(gameDataObject.PeriodDescriptor.Number)+" "+gameDataObject.Clock.TimeRemaining),
		2, 1, core.Qt__AlignCenter)

	for position, obj := range awayTeamObjects {
		layout.AddWidget2(obj, position, 2, core.Qt__AlignCenter)
	}
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

func CreateGamesWidget(gameDataObjects []models.GameData, sweaterColors map[string][]string) *widgets.QStackedWidget {
	gameStackWidget := widgets.NewQStackedWidget(nil)
	for _, gameDataObject := range gameDataObjects {
		gameWidget := CreateGameWidgetFromGameDataObject(gameDataObject, sweaterColors)
		gameStackWidget.AddWidget(gameWidget)
	}
	gameStackWidget.SetCurrentIndex(0)
	return gameStackWidget
}

func CreateGameDetailsWidget() *widgets.QGroupBox {
	return nil
}

func CreateGameDropdownsWidget(gameDataObjects []models.GameData, gamesStack *widgets.QStackedWidget) *widgets.QComboBox {
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

func CreateGameManagerWidget(gameDataObjects []models.GameData, sweaterColors map[string][]string) *widgets.QGroupBox {
	gameManager := widgets.NewQGroupBox(nil)
	topbarLayout := widgets.NewQVBoxLayout()
	gameStackLayout := widgets.NewQStackedLayout()
	gameManagerLayout := widgets.NewQVBoxLayout2(gameManager)

	gameManagerLayout.AddLayout(topbarLayout, 1)
	gameManagerLayout.AddLayout(gameStackLayout, 1)
	gamesStack := CreateGamesWidget(gameDataObjects, sweaterColors)
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
	window := widgets.NewQMainWindow(nil, 0)
	gameDataObjects := game.UIGetGameDataObjects()
	sweaterColors := game.GetSweaterColors()
	//ToDo : Pass in Window to all the things
	gameManager := CreateGameManagerWidget(gameDataObjects, sweaterColors)
	window.SetCentralWidget(gameManager)
	loadingScreen.Finish(nil)
	window.Show()
	app.Exec()
}
