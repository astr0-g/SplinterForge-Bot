from selenium import webdriver
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
    def success(message):
        printGreen(f"[{log_info.time()}] {userName}: {message}\n")

    @staticmethod
    def error(message):
        printRed(f"[{log_info.time()}] {userName}: {message}\n")

    @staticmethod
    def alerts(message):
        print(f"[{log_info.time()}] {userName}: {message}")

    @staticmethod
    def status(message):
        printYellow(f"[{log_info.time()}] {userName}: {message}\n")

    @staticmethod
    def verify(message):
        printDarkBlue(f"[{log_info.time()}] {userName}: {message}\n")


def start_font():
    f = Figlet(font='smkeyboard', width=150)
    text = f.renderText('SplinterForge\nBot-Beta\n@lil_astr_0')
    return f"{text}"


def _init():
    global playingSummoners, playingMonster, userName, postingKey, timeSleepInMinute, bossId, skipCardIfSelectedError
    printSkyBlue(start_font())
    print("Welcome using SplinterForge Bot, for more infos about update or guide please visit: https://github.com/Astr0-G/SplinterForge-Bot")
    playingSummoners = []
    playingMonster = []
    try:
        f = open('config.json')
    except:
        print("error loading config.json, please create config.json by checking config-example.json")
        print("Closing in 10 seconds...")
        time.sleep(10)
        sys.exit()
    config = json.load(f)
    for i in config['playingSummoners']:
        playingSummoners.append(f"//div/img[@id='{i}']")
    for i in config['playingMonster']:
        playingMonster.append({
            "playingMonterDiv": f"//div/img[@id='{i}']",
            "playingMonsterId": f"{i}"
        })
    userName = config['userName']
    postingKey = config['postingKey']
    if userName == "" or postingKey == "":
        log_info.error(
            "error loading config.json, please add user name or posting key.")
        log_info.error(
            "Closing in 10 seconds...")
        time.sleep(10)
        sys.exit()
    else:
        log_info.alerts("Loading config.json")
    bossId = f"//div[@tabindex='{config['bossId']}']"
    skipCardIfSelectedError = config['skipCardIfSelectedError']
    timeSleepInMinute = int(config['timeSleepInMinute']) * 60
    f.close()
    chromedriver_autoinstaller.install()


def selectMonsterCards(i, cardId, cardDiv):
    selectMana(checkCardMana(cardId))
    WebDriverWait(driver, 10).until(
        EC.element_to_be_clickable((By.XPATH, cardDiv))).click()
    # ----------Speed up Process----------
    selectMana(checkCardMana(cardId))
    # log_info.success(f"Monster card ID {cardId} selected successful!")


def checkCard():
    for i in range(8):
        driver.execute_script("window.scrollBy(0, 10000)")
        time.sleep(0.4)
    for i in range(2):
        driver.execute_script("window.scrollBy(0, -10000)")


def selectMana(mana):
    if int(mana) < 10:
        WebDriverWait(driver, 2).until(
            EC.element_to_be_clickable((By.XPATH, f"/html/body/app/div[1]/slcards/div[5]/section[1]/div/div[2]/div/div[1]/sl-card-filters/div[1]/div[{int(mana)+1}]"))).click()
    else:
        WebDriverWait(driver, 2).until(
            EC.element_to_be_clickable((By.XPATH, f"/html/body/app/div[1]/slcards/div[5]/section[1]/div/div[2]/div/div[1]/sl-card-filters/div[1]/div[11]"))).click()


def checkCardMana(cardid):
    f = open('data/cardsDetails.json')
    mana = json.load(f)
    for i in range(len(mana)):
        if int(cardid) == int(mana[i]['id']):
            return(mana[i]['stats']['mana'][0])


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


def start():
    global driver

    _init()
    executable_path = "/webdrivers"
    os.environ["webdriver.chrome.driver"] = executable_path
    options = Options()
    options.add_extension('data/hivekeychain.crx')
    options.add_argument("--mute-audio")
    options.add_argument("--window-size=300,1080")
    options.add_argument("--headless=new")
    options.add_argument('--ignore-certificate-errors')
    options.add_argument('--allow-running-insecure-content')
    options.add_argument('--disable-blink-features=AutomationControlled')
    options.add_experimental_option("excludeSwitches", ["enable-logging"])
    options.add_experimental_option('useAutomationExtension', False)
    driver = webdriver.Chrome(options=options)
    driver.get("chrome-extension://jcacnejopjdphbnjgfaaobbfafkihpep/popup.html")
    log_info.alerts(
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
            log_info.success(
                "Account successful login!")
    except:
        log_info.error(
            "login error! check your useranme or posting keys in config.json file and retry.")
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
    finishRound = 0
    forgebalance = WebDriverWait(driver, 5).until(
        EC.presence_of_element_located((By.XPATH, "/html/body/app/div[1]/div[1]/app-header/section/div[4]/div[2]/div[2]/div[1]/a[1]/div[1]/span"))).text
    log_info.success(
        f"Sleeping for 20 seconds to start...")
    ctypes.windll.kernel32.SetConsoleTitleW(
        f"SplinterForge Bot | Forge Balance : {forgebalance} | Finished Round : {finishRound}")
    time.sleep(20)
    log_info.success(
        f"Starting playing...")
    while True:
        try:
            driver.get("https://splinterforge.io/#/slcards")
            check()
            WebDriverWait(driver, 5).until(
                EC.element_to_be_clickable((By.XPATH, bossId))).click()
            WebDriverWait(driver, 5).until(
                EC.element_to_be_clickable((By.XPATH, "/html/body/app/div[1]/slcards/div[5]/section[1]/div/div[1]/div[2]/button"))).click()
            log_info.alerts("Selecting summoners and monsters...")
            for i in playingSummoners:
                driver.execute_script("window.scrollBy(0, -4000)")
                WebDriverWait(driver, 5).until(
                    EC.element_to_be_clickable((By.XPATH, i))).click()
                time.sleep(1)
            log_info.success("Summoners selected successful!")
            checkCard()
            for i in playingMonster:
                cardId = i["playingMonsterId"]
                cardDiv = i["playingMonterDiv"]
                if skipCardIfSelectedError:
                    try:
                        selectMonsterCards(i, cardId, cardDiv)
                    except:
                        driver.get_screenshot_as_file(
                            f"errorSelectedMonterCardId{cardId}.png")
                        log_info.alerts(
                            f"Error select card ID: {cardId}, skipped this card, check errorSelectedMonterCardId{cardId}.png to fix...")
                        pass
                else:
                    selectMonsterCards(i)
            log_info.success("Monster selected successful!")
            manaUsed = WebDriverWait(driver, 5).until(
                EC.presence_of_element_located((By.XPATH, "/html/body/app/div[1]/slcards/div[5]/div[2]/div[1]/div[1]/button/span"))).text
            totalManaHave = WebDriverWait(driver, 5).until(
                EC.presence_of_element_located((By.XPATH, "/html/body/app/div[1]/div[1]/app-header/section/div[4]/div[2]/div[2]/div[1]/a[2]/div[1]/span"))).text
            if int(totalManaHave.split('/')[0]) > int(manaUsed):
                WebDriverWait(driver, 5).until(
                    EC.element_to_be_clickable((By.XPATH, "/html/body/app/div[1]/slcards/div[5]/button[1]/div[2]/span"))).click()
                time.sleep(4)
                log_info.success("Battle finished! Sleep for 30s...")
                finishRound += 1
                forgebalance = WebDriverWait(driver, 10).until(
                    EC.presence_of_element_located((By.XPATH, "/html/body/app/div[1]/div[1]/app-header/section/div[4]/div[2]/div[2]/div[1]/a[1]/div[1]/span"))).text
                ctypes.windll.kernel32.SetConsoleTitleW(
                    f"SplinterForge Bot | Forge Balance : {forgebalance} | Finished Round : {finishRound}")
                time.sleep(30)
                if timeSleepInMinute != 0:
                    log_info.alerts(
                        f"Sleep for {int(timeSleepInMinute/60)} mins...")
                time.sleep(timeSleepInMinute)
            else:
                log_info.alerts("Not enough stamina, sleep for 1 hour...")
                time.sleep(1800)
            while True:
                try:
                    driver.get("https://splinterforge.io/#/")
                    time.sleep(1)
                    check()
                    forgebalance = WebDriverWait(driver, 5).until(
                        EC.presence_of_element_located((By.XPATH, "/html/body/app/div[1]/div[1]/app-header/section/div[4]/div[2]/div[2]/div[1]/a[1]/div[1]/span"))).text
                    break
                except:
                    pass
        except:
            log_info.error(
                "There might be some erros with the server or your playing cards, check config and retry.")
            time.sleep(10)
            pass


if __name__ == '__main__':
    start()
