package views

import (
	"log"
	"quickRadio/controllers"
	"strings"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

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

func (widget *GamestateAndStatsWidget) createStaticDataLabel(name string, data string, parentWidget *widgets.QGroupBox, fontSize int) *widgets.QLabel {
	label := widgets.NewQLabel2(data, parentWidget, core.Qt__Widget)
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
	gamestateLayout.AddWidget(widget.createDynamicDataLabel("GAMESTATE", widget.gameController.GetGamestateString(), widget.gameController.Landinglink, widget.UIWidget, fontSize), 0, core.Qt__AlignCenter)
	return gamestateLayout
}

func (widget *GamestateAndStatsWidget) createTeamGameStatLayout(homeStatData string, categoryName string, awayStatData string, homeHandle bool) *widgets.QVBoxLayout {
	//For this one we will use each of the Team Game Stats strucs and turn them into 3 lines, 1. 2 dynamic labels and 1 static labe. 2. A slider. 3. Two Dynamic labels Or a spacer that holds the idea od two dynamic labels
	fontSize := 12
	gameStatLayout := widgets.NewQVBoxLayout()
	gameStatLabelLayout := widgets.NewQHBoxLayout()
	gameStatLabelLayout.AddWidget(widget.createDynamicDataLabel("homeStat", homeStatData, widget.gameController.Landinglink, widget.UIWidget, fontSize), 0, core.Qt__AlignLeft)
	gameStatLabelLayout.AddWidget(widget.createStaticDataLabel("category", categoryName, widget.UIWidget, fontSize), 0, core.Qt__AlignCenter)
	gameStatLabelLayout.AddWidget(widget.createDynamicDataLabel("awayStat", awayStatData, widget.gameController.Landinglink, widget.UIWidget, fontSize), 0, core.Qt__AlignCenter)
	slider := widgets.NewQSlider(widget.UIWidget)
	slider.SetStyleSheet(CreateSliderStylesheet(*widget.gameController.HomeTeamController.Sweater, *widget.gameController.AwayTeamController.Sweater, homeHandle))
	slider.SetEnabled(false)
	slider.ConnectTimerEvent(func(event *core.QTimerEvent) {
		//Here every second we redraw the slider to update it from the files we have.
		//We do this by checking which team is winning and setting the style sheet accordingly
		//Then repaint
		widget.updateMap[slider.ObjectName()] = true
		homeHandle = true //Update this if some stat is greater than another. Should be a simple > or < cmp, since these are all number.
		slider.SetStyleSheet(CreateSliderStylesheet(*widget.gameController.HomeTeamController.Sweater, *widget.gameController.AwayTeamController.Sweater, homeHandle))
		slider.Repaint()
	})
	slider.StartTimer(widget.LabelTimer, core.Qt__PreciseTimer)
	gameStatLayout.AddLayout(gameStatLabelLayout, 0)
	gameStatLayout.AddWidget(slider, 0, core.Qt__AlignCenter)
	return gameStatLayout
}

func (widget *GamestateAndStatsWidget) createTeamGameStatsLayout() *widgets.QVBoxLayout {
	//Collection of team game states in a Vertical Layout For our test or default, we want to create 3 things,
	//One where the home team is winning, one where its a tie and one where away is winning.
	teamStatsLayout := widgets.NewQVBoxLayout()
	for _, gameStatObject := range widget.gameController.GetTeamGameStatsObjects() {
		if gameStatObject.HomeValue >= gameStatObject.AwayValue { //We will need to figure out the types in a switch for proper compare
			teamStatsLayout.AddLayout(widget.createTeamGameStatLayout(gameStatObject.HomeValue, gameStatObject.Category, gameStatObject.AwayValue, true), 0)
		} else {
			teamStatsLayout.AddLayout(widget.createTeamGameStatLayout(gameStatObject.HomeValue, gameStatObject.Category, gameStatObject.AwayValue, false), 0)
		}
	}
	return teamStatsLayout
}

func (widget *GamestateAndStatsWidget) createGamestateAndStatsWidget() {
	//Create main layout and widget
	gamestateAndStatsLayout := widgets.NewQVBoxLayout()
	gamestateAndStatsWidget := widgets.NewQGroupBox(widget.gameWidget)
	//Create Child layouts
	gamestateLayout := widget.createGamestateLayout()
	teamGameStatsLayout := widget.createTeamGameStatsLayout()
	//Add Child Layouts
	gamestateAndStatsLayout.AddLayout(gamestateLayout, 0)
	gamestateAndStatsLayout.AddLayout(teamGameStatsLayout, 0)
	//Set Size and Stylesheet - Work off a scaling factor - base = 100 (base*1.77)*ScalingFactor and base*scalingFactor ::Scaling Factor is 2. :: 1.77 is Desired Aspect Ratio.
	gamestateAndStatsWidget.SetLayout(gamestateAndStatsLayout)
	//Set widget UI
	widget.UILayout = gamestateAndStatsLayout
	widget.UIWidget = gamestateAndStatsWidget
}

func CreateNewGamestateAndStatsWidget(labelTimer int, gameIndex int, IsActive bool, controller *controllers.GameController, gameWidget *widgets.QGroupBox) *GamestateAndStatsWidget {
	widget := GamestateAndStatsWidget{}
	widget.updateMap = map[string]bool{}
	widget.LabelTimer = labelTimer
	widget.gameController = controller
	widget.gameWidget = gameWidget
	widget.createGamestateAndStatsWidget()
	return &widget
}
