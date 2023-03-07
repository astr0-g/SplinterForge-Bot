package RequestFunc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/levigross/grequests"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"splinterforge/SpStruct"
	"strings"
	"time"
)

func FetchselectHero(publicAPIEndpoint string, bossName string, splinterforgeAPIEndpoint string) (string, error) {
	bossID := FetchBossID(bossName, splinterforgeAPIEndpoint)

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
func FetchBossID(bossName string, splinterforgeAPIEndpoint string) string {
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
			return bossData["id"].(string)
		}
	}

	return ""
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

	bossId := FetchBossID(bossName, splinterlandAPIEndpoint)

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
