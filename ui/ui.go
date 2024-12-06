package ui

import (
	"fmt"
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
	homeTeam := "LAK"
	awayTeam := "LAK"
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
