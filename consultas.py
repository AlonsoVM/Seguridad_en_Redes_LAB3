#!/usr/bin/env python3

import json
import requests
USERNAME = "username"
PASS = "password"
AUTH_TOKEN = "token"
ACCESS_TOKEN = "access_token"
AUTHORIZATION = "Authorization"
DOCUMENT = "doc_content"
data = {
"capacidad_bateria": 20,
"dispositivos_iot": [
{
"nombre": "Router",
"consumo": 0.6
},
{
"nombre": "Camara",
"consumo": 0.5
},
{
"nombre": "Termostato",
"consumo": 0.3
},
{
"nombre": "Sensor",
"consumo": 0.1
},
{
"nombre": "Sensor",
"consumo": 0.1
},
{
"nombre": "Sensor",
"consumo": 0.1
},
{
"nombre": "Sensor",
"consumo": 0.1
}
]
}

def version(token: str):
    dit = {}
    dit[AUTH_TOKEN] = token
    resp = requests.get("http://myserver.local:5000/version", headers= dit)
    print(f'{resp.content.decode()} {resp.status_code}')

def singup(username : str, password : str):
    dit = {}
    dit[USERNAME] = username
    dit[PASS] = password
    resp = requests.post("http://myserver.local:5000/signup", data=json.dumps(dit))
    print(resp.status_code)
    if resp.status_code == 200: print(json.loads(resp.content))
    else: print(resp.content.decode())
    return json.loads(resp.content)[ACCESS_TOKEN]

def login(username : str, password : str):
    dit = {}
    dit[USERNAME] = username
    dit[PASS] = password
    resp = requests.post("http://myserver.local:5000/login", data=json.dumps(dit))
    print(resp.status_code) 
    if resp.status_code == 200: print(json.loads(resp.content))
    else: print(resp.content.decode())
    return json.loads(resp.content)[ACCESS_TOKEN]

def PostDoc(token : str):
    dit = {}
    dit[AUTHORIZATION] = AUTH_TOKEN + " " + token
    ditAux = {}
    ditAux["Perejil"] = 5
    ditAux["Tomillo"] = 2
    ditAux["Acciones"] = ["Remover", "limpiar"]
    ditAux2 = {}
    f = open("iot.json", "rb")
    data = f.read()
    jsontodata = json.loads(data)
    dataIndented = json.dumps(jsontodata, indent=2)
    ditAux2[DOCUMENT] = jsontodata
    resp = requests.post("http://myserver.local:5000/AlonsoVilla2233/iot", headers=dit, data=json.dumps(ditAux2, indent=2))
    print(resp.content.decode())

def PutDoc(token : str):
    dit = {}
    dit[AUTHORIZATION] = AUTH_TOKEN + " " + token
    ditAux = {}
    ditAux["sal"] = 5
    ditAux["Azucar"] = 2
    ditAux2 = {}
    f = open("iot.json", "rb")
    data = f.read()
    jsontodata = json.loads(data)
    dataIndented = json.dumps(jsontodata, indent=2)
    ditAux2[DOCUMENT] = jsontodata
    len(json.dumps(ditAux2))

    resp = requests.put("http://myserver.local:5000/AlonsoVilla2233/iot", headers=dit, data=json.dumps(ditAux2, indent=2))
    print(resp.content.decode())

def GETDoc(token : str):
    dit = {}
    dit[AUTHORIZATION] = AUTH_TOKEN + " " + token
    resp = requests.get("http://myserver.local:5000/AlonsoVilla2233/iot", headers=dit)
    print(resp.content.decode())

def DeleteDoc(token : str):
    dit = {}
    dit[AUTHORIZATION] = AUTH_TOKEN + " " + token
    resp = requests.delete("http://myserver.local:5000/AlonsoVilla2233/iot", headers=dit)
    print(resp.content.decode())

def GETALLDoc(token : str):
    dit = {}
    dit[AUTHORIZATION] = AUTH_TOKEN + " " + token
    resp = requests.get("http://myserver.local:5000/AlonsoVilla2233/_all_docs", headers=dit)
    print(resp.content.decode())
token = login("AlonsoVilla2233", "qwe")
DeleteDoc(token)