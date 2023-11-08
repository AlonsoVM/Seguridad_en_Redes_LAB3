#!/usr/bin/env python3

import json
import requests
USERNAME = "username"
PASS = "password"
AUTH_TOKEN = "token"
ACCESS_TOKEN = "access_token"
AUTHORIZATION = "Authorization"

def version(token: str):
    dit = {}
    dit[AUTH_TOKEN] = token
    resp = requests.get("http://127.0.0.1:8080/version", headers= dit)
    print(f'{resp.content.decode()} {resp.status_code}')
def singup(username : str, password : str):
    dit = {}
    dit[USERNAME] = username
    dit[PASS] = password
    resp = requests.post("http://127.0.0.1:8080/singup", data=json.dumps(dit))
    print(resp.status_code)
    if resp.status_code == 200: print(json.loads(resp.content))
    else: print(resp.content.decode())
    return json.loads(resp.content)[ACCESS_TOKEN]

def login(username : str, password : str):
    dit = {}
    dit[USERNAME] = username
    dit[PASS] = password
    resp = requests.post("http://127.0.0.1:8080/login", data=json.dumps(dit))
    print(resp.status_code) 
    if resp.status_code == 200: print(json.loads(resp.content))
    else: print(resp.content.decode())
    return json.loads(resp.content)[ACCESS_TOKEN]

def PostDoc(token : str):
    dit = {}
    dit[AUTHORIZATION] = AUTH_TOKEN + " " + token
    ditAux = {}
    ditAux["Sal"] = 5
    ditAux["Azucar"] = 2
    ditAux["Acciones"] = ["Remover", "limpiar"]
    resp = requests.post("http://127.0.0.1:8080/Alonso3/proporciones", headers=dit, data=json.dumps(ditAux))
    print(resp.content.decode())

def GETDoc(token : str):
    dit = {}
    dit[AUTHORIZATION] = AUTH_TOKEN + " " + token
    resp = requests.get("http://127.0.0.1:8080/Alonso3/proporciones", headers=dit)
    print(resp.content.decode())
token = login("Alonso3", "qwe")
GETDoc(token)