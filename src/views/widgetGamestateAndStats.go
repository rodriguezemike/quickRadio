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
	UIWidget       *widgets.QGroupBox
	UILayout       *widgets.QVBoxLayout
	parentWidget   *widgets.QGroupBox
	gameController *controllers.GameController
	updateMap      map[string]bool
}

func (widget *GamestateAndStatsWidget) ClearUpdateMap() {
	for key := range widget.updateMap {
		widget.updateMap[key] = false
	}
}

func (widget *GamestateAndStatsWidget) IsUpdated() bool {
	for _, v := range widget.updateMap {
		if !v {
			return false
		}
	}
	return true
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
	label.SetWordWrap(true)
	label.ConnectTimerEvent(func(event *core.QTimerEvent) {
		val, ok := widget.updateMap[label.ObjectName()]
		if !ok || !val {
			if strings.Contains(label.ObjectName(), "GAMESTATE") {
				label.SetText(widget.gameController.GetGamestateString())
			} else {
				label.SetText(widget.gameController.GetTextFromObjectNameFilepath(label.ObjectName()))
			}
			label.Repaint()
			widget.updateMap[label.ObjectName()] = true
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

func (widget *GamestateAndStatsWidget) createStatSlider(categoryName string, statMaxDatum int, homeStatDatum int, awayStatDatum int, homeHandle bool) *widgets.QSlider {
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
	return slider
}
func (widget *GamestateAndStatsWidget) createHomeAndAwayStatLabels(categoryName string, homeStatDatum int, awayStatDatum int, fontSize int) (*widgets.QLabel, *widgets.QLabel) {
	homeStatLabel := widget.createDynamicDataLabel(
		(widget.setDynamicUIObjectName(models.DEFAULT_HOME_PREFIX, categoryName, models.NAME_DELIMITER)),
		strconv.Itoa(homeStatDatum), widget.gameController.Landinglink, widget.UIWidget, fontSize)
	awayStatLabel := widget.createDynamicDataLabel(
		(widget.setDynamicUIObjectName(models.DEFAULT_AWAY_PREFIX, categoryName, models.NAME_DELIMITER)),
		strconv.Itoa(awayStatDatum), widget.gameController.Landinglink, widget.UIWidget, fontSize)
	return homeStatLabel, awayStatLabel
}

func (widget *GamestateAndStatsWidget) createTeamGameStatLayout(homeStatDatum int, categoryName string, awayStatDatum int, statMaxDatum int, homeHandle bool) *widgets.QVBoxLayout {
	//Horizontal and vertical 2 row layout, 1 for labels and 1 for slider.
	fontSize := 12
	gameStatLayout := widgets.NewQVBoxLayout()
	gameStatLabelLayout := widgets.NewQHBoxLayout()
	homeStatLabel, awayStatLabel := widget.createHomeAndAwayStatLabels(categoryName, homeStatDatum, awayStatDatum, fontSize)
	//Home stat dynamic label
	gameStatLabelLayout.AddWidget(homeStatLabel, 0, core.Qt__AlignLeft)
	//Stat category static label
	gameStatLabelLayout.AddWidget(widget.createStaticDataLabel("category", categoryName, widget.UIWidget, fontSize), 0, core.Qt__AlignCenter)
	//Away stat dynamic label
	gameStatLabelLayout.AddWidget(awayStatLabel, 0, core.Qt__AlignCenter)
	//Stat slider
	slider := widget.createStatSlider(categoryName, statMaxDatum, homeStatDatum, awayStatDatum, homeHandle)
	//Add label layout and slider widget
	gameStatLayout.AddLayout(gameStatLabelLayout, 0)
	gameStatLayout.AddWidget(slider, 0, core.Qt__AlignCenter)
	return gameStatLayout
}

func (widget *GamestateAndStatsWidget) createTeamGameStatsLayout() *widgets.QVBoxLayout {
	//Collection of team game states in a Vertical Layout For our test or default, we want to create 3 things,
	//One where the home team is winning, one where its a tie and one where away is winning.
	teamStatsLayout := widgets.NewQVBoxLayout()
	for _, gameStatObject := range widget.gameController.GetTeamGameStatsObjects() {
		anyAwayValue, _ := gameStatObject.AwayValue.(string)
		anyHomeValue, _ := gameStatObject.HomeValue.(string)
		awayValue, err := strconv.Atoi(anyAwayValue)
		radioErrors.ErrorLog(err)
		homeValue, err := strconv.Atoi(anyHomeValue)
		radioErrors.ErrorLog(err)
		maxValue := homeValue + awayValue
		if anyHomeValue >= anyAwayValue { //We will need to figure out the types in a switch for proper compare
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
	gamestateAndStatsWidget := widgets.NewQGroupBox(widget.parentWidget)
	gamestateAndStatsWidget.SetProperty("widget-type", core.NewQVariant12("gamestatsAndGamestate"))
	//Create Child layouts
	gamestateLayout := widget.createGamestateLayout()
	//teamGameStatsLayout := widget.createTeamGameStatsLayout() For next Version
	//Add Child Layouts
	gamestateAndStatsLayout.AddLayout(gamestateLayout, 0)
	//gamestateAndStatsLayout.AddLayout(teamGameStatsLayout, 0) For next Version
	//Set Size and Stylesheet - Work off a scaling factor - base = 100 (base*1.77)*ScalingFactor and base*scalingFactor ::Scaling Factor is 2. :: 1.77 is Desired Aspect Ratio.
	gamestateAndStatsWidget.SetLayout(gamestateAndStatsLayout)
	gamestateAndStatsWidget.SetMinimumSize(core.NewQSize2(354, 200))
	gamestateAndStatsWidget.SetMaximumSize(core.NewQSize2(885, 500))
	gamestateAndStatsWidget.SetStyleSheet(CreateGameStatsAndGamestateStylesheet())
	//Set widget UI
	widget.UILayout = gamestateAndStatsLayout
	widget.UIWidget = gamestateAndStatsWidget
}

func CreateNewGamestateAndStatsWidget(labelTimer int, controller *controllers.GameController, gameWidget *widgets.QGroupBox) *GamestateAndStatsWidget {
	widget := GamestateAndStatsWidget{}
	widget.updateMap = map[string]bool{}
	widget.LabelTimer = labelTimer
	widget.gameController = controller
	widget.parentWidget = gameWidget
	widget.createGamestateAndStatsWidget()
	return &widget
}
