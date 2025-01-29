package views

import (
	"encoding/json"
	"log"
	"quickRadio/models"
	"quickRadio/quickio"
	"quickRadio/radioErrors"
	"strconv"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

func CreatePlayersOnIceWidgetForTeam(teamOnIce models.TeamOnIce, teamAbbrev string, gamecenterLink string, gameWidget *widgets.QGroupBox) *widgets.QGroupBox {
	//Every 10k Seconds update by replacing all labels in the Widget. with New Data.
	//When we refactor think of a better way to do this, including the many files approach.
	log.Println("GUI::CreatePlayersOnIceWidgetForTeam - Line 204")
	log.Println("GUI::CreatePlayersOnIceWidgetForTeam ", teamAbbrev)
	players := []models.PlayerOnIce{}
	players = append(append(append(players, teamOnIce.Forwards...), teamOnIce.Defensemen...), teamOnIce.Goalies...)
	playersOnIceLayout := widgets.NewQVBoxLayout()
	playersOnIceWidget := widgets.NewQGroupBox(gameWidget)
	playersOnIceWidget.SetAccessibleName(gamecenterLink + " " + teamAbbrev)
	log.Println("GUI::CreatePlayersOnIceWidgetForTeam - Line 210")
	onIceLabel := widgets.NewQLabel2("OnIce", playersOnIceWidget, core.Qt__Widget)
	onIceLabel.SetText("Players On Ice")
	playersOnIceLayout.AddWidget(onIceLabel, 0, core.Qt__AlignCenter)
	log.Println("GUI::CreatePlayersOnIceWidgetForTeam - Line 214")
	for _, player := range players {
		text := strconv.Itoa(player.SweaterNumber) + " " + player.Name.Default + player.PositionCode
		playerLabel := widgets.NewQLabel2(player.Name.Default, playersOnIceWidget, core.Qt__Widget)
		playerLabel.SetText(text)
		playersOnIceLayout.AddWidget(playerLabel, 0, core.Qt__AlignCenter)
		log.Println("GUI::CreatePlayersOnIceWidgetForTeam - Line 220")
	}
	penalityBoxLabel := widgets.NewQLabel2("PenalityBox", playersOnIceWidget, core.Qt__Widget)
	penalityBoxLabel.SetText("Penality Box")
	playersOnIceLayout.AddWidget(penalityBoxLabel, 0, core.Qt__AlignCenter)
	log.Println("GUI::CreatePlayersOnIceWidgetForTeam - Line 225")
	for _, player := range teamOnIce.PenaltyBox {
		text := strconv.Itoa(player.SweaterNumber) + " " + player.Name.Default + player.PositionCode
		playerLabel := widgets.NewQLabel2(player.Name.Default, playersOnIceWidget, core.Qt__Widget)
		playerLabel.SetText(text)
		playersOnIceLayout.AddWidget(playerLabel, 0, core.Qt__AlignCenter)
		log.Println("GUI::CreatePlayersOnIceWidgetForTeam - Line 231")
	}
	playersOnIceWidget.SetLayout(playersOnIceLayout)
	log.Println("GUI::CreatePlayersOnIceWidgetForTeam - Line 235")
	return playersOnIceWidget
}

func CreateGameDetailsWidgetFromGameDataObject(gamedataObject models.GameData, gamecenterLink string, gameWidget *widgets.QGroupBox) *widgets.QGroupBox {
	gameDetailsJson, err := json.MarshalIndent(gamedataObject, "", "\t")
	radioErrors.ErrorCheck(err)
	gameDetailsLayout := widgets.NewQGridLayout(gameWidget)
	gameDetailsWidget := widgets.NewQGroupBox(gameWidget)
	scrollableArea := widgets.NewQScrollArea(gameDetailsWidget)
	jsonDumpLabel := widgets.NewQLabel2("Test", scrollableArea, core.Qt__Window)
	jsonDumpLabel.SetWordWrap(true)
	jsonDumpLabel.SetText(string(gameDetailsJson))
	jsonDumpLabel.ConnectTimerEvent(func(event *core.QTimerEvent) {
		if jsonDumpLabel.IsVisible() {
			gdo := quickio.GetGameDataObject(gamecenterLink)
			gameDetailsJson, _ := json.MarshalIndent(gdo, "", " ")
			jsonDumpLabel.SetText(string(gameDetailsJson))
			jsonDumpLabel.Repaint()
		}
	})
	jsonDumpLabel.StartTimer(30000, core.Qt__VeryCoarseTimer)
	scrollableArea.SetWidgetResizable(true)
	scrollableArea.SetWidget(jsonDumpLabel)
	homeTeamPlayersOnIce := CreatePlayersOnIceWidgetForTeam(gamedataObject.Summary.IceSurface.HomeTeam, gamedataObject.HomeTeam.Abbrev, gamecenterLink, gameWidget)
	//awayTeamPlayersOnIce := CreatePlayersOnIceWidgetForTeam(gamedataObject.Summary.IceSurface.AwayTeam, gamedataObject.AwayTeam.Abbrev, gamecenterLink, gameWidget)
	log.Println("GUI::CreateGameDetailsWidgetFromGameDataObject line 261")
	gameDetailsLayout.AddWidget(homeTeamPlayersOnIce)
	log.Println("GUI::CreateGameDetailsWidgetFromGameDataObject line 263")
	gameDetailsLayout.AddWidget(gameDetailsWidget)
	//gameDetailsLayout.AddWidget(awayTeamPlayersOnIce)
	gameDetailsWidget.SetLayout(gameDetailsLayout)
	log.Println("GUI::CreateGameDetailsWidgetFromGameDataObject line 265")
	return gameDetailsWidget
}
