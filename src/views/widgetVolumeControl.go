package views

import (
	"fmt"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

type VolumeControlWidget struct {
	UIWidget     *widgets.QGroupBox
	UILayout     *widgets.QHBoxLayout
	volumeSlider *widgets.QSlider
	volumeLabel  *widgets.QLabel
	muteButton   *widgets.QPushButton
	parentWidget *widgets.QGroupBox
	currentGame  *GameView
	isMuted      bool
	savedVolume  float64
}

func (widget *VolumeControlWidget) SetCurrentGame(gameView *GameView) {
	widget.currentGame = gameView
	widget.updateVolumeDisplay()
}

func (widget *VolumeControlWidget) updateVolumeDisplay() {
	if widget.currentGame == nil {
		widget.volumeSlider.SetEnabled(false)
		widget.muteButton.SetEnabled(false)
		widget.volumeLabel.SetText("Volume: --")
		return
	}

	// Check if any team is playing radio
	homeTeamPlaying := widget.currentGame.HomeTeamWidget.isPlaying
	awayTeamPlaying := widget.currentGame.AwayTeamWidget.isPlaying

	if homeTeamPlaying || awayTeamPlaying {
		widget.volumeSlider.SetEnabled(true)
		widget.muteButton.SetEnabled(true)

		var volume float64
		if homeTeamPlaying && widget.currentGame.HomeTeamWidget.radioController != nil {
			volume = widget.currentGame.HomeTeamWidget.radioController.GetVolume()
		} else if awayTeamPlaying && widget.currentGame.AwayTeamWidget.radioController != nil {
			volume = widget.currentGame.AwayTeamWidget.radioController.GetVolume()
		} else {
			volume = 1.0
		}

		widget.volumeSlider.SetValue(int(volume * 100))
		widget.volumeLabel.SetText(fmt.Sprintf("Volume: %d%%", int(volume*100)))
	} else {
		widget.volumeSlider.SetEnabled(false)
		widget.muteButton.SetEnabled(false)
		widget.volumeLabel.SetText("Volume: --")
	}
}

func (widget *VolumeControlWidget) setVolume(volume float64) {
	if widget.currentGame == nil {
		return
	}

	if widget.currentGame.HomeTeamWidget.isPlaying && widget.currentGame.HomeTeamWidget.radioController != nil {
		widget.currentGame.HomeTeamWidget.radioController.SetVolume(volume)
	}
	if widget.currentGame.AwayTeamWidget.isPlaying && widget.currentGame.AwayTeamWidget.radioController != nil {
		widget.currentGame.AwayTeamWidget.radioController.SetVolume(volume)
	}
}

func (widget *VolumeControlWidget) toggleMute() {
	if widget.isMuted {
		// Unmute
		widget.setVolume(widget.savedVolume)
		widget.volumeSlider.SetValue(int(widget.savedVolume * 100))
		widget.muteButton.SetText("🔊")
		widget.isMuted = false
	} else {
		// Mute
		widget.savedVolume = float64(widget.volumeSlider.Value()) / 100.0
		widget.setVolume(0.0)
		widget.volumeSlider.SetValue(0)
		widget.muteButton.SetText("🔇")
		widget.isMuted = true
	}
	widget.updateVolumeDisplay()
}

func (widget *VolumeControlWidget) createVolumeControlWidget() {
	widget.UILayout = widgets.NewQHBoxLayout()
	widget.UIWidget = widgets.NewQGroupBox(widget.parentWidget)
	widget.UIWidget.SetProperty("widget-type", core.NewQVariant12("volumeControl"))

	// Volume label
	widget.volumeLabel = widgets.NewQLabel2("Volume: --", widget.UIWidget, core.Qt__Widget)
	widget.volumeLabel.SetObjectName("volumeLabel")
	widget.volumeLabel.SetMinimumWidth(80)

	// Mute button
	widget.muteButton = widgets.NewQPushButton2("🔊", widget.UIWidget)
	widget.muteButton.SetObjectName("muteButton")
	widget.muteButton.SetMaximumSize(core.NewQSize2(30, 25))
	widget.muteButton.SetEnabled(false)
	widget.muteButton.ConnectClicked(func(checked bool) {
		widget.toggleMute()
	})

	// Volume slider
	widget.volumeSlider = widgets.NewQSlider2(core.Qt__Horizontal, widget.UIWidget)
	widget.volumeSlider.SetObjectName("volumeSlider")
	widget.volumeSlider.SetRange(0, 200) // 0% to 200%
	widget.volumeSlider.SetValue(100)    // Default 100%
	widget.volumeSlider.SetMaximumWidth(100)
	widget.volumeSlider.SetEnabled(true)
	widget.volumeSlider.ConnectValueChanged(func(value int) {
		if !widget.isMuted {
			volume := float64(value) / 100.0
			widget.setVolume(volume)
			widget.volumeLabel.SetText(fmt.Sprintf("Volume: %d%%", value))
		}
	})

	// Set slider style using centralized stylesheet
	widget.volumeSlider.SetStyleSheet(CreateVolumeSliderStylesheet())

	// Add widgets to layout
	widget.UILayout.AddWidget(widget.volumeLabel, 0, core.Qt__AlignLeft)
	widget.UILayout.AddWidget(widget.muteButton, 0, core.Qt__AlignLeft)
	widget.UILayout.AddWidget(widget.volumeSlider, 0, core.Qt__AlignLeft)

	widget.UIWidget.SetLayout(widget.UILayout)
	widget.UIWidget.SetMaximumWidth(250)
	widget.UIWidget.SetStyleSheet(CreateVolumeControlWidgetStylesheet())

	// Initialize state
	widget.isMuted = false
	widget.savedVolume = 1.0
}

func CreateVolumeControlWidget(parentWidget *widgets.QGroupBox) *VolumeControlWidget {
	widget := &VolumeControlWidget{}
	widget.parentWidget = parentWidget
	widget.createVolumeControlWidget()
	return widget
}
