import pyautogui
import time

def click_every_five_seconds():
    try:
        while True:
            pyautogui.click()  # Simulate a left click
            time.sleep(1)  # Wait for 5 seconds
    except KeyboardInterrupt:
        print("Script stopped.")

if __name__ == "__main__":
    time.sleep(10)
    click_every_five_seconds()
