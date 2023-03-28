package ProcedureFunc

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"splinterforge/DriverAction"
	"splinterforge/LogFunc"
	"splinterforge/SpStruct"
	"strconv"
	"strings"
	"sync"
	"time"
	"github.com/levigross/grequests"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"github.com/tebeka/selenium/log"

	"splinterforge/ColorPrint"
	"splinterforge/GameFunc"
	"splinterforge/ReadFunc"
	"splinterforge/RequestFunc"
)

var (
	PcPlatForm    = runtime.GOOS
	ExtensionPath = ""
	RealPath      = ""
)


func init() {
	if PcPlatForm == "windows" {
		ExtensionPath = "data/hivekeychain.crx"
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
		ExtensionPath = RealPath + "/data/hivekeychain.crx"
	}
}

func AccountLogin(userName string, postingKey string, wd selenium.WebDriver) bool {
	err := wd.SetImplicitWaitTimeout(5 * time.Second)
	if err != nil {
		panic(err)
	}
	DriverAction.DriverGet("chrome-extension://jcacnejopjdphbnjgfaaobbfafkihpep/popup.html", wd)

	DriverAction.DriverElementWaitAndClick(wd, "/html/body/div/div/div[4]/div[2]/div[5]/button")

	el, err := wd.FindElement(selenium.ByXPATH, "/html/body/div/div/div[1]/div/div[1]/div/input")
	if err != nil {
		ColorPrint.PrintRed(userName, "Login Error!")
		return false
	}
	err = el.SendKeys("Aa123Aa123!!")
	if err != nil {
		ColorPrint.PrintRed(userName, "Login Error!")
		return false
	}
	el, err = wd.FindElement(selenium.ByXPATH, "/html/body/div/div/div[1]/div/div[2]/div/input")
	if err != nil {
		ColorPrint.PrintRed(userName, "Login Error!")
		return false
	}
	err = el.SendKeys("Aa123Aa123!!")
	if err != nil {
		ColorPrint.PrintRed(userName, "Login Error!")
		return false
	}
	el, err = wd.FindElement(selenium.ByXPATH, "/html/body/div/div/div[1]/button/div")
	if err != nil {
		ColorPrint.PrintRed(userName, "Login Error!")
		return false
	}
	err = el.Click()
	if err != nil {
		ColorPrint.PrintRed(userName, "Login Error!")
		return false
	}
	el, err = wd.FindElement(selenium.ByXPATH, "/html/body/div/div/div[1]/div[2]/div/div[2]/button[1]/div")
	if err != nil {
		ColorPrint.PrintRed(userName, "Login Error!")
		return false
	}
	err = el.Click()
	if err != nil {
		ColorPrint.PrintRed(userName, "Login Error!")
		return false
	}
	el, err = wd.FindElement(selenium.ByXPATH, "/html/body/div/div/div[1]/div[2]/div/div[2]/div[1]/div/input")
	if err != nil {
		ColorPrint.PrintRed(userName, "Login Error!")
		return false
	}
	err = el.SendKeys(userName)
	if err != nil {
		ColorPrint.PrintRed(userName, "Login Error!")
		return false
	}
	el, err = wd.FindElement(selenium.ByXPATH, "/html/body/div/div/div[1]/div[2]/div/div[2]/div[2]/div/input")
	if err != nil {
		ColorPrint.PrintRed(userName, "Login Error!")
		return false
	}
	err = el.SendKeys(postingKey)
	if err != nil {
		ColorPrint.PrintRed(userName, "Login Error!")
		return false
	}
	el, err = wd.FindElement(selenium.ByXPATH, "/html/body/div/div/div[1]/div[2]/div/div[2]/div[2]/div/input")
	if err != nil {
		ColorPrint.PrintRed(userName, "Login Error!")
		return false
	}
	time.Sleep(2 * time.Second)
	err = el.SendKeys(selenium.ReturnKey)
	if err != nil {
		ColorPrint.PrintRed(userName, "Login Error!")
		return false
	}
	el, err = wd.FindElement(selenium.ByXPATH, "/html/body/div/div/div[4]/div")
	if err != nil {
		ColorPrint.PrintRed(userName, "Login failure! Please check your accounts.txt file or it is a server error.")
		return false
	}
	hiveKeyChainLogin, err := el.Text()
	if err != nil {
		ColorPrint.PrintRed(userName, "Login Error!")
		return false
	}
	if hiveKeyChainLogin != "HIVE KEYCHAIN" {
		ColorPrint.PrintRed(userName, "Login failure! Please check your accounts.txt file or it is a server error.")
		return false
	} else {
		err = wd.ResizeWindow("bigger", 1565, 1080)
		if err != nil {
			println("can not change size")
		}
		DriverAction.DriverGet("https://splinterforge.io/#/", wd)
		GameFunc.CheckPopUp(wd, 1000)
		el, err = wd.FindElement(selenium.ByXPATH, "/html/body/app/div[1]/div[1]/app-header/success-modal/section/div[1]/div[4]/div/button")
		if err != nil {
			ColorPrint.PrintRed(userName, "Login Error!")
			return false
		}
		el.Click()
		el, err = wd.FindElement(selenium.ByXPATH, "/html/body/app/div[1]/div[1]/app-header/section/div[4]/div[2]/div/div/a/div[1]")
		if err != nil {
			ColorPrint.PrintRed(userName, "Login Error!")
			return false
		}
		el.Click()
		el, err = wd.FindElement(selenium.ByXPATH, "/html/body/app/div[1]/login-modal/div/div/div/div[2]/div[2]/input")
		if err != nil {
			ColorPrint.PrintRed(userName, "Login Error!")
			return false
		}
		el.SendKeys(userName)
		ellogin, err := wd.FindElement(selenium.ByXPATH, "/html/body/app/div[1]/login-modal/div/div/div/div[2]/div[3]/button")
		if err != nil {
			ColorPrint.PrintRed(userName, "Login Error!")
			return false
		}
		CheckloginButton(wd)
		ellogin.Click()
		for {
			handles, _ := wd.WindowHandles()
			if len(handles) == 2 {
				break
			}
		}
		handles, _ := wd.WindowHandles()
		wd.SwitchWindow(handles[1])
		el, _ = wd.FindElement(selenium.ByXPATH, "/html/body/div/div/div/div[3]/div[1]/div/div")
		el.Click()
		el, _ = wd.FindElement(selenium.ByXPATH, "/html/body/div/div/div/div[3]/div[2]/button[2]/div")
		el.Click()
		wd.SwitchWindow(handles[0])
		Checklogin(userName, wd)
		ColorPrint.PrintGreen(userName, "Login successful!")
		return true
	}
}

func CheckloginButton(wd selenium.WebDriver) {
	defer func() {
		if r := recover(); r != nil {
		}
	}()
	wd.SetImplicitWaitTimeout(1 * time.Second)
	el, _ := wd.FindElement(selenium.ByXPATH, "/html/body/app/div[1]/login-modal/div/div/div/div[2]/div[3]/div/input")
	el.Click()
	// ellogin, _ := wd.FindElement(selenium.ByXPATH, "/html/body/app/div[1]/login-modal/div/div/div/div[2]/div[3]/button")
	// ellogin.Click()
}

func Checklogin(userName string, wd selenium.WebDriver) bool {
	for i := 0; i < 10; i++ {
		el, err := wd.FindElement(selenium.ByXPATH, "/html/body/app/div[1]/div[1]/app-header/section/div[4]/div[2]/div[2]/div[2]/a/div[2]")
		if err != nil {
			time.Sleep(3 * time.Second)
			continue
		}
		text, _ := el.Text()
		if text == userName {
			return true
		}
		time.Sleep(1 * time.Second)
	}
	return false
}

func AccountBattle(wait bool, wd selenium.WebDriver, userName string, bossId string, headless bool, heroesType string, timeSleepInMinute int, cardSelection []SpStruct.CardSelection, autoSelectHero bool, autoSelectCard bool, autoSelectSleepTime bool, splinterlandAPIEndpoint string, publicAPIEndpoint string, splinterforgeAPIEndpoint string, showForgeReward bool, showAccountDetails bool, waitForBossRespawn bool, shareBattleLog bool, s *sync.WaitGroup, w *sync.WaitGroup, accountLists []SpStruct.UserData) {
	starttimestamp := time.Now().Unix()
	CookiesStatus := true
	Unexpected := false
	name, _ := wd.ExecuteScript("return localStorage.getItem('forge:username');", nil)
	key, _ := wd.ExecuteScript("return localStorage.getItem('forge:key');", nil)
	bossName, _, err := GameFunc.SelectBoss(userName, bossId, waitForBossRespawn, wd)
	if err == nil {
		cardSelection, _, _ = GameFunc.SelectCards(cardSelection, bossName, userName, splinterlandAPIEndpoint, publicAPIEndpoint, autoSelectCard)
		seletedNumOfSummoners := 1
		el, _ := wd.FindElement(selenium.ByXPATH, "/html/body/app/div[1]/slcards/div[5]/section[1]/div/div[1]/div[2]/button")
		el.Click()
		ColorPrint.PrintWhite(userName, "Participating in battles...")
		printData := [][]string{}
		selectResult := true
		bossLeague, bossAbilities, bossRandomAbilities := RequestFunc.FetchBossAbilities(userName, key.(string), bossName, splinterforgeAPIEndpoint)
		heroTypechoosed := GameFunc.SelectHero(heroesType, name.(string), key.(string), bossRandomAbilities, wd, autoSelectHero, publicAPIEndpoint, bossName, splinterforgeAPIEndpoint)
		for _, selection := range cardSelection {
			for i, playingSummoner := range selection.PlayingSummoners {
				result := GameFunc.SelectSummoners(userName, seletedNumOfSummoners, playingSummoner.PlayingSummonersDiv, wd)
				if result {
					seletedNumOfSummoners++
					printData = append(printData, []string{fmt.Sprintf("Summoners #%d", i+1), playingSummoner.PlayingSummonersName, playingSummoner.PlayingSummonersID, "success"})
				}
			}
			seletedNumOfMonsters := 1
			for j, playingMonster := range selection.PlayingMonsters {
				result := GameFunc.SelectMonsters(userName, seletedNumOfMonsters, playingMonster.PlayingMontersDiv, wd)
				if result {
					seletedNumOfMonsters++
					printData = append(printData, []string{fmt.Sprintf("Monsters #%d", j+1), playingMonster.PlayingMonstersName, playingMonster.PlayingMonstersID, "success"})
				} else {
					printData = append(printData, []string{fmt.Sprintf("Monsters #%d", j+1), playingMonster.PlayingMonstersName, playingMonster.PlayingMonstersID, "error"})
					selectResult = false
				}
			}
		}
		ReadFunc.GetAccountDetails(name, key, splinterforgeAPIEndpoint)
		el, _ = wd.FindElement(selenium.ByXPATH, "/html/body/app/div[1]/slcards/div[5]/div[2]/div[1]/div[1]/button/span")
		mana, _ := el.Text()
		manaused, _ := strconv.Atoi(mana)
		CurrentStamina := ReadFunc.GetAccountDetails(name, key, splinterforgeAPIEndpoint)
		if CurrentStamina > manaused && manaused >= 30 {
			LogFunc.PrintResultBox(userName, printData, selectResult, bossName, bossLeague, heroTypechoosed, bossAbilities, bossRandomAbilities)
			battletimestamp := time.Now().Unix()
			if battletimestamp-starttimestamp < 30 {
				time.Sleep(time.Duration(starttimestamp+30-battletimestamp) * time.Second)
			}
			returnJsonResult := false
			postJsonResult := !showForgeReward
			fetchTime := 0
			fitPostData := SpStruct.FitBossPostData{}
			fitRes := SpStruct.FitBossRequestsData{}
			for {
				el, err := wd.FindElement(selenium.ByXPATH, "/html/body/app/div[1]/slcards/div[5]/button[1]/div[2]/span")
				if err == nil {
					el.Click()
				} else {
					netLogs, _ := wd.Log("performance")
					for _, netLog := range netLogs {
						if returnJsonResult == false && strings.Contains(netLog.Message, fmt.Sprintf("%s/boss/fight_boss", splinterforgeAPIEndpoint)) && strings.Contains(netLog.Message, "\"method\":\"Network.requestWillBeSent\"") {
							json.Unmarshal([]byte(netLog.Message), &fitRes)
							json.Unmarshal([]byte(fitRes.Message.Params.Request.PostData), &fitPostData)
							ColorPrint.PrintGreen(userName, "Battle was successful!")
							returnJsonResult = true
						}
						if showForgeReward == true && postJsonResult == false && strings.Contains(netLog.Message, fmt.Sprintf("%s/boss/fight_boss", splinterforgeAPIEndpoint)) && strings.Contains(netLog.Message, "\"method\":\"Network.responseReceived\"") {
							defer func() {
								if r := recover(); r != nil {
									ColorPrint.PrintRed(userName, "Encountering difficulty in reading the game results, but the battle has ended.")
								}
							}()
							var GetResponseBody = SpStruct.GetResponseBody{}
							var GetRewardBody = SpStruct.CDPFitReturnData{}
							resString, err := RequestFunc.GetReponseBody(wd.SessionID(), fitRes.Message.Params.RequestID, userName)
							if err == nil && resString != "" {
								json.Unmarshal([]byte(resString), &GetResponseBody)
								json.Unmarshal([]byte(GetResponseBody.Value.Body), &GetRewardBody)
								if GetRewardBody.TotalDmg >= 0 && GetRewardBody.Points >= 0 && GetRewardBody.Rewards[0].Qty >= 0 && GetRewardBody.Rewards[1].Qty >= 0 {
									ColorPrint.PrintCyan(userName, fmt.Sprintf("You made battle damage %s, battle points %s, reward Forgium %0.3f, reward Electrum %0.2f.", strconv.Itoa(GetRewardBody.TotalDmg), strconv.Itoa(GetRewardBody.Points), GetRewardBody.Rewards[0].Qty, GetRewardBody.Rewards[1].Qty))
									postJsonResult = true
									if shareBattleLog {
										go func() {
											var ShareLog = SpStruct.ShareCDPFitReturnData{}
											json.Unmarshal([]byte(GetResponseBody.Value.Body), &ShareLog)
											RequestFunc.ShareLogToApi(ShareLog, bossName, bossAbilities, bossRandomAbilities, bossLeague, heroTypechoosed, GetRewardBody.TotalDmg)
										}()
									}
								}
							}
						}
					}
					if returnJsonResult == true && postJsonResult == true {
						break
					} else if fetchTime >= 5 {
						break
					} else {
						time.Sleep(1 * time.Second)
						fetchTime++
						continue
					}
				}

			}
			checkTime := 0
			if postJsonResult == false && showForgeReward == true {
				err := DriverAction.DriverElementWaitAndClick(wd, "/html/body/app/div[1]/slcards/div[4]/div[2]/button[2]")
				if err == nil {
					for {
						el, _ = wd.FindElement(selenium.ByXPATH, "/html/body/app/div[1]/slcards/div[5]/div[1]/replay/section/rewards-modal/section/div[1]/div[1]/div[2]/span[1]/span")
						resultdmg, _ := el.Text()
						el, _ = wd.FindElement(selenium.ByXPATH, "/html/body/app/div[1]/slcards/div[5]/div[1]/replay/section/rewards-modal/section/div[1]/div[1]/div[2]/span[2]/span[2]")
						resultpoints, _ := el.Text()
						el, _ = wd.FindElement(selenium.ByXPATH, "/html/body/app/div[1]/slcards/div[5]/div[1]/replay/section/rewards-modal/section/div[1]/div[1]/p[1]")
						resultsstring, _ := el.Text()
						if resultpoints != "" && resultdmg != "" {
							parts := strings.Split(resultsstring, " ")
							forgiumStr := parts[3]
							electrumStr := parts[6]
							forgium, err := strconv.ParseFloat(forgiumStr, 64)
							if err != nil {
								ColorPrint.PrintRed(userName, "Encountering difficulty in reading the game results, but the battle has ended.")
								break
							}
							electrum, err := strconv.ParseFloat(electrumStr, 64)
							if err != nil {
								ColorPrint.PrintRed(userName, "Encountering difficulty in reading the game results, but the battle has ended.")
								break
							}
							ColorPrint.PrintCyan(userName, fmt.Sprintf("You made battle damage %s, battle points %s, reward Forgium %0.3f, reward Electrum %0.2f.", resultdmg, resultpoints, forgium, electrum))
							break
						} else if checkTime > 5 {
							ColorPrint.PrintRed(userName, "Encountering difficulty in reading the game results, but the battle has ended.")
							break
						} else {
							time.Sleep(2 * time.Second)
							checkTime++
							continue
						}
					}
				} else {
					ColorPrint.PrintRed(userName, "Encountering difficulty in reading the game results, but the battle has ended.")
				}
			}
			if showAccountDetails {
				LogFunc.PrintAccountDetails(userName, name, key, splinterforgeAPIEndpoint)
			}
			ColorPrint.PrintGold(userName, "Successful generated Cookies, the account will continue play with this setup.")
			wd.Close()
			s.Add(1)
			enoughMana := true
			go func() {
				for {
					if autoSelectSleepTime == true && enoughMana == true {
						ColorPrint.PrintWhite(userName, fmt.Sprintf("this account will enter a state of inactivity for %s minutes based on auto selected info.", mana))
						time.Sleep(time.Duration(manaused) * time.Minute)
					} else if autoSelectSleepTime == false && enoughMana == true {
						ColorPrint.PrintWhite(userName, fmt.Sprintf("According to your configuration, this account will enter a state of inactivity for %s minutes.", strconv.Itoa(timeSleepInMinute)))
						time.Sleep(time.Duration(timeSleepInMinute) * time.Minute)
					}
					BossLeague, BossAbilities, BossRandomAbilities := RequestFunc.FetchBossAbilities(name.(string), key.(string), bossName, splinterforgeAPIEndpoint)
					heroTypes := [3]string{"Warrior", "Wizard", "Ranger"}
					heroIndex, _ := strconv.Atoi(heroesType)
					var resHeroName string
					var FetchHeroErr error
					if autoSelectHero {
						resHeroName, FetchHeroErr = RequestFunc.FetchselectHero(BossRandomAbilities, name.(string), key.(string), publicAPIEndpoint, bossName, splinterforgeAPIEndpoint)
						if FetchHeroErr == nil {
							TeamIndex := len(fitPostData.Team) - 1
							TeamHero := fitPostData.Team[TeamIndex]
							stats, ok := TeamHero["stats"].(map[string]interface{})
							if ok {
								stats["attack"] = [4]int{0, 0, 0, 0}
								stats["ranged"] = [4]int{0, 0, 0, 0}
								stats["magic"] = [4]int{0, 0, 0, 0}
								TeamHero["stats"] = stats
								TeamHero["uid"] = resHeroName + " hero"
							}
						} else {
							resHeroName := heroTypes[heroIndex-1]
							TeamIndex := len(fitPostData.Team) - 1
							TeamHero := fitPostData.Team[TeamIndex]
							stats, ok := TeamHero["stats"].(map[string]interface{})
							if ok {
								stats["attack"] = [4]int{0, 0, 0, 0}
								stats["ranged"] = [4]int{0, 0, 0, 0}
								stats["magic"] = [4]int{0, 0, 0, 0}
								TeamHero["stats"] = stats
								TeamHero["uid"] = resHeroName + " hero"
							}
						}
					}
					ColorPrint.PrintWhite(userName, "Participating in battles with current setup...")
					reFit, err := grequests.Post(fmt.Sprintf("%s/boss/fight_boss", splinterforgeAPIEndpoint), &grequests.RequestOptions{
						JSON: fitPostData,
						Headers: map[string]string{
							"Content-Type": fitRes.Message.Params.Request.Headers.ContentType,
							"User-Agent":   fitRes.Message.Params.Request.Headers.UserAgent,
						},
					})
					if err == nil {
						LogFunc.PrintResultBox(userName, printData, selectResult, bossName, BossLeague, resHeroName, BossAbilities, BossRandomAbilities)
						if strings.Contains(reFit.String(), "not enough mana!") {
							ColorPrint.PrintYellow(userName, "Insufficient stamina, entering a rest state of inactivity for 1 hour...")
							time.Sleep(1 * time.Hour)
							enoughMana = false
							continue
						} else if strings.Contains(reFit.String(), "decoded message was invalid") {
							CookiesStatus = false
							ColorPrint.PrintRed(userName, "Cookies Error, Restarting...")
							break
						} else if strings.Contains(reFit.String(), "totalDmg") && strings.Contains(reFit.String(), "points") {
							ColorPrint.PrintGreen(userName, "Battle was successful!")
							var fitReturnData = SpStruct.FitReturnData{}
							json.Unmarshal(reFit.Bytes(), &fitReturnData)
							if showForgeReward {
								ColorPrint.PrintCyan(userName, fmt.Sprintf("You made battle damage %s, battle points %s, reward Forgium %0.3f, reward Electrum %0.2f.", strconv.Itoa(fitReturnData.TotalDmg), strconv.Itoa(fitReturnData.Points), fitReturnData.Rewards[0].Qty, fitReturnData.Rewards[1].Qty))
							}
							if shareBattleLog {
								go func() {
									var ShareLog = SpStruct.ShareCDPFitReturnData{}
									json.Unmarshal(reFit.Bytes(), &ShareLog)
									RequestFunc.ShareLogToApi(ShareLog, bossName, BossAbilities, BossRandomAbilities, BossLeague, resHeroName, fitReturnData.TotalDmg)
								}()
							}
							time.Sleep(5 * time.Second)
							if showAccountDetails {
								LogFunc.PrintAccountDetails(userName, name, key, splinterforgeAPIEndpoint)
							}
							enoughMana = true
							continue
						} else {
							Unexpected = true
							ColorPrint.PrintRed(userName, "Unexpected Error, Restarting...")
							break
						}
					} else {
						ColorPrint.PrintRed(userName, "Game server API error, Restarting...")
						Unexpected = true
						break
					}
				}
				if !CookiesStatus || Unexpected {
					AccountRestartCoroutine(accountLists, false, userName, headless, showForgeReward, showAccountDetails, autoSelectCard, autoSelectHero, autoSelectSleepTime, waitForBossRespawn, shareBattleLog, splinterforgeAPIEndpoint, splinterlandAPIEndpoint, publicAPIEndpoint, s, w)
					s.Done()
				}

			}()

			if wait == true {
				w.Done()
			}
		} else {
			if wait == true {
				w.Done()
			}
			if CurrentStamina < manaused {
				wd.Close()
				ColorPrint.PrintYellow(userName, "Insufficient stamina, entering a rest state of inactivity for 1 hour...")
				time.Sleep(1 * time.Hour)
				AccountRestartCoroutine(accountLists, false, userName, headless, showForgeReward, showAccountDetails, autoSelectCard, autoSelectHero, autoSelectSleepTime, waitForBossRespawn, shareBattleLog, splinterforgeAPIEndpoint, splinterlandAPIEndpoint, publicAPIEndpoint, s, w)
			} else if manaused < 30 {
				wd.Close()
				ColorPrint.PrintYellow(userName, "Card Selected not meet 30 mana requirements, restarting...")
				AccountRestartCoroutine(accountLists, false, userName, headless, showForgeReward, showAccountDetails, autoSelectCard, autoSelectHero, autoSelectSleepTime, waitForBossRespawn, shareBattleLog, splinterforgeAPIEndpoint, splinterlandAPIEndpoint, publicAPIEndpoint, s, w)
			}
		}
	} else {
		if wait == true {
			w.Done()
		}
		wd.Close()
		ColorPrint.PrintYellow(userName, "Boss is dead, trying in 30 minutes...")
		time.Sleep(30 * time.Minute)
		AccountRestartCoroutine(accountLists, false, userName, headless, showForgeReward, showAccountDetails, autoSelectCard, autoSelectHero, autoSelectSleepTime, waitForBossRespawn, shareBattleLog, splinterforgeAPIEndpoint, splinterlandAPIEndpoint, publicAPIEndpoint, s, w)
	}

}

func InitializeDriver(wait bool, userData SpStruct.UserData, headless bool, showForgeReward bool, showAccountDetails bool, autoSelectCard bool, autoSelectHero bool, autoSelectSleepTime bool, waitForBossRespawn bool, shareBattleLog bool, splinterforgeAPIEndpoint string, splinterlandAPIEndpoint string, publicAPIEndpoint string, accountLists []SpStruct.UserData, s *sync.WaitGroup, w *sync.WaitGroup) {
	ColorPrint.PrintWhite(userData.UserName, "Initializing...")
	extensionData, err := ioutil.ReadFile(ExtensionPath)
	if err != nil {
		println((1))
	}
	t := true
	extensionBase64 := base64.StdEncoding.EncodeToString(extensionData)
	chromeOptions := chrome.Capabilities{
		PerfLoggingPrefs: &chrome.PerfLoggingPreferences{
			EnableNetwork: &t,
			EnablePage:    &t},
		W3C:  false,
		Path: "",
		Args: []string{
			"--no-sandbox",
			"--disable-dev-shm-usage",
			"--disable-setuid-sandbox",
			"--disable-backgrounding-occluded-windows",
			"--disable-background-timer-throttling",
			"--disable-translate",
			"--disable-popup-blocking",
			"--disable-infobars",
			// "--disable-gpu",
			"--disable-blink-features=AutomationControlled",
			"--mute-audio",
			"--ignore-certificate-errors",
			"--allow-running-insecure-content",
			"--window-size=300,600",
		},

		Extensions: []string{extensionBase64},
		Prefs: map[string]interface{}{
			"profile.managed_default_content_settings.images":       1,
			"profile.managed_default_content_settings.cookies":      1,
			"profile.managed_default_content_settings.javascript":   1,
			"profile.managed_default_content_settings.plugins":      1,
			"profile.default_content_setting_values.notifications":  2,
			"profile.managed_default_content_settings.stylesheets":  2,
			"profile.managed_default_content_settings.popups":       2,
			"profile.managed_default_content_settings.geolocation":  2,
			"profile.managed_default_content_settings.media_stream": 2,
		},
		ExcludeSwitches: []string{
			"enable-automation",
			"enable-logging",
		},
	}
	if headless {
		chromeOptions.Args = append(chromeOptions.Args, "--headless=new")
	}
	caps := selenium.Capabilities{}
	caps.AddChrome(chromeOptions)
	caps.SetLogLevel(log.Performance, log.All)

	driver, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", 9515))
	if err != nil {
		fmt.Printf("Failed to create WebDriver: %s\n", err)
		os.Exit(1)
	}
	defer driver.Quit()
	loginResult := AccountLogin(userData.UserName, userData.PostingKey, driver)
	if loginResult {
		GameFunc.CheckPopUp(driver, 1000)
		AccountBattle(wait, driver, userData.UserName, userData.BossID, headless, userData.HeroesType, userData.TimeSleepInMinute, userData.CardSelection, autoSelectHero, autoSelectCard, autoSelectSleepTime, splinterlandAPIEndpoint, publicAPIEndpoint, splinterforgeAPIEndpoint, showForgeReward, showAccountDetails, waitForBossRespawn, shareBattleLog, s, w, accountLists)
	} else {
		driver.Close()
		if wait == true {
			w.Done()
		}
		ColorPrint.PrintYellow(userData.UserName, "Retrying in 30 seconds...")
		time.Sleep(30 * time.Second)
		AccountRestartCoroutine(accountLists, false, userData.UserName, headless, showForgeReward, showAccountDetails, autoSelectCard, autoSelectHero, autoSelectSleepTime, waitForBossRespawn, shareBattleLog, splinterforgeAPIEndpoint, splinterlandAPIEndpoint, publicAPIEndpoint, s, w)
	}
}

func AccountRestartCoroutine(accountLists []SpStruct.UserData, wait bool, userName string, headless bool, showForgeReward bool, showAccountDetails bool, autoSelectCard bool, autoSelectHero bool, autoSelectSleepTime bool, waitForBossRespawn bool, shareBattleLog bool, splinterforgeAPIEndpoint string, splinterlandAPIEndpoint string, publicAPIEndpoint string, s *sync.WaitGroup, w *sync.WaitGroup) {
	for _, v := range accountLists {
		if v.UserName == userName {
			InitializeDriver(wait, v, headless, showForgeReward, showAccountDetails, autoSelectCard, autoSelectHero, autoSelectSleepTime, waitForBossRespawn, shareBattleLog, splinterforgeAPIEndpoint, splinterlandAPIEndpoint, publicAPIEndpoint, accountLists, s, w)
		}
	}
}
