package views

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

type GameManagerView struct {
	DropdownWidth     int
	LabelTimer        int
	UIWidget          *widgets.QGroupBox
	UILayout          *widgets.QVBoxLayout
	topControlsLayout *widgets.QHBoxLayout
	gamesDropdown     *GamesDropdownWidget
	gamesStack        *widgets.QStackedWidget
	nowPlayingLabel   *widgets.QLabel
	volumeControls    *VolumeControlWidget
	games             []*GameView
	parentWidget      *widgets.QGroupBox
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
				if game.gameController.IsLive() || game.gameController.IsFuture() || game.gameController.IsPregame() || game.gameController.IsFinal() {
					gamesToUpdate = append(gamesToUpdate, game)
				}
			}
			for i := range gamesToUpdate {
				workGroup.Add(1)
				go func(game *GameView) {
					defer workGroup.Done()
					game.gameController.ProduceGameData()
					game.ClearUpdateMaps()
				}(gamesToUpdate[i])
			}
			workGroup.Wait()

			// Update now playing status for current game
			currentIndex := view.gamesStack.CurrentIndex()
			if currentIndex >= 0 && currentIndex < len(view.games) {
				currentGame := view.games[currentIndex]
				isPlaying := currentGame.IsRadioPlaying()
				if isPlaying {
					homeTeam, awayTeam := currentGame.GetPlayingTeams()
					view.UpdateNowPlayingLabel(true, homeTeam, awayTeam)
				} else {
					view.UpdateNowPlayingLabel(false, "", "")
				}

				// Update volume controls
				if view.volumeControls != nil {
					view.volumeControls.updateVolumeDisplay()
				}
			}

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

func (view *GameManagerView) createTopControlsLayout() *widgets.QHBoxLayout {
	topLayout := widgets.NewQHBoxLayout()

	// Add games dropdown on the left
	topLayout.AddWidget(view.gamesDropdown.UIWidget, 0, core.Qt__AlignLeft)

	// Add spacer to push controls to the right
	topLayout.AddStretch(1)

	// Add now playing label in the center-right
	view.nowPlayingLabel = widgets.NewQLabel2(" ", view.UIWidget, core.Qt__Window)
	view.nowPlayingLabel.SetObjectName("nowPlayingLabel")
	view.nowPlayingLabel.SetAlignment(core.Qt__AlignCenter)
	view.nowPlayingLabel.SetStyleSheet("QLabel { color: white; font-size: 16px; font-weight: bold; }")
	topLayout.AddWidget(view.nowPlayingLabel, 0, core.Qt__AlignCenter)

	// Add volume controls on the right
	view.volumeControls = CreateVolumeControlWidget(view.UIWidget)
	topLayout.AddWidget(view.volumeControls.UIWidget, 0, core.Qt__AlignRight)

	return topLayout
}

func (view *GameManagerView) UpdateNowPlayingLabel(isPlaying bool, homeTeam, awayTeam string) {
	if isPlaying && homeTeam != "" && awayTeam != "" {
		view.nowPlayingLabel.SetText(fmt.Sprintf("Now playing - %s vs %s", homeTeam, awayTeam))
	} else {
		view.nowPlayingLabel.SetText("")
	}
}

func (view *GameManagerView) SetVolumeControlsGame(gameView *GameView) {
	if view.volumeControls != nil {
		view.volumeControls.SetCurrentGame(gameView)
	}
}

func (view *GameManagerView) createGameManagerView() {
	//Set UI widget and Layout
	view.UILayout = widgets.NewQVBoxLayout()
	view.UIWidget = widgets.NewQGroupBox(view.parentWidget)
	//Create Child Widgets
	view.gamesStack = view.createGamesWidget()
	view.gamesDropdown = CreateNewGamesDropdownWidget(view.DropdownWidth, view.games, view.gamesStack)

	// Create top controls layout
	view.topControlsLayout = view.createTopControlsLayout()

	// Connect dropdown change to update volume controls
	view.gamesDropdown.dropdown.ConnectCurrentIndexChanged(func(index int) {
		view.gamesStack.SetCurrentIndex(index)
		if index >= 0 && index < len(view.games) {
			view.SetVolumeControlsGame(view.games[index])
		}
	})

	//Add Child Widget
	view.UILayout.AddLayout(view.topControlsLayout, 0)
	view.UILayout.AddWidget(view.gamesStack, 0, core.Qt__AlignTop)
	//Set Size and Stylesheet - Work off a scaling factor - base = 100 (base*1.77)*ScalingFactor and base*scalingFactor ::Scaling Factor is 2. :: 1.77 is Desired Aspect Ratio.
	view.UIWidget.SetLayout(view.UILayout)
	view.UIWidget.SetStyleSheet(CreateGameManagerStyleSheet())
	view.UIWidget.SetMaximumWidth(1920)

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
