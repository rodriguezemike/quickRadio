package quickio

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
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
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
)

func GetRadioFormatLinkAndDirectory(radioLink string) (string, string, []string) {
	qualitySlug := GetQualityStreamSlug(radioLink, "192K")
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

func GetSweaterColors() map[string][]string {
	sweaterColors := make(map[string][]string)
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filepath.Dir(filename))
	path := filepath.Join(dir, "assets", "teams", "sweater_colors.txt")

	fileObject, err := os.Open(path)
	radioErrors.ErrorLog(err)
	defer fileObject.Close()
	scanner := bufio.NewScanner(fileObject)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		lineColors := strings.Split(scanner.Text(), ";")
		sweaterColors[lineColors[0]] = lineColors[1:]
	}
	return sweaterColors
}

func GetSweaters() map[string]models.Sweater {
	sweaters := make(map[string]models.Sweater)
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filepath.Dir(filename))
	path := filepath.Join(dir, "assets", "teams", "sweater_colors.txt")

	fileObject, err := os.Open(path)
	radioErrors.ErrorLog(err)
	defer fileObject.Close()
	scanner := bufio.NewScanner(fileObject)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		var sweater models.Sweater
		lineColors := strings.Split(scanner.Text(), ";")
		sweater.TeamAbbrev = lineColors[0]
		sweater.PrimaryColor = lineColors[1]
		sweater.SecondaryColor = lineColors[2]
		sweaters[sweater.TeamAbbrev] = sweater
	}

	return sweaters
}

func GetDataFromFile(path string) []byte {
	fileObject, err := os.Open(path)
	radioErrors.ErrorLog(err)
	defer fileObject.Close()
	byteValue, err := io.ReadAll(fileObject)
	radioErrors.ErrorLog(err)
	return byteValue
}

func GetLinksJson() map[string]interface{} {
	var linksMap map[string]interface{}
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filepath.Dir(filename))
	jsonPath := filepath.Join(dir, "assets", "links", "links.json")

	jsonFileObject, err := os.Open(jsonPath)
	radioErrors.ErrorLog(err)
	defer jsonFileObject.Close()
	byteValue, _ := io.ReadAll(jsonFileObject)

	json.Unmarshal([]byte(byteValue), &linksMap)
	return linksMap
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

func GetDataFromResponse(url string) ([]byte, io.ReadCloser) {
	//Seperate this into two funcs one to store in IO and the other to that is used here.
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

func GetGameDataObjectFromLandingLinks(landingLinks []string) []models.GameData {
	var gameDataObjects []models.GameData
	for _, landingLink := range landingLinks {
		gameDataObjects = append(gameDataObjects, GetGameDataObject(landingLink))
	}
	return gameDataObjects
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
	gameLandingLink = strings.Replace(gameLandingLink, "landing", "right-rail", 1)
	byteValue, bodyCloser := GetDataFromResponse(gameLandingLink)
	err := json.Unmarshal(byteValue, versesData)
	radioErrors.ErrorLog(err)
	bodyCloser.Close()
	return *versesData
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
			_ = os.Remove(localpath)
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
	CreateTmpDirectory(filepath.Join(quickRadioTempDirectory, gameSubDirectory))
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

func TranscodeToWave(aacFilepath string) string {
	wavFilepath := strings.Replace(aacFilepath, ".aac", ".wav", 1)
	if !DoesFileExist(wavFilepath) {
		_ = ffmpeg_go.Input(aacFilepath).Output(wavFilepath).OverWriteOutput().Run()
	}
	return wavFilepath
}

func UpdateRadioWavs(qualityLink string) {
	aacPaths := GetAACPaths(qualityLink)
	GoDownloadAndTranscodeAACs(aacPaths)
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

func GetTestFileObject(desiredFilename string) *os.File {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filepath.Dir(filename))
	testFilePath := filepath.Join(dir, "assets", "tests", desiredFilename)
	fileObject, _ := os.Open(testFilePath)
	return fileObject
}

func CreateTmpDirectory(tmpDirectory string) {
	if _, err := os.Stat(tmpDirectory); err == nil {
		return
	} else if os.IsNotExist(err) {
		err := os.MkdirAll(tmpDirectory, 0777)
		radioErrors.ErrorLog(err)
		return
	} else {
		radioErrors.ErrorLog(err)
		return
	}
}

func DoesFileExistErr(filepath string) error {
	if _, err := os.Stat(filepath); err == nil {
		return nil
	} else if os.IsNotExist(err) {
		return errors.New("file does not exist")
	} else {
		radioErrors.ErrorLog(err)
		return err
	}
}

func DoesFileExist(filepath string) bool {
	if _, err := os.Stat(filepath); err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	} else {
		radioErrors.ErrorLog(err)
		return false
	}
}

func GetQuickTmpFolder() string {
	tempDirectory := os.TempDir()
	return filepath.Join(tempDirectory, "QuickRadio")
}

func GetActiveGameDirectory() string {
	activeGameDirectory := filepath.Join(GetQuickTmpFolder(), "ActiveGame")
	CreateTmpDirectory(activeGameDirectory)
	return activeGameDirectory
}

func EmptyTmpFolder() {
	tempDirectory := GetQuickTmpFolder()
	os.RemoveAll(tempDirectory)
}

func EmptyRadioDirectory(radioDirectory string) {
	os.RemoveAll(radioDirectory)
}

func EmptyActiveGameDirectory(activeGameDirectory string) {
	files, _ := os.ReadDir(activeGameDirectory)
	for _, f := range files {
		info, _ := f.Info()
		os.Remove(filepath.Join(activeGameDirectory, info.Name()))
	}
}

func GetProjectDir() string {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filepath.Dir(filename))
	return dir
}

func GetLogoPath(logoFilename string) string {
	dir := GetProjectDir()
	path := filepath.Join(dir, "assets", "svgs", "logos", logoFilename)
	return path
}

func GetLockpath(teamAbbrev string) string {
	filename := teamAbbrev + ".RADIOLOCK"
	return filepath.Join(GetQuickTmpFolder(), filename)
}

func TouchFile(filepath string) {
	f, err := os.Create(filepath)
	radioErrors.ErrorLog(err)
	f.Close()
}

func WriteFile(filepath string, data string) {
	f, err := os.Create(filepath)
	radioErrors.ErrorLog(err)
	f.WriteString(data)
	f.Close()
}

func CreateRadioLock(teamAbbrev string) {
	lockPath := GetLockpath(teamAbbrev)
	f, _ := os.Create(lockPath)
	f.Close()
}

func IsRadioLocked() bool {
	files, _ := os.ReadDir(GetQuickTmpFolder())
	for _, f := range files {
		info, _ := f.Info()
		if strings.HasSuffix(info.Name(), ".RADIOLOCK") {
			return true
		}
	}
	return false
}

func IsOurRadioLocked(teamAbbrev string) bool {
	lockPath := GetLockpath(teamAbbrev)
	return DoesFileExist(lockPath)
}

func DeleteRadioLock(teamAbbrev string) {
	log.Println("DeleteRadioLock", teamAbbrev)
	lockPath := GetLockpath(teamAbbrev)
	log.Println("DeleteRadioLock", lockPath)
	os.Remove(lockPath)
}

func DeleteAnyRadioLock() {
	tmpDirectory := GetQuickTmpFolder()
	files, _ := os.ReadDir(tmpDirectory)
	for _, f := range files {
		info, _ := f.Info()
		if strings.HasSuffix(info.Name(), ".RADIOLOCK") {
			os.Remove(path.Join(tmpDirectory, info.Name()))
		}
	}
}
