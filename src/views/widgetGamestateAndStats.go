package views

import (
	"log"
	"quickRadio/controllers"
	"strings"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

//This should hold up to date GameState information like hits, possession time, Face-off %, Power Play %
//penality minutes, hits, blocked shots, giveaways, takaes, if we have room for above put season series info
//Sliders for all of these that you cant interact with. With Dot or shape showing whos winning
//Primary colars for each side of the bar. Do this by setting home collar for bar and then fill up with away color
//Or vice versa for effect.
//Each one will have 3 Horizontal layouts with last one maybe being a spacer if no text is needed.
//Should also time game state, time left, period, Situation and should be able to change colors ased on game state
//Who has a Power play, are we in the crtical game state? Etc.
//This may also be the location for the mini video player for the goal with edget animation below on ice.
//This last one is for next season need to stand up version 2 of ui first.

//Leave On Ice stuff for the Ice widget that will be blow this.

type GamestateAndStatsWidget struct {
	LabelTimer     int
	GameIndex      int
	IsActive       bool
	UIWidget       *widgets.QGroupBox
	UILayout       *widgets.QVBoxLayout
	gameWidget     *widgets.QGroupBox
	gameController *controllers.GameController
	updateMap      map[string]bool
}

func (widget *GamestateAndStatsWidget) createStaticDataLabel(name string, data string, gameWidget *widgets.QGroupBox, fontSize int) *widgets.QLabel {
	label := widgets.NewQLabel2(data, gameWidget, core.Qt__Widget)
	label.SetObjectName(name)
	label.SetProperty("label-type", core.NewQVariant12("static"))
	label.SetStyleSheet(CreateStaticDataLabelStylesheet(fontSize))
	log.Println("Static Label stylesheet", label.StyleSheet())
	return label
}

func (widget *GamestateAndStatsWidget) createDynamicDataLabel(name string, data string, gamecenterLink string, parentWidget *widgets.QGroupBox, fontSize int) *widgets.QLabel {
	label := widgets.NewQLabel2(data, parentWidget, core.Qt__Widget)
	label.SetObjectName(name)
	label.SetAccessibleDescription(gamecenterLink)
	label.SetProperty("label-type", core.NewQVariant12("dynamic"))
	label.SetStyleSheet(CreateDynamicDataLabelStylesheet(fontSize))
	label.ConnectTimerEvent(func(event *core.QTimerEvent) {
		if label.IsVisible() {
			val, ok := widget.updateMap[label.ObjectName()]
			if !ok || !val {
				if strings.Contains(label.ObjectName(), "GAMESTATE") {
					label.SetText(widget.gameController.GetGamestateString())
					label.Repaint()
					widget.updateMap[label.ObjectName()] = true
				}
			}
		}
	})
	label.StartTimer(widget.LabelTimer, core.Qt__PreciseTimer)
	log.Println("Dynamic Label stylesheet", label.StyleSheet())
	return label
}

func (widget *GamestateAndStatsWidget) createGamestateLayout() *widgets.QHBoxLayout {
	fontSize := 32
	gamestateLayout := widgets.NewQHBoxLayout()
	gamestateWidget := widgets.NewQGroupBox(widget.gameWidget)
	gamestateLayout.AddWidget(widget.createDynamicDataLabel("GAMESTATE", widget.gameController.GetGamestateString(), widget.gameController.Landinglink, gamestateWidget, fontSize), 0, core.Qt__AlignCenter)
	gamestateWidget.SetLayout(gamestateLayout)
	return gamestateLayout
}

func (widget *GamestateAndStatsWidget) createTeamGameStatLayout() *widgets.QHBoxLayout {
	//For this one we will use each of the Team Game Stats strucs and turn them into 3 lines, 1. 2 dynamic labels and 1 static labe. 2. A slider. 3. Two Dynamic labels Or a spacer that holds the idea od two dynamic labels
	return widgets.NewQHBoxLayout()
}

func (widget *GamestateAndStatsWidget) createTeamGameStatsLayout() *widgets.QVBoxLayout {
	//Collection of team game states in a Vertical Layout
	return widgets.NewQVBoxLayout()
}

func (widget *GamestateAndStatsWidget) createGamestateAndStatsWidget() {
	//Here for Live Game layout, call the game state layout call each of the game stats and the game state layout and return them in a vertical box We need to match the size to the team layout even if we need to add more room on the team side
	//We may be also be able to allow the grid layout of the GameView to do this for us.
	gamestateAndStatsLayout := widgets.NewQVBoxLayout()
	gamestateAndStatsWidget := widgets.NewQGroupBox(widget.gameWidget)
	//Create Child layouts
	gamestateLayout := widget.createGamestateLayout()
	//teamGameStatsLayout := widget.createTeamGameStatsLayout()
	//Add Child Layouts
	gamestateAndStatsLayout.AddLayout(gamestateLayout, 0)
	//gamestateAndStatsLayout.AddLayout(teamGamestatsLayout, 0)
	//Set Size and Stylesheet - Work off a scaling factor - base = 100 (base*1.77)*ScalingFactor and base*scalingFactor ::Scaling Factor is 2. :: 1.77 is Desired Aspect Ratio.
	gamestateAndStatsWidget.SetLayout(gamestateAndStatsLayout)
	//Set widget UI
	widget.UILayout = gamestateAndStatsLayout
	widget.UIWidget = gamestateAndStatsWidget
}

func CreateNewGamestateAndStatsWidget(labelTimer int, gameIndex int, IsActive bool, controller *controllers.GameController, gameWidget *widgets.QGroupBox) *GamestateAndStatsWidget {
	widget := GamestateAndStatsWidget{}
	widget.LabelTimer = labelTimer
	widget.gameController = controller
	widget.gameWidget = gameWidget
	widget.createGamestateAndStatsWidget()
	return &widget
}
