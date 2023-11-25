// main.go
package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var userManager UserManager

var tokenManager TokenManager

var MemManager MemoryManager

func createDocResponse(bytesWritten int) map[string]interface{} {
	data := map[string]interface{}{
		"size": bytesWritten,
	}
	return data
}

func createTokenResponse(token string) map[string]interface{} {
	data := map[string]interface{}{
		"access_token": token,
	}
	return data
}

func SignupHandler(c *gin.Context) {
	var UserAux, err = userManager.createUser(c.Request.Body)
	if err != nil {
		c.String(http.StatusConflict, err.Error())
		return
	}

	token := tokenManager.getToken(UserAux.UserName)
	response := createTokenResponse(token)
	c.IndentedJSON(http.StatusOK, response)

}

func LoginHandler(c *gin.Context) {
	var UserAux, err = userManager.logUser(c.Request.Body)
	if err != nil {
		c.String(http.StatusUnauthorized, err.Error())
		return
	}
	token := tokenManager.getToken(UserAux.UserName)
	response := createTokenResponse(token)
	c.IndentedJSON(http.StatusOK, response)
}

func VersionHandler(c *gin.Context) {
	c.String(http.StatusOK, "0.1.0")
}

func parseHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", &MissingAuthHeader{"Missing Authorization Header in the request"}
	}
	words := strings.Split(authHeader, " ")
	if len(words) != 2 {
		return "", &BadAuthHeader{"Malformed Authorization header, type is token token_id"}
	}
	return words[1], nil
}

func parseParams(params gin.Params, token string) (string, string, error) {
	var username string
	var docId string
	if len(params) != 1 {
		username = params[0].Value
		docId = params[1].Value
	} else {
		username = params[0].Value
	}
	if !userManager.UserExist(username) {
		return "", "", &UserNotExists{username}
	}
	if !tokenManager.tokenExists(token) {
		return "", "", &TokenExpired{token}
	}
	if username != tokenManager.getTokenOwner(token) {
		return "", "", &NotOwner{username, token}
	}
	return username, docId, nil
}

func parseBody(body io.ReadCloser) ([]byte, error) {
	datosJson, _ := io.ReadAll(body)
	var jsonFormat map[string]interface{}
	var dataToSave []byte

	json.Unmarshal(datosJson, &jsonFormat)
	tempData := jsonFormat["doc_content"]
	if tempData == nil {
		return dataToSave, &MissingDocContent{"Missing doc_content"}
	}

	dataToSave, _ = json.Marshal(tempData)

	return dataToSave, nil

}

func DocHandler(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	token, err := parseHeader(authHeader)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	username, docId, err := parseParams(c.Params, token)
	if err != nil {
		c.String(http.StatusUnauthorized, err.Error())
		return
	}

	if c.Request.Method == "POST" {
		bytes, err := parseBody(c.Request.Body)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		bytesWritten, err := MemManager.saveInfo(username, docId, bytes)

		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		response := createDocResponse(bytesWritten)
		c.IndentedJSON(http.StatusOK, response)
	} else if c.Request.Method == "PUT" {
		bytes, err := parseBody(c.Request.Body)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
		}

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
func AllDocsHandler(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	token, err := parseHeader(authHeader)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	username, _, err := parseParams(c.Params, token)
	if err != nil {
		c.String(http.StatusUnauthorized, err.Error())
		return
	}

	data, err := MemManager.getAllDoc(username)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.IndentedJSON(http.StatusOK, data)
}

func main() {

	r := gin.Default()
	ruta := "Usuarios.json"
	storage := "StorageDir"
	userManager.file = ruta
	MemManager.StorageDir = storage
	userManager.loadUsers()
	go tokenManager.deleteOldTokens()
	r.POST("/singup", SignupHandler)
	r.POST("/login", LoginHandler)
	r.GET("/version", VersionHandler)
	r.POST("/:username/:doc_id", DocHandler)
	r.PUT("/:username/:doc_id", DocHandler)
	r.GET("/:username/:doc_id", DocHandler)
	r.DELETE("/:username/:doc_id", DocHandler)
	r.GET("/:username/_all_docs", AllDocsHandler)
	r.Run("myserver.local:5000") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
