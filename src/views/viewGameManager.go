package views

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

func (view *QuickRadioView) CreateGameManagerWidget() *widgets.QGroupBox {
	gameManager := widgets.NewQGroupBox(nil)
	topbarLayout := widgets.NewQVBoxLayout()
	gameStackLayout := widgets.NewQStackedLayout()
	gameManagerLayout := widgets.NewQVBoxLayout2(gameManager)

	gameManagerLayout.AddLayout(topbarLayout, 1)
	gameManagerLayout.AddLayout(gameStackLayout, 1)
	view.gamesStackWidget = view.CreateGamesWidget()
	gameDropdown := view.CreateGameDropdownsWidget()

	topbarLayout.AddWidget(gameDropdown, 1, core.Qt__AlignAbsolute)
	gameStackLayout.AddWidget(view.gamesStackWidget)
	gameManager.SetStyleSheet(CreateGameManagerStyleSheet())
	return gameManager
}
