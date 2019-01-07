#!/ur/bin/env python3
import argparse
import collections
import getpass
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
    parser = argparse.ArgumentParser()
    parser.add_argument("-e", "--send-email", help="email to send from", required=True)
    parser.add_argument("-r", "--recv-email", help="email to send to", required=True)
    parser.add_argument("-s", "--smtp-server", help="smtp server to talk to", required=True)
    password = getpass.getpass()

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
    title = []
    found = False
    for regex in searches.keys():
        matches = re.match(regex, paste)
        if matches:
            #yay we found a match! as of now we don't actually use any of the matched stuff
            title.append(searches[regex]) #adds the match description to the title
            found = True
    if found:
        return (" ".join(title),paste)
    return None

def send_results(results):
    pass

def setup_email(email, password):
    context = ssl.create_default_context()
    server =  smtplib.SMTP(smtp_server, port)
    server.starttls(context=context)
    server.login(sender_email, password)
    server.sendmail(sender_email, receiver_email, message)

if __name__ == "__main__":
    main()
