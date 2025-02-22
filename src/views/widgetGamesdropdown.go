package views

import "github.com/therecipe/qt/widgets"

type GamesDropdownWidget struct {
	Width        int
	UIWidget     *widgets.QComboBox
	UILayout     *widgets.QHBoxLayout
	parentWidget *widgets.QStackedWidget
	games        []*GameView
	gameNames    []string
}

func (widget *GamesDropdownWidget) createGamesDropdownWidget() {
	for _, game := range widget.games {
		widget.gameNames = append(widget.gameNames, game.GetGameName())
	}
	widget.UILayout = widgets.NewQHBoxLayout()
	widget.UIWidget = widgets.NewQComboBox(widget.parentWidget)
	widget.UIWidget.SetStyleSheet(CreateDropdownStyleSheet())
	widget.UIWidget.SetFixedWidth(widget.Width)
	widget.UIWidget.AddItems(widget.gameNames)
	widget.UIWidget.ConnectCurrentIndexChanged(func(index int) {
		//Sort out update maps.
		widget.games[index].ClearUpdateMaps()
		widget.parentWidget.SetCurrentIndex(index)
	})
}

func CreateNewGamesDropdownWidget(width int, games []*GameView) *GamesDropdownWidget {
	widget := GamesDropdownWidget{}
	widget.Width = width
	widget.games = games
	widget.createGamesDropdownWidget()
	return &widget
}
