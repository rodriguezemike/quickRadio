package views

import (
	"log"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

type GamesDropdownWidget struct {
	Width        int
	UIWidget     *widgets.QGroupBox
	UILayout     *widgets.QHBoxLayout
	parentWidget *widgets.QStackedWidget
	dropdown     *widgets.QComboBox
	games        []*GameView
	gameNames    []string
}

func (widget *GamesDropdownWidget) createGamesDropdownWidget() {
	for _, game := range widget.games {
		widget.gameNames = append(widget.gameNames, game.GetGameName())
	}
	widget.UILayout = widgets.NewQHBoxLayout()
	widget.UIWidget = widgets.NewQGroupBox(widget.parentWidget)
	widget.dropdown = widgets.NewQComboBox(widget.UIWidget)
	widget.dropdown.SetStyleSheet(CreateDropdownStyleSheet())
	widget.dropdown.SetFixedWidth(widget.Width)
	widget.dropdown.AddItems(widget.gameNames)
	widget.dropdown.ConnectCurrentIndexChanged(func(index int) {
		log.Println("Index", index)
		go widget.games[index].gameController.ProduceGameData()
		log.Println("Clearing Update Maps")
		widget.games[index].ClearUpdateMaps()
		log.Println("Setting Index on parent widget which should be a qstack widget")
		widget.parentWidget.SetCurrentIndex(index)
		log.Println("Did it.")
	})
	widget.UILayout.AddWidget(widget.dropdown, 0, core.Qt__AlignTop)
	widget.UIWidget.SetLayout(widget.UILayout)
}

func CreateNewGamesDropdownWidget(width int, games []*GameView, parentWidget *widgets.QStackedWidget) *GamesDropdownWidget {
	widget := GamesDropdownWidget{}
	widget.Width = width
	widget.games = games
	widget.parentWidget = parentWidget
	widget.createGamesDropdownWidget()
	return &widget
}
