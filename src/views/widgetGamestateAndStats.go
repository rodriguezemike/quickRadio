package views

import (
	"log"
	"math"
	"quickRadio/controllers"
	"quickRadio/models"
	"strconv"
	"strings"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

type GamestateAndStatsWidget struct {
	LabelTimer          int
	UIWidget            *widgets.QGroupBox
	UILayout            *widgets.QVBoxLayout
	pregameStatsLayout  *widgets.QVBoxLayout
	liveGameStatsLayout *widgets.QVBoxLayout
	pregameStatsWidget  *widgets.QGroupBox
	liveGameStatsWidget *widgets.QGroupBox
	parentWidget        *widgets.QGroupBox
	gameController      *controllers.GameController
	IsFuture            bool
	updateMap           map[string]bool
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
				label.SetText(widget.gameController.GetTextFromObjectNameFilepath(label.ObjectName(), label.Text()))
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
	if !homeHandle {
		slider.SetRange(-1*statMaxDatum, 0)
		slider.SetValue(-1 * awayStatDatum)
	} else {
		slider.SetRange(0, statMaxDatum)
		slider.SetValue(homeStatDatum)
	}
	slider.SetTickInterval(1)
	slider.SetTickPosition(widgets.QSlider__NoTicks)
}

func (widget *GamestateAndStatsWidget) setSliderValuesFloat(statMaxDatum float64, homeStatDatum float64, awayStatDatum float64, homeHandle bool, slider *widgets.QSlider) {
	if !homeHandle {
		slider.SetRange(-1*int(math.Ceil(statMaxDatum)), 0)
		slider.SetValue(-1 * int(math.Ceil(awayStatDatum)))
	} else {
		slider.SetRange(0, int(math.Ceil(statMaxDatum)))
		slider.SetValue(int(math.Ceil(homeStatDatum)))
	}
	slider.SetTickInterval(1)
	slider.SetTickPosition(widgets.QSlider__NoTicks)
}

func (widget *GamestateAndStatsWidget) createFloatStatSlider(categoryName string, homeStatDatum string, awayStatDatum string, homeHandle bool, dynamic bool) *widgets.QSlider {
	slider := widgets.NewQSlider2(core.Qt__Horizontal, widget.UIWidget)
	slider.SetObjectName(widget.setDynamicUIObjectName("slider", categoryName, models.NAME_DELIMITER))
	homeStat, _ := strconv.ParseFloat(homeStatDatum, 64)
	awayStat, _ := strconv.ParseFloat(awayStatDatum, 64)
	homeStat = homeStat * 100.0
	awayStat = awayStat * 100.0
	maxStat := homeStat + awayStat
	widget.setSliderValuesFloat(maxStat, homeStat, awayStat, homeHandle, slider)
	slider.SetStyleSheet(CreateSliderStylesheet(*widget.gameController.HomeTeamController.Sweater, *widget.gameController.AwayTeamController.Sweater, homeHandle))
	slider.SetEnabled(false)
	if dynamic {
		slider.ConnectTimerEvent(func(event *core.QTimerEvent) {
			val, ok := widget.updateMap[slider.ObjectName()]
			if !val || !ok {
				homeStat, awayStat, maxStat, homeHandle = widget.gameController.GetGameStatFloatsFromFilepath(categoryName)
				homeStat = homeStat * 100.0
				awayStat = awayStat * 199.0
				maxStat = homeStat + awayStat
				widget.setSliderValuesFloat(maxStat, homeStat, awayStat, homeHandle, slider)
				slider.SetStyleSheet(CreateSliderStylesheet(*widget.gameController.HomeTeamController.Sweater, *widget.gameController.AwayTeamController.Sweater, homeHandle))
				slider.Repaint()
				widget.updateMap[slider.ObjectName()] = true
			}
		})
		slider.StartTimer(widget.LabelTimer, core.Qt__PreciseTimer)
	}
	return slider
}

func (widget *GamestateAndStatsWidget) createIntStatSlider(categoryName string, homeStatDatum string, awayStatDatum string, homeHandle bool, dynamic bool) *widgets.QSlider {
	slider := widgets.NewQSlider2(core.Qt__Horizontal, widget.UIWidget)
	slider.SetObjectName(widget.setDynamicUIObjectName("slider", categoryName, models.NAME_DELIMITER))
	homeStat, _ := strconv.Atoi(homeStatDatum)
	awayStat, _ := strconv.Atoi(awayStatDatum)
	maxStat := homeStat + awayStat
	widget.setSliderValues(maxStat, homeStat, awayStat, homeHandle, slider)
	slider.SetStyleSheet(CreateSliderStylesheet(*widget.gameController.HomeTeamController.Sweater, *widget.gameController.AwayTeamController.Sweater, homeHandle))
	slider.SetEnabled(false)
	if dynamic {
		slider.ConnectTimerEvent(func(event *core.QTimerEvent) {
			val, ok := widget.updateMap[slider.ObjectName()]
			if !val || !ok {
				homeStat, awayStat, maxStat, homeHandle = widget.gameController.GetGameStatFromFilepath(categoryName)
				widget.setSliderValues(maxStat, homeStat, awayStat, homeHandle, slider)
				slider.SetStyleSheet(CreateSliderStylesheet(*widget.gameController.HomeTeamController.Sweater, *widget.gameController.AwayTeamController.Sweater, homeHandle))
				slider.Repaint()
				widget.updateMap[slider.ObjectName()] = true
			}
		})
		slider.StartTimer(widget.LabelTimer, core.Qt__PreciseTimer)
	}
	return slider
}

func (widget *GamestateAndStatsWidget) createStaticHomeAndAwayStatLabels(categoryName string, homeStatDatum string, awayStatDatum string, fontSize int) (*widgets.QLabel, *widgets.QLabel) {
	homeStatLabel := widget.createStaticDataLabel(
		(widget.setDynamicUIObjectName(models.DEFAULT_HOME_PREFIX, categoryName, models.NAME_DELIMITER)),
		homeStatDatum, widget.UIWidget, fontSize)
	awayStatLabel := widget.createStaticDataLabel(
		(widget.setDynamicUIObjectName(models.DEFAULT_AWAY_PREFIX, categoryName, models.NAME_DELIMITER)),
		awayStatDatum, widget.UIWidget, fontSize)
	return homeStatLabel, awayStatLabel
}

func (widget *GamestateAndStatsWidget) createHomeAndAwayStatLabels(categoryName string, homeStatDatum string, awayStatDatum string, fontSize int) (*widgets.QLabel, *widgets.QLabel) {
	homeStatLabel := widget.createDynamicDataLabel(
		(widget.setDynamicUIObjectName(models.DEFAULT_HOME_PREFIX, categoryName, models.NAME_DELIMITER)),
		homeStatDatum, widget.gameController.Landinglink, widget.UIWidget, fontSize)
	awayStatLabel := widget.createDynamicDataLabel(
		(widget.setDynamicUIObjectName(models.DEFAULT_AWAY_PREFIX, categoryName, models.NAME_DELIMITER)),
		awayStatDatum, widget.gameController.Landinglink, widget.UIWidget, fontSize)
	return homeStatLabel, awayStatLabel
}

func (widget *GamestateAndStatsWidget) createPregameTeamGameStatLayout(homeStatDatum string, categoryName string, awayStatDatum string) *widgets.QVBoxLayout {
	fontSize := 12
	gameStatLayout := widgets.NewQVBoxLayout()
	gameStatLabelLayout := widgets.NewQHBoxLayout()
	var slider *widgets.QSlider
	log.Println("widgetGameStateAndStats::createTeamGameStatLayout:: ", homeStatDatum, categoryName, awayStatDatum)
	val, ok := models.PREGAME_LABEL_STATS_MAP[categoryName]
	if ok {
		categoryName = val
	}
	//Stat slider
	if strings.Contains(homeStatDatum, ".") {
		//Float based static labels and slider
		homeStat, _ := strconv.ParseFloat(homeStatDatum, 64)
		awayStat, _ := strconv.ParseFloat(awayStatDatum, 64)
		if strings.Contains(categoryName, models.PERCENTAGE_CONST) {
			homeStat = homeStat * 100.0
			awayStat = awayStat * 100.0
			homeStatDatum = strconv.FormatFloat(homeStat, 'f', 2, 64) + "%"
			awayStatDatum = strconv.FormatFloat(awayStat, 'f', 2, 64) + "%"
		}
		homeStatLabel, awayStatLabel := widget.createStaticHomeAndAwayStatLabels(categoryName, homeStatDatum, awayStatDatum, fontSize)
		gameStatLabelLayout.AddWidget(homeStatLabel, 0, core.Qt__AlignLeft)
		gameStatLabelLayout.AddWidget(widget.createStaticDataLabel("category", categoryName, widget.UIWidget, fontSize), 0, core.Qt__AlignCenter)
		gameStatLabelLayout.AddWidget(awayStatLabel, 0, core.Qt__AlignCenter)
		//create float slider
		if homeStat > awayStat {
			slider = widget.createIntStatSlider(categoryName, "1", "0", true, false)
		} else {
			slider = widget.createIntStatSlider(categoryName, "0", "1", false, false)
		}
	} else {
		log.Println("widgetGameStateAndStats::createTeamGameStatLayout::INT SLIDER", homeStatDatum, categoryName, awayStatDatum)
		//Assume Int Slider
		homeStatLabel, awayStatLabel := widget.createStaticHomeAndAwayStatLabels(categoryName, homeStatDatum, awayStatDatum, fontSize)
		homeStat, _ := strconv.Atoi(homeStatDatum)
		awayStat, _ := strconv.Atoi(awayStatDatum)
		//Home stat dynamic label
		gameStatLabelLayout.AddWidget(homeStatLabel, 0, core.Qt__AlignLeft)
		//Stat category static label
		gameStatLabelLayout.AddWidget(widget.createStaticDataLabel("category", categoryName, widget.UIWidget, fontSize), 0, core.Qt__AlignCenter)
		//Away stat dynamic label
		gameStatLabelLayout.AddWidget(awayStatLabel, 0, core.Qt__AlignCenter)
		if strings.Contains(categoryName, "Rank") {
			if homeStat > awayStat {
				slider = widget.createIntStatSlider(categoryName, "0", "1", false, false)
			} else {
				slider = widget.createIntStatSlider(categoryName, "1", "0", true, false)

			}
		} else {
			if homeStat > awayStat {
				slider = widget.createIntStatSlider(categoryName, homeStatDatum, awayStatDatum, true, false)
			} else {
				slider = widget.createIntStatSlider(categoryName, homeStatDatum, awayStatDatum, false, false)

			}
		}
	}
	//Add label layout and slider widget
	gameStatLayout.AddLayout(gameStatLabelLayout, 0)
	gameStatLayout.AddWidget(slider, 0, core.Qt__AlignCenter)
	return gameStatLayout
}

func (widget *GamestateAndStatsWidget) createLiveTeamGameStatLayout(homeStatDatum string, categoryName string, awayStatDatum string, homeHandle bool) *widgets.QVBoxLayout {
	//Horizontal and vertical 2 row layout, 1 for labels and 1 for slider.
	fontSize := 12
	gameStatLayout := widgets.NewQVBoxLayout()
	gameStatLabelLayout := widgets.NewQHBoxLayout()
	var slider *widgets.QSlider
	log.Println("widgetGameStateAndStats::createTeamGameStatLayout:: ", homeStatDatum, categoryName, awayStatDatum)
	homeStatLabel, awayStatLabel := widget.createHomeAndAwayStatLabels(categoryName, homeStatDatum, awayStatDatum, fontSize)
	//Home stat dynamic label
	gameStatLabelLayout.AddWidget(homeStatLabel, 0, core.Qt__AlignLeft)
	//Stat category static label
	gameStatLabelLayout.AddWidget(widget.createStaticDataLabel("category", categoryName, widget.UIWidget, fontSize), 0, core.Qt__AlignCenter)
	//Away stat dynamic label
	gameStatLabelLayout.AddWidget(awayStatLabel, 0, core.Qt__AlignCenter)
	//Stat slider
	if strings.Contains(homeStatDatum, ".") {
		//create float slider
		slider = widget.createFloatStatSlider(categoryName, homeStatDatum, awayStatDatum, homeHandle, true)
	} else if strings.Contains(homeStatDatum, "/") {
		//Create float slider but first do the math to convert over to decimal
		slider = widget.createFloatStatSlider(categoryName, homeStatDatum, awayStatDatum, homeHandle, true)
	} else {
		//Assume Int Slider
		slider = widget.createIntStatSlider(categoryName, homeStatDatum, awayStatDatum, homeHandle, true)
	}
	//Add label layout and slider widget
	gameStatLayout.AddLayout(gameStatLabelLayout, 0)
	gameStatLayout.AddWidget(slider, 0, core.Qt__AlignCenter)
	return gameStatLayout
}

//For this Were gonna have to create a could different widgets for viz
//Ideally wed have a stacked widget a default stack of sliders and stats for pregame if live
//Or we have a two layouts where pregame triggers live creation and switch
//For now were gonna handle this with vis switches on two different widgets and let redraw sort it out.

func (widget *GamestateAndStatsWidget) createPregameStatsLayout() {
	var teamGameObjects []models.TeamGameStat
	widget.pregameStatsLayout = widgets.NewQVBoxLayout()
	widget.IsFuture = true
	teamGameObjects = models.ConvertTeamSeasonStatsToTeamGameStats(widget.gameController.GetTeamSeasonStatObject())
	for _, gameStatObject := range teamGameObjects {
		awayValueString := widget.gameController.ConvertAnyStatToGameStatString(&gameStatObject, false)
		homeValueString := widget.gameController.ConvertAnyStatToGameStatString(&gameStatObject, true)
		widget.pregameStatsLayout.AddLayout(widget.createPregameTeamGameStatLayout(homeValueString, gameStatObject.Category, awayValueString), 0)
	}
	widget.pregameStatsLayout.ConnectTimerEvent(func(event *core.QTimerEvent) {
		if widget.IsFuture && widget.gameController.IsLive() {
			widget.IsFuture = false
			widget.createLiveGamestatsLayout()
			widget.liveGameStatsWidget.SetLayout(widget.liveGameStatsLayout)
			widget.pregameStatsWidget.SetVisible(false)
			widget.liveGameStatsWidget.SetVisible(true)
			widget.pregameStatsLayout.DisconnectTimerEvent()
		}
	})
	widget.pregameStatsLayout.StartTimer(widget.LabelTimer, core.Qt__CoarseTimer)
}

func (widget *GamestateAndStatsWidget) createLiveGamestatsLayout() {
	var teamGameObjects []models.TeamGameStat
	widget.liveGameStatsLayout = widgets.NewQVBoxLayout()
	teamGameObjects = widget.gameController.GetTeamGameStatsObjects()
	for _, gameStatObject := range teamGameObjects {
		awayValueString := widget.gameController.ConvertAnyStatToGameStatString(&gameStatObject, false)
		homeValueString := widget.gameController.ConvertAnyStatToGameStatString(&gameStatObject, true)
		if homeValueString >= awayValueString { //We will need to figure out the types in a switch for proper compare
			widget.liveGameStatsLayout.AddLayout(widget.createLiveTeamGameStatLayout(homeValueString, gameStatObject.Category, awayValueString, true), 0)
		} else {
			widget.liveGameStatsLayout.AddLayout(widget.createLiveTeamGameStatLayout(awayValueString, gameStatObject.Category, homeValueString, false), 0)
		}
	}
}

func (widget *GamestateAndStatsWidget) createGamestateAndStatsWidget() {
	//Create main layout and widget
	gamestateAndStatsLayout := widgets.NewQVBoxLayout()
	gamestateAndStatsWidget := widgets.NewQGroupBox(widget.parentWidget)
	gamestateAndStatsWidget.SetProperty("widget-type", core.NewQVariant12("gamestatsAndGamestate"))
	//Create Child layouts and Widgets
	widget.pregameStatsWidget = widgets.NewQGroupBox(widget.parentWidget)
	widget.liveGameStatsWidget = widgets.NewQGroupBox(widget.parentWidget)
	if widget.gameController.IsFuture() || widget.gameController.IsPregame() {
		widget.createPregameStatsLayout()
		widget.pregameStatsWidget.SetLayout(widget.pregameStatsLayout)
		widget.pregameStatsWidget.SetVisible(true)
		widget.liveGameStatsWidget.SetVisible(false)
	} else {
		widget.createLiveGamestatsLayout()
		widget.liveGameStatsWidget.SetLayout(widget.liveGameStatsLayout)
		widget.liveGameStatsWidget.SetVisible(true)
		widget.pregameStatsWidget.SetVisible(false)
	}
	gamestateLayout := widget.createGamestateLayout()
	//Add Child Layouts and widgets
	gamestateAndStatsLayout.AddLayout(gamestateLayout, 0)
	gamestateAndStatsLayout.AddWidget(widget.pregameStatsWidget, 0, core.Qt__AlignCenter)
	gamestateAndStatsLayout.AddWidget(widget.liveGameStatsWidget, 0, core.Qt__AlignCenter)
	//Set Size and Stylesheet - Work off a scaling factor - base = 100 (base*1.77)*ScalingFactor and base*scalingFactor ::Scaling Factor is 2. :: 1.77 is Desired Aspect Ratio.
	gamestateAndStatsWidget.SetLayout(gamestateAndStatsLayout)
	gamestateAndStatsWidget.SetStyleSheet(CreateGameStatsAndGamestateStylesheet())
	//Set widget UI
	widget.UILayout = gamestateAndStatsLayout
	widget.UIWidget = gamestateAndStatsWidget
	widget.UIWidget.SetMinimumSize(core.NewQSize2(100, 770))
	widget.UIWidget.SetMaximumSize(core.NewQSize2(500, 1080))

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
