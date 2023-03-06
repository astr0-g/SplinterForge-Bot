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
	"github.com/theckman/yacspin"

	"splinterforge/spstruct"
)

var (
	accountLists                                                                                                                                                                          = []spstruct.UserData{}
	q                                                                                                                                                                                     = &sync.WaitGroup{}
	r                                                                                                                                                                                     = &sync.WaitGroup{}
	w                                                                                                                                                                                     = &sync.WaitGroup{}
	s                                                                                                                                                                                     = &sync.WaitGroup{}
	headless, startThread, showForgeReward, showAccountDetails, autoSelectCard, autoSelectHero, autoSelectSleepTime, splinterforgeAPIEndpoint, splinterlandAPIEndpoint, publicAPIEndpoint = getConfig("config/config.txt")
)

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
func printInfo() {
	fmt.Println("+-------------------------------------------------+")
	fmt.Println("|          Welcome to SplinterForge Bot!          |")
	fmt.Println("|          Open source Github repository          |")
	fmt.Println("|   https://github.com/Astr0-G/SplinterForge-Bot  |")
	fmt.Println("|                  Discord server                 |")
	fmt.Println("|          https://discord.gg/pm8SGZkYcD          |")
	fmt.Println("+-------------------------------------------------+")
}
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
func PrintAccountDetails(userName string, name interface{}, key interface{}) {
	res, _ := grequests.Post(fmt.Sprintf("%s/users/keyLogin", splinterforgeAPIEndpoint), &grequests.RequestOptions{
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
	PrintWhite(userName, fmt.Sprintf("Account balance %s Forge, current mana %s / %s.", strconv.FormatFloat(powerRes.Sc.Balance, 'f', 2, 64), strconv.Itoa(powerRes.Stamina.Current), strconv.Itoa(powerRes.Stamina.Max)))
}
func printConfigSettings(totalAccounts int, headless bool, startThread int, showForgeReward bool, showAccountDetails bool, autoSelectCard bool, autoSelectHero bool, autoSelectSleepTime bool) {
	data := [][]string{
		{"TOTAL_ACCOUNTS_LOADED", fmt.Sprint(totalAccounts)},
		{"HEADLESS", fmt.Sprint(headless)},
		{"THREADING", fmt.Sprint(startThread)},
		{"SHOW_FORGE_REWARD", fmt.Sprint(showForgeReward)},
		{"SHOW_ACCOUNT_DETAILS", fmt.Sprint(showAccountDetails)},
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
func getConfig(filePath string) (bool, int, bool, bool, bool, bool, bool, string, string, string) {
	file, err := os.Open(filePath)
	if err != nil {
		PrintRed("SF", "Error Reading Config.txt file")
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var (
		headless                 bool
		startThread              int
		showForgeReward          bool
		showAccountDetails       bool
		autoSelectCard           bool
		autoSelectHero           bool
		autoSelectSleepTime      bool
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
				startThread, _ = strconv.Atoi(value)
			case "SHOW_FORGE_REWARD":
				showForgeReward = value == "true"
			case "SHOW_ACCOUNT_DETAILS":
				showAccountDetails = value == "true"
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
	return headless, startThread, showForgeReward, showAccountDetails, autoSelectCard, autoSelectHero, autoSelectSleepTime, splinterforgeAPIEndpoint, splinterlandAPIEndpoint, publicAPIEndpoint
}
func getLines(filePath string) (int, error) {
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
func getAccountDetails(name interface{}, key interface{}) int {
	res, _ := grequests.Post(fmt.Sprintf("%s/users/keyLogin", splinterforgeAPIEndpoint), &grequests.RequestOptions{
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
	reusltTime, _ := getTimeDiff(powerRes.Stamina.Last)
	CurrentStamina := powerRes.Stamina.Current + reusltTime
	if CurrentStamina > powerRes.Stamina.Max {
		CurrentStamina = powerRes.Stamina.Max
	}
	return CurrentStamina
}
func getTimeDiff(oldTime string) (int, error) {
	now := time.Now()
	t, err := time.Parse(time.RFC3339, oldTime)
	if err != nil {
		return 0, err
	}
	diff := int(now.Unix() - t.Unix())
	diffInMinutes := diff / 60
	return diffInMinutes, nil
}
func GetReponseBody(sessionId string, requestId string, userName string) (string, error) {
	res, err := grequests.Post(fmt.Sprintf("http://localhost:9515/wd/hub/session/%s/goog/cdp/execute", sessionId),
		&grequests.RequestOptions{
			JSON: map[string]interface{}{
				"cmd": "Network.getResponseBody",
				"params": map[string]string{
					"requestId": requestId,
				},
			},
		})
	if err == nil {
		var fitResponseData = spstruct.GetResponseBody{}
		json.Unmarshal(res.Bytes(), &fitResponseData)
		fmt.Println(fitResponseData.Value.Body)
		// fmt.Println(fitReturnData.TotalDmg)
		// fmt.Println(fitReturnData.Points)
		// fmt.Println(fitReturnData.Rewards[0])
		// fmt.Println(fitReturnData.Rewards[1])
		// if showForgeReward {
		// 	PrintYellow(userName, fmt.Sprintf("You made battle damage %s, battle points %s, reward Forgium %0.3f, reward Electrum %0.2f.", strconv.Itoa(fitReturnData.TotalDmg), strconv.Itoa(fitReturnData.Points), fitReturnData.Rewards[0].Qty, fitReturnData.Rewards[1].Qty))
		// }
		time.Sleep(5 * time.Second)
		return res.String(), nil

	} else {
		fmt.Println("GetReponseBody error > ", err)
		return "", err
	}

}
func fetchselectHero(publicAPIEndpoint string, bossName string) (string, error) {
	bossID := fetchBossID(bossName)

	requestBody := map[string]interface{}{
		"bossId": bossID,
	}

	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	endpoint := fmt.Sprintf("%s/heroselection", publicAPIEndpoint)
	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(requestBodyBytes))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

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
	url := fmt.Sprintf("%s/boss/getBosses", splinterforgeAPIEndpoint)

	resp, err := http.Get(url)
	if err != nil {
		glog.Fatal(err)
	}
	defer resp.Body.Close()

	var responseData []map[string]interface{}
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&responseData)

	for _, bossData := range responseData {
		if strings.EqualFold(bossData["name"].(string), bossName) {
			return bossData["id"].(string)
		}
	}

	return ""
}
func fetchPlayerCard(userName string, splinterlandAPIEndpoint string) ([]int, error) {
	url := fmt.Sprintf("%s/cards/collection/%s", splinterlandAPIEndpoint, userName)
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
	err = json.Unmarshal(jsonString, &cardList)
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
	return cardDetailIDs, nil
}
func fetchBattleCards(bossName string, userName string, splinterlandAPIEndpoint string, publicAPIEndpoint string) (string, error) {
	cardsId, err := fetchPlayerCard(userName, splinterlandAPIEndpoint)
	if err != nil {
		return "", err
	}

	bossId := fetchBossID(bossName)

	postData := spstruct.BattleCardsRequestBody{
		BossId:   bossId,
		BossName: bossName,
		Team:     cardsId,
	}

	url := fmt.Sprintf("%s/teamselection/", publicAPIEndpoint)
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

	if !gresponse.Ok {
		return "", fmt.Errorf("HTTP error: %d - %s", gresponse.StatusCode, gresponse.String())
	}

	var responseData map[string]interface{}
	err = json.Unmarshal(gresponse.Bytes(), &responseData)
	if err != nil {
		return "", err
	}

	jsonResponse, err := json.Marshal(responseData)
	if err != nil {
		return "", err
	}

	return string(jsonResponse), nil
}
func selectCards(cardSelection []spstruct.CardSelection, bossName string, userName string, splinterlandAPIEndpoint string, publicAPIEndpoint string, autoSelectCard bool) ([]spstruct.CardSelection, bool, error) {
	if autoSelectCard {
		PrintWhite(
			userName, fmt.Sprintf("Auto selecting playing cards for desire boss: %s", bossName))
		playingDeck, err := fetchBattleCards(bossName, userName, splinterlandAPIEndpoint, publicAPIEndpoint)
		if err != nil {
			PrintYellow(userName, "API Error, Fetch Battle cards failed.")
			return cardSelection, false, err
		}
		playingDeckBytes := []byte(playingDeck)
		var playingDeckMap map[string]interface{}
		err = json.Unmarshal(playingDeckBytes, &playingDeckMap)
		if err != nil {
			return cardSelection, false, err
		}
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
			return cardSelectionList, true, nil
		} else {
			PrintYellow(userName, "Auto selecting playing cards deck failed, you might have too less cards in the account, will continue play with your card setting.")
			return cardSelection, false, nil
		}
	} else {
		return cardSelection, false, nil
	}
}
func selectHero(heroesType string, userName string, wd selenium.WebDriver, auto_select_hero bool, publicAPIEndpoint string, bossName string) {
	checkPopUp(wd, 1000)
	heroTypes := [3]string{"Warrior", "Wizard", "Ranger"}
	if auto_select_hero {
		hero_type, err := fetchselectHero(publicAPIEndpoint, bossName)
		if err == nil {
			PrintWhite(userName, fmt.Sprintf("Auto selecting heroes type: %s for desired boss: %s", hero_type, bossName))
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
		PrintWhite(userName, "Selecting heroes type...")
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
	PrintWhite(userName, fmt.Sprintf("Selected hero type: %s", hero_type))
}
func selectBoss(userName string, bossIdToSelect string, wd selenium.WebDriver) (string, string, error) {
	wd.SetImplicitWaitTimeout(2 * time.Second)
	for {
		el, _ := wd.FindElement(selenium.ByXPATH, "/html/body/app/div[1]/div[1]/app-header/section/div[4]/div[2]/div[1]/a[5]/div[1]")
		el.Click()
		time.Sleep(1 * time.Second)
		bossSelector := fmt.Sprintf("//div[@tabindex='%s']", bossIdToSelect)
		el, _ = wd.FindElement(selenium.ByXPATH, bossSelector)
		el.Click()
		element, err := wd.FindElement(selenium.ByXPATH, "/html/body/app/div[1]/slcards/div[5]/section[1]/div/div[1]/div[2]/button")
		if err == nil {
			text, err := element.Text()
			if err == nil {
				if text != "BOSS IS DEAD" {
					time.Sleep(1 * time.Second)

					if element, err := wd.FindElement(selenium.ByXPATH, "/html/body/app/div[1]/slcards/div[5]/section[1]/div/div[1]/div[1]/div[2]/div[3]/h3"); err == nil {
						bossName, _ := element.Text()
						return bossName, bossIdToSelect, nil
					}
				} else {
					PrintRed(userName, "The selected boss has been defeated, selecting another one automatically...")

					if bossIdToSelect < "17" {
						bossIdInt, err := strconv.Atoi(bossIdToSelect)
						if err != nil {

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
func selectSummoners(userName string, seletedNumOfSummoners int, cardDiv string, wd selenium.WebDriver) bool {
	scroolTime := 0
	clickedTime := 0
	result := false
	time.Sleep(1 * time.Second)
	for scroolTime < 5 && clickedTime < 5 {
		el, err := wd.FindElement(selenium.ByXPATH, cardDiv)
		if err != nil {
			continue
		} else {
			el.Click()
			time.Sleep(1 * time.Second)
			checkCardDiv := fmt.Sprintf("/html/body/app/div[1]/slcards/div[5]/section[1]/div/div[1]/div[2]/div[1]/div[2]/div[%s]", strconv.Itoa(seletedNumOfSummoners))
			success, _ := DriverwaitForElement(wd, checkCardDiv)
			if success {
				result = true
				break
			} else {
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
			wd.ExecuteScript("window.scrollBy(0, 450)", nil)
			scroolTime++
			continue
		} else {
			el.Click()
			checkCardDiv := fmt.Sprintf("/html/body/app/div[1]/slcards/div[5]/section[1]/div/div[1]/div[2]/div[2]/div[2]/div[%s]", strconv.Itoa(seletedNumOfMonsters))
			success, _ := DriverwaitForElement(wd, checkCardDiv)
			if success {
				result = true
				break
			} else {
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
func DriverElementWaitAndClick(wd selenium.WebDriver, xpath string) {
	byXpath := selenium.ByXPATH
	for {
		element, err := wd.FindElement(byXpath, xpath)
		isEnabled, err1 := element.IsEnabled()
		if err == nil && err1 == nil && isEnabled {
			err = element.Click()
			if err != nil {
				continue
			} else {
				return
			}

		}
	}
}
func DriverwaitForElement(wd selenium.WebDriver, xpath string) (bool, error) {
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
func checkPopUp(wd selenium.WebDriver, millisecond int) {
	defer func() {
	}()
	if element, err := wd.FindElement(selenium.ByXPATH, "/html/body/app/div[1]/div[1]/app-header/success-modal/section/div[1]/div[4]/div/button"); err == nil {
		if err = element.Click(); err != nil {

		}
	}
	if element, err := wd.FindElement(selenium.ByXPATH, "/html/body/app/div[1]/login-modal/div/div/div/div[2]/div[3]/button"); err == nil {
		if err = element.Click(); err != nil {

		}
	}
	duration := time.Duration(millisecond) * time.Millisecond
	time.Sleep(duration)
}
func checklogin(userName string, wd selenium.WebDriver) bool {
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
func accountRestartCoroutine(wait bool, userName string) {
	for _, v := range accountLists {
		if v.UserName == userName {
			initializeDriver(wait, v, headless, showForgeReward, showAccountDetails, autoSelectCard, autoSelectHero, autoSelectSleepTime, splinterforgeAPIEndpoint, splinterlandAPIEndpoint, publicAPIEndpoint)
		}
	}
}
func accountLogin(userName string, postingKey string, wd selenium.WebDriver) bool {
	err := wd.SetImplicitWaitTimeout(5 * time.Second)
	if err != nil {
		panic(err)
	}
	DriverGet("chrome-extension://jcacnejopjdphbnjgfaaobbfafkihpep/popup.html", wd)

	DriverElementWaitAndClick(wd, "/html/body/div/div/div[4]/div[2]/div[5]/button")

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
	el, err = wd.FindElement(selenium.ByXPATH, "/html/body/div/div/div[4]/div")
	if err != nil {
		PrintRed(userName, "Login failure! Please check your accounts.txt file or it is a server error.")
		return false
	}
	hiveKeyChainLogin, _ := el.Text()
	if hiveKeyChainLogin != "HIVE KEYCHAIN" {
		PrintRed(userName, "Login failure! Please check your accounts.txt file or it is a server error.")
		return false
	} else {
		err = wd.ResizeWindow("bigger", 1565, 1080)
		if err != nil {
			println("can not change size")
		}
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
		PrintGreen(userName, "Login successful!")
		return true
	}
}
func accountBattle(wait bool, wd selenium.WebDriver, userName string, bossId string, heroesType string, timeSleepInMinute int, cardSelection []spstruct.CardSelection, autoSelectHero bool, autoSelectCard bool, autoSelectSleepTime bool, splinterlandAPIEndpoint string, publicAPIEndpoint string) {
	CookiesStatus := true
	Unexpected := false

	name, _ := wd.ExecuteScript("return localStorage.getItem('forge:username');", nil)
	key, _ := wd.ExecuteScript("return localStorage.getItem('forge:key');", nil)

	starttimestamp := time.Now().Unix()
	bossName, bossIdToSelect, _ := selectBoss(userName, bossId, wd)
	selectHero(heroesType, userName, wd, autoSelectHero, publicAPIEndpoint, bossName)
	selectBoss(userName, bossIdToSelect, wd)
	cardSelection, _, _ = selectCards(cardSelection, bossName, userName, splinterlandAPIEndpoint, publicAPIEndpoint, autoSelectCard)
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
	getAccountDetails(name, key)
	el, _ = wd.FindElement(selenium.ByXPATH, "/html/body/app/div[1]/slcards/div[5]/div[2]/div[1]/div[1]/button/span")
	mana, _ := el.Text()
	manaused, _ := strconv.Atoi(mana)
	CurrentStamina := getAccountDetails(name, key)
	if CurrentStamina > manaused && manaused >= 15 {
		printResultBox(userName, printData, selectResult)
		battletimestamp := time.Now().Unix()
		if battletimestamp-starttimestamp < 30 {
			time.Sleep((time.Duration(battletimestamp) - time.Duration(starttimestamp) + 1) * time.Second)
		}
		el, _ = wd.FindElement(selenium.ByXPATH, "/html/body/app/div[1]/slcards/div[5]/button[1]/div[2]/span")
		el.Click()
		returnJsonResult := false
		postJsonResult := false
		fitPostData := spstruct.FitBossPostData{}
		fitRes := spstruct.FitBossRequestsData{}
		for {
			netLogs, _ := wd.Log("performance")
			for _, netLog := range netLogs {
				if strings.Contains(netLog.Message, fmt.Sprintf("%s/boss/fight_boss", splinterforgeAPIEndpoint)) && strings.Contains(netLog.Message, "\"method\":\"Network.requestWillBeSent\"") {
					json.Unmarshal([]byte(netLog.Message), &fitRes)
					json.Unmarshal([]byte(fitRes.Message.Params.Request.PostData), &fitPostData)
					fmt.Println(fitRes.Message.Params.RequestID)
					returnJsonResult = true
				}
				if strings.Contains(netLog.Message, fmt.Sprintf("%s/boss/fight_boss", splinterforgeAPIEndpoint)) && strings.Contains(netLog.Message, "\"method\":\"Network.responseReceived\"") {
					var GetResponseBody = spstruct.GetResponseBody{}
					var GetRewardBody = spstruct.CDPFitReturnData{}
					resString, err := GetReponseBody(wd.SessionID(), fitRes.Message.Params.RequestID, userName)
					if err == nil {
						json.Unmarshal([]byte(resString), &GetResponseBody)
						fmt.Println("resString", resString)
						fmt.Println("GetResponseBody", GetResponseBody)
						fmt.Println("GetResponseBody.Value", GetResponseBody.Value)
						fmt.Println("GetResponseBody.Value.Body", GetResponseBody.Value.Body)
						PrintWhite(userName, "Battle was successful!")
						err1 := json.Unmarshal([]byte(GetResponseBody.Value.Body), &GetRewardBody)
						PrintYellow(userName, fmt.Sprintf("You made battle damage %s, battle points %s, reward Forgium %0.3f, reward Electrum %0.2f.", strconv.Itoa(GetRewardBody.TotalDmg), strconv.Itoa(GetRewardBody.Points), GetRewardBody.Rewards[0].Qty, GetRewardBody.Rewards[1].Qty))
						fmt.Println("--------------------------------------------------------------------------------------")
						fmt.Println("TotalDmg", GetRewardBody.TotalDmg, err1)
						fmt.Println("Rewards 1 > ", GetRewardBody.Rewards[0], err1)
						fmt.Println("Rewards 1 Details > ", GetRewardBody.Rewards[0].Name, GetRewardBody.Rewards[0].Qty, GetRewardBody.Rewards[0].Type)
						fmt.Println("Rewards 2 > ", GetRewardBody.Rewards[1], err1)
						fmt.Println("Rewards 2 Details > ", GetRewardBody.Rewards[1].Name, GetRewardBody.Rewards[1].Qty, GetRewardBody.Rewards[1].Type)
						fmt.Println("Points", GetRewardBody.Points, err1)
						postJsonResult = true
					}
				}
			}
			if returnJsonResult == true && postJsonResult == true {
				break
			} else {
				time.Sleep(2 * time.Second)
				continue
			}
		}
		if showForgeReward {
			DriverElementWaitAndClick(wd, "/html/body/app/div[1]/slcards/div[4]/div[2]/button[2]")
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
					PrintYellow(userName, fmt.Sprintf("You made battle damage %s, battle points %s, reward Forgium %0.3f, reward Electrum %0.2f.", resultdmg, resultpoints, forgium, electrum))
					break
				} else {
					time.Sleep(5 * time.Second)
					continue
				}
			}
		}
		if showAccountDetails {
			PrintAccountDetails(userName, name, key)
		}
		wd.Close()
		PrintWhite(userName, "Successful generated Cookies, the account will continue play with this setup.")
		s.Add(1)
		go func() {
			for {
				if autoSelectSleepTime {
					PrintWhite(userName, fmt.Sprintf("this account will enter a state of inactivity for %s minutes based on auto selected info.", mana))
					time.Sleep(time.Duration(manaused) * time.Minute)
				} else {
					PrintWhite(userName, fmt.Sprintf("According to your configuration, this account will enter a state of inactivity for %s minutes.", strconv.Itoa(timeSleepInMinute)))
					time.Sleep(time.Duration(timeSleepInMinute) * time.Minute)
				}
				PrintWhite(userName, "Participating in battles with current setup...")
				reFit, err := grequests.Post(fmt.Sprintf("%s/boss/fight_boss", splinterforgeAPIEndpoint), &grequests.RequestOptions{
					JSON: fitPostData,
					Headers: map[string]string{
						"Content-Type": fitRes.Message.Params.Request.Headers.ContentType,
						"User-Agent":   fitRes.Message.Params.Request.Headers.UserAgent,
					},
				})
				if err == nil {
					printResultBox(userName, printData, selectResult)
					if strings.Contains(reFit.String(), "not enough mana!") {
						PrintYellow(userName, "Insufficient stamina, entering a rest state of inactivity for 1 hour...")
						time.Sleep(1 * time.Hour)
						continue
					} else if strings.Contains(reFit.String(), "decoded message was invalid") {
						CookiesStatus = false
						PrintRed(userName, "Cookies Error, Restarting...")
						break
					} else if strings.Contains(reFit.String(), "totalDmg") && strings.Contains(reFit.String(), "points") {
						PrintWhite(userName, "Battle was successful!")
						var fitReturnData = spstruct.FitReturnData{}
						json.Unmarshal(reFit.Bytes(), &fitReturnData)
						if showForgeReward {
							PrintYellow(userName, fmt.Sprintf("You made battle damage %s, battle points %s, reward Forgium %0.3f, reward Electrum %0.2f.", strconv.Itoa(fitReturnData.TotalDmg), strconv.Itoa(fitReturnData.Points), fitReturnData.Rewards[0].Qty, fitReturnData.Rewards[1].Qty))
						}
						time.Sleep(5 * time.Second)
						if showAccountDetails {
							PrintAccountDetails(userName, name, key)
						}
						continue
					} else {
						Unexpected = true
						PrintRed(userName, "Unexpected Error, Restarting...")
						break
					}
				} else {
					PrintRed(userName, "Game server API error, Restarting...")
					Unexpected = true
					break
				}
			}
			if !CookiesStatus || Unexpected {
				accountRestartCoroutine(wait, userName)
				s.Done()
			}

		}()
		if wait == true {
			w.Done()
		} else {
			q.Done()
		}
	} else {
		if CurrentStamina < manaused {
			wd.Close()
			PrintYellow(userName, "Insufficient stamina, entering a rest state of inactivity for 1 hour...")
			time.Sleep(1 * time.Hour)
			accountRestartCoroutine(wait, userName)
		} else if manaused < 15 {
			wd.Close()
			PrintYellow(userName, "Card Selected not meet 15 mana requirements, restarting...")
			accountRestartCoroutine(wait, userName)
		}
	}
}
func initializeDriver(Wait bool, userData spstruct.UserData, headless bool, showForgeReward bool, showAccountDetails bool, autoSelectCard bool, autoSelectHero bool, autoSelectSleepTime bool, splinterforgeAPIEndpoint string, splinterlandAPIEndpoint string, publicAPIEndpoint string) {
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
	loginResult := accountLogin(userData.UserName, userData.PostingKey, driver)
	if loginResult {
		checkPopUp(driver, 1000)
		accountBattle(Wait, driver, userData.UserName, userData.BossID, userData.HeroesType, userData.TimeSleepInMinute, userData.CardSelection, autoSelectHero, autoSelectCard, autoSelectSleepTime, splinterlandAPIEndpoint, publicAPIEndpoint)
	} else {
		driver.Close()
		PrintYellow(userData.UserName, "Retrying in 30 seconds...")
		time.Sleep(30 * time.Second)
		accountRestartCoroutine(Wait, userData.UserName)
	}
}
func initializeAccount(accountNo int) (string, string, string, string, []spstruct.CardSelection, int) {

	userName, postingKey, err := getAccountData("config/accounts.txt", accountNo)
	if err != nil || userName == "" || postingKey == "" {
		PrintRed("ERROR", "Error in loading accounts.txt, please add username or posting key and try again.")
	}
	heroesType, bossId, playingSummoners, playingMonster, timeSleepInMinute, err := getCardSettingData("config/cardSettings.txt", accountNo)
	if err != nil {
		PrintRed("ERROR", "Error loading cardSettings.txt file")
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

	return userName, postingKey, heroesType, bossId, cardSelectionList, timeSleepInMinute
}

func initializeUserData() {
	cfg := yacspin.Config{
		Frequency:       100 * time.Millisecond,
		CharSet:         yacspin.CharSets[59],
		Suffix:          " Initialization Driver",
		SuffixAutoColon: true,
		StopCharacter:   " ✓",
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
		StopCharacter:   " ✓",
		StopColors:      []string{"fgGreen"},
		StopMessage:     "Success!",
	}
	spinner, _ = yacspin.New(cfg)
	spinner.Start()
	spinner.Message("reading accounts.txt...")
	lineCount, errCountLines := getLines("config/accounts.txt")
	time.Sleep(500 * time.Millisecond)
	spinner.Message("reading cardSettings.txt...")
	if errCountLines == nil && lineCount > 1 {
		for i := 0; i < lineCount-1; i++ {
			r.Add(1)
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

				r.Done()
			}(i)
		}
		r.Wait()
		spinner.Message("reading config.txt...")
		time.Sleep(500 * time.Millisecond)
		spinner.Stop()
		printConfigSettings(lineCount-1, headless, startThread, showForgeReward, showAccountDetails, autoSelectCard, autoSelectHero, autoSelectSleepTime)
		//判断 > 如果当前账户数小于启动线程数，那么就按照账户数启动线程
		if len(accountLists) <= startThread {
			for i := 0; i < len(accountLists); i++ {
				q.Add(1)
				go initializeDriver(false, accountLists[i], headless, showForgeReward, showAccountDetails, autoSelectCard, autoSelectHero, autoSelectSleepTime, splinterforgeAPIEndpoint, splinterlandAPIEndpoint, publicAPIEndpoint)
			}
			q.Wait()
		} else {
			//如果当前账户数大于启动线程数，那么就按照启动线程数启动线程
			for i := 0; i < len(accountLists); i += startThread {
				for j := 0; j < startThread; j++ {
					w.Add(1)
					go initializeDriver(true, accountLists[i+j], headless, showForgeReward, showAccountDetails, autoSelectCard, autoSelectHero, autoSelectSleepTime, splinterforgeAPIEndpoint, splinterlandAPIEndpoint, publicAPIEndpoint)
				}
				w.Wait()
			}

		}
		s.Wait()
	} else {
		fmt.Print("Please add accounts in accounts.txt\n")
		os.Exit(1)
	}
}
func main() {
	printInfo()
	initializeUserData()
}
