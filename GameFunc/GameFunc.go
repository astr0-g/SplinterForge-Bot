package GameFunc

import (
	"encoding/json"
	"fmt"
	"splinterforge/ColorPrint"
	"splinterforge/SpStruct"
	"strconv"
	"time"

	"github.com/tebeka/selenium"

	"splinterforge/DriverAction"
	"splinterforge/ReadFunc"
	"splinterforge/RequestFunc"
)

func CheckPopUp(wd selenium.WebDriver, millisecond int) {
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

func SelectCards(cardSelection []SpStruct.CardSelection, bossName string, userName string, splinterlandAPIEndpoint string, publicAPIEndpoint string, autoSelectCard bool) ([]SpStruct.CardSelection, bool, error) {
	if autoSelectCard {
		ColorPrint.PrintWhite(
			userName, fmt.Sprintf("Auto selecting playing cards for desire boss: %s", bossName))
		playingDeck, err := RequestFunc.FetchBattleCards(bossName, userName, splinterlandAPIEndpoint, publicAPIEndpoint)
		if err != nil {
			ColorPrint.PrintPurple(userName, "API Error, Fetch Battle cards failed.")
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

			playingSummonersList := make([]SpStruct.Summoners, 0, len(summonerIds))
			playingMonsterList := make([]SpStruct.MonsterId, 0, len(monsterIds))
			for _, i := range summonerIds {
				summonerId := fmt.Sprintf("%v", i)
				cardName, _ := ReadFunc.GetCardName(summonerId)
				playingSummonersList = append(playingSummonersList, SpStruct.Summoners{
					PlayingSummonersID:   summonerId,
					PlayingSummonersName: cardName,
					PlayingSummonersDiv:  fmt.Sprintf("//div/img[@id='%s']", summonerId),
				})
			}

			for _, i := range monsterIds {
				monsterId := fmt.Sprintf("%v", i)
				cardName, _ := ReadFunc.GetCardName(monsterId)
				playingMonsterList = append(playingMonsterList, SpStruct.MonsterId{
					PlayingMonstersID:   monsterId,
					PlayingMonstersName: cardName,
					PlayingMontersDiv:   fmt.Sprintf("//div/img[@id='%s']", monsterId),
				})
			}

			var cardSelectionList = []SpStruct.CardSelection{}
			cardSelection := SpStruct.CardSelection{
				PlayingSummoners: playingSummonersList,
				PlayingMonsters:  playingMonsterList,
			}
			cardSelectionList = append(cardSelectionList, cardSelection)
			return cardSelectionList, true, nil
		} else {
			ColorPrint.PrintPurple(userName, "Auto selecting playing cards deck failed, you might have too less cards in the account, will continue play with your card setting.")
			return cardSelection, false, nil
		}
	} else {
		return cardSelection, false, nil
	}
}

func SelectHero(heroesType string, userName string, wd selenium.WebDriver, auto_select_hero bool, publicAPIEndpoint string, bossName string, splinterforgeAPIEndpoint string) {
	heroTypes := [3]string{"Warrior", "Wizard", "Ranger"}
	if auto_select_hero {
		hero_type, err := RequestFunc.FetchselectHero(publicAPIEndpoint, bossName, splinterforgeAPIEndpoint)
		if err == nil {
			ColorPrint.PrintWhite(userName, fmt.Sprintf("Auto selecting heroes type: %s for desired boss: %s", hero_type, bossName))
			for i, val := range heroTypes {
				if val == hero_type {
					heroesType = strconv.Itoa(i + 1)
					break
				}
			}
		} else {
			ColorPrint.PrintRed(userName, "Auto selecting heroes type failed due to API error.")
		}
	} else {
		ColorPrint.PrintWhite(userName, "Selecting heroes type...")
	}
	defer func() {
		if err := recover(); err != nil {
			ColorPrint.PrintRed(userName, "Error in selecting hero type, continue...")
		}
	}()
	heroIndex, _ := strconv.Atoi(heroesType)
	hero_type := heroTypes[heroIndex-1]
	bossXpath := "/html/body/app/div[1]/slcards/div[5]/section[1]/div/div[1]/div[2]/section/div"
	el, _ := wd.FindElement(selenium.ByXPATH, bossXpath)
	el.Click()
	bossSelectXpath := fmt.Sprintf("%s/ul/li[%s]", bossXpath, heroesType)
	el, _ = wd.FindElement(selenium.ByXPATH, bossSelectXpath)
	el.Click()
	ColorPrint.PrintWhite(userName, fmt.Sprintf("Selected hero type: %s", hero_type))
}

func SelectBoss(userName string, bossIdToSelect string, wd selenium.WebDriver) (string, string, error) {
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
					ColorPrint.PrintRed(userName, "The selected boss has been defeated, selecting another one automatically...")

					if bossIdToSelect < "17" {
						bossIdInt, _ := strconv.Atoi(bossIdToSelect)
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

func SelectSummoners(userName string, seletedNumOfSummoners int, cardDiv string, wd selenium.WebDriver) bool {
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
			err = el.Click()
			if err == nil {
				time.Sleep(1 * time.Second)
				checkCardDiv := fmt.Sprintf("/html/body/app/div[1]/slcards/div[5]/section[1]/div/div[1]/div[2]/div[1]/div[2]/div[%s]", strconv.Itoa(seletedNumOfSummoners))
				success, _ := DriverAction.DriverwaitForElement(wd, checkCardDiv)
				if success {
					result = true
					break
				} else {
					clickedTime++
				}
			} else {
				wd.ExecuteScript("window.scrollBy(0, 450)", nil)
				scroolTime++
				continue
			}
			
		}

	}
	if clickedTime%2 == 1 {
		result = true
	}
	if !result {
		ColorPrint.PrintRed(userName, "Error in selecting summoners, skipped...")
	}
	wd.ExecuteScript("window.scrollBy(0, -4000)", nil)
	time.Sleep(1 * time.Second)
	return result
}

func SelectMonsters(userName string, seletedNumOfMonsters int, cardDiv string, wd selenium.WebDriver) bool {
	scroolTime := 0
	clickedTime := 0
	result := false
	time.Sleep(1 * time.Second)
	for scroolTime < 10 && clickedTime < 5 {
		el, err := wd.FindElement(selenium.ByXPATH, cardDiv)
		if err != nil {
			wd.ExecuteScript("window.scrollBy(0, 450)", nil)
			scroolTime++
			continue
		} else {
			err = el.Click()
			if err == nil {
				checkCardDiv := fmt.Sprintf("/html/body/app/div[1]/slcards/div[5]/section[1]/div/div[1]/div[2]/div[2]/div[2]/div[%s]", strconv.Itoa(seletedNumOfMonsters))
				success, _ := DriverAction.DriverwaitForElement(wd, checkCardDiv)
				if success {
					result = true
					break
				} else {
					clickedTime++
					wd.ExecuteScript("window.scrollBy(0, 450)", nil)
					continue
				}
			} else {
				wd.ExecuteScript("window.scrollBy(0, 450)", nil)
				scroolTime++
				continue
			}
		}

	}
	if clickedTime%2 == 1 {
		result = true
	}
	if !result {
		ColorPrint.PrintRed(userName, "Error in selecting Monsters, skipped...")
	}
	wd.ExecuteScript("window.scrollBy(0, -4000)", nil)
	time.Sleep(1 * time.Second)
	return result
}
