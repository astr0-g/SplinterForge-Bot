from selenium import webdriver
import multiprocessing
import linecache
from selenium.webdriver.common.keys import Keys
import chromedriver_autoinstaller
from selenium.webdriver.chrome.options import Options
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
from prettytable import PrettyTable
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


# dark blue
def printDarkBlue(message):
    set_cmd_text_color(FOREGROUND_DARKBLUE)
    sys.stdout.write(message)
    resetColor()


# green
def printGreen(message):
    set_cmd_text_color(FOREGROUND_GREEN)
    sys.stdout.write(message)
    resetColor()

# sky blue


def printSkyBlue(message):
    set_cmd_text_color(FOREGROUND_SKYBLUE)
    sys.stdout.write(message)
    resetColor()


# red
def printRed(message):
    set_cmd_text_color(FOREGROUND_DARKRED)
    sys.stdout.write(message)
    resetColor()


# yellow
def printYellow(message):
    set_cmd_text_color(FOREGROUND_YELLOW)
    sys.stdout.write(message)
    resetColor()

# pink


def printPinky(message):
    set_cmd_text_color(FOREGROUND_Pink)
    sys.stdout.write(message)
    resetColor()


# reset white
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


def print_result_box(userName, data):
    table = PrettyTable(["Card", "ID", "Name", "Selection Results"])
    for row in data:
        table.add_row(row)
    dataToPrint = table.get_string(sortby="Card")
    log_info.alerts(userName, "Card selection results:")
    print(dataToPrint)


def start_font():
    f = Figlet(font='smkeyboard', width=150)
    text = f.renderText('SplinterForge\nBot-Beta\n@lil_astr_0')
    return f"{text}"


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
                           "error loading accounts.txt, please add user name or posting key.")
            log_info.error(userName,
                           "Closing in 10 seconds...")
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
        # bossId = f"//div[@tabindex='{bossId}']"
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
        print("error loading cardSettings.txt")
        print("Closing in 10 seconds...")
        time.sleep(10)
        sys.exit()
    return userName, postingKey, heroesType, cardSelection, timeSleepInMinute


def selectSummoners(userName, seletedNumOfSummoners, cardDiv, driver):
    scroolTime = 0
    result = False
    time.sleep(1)
    while scroolTime < 5:
        try:
            WebDriverWait(driver, 1).until(
                EC.element_to_be_clickable((By.XPATH, cardDiv))).click()
            time.sleep(1)
            selectNumber = driver.find_element(
                By.XPATH, "/html/body/app/div[1]/slcards/div[5]/section[1]/div/div[1]/div[2]/div[1]/div[1]/h3/span[2]").text
            if seletedNumOfSummoners == int(selectNumber.split("/")[0]):
                result = True
                break
        except:
            driver.execute_script("window.scrollBy(0, 180)")
            scroolTime += 1
    if not result:
        log_info.error(userName, "Error select summoners, will retrying...")
    driver.execute_script("window.scrollBy(0, -4000)")
    time.sleep(1)
    return result


def selectMonsterCards(userName, seletedNumOfMonsters, cardId, cardDiv, driver):
    scroolTime = 0
    while scroolTime < 15:
        try:
            WebDriverWait(driver, 0.5).until(
                EC.element_to_be_clickable((By.XPATH, cardDiv))).click()
            time.sleep(1)
            selectNumber = driver.find_element(
                By.XPATH, "/html/body/app/div[1]/slcards/div[5]/div[1]/h3/span[2]/div[1]/button/span").text
            if seletedNumOfMonsters == int(selectNumber):
                result = True
                # log_info.success(
                #     userName, f"Monster card ID {cardId} selected successful!")
                break
        except:
            driver.execute_script("window.scrollBy(0, 450)")
            scroolTime += 1
            pass
    else:
        result = False
        # log_info.error(userName,
        #                f"Error select card ID: {cardId}, skipped this card...")

    driver.execute_script("window.scrollBy(0, -10000)")
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


def heroSelect(heroesType, userName, driver):
    check(driver)
    try:
        heroesType = int(heroesType)
        if heroesType != 1 or heroesType != 2 or heroesType != 3:
            WebDriverWait(driver, 10).until(
                EC.presence_of_element_located((By.XPATH, "/html/body/app/div[1]/div[1]/app-header/section/div[4]/div[2]/div[1]/a[4]"))).click()
            WebDriverWait(driver, 10).until(
                EC.presence_of_element_located((By.XPATH, "/html/body/app/div[1]/splinterforge-heros/div[3]/section/div/div/div[2]/div[1]"))).click()
            WebDriverWait(driver, 10).until(
                EC.presence_of_element_located((By.XPATH, f"/html/body/app/div[1]/splinterforge-heros/div[3]/section/div/div/div[2]/div[1]/ul/li[{heroesType}]"))).click()
            if heroesType == 1:
                log_info.success(userName,
                                 f"Selected hero type: Warrior")
            elif heroesType == 2:
                log_info.success(userName,
                                 f"Selected hero type: Wizard")
            elif heroesType == 3:
                log_info.success(userName,
                                 f"Selected hero type: Ranger")
    except:
        log_info.error(userName,
                       f"Error selecting heros type, continue...")
        pass


def login(userName, postingKey, driver):
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
        time.sleep(2)
        WebDriverWait(driver, 5).until(
            EC.element_to_be_clickable((By.XPATH, "/html/body/div/div/div[1]/div[2]/div/div[2]/div[2]/div/input"))).send_keys(Keys.ENTER)
        if WebDriverWait(driver, 5).until(
                EC.presence_of_element_located((By.XPATH, "/html/body/div/div/div[4]/div"))).text == "HIVE KEYCHAIN":
            log_info.success(userName,
                             "Account successful login!")
    except:
        log_info.error(userName,
                       "login error! check your useranme or posting keys in config.json file and retry.")
        time.sleep(15)
        driver.close()
    driver.get("https://splinterforge.io/#/")
    driver.set_window_size(1920, 1080)
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
            break
        except:
            pass


def battle(cardSelection, userName, timeSleepInMinute, driver):
    try:
        check(driver)
        WebDriverWait(driver, 10).until(
            EC.element_to_be_clickable((By.XPATH, "/html/body/app/div[1]/div[1]/app-header/section/div[4]/div[2]/div[1]/a[5]/div[1]"))).click()
        for j in range(len(cardSelection)):
            bossIdToSelect = cardSelection[j]['bossId']
            while True:
                WebDriverWait(driver, 5).until(
                    EC.element_to_be_clickable((By.XPATH, f"//div[@tabindex='{bossIdToSelect}']"))).click()
                if "BOSS IS DEAD" != WebDriverWait(driver, 5).until(
                        EC.presence_of_element_located((By.XPATH, "/html/body/app/div[1]/slcards/div[5]/section[1]/div/div[1]/div[2]/button"))).text:
                    time.sleep(1)
                    break
                else:
                    log_info.error(userName,
                                   "Boss you selected is dead, automatically selecting another one...")
                    if int(bossIdToSelect) < 17:
                        bossIdToSelect = str(int(bossIdToSelect) + 1)
                    else:
                        bossIdToSelect = "14"
                    WebDriverWait(driver, 10).until(
                        EC.element_to_be_clickable((By.XPATH, "/html/body/app/div[1]/div[1]/app-header/section/div[4]/div[2]/div[1]/a[5]/div[1]"))).click()
            WebDriverWait(driver, 5).until(
                EC.element_to_be_clickable((By.XPATH, "/html/body/app/div[1]/slcards/div[5]/section[1]/div/div[1]/div[2]/button"))).click()
            log_info.alerts(
                userName, "Selecting summoners and monsters...")
            printData = []
            seletedNumOfSummoners = 1
            for i in cardSelection[j]['playingSummoners']:
                cardId = i["playingSummonersId"]
                cardDiv = i["playingSummonersDiv"]
                cardName = i["playingSummonersName"]
                resultSeletetSummoners = selectSummoners(
                    userName, seletedNumOfSummoners, cardDiv, driver)
                seletedNumOfSummoners += 1
                if resultSeletetSummoners:
                    printData.append(
                        [f"Summoners #{int(cardSelection[j]['playingSummoners'].index(i)) + 1}", cardId, cardName, "✔️"])
                else:
                    log_info("restarting...")
            log_info.success(
                userName, "Summoners selected successful!")
            seletedNumOfMonsters = 1
            for i in cardSelection[j]['playingMonsterId']:
                cardId = i["playingMonstersId"]
                cardDiv = i["playingMontersDiv"]
                cardName = i["playingMonstersName"]
                resultSeletetMonsters = selectMonsterCards(
                    userName, seletedNumOfMonsters, cardId, cardDiv, driver)
                time.sleep(1)
                if resultSeletetMonsters:
                    printData.append(
                        [f"Monsters #{int(cardSelection[j]['playingMonsterId'].index(i))+1}", cardId, cardName, "✔️"])
                    seletedNumOfMonsters += 1
                else:
                    printData.append(
                        [f"Monsters #{int(cardSelection[j]['playingMonsterId'].index(i))+1}", cardId, cardName, "❌"])
            print_result_box(userName, printData)
            manaUsed = WebDriverWait(driver, 5).until(
                EC.presence_of_element_located((By.XPATH, "/html/body/app/div[1]/slcards/div[5]/div[2]/div[1]/div[1]/button/span"))).text
            totalManaHave = WebDriverWait(driver, 5).until(
                EC.presence_of_element_located((By.XPATH, "/html/body/app/div[1]/div[1]/app-header/section/div[4]/div[2]/div[2]/div[1]/a[2]/div[1]/span"))).text
            if int(manaUsed) <= 15:
                log_info.error(
                    userName, "Monster cards selected not met required mana, change your cardSettings.txt, but the bot will retrying...")
            elif int(totalManaHave.split('/')[0]) > int(manaUsed):
                try:
                    WebDriverWait(driver, 5).until(
                        EC.element_to_be_clickable((By.XPATH, "/html/body/app/div[1]/slcards/div[5]/button[1]/div[2]/span"))).click()
                    forgebalance = WebDriverWait(driver, 10).until(
                        EC.presence_of_element_located((By.XPATH, "/html/body/app/div[1]/div[1]/app-header/section/div[4]/div[2]/div[2]/div[1]/a[1]/div[1]/span"))).text
                    WebDriverWait(driver, 5).until(
                        EC.element_to_be_clickable((By.XPATH, "/html/body/app/div[1]/slcards/div[4]/div[2]/button[2]/span"))).click()
                    WebDriverWait(driver, 20).until(
                        EC.presence_of_element_located((By.XPATH, "/html/body/app/div[1]/slcards/div[5]/div[1]/replay/section/rewards-modal/section/div[1]/div[1]/p[1]")))
                    time.sleep(2)
                    reward = driver.find_element(
                        By.XPATH, "/html/body/app/div[1]/slcards/div[5]/div[1]/replay/section/rewards-modal/section/div[1]/div[1]/p[1]").text
                    log_info.status(userName, reward)
                    log_info.status(
                        userName, f"You have total balance of {forgebalance} Forge token.")
                    log_info.success(
                        userName, "Battle finished! Sleep for 30s...")
                    time.sleep(30)
                except:
                    log_info.success(
                        userName, "Having difficulty reading rewards but battle finished! Sleep for 30s...")

                    time.sleep(30)
                    continue
                if timeSleepInMinute != 0:
                    log_info.alerts(userName,
                                    f"Sleep for {int(timeSleepInMinute/60)} mins...")
                    time.sleep(timeSleepInMinute)

            else:
                log_info.alerts(
                    userName, "Not enough stamina, sleep for 1 hour...")
                time.sleep(1800)
            while True:
                try:
                    driver.get("https://splinterforge.io/#/")
                    time.sleep(1)
                    driver.get("https://splinterforge.io/#/")
                    check(driver)
                    break
                except:
                    pass
    except:
        driver.get("https://splinterforge.io/#/")
        log_info.error(userName,
                       "There might be some errors with the server or your playing cards, retrying in 30 seconds...")
        time.sleep(15)
        driver.get("https://splinterforge.io/#/")
        time.sleep(15)
        pass


def start(i, accountNo):

    while True:
        try:
            userName, postingKey, heroesType, cardSelection, timeSleepInMinute = _init(
                accountNo)
            log_info.alerts(userName, "Initializing...")
            prefs = {"profile.managed_default_content_settings.images": 1,
                     "profile.managed_default_content_settings.cookies": 1,
                     "profile.managed_default_content_settings.javascript": 1,
                     "profile.managed_default_content_settings.plugins": 1,
                     "profile.default_content_setting_values.notifications": 2,
                     "profile.managed_default_content_settings.stylesheets": 2,
                     "profile.managed_default_content_settings.popups": 2,
                     "profile.managed_default_content_settings.geolocation": 2,
                     "profile.managed_default_content_settings.media_stream": 2,
                     }
            executable_path = "/webdrivers"
            os.environ["webdriver.chrome.driver"] = executable_path
            options = Options()
            options.add_extension('data/hivekeychain.crx')
            options.add_argument("--mute-audio")
            options.add_argument("--window-size=300,1080")
            # options.add_argument("--headless=new")
            options.add_argument('--ignore-certificate-errors')
            options.add_argument('--allow-running-insecure-content')
            options.add_argument(
                '--disable-blink-features=AutomationControlled')
            options.add_experimental_option(
                "excludeSwitches", ["enable-logging"])
            options.add_experimental_option('useAutomationExtension', False)
            options.experimental_options["prefs"] = prefs
            driver = webdriver.Chrome(options=options)
            thisDriver = driver
            driver.get(
                "chrome-extension://jcacnejopjdphbnjgfaaobbfafkihpep/popup.html")
            login(userName, postingKey, driver)
            log_info.alerts(userName,
                            f"Selecting heros type...")
            heroSelect(heroesType, userName, driver)
            time.sleep(10)
            log_info.success(userName,
                             f"Entering battles...")
            while True:
                battle(cardSelection, userName, timeSleepInMinute, driver)
        except:
            thisDriver.quit()
            log_info.error(userName,
                           "Process error, restarting...")
            pass


def startMulti():
    multiprocessing.freeze_support()
    chromedriver_autoinstaller.install()
    printSkyBlue(start_font())
    print("Welcome using SplinterForge Bot, for more infos about update or guide please visit: https://github.com/Astr0-G/SplinterForge-Bot")
    totallaccounts = int(file_len("config/accounts.txt"))
    print(f"Total accounts loaded: {totallaccounts - 1}")
    if totallaccounts <= 1:
        print("Please add accounts to accounts.txt in config folder")
    # os.system("taskkill /im chromedriver.exe /f")
    workers = []
    for i in range(totallaccounts - 1):
        a = str(i + 1)
        workers.append(multiprocessing.Process(
            target=start, args=(a, a)))
    for w in workers:
        try:
            multiprocessing.freeze_support()
            w.start()
            time.sleep(totallaccounts)
        except:
            pass


if __name__ == '__main__':
    startMulti()
