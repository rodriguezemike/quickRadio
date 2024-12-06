package main

import (
	"os"
	"quickRadio/ui"

	"github.com/therecipe/qt/widgets"
)

func main() {
	/*
		linksMap := internals.GetLinksJson()
		gamecenterBase := fmt.Sprintf("%v", linksMap["gamecenter_api_base"])
		gamecenterLanding := fmt.Sprintf("%v", linksMap["gamecenter_api_slug"])
		team := fmt.Sprintf("%v", linksMap["team_abbrev"])
		gameRegexs := []string{fmt.Sprintf("%v", linksMap["home_game_regex"]), fmt.Sprintf("%v", linksMap["away_game_regex"])}
		html := GetGameHtml(linksMap)
		landingLink, err := internals.GetGameLandingLink(html, gamecenterBase, gamecenterLanding, gameRegexs)
		internals.ErrorCheck(err)
		gameDataObject := internals.GetGameDataObjectFromResponse(landingLink)
		radioLink, err := internals.GetRadioLink(gameDataObject, team)
		internals.ErrorCheck(err)
		log.Println(radioLink)
	*/
	//awayTeam := "LAK"
	app := widgets.NewQApplication(len(os.Args), os.Args)

	window := widgets.NewQMainWindow(nil, 0)
	gameWidget := ui.CreateGameWidget()
	window.SetCentralWidget(gameWidget)
	window.Show()
	app.Exec()
}
