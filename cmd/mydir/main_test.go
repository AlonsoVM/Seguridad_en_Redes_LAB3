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
const extensionSingUp = "/signup"
const extensionLogin = "/login"
const extensionVersion = "/version"
const extensionPost = "/test/testjson"
const username = "username"
const password = "password"
const passwordTest = "test123"
const BadPasswordTest = "test1234"
const userTest = "test"
const tokenTest = "tokenTest"
const badTokenTest = "BadtokenTest"
const BadUserTest = "badtest"

const stringToUpload = "\"Prueba123\""

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
	r.POST("/signup", SignupHandler)

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
		t.Error("error in the signup petition")
	}

	// repeating the same request, now it shoul faild cause the user already exits
	recorder = httptest.NewRecorder()
	request, _ = http.NewRequest("POST", extensionSingUp, bytes.NewBuffer(bytesToSend))
	r.ServeHTTP(recorder, request)
	status = recorder.Code
	if status != 409 {
		t.Error("error in the signup petition")
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

	body := map[string]interface{}{
		"doc_content": "Prueba123",
	}
	bytesToSend, _ := json.Marshal(body)
	request, _ := http.NewRequest("POST", "/test/testjson", bytes.NewBuffer(bytesToSend))

	token := "token" + " " + tokenTest
	request.Header.Set("Authorization", token)

	recorder := httptest.NewRecorder()

	r.ServeHTTP(recorder, request)

	if recorder.Code != 200 {
		t.Error("Error in POST document")
	}

}
func TestDocHandlerPUT(t *testing.T) {
	createToken()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.PUT("/:username/:doc_id", DocHandler)
	initialice()

	body := map[string]interface{}{
		"doc_content": "Prueba123",
	}
	bytesToSend, _ := json.Marshal(body)
	request, _ := http.NewRequest("PUT", "/test/testjson", bytes.NewBuffer(bytesToSend))

	token := "token" + " " + tokenTest
	request.Header.Set("Authorization", token)

	recorder := httptest.NewRecorder()

	r.ServeHTTP(recorder, request)

	if recorder.Code != 200 {
		t.Error("Error in POST document")
	}

}
func TestDocHandlerGET(t *testing.T) {
	createToken()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/:username/:doc_id", DocHandler)
	initialice()

	request, _ := http.NewRequest("GET", "/test/testjson", nil)
	token := "token" + " " + tokenTest
	request.Header.Set("Authorization", token)

	recorder := httptest.NewRecorder()

	r.ServeHTTP(recorder, request)
	if recorder.Code != 200 {
		t.Error("Error in GET document")
	}
}

func TestDocHandlerPOSTRepeated(t *testing.T) {
	createToken()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/:username/:doc_id", DocHandler)
	initialice()

	body := map[string]interface{}{
		"doc_content": "Prueba123",
	}
	bytesToSend, _ := json.Marshal(body)
	request, _ := http.NewRequest("POST", "/test/testjson", bytes.NewBuffer(bytesToSend))
	token := "token" + " " + tokenTest
	request.Header.Set("Authorization", token)

	recorder := httptest.NewRecorder()

	r.ServeHTTP(recorder, request)

	if recorder.Code != 500 {
		t.Error("Error in POST  repeated document")
	}
}

func TestDocHandlerPOSTBadToken(t *testing.T) {
	createToken()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/:username/:doc_id", DocHandler)
	initialice()

	body := map[string]interface{}{
		"doc_content": "Prueba123",
	}
	bytesToSend, _ := json.Marshal(body)
	request, _ := http.NewRequest("POST", "/test/testjson", bytes.NewBuffer(bytesToSend))
	token := "token" + " " + badTokenTest
	request.Header.Set("Authorization", token)

	recorder := httptest.NewRecorder()

	r.ServeHTTP(recorder, request)

	if recorder.Code != 401 {
		t.Error("Error in POST document with bad token")
	}
}

func TestDocHandlerGETBadToken(t *testing.T) {
	createToken()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/:username/:doc_id", DocHandler)
	initialice()

	request, _ := http.NewRequest("GET", "/test/testjson", nil)
	token := "token" + " " + badTokenTest
	request.Header.Set("Authorization", token)

	recorder := httptest.NewRecorder()

	r.ServeHTTP(recorder, request)
	if recorder.Code != 401 {
		t.Error("Error in GET document")
	}
}

func TestDocHandlerGETALLDOCS(t *testing.T) {
	createToken()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/:username/_all_docs", AllDocsHandler)
	initialice()

	request, _ := http.NewRequest("GET", "/test/_all_docs", nil)
	token := "token" + " " + tokenTest
	request.Header.Set("Authorization", token)

	recorder := httptest.NewRecorder()
	fmt.Println(recorder.Body.String())
	r.ServeHTTP(recorder, request)
	if recorder.Code != 200 {
		t.Error("Error in GET all document")
	}
}

func TestDocHandlerDELETE(t *testing.T) {
	createToken()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.DELETE("/:username/:doc_id", DocHandler)
	initialice()

	request, _ := http.NewRequest("DELETE", "/test/testjson", nil)
	token := "token" + " " + tokenTest
	request.Header.Set("Authorization", token)

	recorder := httptest.NewRecorder()

	r.ServeHTTP(recorder, request)
	if recorder.Code != 200 {
		t.Error("Error in DELETE document")
	}

}
