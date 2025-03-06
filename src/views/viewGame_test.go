package views_test

import (
	"flag"
	"os"
	"quickRadio/views"
	"testing"

	"github.com/therecipe/qt/widgets"
)

func createTestGameView() (*widgets.QApplication, *views.GameView) {
	os.Setenv("DISPLAY", ":99") //Mainly for Ubuntu CI. Theres no harm in running every time.
	app := widgets.NewQApplication(len(os.Args), os.Args)
	view := views.CreateNewDefaultGameView()
	return app, view
}

func TestGameViewUI(t *testing.T) {
	flag.Parse()
	if !*flagCI {
		app, view := createTestGameView()
		app.SetApplicationDisplayName("TestGameView")
		window := widgets.NewQMainWindow(nil, 0)
		window.SetCentralWidget(view.UIWidget)
		window.Show()
		app.Exec()
		app.Quit()
	} else {
		t.Skip("We are in a CI env and skipping Visual based test.")
	}
}
