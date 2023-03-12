package main

import (
	_ "embed"
	"fmt"
	"os"
	"runtime"
	"splinterforge/ColorPrint"
	"splinterforge/LogFunc"
	"splinterforge/ProcedureFunc"
	"splinterforge/ReadFunc"
	"splinterforge/SpStruct"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/selenium-Driver-Check/SeleniumDriverCheck"
	"github.com/tebeka/selenium"
	"github.com/theckman/yacspin"
)

var (
	accountLists = []SpStruct.UserData{}
	r            = &sync.WaitGroup{}
	w            = &sync.WaitGroup{}
	s            = &sync.WaitGroup{}
	//headless, threadingLimit, showForgeReward, showAccountDetails, autoSelectCard, autoSelectHero, autoSelectSleepTime, splinterforgeAPIEndpoint, splinterlandAPIEndpoint, publicAPIEndpoint = ReadFunc.GetConfig("./config/config.txt")
	headless                 = false
	threadingLimit           = 0
	showForgeReward          = false
	showAccountDetails       = false
	autoSelectCard           = false
	autoSelectHero           = false
	autoSelectSleepTime      = false
	splinterforgeAPIEndpoint = ""
	splinterlandAPIEndpoint  = ""
	publicAPIEndpoint        = ""
	RealPath                 = ""

	ConfigAccountsPath    = ""
	ConfigCardSettingPath = ""
)

//检测操作系统，以区分读取文件路径
func init() {
	PcPlatForm := runtime.GOOS
	if PcPlatForm == "windows" {
		ConfigAccountsPath = "config/accounts.txt"
		ConfigCardSettingPath = "config/cardSettings.txt"
		headless, threadingLimit, showForgeReward, showAccountDetails, autoSelectCard, autoSelectHero, autoSelectSleepTime, splinterforgeAPIEndpoint, splinterlandAPIEndpoint, publicAPIEndpoint = ReadFunc.GetConfig("config/config.txt")
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
		ConfigAccountsPath = RealPath + "/config/accounts.txt"
		ConfigCardSettingPath = RealPath + "/config/cardSettings.txt"
		headless, threadingLimit, showForgeReward, showAccountDetails, autoSelectCard, autoSelectHero, autoSelectSleepTime, splinterforgeAPIEndpoint, splinterlandAPIEndpoint, publicAPIEndpoint = ReadFunc.GetConfig(RealPath + "/config/config.txt")

	}
}

func initializeAccount(accountNo int) (string, string, string, string, []SpStruct.CardSelection, int) {
	userName, postingKey, err := ReadFunc.GetAccountData(ConfigAccountsPath, accountNo)
	if err != nil || userName == "" || postingKey == "" {
		ColorPrint.PrintRed("ERROR", "Error in loading accounts.txt, please add username or posting key and try again.")
	}
	heroesType, bossId, playingSummoners, playingMonster, timeSleepInMinute, err := ReadFunc.GetCardSettingData(ConfigCardSettingPath, accountNo)
	if heroesType == "" || bossId == "" {
		fmt.Println("")
		ColorPrint.PrintRed("SF", fmt.Sprintf("Error loading cardSettings.txt file for account %s", strconv.Itoa(accountNo)))
		ColorPrint.PrintWhite("SF", "Terminating in 10 seconds...")
		time.Sleep(10 * time.Second)
		os.Exit(1)
	}
	if err != nil {
		ColorPrint.PrintRed("ERROR", "Error loading cardSettings.txt file")
	}
	playingSummonersList := make([]SpStruct.Summoners, 0, len(playingSummoners))
	playingMonsterList := make([]SpStruct.MonsterId, 0, len(playingMonster))
	for _, i := range playingSummoners {
		cardName, _ := ReadFunc.GetCardName(i)
		playingSummonersList = append(playingSummonersList, SpStruct.Summoners{
			PlayingSummonersDiv:  fmt.Sprintf("//div/img[@id='%s']", i),
			PlayingSummonersID:   i,
			PlayingSummonersName: cardName,
		})
	}
	for _, i := range playingMonster {
		cardName, _ := ReadFunc.GetCardName(i)
		playingMonsterList = append(playingMonsterList, SpStruct.MonsterId{
			PlayingMontersDiv:   fmt.Sprintf("//div/img[@id='%s']", i),
			PlayingMonstersID:   i,
			PlayingMonstersName: cardName,
		})
	}
	var cardSelectionList = []SpStruct.CardSelection{}
	cardSelection := SpStruct.CardSelection{
		PlayingSummoners: playingSummonersList,
		PlayingMonsters:  playingMonsterList,
	}
	cardSelectionList = append(cardSelectionList, cardSelection)

	return userName, postingKey, heroesType, bossId, cardSelectionList, timeSleepInMinute
}
func initializeUserData() {
	LogFunc.PrintInfo()
	cfg := yacspin.Config{
		Frequency:       100 * time.Millisecond,
		CharSet:         yacspin.CharSets[59],
		Suffix:          " Initialization Driver",
		SuffixAutoColon: true,
		StopCharacter:   " pass",
		StopColors:      []string{"fgGreen"},
		StopMessage:     "Success!",
	}
	spinner, _ := yacspin.New(cfg)
	spinner.Start()
	spinner.Message("Checking...")
	wd, err := selenium.NewChromeDriverService(SeleniumDriverCheck.AutoDownload_ChromeDriver(false), 9515)
	if err != nil {
		fmt.Printf("Failed to create ChromeDriver service: %s\n", err)
		os.Exit(1)
	}
	defer wd.Stop()
	spinner.Stop()
	cfg = yacspin.Config{
		Frequency:       100 * time.Millisecond,
		CharSet:         yacspin.CharSets[59],
		Suffix:          " Reading config folder",
		SuffixAutoColon: true,
		StopCharacter:   " pass",
		StopColors:      []string{"fgGreen"},
		StopMessage:     "Success!",
	}
	spinner, _ = yacspin.New(cfg)
	spinner.Start()
	spinner.Message("reading accounts.txt...")
	lineCount, errCountLines := ReadFunc.GetLines(ConfigAccountsPath)
	time.Sleep(500 * time.Millisecond)
	if errCountLines == nil && lineCount > 1 {
		spinner.Message("reading cardSettings.txt...")
		for i := 0; i < lineCount-1; i++ {
			r.Add(1)
			time.Sleep(50 * time.Millisecond)
			go func(num int) {
				userName, postingKey, heroesType, bossId, cardSelection, timeSleepInMinute := initializeAccount(num + 1)
				accountLists = append(accountLists, SpStruct.UserData{
					UserName:          userName,
					PostingKey:        postingKey,
					BossID:            bossId,
					HeroesType:        heroesType,
					CardSelection:     cardSelection,
					TimeSleepInMinute: timeSleepInMinute,
				})
				r.Done()
			}(i)
		}
		r.Wait()
		spinner.Message("reading config.txt..")
		time.Sleep(500 * time.Millisecond)
		spinner.Stop()
		LogFunc.PrintConfigSettings(lineCount-1, headless, threadingLimit, showForgeReward, showAccountDetails, autoSelectCard, autoSelectHero, autoSelectSleepTime)
		for i := 0; i < len(accountLists); i += threadingLimit {
			if len(accountLists)-i < threadingLimit {
				threadingLimit = len(accountLists) - i
			}
			for j := 0; j < threadingLimit; j++ {
				w.Add(1)
				go ProcedureFunc.InitializeDriver(true, accountLists[i+j], headless, showForgeReward, showAccountDetails, autoSelectCard, autoSelectHero, autoSelectSleepTime, splinterforgeAPIEndpoint, splinterlandAPIEndpoint, publicAPIEndpoint, accountLists, s, w)
			}
			w.Wait()
			LogFunc.PrintInfo()
		}
		s.Wait()
	} else {
		fmt.Println("")
		ColorPrint.PrintRed("SF", "Please add accounts in accounts.txt")
		ColorPrint.PrintWhite("SF", "Terminating in 10 seconds...")
		time.Sleep(10 * time.Second)
		os.Exit(1)
	}
}
func main() {
	initializeUserData()

}
