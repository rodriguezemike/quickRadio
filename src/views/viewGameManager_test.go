package views_test

import (
	"context"
	"flag"
	"os"
	"quickRadio/views"
	"testing"
	"time"

	"github.com/therecipe/qt/widgets"
)

var flagCIGameManager = flag.Bool("skip-for-ci-game-manager", false, "Skip Visual UI test for ci envs")

func createTestGameManagerView() (*widgets.QApplication, []*views.GameView, *views.GameManagerView) {
	os.Setenv("DISPLAY", ":99") //Mainly for Ubuntu CI. Theres no harm in running every time.
	app := widgets.NewQApplication(len(os.Args), os.Args)

	// Create test games
	games := createTestGames()

	// Create parent widget
	parentWidget := widgets.NewQGroupBox(nil)

	// Create GameManagerView
	dropdownWidth := 200
	labelTimer := 1000
	gameManagerView := views.CreateNewGameManagerView(dropdownWidth, games, parentWidget, labelTimer)

	return app, games, gameManagerView
}

func createTestGames() []*views.GameView {
	var games []*views.GameView

	// Create a few test games using the default constructor
	for i := 0; i < 3; i++ {
		gameView := views.CreateNewDefaultGameView()
		games = append(games, gameView)
	}

	return games
}

func TestGameManagerViewConstructor(t *testing.T) {
	_, _, gameManagerView := createTestGameManagerView()

	// Test basic properties (public fields)
	if gameManagerView.DropdownWidth != 200 {
		t.Fatalf("Expected DropdownWidth 200, got %d", gameManagerView.DropdownWidth)
	}

	if gameManagerView.LabelTimer != 1000 {
		t.Fatalf("Expected LabelTimer 1000, got %d", gameManagerView.LabelTimer)
	}

	// Test UI components were created (public fields)
	if gameManagerView.UIWidget == nil {
		t.Fatalf("UIWidget should not be nil")
	}

	if gameManagerView.UILayout == nil {
		t.Fatalf("UILayout should not be nil")
	}
}

func TestGameManagerViewUpdateNowPlayingLabel(t *testing.T) {
	_, _, gameManagerView := createTestGameManagerView()

	// Test that UpdateNowPlayingLabel method exists and can be called
	// Since nowPlayingLabel is private, we can't directly verify the text
	// but we can ensure the method doesn't panic
	gameManagerView.UpdateNowPlayingLabel(true, "LAK", "VGK")
	gameManagerView.UpdateNowPlayingLabel(false, "", "")
	gameManagerView.UpdateNowPlayingLabel(true, "LAK", "")

	// If we reach here, the method calls succeeded
}

func TestGameManagerViewSetVolumeControlsGame(t *testing.T) {
	_, games, gameManagerView := createTestGameManagerView()

	if len(games) == 0 {
		t.Skip("No games available for testing")
	}

	// Test that SetVolumeControlsGame method exists and can be called
	// Since volumeControls is private, we can't verify the result directly
	gameManagerView.SetVolumeControlsGame(games[0])

	// If we reach here, the method call succeeded
}

func TestGameManagerViewGoUpdateGames(t *testing.T) {
	_, _, gameManagerView := createTestGameManagerView()

	// Test that GoUpdateGames can be started and cancelled
	ctx, cancel := context.WithCancel(context.Background())

	// Start the update goroutine
	go gameManagerView.GoUpdateGames(ctx)

	// Let it run briefly
	time.Sleep(100 * time.Millisecond)

	// Cancel the context
	cancel()

	// Give it time to exit
	time.Sleep(100 * time.Millisecond)

	// If we reach here without hanging, the test passes
	// (The goroutine should have exited when context was cancelled)
}

func TestGameManagerViewUI(t *testing.T) {
	flag.Parse()
	if !*flagCIGameManager {
		app, _, gameManagerView := createTestGameManagerView()
		app.SetApplicationDisplayName("TestGameManagerView")
		window := widgets.NewQMainWindow(nil, 0)
		window.SetCentralWidget(gameManagerView.UIWidget)
		window.Resize2(1200, 800)
		window.Show()

		// Run briefly for visual inspection
		go func() {
			time.Sleep(2 * time.Second)
			app.Quit()
		}()

		app.Exec()
	} else {
		t.Skip("We are in a CI env and skipping Visual based test.")
	}
}

func TestGameManagerViewMemoryCleanup(t *testing.T) {
	_, _, gameManagerView := createTestGameManagerView()

	// Test that we can create and cleanup multiple instances
	for i := 0; i < 5; i++ {
		games := createTestGames()
		parentWidget := widgets.NewQGroupBox(nil)
		tempView := views.CreateNewGameManagerView(200, games, parentWidget, 1000)

		// Basic validation
		if tempView.UIWidget == nil {
			t.Fatalf("Iteration %d: UIWidget should not be nil", i)
		}

		// Cleanup
		tempView.UIWidget.DestroyQObject()
	}

	// Original view should still be valid
	if gameManagerView.UIWidget == nil {
		t.Fatalf("Original gameManagerView should still be valid")
	}
}
