#!/usr/bin/env python3

import json
import requests
USERNAME = "username"
PASS = "password"
AUTH_TOKEN = "access_token"

def version():
    resp = requests.get("http://127.0.0.1:8080/version")
    print(f'{resp.content.decode()} {resp.status_code}')
def singup(username : str, password : str):
    dit = {}
    dit[USERNAME] = username
    dit[PASS] = password
    resp = requests.post("http://127.0.0.1:8080/singup", data=json.dumps(dit))
    print(resp.status_code)
    if resp.status_code == 200: print(json.loads(resp.content))
    else: print(resp.content.decode())

def login(username : str, password : str):
    dit = {}
    dit[USERNAME] = username
    dit[PASS] = password
    resp = requests.post("http://127.0.0.1:8080/login", data=json.dumps(dit))
    print(resp.status_code)
    if resp.status_code == 200: print(json.loads(resp.content))
    else: print(resp.content.decode())

version()
singup("Benito", "Contrasenna12345_")