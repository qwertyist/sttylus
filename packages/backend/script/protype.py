#!/usr/bin/python
import os
import sys
import json
import contextlib

def isLatin(s):
    with contextlib.suppress(UnicodeEncodeError):
        try:
            test = bytearray(s, "latin-1")
        except UnicodeDecodeError:
            return False
        else:
            return True

if len(sys.argv) < 2:
    print("No list provided, reading 'list.json'")
    with open("list.json") as f:
        abbs = json.loads(f.read())["abbs"]
    output = "tmp"
    listName = "test"
else:
    print("Creating ProType list for file '",sys.argv[1],"'")
    with open(sys.argv[1]) as f:
        abbs = json.loads(f.read())
    output = os.path.dirname(sys.argv[1])
    print("Output directory:", output)
    listName = sys.argv[2]

encoding = "latin-1"
wordlist = "wordlist.dat"
settings = "settings.dat"

export = []
for a in abbs:
    if isLatin(a["abb"]+a["word"]):
        abb = { 
                "aLen": len(a["abb"]),
                "abb": a["abb"], 
                "wLen": len(a["word"]),
                "word": a["word"]
            }
        export.append(abb)
n = len(export)
print("loaded", len(abbs) ,"abbs and omitted", len(abbs)-len(export), "rows")
with open(output + "/settings.dat", "wb") as f:
    f.write(bytearray("Standard", encoding))

os.mkdir(output + "/" + listName)
with open(output + "/" + listName + "/wordlist.dat", "wb") as f:
    f.write(bytes(n.to_bytes(2,"little")))
    for entry in export:
        try:
            f.write(bytes([entry["aLen"]]))
            f.write(bytearray(entry["abb"], encoding))
            f.write(bytes([entry["wLen"]]))
            f.write(bytearray(entry["word"], encoding))
        except UnicodeEncodeError:
            print("failed at:", entry)
            continue



