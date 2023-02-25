package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/selenium-Driver-Check/SeleniumDriverCheck"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

func getCardName(cardId string) (string, error) {
	// Open the JSON file containing the card names and IDs.
	file, err := os.Open("data/cardMapping.json")
	if err != nil {
		return "", fmt.Errorf("error opening JSON file: %w", err)
	}
	defer file.Close()

	// Decode the JSON file into a slice of maps.
	var cards []map[string]string
	err = json.NewDecoder(file).Decode(&cards)
	if err != nil {
		return "", fmt.Errorf("error decoding JSON file: %w", err)
	}

	// Find the card ID for the given card name.
	for _, card := range cards {
		if id, ok := card[cardId]; ok {
			return id, nil
		}
	}

	// Return an error if the card name was not found.
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
            continue // skip the first line
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

func elementWaitAndClick(wd selenium.WebDriver, xpath string){
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
func login(userName string, postingKey string, wd selenium.WebDriver,err error) {
	
    err = wd.SetImplicitWaitTimeout(5 * time.Second)
    if err != nil {
        panic(err)
    }
    err = wd.Get("chrome-extension://jcacnejopjdphbnjgfaaobbfafkihpep/popup.html")
    if err != nil {
        panic(err)
    }
    
    elementWaitAndClick(wd,"/html/body/div/div/div[4]/div[2]/div[5]/button")

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
	time.Sleep(1*time.Second)
    el.SendKeys("\ue007")
    err = wd.ResizeWindow("bigger",1565,1080)
    if err != nil{
        println("can not change size")
    }

	
	// wd.SetWindowSize(1565, 1080)
	wd.Get("https://splinterforge.io/#/")
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
    println("success log in")
    fmt.Println(time.Now())
}

func initializeAccount(accountNo int) (string, string, string, string,[]map[string]interface{}, int) {
    userName, postingKey, err := getAccountData("config/accounts.txt", accountNo)
    if err != nil || userName == "" || postingKey == "" {
        fmt.Println("Error in loading accounts.txt, please add username or posting key and try again.")
    }

    heroesType, bossId, playingSummoners, playingMonster, timeSleepInMinute, err := getCardSettingData("config/cardSettings.txt", accountNo)
    if err != nil {
        fmt.Println("Error loading cardSettings.txt file")
    }

    playingSummonersList := make([]map[string]string, 0, len(playingSummoners))
    playingMonsterList := make([]map[string]string, 0, len(playingMonster))

    for _, i := range playingSummoners {
        cardName, _ := getCardName(i)
        playingSummonersList = append(playingSummonersList, map[string]string{
            "playingSummonersDiv":  fmt.Sprintf("//div/img[@id='%s']", i),
            "playingSummonersId":   i,
            "playingSummonersName": cardName,
        })
    }
    for _, i := range playingMonster {
        cardName, _ := getCardName(i)
        playingMonsterList = append(playingMonsterList, map[string]string{
            "playingMontersDiv":   fmt.Sprintf("//div/img[@id='%s']", i),
            "playingMonstersId":   i,
            "playingMonstersName": cardName,
        })
    }

    cardSelection := []map[string]interface{}{
        {
            "playingSummoners": playingSummonersList,
            "playingMonsterId": playingMonsterList,
        },
    }

    timeSleepInMinute *= 60

    return userName, postingKey, heroesType, bossId, cardSelection, timeSleepInMinute
}

func initializeDriver(accountData []map[string]interface{}){
    extensionData, err := ioutil.ReadFile("data/hivekeychain.crx")
    if err != nil {
        println((1))
    }

    extensionBase64 := base64.StdEncoding.EncodeToString(extensionData)
    chromeOptions := chrome.Capabilities{
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
            "--disable-gpu",
            "--disable-blink-features=AutomationControlled",
            "--mute-audio",
            "--ignore-certificate-errors",
            "--allow-running-insecure-content",
            "--window-size=300,600",
            // "--headless=new",
        },
        Extensions: []string{extensionBase64},
        Prefs: map[string]interface{}{
            "profile.managed_default_content_settings.images":          1,
            "profile.managed_default_content_settings.cookies":         1,
            "profile.managed_default_content_settings.javascript":      1,
            "profile.managed_default_content_settings.plugins":         1,
            "profile.default_content_setting_values.notifications":     2,
            "profile.managed_default_content_settings.stylesheets":     2,
            "profile.managed_default_content_settings.popups":          2,
            "profile.managed_default_content_settings.geolocation":     2,
            "profile.managed_default_content_settings.media_stream":    2,
        },
        ExcludeSwitches: []string{
            "enable-automation",
            "enable-logging",
        },
    }

    caps := selenium.Capabilities{}
    caps.AddChrome(chromeOptions)


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

    userName, _ := accountData[0]["userName"].(string)
    postingKey, _ := accountData[0]["postingKey"].(string)
    // bossId, _ := accountData[0]["bossId"].(string)
    // heroesType, _ := accountData[0]["heroesType"].(string)
    // cardSelection, _ := accountData[0]["cardSelection"].([]map[string]interface{})
    login(userName,postingKey,driver,err)
    // battle(accountData,driver,err)
    screenshot, err := driver.Screenshot()
    if err != nil {
        fmt.Printf("Failed to take screenshot: %s\n", err)
        os.Exit(1)
    }

    // write the screenshot to a file
    if err := ioutil.WriteFile("screenshot.png", screenshot, 0644); err != nil {
        fmt.Printf("Failed to write screenshot to file: %s\n", err)
        os.Exit(1)
    }
}
func countLines(filePath string) (int, error) {
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

func initializeUserData() {
    lineCount, errCountLines := countLines("config/accounts.txt")
    var accountLists []map[string]interface{}
    if errCountLines == nil && lineCount > 1 {
        for i := 0; i < lineCount-1; i++ {
            userName, postingKey, heroesType, bossId, cardSelection, timeSleepInMinute := initializeAccount(i+1)
            accountLists = append(accountLists, map[string]interface{}{
                "userName":          userName,
                "postingKey":        postingKey,
                "bossId":            bossId,
                "heroesType":        heroesType,
                "cardSelection":     cardSelection,
                "timeSleepInMinute": timeSleepInMinute,
            })
        }
    } else {
        fmt.Print("Please add accounts in accounts.txt\n")
        os.Exit(1)
    }
    fmt.Println(accountLists)
    initializeDriver([]map[string]interface{}{accountLists[0]})
}
func main() {
    initializeUserData()
}