package views

import (
	"log"
	"quickRadio/controllers"
	"quickRadio/models"
	"quickRadio/radioErrors"
	"strconv"
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

func (widget *GamestateAndStatsWidget) setDynamicUIObjectName(prefix string, suffix string, delimiter string) string {
	return prefix + delimiter + suffix
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
				} else {
					label.SetText(widget.gameController.GetTextFromObjectName(label.ObjectName()))
				}
				label.Repaint()
				widget.updateMap[label.ObjectName()] = true
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

func (widget *GamestateAndStatsWidget) setSliderValues(statMaxDatum int, homeStatDatum int, awayStatDatum int, homeHandle bool, slider *widgets.QSlider) {
	if homeHandle {
		slider.SetRange(-1*statMaxDatum, 0)
		slider.SetValue(-1 * homeStatDatum)
	} else {
		slider.SetRange(0, statMaxDatum)
		slider.SetValue(awayStatDatum)
	}
	slider.SetTickInterval(1)
	slider.SetTickPosition(widgets.QSlider__NoTicks)
}

func (widget *GamestateAndStatsWidget) createTeamGameStatLayout(homeStatDatum int, categoryName string, awayStatDatum int, statMaxDatum int, homeHandle bool) *widgets.QVBoxLayout {
	//Horizontal and vertical 2 row layout, 1 for labels and 1 for slider.
	fontSize := 12
	gameStatLayout := widgets.NewQVBoxLayout()
	gameStatLabelLayout := widgets.NewQHBoxLayout()
	//Home stat dynamic label
	gameStatLabelLayout.AddWidget(
		widget.createDynamicDataLabel(
			(widget.setDynamicUIObjectName(models.DEFAULT_HOME_PREFIX, categoryName, models.NAME_DELIMITER)),
			strconv.Itoa(homeStatDatum), widget.gameController.Landinglink, widget.UIWidget, fontSize),
		0, core.Qt__AlignLeft,
	)
	//Stat categpru static label
	gameStatLabelLayout.AddWidget(widget.createStaticDataLabel("category", categoryName, widget.UIWidget, fontSize), 0, core.Qt__AlignCenter)
	//Away stat dynamic label
	gameStatLabelLayout.AddWidget(
		widget.createDynamicDataLabel(
			(widget.setDynamicUIObjectName(models.DEFAULT_AWAY_PREFIX, categoryName, models.NAME_DELIMITER)),
			strconv.Itoa(awayStatDatum), widget.gameController.Landinglink, widget.UIWidget, fontSize),
		0, core.Qt__AlignCenter)
	//Stat slider
	slider := widgets.NewQSlider2(core.Qt__Horizontal, widget.UIWidget)
	slider.SetObjectName(widget.setDynamicUIObjectName("slider", categoryName, models.NAME_DELIMITER))
	widget.setSliderValues(statMaxDatum, homeStatDatum, awayStatDatum, homeHandle, slider)
	slider.SetStyleSheet(CreateSliderStylesheet(*widget.gameController.HomeTeamController.Sweater, *widget.gameController.AwayTeamController.Sweater, homeHandle))
	slider.SetEnabled(false)
	slider.ConnectTimerEvent(func(event *core.QTimerEvent) {
		val, ok := widget.updateMap[slider.ObjectName()]
		if !val || !ok {
			homeStatDatum, awayStatDatum, statMaxDatum, homeHandle = widget.gameController.GetGameStatFromFilepath(categoryName)
			widget.setSliderValues(statMaxDatum, homeStatDatum, awayStatDatum, homeHandle, slider)
			slider.SetStyleSheet(CreateSliderStylesheet(*widget.gameController.HomeTeamController.Sweater, *widget.gameController.AwayTeamController.Sweater, homeHandle))
			slider.Repaint()
			widget.updateMap[slider.ObjectName()] = true
		}
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
		awayValue, err := strconv.Atoi(gameStatObject.AwayValue)
		radioErrors.ErrorLog(err)
		homeValue, err := strconv.Atoi(gameStatObject.HomeValue)
		radioErrors.ErrorLog(err)
		maxValue := homeValue + awayValue
		if gameStatObject.HomeValue >= gameStatObject.AwayValue { //We will need to figure out the types in a switch for proper compare
			teamStatsLayout.AddLayout(widget.createTeamGameStatLayout(homeValue, gameStatObject.Category, awayValue, maxValue, true), 0)
		} else {
			teamStatsLayout.AddLayout(widget.createTeamGameStatLayout(homeValue, gameStatObject.Category, awayValue, maxValue, false), 0)
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
