import os
import sys
import time
import ctypes
import re
import shutil
import tempfile
ctypes.windll.kernel32.SetConsoleTitleW("SplinterForge Bot Cleanup Script")


def process():
    try:
        os.system("taskkill /im chromedriver.exe /f")
        os.system("taskkill /im chrome.exe /f")
        os.system("taskkill /im python.exe /f")

    except:
        print("Error Encountered while running script")


def clean():
    cache_dir = tempfile.gettempdir() + "\\chromedriver\\Cache"
    if os.path.exists(cache_dir):
        shutil.rmtree(cache_dir)
    profile_dir = tempfile.gettempdir() + "\\chromedriver\\Profile"
    if os.path.exists(profile_dir):
        shutil.rmtree(profile_dir)
    chromedriver_exe = tempfile.gettempdir() + "\\chromedriver.exe"
    if os.path.exists(chromedriver_exe):
        os.remove(chromedriver_exe)
    program_files = "C:\\Program Files (x86)"

    for dirpath, dirnames, filenames in os.walk(program_files):
        for dirname in dirnames:
            if re.match(r"^scoped_dir\d+_\d+$", dirname):
                dir_path = os.path.join(dirpath, dirname)
                try:
                    shutil.rmtree(dir_path)
                    print(f"Deleted directory: {dir_path}")
                except Exception as e:
                    print(f"Failed to delete directory: {dir_path}")
                    print(f"Error: {e}")

    temp_directory = os.environ["TEMP"]

    for filename in os.listdir(temp_directory):
        file_path = os.path.join(temp_directory, filename)
        try:
            if os.path.isfile(file_path):
                os.unlink(file_path)
            elif os.path.isdir(file_path):
                shutil.rmtree(file_path)
        except Exception as e:
            print(f"Failed to delete file or directory: {file_path}")
            print(f"Error: {e}")


print("Please make sure you are ok about cleaning up all chrome.exe, chromedriver.exe, python.exe process currently running in this computer and also files were generated from running the bot, which means SplinterForge Bot, any running chrome browser, and this script will be closing after performed. please type 'yes' and then press 'ENTER' on keyboard to confirm:")
text = input().lower()
if text == "yes":
    clean()
    process()
else:
    print("Cleanup cancelled... Exit in 5 seconds...")
    time.sleep(5)
    sys.exit()
