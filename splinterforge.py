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
import time
import os
import sys
import ctypes

from colorama import init

init(convert=True)
STD_OUTPUT_HANDLE = -11
std_out_handle = ctypes.windll.kernel32.GetStdHandle(STD_OUTPUT_HANDLE)
ctypes.windll.kernel32.SetConsoleTitleW("SplinterForge Bot")


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
def printDarkBlue(mess):
    set_cmd_text_color(FOREGROUND_DARKBLUE)
    sys.stdout.write(mess)
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


def start_font():
    f = Figlet(font='smkeyboard', width=150)
    text = f.renderText('SplinterForge\nBot-Beta\n@lil_astr_0')
    return f"{text}"


def file_len(file_path):
    with open(file_path, 'r') as rf:
        record = 0
        for line in rf:
            if line != "\n":
                record += 1
        return record


def getAccountData(file_path, line_number):
    acctinfo = linecache.getline(file_path, int(line_number)+1).strip()
    time.sleep(1)
    userName = acctinfo.split(":")[0]
    postingKey = acctinfo.split(":")[1]
    return userName, postingKey


def getCardDetails(file_path, line_number):
    acctinfo = linecache.getline(file_path, int(line_number)+1).strip()
    bossId = acctinfo.split(":")[0]
    playingSummoners = acctinfo.split(":")[1].split(',')
    playingMonster = acctinfo.split(":")[2].split(',')
    timeSleepInMinute = acctinfo.split(":")[3]
    return bossId, playingSummoners, playingMonster, timeSleepInMinute


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
        else:
            log_info.alerts(userName, "Loaded account.")
    except:
        print("error loading accounts.txt, retrying...")
        time.sleep(10)
        sys.exit()
    try:
        cardSelection = []
        playingSummonersList = []
        playingMonsterList = []
        bossId, playingSummoners, playingMonster, timeSleepInMinute = getCardDetails(
            "config/cardSettings.txt", accountNo)
        bossId = f"//div[@tabindex='{bossId}']"
        for i in list(playingSummoners):
            playingSummonersList.append(f"//div/img[@id='{i}']")
        for i in playingMonster:
            playingMonsterList.append({
                "playingMonterDiv": f"//div/img[@id='{i}']",
                "playingMonsterId": f"{i}"
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
    log_info.alerts(userName, "Loaded card settings.")
    return userName, postingKey, cardSelection, timeSleepInMinute


def selectMonsterCards(userName, i, cardId, cardDiv):
    scroolTime = 0
    while True:
        try:
            if scroolTime < 9:
                WebDriverWait(driver, 2).until(
                    EC.element_to_be_clickable((By.XPATH, cardDiv))).click()
                time.sleep(1)
                selectNumber = driver.find_element(
                    By.XPATH, "/html/body/app/div[1]/slcards/div[5]/div[1]/h3/span[2]/div[1]/button/span").text
                if i == int(selectNumber):
                    result = True
                    log_info.success(
                        userName, f"Monster card ID {cardId} selected successful!")
                    break
            else:
                log_info.error(userName,
                               f"Error select card ID: {cardId}, skipped this card...")
                result = False
                break
        except:
            driver.execute_script("window.scrollBy(0, 1800)")
            scroolTime += 1
            pass
    driver.execute_script("window.scrollBy(0, -10000)")
    return result


def check():
    try:
        WebDriverWait(driver, 2).until(
            EC.element_to_be_clickable((By.XPATH, "/html/body/app/div[1]/div[1]/app-header/success-modal/section/div[1]/div[4]/div/button"))).click()
    except:
        pass
    try:
        WebDriverWait(driver, 2).until(
            EC.element_to_be_clickable((By.XPATH, "/html/body/app/div[1]/login-modal/div/div/div/div[2]/div[3]/button"))).click()
    except:
        pass


def start(i, accountNo):
    global driver
    while True:
        try:
            userName, postingKey, cardSelection, timeSleepInMinute = _init(
                accountNo)
            executable_path = "/webdrivers"
            os.environ["webdriver.chrome.driver"] = executable_path
            options = Options()
            options.add_extension('data/hivekeychain.crx')
            options.add_argument("--mute-audio")
            options.add_argument("--window-size=300,1080")
            options.add_argument("--headless=new")
            options.add_argument('--ignore-certificate-errors')
            options.add_argument('--allow-running-insecure-content')
            options.add_argument(
                '--disable-blink-features=AutomationControlled')
            options.add_experimental_option(
                "excludeSwitches", ["enable-logging"])
            options.add_experimental_option('useAutomationExtension', False)
            driver = webdriver.Chrome(options=options)
            thisDriver = driver
            driver.get(
                "chrome-extension://jcacnejopjdphbnjgfaaobbfafkihpep/popup.html")
            log_info.alerts(userName,
                            f"SplinterForge Bot is starting...")
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
                time.sleep(1)
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
            log_info.success(userName,
                             f"Sleeping for 20 seconds to start...")
            time.sleep(20)
            log_info.success(userName,
                             f"Starting playing...")
            while True:
                try:
                    driver.get("https://splinterforge.io/#/slcards")
                    check()
                    for j in range(len(cardSelection)):
                        WebDriverWait(driver, 5).until(
                            EC.element_to_be_clickable((By.XPATH, cardSelection[j]['bossId']))).click()
                        if "BOSS IS DEAD" == WebDriverWait(driver, 5).until(
                                EC.presence_of_element_located((By.XPATH, "/html/body/app/div[1]/slcards/div[5]/section[1]/div/div[1]/div[2]/button"))).text:
                            log_info.error(userName,
                                           "Boss you selected is dead, please change your boss Id in config.json and restart the bot.")
                            time.sleep(10*18)
                        WebDriverWait(driver, 5).until(
                            EC.element_to_be_clickable((By.XPATH, "/html/body/app/div[1]/slcards/div[5]/section[1]/div/div[1]/div[2]/button"))).click()
                        log_info.alerts(
                            userName, "Selecting summoners and monsters...")
                        for i in cardSelection[j]['playingSummoners']:
                            driver.execute_script("window.scrollBy(0, -4000)")
                            WebDriverWait(driver, 5).until(
                                EC.element_to_be_clickable((By.XPATH, i))).click()
                        log_info.success(
                            userName, "Summoners selected successful!")
                        seletedNum = 1
                        for i in cardSelection[j]['playingMonsterId']:
                            cardId = i["playingMonsterId"]
                            cardDiv = i["playingMonterDiv"]
                            result = selectMonsterCards(
                                userName, seletedNum, cardId, cardDiv)
                            time.sleep(1)
                            if result:
                                seletedNum += 1
                        log_info.success(
                            userName, "Monsters selected successful!")
                        manaUsed = WebDriverWait(driver, 5).until(
                            EC.presence_of_element_located((By.XPATH, "/html/body/app/div[1]/slcards/div[5]/div[2]/div[1]/div[1]/button/span"))).text
                        totalManaHave = WebDriverWait(driver, 5).until(
                            EC.presence_of_element_located((By.XPATH, "/html/body/app/div[1]/div[1]/app-header/section/div[4]/div[2]/div[2]/div[1]/a[2]/div[1]/span"))).text
                        if int(totalManaHave.split('/')[0]) > int(manaUsed):
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
                            log_info.status(userName, f"You have total balance of {forgebalance} Forge token.")
                            log_info.success(
                                userName, "Battle finished! Sleep for 30s...")
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
                                check()
                                break
                            except:
                                pass
                except:
                    driver.get("https://splinterforge.io/#/")
                    time.sleep(1)
                    log_info.error(userName,
                                   "There might be some errors with the server or your playing cards, retrying in 10 seconds...")
                    time.sleep(10)
                    pass
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
    print(f"Loaded {(totallaccounts - 1)} accounts, and {(totallaccounts - 1)} seconds waiting time will be applied between each account starting.")
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
