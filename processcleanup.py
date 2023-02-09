import os
import sys
import time
import ctypes

ctypes.windll.kernel32.SetConsoleTitleW("SplinterForge Cleanup Script")


def process():
    try:
        os.system("taskkill /im chromedriver.exe /f")
        os.system("taskkill /im chrome.exe /f")
        os.system("taskkill /im python.exe /f")

    except:
        print("Error Encountered while running script")


print("Please make sure you are ok about cleaning up all chrome.exe, chromedriver.exe, python.exe process currently running in this computer, which means SplinterForge Bot, any running chrome browser, and this script will be closing after performed. please type 'yes' and then press 'ENTER' on keyboard to confirm:")
text = input().lower()
if text == "yes":
    process()
else:
    print("Cleanup cancelled... Exit in 5 seconds...")
    time.sleep(5)
    sys.exit()
