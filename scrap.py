#!/ur/bin/env python3
import collections
import json
import re
import threading
import time

import requests

searches = [
    r"[\s:]*([0-9a-zA-Z]{64})[\s]*", #btc private key
    r"[\s:]*([0-9a-zA-Z]{32})[\s]*", #btc address
    r"[\s:]*([59][0-9a-zA-Z]{50})[\s]*" #Base58 Wallet Import format
]

pastbin_listing_url = "https://scrape.pastebin.com/api_scraping.php"
pastbin_listing_params = {"limit":"100"}
pastbin_scrape_url = "https://scrape.pastebin.com/api_scrape_item.php?i="


def main():
    old_listing = set()
    current_listing = set()
    new_elements = set()
    while (1):
        current_listing = get_updates() # slow operation, but it has to be done
        new_elements = current_listing - old_listing
        old_listing = current_listing

        for item in new_elements:
            threading.Thread(target=parse_paste, args=(item,), daemon=True)

        time.sleep(1) #sleep for a second no matter what. pastes come in slow most of the time

def get_updates():
    ret = set()
    listing = requests.get(url=pastbin_listing_url, params=pastbin_listing_params)
    if listing.status_code == 200:
        pastes = json.loads(listing.text)
        for items in pastes:
            ret.add(items["key"])
    else:
        raise ConnectionError("Status code: {}".format(listing.status_code))
    return ret

def parse_paste():
    pass


if __name__ == "__main__":
    main()
