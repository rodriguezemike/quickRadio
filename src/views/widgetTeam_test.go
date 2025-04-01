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

func createTestTeamWidget() (*widgets.QApplication, *controllers.TeamController, int, *sync.Mutex, *views.TeamWidget) {
	os.Setenv("DISPLAY", ":99") //Mainly for Ubuntu CI. Theres no harm in running every time.
	app := widgets.NewQApplication(len(os.Args), os.Args)
	teamController := controllers.CreateNewDefaultTeamController()
	labelTimer := 1000
	radioLock := sync.Mutex{}
	gameWidget := widgets.NewQGroupBox(nil)
	widget := views.CreateNewTeamWidget(labelTimer, teamController, &radioLock, gameWidget, nil)
	return app, teamController, labelTimer, &radioLock, widget
}

func TestTeamWidgetConstructor(t *testing.T) {
	_, teamController, labelTimer, radioLock, widget := createTestTeamWidget()

	if !widget.RadioLockReferenceTest(radioLock) {
		t.Fatalf(`!widget.RadioLockReferenceTest(&teamController) | Wanted address of Radiolock %v | check to see if pass by copy in constructor.`, radioLock)
	}
	if !widget.GameControllerReferenceTest(teamController) {
		t.Fatalf(`!widget.GameControllerReferenceTest(&teamController) | Wanted address of teamController %v | check to see if pass by copy in constructor.`, teamController)
	}
	if !widget.LabelTimerTest(labelTimer) {
		t.Fatalf(`!widget.LabelTimerTest(labelTimer) | Wanted Label Timer Value %v | Got %v`, labelTimer, widget.LabelTimer)
	}

}

func TestTeamWidgetUI(t *testing.T) {
	flag.Parse()
	if !*flagCI {
		app, _, _, _, widget := createTestTeamWidget()
		app.SetApplicationDisplayName("TestTeamWidgetU")
		window := widgets.NewQMainWindow(nil, 0)
		window.SetCentralWidget(widget.UIWidget)
		window.Show()
		app.Exec()
		app.Quit()
	} else {
		t.Skip("We are in a CI env and skipping Visual based test.")
	}

}
