#!/ur/bin/env python3
import collections
import json
import re
import smtplib
import ssl
import threading
import time

import requests

searches = {
    r"[\s:]*([0-9a-zA-Z]{64})[\s]*": "btc private key",
    r"[\s:]*([0-9a-zA-Z]{32})[\s]*": "btc address",
    r"[\s:]*([59][0-9a-zA-Z]{50})[\s]*": "Base58 Wallet Import format"
}

pastebin_listing_url = "https://scrape.pastebin.com/api_scraping.php"
pastebin_listing_params = {"limit":"100"}
pastebin_scrape_url = "https://scrape.pastebin.com/api_scrape_item.php?i={}"


def main():
    old_listing = set()
    current_listing = set()
    new_elements = set()
    while (1):
        current_listing = get_updates() # slow operation, but it has to be done
        new_elements = current_listing - old_listing
        old_listing = current_listing

        for item in new_elements:
            new_thread = threading.Thread(target=parse_paste, args=(item,), daemon=True)
            new_thread.start()

        time.sleep(1) #sleep for a second no matter what. pastes come in slow most of the time

def get_updates():
    ret = set()
    listing = requests.get(url=pastebin_listing_url, params=pastebin_listing_params)
    if listing.status_code == 200:
        pastes = json.loads(listing.text)
        for items in pastes:
            ret.add(items["key"])
    else:
        raise ConnectionError("Status code: {}".format(listing.status_code))
    return ret

def parse_paste(paste_id):
    paste = get_paste(paste_id)
    if not paste:
        print("[-] Error: Unable to retrieve paste {}".format(paste_id))
        return None #currently we just fail if we have a problem when trying to get a paste
    results = search_paste(paste)
    if results:
        send_results(results)

    return None

def get_paste(paste_id):
    try:
        paste = requests.get(url=pastebin_scrape_url.format(paste_id))
        if paste.status_code == 200:
            return paste.text
    except Exception as e:
        print("Connection error: {}".format(e))
    return None #we get here if any kind of error occurred

def search_paste(paste):
    return None

def send_results(results):
    pass

if __name__ == "__main__":
    main()
