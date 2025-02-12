package quickio

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path"
	"path/filepath"
	"quickRadio/models"
	"quickRadio/radioErrors"
	"runtime"
	"strings"

	ffmpeg_go "github.com/u2takey/ffmpeg-go"
)

func TouchFile(path string) {
	CreateDirectory(filepath.Dir(path))
	f, err := os.Create(path)
	radioErrors.ErrorLog(err)
	f.Close()
}

func WriteFile(path string, data string) {
	CreateDirectory(filepath.Dir(path))
	f, err := os.Create(path)
	radioErrors.ErrorLog(err)
	f.WriteString(data)
	f.Close()
}

func GetProjectDir() string {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filepath.Dir(filepath.Dir(filename)))
	return dir
}

func GetProjectDirWithForwadSlash() string {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Dir(path.Dir(path.Dir(filename)))
	return dir
}

func GetQuickTmpFolder() string {
	tempDirectory := os.TempDir()
	return filepath.Join(tempDirectory, "QuickRadio")
}

func CreateDirectory(directoryPath string) {
	if _, err := os.Stat(directoryPath); err == nil {
		return
	} else if os.IsNotExist(err) {
		err := os.MkdirAll(directoryPath, 0777)
		radioErrors.ErrorLog(err)
		return
	} else {
		radioErrors.ErrorLog(err)
		return
	}
}

func EmptyActiveGameDirectory(activeGameDirectory string) {
	files, _ := os.ReadDir(activeGameDirectory)
	for _, f := range files {
		info, _ := f.Info()
		os.Remove(filepath.Join(activeGameDirectory, info.Name()))
	}
}

func EmptyTmpFolder() {
	tempDirectory := GetQuickTmpFolder()
	os.RemoveAll(tempDirectory)
}

func EmptyRadioDirectory(radioDirectory string) {
	os.RemoveAll(radioDirectory)
}

func GetActiveGameDirectory() string {
	activeGameDirectory := filepath.Join(GetQuickTmpFolder(), "ActiveGame")
	CreateDirectory(activeGameDirectory)
	return activeGameDirectory
}

func OpenAndGetFileHander(path string) *os.File {
	fileObject, err := os.Open(path)
	radioErrors.ErrorLog(err)
	return fileObject
}

func GetDataFromFile(path string) ([]byte, *os.File) {
	fileObject := OpenAndGetFileHander(path)
	byteValue, err := io.ReadAll(fileObject)
	radioErrors.ErrorLog(err)
	return byteValue, fileObject
}

func GetLinksJson() map[string]interface{} {
	var linksMap map[string]interface{}
	dir := GetProjectDir()
	jsonPath := filepath.Join(dir, "assets", "links", "links.json")

	jsonFileObject, err := os.Open(jsonPath)
	radioErrors.ErrorLog(err)
	defer jsonFileObject.Close()
	byteValue, _ := io.ReadAll(jsonFileObject)

	json.Unmarshal([]byte(byteValue), &linksMap)
	return linksMap
}

func GetSweaters() map[string]models.Sweater {
	sweaters := make(map[string]models.Sweater)
	dir := GetProjectDir()
	path := filepath.Join(dir, "assets", "teams", "sweater_colors.txt")
	fileObject := OpenAndGetFileHander(path)
	scanner := bufio.NewScanner(fileObject)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		var sweater models.Sweater
		lineColors := strings.Split(scanner.Text(), ";")
		sweater.TeamAbbrev = strings.ReplaceAll(lineColors[0], " ", "")
		sweater.PrimaryColor = strings.ReplaceAll(lineColors[1], " ", "")
		sweater.SecondaryColor = strings.ReplaceAll(lineColors[2], " ", "")
		sweaters[sweater.TeamAbbrev] = sweater
	}
	fileObject.Close()
	return sweaters
}

func TranscodeToWave(aacFilepath string) string {
	wavFilepath := strings.Replace(aacFilepath, ".aac", ".wav", 1)
	if !DoesFileExist(wavFilepath) {
		_ = ffmpeg_go.Input(aacFilepath).Output(wavFilepath).OverWriteOutput().Run()
	}
	return wavFilepath
}

func GetTestFileObject(desiredFilename string) *os.File {
	dir := GetProjectDir()
	testFilePath := filepath.Join(dir, "assets", "tests", desiredFilename)
	fileObject, _ := os.Open(testFilePath)
	return fileObject
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

func GetLogoPath(logoFilename string) string {
	dir := GetProjectDir()
	path := filepath.Join(dir, "assets", "svgs", "logos", logoFilename)
	return path
}
