package views

import (
	"context"
	"log"
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
			var workGroup sync.WaitGroup
			for _, game := range view.games {
				if game.gameController.IsLive() || game.gameController.IsFuture() || game.gameController.IsPregame() {
					gamesToUpdate = append(gamesToUpdate, game)
				}
			}
			log.Println(view.games)
			for i := range gamesToUpdate {
				workGroup.Add(1)
				go func(game *GameView) {
					defer workGroup.Done()
					game.gameController.ProduceGameData()
					game.ClearUpdateMaps()
				}(gamesToUpdate[i])
			}
			workGroup.Wait()
			time.Sleep(time.Duration(view.LabelTimer-1000) * time.Millisecond)
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
	view.gamesDropdown = CreateNewGamesDropdownWidget(view.DropdownWidth, view.games, view.gamesStack)
	//Add Child Widget
	view.UILayout.AddWidget(view.gamesDropdown.UIWidget, 0, core.Qt__AlignTop)
	view.UILayout.AddWidget(view.gamesStack, 0, core.Qt__AlignTop)
	//Set Size and Stylesheet - Work off a scaling factor - base = 100 (base*1.77)*ScalingFactor and base*scalingFactor ::Scaling Factor is 2. :: 1.77 is Desired Aspect Ratio.
	view.UIWidget.SetLayout(view.UILayout)
}

func CreateNewGameManagerView(dropdownWidth int, games []*GameView, parentWidget *widgets.QGroupBox, labelTimer int) *GameManagerView {
	view := GameManagerView{}
	view.DropdownWidth = dropdownWidth
	view.games = games
	view.parentWidget = parentWidget
	view.LabelTimer = labelTimer
	view.createGameManagerView()
	return &view
}
