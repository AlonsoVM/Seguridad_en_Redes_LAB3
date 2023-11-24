// main_test.go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

const serverEndPoint = "http://myserver.local:5000"
const extensionSingUp = "/singup"
const extensionLogin = "/login"
const extensionVersion = "/version"
const extensionPost = "/test/testjson"
const username = "username"
const password = "password"
const passwordTest = "test123"
const BadPasswordTest = "test1234"
const userTest = "test"
const tokenTest = "tokenTest"
const BadUserTest = "badtest"

func initialice() {
	ruta := "UsuariosTest.json"
	storage := "StorageDirTest"
	userManager.file = ruta
	MemManager.StorageDir = storage
	userManager.loadUsers()
}

func createToken() {
	var tokenToUse VolatileToken
	tokenToUse.userName = userTest
	tokenToUse.Time = time.Now()
	tokenToUse.Token = "tokenTest"
	tokenManager.VolatileTokens = append(tokenManager.VolatileTokens, tokenToUse)
}
func TestVersionHandler(t *testing.T) {
	initialice()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/version", VersionHandler)

	request, _ := http.NewRequest("GET", extensionVersion, nil)

	recorder := httptest.NewRecorder()

	r.ServeHTTP(recorder, request)

	if recorder.Code != 200 {
		t.Error("Error in the version handler")
	}
}
func TestSignupHandler(t *testing.T) {
	//initialice gin y testMode
	initialice()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/singup", SignupHandler)

	//Creating the body of the singup request
	body := make(map[string]interface{})
	body["username"] = userTest
	body["password"] = passwordTest
	bytesToSend, _ := json.Marshal(body)

	request, err := http.NewRequest("POST", extensionSingUp, bytes.NewBuffer(bytesToSend))
	if err != nil {
		fmt.Println("error creating petition")
	}

	recorder := httptest.NewRecorder()
	//Simulating petition to the server should works correct
	r.ServeHTTP(recorder, request)
	status := recorder.Code
	if status != 200 {
		t.Error("error in the singup petition")
	}

	// repeating the same request, now it shoul faild cause the user already exits
	recorder = httptest.NewRecorder()
	request, _ = http.NewRequest("POST", extensionSingUp, bytes.NewBuffer(bytesToSend))
	r.ServeHTTP(recorder, request)
	status = recorder.Code
	if status != 409 {
		t.Error("error in the singup petition")
	}

}

func TestLoginHandler(t *testing.T) {
	initialice()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/login", LoginHandler)
	body := make(map[string]interface{})
	body["username"] = userTest
	body["password"] = passwordTest
	bytesToSend, _ := json.Marshal(body)

	request, err := http.NewRequest("POST", extensionLogin, bytes.NewBuffer(bytesToSend))
	if err != nil {
		t.Fatal("error creating petition")
	}

	recorder := httptest.NewRecorder()

	r.ServeHTTP(recorder, request)
	if recorder.Code != 200 {
		t.Error("error in the login petition")
	}

}

func TestLoginHandlerBadUser(t *testing.T) {
	initialice()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/login", LoginHandler)

	body := make(map[string]interface{})
	body["username"] = BadUserTest
	body["password"] = passwordTest
	bytesToSend, _ := json.Marshal(body)
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", extensionLogin, bytes.NewBuffer(bytesToSend))
	r.ServeHTTP(recorder, request)

	if recorder.Code != 401 {
		t.Error("Error in the login ")
	}

}

func TestLoginHandlerBadPassword(t *testing.T) {
	initialice()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/login", LoginHandler)

	body := make(map[string]interface{})
	body["username"] = userTest
	body["password"] = BadPasswordTest
	bytesToSend, _ := json.Marshal(body)
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", extensionLogin, bytes.NewBuffer(bytesToSend))
	r.ServeHTTP(recorder, request)

	if recorder.Code != 401 {
		t.Error("Error in the login ")
	}
}

func TestDocHandlerPOST(t *testing.T) {
	createToken()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/:username/:doc_id", DocHandler)
	initialice()
	body := make(map[string]interface{})
	body["doc_content"] = "Prueba"
	bytesToSend, _ := json.Marshal(body)
	token := "token" + " " + tokenTest
	fmt.Println(token)
	request, _ := http.NewRequest("POST", "/test/testjson", bytes.NewBuffer(bytesToSend))
	request.Header.Set("Authorization", token)

	recorder := httptest.NewRecorder()

	r.ServeHTTP(recorder, request)

	fmt.Println(recorder.Body.String())

}
