#!/ur/bin/env python3
import collections
import json
import re

import requests

searches = [
    r"[\s:]*([0-9a-zA-Z]{64})[\s]*", #btc private key
    r"[\s:]*([0-9a-zA-Z]{32})[\s]*", #btc address
    r"[\s:]*([59][0-9a-zA-Z]{50})[\s]*" #Base58 Wallet Import format
]

pastbin_scrape_url = "https://scrape.pastebin.com/api_scraping.php"
pastbin_scrape_params = {"limit":"100"}

class OrderedSet(collections.OrderedDict, collections.MutableSet):
    def __init__(self, *args):
        super.__init__()
        self.update(*args)

    def update(self, *args, **kwargs):
        if kwargs:
            raise TypeError("update() takes no keyword arguments")

        for s in args:
            for e in s:
                 self.add(e)

    def add(self, elem):
        self[elem] = None

    def discard(self, elem):
        self.pop(elem, None)

    def __le__(self, other):
        return all(e in other for e in self)

    def __lt__(self, other):
        return self <= other and self != other

    def __ge__(self, other):
        return all(e in self for e in other)

    def __gt__(self, other):
        return self >= other and self != other

    def __repr__(self):
        return 'OrderedSet([%s])' % (', '.join(map(repr, self.keys())))

    def __str__(self):
        return '{%s}' % (', '.join(map(repr, self.keys())))

    difference = property(lambda self: self.__sub__)
    difference_update = property(lambda self: self.__isub__)
    intersection = property(lambda self: self.__and__)
    intersection_update = property(lambda self: self.__iand__)
    issubset = property(lambda self: self.__le__)
    issuperset = property(lambda self: self.__ge__)
    symmetric_difference = property(lambda self: self.__xor__)
    symmetric_difference_update = property(lambda self: self.__ixor__)
    union = property(lambda self: self.__or__)


def main():
    pass

def get_updates():
    r = requests.get(url=pastbin_scrape_url, params=pastbin_scrape_params)
    if r.status_code == 200:
        
    else:
        raise ConnectionError("Status code: {}".format(r.status_code))
    

if __name__ == "__main__":
    main()
