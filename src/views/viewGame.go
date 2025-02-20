package views

import (
	"github.com/therecipe/qt/widgets"
)

//This a Link and creates two team widgies and 1 gamestate and stats widgie
//Then drops em in a HBOX Layout, this calls produce data in the game controller and calls conume data
//After checking all the update maps for the various widgies
//If we have to we can add the players on ice widget or add em under the team
//Currents players on ICe with a timer for time spent on ice.

type GameView struct {
	LabelTimer              int
	UIWidget                *widgets.QGroupBox
	UILayout                *widgets.QHBoxLayout
	AwayTeamWidget          *TeamWidget
	HomeTeamWidget          *TeamWidget
	GamestateAndStatsWidget *GamestateAndStatsWidget
	parentWidget            *widgets.QGroupBox
}

func (view *GameView) createGameView(gamecenterLink string, parentWidget *widgets.QGroupBox) {
	view.parentWidget = parentWidget
	//Set UI widget and Layout
	viewLayout := widgets.NewQHBoxLayout()
	viewGroupBox := widgets.NewQGroupBox(view.parentWidget)
	view.UILayout = viewLayout
	view.UIWidget = viewGroupBox
	//Create Child layouts
	//Add Child Layouts
	//Set Size and Stylesheet - Work off a scaling factor - base = 100 (base*1.77)*ScalingFactor and base*scalingFactor ::Scaling Factor is 2. :: 1.77 is Desired Aspect Ratio.
}

func CreateNewGameView(gamecenterLink string, parentWidget *widgets.QGroupBox) *GameView {
	//Create team layout, groupbox and set custom properties

	GameView := GameView{}
	GameView.createGameView(gamecenterLink, parentWidget)
	return &GameView
}
