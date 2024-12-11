package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"quickRadio/game"
	"runtime"
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

func CreateTeamRadioStreamButton(teamAbbrev string) *widgets.QPushButton {
	teamIcon := GetTeamIcon(teamAbbrev)
	button := widgets.NewQPushButton(nil)
	button.SetIcon(teamIcon)
	button.SetFixedHeight(320)
	button.SetFixedWidth(320)
	button.SetIconSize(button.FrameSize())
	button.ConnectClicked(func(bool) {
		linksMap := game.GetLinksJson()
		html := game.GetGameHtml(linksMap)
		gamecenterBase := fmt.Sprintf("%v", linksMap["gamecenter_api_base"])
		gamecenterLanding := fmt.Sprintf("%v", linksMap["gamecenter_api_slug"])
		gameRegexs := []string{fmt.Sprintf("%v", linksMap["home_game_regex"]), fmt.Sprintf("%v", linksMap["away_game_regex"])}
		landingLink, err := game.GetGameLandingLink(html, gamecenterBase, gamecenterLanding, gameRegexs, teamAbbrev)
		if err != nil {
			widgets.QMessageBox_Information(nil, "error", err.Error(), widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
		} else {
			widgets.QMessageBox_Information(nil, "title", landingLink, widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
		}
	})
	return button
}

func CreateGamelandingButton() *widgets.QPushButton {
	button := widgets.NewQPushButton2("Game Landing", nil)
	button.ConnectClicked(func(bool) {
		linksMap := game.GetLinksJson()
		html := game.GetGameHtml(linksMap)
		gamecenterBase := fmt.Sprintf("%v", linksMap["gamecenter_api_base"])
		gamecenterLanding := fmt.Sprintf("%v", linksMap["gamecenter_api_slug"])
		gameRegexs := []string{fmt.Sprintf("%v", linksMap["home_game_regex"]), fmt.Sprintf("%v", linksMap["away_game_regex"])}
		landingLink, err := game.GetGameLandingLink(html, gamecenterBase, gamecenterLanding, gameRegexs, fmt.Sprintf("%v", linksMap["default_team_abbrev"]))
		if err != nil {
			widgets.QMessageBox_Information(nil, "error", err.Error(), widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
		} else {
			widgets.QMessageBox_Information(nil, "title", landingLink, widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
		}
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
	gameLandingButton := CreateGamelandingButton()
	homeTeamButton := CreateTeamRadioStreamButton(homeTeam)
	awayTeamButton := CreateTeamRadioStreamButton(awayTeam)
	iceRinkLabel := CreateIceRinklabel()
	layout.AddWidget2(homeTeamButton, 0, 0, core.Qt__AlignLeft)
	layout.AddWidget2(iceRinkLabel, 0, 1, core.Qt__AlignCenter)
	layout.AddWidget2(awayTeamButton, 0, 2, core.Qt__AlignRight)
	layout.AddWidget2(gameLandingButton, 1, 1, core.Qt__AlignCenter)
	gameWidget.SetLayout(layout)
	return gameWidget
}

func CreateGameWidgetFromLandinglink(landingLink string) *widgets.QGroupBox {
	//toDo - Pull the team abbrev from the link and capitalize
	homeTeam := "NHL"
	awayTeam := "NHL"
	layout := widgets.NewQGridLayout(nil)
	gameWidget := widgets.NewQGroupBox(nil)
	gameLandingButton := CreateGamelandingButton()
	homeTeamButton := CreateTeamRadioStreamButton(homeTeam)
	awayTeamButton := CreateTeamRadioStreamButton(awayTeam)
	iceRinkLabel := CreateIceRinklabel()
	layout.AddWidget2(homeTeamButton, 0, 0, core.Qt__AlignLeft)
	layout.AddWidget2(iceRinkLabel, 0, 1, core.Qt__AlignCenter)
	layout.AddWidget2(awayTeamButton, 0, 2, core.Qt__AlignRight)
	layout.AddWidget2(gameLandingButton, 1, 1, core.Qt__AlignCenter)
	gameWidget.SetLayout(layout)
	//Here we wnt to create the GameDetails widget which will be updated every so often
	//Lastly we want to make sure our audio player 1. works and that the buttons are set up.
	return gameWidget
}

func CreateGamesWidget(landingLinks []string) *widgets.QStackedLayout {
	gameStackLayout := widgets.NewQStackedLayout()
	for _, landingLink := range landingLinks {
		gameWidget := CreateGameWidgetFromLandinglink(landingLink)
		gameStackLayout.AddWidget(gameWidget)
	}
	gameStackLayout.SetCurrentIndex(0)
	return gameStackLayout
}

func CreateGameDetailsWidget() *widgets.QGroupBox {
	return nil
}

func CreateGameDropdowns(landingLinks []string) *widgets.QVBoxLayout {
	dropdown := widgets.NewQComboBox(nil)
	dropdown.AddItems(landingLinks)
	vboxLayout := widgets.NewQVBoxLayout()
	vboxLayout.AddWidget(dropdown, 0, core.Qt__AlignCenter)
	//Here we Need to set the dropbox to switch the stacked widgets and call an update for fresh game data
	return vboxLayout
}

func CreateGameManagerWidget(landingLinks []string) *widgets.QGroupBox {
	gameDropdown := CreateGameDropdowns(landingLinks)
	games := CreateGamesWidget(landingLinks)
	gameManager := widgets.NewQGroupBox(nil)
	gameManagerLayout := widgets.NewQVBoxLayout()
	gameManagerLayout.AddChildLayout(gameDropdown)
	gameManagerLayout.AddChildLayout(games)
	gameManager.SetLayout(gameManagerLayout)
	return gameManager
}

func CreateLoadingScreen() *widgets.QSplashScreen {
	pixmap := GetTeamPixmap("LAK")
	splash := widgets.NewQSplashScreen(pixmap, core.Qt__Widget)
	return splash
}

func CreateAndRunUI() {
	app := widgets.NewQApplication(len(os.Args), os.Args)
	loadingScreen := CreateLoadingScreen()
	loadingScreen.Show()
	landingLinks := game.UIGetGameLandingLinks()
	gameManager := CreateGameManagerWidget(landingLinks)
	loadingScreen.Finish(nil)
	window := widgets.NewQMainWindow(nil, 0)
	window.SetCentralWidget(gameManager)
	window.Show()
	app.Exec()
	//Here we want to do all the polling and under the hood stuff to update whatever widet we have.
}
