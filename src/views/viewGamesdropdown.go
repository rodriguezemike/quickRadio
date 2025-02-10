//go:build ignore
// +build ignore

package views

import "github.com/therecipe/qt/widgets"

func (view *QuickRadioView) CreateGameDropdownsWidget() *widgets.QComboBox {
	var gameNames []string
	for _, gameDataObject := range view.gameController.GetGameDataObjects() {
		gameNames = append(gameNames, gameDataObject.HomeTeam.Abbrev+" vs "+gameDataObject.AwayTeam.Abbrev)
	}
	dropdown := widgets.NewQComboBox(view.gameManagerWidget)
	dropdown.SetStyleSheet(CreateDropdownStyleSheet())
	dropdown.SetFixedWidth(600)
	dropdown.AddItems(gameNames)
	dropdown.ConnectCurrentIndexChanged(func(index int) {
		view.activeGameDataUpdateMap = nil
		view.activeGameDataUpdateMap = map[string]bool{}
		go view.gameController.SwitchActiveGame(index)
		view.gamesStackWidget.SetCurrentIndex(index)
		view.activeGameWidget = view.gamesStackWidget.CurrentWidget()
	})
	return dropdown
}
