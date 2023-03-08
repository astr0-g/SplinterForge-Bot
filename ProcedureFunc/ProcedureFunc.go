package ProcedureFunc

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
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
	err = el.SendKeys("\ue007")
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
		el, err = wd.FindElement(selenium.ByXPATH, "/html/body/app/div[1]/login-modal/div/div/div/div[2]/div[3]/button")
		if err != nil {
			ColorPrint.PrintRed(userName, "Login Error!")
			return false
		}
		el.Click()
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

func AccountBattle(wait bool, wd selenium.WebDriver, userName string, bossId string, headless bool, heroesType string, timeSleepInMinute int, cardSelection []SpStruct.CardSelection, autoSelectHero bool, autoSelectCard bool, autoSelectSleepTime bool, splinterlandAPIEndpoint string, publicAPIEndpoint string, splinterforgeAPIEndpoint string, showForgeReward bool, showAccountDetails bool, s *sync.WaitGroup, w *sync.WaitGroup, accountLists []SpStruct.UserData) {
	starttimestamp := time.Now().Unix()
	CookiesStatus := true
	Unexpected := false
	name, _ := wd.ExecuteScript("return localStorage.getItem('forge:username');", nil)
	key, _ := wd.ExecuteScript("return localStorage.getItem('forge:key');", nil)
	bossName, _, _ := GameFunc.SelectBoss(userName, bossId, wd)
	cardSelection, _, _ = GameFunc.SelectCards(cardSelection, bossName, userName, splinterlandAPIEndpoint, publicAPIEndpoint, autoSelectCard)
	seletedNumOfSummoners := 1
	el, _ := wd.FindElement(selenium.ByXPATH, "/html/body/app/div[1]/slcards/div[5]/section[1]/div/div[1]/div[2]/button")
	el.Click()
	ColorPrint.PrintWhite(userName, "Participating in battles...")
	printData := [][]string{}
	selectResult := true
	GameFunc.SelectHero(heroesType, userName, wd, autoSelectHero, publicAPIEndpoint, bossName, splinterforgeAPIEndpoint)
	for _, selection := range cardSelection {
		for i, playingSummoner := range selection.PlayingSummoners {
			result := GameFunc.SelectSummoners(userName, seletedNumOfSummoners, playingSummoner.PlayingSummonersDiv, wd)
			if result {
				seletedNumOfSummoners++
				printData = append(printData, []string{fmt.Sprintf("Summoners #%d", i+1), playingSummoner.PlayingSummonersID, playingSummoner.PlayingSummonersName, "success"})
			}
		}
		seletedNumOfMonsters := 1
		for j, playingMonster := range selection.PlayingMonsters {
			result := GameFunc.SelectMonsters(userName, seletedNumOfMonsters, playingMonster.PlayingMontersDiv, wd)
			if result {
				seletedNumOfMonsters++
				printData = append(printData, []string{fmt.Sprintf("Monsters #%d", j+1), playingMonster.PlayingMonstersID, playingMonster.PlayingMonstersName, "success"})
			} else {
				printData = append(printData, []string{fmt.Sprintf("Monsters #%d", j+1), playingMonster.PlayingMonstersID, playingMonster.PlayingMonstersName, "error"})
				selectResult = false
			}
		}
	}
	ReadFunc.GetAccountDetails(name, key, splinterforgeAPIEndpoint)
	el, _ = wd.FindElement(selenium.ByXPATH, "/html/body/app/div[1]/slcards/div[5]/div[2]/div[1]/div[1]/button/span")
	mana, _ := el.Text()
	manaused, _ := strconv.Atoi(mana)
	CurrentStamina := ReadFunc.GetAccountDetails(name, key, splinterforgeAPIEndpoint)
	if CurrentStamina > manaused && manaused >= 15 {
		LogFunc.PrintResultBox(userName, printData, selectResult)
		battletimestamp := time.Now().Unix()
		if battletimestamp - starttimestamp < 30 {
			time.Sleep(time.Duration(starttimestamp + 30 - battletimestamp) *time.Second)
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
				continue
			}
			netLogs, _ := wd.Log("performance")
			for _, netLog := range netLogs {
				if returnJsonResult == false && strings.Contains(netLog.Message, fmt.Sprintf("%s/boss/fight_boss", splinterforgeAPIEndpoint)) && strings.Contains(netLog.Message, "\"method\":\"Network.requestWillBeSent\"") {
					json.Unmarshal([]byte(netLog.Message), &fitRes)
					json.Unmarshal([]byte(fitRes.Message.Params.Request.PostData), &fitPostData)
					ColorPrint.PrintGreen(userName, "Battle was successful!")
					returnJsonResult = true
				}
				if showForgeReward == true && postJsonResult == false && strings.Contains(netLog.Message, fmt.Sprintf("%s/boss/fight_boss", splinterforgeAPIEndpoint)) && strings.Contains(netLog.Message, "\"method\":\"Network.responseReceived\"") {
					var GetResponseBody = SpStruct.GetResponseBody{}
					var GetRewardBody = SpStruct.CDPFitReturnData{}
					resString, err := RequestFunc.GetReponseBody(wd.SessionID(), fitRes.Message.Params.RequestID, userName)
					if err == nil && resString != "" {
						json.Unmarshal([]byte(resString), &GetResponseBody)
						json.Unmarshal([]byte(GetResponseBody.Value.Body), &GetRewardBody)
						if GetRewardBody.TotalDmg >= 0 && GetRewardBody.Points >= 0  && GetRewardBody.Rewards[0].Qty >= 0 && GetRewardBody.Rewards[1].Qty >= 0 {
							ColorPrint.PrintCyan(userName, fmt.Sprintf("You made battle damage %s, battle points %s, reward Forgium %0.3f, reward Electrum %0.2f.", strconv.Itoa(GetRewardBody.TotalDmg), strconv.Itoa(GetRewardBody.Points), GetRewardBody.Rewards[0].Qty, GetRewardBody.Rewards[1].Qty))
							postJsonResult = true
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
						forgium, _ := strconv.ParseFloat(forgiumStr, 64)
						electrum, _ := strconv.ParseFloat(electrumStr, 64)
						ColorPrint.PrintCyan(userName, fmt.Sprintf("You made battle damage %s, battle points %s, reward Forgium %0.3f, reward Electrum %0.2f.", resultdmg, resultpoints, forgium, electrum))
						break
					} else if checkTime > 5 {
						ColorPrint.PrintRed(userName, "Encountering difficulty in reading the game results, but the battle has ended.")
						break
					} else{
						time.Sleep(5 * time.Second)
						checkTime ++
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
		go func() {
			for {
				if autoSelectSleepTime {
					ColorPrint.PrintWhite(userName, fmt.Sprintf("this account will enter a state of inactivity for %s minutes based on auto selected info.", mana))
					time.Sleep(time.Duration(manaused) * time.Minute)
				} else {
					ColorPrint.PrintWhite(userName, fmt.Sprintf("According to your configuration, this account will enter a state of inactivity for %s minutes.", strconv.Itoa(timeSleepInMinute)))
					time.Sleep(time.Duration(timeSleepInMinute) * time.Minute)
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
					LogFunc.PrintResultBox(userName, printData, selectResult)
					if strings.Contains(reFit.String(), "not enough mana!") {
						ColorPrint.PrintYellow(userName, "Insufficient stamina, entering a rest state of inactivity for 1 hour...")
						time.Sleep(1 * time.Hour)
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
						time.Sleep(5 * time.Second)
						if showAccountDetails {
							LogFunc.PrintAccountDetails(userName, name, key, splinterforgeAPIEndpoint)
						}
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
				AccountRestartCoroutine(accountLists, false, userName, headless, showForgeReward, showAccountDetails, autoSelectCard, autoSelectHero, autoSelectSleepTime, splinterforgeAPIEndpoint, splinterlandAPIEndpoint, publicAPIEndpoint, s, w)
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
			AccountRestartCoroutine(accountLists, false, userName, headless, showForgeReward, showAccountDetails, autoSelectCard, autoSelectHero, autoSelectSleepTime, splinterforgeAPIEndpoint, splinterlandAPIEndpoint, publicAPIEndpoint, s, w)
		} else if manaused < 15 {
			wd.Close()
			ColorPrint.PrintYellow(userName, "Card Selected not meet 15 mana requirements, restarting...")
			AccountRestartCoroutine(accountLists, false, userName, headless, showForgeReward, showAccountDetails, autoSelectCard, autoSelectHero, autoSelectSleepTime, splinterforgeAPIEndpoint, splinterlandAPIEndpoint, publicAPIEndpoint, s, w)
		}
	}
}

func InitializeDriver(wait bool, userData SpStruct.UserData, headless bool, showForgeReward bool, showAccountDetails bool, autoSelectCard bool, autoSelectHero bool, autoSelectSleepTime bool, splinterforgeAPIEndpoint string, splinterlandAPIEndpoint string, publicAPIEndpoint string, accountLists []SpStruct.UserData, s *sync.WaitGroup, w *sync.WaitGroup) {
	ColorPrint.PrintWhite(userData.UserName, "Initializing...")
	extensionData, err := ioutil.ReadFile("data/hivekeychain.crx")
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
		AccountBattle(wait, driver, userData.UserName, userData.BossID, headless, userData.HeroesType, userData.TimeSleepInMinute, userData.CardSelection, autoSelectHero, autoSelectCard, autoSelectSleepTime, splinterlandAPIEndpoint, publicAPIEndpoint, splinterforgeAPIEndpoint, showForgeReward, showAccountDetails, s, w, accountLists)
	} else {
		driver.Close()
		if wait == true {
			w.Done()
		}
		ColorPrint.PrintYellow(userData.UserName, "Retrying in 30 seconds...")
		time.Sleep(30 * time.Second)
		AccountRestartCoroutine(accountLists, false, userData.UserName, headless, showForgeReward, showAccountDetails, autoSelectCard, autoSelectHero, autoSelectSleepTime, splinterforgeAPIEndpoint, splinterlandAPIEndpoint, publicAPIEndpoint, s, w)
	}
}

func AccountRestartCoroutine(accountLists []SpStruct.UserData, wait bool, userName string, headless bool, showForgeReward bool, showAccountDetails bool, autoSelectCard bool, autoSelectHero bool, autoSelectSleepTime bool, splinterforgeAPIEndpoint string, splinterlandAPIEndpoint string, publicAPIEndpoint string, s *sync.WaitGroup, w *sync.WaitGroup) {
	for _, v := range accountLists {
		if v.UserName == userName {
			InitializeDriver(wait, v, headless, showForgeReward, showAccountDetails, autoSelectCard, autoSelectHero, autoSelectSleepTime, splinterforgeAPIEndpoint, splinterlandAPIEndpoint, publicAPIEndpoint, accountLists, s, w)
		}
	}
}
