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

func GetGameHtml(linksMap map[string]interface{}) string {
	var html string
	sleepTimer := linksMap["load_sleep_timer"].(float64)
	baseUrl := fmt.Sprintf("%v", linksMap["base"])

	ctx, cancel := chromedp.NewContext(context.Background())
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
	radioErrors.ErrorCheck(err)
	return html
}

func GetSweaterColors() map[string][]string {
	sweaterColors := make(map[string][]string)
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filepath.Dir(filename))
	path := filepath.Join(dir, "assets", "teams", "sweater_colors.txt")

	fileObject, err := os.Open(path)
	radioErrors.ErrorCheck(err)
	defer fileObject.Close()
	scanner := bufio.NewScanner(fileObject)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		lineColors := strings.Split(scanner.Text(), ";")
		sweaterColors[lineColors[0]] = lineColors[1:]
	}

	return sweaterColors
}

func GetLinksJson() map[string]interface{} {
	var linksMap map[string]interface{}
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filepath.Dir(filename))
	jsonPath := filepath.Join(dir, "assets", "links", "links.json")

	jsonFileObject, err := os.Open(jsonPath)
	radioErrors.ErrorCheck(err)
	defer jsonFileObject.Close()
	byteValue, _ := io.ReadAll(jsonFileObject)

	json.Unmarshal([]byte(byteValue), &linksMap)
	return linksMap
}

func GetGameLandingLinks(html string, gamecenterBase string, gamecenterLanding string, gameRegex string) ([]string, error) {
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
	if len(gameLandingLinks) == 0 {
		return nil, errors.New("couldnt find any game for today")
	}
	return gameLandingLinks, nil
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

func GetDataFromResponse(url string) []byte {
	//Seperate this into two funcs one to store in IO and the other to that is used here.
	resp, err := http.Get(url)
	radioErrors.ErrorCheck(err)
	defer resp.Body.Close()
	byteValue, err := io.ReadAll(resp.Body)
	radioErrors.ErrorCheck(err)
	return byteValue

}

func GetQualityStreamSlug(radioLink string, audioQuality string) string {
	byteValue := GetDataFromResponse(radioLink)
	audioQualitySlug, err := ExtractQualityStreamSlug(string(byteValue), audioQuality)
	radioErrors.ErrorCheck(err)
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
	audioFilepaths := []string{}
	resp, err := http.Get(qualityRadioPath)
	radioErrors.ErrorCheck(err)
	defer resp.Body.Close()
	byteValue, err := io.ReadAll(resp.Body)
	radioErrors.ErrorCheck(err)
	audioSlugs, err := GetAACSlugsFromQualityFile(string(byteValue))
	for _, audioSlug := range audioSlugs {
		audioFilepaths = append(audioFilepaths, BuildAACRadioPath(qualityRadioPath, audioSlug))
	}
	radioErrors.ErrorCheck(err)
	return audioFilepaths
}

func DownloadAndTranscodeAACs(paths []string) []string {
	log.Println("FILE IO - FUNC - DownloadAndTranscodeAACs")
	wavpaths := make([]string, len(paths))
	var workGroup sync.WaitGroup
	for i := 0; i < len(paths); i++ {
		workGroup.Add(1)
		go func(path string) {
			defer workGroup.Done()
			localpath, err := DownloadAAC(path)
			radioErrors.ErrorCheck(err)
			wavPath := TranscodeToWave(localpath)
			wavpaths = append(wavpaths, wavPath)
		}(paths[i])
		workGroup.Wait()
	}
	log.Println("VAR - wavpaths", wavpaths)
	return wavpaths
}

func DownloadAAC(aacRequestPath string) (string, error) {
	log.Println("FUNC - DownloadAAC")
	quickRadioTempDirectory := GetQuickTmpFolder()
	filename := strings.Split(aacRequestPath, "/")[len(strings.Split(aacRequestPath, "/"))-1]
	gameSubDirectory := filepath.Join(strings.Split(aacRequestPath, "/")[4 : len(strings.Split(aacRequestPath, "/"))-2]...)
	CreateTmpDirectory(filepath.Join(quickRadioTempDirectory, gameSubDirectory))
	filepath := filepath.Join(quickRadioTempDirectory, gameSubDirectory, filename)
	log.Println("VAR - quickRadioTempDirectory", quickRadioTempDirectory)
	log.Println("VAR - filename", filename)
	log.Println("VAR - gameSubDirectory", gameSubDirectory)
	log.Println("VAR - filepath", filepath)
	log.Println("VAR - aacRequestPath", aacRequestPath)
	if runtime.GOOS == "windows" {
		if !DoesFileExist(filepath) {
			cmd := exec.Command("curl", "-o", filepath, aacRequestPath)
			_, err := cmd.CombinedOutput()
			radioErrors.ErrorCheck(err)
			err = DoesFileExistErr(filepath)
			radioErrors.ErrorCheck(err)
		}
	} else if runtime.GOOS == "linux" {
		if !DoesFileExist(filepath) {
			cmd := exec.Command("wget", aacRequestPath, filepath)
			_, err := cmd.CombinedOutput()
			radioErrors.ErrorCheck(err)
			err = DoesFileExistErr(filepath)
			radioErrors.ErrorCheck(err)
		}
	} else {
		if !DoesFileExist(filepath) {
			cmd := exec.Command("wget", aacRequestPath, filepath)
			_, err := cmd.CombinedOutput()
			radioErrors.ErrorCheck(err)
			err = DoesFileExistErr(filepath)
			radioErrors.ErrorCheck(err)
		}
	}
	return filepath, nil
}

func TranscodeToWave(aacFilepath string) string {
	wavFilepath := strings.Replace(aacFilepath, ".aac", ".wav", 1)
	if !DoesFileExist(wavFilepath) {
		err := ffmpeg_go.Input(aacFilepath).Output(wavFilepath).OverWriteOutput().ErrorToStdOut().Run()
		radioErrors.ErrorCheck(err)
	}
	return wavFilepath
}

func UpdateRadioWavs(qualityLink string) {
	aacPaths := GetAACPaths(qualityLink)
	DownloadAndTranscodeAACs(aacPaths)
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
		radioErrors.ErrorCheck(err)
		return
	} else {
		radioErrors.ErrorCheck(err)
		return
	}
}

func DoesFileExistErr(filepath string) error {
	if _, err := os.Stat(filepath); err == nil {
		return nil
	} else if os.IsNotExist(err) {
		return errors.New("file does not exist")
	} else {
		radioErrors.ErrorCheck(err)
		return err
	}
}

func DoesFileExist(filepath string) bool {
	if _, err := os.Stat(filepath); err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	} else {
		radioErrors.ErrorCheck(err)
		return false
	}
}

func GetQuickTmpFolder() string {
	log.Println("FUNC - GetQuickTmpFolder")
	tempDirectory := os.TempDir()
	return filepath.Join(tempDirectory, "QuickRadio")
}

func EmptyTmpFolder() {
	log.Println("FUNC - EmptyTmpFolder")
	tempDirectory := GetQuickTmpFolder()
	log.Println("Removing Temp Directory ", tempDirectory)
	os.RemoveAll(tempDirectory)
}

func CreateRadioLock() {
	lockPath := path.Join(GetQuickTmpFolder(), "LOCK.RADIO")
	os.Create(lockPath)
}

func IsRadioLocked() bool {
	lockPath := path.Join(GetQuickTmpFolder(), "LOCK.RADIO")
	return DoesFileExist(lockPath)
}

func DeleteRadioLock() {
	lockPath := path.Join(GetQuickTmpFolder(), "LOCK.RADIO")
	DoesFileExistErr(lockPath)
	os.Remove(lockPath)
}
