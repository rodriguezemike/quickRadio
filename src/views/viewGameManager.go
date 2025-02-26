package views

import (
	"context"
	"sync"
	"time"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

type GameManagerView struct {
	DropdownWidth int
	LabelTimer    int
	UIWidget      *widgets.QGroupBox
	UILayout      *widgets.QVBoxLayout
	gamesDropdown *GamesDropdownWidget
	gamesStack    *widgets.QStackedWidget
	games         []*GameView
	parentWidget  *widgets.QGroupBox
}

func (view *GameManagerView) GoUpdateGames(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			var gamesToUpdate []*GameView
			for _, game := range view.games {
				if game.gameController.IsLive() || game.gameController.IsFuture() {
					gamesToUpdate = append(gamesToUpdate, game)
				}
			}
			var workGroup sync.WaitGroup
			for i := range view.games {
				workGroup.Add(i)
				go func(game *GameView) {
					defer workGroup.Done()
					game.gameController.ProduceGameData()
					game.ClearUpdateMaps()
				}(gamesToUpdate[i])
			}
			workGroup.Wait()
			time.Sleep(time.Duration(view.LabelTimer-1000) * time.Second)
		}
	}
}

func (view *GameManagerView) createGamesWidget() *widgets.QStackedWidget {
	gamesStack := widgets.NewQStackedWidget(view.UIWidget)
	for _, gameView := range view.games {
		gamesStack.AddWidget(gameView.UIWidget)
	}
	return gamesStack
}

func (view *GameManagerView) createGameManagerView() {
	//Set UI widget and Layout
	view.UILayout = widgets.NewQVBoxLayout()
	view.UIWidget = widgets.NewQGroupBox(view.parentWidget)
	//Create Child Widgets
	view.gamesStack = view.createGamesWidget()
	view.gamesDropdown = CreateNewGamesDropdownWidget(view.DropdownWidth, view.games)
	//Add Child Widget
	view.UILayout.AddWidget(view.gamesDropdown.UIWidget, 0, core.Qt__AlignTop)
	view.UILayout.AddWidget(view.gamesStack, 0, core.Qt__AlignTop)
	//Set Size and Stylesheet - Work off a scaling factor - base = 100 (base*1.77)*ScalingFactor and base*scalingFactor ::Scaling Factor is 2. :: 1.77 is Desired Aspect Ratio.
	view.UIWidget.SetLayout(view.UILayout)
}

func CreateNewGameManagerView(dropdownWidth int, games []*GameView, parentWidget *widgets.QGroupBox) *GameManagerView {
	view := GameManagerView{}
	view.DropdownWidth = dropdownWidth
	view.games = games
	view.parentWidget = parentWidget
	view.createGameManagerView()
	return &view
}
