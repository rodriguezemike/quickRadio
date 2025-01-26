package models

type Sweater struct {
	TeamAbbrev     string
	PrimaryColor   string
	SecondaryColor string
}

type TeamData struct {
	Id   int `json:"id"`
	Name struct {
		Default string `json:"default"`
	} `json:"name"`
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
}

type PlayerOnIce struct {
	PlayerId int `json:"playerId"`
	Name     struct {
		Default string `json:"default"`
	}
	SweaterNumber int    `json:"sweaterNumber"`
	PositionCode  string `json:"positionCode"`
	Headshot      string `json:"headShot"`
	TotalSOI      int    `json:"totalSOI"`
}

type TeamOnIce struct {
	Forwards   []PlayerOnIce `json:"forwards"`
	Defensemen []PlayerOnIce `json:"defensemen"`
	Goalies    []PlayerOnIce `json:"goalies"`
	PenaltyBox []PlayerOnIce `json:"penaltyBox"`
}
type GameVersesData struct {
	SeasonSeries []struct {
		Id                int    `json:"id"`
		Season            int    `json:"season"`
		GameType          int    `json:"gameType"`
		GameDate          string `json:"gameDate"`
		StartTimeUTC      string `json:"startTimeUTC"`
		EasternUTCOffset  string `json:"easternUTCOffset"`
		VenueUTCOffset    string `json:"venueUTCOffset"`
		GameState         string `json:"gameState"`
		GameScheduleState string `json:"OK"`
		AwayTeam          struct {
			Id     int    `json:"id"`
			Abbrev string `json:"abbrev"`
			Logo   string `json:"logo"`
			Score  int    `json:"score"`
		} `json:"awayTeam"`
		HomeTeam struct {
			Id     int    `json:"id"`
			Abbrev string `json:"abbrev"`
			Logo   string `json:"logo"`
			Score  int    `json:"score"`
		} `json:"homeTeam"`
		GameCenterLink string `json:"gameCenterLink"`
	} `json:"seasonSeries"`
	SeasonSeriesWins struct {
		AwayTeamWins int `json:"awayTeamWins"`
		HomeTeamWins int `json:"homeTeamWins"`
	} `json:"seasonSeriesWins"`
	GameInfo struct {
		Referees []struct {
			Default string `json:"default"`
		} `json:"referees"`
		Linesmen []struct {
			Default string `json:"default"`
		} `json:"linesmen"`
		AwayTeam struct {
			HeadCoach struct {
				Default string `json:"default"`
			} `json:"headCoach"`
			Scratches []struct {
				Id        int `json:"id"`
				FirstName struct {
					Default string `json:"default"`
				} `json:"firstName"`
				LastName struct {
					Default string `json:"default"`
				} `json:"lastName"`
			}
		} `json:"awayTeam"`
		HomeTeam struct {
			HeadCoach struct {
				Default string `json:"default"`
			} `json:"headCoach"`
			Scratches []struct {
				Id        int `json:"id"`
				FirstName struct {
					Default string `json:"default"`
				} `json:"firstName"`
				LastName struct {
					Default string `json:"default"`
				} `json:"lastName"`
			}
		} `json:"homeTeam"`
		Linescore struct {
			ByPeriod []struct {
				PeriodDescriptor struct {
					Number               int    `json:"number"`
					PeriodType           string `json:"periodType"`
					MaxRegulationPeriods int    `json:"maxRegulationPeriods"`
				} `json:"PeriodDescriptor"`
				Away int `json:"away"`
				Home int `json:"home"`
			} `json:"byPeriod"`
			Totals struct {
				Away int `json:"away"`
				Home int `json:"home"`
			} `json:"totals"`
		} `json:"linescore"`
		ShotsByPeriod []struct {
			PeriodDescriptor struct {
				Number               int    `json:"number"`
				PeriodType           string `json:"periodType"`
				MaxRegulationPeriods int    `json:"maxRegulationPeriods"`
			} `json:"PeriodDescriptor"`
			Away int `json:"away"`
			Home int `json:"home"`
		} `json:"shotsByPeriod"`
		TeamGameStats []struct {
			Category  string `json:"category"`
			AwayValue string `json:"awayValue,string"`
			HomeValue string `json:"homeValue,string"`
		} `json:"teamGameStats"`
		GameReports struct {
			GameSummary       string `json:"gameSummary"`
			EventSummary      string `json:"eventSummary"`
			PlayByPlay        string `json:"playByPlay"`
			FaceoffSummary    string `json:"faceoffSummary"`
			FaceoffComparison string `json:"faceoffComparison"`
			Rosters           string `json:"rosters"`
			ShotSummary       string `json:"shotSummary"`
			ToiAway           string `json:"toiAway"`
			ToiHome           string `json:"toiHome"`
		} `json:"gameReports"`
	} `json:"gameInfo"`
}

type GameData struct {
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
		MaxRegulationPeriods int    `json:"maxRegulationPeriods"`
	} `json:"periodDescriptor"`
	TvBroadcasts []struct {
		Id             int    `json:"id"`
		Market         string `json:"market"`
		CountryCode    string `json:"countryCode"`
		Network        string `json:"network"`
		SequenceNumber int    `json:"sequenceNumber"`
	} `json:"tvBroadcasts"`
	GameState         string   `json:"gameState"`
	GameScheduleState string   `json:"gameScheduleState"`
	AwayTeam          TeamData `json:"awayTeam"`
	HomeTeam          TeamData `json:"homeTeam"`
	ShootoutInuse     bool     `json:"shootoutInuse"`
	MaxPeriods        int      `json:"MaxPeriods"`
	RegPeriods        int      `json:"regPeriods"`
	OtInUse           bool     `json:"otInUse"`
	TiesInUse         bool     `json:"tiesInUse"`
	Summary           struct {
		IceSurface struct {
			AwayTeam TeamOnIce `json:"awayTeam"`
			HomeTeam TeamOnIce `json:"homeTeam"`
		} `json:"IceSurface"`
		Scoring []struct {
			PeriodDescriptor struct {
				Number               int    `json:"number"`
				PeriodType           string `json:"periodType"`
				MaxRegulationPeriods int    `json:"maxRegulationPeriods"`
			} `json:"PeriodDescriptor"`
			Goals []struct {
				SituationCode string `json:"situationCode"`
				Strength      string `json:"strength"`
				PlayerId      int    `json:"playerId"`
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
				HighlightClip           int    `json:"highlightClip"`
				DiscreteClip            int    `json:"discreteClip"`
				GoalsToDate             int    `json:"goalsToDate"`
				AwayScore               int    `json:"awayScore"`
				HomeScore               int    `json:"homeScore"`
				LeadingTeamAbbrev       struct {
					Default string `json:"default"`
				}
				TimeInPeriod string `json:"timeInPeriod"`
				ShotType     string `json:"snap"`
				GoalModifier string `json:"goalModifier"`
				Assists      []struct {
					PlayerId  int `json:"playerId"`
					FirstName struct {
						Default string `json:"default"`
					} `json:"firstName"`
					LastName struct {
						Default string `json:"default"`
					} `json:"lastName"`
					Name struct {
						Default string `json:"default"`
					} `json:"name"`
					SweaterNumber int `json:"sweaterNumber"`
					AssistsToDate int `json:"assistsToDate"`
				} `json:"assists"`
				PPTReplayUrl          string `json:"pptReplayUrl"`
				HomeTeamDefendingSide string `json:"homeTeamDefendingSide"`
			} `json:"goals"`
		} `json:"scoring"`
		Shootout []struct {
			Sequence   int    `json:"sequence"`
			PlayerId   int    `json:"playerId"`
			TeamAbbrev string `json:"teamAbbrev"`
			FirstName  string `json:"firstName"`
			ShotType   string `json:"shotType"`
			Result     string `json:"result"`
			Headshot   string `json:"headshot"`
			GameWinner bool   `json:"gameWinner"`
		} `json:"shootout"`
		ThreeStars []struct {
			Star     int `json:"star"`
			PlayerId int `json:"playerId"`
			Name     struct {
				Default string `json:"default"`
			}
			SweaterNumber       int     `json:"sweaterNumber"`
			HeadShot            string  `json:"headShot"`
			Position            string  `json:"position"`
			Goals               int     `json:"goals"`
			Assists             int     `json:"assists"`
			Points              int     `json:"points"`
			TeamAbbrev          string  `json:"teamAbbrev"`
			GoalsAgainstAverage float32 `json:"goalsAgainstAverage"`
			SavePctg            float32 `json:"savePctg"`
		} `json:"threeStars"`
		Penalties []struct {
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
		SecondsRemaining int    `json:"secondsRemaining"`
		Running          bool   `json:"running"`
		InIntermission   bool   `json:"inIntermission"`
	} `json:"clock"`
	Situation struct {
		HomeTeam struct {
			Abbrev                string   `json:"abbrev"`
			SituationDescriptions []string `json:"situationDescriptions"`
			Strength              int      `json:"strength"`
		}
		AwayTeam struct {
			Abbrev                string   `json:"abbrev"`
			SituationDescriptions []string `json:"situationDescriptions"`
			Strength              int      `json:"strength"`
		}
		SituationCode    string `json:"situationCode"`
		TimeRemaining    string `json:"timeRemaining"`
		SecondsRemaining int    `json:"secondsRemaining"`
	}
}
