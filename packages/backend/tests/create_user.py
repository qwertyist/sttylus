import requests
import json
from dotenv import dotenv_values

config = dotenv_values(".env")
admin_id = config["STTYLUS_ADMIN_ID"]
user_id = config["STTYLUS_USER_ID"]
tester_id = config["STTYLUS_TESTER_ID"]
api = config["STTYLUS_LOCAL_API"]
print("Testing against API:", api)
print("Admin:", admin_id)
print("User:", user_id)
print("Tester:", tester_id)

headers = { "admin": { "X-Id-Token": admin_id },
        "user" : { "X-Id-Token": user_id },
        "tester" : { "X-Id-Token": tester_id }
        }

def get_users(headers):
    return requests.get(api + "/users", headers=headers)

def get_user(headers, id): 
    return requests.get(api + "/user/" + id, headers=headers)

def create_user(headers, data):
    return requests.post(api + "/user", data=json.dumps(data), headers=headers)

def update_user(headers, data):
    return requests.put(api + "/user/" + data["id"], data=json.dumps(data), headers=headers)

print("Get initial user data")
r = get_user(headers["user"], user_id)
if r.status_code == 200:
    user_data = json.loads(r.text)
    print("User email:", user_data["email"])
else:
    print("Couldn't retrieve user data:", r.status_code)

r = get_user(headers["tester"], tester_id)
if r.status_code == 200:
    tester_data = json.loads(r.text)
    print("Tester email:", tester_data["email"])
else:
    print("Couldn't retrieve tester data:", r.status_code)

r = get_user(headers["admin"], admin_id)
if r.status_code == 200:
    admin_data = json.loads(r.text)
    print("Admin email:", admin_data["email"])
else:
    print("Couldn't retrieve admin data:", r.status_code)

r = get_user(headers["user"], admin_id)
if r.status_code == 200:
    user_data = json.loads(r.text)
    print("Admin data:", admin_data)
else:
    print("Couldn't retreive admin data:", r.status_code)
    
print("Try to get all users without providing ID")
r = get_users({})
if r.status_code == 200:
    print("OK")
else:
    print("Couldn't get all users", r.status_code)

print("Try to get all users as user")
r = get_users(headers["user"])
if r.status_code == 200:
    print("OK", len(json.loads(r.text)))
else:
    print("Couldn't get all users", r.status_code)

print("Try to get all users as admin")
r = get_users(headers["admin"])
if r.status_code == 200:
    print("OK:", len(json.loads(r.text)), "users")
else:
    print("Couldn't get all users", r.status_code)

print("Creating user without providing ID")
data = { "role":"tester", "name":"testkonto", "email":"tester@sttylus.se"}
r = create_user({}, data)
if r.status_code == 200:
    print("OK:", json.loasd(r.text)["updated"])
else:
    print("Failed:", r.status_code, r.text)

print("Creating user as admin")
r = create_user(headers["admin"], data)
if r.status_code == 401:
    print("Not allowed")
elif r.status_code == 409: 
    print("Email already in use")

print("Creating user as user")
r = create_user(headers["admin"], data)
if r.status_code == 401:
    print("Not allowed")
elif r.status_code == 409: 
    print("Email already in use")

print("Updating user without providing ID")
r = update_user({}, tester_data)
if r.status_code == 200:
    print("Updated user")
else:
    print("Failed:", r.status_code, r.text)

print("Updating user as admin")
r = update_user(headers["admin"], user_data)
if r.status_code == 200:
    print("Updated user", json.loads(r.text)["updated"])
else:
    print("Failed:", r.status_code, r.text)

print("Updating user as same user") 
r = update_user(headers["user"], user_data)
if r.status_code == 200:
    print("Updated user", r.text)
else:
    print("Failed:", r.status_code, r.text)

print("Updating user as other user") 
r = update_user(headers["tester"], user_data)
if r.status_code == 200:
    print("Updated user:", r.text)
else:
    print("Failed:", r.status_code, r.text)

print("Updating tester as tester")
r = update_user(headers["tester"], tester_data)
if r.status_code == 200:
    print("Updated user:", r.text)
else:
    print("Failed:", r.status_code, r.text)
