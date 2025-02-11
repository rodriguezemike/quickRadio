package views_test

import (
	"flag"
	"os"
	"quickRadio/controllers"
	"quickRadio/views"
	"sync"
	"testing"

	"github.com/therecipe/qt/widgets"
)

var flagCI = flag.Bool("skip-for-ci", false, "Skip Visual UI test for ci envs")

func createTestTeamWidget() (*widgets.QApplication, *controllers.GameController, int, *sync.Mutex, *views.TeamWidget) {
	os.Setenv("DISPLAY", ":99")
	app := widgets.NewQApplication(len(os.Args), os.Args)
	gameController := controllers.NewGameController()
	labelTimer := 1000
	radioLock := sync.Mutex{}
	gameWidget := widgets.NewQGroupBox(nil)
	widget := views.CreateNewTeamWidget(labelTimer, -1, false, gameController, &radioLock, gameWidget)
	return app, gameController, labelTimer, &radioLock, widget
}

func TestTeamWidgetConstructor(t *testing.T) {
	_, gameController, labelTimer, radioLock, widget := createTestTeamWidget()

	if !widget.RadioLockReferenceTest(radioLock) {
		t.Fatalf(`!widget.RadioLockReferenceTest(&radioLock) | Wanted address of Radiolock %v | check to see if pass by copy in constructor.`, radioLock)
	}
	if !widget.GameControllerReferenceTest(gameController) {
		t.Fatalf(`!widget.GameControllerReferenceTest(&gameController) | Wanted address of gameController %v | check to see if pass by copy in constructor.`, gameController)
	}
	if !widget.LabelTimerTest(labelTimer) {
		t.Fatalf(`!widget.LabelTimerTest(labelTimer) | Wanted Label Timer Value %v | Got %v`, labelTimer, widget.LabelTimer)
	}

}

// Here we can automate the UI visual test with a image compare of the widget with some design expectation image.
// Functionality of the button can be tested with UI automation, but if the radio controller is working then we have already isolated the test
// That needs to be ran. Unless we do this in parallel to test multiple tonditions it dont make that much of a diff compared to just manually testing the UI.
// Plus we need eyes on our UI for this level.
func TestTeamWidgetUI(t *testing.T) {
	flag.Parse()
	if !*flagCI {
		app, _, _, _, widget := createTestTeamWidget()
		app.SetApplicationDisplayName("TestTeamWidgetUI")
		window := widgets.NewQMainWindow(nil, 0)
		window.SetCentralWidget(widget.UI)
		window.Show()
		app.Exec()
		app.Quit()
	} else {
		t.Skip("We are in a CI env and skipping Visual based test.")
	}

}
