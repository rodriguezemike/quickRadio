//go:build ignore
// +build ignore

package views

import (
	"quickRadio/controllers"

	"github.com/therecipe/qt/widgets"
)

type QuickRadioView struct {
	gameController    *controllers.GameController
	radioController   *controllers.RadioController
	app               *widgets.QApplication
	window            *widgets.QMainWindow
	gameManagerWidget *widgets.QGroupBox
	gamesWidget       *widgets.QStackedWidget
	activeGameWidget  *widgets.QGroupBox
	activeGameIndex   int
}
