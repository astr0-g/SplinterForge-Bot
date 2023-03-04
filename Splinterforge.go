package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/golang/glog"
	"github.com/levigross/grequests"
	"github.com/olekukonko/tablewriter"
	"github.com/selenium-Driver-Check/SeleniumDriverCheck"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"github.com/tebeka/selenium/log"

	"splinterforge/spstruct"
)

func printResultBox(userName string, data [][]string, selectResult bool) {
	sort.Slice(data, func(i, j int) bool {
		return data[i][0] < data[j][0]
	})

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Card", "ID", "Name", "Results"})
	for _, row := range data {
		table.Append(row)
	}

	if selectResult {
		PrintGreen(userName, "Card selection results:")
		table.Render()
		color.Set(color.FgWhite)
	} else {
		PrintYellow(userName, "Card selection results:")
		table.Render()
		color.Set(color.FgWhite)
	}
}
func printConfigSettings(totalAccounts int, headless bool, closeDriverWhileSleeping bool, startThread int, startThreadInterval int, showForgeReward bool, showTotalForgeBalance bool, printSystemUsage bool, checkSystemUsageFrequency int, autoSelectCard bool, autoSelectHero bool, autoSelectSleepTime int) {
	data := [][]string{
		{"TOTAL_ACCOUNTS_LOADED", fmt.Sprint(totalAccounts)},
		{"HEADLESS", fmt.Sprint(headless)},
		{"CLOSE_DRIVER_WHILE_SLEEPING", fmt.Sprint(closeDriverWhileSleeping)},
		{"START_THREAD", fmt.Sprint(startThread)},
		{"START_THREAD_INTERVAL(seconds)", fmt.Sprint(startThreadInterval)},
		{"SHOW_FORGE_REWARD", fmt.Sprint(showForgeReward)},
		{"SHOW_TOTAL_FORGE_BALANCE", fmt.Sprint(showTotalForgeBalance)},
		{"PRINT_SYSTEM_USAGE", fmt.Sprint(printSystemUsage)},
		{"CHECK_SYSTEM_USAGE_FREQUENCY(seconds)", fmt.Sprint(checkSystemUsageFrequency)},
		{"AUTO_SELECT_CARD", fmt.Sprint(autoSelectCard)},
		{"AUTO_SELECT_HERO", fmt.Sprint(autoSelectHero)},
		{"AUTO_SELECT_SLEEPTIME", fmt.Sprint(autoSelectSleepTime)},
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i][0] < data[j][0]
	})

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Settings", "Value"})
	for _, row := range data {
		table.Append(row)
	}
	table.Render()
}

func PrintYellow(username string, message string) {
	now := time.Now()
	color.Set(color.FgYellow)
	fmt.Println("["+now.Format("2006-01-02 15:04:05")+"]", username+":", message)
}
func PrintRed(username string, message string) {
	now := time.Now()
	color.Set(color.FgRed)
	fmt.Println("["+now.Format("2006-01-02 15:04:05")+"]", username+":", message)
}
func PrintGreen(username string, message string) {
	now := time.Now()
	color.Set(color.FgGreen)
	fmt.Println("["+now.Format("2006-01-02 15:04:05")+"]", username+":", message)
}
func PrintBlue(username string, message string) {
	now := time.Now()
	color.Set(color.FgBlue)
	fmt.Println("["+now.Format("2006-01-02 15:04:05")+"]", username+":", message)
}
func PrintWhite(username string, message string) {
	now := time.Now()
	color.Set(color.FgWhite)
	fmt.Println("["+now.Format("2006-01-02 15:04:05")+"]", username+":", message)
}
func getCardName(cardId string) (string, error) {
	file, err := os.Open("data/cardMapping.json")
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

func getAccountData(filePath string, lineNumber int) (string, string, error) {
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

func getCardSettingData(filePath string, lineNumber int) (string, string, []string, []string, int, error) {
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
func elementWaitAndClick(wd selenium.WebDriver, xpath string) {
	byXpath := selenium.ByXPATH
	for {
		element, err := wd.FindElement(byXpath, xpath)
		if err != nil {
			panic(err)
		}
		isEnabled, err := element.IsEnabled()
		if err != nil {
			panic(err)
		}
		if isEnabled {
			err = element.Click()
			if err != nil {
				panic(err)
			}
			break
		}
		time.Sleep(1 * time.Second)
	}
}
func fetchHeroSelect(publicAPIEndpoint string, bossName string) (string, error) {
	bossID := fetchBossID(bossName)

	// Set up the request body
	requestBody := map[string]interface{}{
		"bossId": bossID,
	}

	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	// Set up the HTTP request
	endpoint := fmt.Sprintf("%s/heroselection", publicAPIEndpoint)
	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(requestBodyBytes))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the request and parse the response
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var responseData map[string]interface{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&responseData)
	if err != nil {
		return "", err
	}

	heroTypeToChoose := strings.Split(responseData["heroTypes"].(string), " ")[0]
	return heroTypeToChoose, nil
}

func fetchBossID(bossName string) string {
	url := "https://splinterforge.io/boss/getBosses"

	// Send the HTTP GET request and parse the JSON response
	resp, err := http.Get(url)
	if err != nil {
		glog.Fatal(err)
	}
	defer resp.Body.Close()

	var responseData []map[string]interface{}
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&responseData)

	// Search for the boss ID by name
	for _, bossData := range responseData {
		if strings.EqualFold(bossData["name"].(string), bossName) {
			return bossData["id"].(string)
		}
	}

	return ""
}
func checkPopUp(wd selenium.WebDriver, millisecond int) {
	defer func() {
		if err := recover(); err != nil {
			// Handle any panic that occurs during the execution of the function
		}
	}()
	if element, err := wd.FindElement(selenium.ByXPATH, "/html/body/app/div[1]/div[1]/app-header/success-modal/section/div[1]/div[4]/div/button"); err == nil {
		if err = element.Click(); err != nil {
			// Handle any errors that occur during the click operation
		}
	}
	if element, err := wd.FindElement(selenium.ByXPATH, "/html/body/app/div[1]/login-modal/div/div/div/div[2]/div[3]/button"); err == nil {
		if err = element.Click(); err != nil {
			// Handle any errors that occur during the click operation
		}
	}
	duration := time.Duration(millisecond) * time.Millisecond
	time.Sleep(duration)
}
func DriverGet(URL string, wd selenium.WebDriver) {
	err := wd.Get(URL)
	if err != nil {
		panic(err)
	}
	script := `
	    var imgs = document.getElementsByTagName('img');
	    for (var i = 0; i < imgs.length; i++) {
	        imgs[i].parentNode.removeChild(imgs[i]);
	    }
	    var style = document.createElement('style');
	    style.innerHTML = 'img { opacity: 0 }';
	    document.head.appendChild(style);
	    var style = document.createElement('style');
	    style.innerHTML = '* { background-image: none !important; }';
	    document.head.appendChild(style);
	    var style = document.createElement('style');
	    style.innerHTML = '* { color: transparent !important; }';
	    document.head.appendChild(style);
	    var style = document.createElement('style');
	    style.innerHTML = 'img.fade_image { display: none !important; }';
	    document.head.appendChild(style);
	    var style = document.createElement('style');
	    style.innerHTML = '* { transition: paused !important; }';
	    document.head.appendChild(style);
	`
	wd.ExecuteScript(script, nil)
}
func checklogin(userName string, wd selenium.WebDriver) {
	for {
		el, err := wd.FindElement(selenium.ByXPATH, "/html/body/app/div[1]/div[1]/app-header/section/div[4]/div[2]/div[2]/div[2]/a/div[2]")
		if err != nil {
			time.Sleep(3 * time.Second)
			continue
		} else {
			text, _ := el.Text()
			if text == userName {
				break
			} else {
				fmt.Println("username erro")
			}
		}
	}
}
func login(userName string, postingKey string, wd selenium.WebDriver) {

	err := wd.SetImplicitWaitTimeout(5 * time.Second)
	if err != nil {
		panic(err)
	}
	DriverGet("chrome-extension://jcacnejopjdphbnjgfaaobbfafkihpep/popup.html", wd)

	elementWaitAndClick(wd, "/html/body/div/div/div[4]/div[2]/div[5]/button")

	el, _ := wd.FindElement(selenium.ByXPATH, "/html/body/div/div/div[1]/div/div[1]/div/input")
	el.SendKeys("Aa123Aa123!!")
	el, _ = wd.FindElement(selenium.ByXPATH, "/html/body/div/div/div[1]/div/div[2]/div/input")
	el.SendKeys("Aa123Aa123!!")
	el, _ = wd.FindElement(selenium.ByXPATH, "/html/body/div/div/div[1]/button/div")
	el.Click()
	el, _ = wd.FindElement(selenium.ByXPATH, "/html/body/div/div/div[1]/div[2]/div/div[2]/button[1]/div")
	el.Click()
	el, _ = wd.FindElement(selenium.ByXPATH, "/html/body/div/div/div[1]/div[2]/div/div[2]/div[1]/div/input")
	el.SendKeys(userName)
	el, _ = wd.FindElement(selenium.ByXPATH, "/html/body/div/div/div[1]/div[2]/div/div[2]/div[2]/div/input")
	el.SendKeys(postingKey)
	el, _ = wd.FindElement(selenium.ByXPATH, "/html/body/div/div/div[1]/div[2]/div/div[2]/div[2]/div/input")
	time.Sleep(2 * time.Second)
	el.SendKeys("\ue007")
	err = wd.ResizeWindow("bigger", 1565, 1080)
	if err != nil {
		println("can not change size")
	}

	// wd.SetWindowSize(1565, 1080)
	DriverGet("https://splinterforge.io/#/", wd)

	el, _ = wd.FindElement(selenium.ByXPATH, "/html/body/app/div[1]/div[1]/app-header/success-modal/section/div[1]/div[4]/div/button")
	el.Click()
	el, _ = wd.FindElement(selenium.ByXPATH, "/html/body/app/div[1]/div[1]/app-header/section/div[4]/div[2]/div/div/a/div[1]")
	el.Click()
	el, _ = wd.FindElement(selenium.ByXPATH, "/html/body/app/div[1]/login-modal/div/div/div/div[2]/div[2]/input")
	el.SendKeys(userName)
	el, _ = wd.FindElement(selenium.ByXPATH, "/html/body/app/div[1]/login-modal/div/div/div/div[2]/div[3]/button")
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
	checklogin(userName, wd)
	PrintGreen(userName, "success log in")
}
func initializeAccount(accountNo int) (string, string, string, string, []spstruct.CardSelection, int) {

	userName, postingKey, err := getAccountData("config/accounts.txt", accountNo)
	if err != nil || userName == "" || postingKey == "" {
		fmt.Println("Error in loading accounts.txt, please add username or posting key and try again.")
	}
	heroesType, bossId, playingSummoners, playingMonster, timeSleepInMinute, err := getCardSettingData("config/cardSettings.txt", accountNo)
	if err != nil {
		fmt.Println("Error loading cardSettings.txt file")
	}
	playingSummonersList := make([]spstruct.Summoners, 0, len(playingSummoners))
	playingMonsterList := make([]spstruct.MonsterId, 0, len(playingMonster))
	for _, i := range playingSummoners {
		cardName, _ := getCardName(i)
		playingSummonersList = append(playingSummonersList, spstruct.Summoners{
			PlayingSummonersDiv:  fmt.Sprintf("//div/img[@id='%s']", i),
			PlayingSummonersID:   i,
			PlayingSummonersName: cardName,
		})
	}
	for _, i := range playingMonster {
		cardName, _ := getCardName(i)
		playingMonsterList = append(playingMonsterList, spstruct.MonsterId{
			PlayingMontersDiv:   fmt.Sprintf("//div/img[@id='%s']", i),
			PlayingMonstersID:   i,
			PlayingMonstersName: cardName,
		})
	}
	var cardSelectionList = []spstruct.CardSelection{}
	cardSelection := spstruct.CardSelection{
		PlayingSummoners: playingSummonersList,
		PlayingMonsters:  playingMonsterList,
	}
	cardSelectionList = append(cardSelectionList, cardSelection)

	timeSleepInMinute *= 60

	return userName, postingKey, heroesType, bossId, cardSelectionList, timeSleepInMinute
}
func fetchPlayerCard(userName string, splinterland_api_endpoint string) ([]int, error) {
	url := fmt.Sprintf("%s/cards/collection/%s", splinterland_api_endpoint, userName)
	method := "GET"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	err := writer.Close()
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		panic(err)
	}

	req.Header.Add("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var jsonData interface{}
	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		panic(err)
	}

	jsonString, err := json.Marshal(jsonData)
	if err != nil {
		panic(err)
	}
	var cardList spstruct.CardList
	err = json.Unmarshal([]byte(jsonString), &cardList)
	if err != nil {
		panic(err)
	}

	seen := map[int]bool{}
	cardDetailIDs := []int{}
	for _, card := range cardList.Cards {
		if !seen[card.CardDetailID] {
			cardDetailIDs = append(cardDetailIDs, card.CardDetailID)
			seen[card.CardDetailID] = true
		}
	}
	// fmt.Println(cardDetailIDs)
	return cardDetailIDs, nil
}
func fetchBattleCards(bossName string, userName string, splinterland_api_endpoint string, public_api_endpoint string) (string, error) {
	// Fetch player cards
	cardsId, err := fetchPlayerCard(userName, splinterland_api_endpoint)
	if err != nil {
		return "", err
	}

	// Fetch boss ID
	bossId := fetchBossID(bossName)

	// Prepare the request body
	postData := spstruct.BattleCardsRequestBody{
		BossId:   bossId,
		BossName: bossName,
		Team:     cardsId,
	}

	// Make the HTTP request
	url := fmt.Sprintf("%s/teamselection/", public_api_endpoint)
	requestOptions := &grequests.RequestOptions{
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		JSON: postData,
	}
	gresponse, err := grequests.Post(url, requestOptions)
	if err != nil {
		return "", err
	}

	// Check if the response was successful
	if !gresponse.Ok {
		return "", fmt.Errorf("HTTP error: %d - %s", gresponse.StatusCode, gresponse.String())
	}

	// Parse the response JSON
	var responseData map[string]interface{}
	err = json.Unmarshal(gresponse.Bytes(), &responseData)
	if err != nil {
		return "", err
	}

	// Encode the response as JSON
	jsonResponse, err := json.Marshal(responseData)
	if err != nil {
		return "", err
	}

	// Return the response as a JSON-encoded string
	return string(jsonResponse), nil
}
func autoSelectCards(cardSelection []spstruct.CardSelection, bossName string, userName string, splinterland_api_endpoint string, public_api_endpoint string, autoSelectCard bool) ([]spstruct.CardSelection, bool, error) {
	if autoSelectCard {
		PrintWhite(
			userName, fmt.Sprintf("Auto selecting playing cards for desire boss: %s", bossName))
		playingDeck, err := fetchBattleCards(bossName, userName, splinterland_api_endpoint, public_api_endpoint)
		if err != nil {
			// fmt.Println(err)
			return cardSelection, false, err
		}
		playingDeckBytes := []byte(playingDeck)
		var playingDeckMap map[string]interface{}
		err = json.Unmarshal(playingDeckBytes, &playingDeckMap)
		if err != nil {
			return cardSelection, false, err
		}
		// fmt.Println(playingDeckMap["RecommendTeam"].(bool))
		// fmt.Println(playingDeckMap)
		if playingDeckMap["RecommendTeam"].(bool) {
			summonerIds, ok := playingDeckMap["summonerIds"].([]interface{})
			if !ok {
				return cardSelection, false, fmt.Errorf("summonerIds is not an array")
			}

			monsterIds, ok := playingDeckMap["monsterIds"].([]interface{})
			if !ok {
				return cardSelection, false, fmt.Errorf("summonerIds is not an array")
			}

			playingSummonersList := make([]spstruct.Summoners, 0, len(summonerIds))
			playingMonsterList := make([]spstruct.MonsterId, 0, len(monsterIds))
			for _, i := range summonerIds {
				summonerId := fmt.Sprintf("%v", i)
				cardName, _ := getCardName(summonerId)
				playingSummonersList = append(playingSummonersList, spstruct.Summoners{
					PlayingSummonersID:   summonerId,
					PlayingSummonersName: cardName,
					PlayingSummonersDiv:  fmt.Sprintf("//div/img[@id='%s']", summonerId),
				})
			}

			for _, i := range monsterIds {
				monsterId := fmt.Sprintf("%v", i)
				cardName, _ := getCardName(monsterId)
				playingMonsterList = append(playingMonsterList, spstruct.MonsterId{
					PlayingMonstersID:   monsterId,
					PlayingMonstersName: cardName,
					PlayingMontersDiv:   fmt.Sprintf("//div/img[@id='%s']", monsterId),
				})
			}

			var cardSelectionList = []spstruct.CardSelection{}
			cardSelection := spstruct.CardSelection{
				PlayingSummoners: playingSummonersList,
				PlayingMonsters:  playingMonsterList,
			}
			cardSelectionList = append(cardSelectionList, cardSelection)
			// fmt.Println(cardSelectionList)
			return cardSelectionList, true, nil
		} else {
			// log_info.status(
			//  userName, "Auto selecting playing cards deck failed, you might have too less cards in the account, will continue play with your card setting.")
			return cardSelection, false, nil
		}
	} else {
		return cardSelection, false, nil
	}
}

// return cardSelection, autoSelectResult
func heroSelect(heroesType string, userName string, wd selenium.WebDriver, auto_select_hero bool, public_api_endpoint string, bossName string) {
	checkPopUp(wd, 1000)
	heroTypes := [3]string{"Warrior", "Wizard", "Ranger"}
	if auto_select_hero {
		hero_type, err := fetchHeroSelect(public_api_endpoint, bossName)
		if err == nil {
			PrintYellow(userName, fmt.Sprintf("Auto selecting heroes type: %s for desired boss: %s", hero_type, bossName))
			for i, val := range heroTypes {
				if val == hero_type {
					heroesType = strconv.Itoa(i + 1)
					break
				}
			}
		} else {
			PrintRed(userName, "Auto selecting heroes type failed due to API error.")
		}
	} else {
		PrintYellow(userName, "Selecting heroes type...")
	}
	defer func() {
		if err := recover(); err != nil {
			PrintRed(userName, "Error in selecting hero type, continue...")
		}
	}()
	heroIndex, _ := strconv.Atoi(heroesType)
	hero_type := heroTypes[heroIndex-1]
	el, _ := wd.FindElement(selenium.ByXPATH, "/html/body/app/div[1]/div[1]/app-header/section/div[4]/div[2]/div[1]/a[4]")
	el.Click()
	bossXpath := "/html/body/app/div[1]/splinterforge-heros/div[3]/section/div/div/div[2]/div[1]"
	el, _ = wd.FindElement(selenium.ByXPATH, bossXpath)
	el.Click()
	bossSelectXpath := fmt.Sprintf("%s/ul/li[%s]", bossXpath, heroesType)
	el, _ = wd.FindElement(selenium.ByXPATH, bossSelectXpath)
	el.Click()
	PrintGreen(userName, fmt.Sprintf("Selected hero type: %s", hero_type))
}
func bossSelect(userName string, bossIdToSelect string, wd selenium.WebDriver) (string, string, error) {
	wd.SetImplicitWaitTimeout(2 * time.Second)
	// Click on the "Bosses" button

	// Loop until the boss is defeated or a timeout occurs
	for {
		if element, err := wd.FindElement(selenium.ByXPATH, "/html/body/app/div[1]/div[1]/app-header/section/div[4]/div[2]/div[1]/a[5]/div[1]"); err == nil {
			if err = element.Click(); err != nil {
				// fmt.Println(err)
				// Handle any errors that occur during the click operation
			}
		}
		time.Sleep(1 * time.Second)
		// Click on the boss to select it
		bossSelector := fmt.Sprintf("//div[@tabindex='%s']", bossIdToSelect)
		if element, err := wd.FindElement(selenium.ByXPATH, bossSelector); err == nil {
			if err = element.Click(); err != nil {
				// fmt.Println(err)
			}
		}
		// Check if the boss is defeated
		element, err := wd.FindElement(selenium.ByXPATH, "/html/body/app/div[1]/slcards/div[5]/section[1]/div/div[1]/div[2]/button")
		if err == nil {
			text, err := element.Text()
			if err == nil {
				if text != "BOSS IS DEAD" {
					time.Sleep(1 * time.Second)

					// Get the boss name
					if element, err := wd.FindElement(selenium.ByXPATH, "/html/body/app/div[1]/slcards/div[5]/section[1]/div/div[1]/div[1]/div[2]/div[3]/h3"); err == nil {
						bossName, _ := element.Text()
						return bossName, bossIdToSelect, nil
					}
				} else {
					PrintRed(userName, "The selected boss has been defeated, selecting another one automatically...")
					// Select the next boss
					if bossIdToSelect < "17" {
						bossIdInt, err := strconv.Atoi(bossIdToSelect)
						if err != nil {
							// handle the error
						}
						bossIdInt++
						bossIdToSelect = strconv.Itoa(bossIdInt)
					} else {
						bossIdToSelect = "14"
					}
				}
			}
		}

	}
}
func waitForElement(wd selenium.WebDriver, xpath string) (bool, error) {
	for i := 0; i < 5; i++ {
		_, err := wd.FindElement(selenium.ByXPATH, xpath)
		if err != nil {
			time.Sleep(100 * time.Millisecond)
			continue
		} else {
			return true, nil
		}

	}
	return false, nil
}
func selectSummoners(userName string, seletedNumOfSummoners int, cardDiv string, wd selenium.WebDriver) bool {
	scroolTime := 0
	clickedTime := 0
	result := false
	time.Sleep(1 * time.Second)
	for scroolTime < 5 && clickedTime < 5 {
		el, err := wd.FindElement(selenium.ByXPATH, cardDiv)
		if err != nil {
			// fmt.Println(err)
			continue
		} else {
			el.Click()
			time.Sleep(1 * time.Second)
			checkCardDiv := fmt.Sprintf("/html/body/app/div[1]/slcards/div[5]/section[1]/div/div[1]/div[2]/div[1]/div[2]/div[%s]", strconv.Itoa(seletedNumOfSummoners))
			success, err := waitForElement(wd, checkCardDiv)
			if err != nil {
				// panic(err)
			}

			if success {
				// fmt.Println("Button clicked!")
				result = true
				break
			} else {
				// fmt.Println("Button not clicked!")
				clickedTime++
				scroolTime++
			}
		}

	}
	if clickedTime%2 == 1 {
		result = true
	}
	if !result {
		PrintRed(userName, "Error in selecting summoners, skipped...")
	}
	wd.ExecuteScript("window.scrollBy(0, -4000)", nil)
	time.Sleep(1 * time.Second)
	return result
}
func selectMonsters(userName string, seletedNumOfMonsters int, cardDiv string, wd selenium.WebDriver) bool {
	scroolTime := 0
	clickedTime := 0
	result := false
	time.Sleep(1 * time.Second)
	for scroolTime < 5 && clickedTime < 5 {
		el, err := wd.FindElement(selenium.ByXPATH, cardDiv)
		if err != nil {
			// fmt.Println(err)
			wd.ExecuteScript("window.scrollBy(0, 450)", nil)
			scroolTime++
			continue
		} else {
			el.Click()
			checkCardDiv := fmt.Sprintf("/html/body/app/div[1]/slcards/div[5]/section[1]/div/div[1]/div[2]/div[2]/div[2]/div[%s]", strconv.Itoa(seletedNumOfMonsters))
			// fmt.Println(checkCardDiv)
			success, _ := waitForElement(wd, checkCardDiv)
			if success {
				// fmt.Println("Button clicked!")
				// fmt.Println("good")
				result = true
				break
			} else {
				// fmt.Println("Button not clicked!")
				clickedTime++
				wd.ExecuteScript("window.scrollBy(0, 450)", nil)
				continue
			}
		}

	}
	if clickedTime%2 == 1 {
		result = true
	}
	if !result {
		PrintRed(userName, "Error in selecting Monsters, skipped...")
	}
	wd.ExecuteScript("window.scrollBy(0, -4000)", nil)
	time.Sleep(1 * time.Second)
	return result

}
func getConfig(filePath string) (bool, bool, int, int, bool, bool, bool, int, bool, bool, int, string, string, string) {
	file, err := os.Open(filePath)
	if err != nil {
		PrintRed("SF", "Error Reading Config.txt file")
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var (
		headless                 bool
		closeDriverWhileSleeping bool
		startThread              int
		startThreadInterval      int
		showForgeReward          bool
		showTotalForgeBalance    bool
		printSystemUsage         bool
		checkSystemUsageFreq     int
		autoSelectCard           bool
		autoSelectHero           bool
		autoSelectSleepTime      int
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
			case "CLOSE_DRIVER_WHILE_SLEEPING":
				closeDriverWhileSleeping = value == "true"
			case "START_THREAD":
				startThread, _ = strconv.Atoi(value)
			case "START_THREAD_INTERVAL":
				startThreadInterval, _ = strconv.Atoi(value)
			case "SHOW_FORGE_REWARD":
				showForgeReward = value == "true"
			case "SHOW_TOTAL_FORGE_BALANCE":
				showTotalForgeBalance = value == "true"
			case "PRINT_SYSTEM_USAGE":
				printSystemUsage = value == "true"
			case "CHECK_SYSTEM_USAGE_FREQUENCY":
				checkSystemUsageFreq, _ = strconv.Atoi(value)
			case "AUTO_SELECT_CARD":
				autoSelectCard = value == "true"
			case "AUTO_SELECT_SLEEPTIME":
				autoSelectSleepTime, _ = strconv.Atoi(value)
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
	return headless, closeDriverWhileSleeping, startThread, startThreadInterval, showForgeReward, showTotalForgeBalance, printSystemUsage, checkSystemUsageFreq, autoSelectCard, autoSelectHero, autoSelectSleepTime, splinterforgeAPIEndpoint, splinterlandAPIEndpoint, publicAPIEndpoint
}
func Battle(wd selenium.WebDriver, userName string, bossId string, heroesType string, cardSelection []spstruct.CardSelection, autoSelectHero bool, autoSelectCard bool, splinterlandAPIEndpoint string, publicAPIEndpoint string) {
	MemoStatus := true
	Unexpected := false
	time.Sleep(15 * time.Second)
	bossName, bossIdToSelect, _ := bossSelect(userName, bossId, wd)
	heroSelect(heroesType, userName, wd, autoSelectHero, publicAPIEndpoint, bossName)
	bossSelect(userName, bossIdToSelect, wd)
	cardSelection, _, _ = autoSelectCards(cardSelection, bossName, userName, splinterlandAPIEndpoint, publicAPIEndpoint, autoSelectCard)
	seletedNumOfSummoners := 1
	el, _ := wd.FindElement(selenium.ByXPATH, "/html/body/app/div[1]/slcards/div[5]/section[1]/div/div[1]/div[2]/button")
	el.Click()
	PrintWhite(userName, "Participating in battles...")
	printData := [][]string{}
	selectResult := true
	for _, selection := range cardSelection {
		for i, playingSummoner := range selection.PlayingSummoners {
			result := selectSummoners(userName, seletedNumOfSummoners, playingSummoner.PlayingSummonersDiv, wd)
			if result {
				seletedNumOfSummoners++
				printData = append(printData, []string{fmt.Sprintf("Summoners #%d", i+1), playingSummoner.PlayingSummonersID, playingSummoner.PlayingSummonersName, "success"})
			}
		}
		seletedNumOfMonsters := 1
		for j, playingMonster := range selection.PlayingMonsters {
			result := selectMonsters(userName, seletedNumOfMonsters, playingMonster.PlayingMontersDiv, wd)
			if result {
				seletedNumOfMonsters++
				printData = append(printData, []string{fmt.Sprintf("Monsters #%d", j+1), playingMonster.PlayingMonstersID, playingMonster.PlayingMonstersName, "success"})
			} else {
				printData = append(printData, []string{fmt.Sprintf("Monsters #%d", j+1), playingMonster.PlayingMonstersID, playingMonster.PlayingMonstersName, "error"})
				selectResult = false
			}
		}
	}
	el, _ = wd.FindElement(selenium.ByXPATH, "/html/body/app/div[1]/slcards/div[5]/div[2]/div[1]/div[1]/button/span")
	manaused, _ := el.Text()
	fmt.Println("mana", manaused)
	el, _ = wd.FindElement(selenium.ByXPATH, "/html/body/app/div[1]/slcards/div[5]/div[2]/div[3]/div[1]/button/span")
	manatotal, _ := el.Text()
	fmt.Println("mana", manatotal)
	printResultBox(userName, printData, selectResult)
	el, _ = wd.FindElement(selenium.ByXPATH, "/html/body/app/div[1]/slcards/div[5]/button[1]/div[2]/span")
	el.Click()

	name, _ := wd.ExecuteScript("return localStorage.getItem('forge:username');", nil)
	key, _ := wd.ExecuteScript("return localStorage.getItem('forge:key');", nil)

	returnJsonResult := false
	fitPostData := spstruct.FitBossPostData{}
	fitRes := spstruct.FitBossRequestsData{}

	for {
		d, _ := wd.Log("performance")
		for _, dd := range d {
			if strings.Contains(dd.Message, "https://splinterforge.io/boss/fight_boss") && strings.Contains(dd.Message, "\"method\":\"Network.requestWillBeSent\"") {
				json.Unmarshal([]byte(dd.Message), &fitRes)
				json.Unmarshal([]byte(fitRes.Message.Params.Request.PostData), &fitPostData)
				fmt.Println(fitPostData.Username, "success!")
				returnJsonResult = true
				break
			}
		}

		if returnJsonResult == true {
			break
		} else {
			time.Sleep(2 * time.Second)
			continue
		}
	}
	wd.Close()
	for {
		fitPostData.Memo = "dw9d59wd5298w5d2985dw9"
		reFit, err := grequests.Post("https://splinterforge.io/boss/fight_boss", &grequests.RequestOptions{
			JSON: fitPostData,
			Headers: map[string]string{
				"Content-Type": fitRes.Message.Params.Request.Headers.ContentType,
				"User-Agent":   fitRes.Message.Params.Request.Headers.UserAgent,
			},
		})
		if err == nil {
			fmt.Println(reFit.String())
			if strings.Contains(reFit.String(), "not enough mana!") {
				time.Sleep(1 * time.Hour)
				continue
			} else if strings.Contains(reFit.String(), "decoded message was invalid") {
				MemoStatus = false
				PrintRed(userName, "Memo Error, Restarting...")
				break
			} else if strings.Contains(reFit.String(), "totalDmg") && strings.Contains(reFit.String(), "points") {
				var fitReturnData = spstruct.FitReturnData{}
				json.Unmarshal(reFit.Bytes(), &fitReturnData)
				fmt.Println("fitReturnData.Points > ", fitReturnData.Points)
				fmt.Println("fitReturnData.TotalDmg > ", fitReturnData.TotalDmg)
				fmt.Println("fitReturnData.Rewards[0].Qty > ", fitReturnData.Rewards[0].Qty)
				fmt.Println("fitReturnData.Rewards[1].Qty > ", fitReturnData.Rewards[1].Qty)
				time.Sleep(5 * time.Second)
				res, _ := grequests.Post("https://splinterforge.io/users/keyLogin", &grequests.RequestOptions{
					JSON: map[string]string{
						"username": name.(string),
						"key":      key.(string),
					},
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				})
				var powerRes = spstruct.KeyLoginResData{}
				json.Unmarshal(res.Bytes(), &powerRes)
				count, _ := CaculateTimeDiff(powerRes.Stamina.Last)
				PrintWhite(userName, fmt.Sprintf("powerRes.Stamina.Last = %s", powerRes.Stamina.Last))
				PrintWhite(userName, fmt.Sprintf("count = %s", count))
				PrintWhite(userName, fmt.Sprintf("powerRes.Stamina.Current = %s", strconv.Itoa(powerRes.Stamina.Current)))
				PrintWhite(userName, fmt.Sprintf("powerRes.Stamina.Max = %s", strconv.Itoa(powerRes.Stamina.Max)))
				PrintWhite(userName, fmt.Sprintf("(powerRes.Stamina.Current + count)/20 = %s", (powerRes.Stamina.Current+count)/20))
				PrintWhite(userName, fmt.Sprintf("powerRes.Stamina.Max / 20 = %s", strconv.Itoa(powerRes.Stamina.Max/20)))
				time.Sleep(25 * time.Second)
				continue
			} else {
				Unexpected = true
				PrintRed(userName, "Unexpected Error, Restarting...")
				break
			}
		} else {
			fmt.Println("RequestsFit Err > ", err)
		}
	}
	if MemoStatus == false || Unexpected == true {
		for _, v := range accountLists {
			if v.UserName == userName {
				initializeDriver(v, headless, closeDriverWhileSleeping, showForgeReward, showTotalForgeBalance, autoSelectCard, autoSelectHero, autoSelectSleepTime, splinterforgeAPIEndpoint, splinterlandAPIEndpoint, publicAPIEndpoint)
			}
		}
	}
}

func initializeDriver(userData spstruct.UserData, headless bool, closeDriverWhileSleeping bool, showForgeReward bool, showTotalForgeBalance bool, autoSelectCard bool, autoSelectHero bool, autoSelectSleepTime int, splinterforgeAPIEndpoint string, splinterlandAPIEndpoint string, publicAPIEndpoint string) {
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
	// Start a new ChromeDriver instance
	printLog := false
	wd, err := selenium.NewChromeDriverService(SeleniumDriverCheck.AutoDownload_ChromeDriver(printLog), 9515)
	if err != nil {
		fmt.Printf("Failed to create ChromeDriver service: %s\n", err)
		os.Exit(1)
	}
	defer wd.Stop()

	// Create a new WebDriver instance
	driver, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", 9515))
	if err != nil {
		fmt.Printf("Failed to create WebDriver: %s\n", err)
		os.Exit(1)
	}
	defer driver.Quit()

	userName := userData.UserName
	postingKey := userData.PostingKey
	bossId := userData.BossID
	heroesType := userData.HeroesType
	cardSelection := userData.CardSelection
	login(userName, postingKey, driver)
	checkPopUp(driver, 1000)
	Battle(driver, userName, bossId, heroesType, cardSelection, autoSelectHero, autoSelectCard, splinterlandAPIEndpoint, publicAPIEndpoint)
}

func countLines(filePath string) (int, error) {
	//读取text文件
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

var (
	accountLists                                                                                                                                                                                                                                                                    = []spstruct.UserData{}
	w                                                                                                                                                                                                                                                                               = &sync.WaitGroup{}
	headless, closeDriverWhileSleeping, startThread, startThreadInterval, showForgeReward, showTotalForgeBalance, printSystemUsage, checkSystemUsageFreq, autoSelectCard, autoSelectHero, autoSelectSleepTime, splinterforgeAPIEndpoint, splinterlandAPIEndpoint, publicAPIEndpoint = getConfig("config/config.txt")
)

func initializeUserData() {
	lineCount, errCountLines := countLines("config/accounts.txt")
	if errCountLines == nil && lineCount > 1 {
		for i := 0; i < lineCount-1; i++ {
			w.Add(1)
			go func(num int) {
				userName, postingKey, heroesType, bossId, cardSelection, timeSleepInMinute := initializeAccount(num + 1)
				accountLists = append(accountLists, spstruct.UserData{
					UserName:          userName,
					PostingKey:        postingKey,
					BossID:            bossId,
					HeroesType:        heroesType,
					CardSelection:     cardSelection,
					TimeSleepInMinute: timeSleepInMinute,
				})

				w.Done()
			}(i)
		}
		w.Wait()
		printConfigSettings(lineCount-1, headless, closeDriverWhileSleeping, startThread, startThreadInterval, showForgeReward, showTotalForgeBalance, printSystemUsage, checkSystemUsageFreq, autoSelectCard, autoSelectHero, autoSelectSleepTime)
		for _, v := range accountLists {
			initializeDriver(v, headless, closeDriverWhileSleeping, showForgeReward, showTotalForgeBalance, autoSelectCard, autoSelectHero, autoSelectSleepTime, splinterforgeAPIEndpoint, splinterlandAPIEndpoint, publicAPIEndpoint)
			break
		}
	} else {
		fmt.Print("Please add accounts in accounts.txt\n")
		os.Exit(1)
	}
}

func CaculateTimeDiff(oldTime string) (int, error) {
	now := time.Now()

	// 待计算的时间参数（string格式）
	t, err := time.Parse(time.RFC3339, "2023-03-03T08:04:10.316Z")
	if err != nil {
		return 0, err
	}
	// 转换为Unix时间戳并计算差值
	diff := int(now.Unix() - t.Unix())
	// 转换为分钟
	diffInMinutes := diff / 60
	return diffInMinutes, nil

}

func main() {
	var stats1 runtime.MemStats
	runtime.ReadMemStats(&stats1)
	start := time.Now()

	// Call the function that we want to measure

	initializeUserData()

	// Measure CPU usage after the function is called
	elapsed := time.Since(start)
	var stats2 runtime.MemStats
	runtime.ReadMemStats(&stats2)

	// Calculate and display CPU usage statistics
	cpuTime := time.Duration(stats2.Sys - stats1.Sys)
	fmt.Printf("CPU usage: %.2f%%\n", (float64(cpuTime)/float64(elapsed))*100.0)
	time.Sleep(20 * time.Second)
}
