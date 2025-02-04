package views

import (
	"quickRadio/controllers"
	"quickRadio/models"
	"quickRadio/quickio"
	"strconv"
	"strings"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

func (view *QuickRadioView) getTeamPixmap(teamAbbrev string) *gui.QPixmap {
	teamStringSlice := []string{teamAbbrev, "light.svg"}
	teamLogoFilename := strings.Join(teamStringSlice, "_")
	teamLogoPath := quickio.GetLogoPath(teamLogoFilename)
	teamPixmap := gui.NewQPixmap3(teamLogoPath, "svg", core.Qt__AutoColor)
	return teamPixmap
}

func (view *QuickRadioView) GetTeamIcon(teamAbbrev string) *gui.QIcon {
	teamPixmap := view.getTeamPixmap(teamAbbrev)
	teamIcon := gui.NewQIcon2(teamPixmap)
	return teamIcon
}

func (view *QuickRadioView) SetTeamDataUIObjectName(teamAbbrev string, uiLabel string, delimiter string) string {
	return teamAbbrev + delimiter + uiLabel
}

func (view *QuickRadioView) CreateTeamWidget(team models.TeamData, gamecenterLink string, gameWidget *widgets.QGroupBox) *widgets.QGroupBox {
	teamLayout := widgets.NewQVBoxLayout2(gameWidget)
	teamWidget := widgets.NewQGroupBox(gameWidget)
	teamLayout.AddWidget(view.CreateTeamRadioStreamButton(team.Abbrev, team.RadioLink, gameWidget), 0, core.Qt__AlignCenter)
	teamLayout.AddWidget(view.CreateDataLabel("TeamAbbrev", team.Abbrev, gamecenterLink, gameWidget), 0, core.Qt__AlignCenter)
	teamLayout.AddWidget(view.CreateDataLabel(view.SetTeamDataUIObjectName(team.Abbrev, "SOG", "_"), "SOG: "+strconv.Itoa(team.Sog), gamecenterLink, gameWidget), 0, core.Qt__AlignCenter)
	teamLayout.AddWidget(view.CreateDataLabel(view.SetTeamDataUIObjectName(team.Abbrev, "SCORE", "_"), strconv.Itoa(team.Score), gamecenterLink, gameWidget), 0, core.Qt__AlignCenter)
	teamWidget.SetLayout(teamLayout)
	return teamWidget
}

func (view *QuickRadioView) CreateTeamRadioStreamButton(teamAbbrev string, radioLink string, gameWidget *widgets.QGroupBox) *widgets.QPushButton {
	radioSampleRate := "192K"
	teamIcon := view.GetTeamIcon(teamAbbrev)
	button := widgets.NewQPushButton(gameWidget)
	button.SetStyleSheet(CreateTeamButtonStylesheet(view.gameController.Sweaters[teamAbbrev]))
	button.SetObjectName(view.SetTeamDataUIObjectName(teamAbbrev, "RADIO", "_"))
	button.SetIcon(teamIcon)
	button.SetIconSize(button.FrameSize())
	button.SetCheckable(true)
	button.ConnectToggled(func(onCheck bool) {
		teamAbbrev, _ = view.GetTeamDataFromUIObjectName(button.ObjectName(), "_")
		if onCheck {
			if radioLink != "" {
				view.radioController = controllers.NewRadioController(radioLink, teamAbbrev, radioSampleRate)
				go view.radioController.PlayRadio()
			}
		} else {
			if radioLink != "" {
				go view.radioController.StopRadio(teamAbbrev)
				view.radioController = nil
			}
		}
	})
	return button
}
