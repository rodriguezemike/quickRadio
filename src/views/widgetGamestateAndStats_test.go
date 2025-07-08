package views_test

import (
	"flag"
	"os"
	"quickRadio/controllers"
	"quickRadio/views"
	"testing"

	"github.com/therecipe/qt/widgets"
)

func createTestGamestateAndStatsWidget() (*widgets.QApplication, *controllers.GameController, int, *views.GamestateAndStatsWidget) {
	os.Setenv("DISPLAY", ":99") //Mainly for Ubuntu CI. Theres no harm in running every time.
	app := widgets.NewQApplication(len(os.Args), os.Args)
	gameController := controllers.CreateNewDefaultGameController()
	labelTimer := 1000
	gameWidget := widgets.NewQGroupBox(nil)
	widget := views.CreateNewGamestateAndStatsWidget(labelTimer, gameController, gameWidget)
	return app, gameController, labelTimer, widget
}

func TestPowerplayLabelCreation(t *testing.T) {
	flag.Parse()
	if !*flagCI {
		app, _, _, widget := createTestGamestateAndStatsWidget()

		// Test that the widget was created successfully
		if widget.UIWidget == nil {
			t.Fatalf("GamestateAndStats widget should be created")
		}

		// Note: Testing powerplay functionality requires access to private fields
		// or additional public methods. For now, we test that the widget creates successfully.

		app.Quit()
	} else {
		t.Skip("We are in a CI env and skipping Visual based test.")
	}
}

func TestGamestateAndStatsUI(t *testing.T) {
	flag.Parse()
	if !*flagCI {
		app, _, _, widget := createTestGamestateAndStatsWidget()
		app.SetApplicationDisplayName("TestGamestateAndStats")
		window := widgets.NewQMainWindow(nil, 0)
		window.SetCentralWidget(widget.UIWidget)
		window.Show()
		app.Exec()
		app.Quit()
	} else {
		t.Skip("We are in a CI env and skipping Visual based test.")
	}
}
