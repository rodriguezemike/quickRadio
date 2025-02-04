package views

import (
	"quickRadio/models"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

func (view *QuickRadioView) CreateGamesWidget() *widgets.QStackedWidget {
	gameStackWidget := widgets.NewQStackedWidget(view.gameManagerWidget)
	for i, gameDataObject := range view.gameController.GetGameDataObjects() {
		gameWidget := view.CreateGameWidgetFromGameDataObject(gameDataObject, view.gameController.Landinglinks[i])
		gameStackWidget.AddWidget(gameWidget)
	}
	gameStackWidget.SetCurrentIndex(view.activeGameIndex)
	return gameStackWidget
}

func (view *QuickRadioView) CreateGameWidgetFromGameDataObject(gameDataObject models.GameData, gamecenterLink string) *widgets.QGroupBox {
	layout := widgets.NewQGridLayout(view.gamesStackWidget)
	gameWidget := widgets.NewQGroupBox(view.gamesStackWidget)
	layout.AddWidget2(view.CreateTeamWidget(gameDataObject.HomeTeam, gamecenterLink, gameWidget), 0, 0, core.Qt__AlignCenter)
	layout.AddWidget2(view.CreateIceCenterWidget(gameDataObject, gamecenterLink, gameWidget), 0, 1, core.Qt__AlignCenter)
	layout.AddWidget2(view.CreateTeamWidget(gameDataObject.AwayTeam, gamecenterLink, gameWidget), 0, 2, core.Qt__AlignCenter)
	gameWidget.SetLayout(layout)
	gameWidget.SetStyleSheet(CreateGameStylesheet(gameDataObject.HomeTeam.Abbrev, gameDataObject.AwayTeam.Abbrev))
	return gameWidget
}
