import asyncio
import psutil
import requests
from selenium import webdriver
import multiprocessing
import linecache
from selenium.webdriver.common.keys import Keys
import chromedriver_autoinstaller
from selenium.webdriver.chrome.options import Options
from selenium.common.exceptions import WebDriverException
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from pyfiglet import Figlet
import datetime
import json
import time
import os
import sys
import ctypes
import textwrap
from prettytable import PrettyTable
from tabulate import tabulate
from colorama import init

init(convert=True)
STD_OUTPUT_HANDLE = -11
std_out_handle = ctypes.windll.kernel32.GetStdHandle(STD_OUTPUT_HANDLE)
ctypes.windll.kernel32.SetConsoleTitleW("SplinterForge Bot")
sys.stdout = open(sys.stdout.fileno(), mode='w', encoding='utf8', buffering=1)


def set_cmd_text_color(color, handle=std_out_handle):
    Bool = ctypes.windll.kernel32.SetConsoleTextAttribute(handle, color)
    return Bool


FOREGROUND_YELLOW = 0x0e  # yellow.
FOREGROUND_GREEN = 0x02  # green.
FOREGROUND_RED = 0x04  # red.
FOREGROUND_DARKRED = 0x04  # dark red.
FOREGROUND_SKYBLUE = 0x0b  # skyblue.
FOREGROUND_Pink = 0x0d  # dark gray.
FOREGROUND_BLUE = 0x09  # blue.
FOREGROUND_DARKBLUE = 0x01  # dark blue.


def printDarkBlue(message):
    set_cmd_text_color(FOREGROUND_DARKBLUE)
    sys.stdout.write(message)
    resetColor()


def printGreen(message):
    set_cmd_text_color(FOREGROUND_GREEN)
    sys.stdout.write(message)
    resetColor()


def printSkyBlue(message):
    set_cmd_text_color(FOREGROUND_SKYBLUE)
    sys.stdout.write(message)
    resetColor()


def printRed(message):
    set_cmd_text_color(FOREGROUND_DARKRED)
    sys.stdout.write(message)
    resetColor()


def printYellow(message):
    set_cmd_text_color(FOREGROUND_YELLOW)
    sys.stdout.write(message)
    resetColor()


def printPinky(message):
    set_cmd_text_color(FOREGROUND_Pink)
    sys.stdout.write(message)
    resetColor()


def resetColor():
    set_cmd_text_color(FOREGROUND_RED | FOREGROUND_GREEN | FOREGROUND_BLUE)


class log_info():
    @staticmethod
    def time():
        return str(datetime.datetime.now().replace(microsecond=0))

    @staticmethod
    def normal(message):
        print(f"{message}")

    @staticmethod
    def success(userName, message):
        printGreen(f"[{log_info.time()}] {userName}: {message}\n")

    @staticmethod
    def error(userName, message):
        printRed(f"[{log_info.time()}] {userName}: {message}\n")

    @staticmethod
    def alerts(userName, message):
        print(f"[{log_info.time()}] {userName}: {message}")

    @staticmethod
    def status(userName, message):
        printYellow(f"[{log_info.time()}] {userName}: {message}\n")

    @staticmethod
    def verify(userName, message):
        printDarkBlue(f"[{log_info.time()}] {userName}: {message}\n")


def printResultBox(userName, data, selectResult):
    table = PrettyTable(["Card", "ID", "Name", "Selection Results"])
    for row in data:
        table.add_row(row)
    dataToPrint = table.get_string(sortby="Card")
    if selectResult:
        log_info.success(userName, "Card selection results:")
        printGreen(f"{dataToPrint}\n")
    else:
        log_info.status(userName, "Card selection results:")
        printYellow(f"{dataToPrint}\n")


def printConfigSettings(totallaccounts, headless, close_driver_while_sleeping, start_thread, start_thread_interval, show_forge_reward, show_total_forge_balance, print_system_usage, check_system_usage_frequency, auto_select_card, auto_select_hero):
    data = [['TOTAL_ACCOUNTS_LOADED', totallaccounts],
            ['HEADLESS', headless],
            ['CLOSE_DRIVER_WHILE_SLEEPING', close_driver_while_sleeping],
            ['START_THREAD', start_thread],
            ['START_THREAD_INTERVAL(seconds)', start_thread_interval],
            ['SHOW_FORGE_REWARD', show_forge_reward],
            ['SHOW_TOTAL_FORGE_BALANCE', show_total_forge_balance],
            ['PRINT_SYSTEM_USAGE', print_system_usage],
            ['CHECK_SYSTEM_USAGE_FREQUENCY(seconds)',
             check_system_usage_frequency],
            ['AUTO_SELECT_CARD', auto_select_card],
            ['AUTO_SELECT_HERO', auto_select_hero]
            ]
    print(tabulate(data, headers=['Setting', 'Value'], tablefmt='grid'))


def start_font():
    f = Figlet(font='smkeyboard', width=150)
    text = f.renderText('SplinterForge\n\/\Bot-Beta\/\n\@lil_astr_0/')
    return f"{text}"


def printWelcome(text):
    console_width = 49
    console_height = 0
    lines = text.split('\n')
    wrapped_lines = []
    for line in lines:
        wrapped_lines.extend(textwrap.wrap(line, console_width))
    total_lines = len(wrapped_lines)
    top_pad = (console_height - total_lines) // 2
    bottom_pad = console_height - total_lines - top_pad
    border_width = console_width
    border_char = '-'
    corner_char = '+'
    border_line = corner_char + border_char * border_width + corner_char
    print(border_line)
    for i in range(top_pad):
        print(corner_char + ' ' * border_width + corner_char)
    for line in wrapped_lines:
        print(corner_char + line.center(console_width) + corner_char)
    for i in range(bottom_pad):
        print(corner_char + ' ' * border_width + corner_char)
    print(border_line)


def file_len(file_path):
    with open(file_path, 'r') as rf:
        record = 0
        for line in rf:
            if line.strip():  # Check if the line is not empty after stripping whitespaces
                record += 1
        return record


def getAccountData(file_path, line_number):
    acctinfo = linecache.getline(file_path, int(line_number)+1).strip()
    time.sleep(1)
    userName = acctinfo.split(":")[0]
    postingKey = acctinfo.split(":")[1]
    return userName, postingKey


def getCardSettingData(file_path, line_number):
    acctinfo = linecache.getline(file_path, int(line_number)+1).strip()
    heroesType = acctinfo.split(":")[0]
    bossId = acctinfo.split(":")[1]
    playingSummoners = acctinfo.split(":")[2].split(',')
    playingMonster = acctinfo.split(":")[3].split(',')
    timeSleepInMinute = acctinfo.split(":")[4]
    return heroesType, bossId, playingSummoners, playingMonster, timeSleepInMinute


def getConfig(file_path):
    with open(file_path, 'r') as file:
        lines = file.readlines()
        headless = None
        close_driver_while_sleeping = None
        start_thread = None
        start_thread_interval = None
        show_forge_reward = None
        show_total_forge_balance = None
        print_system_usage = None
        check_system_usage_frequency = None
        auto_select_card = None
        auto_select_hero = None
        splinterland_api_endpoint = None
        public_api_endpoint = None

        for line in lines:
            line = line.strip()
            if '=' in line:
                key, value = line.split('=')
                key = key.strip()
                value = value.strip()
                if value.lower() == 'true':
                    value = True
                elif value.lower() == 'false':
                    value = False
                else:
                    try:
                        value = int(value)
                    except ValueError:
                        value = value

                if key == 'HEADLESS':
                    headless = value
                elif key == 'CLOSE_DRIVER_WHILE_SLEEPING':
                    close_driver_while_sleeping = value
                elif key == 'START_THREAD':
                    start_thread = value
                elif key == 'START_THREAD_INTERVAL':
                    start_thread_interval = value
                elif key == 'SHOW_FORGE_REWARD':
                    show_forge_reward = value
                elif key == 'SHOW_TOTAL_FORGE_BALANCE':
                    show_total_forge_balance = value
                elif key == 'PRINT_SYSTEM_USAGE':
                    print_system_usage = value
                elif key == 'CHECK_SYSTEM_USAGE_FREQUENCY':
                    check_system_usage_frequency = value
                elif key == 'AUTO_SELECT_CARD':
                    auto_select_card = value
                elif key == 'AUTO_SELECT_HERO':
                    auto_select_hero = value
                elif key == 'SPLINTERLAND_API_ENDPOINT':
                    splinterland_api_endpoint = value
                elif key == 'PUBLIC_API_ENDPOINT':
                    public_api_endpoint = value

        return (headless,
                close_driver_while_sleeping,
                start_thread,
                start_thread_interval,
                show_forge_reward,
                show_total_forge_balance,
                print_system_usage,
                check_system_usage_frequency,
                auto_select_card,
                auto_select_hero,
                splinterland_api_endpoint,
                public_api_endpoint)


def getCardData(userName, cardID):
    with open('data/cardsDetails.json') as json_file:
        data = json.load(json_file)
    found = False
    for i in range(len(data)):
        if int(cardID) == int(data[i]['id']):
            found = True
            name = data[i]['name']
            break
    if found:
        return name
    else:
        log_info.error(userName, "Card with ID {} not found".format(cardID))


def _init(accountNo):
    try:
        userName, postingKey = getAccountData("config/accounts.txt", accountNo)
        if userName == "" or postingKey == "":
            log_info.error(userName,
                           "Error in loading accounts.txt, please add username or posting key and try again.")
            log_info.error(userName,
                           "Terminating in 10 seconds...")
            time.sleep(10)
            sys.exit()
    except:
        print("error loading accounts.txt, retrying...")
        time.sleep(10)
        sys.exit()
    try:
        cardSelection = []
        playingSummonersList = []
        playingMonsterList = []
        heroesType, bossId, playingSummoners, playingMonster, timeSleepInMinute = getCardSettingData(
            "config/cardSettings.txt", accountNo)
        for i in playingSummoners:
            cardName = getCardData(userName, i)
            playingSummonersList.append({
                "playingSummonersDiv": f"//div/img[@id='{i}']",
                "playingSummonersId": f"{i}",
                "playingSummonersName": f"{cardName}"
            })
        for i in playingMonster:
            cardName = getCardData(userName, i)
            playingMonsterList.append({
                "playingMontersDiv": f"//div/img[@id='{i}']",
                "playingMonstersId": f"{i}",
                "playingMonstersName": f"{cardName}"
            })
        cardSelection.append({
            "bossId": f"{bossId}",
            "playingSummoners": playingSummonersList,
            "playingMonsterId": playingMonsterList
        })
        timeSleepInMinute = int(timeSleepInMinute) * 60
    except:
        print("error loading cardSettings.txt file")
        print("Terminating in 10 seconds...")
        time.sleep(10)
        sys.exit()
    return userName, postingKey, heroesType, cardSelection, timeSleepInMinute


def collectData(body, public_api_endpoint):
    url = f"{public_api_endpoint}/data/"
    headers = {
        'Content-Type': 'application/json'
    }
    requests.post(url, headers=headers, json=body)


def autoSelectCard(cardSelection, bossName, userName, splinterland_api_endpoint, public_api_endpoint):
    # log_info.alerts(
    #     userName, f"Auto selecting playing cards for desire boss: {bossName}")
    try:
        playingDeck = fetchBattleCards(
            bossName, userName, splinterland_api_endpoint, public_api_endpoint)
        if playingDeck["RecommendTeam"]:
            cardSelection = []
            playingSummonersList = []
            playingMonsterList = []
            for i in playingDeck["summonerIds"]:
                cardName = getCardData(userName, i)
                playingSummonersList.append({
                    "playingSummonersDiv": f"//div/img[@id='{i}']",
                    "playingSummonersId": f"{i}",
                    "playingSummonersName": f"{cardName}"
                })
            for i in playingDeck["monsterIds"]:
                cardName = getCardData(userName, i)
                playingMonsterList.append({
                    "playingMontersDiv": f"//div/img[@id='{i}']",
                    "playingMonstersId": f"{i}",
                    "playingMonstersName": f"{cardName}"
                })
            cardSelection.append({
                "playingSummoners": playingSummonersList,
                "playingMonsterId": playingMonsterList
            })
            log_info.success(
                userName, f"Auto selecting playing cards successful!")
            return cardSelection, True
        else:
            log_info.status(
                userName, f"Auto selecting playing cards failed, you might have too less cards in the account, will continue play with your card setting.")
            return cardSelection, False
    except:
        log_info.alerts(
            userName, f"Auto selecting playing cards failed, API error.")
        return cardSelection, False


def selectSummoners(userName, seletedNumOfSummoners, cardDiv, driver):
    multiprocessing.freeze_support()
    scroolTime = 0
    clickedTime = 0
    result = False
    time.sleep(1)
    while scroolTime < 5 and clickedTime < 5:
        try:
            WebDriverWait(driver, 0.75).until(
                EC.presence_of_element_located((By.XPATH, cardDiv)))
            WebDriverWait(driver, 1).until(
                EC.element_to_be_clickable((By.XPATH, cardDiv))).click()
            time.sleep(1)
            clickedTime += 1

            selectNumber = driver.find_element(
                By.XPATH, "/html/body/app/div[1]/slcards/div[5]/section[1]/div/div[1]/div[2]/div[1]/div[1]/h3/span[2]").text
            if seletedNumOfSummoners == int(selectNumber.split("/")[0]):
                result = True
                break
        except:
            driver.execute_script("window.scrollBy(0, 180)")
            scroolTime += 1
    if clickedTime % 2:
        result = True
    if not result:
        log_info.error(userName, "Error in selecting summoners, retrying...")
    driver.execute_script("window.scrollBy(0, -4000)")
    time.sleep(1)
    return result


def selectMonsterCards(userName, seletedNumOfMonsters, cardId, cardDiv, driver):
    multiprocessing.freeze_support()
    scroolTime = 0
    clickedTime = 0
    result = False
    time.sleep(1)
    while scroolTime < 15 and clickedTime < 15:
        try:
            WebDriverWait(driver, 0.75).until(
                EC.presence_of_element_located((By.XPATH, cardDiv)))
            WebDriverWait(driver, 1).until(
                EC.element_to_be_clickable((By.XPATH, cardDiv))).click()
            clickedTime += 1
            time.sleep(1)
            selectNumber = driver.find_element(
                By.XPATH, "/html/body/app/div[1]/slcards/div[5]/div[1]/h3/span[2]/div[1]/button/span").text
            selectNumber = int(selectNumber)
            if seletedNumOfMonsters == selectNumber:
                result = True
                # log_info.success(
                #     userName, f"Monster card ID {cardId} selected successful!")
                break
        except:
            driver.execute_script("window.scrollBy(0, 450)")
            scroolTime += 1
            pass
        # log_info.error(userName,
        #                f"Error select card ID: {cardId}, skipped this card...")
    if clickedTime % 2:
        result = True
    driver.execute_script("window.scrollBy(0, -10000)")
    time.sleep(1)
    return result


def check(driver):
    try:
        WebDriverWait(driver, 1.2).until(
            EC.element_to_be_clickable((By.XPATH, "/html/body/app/div[1]/div[1]/app-header/success-modal/section/div[1]/div[4]/div/button"))).click()
    except:
        pass
    try:
        WebDriverWait(driver, 1.2).until(
            EC.element_to_be_clickable((By.XPATH, "/html/body/app/div[1]/login-modal/div/div/div/div[2]/div[3]/button"))).click()
    except:
        pass


def fetchBattleCards(bossName, userName, splinterland_api_endpoint, public_api_endpoint):
    cardsId = fetchPlayerCard(userName, splinterland_api_endpoint)
    bossId = fetchBossId(bossName)
    url = f"{public_api_endpoint}/teamselection/"
    headers = {
        'Content-Type': 'application/json'
    }
    body = {
        "bossName": f"{bossName}", "bossId": f"{bossId}", "team": cardsId}
    response = requests.post(url, headers=headers, json=body).json()
    return response


def fetchPlayerCard(userName, splinterland_api_endpoint):
    url = f"{splinterland_api_endpoint}/cards/collection/{userName}"
    headers = {
        'Content-Type': 'application/json'
    }
    response = requests.get(url, headers=headers)
    states = response.json()
    cardsIds = []
    for i in states['cards']:
        if i["card_detail_id"] not in cardsIds:
            cardsIds.append(i["card_detail_id"])
    return cardsIds


def fetchBossId(bossName):
    url = "https://splinterforge.io/boss/getBosses"
    headers = {
        'Content-Type': 'application/json'
    }
    response = requests.get(url, headers=headers)
    states = response.json()
    for i in states:
        if i['name'].lower() == bossName.lower():
            return i['id']


def fetchHeroSelect(public_api_endpoint, bossName):
    bossId = fetchBossId(bossName)
    headers = {
        'Content-Type': 'application/json'
    }
    body = {
        "bossId": f"{bossId}",
    }

    heroTypeToChoose = requests.post(
        f"{public_api_endpoint}/heroselection", headers=headers, json=body).json()['heroTypes'].split(" ")[0]
    return(str(heroTypeToChoose))


def heroSelect(heroesType, userName, driver,
               auto_select_hero, public_api_endpoint, bossName):
    check(driver)
    hero_types = ['Warrior', 'Wizard', 'Ranger']
    if auto_select_hero:
        try:
            hero_type = fetchHeroSelect(public_api_endpoint, bossName)
            log_info.alerts(userName,
                            f"Auto selecting heros type: {hero_type} for desire boss: {bossName}")
            heroesType = hero_types.index(hero_type)+1
        except:
            log_info.error(userName,
                           f"Auto selecting heros type failed due to API error.")
    else:
        log_info.alerts(userName,
                        f"Selecting heros type...")

    try:
        hero_type = hero_types[int(heroesType) - 1]
        xpath = "/html/body/app/div[1]/splinterforge-heros/div[3]/section/div/div/div[2]/div[1]"

        WebDriverWait(driver, 10).until(EC.element_to_be_clickable(
            (By.XPATH, "/html/body/app/div[1]/div[1]/app-header/section/div[4]/div[2]/div[1]/a[4]"))).click()
        WebDriverWait(driver, 10).until(
            EC.element_to_be_clickable((By.XPATH, xpath))).click()
        WebDriverWait(driver, 10).until(EC.element_to_be_clickable(
            (By.XPATH, f"{xpath}/ul/li[{heroesType}]"))).click()
        log_info.success(userName, f"Selected hero type: {hero_type}.")
    except:
        log_info.error(userName, f"Error in selecting hero type, continue...")
        pass
    time.sleep(2)


def bossSelect(bossIdToSelect, userName, driver):
    WebDriverWait(driver, 10).until(
        EC.element_to_be_clickable((By.XPATH, "/html/body/app/div[1]/div[1]/app-header/section/div[4]/div[2]/div[1]/a[5]/div[1]"))).click()
    while True:
        WebDriverWait(driver, 5).until(
            EC.element_to_be_clickable((By.XPATH, f"//div[@tabindex='{bossIdToSelect}']"))).click()
        if "BOSS IS DEAD" != WebDriverWait(driver, 5).until(
                EC.presence_of_element_located((By.XPATH, "/html/body/app/div[1]/slcards/div[5]/section[1]/div/div[1]/div[2]/button"))).text:
            time.sleep(1)
            bossName = WebDriverWait(driver, 5).until(
                EC.presence_of_element_located((By.XPATH, "/html/body/app/div[1]/slcards/div[5]/section[1]/div/div[1]/div[1]/div[2]/div[3]/h3"))).text
            return bossName
        else:
            log_info.error(userName,
                           "The selected boss has been defeated, selecting another one automatically...")
            if int(bossIdToSelect) < 17:
                bossIdToSelect = str(int(bossIdToSelect) + 1)
            else:
                bossIdToSelect = "14"
            WebDriverWait(driver, 10).until(
                EC.element_to_be_clickable((By.XPATH, "/html/body/app/div[1]/div[1]/app-header/section/div[4]/div[2]/div[1]/a[5]/div[1]"))).click()


def login(userName, postingKey, driver):
    driver.get(
        "chrome-extension://jcacnejopjdphbnjgfaaobbfafkihpep/popup.html")
    try:
        WebDriverWait(driver, 5).until(
            EC.element_to_be_clickable((By.XPATH, "/html/body/div/div/div[4]/div[2]/div[5]/button"))).click()
        WebDriverWait(driver, 5).until(
            EC.element_to_be_clickable((By.XPATH, "/html/body/div/div/div[1]/div/div[1]/div/input"))).send_keys("Aa123Aa123!!")
        WebDriverWait(driver, 5).until(
            EC.element_to_be_clickable((By.XPATH, "/html/body/div/div/div[1]/div/div[2]/div/input"))).send_keys("Aa123Aa123!!")
        WebDriverWait(driver, 5).until(
            EC.element_to_be_clickable((By.XPATH, "/html/body/div/div/div[1]/button/div"))).click()
        WebDriverWait(driver, 5).until(
            EC.element_to_be_clickable((By.XPATH, "/html/body/div/div/div[1]/div[2]/div/div[2]/button[1]/div"))).click()
        WebDriverWait(driver, 5).until(
            EC.element_to_be_clickable((By.XPATH, "/html/body/div/div/div[1]/div[2]/div/div[2]/div[1]/div/input"))).send_keys(userName)
        WebDriverWait(driver, 5).until(
            EC.element_to_be_clickable((By.XPATH, "/html/body/div/div/div[1]/div[2]/div/div[2]/div[2]/div/input"))).send_keys(postingKey)
        time.sleep(3)
        WebDriverWait(driver, 5).until(
            EC.element_to_be_clickable((By.XPATH, "/html/body/div/div/div[1]/div[2]/div/div[2]/div[2]/div/input"))).send_keys(Keys.ENTER)
        # if WebDriverWait(driver, 5).until(
        #         EC.presence_of_element_located((By.XPATH, "/html/body/div/div/div[4]/div"))).text == "HIVE KEYCHAIN":
        # log_info.success(userName,
        #                     "Login successful!")
    except:
        log_info.error(userName,
                       "Login failure! Please check your username and posting key in the accounts.txt file and try again.")
        time.sleep(15)
        driver.close()
    driver.set_window_size(1565, 1080)
    driver.get("https://splinterforge.io/#/")
    WebDriverWait(driver, 5).until(
        EC.element_to_be_clickable((By.XPATH, "/html/body/app/div[1]/div[1]/app-header/success-modal/section/div[1]/div[4]/div/button"))).click()
    WebDriverWait(driver, 5).until(
        EC.element_to_be_clickable((By.XPATH, "/html/body/app/div[1]/div[1]/app-header/section/div[4]/div[2]/div/div/a/div[1]"))).click()
    WebDriverWait(driver, 5).until(
        EC.element_to_be_clickable((By.XPATH, "/html/body/app/div[1]/login-modal/div/div/div/div[2]/div[2]/input"))).send_keys(userName)
    time.sleep(1)
    WebDriverWait(driver, 5).until(
        EC.element_to_be_clickable((By.XPATH, "/html/body/app/div[1]/login-modal/div/div/div/div[2]/div[3]/button"))).click()
    while True:
        try:
            driver.switch_to.window(driver.window_handles[1])
            WebDriverWait(driver, 5).until(
                EC.element_to_be_clickable((By.XPATH, "/html/body/div/div/div/div[3]/div[1]/div/div"))).click()
            WebDriverWait(driver, 5).until(
                EC.element_to_be_clickable((By.XPATH, "/html/body/div/div/div/div[3]/div[2]/button[2]/div"))).click()
            driver.switch_to.window(driver.window_handles[0])
            log_info.success(userName,
                             "Login successful!")
            break
        except:
            pass


def battle(cardSelection, userName, heroesType, driver, show_forge_reward, show_total_forge_balance, auto_select_card, auto_select_hero, splinterland_api_endpoint, public_api_endpoint):
    try:
        check(driver)
        bossIdToSelect = cardSelection[0]['bossId']
        bossName = bossSelect(bossIdToSelect, userName, driver)
        heroSelect(heroesType, userName, driver,
                   auto_select_hero, public_api_endpoint, bossName)
        bossSelect(bossIdToSelect, userName, driver)
        log_info.success(userName,
                         f"Participating in battles...")
        WebDriverWait(driver, 5).until(
            EC.element_to_be_clickable((By.XPATH, "/html/body/app/div[1]/slcards/div[5]/section[1]/div/div[1]/div[2]/button"))).click()
        if auto_select_card:
            cardSelection, autoSelectResult = autoSelectCard(
                cardSelection, bossName, userName, splinterland_api_endpoint, public_api_endpoint)
        log_info.alerts(
            userName, "Selecting cards for summoners and monsters, this process could be lengthy...")
        printData = []
        checkDataSummoners = []
        checkDataMonsters = []
        seletedNumOfSummoners = 1
        for i in cardSelection[0]['playingSummoners']:
            summonersInfo = i
            cardId = summonersInfo["playingSummonersId"]
            cardDiv = summonersInfo["playingSummonersDiv"]
            cardName = summonersInfo["playingSummonersName"]
            resultSeletetSummoners = selectSummoners(
                userName, seletedNumOfSummoners, cardDiv, driver)
            seletedNumOfSummoners += 1
            if resultSeletetSummoners:
                printData.append(
                    [f"Summoners #{int(cardSelection[0]['playingSummoners'].index(i)) + 1}", cardId, cardName, "success"])
                checkDataSummoners.append({
                    "_id": f"{seletedNumOfSummoners - 1}",
                    "card_id": f"{cardId}",
                    "card_name": f"{cardName}"
                })
            else:
                log_info("restarting...")
        seletedNumOfMonsters = 1
        selectResult = True
        for i in cardSelection[0]['playingMonsterId']:
            monstersInfo = i
            cardId = monstersInfo["playingMonstersId"]
            cardDiv = monstersInfo["playingMontersDiv"]
            cardName = monstersInfo["playingMonstersName"]
            resultSeletetMonsters = selectMonsterCards(
                userName, seletedNumOfMonsters, cardId, cardDiv, driver)
            if resultSeletetMonsters:
                printData.append(
                    [f"Monsters #{int(cardSelection[0]['playingMonsterId'].index(monstersInfo))+1}", cardId, cardName, "success"])
                checkDataMonsters.append({
                    "_id": f"{seletedNumOfMonsters}",
                    "card_id": f"{cardId}",
                    "card_name": f"{cardName}"
                })
                seletedNumOfMonsters += 1
            else:
                selectResult = False
                printData.append(
                    [f"Monsters #{int(cardSelection[0]['playingMonsterId'].index(monstersInfo))+1}", cardId, cardName, "error"])
        printResultBox(userName, printData, selectResult)
        manaUsed = WebDriverWait(driver, 5).until(
            EC.presence_of_element_located((By.XPATH, "/html/body/app/div[1]/slcards/div[5]/div[2]/div[1]/div[1]/button/span"))).text
        totalManaHave = WebDriverWait(driver, 5).until(
            EC.presence_of_element_located((By.XPATH, "/html/body/app/div[1]/div[1]/app-header/section/div[4]/div[2]/div[2]/div[1]/a[2]/div[1]/span"))).text
        if int(manaUsed) < 15:
            log_info.error(
                userName, "The selected monster cards do not meet the required mana, please adjust your cardSettings.txt, however, bot is retrying...")
        elif int(totalManaHave.split('/')[0]) > int(manaUsed):
            try:
                WebDriverWait(driver, 20).until(
                    EC.element_to_be_clickable((By.XPATH, "/html/body/app/div[1]/slcards/div[5]/button[1]/div[2]/span"))).click()
                if show_forge_reward:
                    WebDriverWait(driver, 20).until(
                        EC.presence_of_element_located((By.XPATH, "/html/body/app/div[1]/slcards/div[4]/div[2]/button[2]/span")))
                    if show_total_forge_balance:
                        forgebalance = WebDriverWait(driver, 10).until(
                            EC.presence_of_element_located((By.XPATH, "/html/body/app/div[1]/div[1]/app-header/section/div[4]/div[2]/div[2]/div[1]/a[1]/div[1]/span"))).text
                    WebDriverWait(driver, 10).until(
                        EC.element_to_be_clickable((By.XPATH, "/html/body/app/div[1]/slcards/div[4]/div[2]/button[2]/span")))
                    WebDriverWait(driver, 20).until(
                        EC.element_to_be_clickable((By.XPATH, "/html/body/app/div[1]/slcards/div[4]/div[2]/button[2]/span"))).click()
                    WebDriverWait(driver, 20).until(
                        EC.presence_of_element_located((By.XPATH, "/html/body/app/div[1]/slcards/div[5]/div[1]/replay/section/rewards-modal/section/div[1]/div[1]/p[1]")))
                    time.sleep(2)
                    points = driver.find_element(
                        By.XPATH, "/html/body/app/div[1]/slcards/div[5]/div[1]/replay/section/rewards-modal/section/div[1]/div[1]/div[2]/span[2]/span[2]").text
                    reward = driver.find_element(
                        By.XPATH, "/html/body/app/div[1]/slcards/div[5]/div[1]/replay/section/rewards-modal/section/div[1]/div[1]/p[1]").text
                    reward = f"Made {points} points! " + reward
                    log_info.status(userName, reward)
                    if show_total_forge_balance:
                        log_info.status(
                            userName, f"Your total balance is {forgebalance} Forge tokens.")
                    # ---------------------------------------
                    if int(points) > 50:
                        try:
                            timenow = datetime.datetime.now().replace(microsecond=0)
                            checkData = {"playername": f"{userName}",
                                         "heroesType": f"{heroesType}",
                                         "bossName": f"{bossName}",
                                         "summoners": checkDataSummoners,
                                         "monsters": checkDataMonsters,
                                         "points": f"{int(points)}",
                                         "time": f"{timenow}"}
                            collectData(checkData, public_api_endpoint)
                        except:
                            pass
                    # ---------------------------------------
                else:
                    WebDriverWait(driver, 20).until(
                        EC.presence_of_element_located((By.XPATH, "/html/body/app/div[1]/slcards/div[4]/div[2]/button[2]/span")))
                    forgebalance = WebDriverWait(driver, 10).until(
                        EC.presence_of_element_located((By.XPATH, "/html/body/app/div[1]/div[1]/app-header/section/div[4]/div[2]/div[2]/div[1]/a[1]/div[1]/span"))).text
                    log_info.status(
                        userName, f"Your total balance is {forgebalance} Forge tokens.")
                log_info.success(
                    userName, "The battle has ended!")
                return 1, manaUsed, autoSelectResult

            except:
                log_info.success(
                    userName, "Encountering difficulty in reading the game results, but the battle has ended.")
                return 1, manaUsed, autoSelectResult

        else:
            log_info.alerts(
                userName, "Insufficient stamina, entering a rest state of inactivity for 1 hour...")
            return 2, manaUsed, autoSelectResult
    except:
        driver.get("https://splinterforge.io/#/")
        log_info.error(userName,
                       "There may be some issues with the server or your playing cards, retrying in 30 seconds...")
        time.sleep(15)
        driver.refresh()
        time.sleep(15)
        return 3, manaUsed, autoSelectResult


def battleLoop(driver, userName, postingKey, heroesType, cardSelection, show_forge_reward, show_total_forge_balance, close_driver_while_sleeping, timeSleepInMinute, auto_select_card, auto_select_hero, splinterland_api_endpoint, public_api_endpoint):
    login(userName, postingKey, driver)
    while True:
        try:
            battleResult, manaUsed, autoSelectResult = battle(
                cardSelection, userName, heroesType, driver, show_forge_reward, show_total_forge_balance, auto_select_card, auto_select_hero, splinterland_api_endpoint, public_api_endpoint)
            if battleResult == 3:
                log_info("restarting...")
            if close_driver_while_sleeping:
                driver.quit()
            if battleResult == 2:
                time.sleep(3600)
            elif battleResult == 1 and timeSleepInMinute != 0:
                if autoSelectResult:
                    log_info.alerts(userName,
                                    f"this account will enter a state of inactivity for {manaUsed} minutes based on auto selected info.")
                    time.sleep(int(manaUsed)*60)
                else:
                    log_info.alerts(userName,
                                    f"According to your configuration, this account will enter a state of inactivity for {int(timeSleepInMinute/60)} minutes.")
                    time.sleep(timeSleepInMinute)
            if close_driver_while_sleeping:
                break

        except:
            pass


def start(i, accountNo, headless, close_driver_while_sleeping, show_forge_reward, show_total_forge_balance, auto_select_card, auto_select_hero, splinterland_api_endpoint, public_api_endpoint):
    while True:
        try:
            try:
                userName, postingKey, heroesType, cardSelection, timeSleepInMinute = _init(
                    accountNo)
                log_info.alerts(userName, "Initializing...")
                executable_path = "/webdrivers"
                os.environ["webdriver.chrome.driver"] = executable_path
                options = Options()
                options.add_extension('data/hivekeychain.crx')
                options.add_argument("--no-sandbox")
                options.add_argument("--disable-dev-shm-usage")
                options.add_argument("--disable-setuid-sandbox")
                options.add_argument(
                    "--disable-backgrounding-occluded-windows")
                options.add_argument("--disable-background-timer-throttling")
                options.add_argument('--disable-translate')
                options.add_argument('--disable-popup-blocking')
                options.add_argument("--disable-infobars")
                # options.add_argument("--disable-gpu")
                options.add_argument(
                    '--disable-blink-features=AutomationControlled')
                options.add_argument("--mute-audio")
                options.add_argument('--ignore-certificate-errors')
                options.add_argument('--allow-running-insecure-content')
                options.add_experimental_option(
                    "excludeSwitches", ["enable-automation"])
                options.add_experimental_option(
                    'useAutomationExtension', False)
                options.add_experimental_option(
                    "prefs", {
                        "profile.managed_default_content_settings.images": 1,
                        "profile.managed_default_content_settings.cookies": 1,
                        "profile.managed_default_content_settings.javascript": 1,
                        "profile.managed_default_content_settings.plugins": 1,
                        "profile.default_content_setting_values.notifications": 2,
                        "profile.managed_default_content_settings.stylesheets": 2,
                        "profile.managed_default_content_settings.popups": 2,
                        "profile.managed_default_content_settings.geolocation": 2,
                        "profile.managed_default_content_settings.media_stream": 2,
                    }
                )
                options.add_argument("--window-size=300,600")
                if headless:
                    options.add_argument("--headless=new")
                options.add_experimental_option(
                    'excludeSwitches', ['enable-logging'])
                try:
                    driver = webdriver.Chrome(options=options)
                except:
                    time.sleep(120)
                    log_info.error(userName,
                                   "Driver ERROR! Restart the problem might fix the issues.")
                    log_info("restarting...")
            except:
                time.sleep(5)
                log_info("restarting...")
            battleLoop(driver, userName, postingKey, heroesType, cardSelection, show_forge_reward,
                       show_total_forge_balance, close_driver_while_sleeping, timeSleepInMinute, auto_select_card, auto_select_hero, splinterland_api_endpoint, public_api_endpoint)

        except:
            try:
                driver.quit()
            except:
                pass
            log_info.error(userName,
                           "Error in this thread, restarting now.")
            time.sleep(5)
            pass


def kill_chromeanddriver():
    for process in psutil.process_iter():
        if process.name() == "chromedriver.exe":
            os.system(f"taskkill /pid {process.pid} /f > NUL")
    for proc in psutil.process_iter():
        if proc.name() == "chrome.exe":
            if "--remote-debugging-port" in " ".join(proc.cmdline()):
                os.system(f"taskkill /pid {proc.pid} /f > NUL")
                break


def startMulti(totallaccounts, headless, close_driver_while_sleeping, start_thread, start_thread_interval, show_forge_reward, show_total_forge_balance, auto_select_card, auto_select_hero, splinterland_api_endpoint, public_api_endpoint):
    # kill_chromeanddriver()
    chromedriver_autoinstaller.install()
    workers = []
    for i in range(totallaccounts):
        a = str(i + 1)
        workers.append(multiprocessing.Process(
            target=start, args=(a, a, headless, close_driver_while_sleeping, show_forge_reward, show_total_forge_balance, auto_select_card, auto_select_hero, splinterland_api_endpoint, public_api_endpoint)))
    current_threads = 0
    while workers or len(multiprocessing.active_children()) > 0:
        if current_threads < start_thread and workers:
            worker = workers.pop(0)
            try:
                multiprocessing.freeze_support()
                worker.start()
                current_threads += 1
                time.sleep(3)
            except:
                pass
        else:
            time.sleep(start_thread_interval)
            current_threads = 0


async def printSystemUsage(check_system_usage_frequency):
    while True:
        cpu_percent = psutil.cpu_percent()
        memory_info = psutil.virtual_memory()
        memory_percent = memory_info.percent
        ctypes.windll.kernel32.SetConsoleTitleW(
            f"SplinterForge Bot | CPU Usage: {cpu_percent}% | Memory Usage: {memory_percent}%")
        await asyncio.sleep(check_system_usage_frequency)


async def main():
    multiprocessing.freeze_support()
    printSkyBlue(start_font())
    time.sleep(1)
    printWelcome('Welcome to SplinterForge Bot!\nOpen source Github respositery\nhttps://github.com/Astr0-G/SplinterForge-Bot\nDiscord server\nhttps://discord.gg/pm8SGZkYcD')
    time.sleep(1)
    try:
        headless, close_driver_while_sleeping, start_thread, start_thread_interval, show_forge_reward, show_total_forge_balance, print_system_usage, check_system_usage_frequency, auto_select_card, auto_select_hero, splinterland_api_endpoint, public_api_endpoint = getConfig(
            'config/config.txt')
        totallaccounts = int(file_len("config/accounts.txt")) - 1
    except:
        print("error reading config.txt located in the config folder.")
        print("Terminating in 10 seconds...")
        time.sleep(10)
        sys.exit()
    if totallaccounts < 1:
        print("You need to add accounts in accounts.txt located in the config folder.")
        print("Terminating in 10 seconds...")
        time.sleep(10)
        sys.exit()
    if totallaccounts < start_thread:
        start_thread = totallaccounts
    printConfigSettings(totallaccounts, headless, close_driver_while_sleeping, start_thread, start_thread_interval, show_forge_reward,
                        show_total_forge_balance, print_system_usage, check_system_usage_frequency, auto_select_card, auto_select_hero)
    process_start_multi = multiprocessing.Process(target=startMulti, args=(
        totallaccounts, headless, close_driver_while_sleeping, start_thread, start_thread_interval, show_forge_reward, show_total_forge_balance, auto_select_card, auto_select_hero, splinterland_api_endpoint, public_api_endpoint))
    process_start_multi.start()
    if print_system_usage:
        await asyncio.gather(printSystemUsage(check_system_usage_frequency))

if __name__ == "__main__":
    asyncio.run(main())
