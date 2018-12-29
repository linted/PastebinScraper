#!/ur/bin/env python3
import re

import requests

searches = [
    r"[\s:]*([0-9a-zA-Z]{64})[\s]*", #btc private key
    r"[\s:]*([0-9a-zA-Z]{32})[\s]*", #btc address
    r"[\s:]*([59][0-9a-zA-Z]{50})[\s]*" #Base58 Wallet Import format
]

pastbin_scrape_url = "https://scrape.pastebin.com/api_scraping.php"
pastbin_scrape_params = {"limit":"100"}

def main():
    pass

def get_updates():
    r = requests.get(url=pastbin_scrape_url, params=pastbin_scrape_params)
    print(r)
    

if __name__ == "__main__":
    main()
