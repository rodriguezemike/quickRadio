//go:build ignore
// +build ignore

package views

import (
	"log"
	"quickRadio/models"
	"strconv"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

func CreateIceRinklabel(gameWidget *widgets.QGroupBox, gameDataObject models.GameData, view *QuickRadioView) *widgets.QLabel {
	homeTeamOnLeft := false
	awayTeam := gameDataObject.AwayTeam
	homeTeam := gameDataObject.HomeTeam
	//Heurstic hometeam starts on the Right.
	if gameDataObject.PeriodDescriptor.Number == 2 {
		homeTeamOnLeft = true
	}
	iceRinkPixmap := view.getIceRinkPixmap()
	homeTeamPixmap := view.getTeamPixmap(homeTeam.Abbrev)
	compositePixmap := gui.NewQPixmap2(iceRinkPixmap.Size())
	compositePixmap.Fill(gui.QColor_FromRgba(0))
	painter := gui.NewQPainter2(compositePixmap)
	painter.DrawPixmap9(0, 0, iceRinkPixmap)
	//painter.DrawPixmap11((iceRinkPixmap.Size().Width()/2)-30, (iceRinkPixmap.Size().Height()/2)-35, 64, 64, homeTeamPixmap)
	painter.DrawPixmap11((iceRinkPixmap.Size().Width()/2)-200, (iceRinkPixmap.Size().Height()/2)-200, 400, 400, homeTeamPixmap) //Hold idea
	//This draws our team sides. Look into 'tinting' the Ice different gradient colors for the team. Red line dividing them.
	if homeTeamOnLeft {
		painter.DrawText3((iceRinkPixmap.Size().Width() / 4), (iceRinkPixmap.Size().Height()/2)-35, homeTeam.Abbrev)
		painter.DrawText3((iceRinkPixmap.Size().Width()/4)*3, (iceRinkPixmap.Size().Height()/2)-35, awayTeam.Abbrev)
	} else {
		painter.DrawText3((iceRinkPixmap.Size().Width() / 4), (iceRinkPixmap.Size().Height()/2)-35, awayTeam.Abbrev)
		painter.DrawText3((iceRinkPixmap.Size().Width()/4)*3, (iceRinkPixmap.Size().Height()/2)-35, homeTeam.Abbrev)
	}
	//Here we want to draw our players.
	//Here we would wanna tint ice.
	//painter.SetOpacity(.5)
	//Draw Gradient
	//Finally we want to make the last ice layer
	painter.SetOpacity(.7)
	painter.DrawPixmap9(0, 0, iceRinkPixmap)
	iceRinkLabel := widgets.NewQLabel2("", gameWidget, core.Qt__Widget)
	iceRinkLabel.SetPixmap(compositePixmap)
	//Look into update func and edge data, if not, maybe load some sort of edge data into the middle for display?
	//This should also be the area for replays of goals when we get into video playing.
	//IceRink.go will be a thing.
	return iceRinkLabel
}

func CreatePlayersOnIceWidgetForTeam(teamOnIce models.TeamOnIce, teamAbbrev string, gamecenterLink string, gameWidget *widgets.QGroupBox) *widgets.QGroupBox {
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
