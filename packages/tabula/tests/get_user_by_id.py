import sys
import json
import requests

r = {}
want = json.loads(sys.argv[2])

try: 
    r = requests.get("http://localhost:8080/user/" + want["id"])
except:
    print(want)
got = json.loads(r.text)
if want == got:
    exit(0)
else:
    print(want)
    exit(1)
