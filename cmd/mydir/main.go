package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var UsersL UserList

var TokenL VolatileTokenList

var MemManager MemoryManager

func createHashedPassword(salt string, pass string) string {
	tempPassword := fmt.Sprintf("%s%s", salt, pass)
	hash := sha256.New()
	return hex.EncodeToString(hash.Sum([]byte(tempPassword)))
}

func createToken() string {
	rand.NewSource(time.Now().UnixMilli())
	random := rand.Int63()
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, uint64(random))
	token := base64.StdEncoding.EncodeToString([]byte(bytes))
	return token
}

func createTokenResponse(token string) map[string]interface{} {
	data := map[string]interface{}{
		"access_token": token,
	}
	//jsonResponse, _ := json.Marshal(data)
	return data
}

func createDocResponse(bytesWritten int) map[string]interface{} {
	data := map[string]interface{}{
		"size": bytesWritten,
	}
	return data
}

func signupHandler(c *gin.Context) {
	rand.NewSource(time.Now().UnixMilli())
	var datosJson, _ = io.ReadAll(c.Request.Body)
	var UserAux User
	json.Unmarshal(datosJson, &UserAux)
	if UsersL.UserExist(UserAux.UserName) {
		c.String(http.StatusUnauthorized, "User exits")
		return
	}
	UserAux.Salt = fmt.Sprintf("%d", rand.Int())
	UserAux.Password = createHashedPassword(UserAux.Salt, UserAux.Password)
	UsersL.saveUsers(UserAux)
	token := createToken()
	response := createTokenResponse(token)
	TokenL.saveToken(token, UserAux.UserName)
	c.IndentedJSON(http.StatusOK, response)

}

func loginHandler(c *gin.Context) {
	var datosJson, _ = io.ReadAll(c.Request.Body)
	var requestInfo User
	json.Unmarshal(datosJson, &requestInfo)
	if !UsersL.UserExist(requestInfo.UserName) {
		c.String(http.StatusBadRequest, "User not exist in the system")
		return
	}
	Password := UsersL.getUserPassword(requestInfo.UserName)
	Salt := UsersL.getUserSalt(requestInfo.UserName)
	tempPass := createHashedPassword(Salt, requestInfo.Password)
	if Password != tempPass {
		c.String(http.StatusUnauthorized, "Failed login bad password")
		return
	}
	token := createToken()
	response := createTokenResponse(token)
	TokenL.saveToken(token, requestInfo.UserName)
	c.JSON(http.StatusOK, response)

}

func versionHandler(c *gin.Context) {
	c.String(http.StatusOK, "0.1.0")
}

func DocHandler(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	words := strings.Split(authHeader, " ")
	token := words[1]
	username := c.Param("username")
	docId := c.Param("doc_id")
	if !UsersL.UserExist(username) {
		c.String(http.StatusUnauthorized, "User not registered in the system")
		return
	}
	if !TokenL.tokenExists(token) {
		c.String(http.StatusUnauthorized, "Your token has expired, refresh it in /login ")
		return
	}
	if username != TokenL.getTokenOwner(token) {
		c.String(http.StatusUnauthorized, "The username given is not the owner of the token")
		return
	}
	if c.Request.Method == "POST" {
		var datosJson, _ = io.ReadAll(c.Request.Body)
		var jsonFormat map[string]interface{}

		json.Unmarshal(datosJson, &jsonFormat)
		dataToSave := jsonFormat["doc_content"]
		bytes, _ := json.Marshal(dataToSave)
		bytesWritten, err := MemManager.saveInfo(username, docId, bytes)

		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		response := createDocResponse(bytesWritten)
		c.IndentedJSON(http.StatusOK, response)
	} else if c.Request.Method == "PUT" {
		var datosJson, _ = io.ReadAll(c.Request.Body)
		var jsonFormat map[string]interface{}

		json.Unmarshal(datosJson, &jsonFormat)
		dataToSave := jsonFormat["doc_content"]
		bytes, _ := json.Marshal(dataToSave)
		bytesWritten, err := MemManager.updateInfo(username, docId, bytes)

		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		response := createDocResponse(bytesWritten)
		c.IndentedJSON(http.StatusOK, response)
	} else if c.Request.Method == "GET" {
		data, err := MemManager.getInfo(username, docId)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		c.IndentedJSON(http.StatusOK, data)
	} else if c.Request.Method == "DELETE" {
		err := MemManager.removeInfo(username, docId)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		c.IndentedJSON(http.StatusOK, "{}")
	}

}
func main() {

	r := gin.Default()
	ruta := "Usuarios.json"
	storage := "StorageDir"
	UsersL.file = ruta
	MemManager.StorageDir = storage
	UsersL.loadUsers()
	go TokenL.deleteOldTokens()
	r.POST("/singup", signupHandler)
	r.POST("/login", loginHandler)
	r.GET("/version", versionHandler)
	r.POST("/:username/:doc_id", DocHandler)
	r.PUT("/:username/:doc_id", DocHandler)
	r.GET("/:username/:doc_id", DocHandler)
	r.DELETE("/:username/:doc_id", DocHandler)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
