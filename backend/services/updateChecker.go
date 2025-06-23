package services

import (
	"bili-audio-downloader/backend/constants"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/tidwall/gjson"
)

// CheckUpdate 通过 GitHub 检查程序更新
// string 为 "0" 代表没有更新，有更新时该位为最新版本号
func CheckUpdate(currentVersion string) (string, error) {
	// If build with a -tag, don't check update
	if hasTag(constants.APP_VERSION) {
		return "-1", nil
	}

	// Get current version tag from github
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", "HIM049/BADownloaderUI")
	resp, err := http.Get(url)
	if err != nil {
		return "0", errors.New(fmt.Sprintln("failed to check update:", err))
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return "0", errors.New(fmt.Sprintln("failed to send request:", resp.StatusCode))
	}

	// Read response and get json data
	bodyString, err := io.ReadAll(resp.Body)
	if err != nil {
		return "0", errors.New(fmt.Sprintln("failed to read response:", err))
	}
	latestVersion := gjson.Get(string(bodyString), "tag_name").String()

	// Compare version number
	if latestVersion > currentVersion {
		return latestVersion, nil
	}
	return "0", nil
}

// Check whether build with a -tag version
func hasTag(s string) bool {
	re := regexp.MustCompile(`-\w+$`)
	return re.MatchString(s)
}
