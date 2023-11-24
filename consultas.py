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
    resp = requests.post("http://myserver.local:5000/singup", data=json.dumps(dit))
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

    ditAux2[DOCUMENT] = ditAux
    resp = requests.post("http://myserver.local:5000/AlonsoVilla2233/test1", headers=dit, data=json.dumps(ditAux2))
    print(resp.content.decode())

def PutDoc(token : str):
    dit = {}
    dit[AUTHORIZATION] = AUTH_TOKEN + " " + token
    ditAux = {}
    ditAux["Sal"] = 5
    ditAux["Azucar"] = 2
    ditAux["Soja"] = "medio vaso"
    ditAux["Acciones"] = ["Remover", "limpiar", "colocar"]
    ditAux2 = {}
    ditAux2[DOCUMENT] = ditAux
    resp = requests.put("http://myserver.local:5000/Alonso3/test", headers=dit, data=json.dumps(ditAux2))
    print(resp.content.decode())

def GETDoc(token : str):
    dit = {}
    dit[AUTHORIZATION] = AUTH_TOKEN + " " + token
    resp = requests.get("http://myserver.local:5000/Alonso3/test2", headers=dit)
    print(resp.content.decode())

def DeleteDoc(token : str):
    dit = {}
    dit[AUTHORIZATION] = AUTH_TOKEN + " " + token
    resp = requests.delete("http://myserver.local:5000/Alonso3/p1", headers=dit)
    print(resp.content.decode())

def GETALLDoc(token : str):
    dit = {}
    dit[AUTHORIZATION] = AUTH_TOKEN + " " + token
    resp = requests.get("http://myserver.local:5000/AlonsoVilla2233/_all_docs", headers=dit)
    print(resp.content.decode())
token = login("AlonsoVilla2233", "qwe")
GETALLDoc(token)