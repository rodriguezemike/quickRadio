package views

import (
	"quickRadio/models"

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
	painter.DrawPixmap11((iceRinkPixmap.Size().Width()/2)-30, (iceRinkPixmap.Size().Height()/2)-35, 400, 400, homeTeamPixmap) //Hold idea
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
