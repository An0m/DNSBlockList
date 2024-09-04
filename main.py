from requests import get
from os import makedirs
import re
from urllib.request import urlopen

def chunkArray(xs, n):
    n = max(1, n)
    return (xs[i:i+n] for i in range(0, len(xs), n))

class ListManager():
    def __init__(self, chunkSize:int) -> None:
        self.all:set[set] = set()
        self.DOMAIN_REGEX = re.compile(r"((?=[a-z0-9-]{1,63}\.)(xn--)?[a-z0-9]+(-[a-z0-9]+)*\.)+[a-z]{2,63}")
        self.load()

    def parseList(self, data:str) -> set[str]:
        res = set()
        for line in data.replace("\t", " ").replace("\r", "").split("\n"):
            if len(line) < 4 or line[0] in "!#":
                continue

            match = self.DOMAIN_REGEX.search(line.removeprefix("0.0.0.0"))
            if not match: continue

            domain = match.group()
            parent = domain.split(".", 1)[-1]

            if parent in res:
                continue
            
            res.add(domain)
        
        if "" in res: res.remove("")
        return res

    def getList(self, url:str) -> set[str]:
        return self.parseList(get(url).text)
    
    def getDifferenceDebug(self, testList:set[str], url:str):
        diff = self.compare(testList)
        return f"{len(diff)}/{len(testList)} new domains found! ({round(len(diff)/len(testList)*100, 2)}%)\t\t\t{url}"

    def load(self) -> None:
        print("Loading lists...")

        # Get web lists
        with open("lists.txt", "r") as f:
            listUrls = f.read().split("\n")

        # Download them
        lists = dict()
        for url in listUrls:
            if url == "" or url[0] == "#": continue
            lists[url] = self.getList(url)

        # Load local list (manual.txt)
        with open("manual.txt", "r") as f:
            l = self.parseList(f.read())
        if len(l) != 0:
            lists["manual"] = l

        # Add to all
        sortedLists = dict(sorted(lists.items(), key=lambda pair: len(pair[1]), reverse=True)) # Sort lists by how many entries are in them
        for url in sortedLists:
            l = lists[url]
            print(self.getDifferenceDebug(l, url))
            self.all |= l
    
    def compare(self, domains:set[str]) -> set[str]:
        return domains.difference(self.all)
    def compareDebug(self, url:str) -> None:
        testList = self.getList(url)
        print(self.getDifferenceDebug(testList, url))

    def genChunks(self, chunkSize:int):
        makedirs("output/chunks", exist_ok=True)
        for i, chunk in enumerate(chunkArray(sorted(self.all), chunkSize)):
            self.saveList(chunk, f"chunks/{i}.csv")
    
    def saveList(self, domains:set[str], filepath="all.txt"):
        makedirs("output", exist_ok=True)
        with open("output/" + filepath, "w") as f:
            f.write("\n".join(domains))
    
    def save(self):
        self.saveList(self.all, "all.txt")
        self.genChunks()

if __name__ == "__main__":
    manager = ListManager()
    manager.save()
    
