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
    global playingSummoners, playingMonster, userName, postingKey, timeSleepInMinute, bossId
    playingSummoners = []
    playingMonster = []
    try:
        f = open('config.json')
    except:
        print("error loading config.txt, please create config.txt by checking config-example.json")
        print("Closing in 10 seconds...")
        time.sleep(10)
        sys.exit()
    config = json.load(f)
    for i in config['playingSummoners']:
        playingSummoners.append(f"//div/img[@id='{i}']")
    for i in config['playingMonster']:
        playingMonster.append(f"//div/img[@id='{i}']")
    userName = config['userName']
    postingKey = config['postingKey']
    if userName == "" or postingKey == "":
        log_info.error(
            "error loading config.txt, please add user name or posting key.")
        log_info.error(
            "Closing in 10 seconds...")
        time.sleep(10)
        sys.exit()
    else:
        log_info.alerts("Loading config.txt")
    bossId = f"//div[@tabindex='{config['bossId']}']"
    timeSleepInMinute = int(config['timeSleepInMinute']) * 60
    f.close()
    chromedriver_autoinstaller.install()
    


def check():
    try:
        WebDriverWait(driver, 2).until(
            EC.element_to_be_clickable((By.XPATH, "/html/body/app/div[1]/div[1]/app-header/success-modal/section/div[1]/div[4]/div/button"))).click()
    except:
        pass


def start():
    global driver

    _init()
    executable_path = "/webdrivers"
    os.environ["webdriver.chrome.driver"] = executable_path
    options = Options()
    options.add_extension('hivekeychain.crx')
    options.add_argument("--mute-audio")
    options.add_experimental_option('useAutomationExtension', False)
    options.add_experimental_option("excludeSwitches", ["enable-logging"])
    options.add_argument('--disable-blink-features=AutomationControlled')
    driver = webdriver.Chrome(options=options)
    driver.get("chrome-extension://jcacnejopjdphbnjgfaaobbfafkihpep/popup.html")
    log_info.alerts(
        f"SplinterForge Bot is starting...")
    log_info.alerts(
        "PLEASE DONT MOVE OR CLICK ANYTHING untill chrome is minimized!")
    driver.set_window_size(350, 1200)
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
    try:
        if WebDriverWait(driver, 5).until(
                EC.presence_of_element_located((By.XPATH, "/html/body/div/div/div[4]/div"))).text == "HIVE KEYCHAIN":
            log_info.success(
                "account successful login!".format(userName))
    except:
        log_info.error(
            "login error! check your useranme or posting keys in config.txt file and retry.")
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
    forgebalance = WebDriverWait(driver, 5).until(
        EC.presence_of_element_located((By.XPATH, "/html/body/app/div[1]/div[1]/app-header/section/div[4]/div[2]/div[2]/div[1]/a[1]/div[1]/span"))).text
    log_info.success(
        f"Your Forge Balance is {forgebalance}, sleeping for 30 seconds to start, minimizing the window now.")
    driver.minimize_window()
    time.sleep(30)
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
            for i in playingMonster:
                driver.execute_script("window.scrollBy(0, -4000)")
                WebDriverWait(driver, 5).until(
                    EC.element_to_be_clickable((By.XPATH, i))).click()
                time.sleep(1)
            WebDriverWait(driver, 5).until(
                EC.element_to_be_clickable((By.XPATH, "/html/body/app/div[1]/slcards/div[5]/button[1]/div[2]/span"))).click()
            try:
                if "Stamina" == WebDriverWait(driver, 5).until(
                        EC.element_to_be_clickable((By.XPATH, "/html/body/app/div[1]/slcards/stamina-elixirs/section/div[1]/h1"))).text:
                    log_info.alerts("Not Enough Stamina, sleep for 1 hour")
                    time.sleep(1770)
            except:
                pass
            log_info.success("battle finished! Sleep for 30s...")
            time.sleep(30)
            if timeSleepInMinute != 0:
                log_info.alerts(f"sleep for {int(timeSleepInMinute/60)} mins")
            time.sleep(timeSleepInMinute)
            while True:
                try:
                    driver.get("https://splinterforge.io/#/")
                    time.sleep(1)
                    check()
                    forgebalance = WebDriverWait(driver, 5).until(
                        EC.presence_of_element_located((By.XPATH, "/html/body/app/div[1]/div[1]/app-header/section/div[4]/div[2]/div[2]/div[1]/a[1]/div[1]/span"))).text
                    log_info.success(f"Your Forge Balance is {forgebalance}.")
                    break
                except:
                    pass
        except:
            log_info.error(
                "There might be some erros with the server or your playing cards, check config and retry.")
            pass


if __name__ == '__main__':
    printSkyBlue(start_font())
    start()
