package views

import (
	"path/filepath"
	"quickRadio/models"
	"quickRadio/quickio"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

func (view *QuickRadioView) getIceRinkPixmap() *gui.QPixmap {
	dir := quickio.GetProjectDir()
	path := filepath.Join(dir, "assets", "svgs", "rink.svg")
	rinkPixmap := gui.NewQPixmap3(path, "svg", core.Qt__AutoColor)
	return rinkPixmap
}

func (view *QuickRadioView) CreateIceCenterWidget(gameDataObject models.GameData, gamecenterLink string, gameWidget *widgets.QGroupBox) *widgets.QGroupBox {
	layout := widgets.NewQVBoxLayout2(gameWidget)
	centerIceWidget := widgets.NewQGroupBox(gameWidget)
	gameStateLabel := view.CreateDataLabel("GAMESTATE", view.gameController.GetGamestateString(&gameDataObject), gamecenterLink, gameWidget)
	layout.AddWidget(CreateIceRinklabel(gameWidget, gameDataObject, view), 0, core.Qt__AlignCenter)
	layout.AddWidget(gameStateLabel, 0, core.Qt__AlignCenter)
	centerIceWidget.SetLayout(layout)
	centerIceWidget.ConnectTimerEvent(func(event *core.QTimerEvent) {
		for _, value := range view.activeGameDataUpdateMap {
			if !value {
				return
			}
		}
		for key := range view.activeGameDataUpdateMap {
			view.activeGameDataUpdateMap[key] = false
			view.gameController.ConsumeActiveGameData()
		}
	})
	centerIceWidget.StartTimer(view.LabelTimer*2, core.Qt__CoarseTimer)
	return centerIceWidget
}
