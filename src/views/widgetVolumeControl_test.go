package views_test

import (
	"flag"
	"os"
	"quickRadio/views"
	"testing"

	"github.com/therecipe/qt/widgets"
)

func createTestVolumeControlWidget() (*widgets.QApplication, *views.VolumeControlWidget) {
	os.Setenv("DISPLAY", ":99") //Mainly for Ubuntu CI. Theres no harm in running every time.
	app := widgets.NewQApplication(len(os.Args), os.Args)
	parentWidget := widgets.NewQGroupBox(nil)
	widget := views.CreateVolumeControlWidget(parentWidget)
	return app, widget
}

func createTestGameViewForVolume() *views.GameView {
	gameView := views.CreateNewDefaultGameView()
	return gameView
}

func TestVolumeControlWidgetCreation(t *testing.T) {
	flag.Parse()
	if !*flagCI {
		app, widget := createTestVolumeControlWidget()

		if widget.UIWidget == nil {
			t.Fatalf("VolumeControlWidget UIWidget should not be nil")
		}

		if widget.UILayout == nil {
			t.Fatalf("VolumeControlWidget UILayout should not be nil")
		}

		app.Quit()
	} else {
		t.Skip("We are in a CI env and skipping Visual based test.")
	}
}

func TestVolumeControlWidgetGameAssignment(t *testing.T) {
	flag.Parse()
	if !*flagCI {
		app, widget := createTestVolumeControlWidget()
		gameView := createTestGameViewForVolume()

		// Test setting current game (we can't access the private field directly,
		// but we can test that SetCurrentGame doesn't panic)
		widget.SetCurrentGame(gameView)

		// Test setting nil game
		widget.SetCurrentGame(nil)

		app.Quit()
	} else {
		t.Skip("We are in a CI env and skipping Visual based test.")
	}
}

func TestVolumeControlWidgetUI(t *testing.T) {
	flag.Parse()
	if !*flagCI {
		app, widget := createTestVolumeControlWidget()
		app.SetApplicationDisplayName("TestVolumeControlWidget")
		window := widgets.NewQMainWindow(nil, 0)
		window.SetCentralWidget(widget.UIWidget)
		window.Show()
		app.Exec()
		app.Quit()
	} else {
		t.Skip("We are in a CI env and skipping Visual based test.")
	}
}
