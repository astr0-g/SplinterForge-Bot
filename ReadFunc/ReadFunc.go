package ReadFunc

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"splinterforge/ColorPrint"
	"splinterforge/SpStruct"
	"strconv"
	"strings"
	"time"

	"github.com/levigross/grequests"
)

var (
	PcPlatForm = runtime.GOOS

	DataCardMappingPath = ""
	RealPath            = ""
)

func init() {
	if PcPlatForm == "windows" {
		DataCardMappingPath = "data/cardMapping.json"
	} else if PcPlatForm == "darwin" {
		//获取当前文件夹
		path, err := os.Executable()
		if err != nil {
			panic(err)
		}
		if strings.Contains(path, "private") || strings.Contains(path, "___go_build") || strings.Contains(path, "folders/") {
			RealPath, _ = os.Getwd()
		} else {
			RealPathLists := strings.Split(path, "/")
			RealPath = strings.Join(RealPathLists[:len(RealPathLists)-1], "/")
		}
		DataCardMappingPath = RealPath + "/data/cardMapping.json"
	}
}

func GetCardName(cardId string) (string, error) {
	file, err := os.Open(DataCardMappingPath)
	if err != nil {
		return "", fmt.Errorf("error opening JSON file: %w", err)
	}
	defer file.Close()

	var cards []map[string]string
	err = json.NewDecoder(file).Decode(&cards)
	if err != nil {
		return "", fmt.Errorf("error decoding JSON file: %w", err)
	}

	for _, card := range cards {
		if id, ok := card[cardId]; ok {
			return id, nil
		}
	}

	return "", fmt.Errorf("card id %s not found", cardId)
}

func GetAccountData(filePath string, lineNumber int) (string, string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for i := 1; scanner.Scan(); i++ {
		if i == 1 {
			continue
		}
		if i == lineNumber+1 {
			acctinfo := scanner.Text()
			time.Sleep(1 * time.Second)
			userName := strings.Split(acctinfo, ":")[0]
			postingKey := strings.Split(acctinfo, ":")[1]
			return userName, postingKey, nil
		}
	}
	if err := scanner.Err(); err != nil {

		return "", "", err
	}
	return "", "", nil
}

func GetCardSettingData(filePath string, lineNumber int) (string, string, []string, []string, int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", "", nil, nil, 0, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for i := 1; scanner.Scan(); i++ {
		if i == lineNumber+1 {
			acctinfo := scanner.Text()
			heroesType := strings.Split(acctinfo, ":")[0]
			bossId := strings.Split(acctinfo, ":")[1]
			playingSummoners := strings.Split(strings.Split(acctinfo, ":")[2], ",")
			playingMonster := strings.Split(strings.Split(acctinfo, ":")[3], ",")
			timeSleepInMinute, err := strconv.Atoi(strings.Split(acctinfo, ":")[4])
			if err != nil {
				return "", "", nil, nil, 0, err
			}
			return heroesType, bossId, playingSummoners, playingMonster, timeSleepInMinute, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", "", nil, nil, 0, err
	}

	return "", "", nil, nil, 0, nil
}

func GetConfig(filePath string) (bool, int, bool, bool, bool, bool, bool,bool, bool, string, string, string) {
	file, err := os.Open(filePath)
	if err != nil {
		ColorPrint.PrintRed("SF", "Error Reading Config.txt file")
		ColorPrint.PrintWhite("SF", "Terminating in 10 seconds...")
		time.Sleep(10 * time.Second)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var (
		headless                 bool
		threadingLimit           int
		showForgeReward          bool
		showAccountDetails       bool
		autoSelectCard           bool
		autoSelectHero           bool
		autoSelectSleepTime      bool
		waitForBossRespawn		 bool
		shareBattleLog           bool
		splinterforgeAPIEndpoint string
		splinterlandAPIEndpoint  string
		publicAPIEndpoint        string
	)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "=") {
			parts := strings.Split(line, "=")
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			switch key {
			case "HEADLESS":
				headless = value == "true"
			case "THREADING":
				threadingLimit, _ = strconv.Atoi(value)
			case "SHOW_FORGE_REWARD":
				showForgeReward = value == "true"
			case "SHOW_ACCOUNT_DETAILS":
				showAccountDetails = value == "true"
			case "AUTO_WAIT_FOR_BOSS_RESPAWN":
				waitForBossRespawn = value == "true"
			case "SHARE_BATTLE_LOG":
				shareBattleLog = value == "true"
			case "AUTO_SELECT_CARD":
				autoSelectCard = value == "true"
			case "AUTO_SELECT_SLEEPTIME":
				autoSelectSleepTime = value == "true"
			case "AUTO_SELECT_HERO":
				autoSelectHero = value == "true"
			case "SPLINTERFORGE_API_ENDPOINT":
				splinterforgeAPIEndpoint = value
			case "SPLINTERLAND_API_ENDPOINT":
				splinterlandAPIEndpoint = value
			case "PUBLIC_API_ENDPOINT":
				publicAPIEndpoint = value
			}
		}
	}
	if threadingLimit == 0 || splinterforgeAPIEndpoint == "" || splinterlandAPIEndpoint == "" || publicAPIEndpoint == "" {
		ColorPrint.PrintRed("SF", "Error reading config.txt file.")
		ColorPrint.PrintWhite("SF", "Terminating in 10 seconds...")
		time.Sleep(10 * time.Second)
		os.Exit(1)
	}
	return headless, threadingLimit, showForgeReward, showAccountDetails, autoSelectCard, autoSelectHero, autoSelectSleepTime, waitForBossRespawn,shareBattleLog, splinterforgeAPIEndpoint, splinterlandAPIEndpoint, publicAPIEndpoint
}

func GetLines(filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}
	return lineCount, nil
}

func GetAccountDetails(name interface{}, key interface{}, splinterforgeAPIEndpoint string) int {
	res, _ := grequests.Post(fmt.Sprintf("%s/users/keyLogin", splinterforgeAPIEndpoint), &grequests.RequestOptions{
		JSON: map[string]string{
			"username": name.(string),
			"key":      key.(string),
		},
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	})
	var powerRes = SpStruct.KeyLoginResData{}
	json.Unmarshal(res.Bytes(), &powerRes)
	reusltTime, _ := GetTimeDiff(powerRes.Stamina.Last)
	CurrentStamina := powerRes.Stamina.Current + reusltTime
	if CurrentStamina > powerRes.Stamina.Max {
		CurrentStamina = powerRes.Stamina.Max
	}
	return CurrentStamina
}

func GetTimeDiff(oldTime string) (int, error) {
	now := time.Now()
	t, err := time.Parse(time.RFC3339, oldTime)
	if err != nil {
		return 0, err
	}
	diff := int(now.Unix() - t.Unix())
	diffInMinutes := diff / 60
	return diffInMinutes, nil
}
