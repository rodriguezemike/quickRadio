package quickio

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"path/filepath"
	"quickRadio/models"
	"quickRadio/radioErrors"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
)

func GetDataFromResponse(url string) ([]byte, io.ReadCloser) {
	resp, err := http.Get(url)
	radioErrors.ErrorLog(err)
	byteValue, err := io.ReadAll(resp.Body)
	radioErrors.ErrorLog(err)
	return byteValue, resp.Body

}

func GetQualityStreamSlug(radioLink string, audioQuality string) string {
	byteValue, bodyCloser := GetDataFromResponse(radioLink)
	audioQualitySlug, err := ExtractQualityStreamSlug(string(byteValue), audioQuality)
	radioErrors.ErrorLog(err)
	bodyCloser.Close()
	return audioQualitySlug
}

func ExtractQualityStreamSlug(m3uContents string, audioQuality string) (string, error) {
	for _, line := range strings.Split(m3uContents, "\n") {
		if strings.Contains(line, audioQuality) && strings.Contains(line, ".m3u8") {
			return line, nil
		}
	}
	return "", errors.New("couldnt find audio quality string")
}

func GetAACSlugsFromQualityFile(m3uContents string) ([]string, error) {
	var audioFiles []string
	for _, line := range strings.Split(m3uContents, "\n") {
		if !strings.Contains(line, "#") && strings.Contains(line, ".aac") {
			audioFiles = append(audioFiles, line)
		}
	}
	if len(audioFiles) == 0 {
		return nil, errors.New("couldnt find audio files")
	}
	return audioFiles, nil
}

func BuildQualityRadioPath(radioLink string, qualitySlug string) string {
	return strings.Join(append(strings.Split(radioLink, "/")[:len(strings.Split(radioLink, "/"))-1], qualitySlug), "/")
}

func BuildAACRadioPath(radioQualityLink string, aacFile string) string {
	return strings.Join(append(strings.Split(radioQualityLink, "/")[:len(strings.Split(radioQualityLink, "/"))-1], aacFile), "/")
}

func GetAACPaths(qualityRadioPath string) []string {
	log.Println("IO::GetAACPaths")
	audioFilepaths := []string{}
	log.Println("IO::GetAACPaths::POLLING ", qualityRadioPath)
	resp, err := http.Get(qualityRadioPath)
	radioErrors.ErrorLog(err)
	if err == nil {
		defer resp.Body.Close()
		byteValue, _ := io.ReadAll(resp.Body)
		audioSlugs, _ := GetAACSlugsFromQualityFile(string(byteValue))
		for _, audioSlug := range audioSlugs {
			audioFilepaths = append(audioFilepaths, BuildAACRadioPath(qualityRadioPath, audioSlug))
		}
	}
	return audioFilepaths
}

func GetRadioFormatLinkAndDirectory(radioLink string, sampleRate string) (string, string, []string) {
	qualitySlug := GetQualityStreamSlug(radioLink, sampleRate)
	radioFormatLink := BuildQualityRadioPath(radioLink, qualitySlug)
	aacPaths := GetAACPaths(radioFormatLink)
	wavPaths := GoDownloadAndTranscodeAACs(aacPaths)
	radioDirectory := filepath.Dir(wavPaths[len(wavPaths)-1])
	return radioFormatLink, radioDirectory, wavPaths
}

func GetGameHtml(linksMap map[string]interface{}) string {
	var html string
	sleepTimer := linksMap["load_sleep_timer"].(float64)
	baseUrl := fmt.Sprintf("%v", linksMap["base"])

	opts := append(chromedp.DefaultExecAllocatorOptions[:], chromedp.Headless, chromedp.NoSandbox)
	allocContext, ctxCancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer ctxCancel()
	ctx, cancel := chromedp.NewContext(allocContext)
	defer cancel()

	err := chromedp.Run(ctx,
		chromedp.Navigate(baseUrl),
		chromedp.Sleep(time.Duration(sleepTimer)*time.Millisecond),
		chromedp.ActionFunc(func(ctx context.Context) error {
			rootNode, err := dom.GetDocument().Do(ctx)
			if err != nil {
				return err
			}
			html, err = dom.GetOuterHTML().WithNodeID(rootNode.NodeID).Do(ctx)
			return err
		}),
	)
	radioErrors.ErrorLog(err)
	return html
}

func GetGameLandingLinks() []string {
	linksMap := GetLinksJson()
	html := GetGameHtml(linksMap)
	gamecenterBase := fmt.Sprintf("%v", linksMap["gamecenter_api_base"])
	gamecenterLanding := fmt.Sprintf("%v", linksMap["gamecenter_api_slug"])
	gameRegex := fmt.Sprintf("%v", linksMap["game_regex"])
	landingLinks := GetGameLandingLinksFromHTML(html, gamecenterBase, gamecenterLanding, gameRegex)
	return landingLinks
}

func GetGameLandingLinksFromHTML(html string, gamecenterBase string, gamecenterLanding string, gameRegex string) []string {
	var gameLandingLinks []string
	gameRegexObject, _ := regexp.Compile(gameRegex)
	allGames := gameRegexObject.FindAllString(html, -1)
	currentDate := strings.ReplaceAll(time.Now().Format(time.DateOnly), "-", "/")
	for _, possibleGame := range allGames {
		if strings.Contains(possibleGame, currentDate) {
			gamecenterLink := strings.Trim(possibleGame, "\"")
			landingLink := gamecenterBase + strings.Split(gamecenterLink, "/")[len(strings.Split(gamecenterLink, "/"))-1] + gamecenterLanding
			gameLandingLinks = append(gameLandingLinks, landingLink)
		}
	}
	return gameLandingLinks
}

func GetGameLandingLink(html string, gamecenterBase string, gamecenterLanding string, gameRegexs []string, teamAbbrev string) (string, error) {
	log.Println()
	log.Println("func getGameLandingLink START")
	var gamecenterLink string
	for _, game := range gameRegexs {
		newGame := strings.Replace(game, "TEAMABBREV", strings.ToLower(teamAbbrev), -1)
		log.Println(teamAbbrev)
		gameRegex, _ := regexp.Compile(newGame)
		allGames := gameRegex.FindAllString(html, -1)
		currentDate := strings.ReplaceAll(time.Now().Format(time.DateOnly), "-", "/")
		for _, possibleGame := range allGames {
			if strings.Contains(possibleGame, currentDate) {
				gamecenterLink = strings.Trim(possibleGame, "\"")
				log.Println("gamecenterLink : ", gamecenterLink)
				break
			}
		}
	}
	if gamecenterLink == "" {
		return "", errors.New("couldnt find game center link")
	}
	gameLandingLink := gamecenterBase + strings.Split(gamecenterLink, "/")[len(strings.Split(gamecenterLink, "/"))-1] + gamecenterLanding
	return gameLandingLink, nil
}

func GetGameDataObject(gameLandingLink string) models.GameData {
	var gameData = &models.GameData{}
	byteValue, bodyCloser := GetDataFromResponse(gameLandingLink)
	err := json.Unmarshal(byteValue, gameData)
	radioErrors.ErrorLog(err)
	bodyCloser.Close()
	return *gameData
}

func GetGameVersesData(gameLandingLink string) models.GameVersesData {
	var versesData = &models.GameVersesData{}
	log.Println(gameLandingLink)
	gameLandingLink = strings.Replace(gameLandingLink, "landing", "right-rail", 1)
	log.Println(gameLandingLink)
	byteValue, bodyCloser := GetDataFromResponse(gameLandingLink)
	err := json.Unmarshal(byteValue, versesData)
	radioErrors.ErrorLog(err)
	bodyCloser.Close()
	return *versesData
}

func GetGameDataObjectFromLandingLinks(landingLinks []string) []models.GameData {
	var gameDataObjects []models.GameData
	for _, landingLink := range landingLinks {
		gameDataObjects = append(gameDataObjects, GetGameDataObject(landingLink))
	}
	return gameDataObjects
}

func GoGetGameDataObjectsFromLandingLinks(landingLinks []string) []models.GameData {
	//If we have to, we will sort by game ID
	var gameDataObjects []models.GameData
	var workGroup sync.WaitGroup
	log.Println("IO::GoGetGameDataObjectsFromLandingLinks")
	log.Println("IO::GoGetGameDataObjectsFromLandingLinks::landingLinks", landingLinks)
	for i := 0; i < len(landingLinks); i++ {
		workGroup.Add(1)
		go func(path string) {
			defer workGroup.Done()
			var gameData = &models.GameData{}
			log.Println("IO::GoGetGameDataObjectsFromLandingLinks::anonFunc::path", path)
			byteValue, bodyCloser := GetDataFromResponse(path)
			err := json.Unmarshal(byteValue, gameData)
			radioErrors.ErrorLog(err)
			gameDataObjects = append(gameDataObjects, *gameData)
			bodyCloser.Close()
		}(landingLinks[i])
		workGroup.Wait()
	}
	log.Println("IO::GoGetGameDataObjectsFromLandingLinks::gameDataObjects", gameDataObjects)
	return gameDataObjects
}

func GoGetGameVersesDataFromLandingLinks(landingLinks []string) []models.GameVersesData {
	//If we have to, we will sort by GameId.
	var gameVersesDataObjects []models.GameVersesData
	var workGroup sync.WaitGroup
	for i := 0; i < len(landingLinks); i++ {
		workGroup.Add(1)
		go func(path string) {
			defer workGroup.Done()
			var versesData = &models.GameVersesData{}
			path = strings.Replace(path, "landing", "right-rail", 1)
			byteValue, bodyCloser := GetDataFromResponse(path)
			err := json.Unmarshal(byteValue, versesData)
			radioErrors.ErrorLog(err)
			gameVersesDataObjects = append(gameVersesDataObjects, *versesData)
			bodyCloser.Close()
		}(landingLinks[i])
		workGroup.Wait()
	}
	return gameVersesDataObjects
}

func GoDownloadAndTranscodeAACs(paths []string) []string {
	//log.Println("IO::DownloadAndTranscodeAACs")
	var wavpaths []string
	var workGroup sync.WaitGroup
	for i := 0; i < len(paths); i++ {
		workGroup.Add(1)
		go func(path string) {
			defer workGroup.Done()
			localpath, err := DownloadAAC(path)
			radioErrors.ErrorLog(err)
			wavPath := TranscodeToWave(localpath)
			//_ = os.Remove(localpath)
			wavpaths = append(wavpaths, strings.TrimSpace(wavPath))
		}(paths[i])
		workGroup.Wait()
	}
	//log.Println("IO::DownloadAndTranscodeAACs::wavpaths", wavpaths)
	return wavpaths
}

func DownloadAAC(aacRequestPath string) (string, error) {
	//log.Println("IO::DownloadAAC")
	quickRadioTempDirectory := GetQuickTmpFolder()
	filename := strings.Split(aacRequestPath, "/")[len(strings.Split(aacRequestPath, "/"))-1]
	gameSubDirectory := filepath.Join(strings.Split(aacRequestPath, "/")[4 : len(strings.Split(aacRequestPath, "/"))-2]...)
	CreateDirectory(filepath.Join(quickRadioTempDirectory, gameSubDirectory))
	filepath := filepath.Join(quickRadioTempDirectory, gameSubDirectory, filename)
	if runtime.GOOS == "windows" {
		if !DoesFileExist(filepath) {
			cmd := exec.Command("curl", "-o", filepath, aacRequestPath)
			_, err := cmd.CombinedOutput()
			radioErrors.ErrorLog(err)
			err = DoesFileExistErr(filepath)
			radioErrors.ErrorLog(err)
		}
	} else if runtime.GOOS == "linux" {
		if !DoesFileExist(filepath) {
			log.Println("IO::DownloadAAC:: ", "wget ", aacRequestPath, "-O", filepath)
			cmd := exec.Command("wget", aacRequestPath, "-O", filepath)
			_, err := cmd.CombinedOutput()
			radioErrors.ErrorLog(err)
			err = DoesFileExistErr(filepath)
			radioErrors.ErrorLog(err)
		}
	} else {
		if !DoesFileExist(filepath) {
			log.Println("IO::DownloadAAC:: ", "wget ", aacRequestPath, "-O", filepath)
			cmd := exec.Command("wget", aacRequestPath, "-O", filepath)
			_, err := cmd.CombinedOutput()
			radioErrors.ErrorLog(err)
			err = DoesFileExistErr(filepath)
			radioErrors.ErrorLog(err)
		}
	}
	//log.Println("IO::DownloadAAC::filepath", filepath)
	return filepath, nil
}

func UpdateRadioWavsWithContext(ctx context.Context, qualityLink string) {
	log.Println("IO::UpdateRadioWavsWithContext")
	running := false
	for {
		select {
		case <-ctx.Done():
			log.Println("IO::UpdateRadioWavsWithContext::Done.")
			return
		default:
			if !running {
				log.Println("IO::UpdateRadioWavsWithContext::Running.")
				aacPaths := GetAACPaths(qualityLink)
				go GoDownloadAndTranscodeAACs(aacPaths)
				running = true
			}
			time.Sleep(1 * time.Second)
		}
	}
}

func UpdateRadioWavs(qualityLink string) {
	aacPaths := GetAACPaths(qualityLink)
	GoDownloadAndTranscodeAACs(aacPaths)
}
