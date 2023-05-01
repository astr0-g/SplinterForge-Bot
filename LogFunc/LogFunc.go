package LogFunc

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"splinterforge/ColorPrint"
	"splinterforge/SpStruct"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/levigross/grequests"
	"github.com/olekukonko/tablewriter"
)

func PrintInfo() {
	color.Set(color.FgWhite)
	fmt.Println("+-------------------------------------------------+")
	fmt.Println("|          Welcome to SplinterForge Bot!          |")
	fmt.Println("|          Open source Github repository          |")
	fmt.Println("|   https://github.com/Astr0-G/SplinterForge-Bot  |")
	fmt.Println("|                  Discord server                 |")
	fmt.Println("|          https://discord.gg/pm8SGZkYcD          |")
	fmt.Println("+-------------------------------------------------+")
}

func PrintResultBox(userName string, data [][]string, selectResult bool,bossName, bossLeague, heroType string, bossAbilities, bossRandomAbilities []string) {
	sort.Slice(data, func(i, j int) bool {
		return data[i][0] < data[j][0]
	})

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"SUMMARY", "Name/Abilities", "ID/League", "Results"})
	table.Append([]string{"Boss", bossName, bossLeague, ""})
	table.Append([]string{"Abilities", strings.Join(bossAbilities, ", "), "", ""})
	table.Append([]string{"Random Abilities", strings.Join(bossRandomAbilities, ", "),"",  ""})
	table.Append([]string{"Hero Chosen", heroType, "",  ""})
	for _, row := range data {
		table.Append(row)
	}
	if selectResult {
		ColorPrint.PrintGreen(userName, "Card selection results:")
		table.Render()
		color.Set(color.FgWhite)
	} else {
		ColorPrint.PrintYellow(userName, "Card selection results:")
		table.Render()
		color.Set(color.FgWhite)
	}
}

func PrintAccountDetails(userName string, name interface{}, key interface{}, splinterforgeAPIEndpoint string) {
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
	ColorPrint.PrintWhite(userName, fmt.Sprintf("Account balance %s Forge, current stamina %s / %s.", strconv.FormatFloat(powerRes.Sc.Balance, 'f', 2, 64), strconv.Itoa(powerRes.Stamina.Current), strconv.Itoa(powerRes.Stamina.Max)))
}

func PrintConfigSettings(totalAccounts int, headless bool, threadingLimit int, showForgeReward bool, showAccountDetails bool, waitForBossRespawn bool, shareBattleLog bool, autoSelectCard bool, autoSelectHero bool, autoSelectSleepTime bool, unwantedAbilities []string) {
	data := [][]string{
		{"TOTAL_ACCOUNTS_LOADED", fmt.Sprint(totalAccounts)},
		{"HEADLESS", fmt.Sprint(headless)},
		{"THREADING", fmt.Sprint(threadingLimit)},
		{"SHOW_FORGE_REWARD", fmt.Sprint(showForgeReward)},
		{"SHOW_ACCOUNT_DETAILS", fmt.Sprint(showAccountDetails)},
		{"SHARE_BATTLE_LOG", fmt.Sprint(shareBattleLog)},
		{"UNWANTED_ABILITIES", fmt.Sprint(unwantedAbilities)},
		{"AUTO_WAIT_FOR_BOSS_RESPAWN", fmt.Sprint(waitForBossRespawn)},
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
