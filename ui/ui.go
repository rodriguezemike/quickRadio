package ui

import (
	"encoding/json"
	"os"
	"path/filepath"
	"quickRadio/audio"
	"quickRadio/game"
	"quickRadio/models"
	"quickRadio/radioErrors"
	"runtime"
	"strconv"
	"strings"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

//Once were Happy with this, we should refactor into a proper singleton UI
//This would make it easier to update as we move towards AI Things in later phasees of the UI.

func GetProjectDir() string {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filepath.Dir(filename))
	return dir
}

func GetLogoPath(logoFilename string) string {
	dir := GetProjectDir()
	path := filepath.Join(dir, "assets", "svgs", "logos", logoFilename)
	return path
}

func GetIceRinkPixmap() *gui.QPixmap {
	dir := GetProjectDir()
	path := filepath.Join(dir, "assets", "svgs", "rink.svg")
	rinkPixmap := gui.NewQPixmap3(path, "svg", core.Qt__AutoColor)
	return rinkPixmap
}

func GetTeamPixmap(teamAbbrev string) *gui.QPixmap {
	teamStringSlice := []string{teamAbbrev, "light.svg"}
	teamLogoFilename := strings.Join(teamStringSlice, "_")
	teamLogoPath := GetLogoPath(teamLogoFilename)
	teamPixmap := gui.NewQPixmap3(teamLogoPath, "svg", core.Qt__AutoColor)
	return teamPixmap
}

func GetTeamIcon(teamAbbrev string) *gui.QIcon {
	teamPixmap := GetTeamPixmap(teamAbbrev)
	teamIcon := gui.NewQIcon2(teamPixmap)
	return teamIcon
}

func CreateTeamRadioStreamButton(teamAbbrev string, radioLink string, sweaterColors map[string][]string, gameWidget *widgets.QGroupBox) *widgets.QPushButton {
	teamIcon := GetTeamIcon(teamAbbrev)
	button := widgets.NewQPushButton(gameWidget)
	button.SetStyleSheet(CreateTeamBackgroundStylesheet(teamAbbrev, sweaterColors))
	button.SetObjectName("radio_" + teamAbbrev)
	button.SetIcon(teamIcon)
	button.SetFixedHeight(320)
	button.SetFixedWidth(320)
	button.SetIconSize(button.FrameSize())
	button.SetCheckable(true)
	button.ConnectToggled(func(onCheck bool) {
		if onCheck {
			go audio.StartRadioFun(radioLink)
		} else {
			go audio.StopRadioFun()
		}
	})
	return button
}

func CreateIceRinklabel(gameWidget *widgets.QGroupBox, homeTeamAbbrev models.TeamData) *widgets.QLabel {
	iceRinkPixmap := GetIceRinkPixmap()
	homeTeamPixmap := GetTeamPixmap(homeTeamAbbrev.Abbrev)
	compositePixmap := gui.NewQPixmap2(iceRinkPixmap.Size())
	compositePixmap.Fill(gui.QColor_FromRgba(0))
	painter := gui.NewQPainter2(compositePixmap)
	painter.DrawPixmap9(0, 0, iceRinkPixmap)
	painter.DrawPixmap11((iceRinkPixmap.Size().Width()/2)-30, (iceRinkPixmap.Size().Height()/2)-35, 64, 64, homeTeamPixmap)
	painter.SetOpacity(.7)
	painter.DrawPixmap9(0, 0, iceRinkPixmap)
	iceRinkLabel := widgets.NewQLabel2("", gameWidget, core.Qt__Widget)
	iceRinkLabel.SetPixmap(compositePixmap)
	//We want new ice pixmap. Find a good one.
	//Here we need to Draw the little circles with hovers on the Ice.
	//Take the total of players then divide by the area on the ice
	//Place approp
	//Do this for each level on the ice.
	//Look into edge data to draw circles. May not be avail.
	//Finish drawing the Abbrevs on the Sides of the ICe. Lets take it to the between the dots instead of the Brodeur Zone.
	//Do we want all mugs and call em from disc?
	//It would be kinda cool to actually Draw on the ice. Maybe for later?
	return iceRinkLabel
}

func CreateDataLabel(name string, data string, genercenterLink string, gameWidget *widgets.QGroupBox) *widgets.QLabel {
	//Here too we want to move away from pulling for every label and move to a cache, checking it with goroutines
	//Only updating that one, once every 10 seconds per visible game. Reducing our API calls to 1 per 10 seconds Rather than 3.
	//Also Note, we wiill want Store a whole bunch of game info in here, we may want to pass in some modified GDO.
	//This is a stop gap for now so that we dont continue to hit the API a bunch of unneccesary times.
	label := widgets.NewQLabel2(data, gameWidget, core.Qt__Widget)
	label.SetObjectName(name)
	label.SetAccessibleDescription(genercenterLink)
	font := label.Font()
	font.SetPointSize(32)
	label.SetFont(font)
	label.SetStyleSheet(CreateLabelStylesheet())
	label.ConnectTimerEvent(func(event *core.QTimerEvent) {
		if label.IsVisible() {
			if strings.Contains(label.ObjectName(), "Score") {
				counter, _ := strconv.Atoi(strings.Split(label.ObjectName(), " ")[1])
				if counter > 10 && counter%10 == 0 {
					gdo := game.GetGameDataObject(label.AccessibleDescription())
					if strings.Contains(label.ObjectName(), gdo.HomeTeam.Abbrev) {
						label.SetText(strconv.Itoa(gdo.HomeTeam.Score))
						label.Repaint()
					} else {
						label.SetText(strconv.Itoa(gdo.AwayTeam.Score))
						label.Repaint()
					}
				} else {
					counter += 1
					label.SetObjectName(strings.Split(label.ObjectName(), " ")[0] + " " + strconv.Itoa(counter))
				}
			} else if strings.Contains(label.ObjectName(), "GameState") {
				//Run clock when clock is running, counting down. Only Resetting When clock is stopped. That should get us somewhat close to instant feedback.
				gdo := game.GetGameDataObject(label.AccessibleDescription())
				if gdo.GameState == "LIVE" || gdo.GameState == "CRIT" {
					if !gdo.Clock.InIntermission {
						label.SetText(gdo.GameState + " - " + "P" + strconv.Itoa(gdo.PeriodDescriptor.Number) + " " + gdo.Clock.TimeRemaining)
					} else {
						label.SetText(gdo.GameState + "INT" + strconv.Itoa(gdo.PeriodDescriptor.Number) + " " + gdo.Clock.TimeRemaining)
					}
				} else {
					label.SetText(gdo.GameState)
				}
				label.Repaint()
			}
		}
	})
	label.StartTimer(1000, core.Qt__PreciseTimer)
	return label
}

func CreateTeamWidget(team models.TeamData, gamecenterLink string, sweaterColors map[string][]string, gameWidget *widgets.QGroupBox) *widgets.QGroupBox {
	teamLayout := widgets.NewQVBoxLayout2(gameWidget)
	teamWidget := widgets.NewQGroupBox(gameWidget)
	teamLayout.AddWidget(CreateTeamRadioStreamButton(team.Abbrev, team.RadioLink, sweaterColors, gameWidget), 0, core.Qt__AlignCenter)
	teamLayout.AddWidget(CreateDataLabel("TeamAbbrev", team.Abbrev, gamecenterLink, gameWidget), 0, core.Qt__AlignCenter)
	teamLayout.AddWidget(CreateDataLabel(team.Abbrev+"Score "+"0", strconv.Itoa(team.Score), gamecenterLink, gameWidget), 0, core.Qt__AlignCenter)
	teamWidget.SetLayout(teamLayout)
	return teamWidget
}

func CreateIceCenterWidget(gameDataObject models.GameData, gamecenterLink string, gameWidget *widgets.QGroupBox) *widgets.QGroupBox {
	layout := widgets.NewQVBoxLayout2(gameWidget)
	centerIceWidget := widgets.NewQGroupBox(gameWidget)
	layout.AddWidget(CreateIceRinklabel(gameWidget, gameDataObject.HomeTeam), 0, core.Qt__AlignCenter)
	layout.AddWidget(CreateDataLabel("GameState "+strconv.Itoa(gameDataObject.Clock.SecondsRemaining), gameDataObject.GameState+" - "+"P"+strconv.Itoa(gameDataObject.PeriodDescriptor.Number)+" "+gameDataObject.Clock.TimeRemaining, gamecenterLink, gameWidget),
		0, core.Qt__AlignCenter)
	centerIceWidget.SetLayout(layout)
	return centerIceWidget
}

func CreateGameWidgetFromGameDataObject(gameDataObject models.GameData, gamecenterLink string, sweaterColors map[string][]string, gameStackWidget *widgets.QStackedWidget) *widgets.QGroupBox {
	layout := widgets.NewQGridLayout(gameStackWidget)
	gameWidget := widgets.NewQGroupBox(gameStackWidget)
	layout.AddWidget2(CreateTeamWidget(gameDataObject.HomeTeam, gamecenterLink, sweaterColors, gameWidget), 0, 0, core.Qt__AlignCenter)
	layout.AddWidget2(CreateIceCenterWidget(gameDataObject, gamecenterLink, gameWidget), 0, 1, core.Qt__AlignCenter)
	layout.AddWidget2(CreateTeamWidget(gameDataObject.AwayTeam, gamecenterLink, sweaterColors, gameWidget), 0, 2, core.Qt__AlignCenter)
	layout.AddWidget2(CreateGameDetailsWidgetFromGameDataObject(gameDataObject, gamecenterLink, gameWidget), 3, 1, core.Qt__AlignCenter)
	gameWidget.SetLayout(layout)
	gameWidget.SetStyleSheet(CreateGameStylesheet(gameDataObject.HomeTeam.Abbrev, gameDataObject.AwayTeam.Abbrev, sweaterColors))
	return gameWidget
}

func CreateGamesWidget(gameDataObjects []models.GameData, gamecenterLinks []string, sweaterColors map[string][]string, gameManager *widgets.QGroupBox) *widgets.QStackedWidget {
	gameStackWidget := widgets.NewQStackedWidget(gameManager)
	for i, gameDataObject := range gameDataObjects {
		gameWidget := CreateGameWidgetFromGameDataObject(gameDataObject, gamecenterLinks[i], sweaterColors, gameStackWidget)
		gameStackWidget.AddWidget(gameWidget)
	}
	gameStackWidget.SetCurrentIndex(0)
	return gameStackWidget
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
			gdo := game.GetGameDataObject(gamecenterLink)
			gameDetailsJson, _ := json.MarshalIndent(gdo, "", " ")
			jsonDumpLabel.SetText(string(gameDetailsJson))
			jsonDumpLabel.Repaint()
		}
	})
	jsonDumpLabel.StartTimer(30000, core.Qt__VeryCoarseTimer)
	scrollableArea.SetWidgetResizable(true)
	scrollableArea.SetWidget(jsonDumpLabel)
	gameDetailsLayout.AddWidget(scrollableArea)
	gameDetailsWidget.SetLayout(gameDetailsLayout)
	return gameDetailsWidget
}

func CreateGameDropdownsWidget(gameDataObjects []models.GameData, gamesStack *widgets.QStackedWidget, gameManager *widgets.QGroupBox) *widgets.QComboBox {
	var gameNames []string
	for _, gameDataObject := range gameDataObjects {
		gameNames = append(gameNames, gameDataObject.HomeTeam.Abbrev+" vs "+gameDataObject.AwayTeam.Abbrev)
	}
	dropdown := widgets.NewQComboBox(gameManager)
	dropdown.SetStyleSheet(CreateDropdownStyleSheet())
	dropdown.SetFixedWidth(600)
	dropdown.AddItems(gameNames)
	dropdown.ConnectCurrentIndexChanged(func(index int) {
		gamesStack.SetCurrentIndex(index)
	})
	return dropdown
}

func CreateGameManagerWidget(gameDataObjects []models.GameData, gamecenterLinks []string, sweaterColors map[string][]string) *widgets.QGroupBox {
	gameManager := widgets.NewQGroupBox(nil)
	topbarLayout := widgets.NewQVBoxLayout()
	gameStackLayout := widgets.NewQStackedLayout()
	gameManagerLayout := widgets.NewQVBoxLayout2(gameManager)

	gameManagerLayout.AddLayout(topbarLayout, 1)
	gameManagerLayout.AddLayout(gameStackLayout, 1)
	gamesStack := CreateGamesWidget(gameDataObjects, gamecenterLinks, sweaterColors, gameManager)
	gameDropdown := CreateGameDropdownsWidget(gameDataObjects, gamesStack, gameManager)

	topbarLayout.AddWidget(gameDropdown, 1, core.Qt__AlignAbsolute)
	gameStackLayout.AddWidget(gamesStack)
	gameManager.SetStyleSheet(CreateGameManagerStyleSheet())
	return gameManager
}

func CreateLoadingScreen() *widgets.QSplashScreen {
	pixmap := GetTeamPixmap("NHL")
	splash := widgets.NewQSplashScreen(pixmap, core.Qt__Widget)
	return splash
}

func KillAllTheFun() {
	audio.RadioKillFun()
}

func CreateAndRunUI() {
	app := widgets.NewQApplication(len(os.Args), os.Args)
	app.SetApplicationDisplayName("QuickRadio")
	app.ConnectAboutToQuit(func() {
		KillAllTheFun()
	})
	app.ConnectDestroyQApplication(func() {
		KillAllTheFun()
	})
	app.SetWindowIcon(GetTeamIcon("NHLF"))
	loadingScreen := CreateLoadingScreen()
	loadingScreen.Show()
	window := widgets.NewQMainWindow(nil, 0)
	gameDataObjects, gamecenterLinks := game.UIGetGameDataObjectsAndGameLandingLinks()
	sweaterColors := game.GetSweaterColors()
	gameManager := CreateGameManagerWidget(gameDataObjects, gamecenterLinks, sweaterColors)
	window.SetCentralWidget(gameManager)
	loadingScreen.Finish(nil)
	window.Show()
	app.Exec()
}
