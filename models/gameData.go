package models

type GameDataStruct struct {
	Id             int    `json:"id"`
	Season         int    `json:"season"`
	GameType       int    `json:"gameType"`
	LimitedScoring bool   `json:"limitedScoring"`
	GameDate       string `json:"gameDate"`
	Venue          struct {
		Default string `json:"default"`
	} `json:"venue"`
	VenueLocation struct {
		Default string `json:"default"`
	} `json:"venueLocation"`
	StartTimeUTC     string `json:"startTimeUTC"`
	EasternUTCOffset string `json:"easternUTCOffset"`
	VenueUTCOffset   string `json:"venueUTCOffset"`
	VenueTimezone    string `json:"venueTimezone"`
	PeriodDescriptor struct {
		Number               int    `json:"number"`
		PeriodType           string `json:"pertiodType"`
		MaxRegulationPeriods int    `json:maxRegulationPeriods`
	} `json:"periodDescriptor"`
	TvBroadcasts []struct {
		Id             int    `json:"id"`
		Market         string `json:"market"`
		CountryCode    string `json:"countryCode"`
		Network        string `json:"network"`
		SequenceNumber int    `json:"sequenceNumber"`
	} `json:"tvBroadcasts"`
	GameState         string `json:"gameState"`
	GameScheduleState string `json:"gameScheduleState"`
	AwayTeam          struct {
		Id   int `json:"id"`
		Name struct {
			Default string `json:"default"`
		} `json:"name`
		Abbrev    string `json:"abbrev"`
		PlaceName struct {
			Default string `json:"default"`
		} `json:"placeName"`
		PlaceNameWithPreposition struct {
			Default string `json:"default"`
			Fr      string `json:"default"`
		}
		Score     int    `json:"score"`
		Sog       int    `json:"sog"`
		Logo      string `json:"logo"`
		DarkLogo  string `json:"darkLogo"`
		RadioLink string `json:"radioLink"`
	} `json:"awayTeam"`
	HomeTeam struct {
		Id   int `json:"id"`
		Name struct {
			Default string `json:"default"`
		} `json:"name`
		Abbrev    string `json:"abbrev"`
		PlaceName struct {
			Default string `json:"default"`
		} `json:"placeName"`
		PlaceNameWithPreposition struct {
			Default string `json:"default"`
			Fr      string `json:"fr"`
		}
		Score     int    `json:"score"`
		Sog       int    `json:"sog"`
		Logo      string `json:"logo"`
		DarkLogo  string `json:"darkLogo"`
		RadioLink string `json:"radioLink"`
	} `json:"homeTeam"`
	ShootoutInuse bool `json:"shootoutInuse"`
	MaxPeriods    int  `json:"shootoutInUse"`
	RegPeriods    int  `json:"regPeriods"`
	OtInUse       bool `json:"otInUse"`
	TiesInUse     bool `json:"tiesInUse"`
	Summary       struct {
		IceSurface struct {
			AwayTeam struct {
				Forwards   []string `json:"forwards"`
				Defensemen []string `json:"defensemen"`
				Goalies    []string `json:"goalies"`
				PenaltyBox []string `json:"penaltyBox"`
			} `json:"awayTeam"`
			HomeTeam struct {
				Forwards   []string `json:"forwards"`
				Defensemen []string `json:"defensemen"`
				Goalies    []string `json:"goalies"`
				PenaltyBox []string `json:"penaltyBox"`
			} `json:"homeTeam"`
		} `json:"IceSurface"`
		Scoring []struct {
			PeriodDescriptor struct {
				Number               int    `json:number`
				PeriodType           string `json:periodType`
				MaxRegulationPeriods int    `json:maxRegulationPeriods`
			} `json:"PeriodDescriptor"`
			Goals []struct {
				SituationCode int    `json:situationCode`
				Strength      string `json:strength`
				PlayerId      int    `json:playerId`
				FirstName     struct {
					Default string `json:"default"`
				}
				LastName struct {
					Default string `json:"default"`
				}
				Name struct {
					Default string `json:"default"`
				}
				TeamAbbrev struct {
					Default string `json:"default"`
				}
				Headshot                string `json:"headshot"`
				HIghlightClipSharingUrl string `json:"highlightClipSharingUrl"`
				HighlightClip           string `json:"highlightClip"`
				DiscreteClip            string `json:"discreteClip"`
				GoalsToDate             int    `json:"goalsToDate"`
				AwayScore               string `json:"awayScore"`
				HomeScore               string `json:"homeScore"`
				LeadingTeamAbbrev       struct {
					Default string `json:"default"`
				}
				TimeInPeriod          string   `json:"timeInPeriod"`
				ShotType              string   `json:"snap"`
				GoalModifier          string   `json:"goalModifier"`
				Assists               []string `json:"assists"`
				PPTReplayUrl          string   `json:"pptReplayUrl"`
				HomeTeamDefendingSide string   `json:"homeTeamDefendingSide"`
			} `json:"goals"`
		} `json:"scoring"`
		Shootout   []string `json:"shootout"`
		ThreeStars []string `json:"threeStars"`
		Penalties  []struct {
			PeriodDescriptor struct {
				Number               int    `json:"number"`
				PeriodType           string `json:"periodType"`
				MaxRegulationPeriods int    `json:"maxRegulationPeriods"`
			} `json:"PeriodDescriptor"`
			Penalties []struct {
				TimeInPeriod      string `json:"timeInPeriod"`
				Type              string `json:"type"`
				Duration          string `json:"duration"`
				CommittedByPlayer string `json:"committedByPlayer"`
				TeamAbbrev        struct {
					Default string `json:"default"`
				}
				DrawnBy string `json:"drawnBy"`
				DescKey string `json:"descKey"`
			} `json:"penalites"`
		} `json:"penalites"`
	}
	Clock struct {
		TimeRemaining    string `json:"timeRemaining"`
		SecondsRemaining int    `json:secondsRemaining`
		Running          bool   `json:"running"`
		InIntermission   bool   `json:"inIntermission"`
	} `json:"clock"`
	Situation struct {
		HomeTeam struct {
			Abbrev                string   `json:"abbrev"`
			SituationDescriptions []string `json:situationDescriptions`
			Strength              int      `json:"strength"`
		}
		AwayTeam struct {
			Abbrev                string   `json:"abbrev"`
			SituationDescriptions []string `json:situationDescriptions`
			Strength              int      `json:"strength"`
		}
		SituationCode    string `json:"situationCode"`
		TimeRemaining    string `json:"timeRemaining"`
		SecondsRemaining int    `json:"secondsRemaining"`
	}
}
