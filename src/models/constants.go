package models

// const game controller data
const DEFAULT_GAMESTATE_STRING = "DEFAULT GAMESTATE STRING"

// const team data
const DEFAULT_NHL_TEAM_ABBREV = "NHLF"
const DEFAULT_ID = -1
const DEFAULT_SCORE = -1
const DEFAULT_SOG = -1
const DEFAULT_RADIO_LINK = ""

// Const player data
const DEFAULT_PLAYER_NAME = " "
const DEFAULT_POSITION_CODE = " "
const DEFAULT_SWEATER_NUMBER = 0
const DEFAULT_SOI = 0

//Const Stat data
const DEFAULT_TOTAL_STAT = "100"
const DEFAULT_TOTAL_STAT_INT = 100
const DEFAULT_WINNING_STAT = "100"
const DEFAULT_WINNING_STAT_INT = 100
const DEFAULT_LOSING_STAT = "0"
const DEFAULT_LOSTING_STAT = 0
const DEFAULT_CATEGORY = "CATEGORY"

//const Delims
const VALUE_DELIMITER = "_"
const NAME_DELIMITER = "_"

//Home and Away prefixs
const DEFAULT_HOME_PREFIX = "home"
const DEFAULT_AWAY_PREFIX = "away"

//Tied prefix
const DEFAULT_TIED_PREFIX = "tied"

//Verses Consts
const DEFAULT_VERSES_STRING = " VS "

//Stats Consts
const PERCENTAGE_CONST = "%"

//Pregame Category label mapping.
var (
	PREGAME_LABEL_STATS_MAP map[string]string
)

func init() {
	PREGAME_LABEL_STATS_MAP = map[string]string{
		"PpPctg":                    "Power Play %",
		"PkPctg":                    "Penality Kill %",
		"FaceoffWinningPctg":        "Face-off %",
		"GoalsForPerGamePlayed":     "Goals Scored Average Per Game",
		"GoalsAgainstPerGamePlayed": "Goals Against Average Per Game",
		"PpPctgRank":                "Power Play Rank",
		"PkPctgRank":                "Penality Kill Rank",
		"FaceoffWinningPctgRank":    "Face-off Win Rank",
		"GoalsForPerGamePlayedRank": "Goals Scored Per Game Rank",
		"GoalsAgainstAverageRank":   "Goals Against Per Game Rank",
	}
}

//Live Game Category label mapping.
var (
	LIVE_GAME_LABEL_STATS_MAP map[string]string
)

func init() {
	LIVE_GAME_LABEL_STATS_MAP = map[string]string{
		"sog":                "Shots On Goal",
		"faceoffWinningPctg": "Face-off %",
		"powerPlay":          "Power Play",
		"powerPlayPctg":      "Power Play %",
		"pim":                "Penality Minutes",
		"hits":               "Hits",
		"blockedShots":       "Blocked Shots",
		"giveaways":          "Giveaways",
		"takeaways":          "Takeaways",
	}
}
