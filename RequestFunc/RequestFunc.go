package RequestFunc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"splinterforge/SpStruct"
	"strconv"
	"strings"
	"time"

	"github.com/levigross/grequests"
)

const (
	ApiUrl = "https://post.splinterforge.xyz/api/data/sp/save"
)

func FetchselectHero(randomAbilities []string, userName string, userKey string, publicAPIEndpoint string, bossName string, splinterforgeAPIEndpoint string) (string, error) {
	bossID, _, abilities, _ := FetchBossID(bossName, splinterforgeAPIEndpoint)
	fullAbilities := append(abilities, randomAbilities...)
	requestBody := map[string]interface{}{
		"bossId": bossID,
		"myList": fullAbilities,
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
	if heroType, ok := responseData["heroType"].(string); ok {
		if responseData["RecommendHero"].(bool) == false {
			return "", errors.New("RecommendHero is false")
		}
		heroTypeToChoose := strings.Split(heroType, " ")[0]
		return heroTypeToChoose, nil
	}

	return "", fmt.Errorf("heroType not found in response data")
}

func FetchBossAbilities(userName string, userKey string, bossName string, splinterforgeAPIEndpoint string) (bossLeague string, abilities []string, randomAbilities []string) {
	_, bossLeague, abilities, _ = FetchBossID(bossName, splinterforgeAPIEndpoint)
	FetchRandomAbilites(userName, bossLeague, splinterforgeAPIEndpoint)
	randomAbilities, _ = FetchRandomAbilitiesForUsername(userName, userKey, bossLeague, splinterforgeAPIEndpoint)
	if bossLeague == "t1" {
		bossLeague = "Bronze"
	} else if bossLeague == "t2" {
		bossLeague = "Silver"
	} else if bossLeague == "t3" {
		bossLeague = "Gold"
	} else if bossLeague == "t4" {
		bossLeague = "Diamond"
	}
	return bossLeague, abilities, randomAbilities
}

func FetchRandomAbilites(userName string, bossLeague string, splinterforgeAPIEndpoint string) (string, error) {
	url := fmt.Sprintf("%s/users/getRules", splinterforgeAPIEndpoint)

	payload := map[string]string{
		"user": userName,
		"type": bossLeague,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func FetchRandomAbilitiesForUsername(name string, key string, bossLeague string, splinterforgeAPIEndpoint string) ([]string, error) {
	res, err := grequests.Post(fmt.Sprintf("%s/users/keyLogin", splinterforgeAPIEndpoint), &grequests.RequestOptions{
		JSON: map[string]string{
			"username": name,
			"key":      key,
		},
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	})
	if err != nil {
		return nil, err
	}

	var powerRes = SpStruct.KeyLoginResData{}
	err = json.Unmarshal(res.Bytes(), &powerRes)
	if err != nil {
		return nil, err
	}
	switch bossLeague {
	case "t1":
		rules := powerRes.UniqueRules.T1.Rules
		return rules, nil
	case "t2":
		rules := powerRes.UniqueRules.T2.Rules
		return rules, nil
	case "t3":
		rules := powerRes.UniqueRules.T3.Rules
		return rules, nil
	case "t4":
		rules := powerRes.UniqueRules.T4.Rules
		return rules, nil
	default:
		return nil, fmt.Errorf("bossLeague not found: %s", bossLeague)
	}

}

func FetchBossID(bossName string, splinterforgeAPIEndpoint string) (string, string, []string, error) {
	url := fmt.Sprintf("%s/boss/getBosses", splinterforgeAPIEndpoint)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var responseData []map[string]interface{}
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&responseData)

	for _, bossData := range responseData {
		if strings.EqualFold(bossData["name"].(string), bossName) {
			abilities := []string{}
			for _, ability := range bossData["ogStats"].(map[string]interface{})["abilities"].([]interface{}) {
				abilities = append(abilities, ability.(string))
			}
			return bossData["id"].(string), bossData["type"].(string), abilities, nil
		}
	}

	return "", "", nil, errors.New("boss not found")
}

func FetchPlayerCard(userName string, splinterlandAPIEndpoint string) ([]int, error) {
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
	var cardList SpStruct.CardList
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

func FetchBattleCards(bossName string, userName string, splinterlandAPIEndpoint string, publicAPIEndpoint string) (string, error) {
	cardsId, err := FetchPlayerCard(userName, splinterlandAPIEndpoint)
	if err != nil {
		return "", err
	}

	bossId, _, _, _ := FetchBossID(bossName, splinterlandAPIEndpoint)

	postData := SpStruct.BattleCardsRequestBody{
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
		var fitResponseData = SpStruct.GetResponseBody{}
		json.Unmarshal(res.Bytes(), &fitResponseData)
		time.Sleep(5 * time.Second)
		return res.String(), nil

	} else {
		fmt.Println("GetReponseBody error > ", err)
		return "", err
	}
}

func ShareLogToApi(shareCDP SpStruct.ShareCDPFitReturnData, BossName string, BossAbilities []string, BossRandomAbilities []string, BossLeague string, resHeroName string, TotalDmg int) {
	shareCDP.AdditionInfo.Boss = BossName
	shareCDP.AdditionInfo.ShareToDiscord = "yes"
	shareCDP.AdditionInfo.Abilities = BossAbilities
	shareCDP.AdditionInfo.RandomAbilities = BossRandomAbilities
	strBossleague, _ := strconv.Atoi(BossLeague)
	shareCDP.AdditionInfo.BoosType = strconv.Itoa(strBossleague * 3)
	shareCDP.AdditionInfo.HeroChosen = resHeroName
	shareCDP.AdditionInfo.TotalDamage = strconv.Itoa(TotalDmg)
	grequests.Post(ApiUrl, &grequests.RequestOptions{
		JSON: shareCDP,
	})
}
