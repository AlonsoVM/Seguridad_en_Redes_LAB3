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
	"time"

	"github.com/gin-gonic/gin"
)

var UsersL UserList

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
	c.JSON(http.StatusOK, response)

}

func versionHandler(c *gin.Context) {
	c.String(http.StatusOK, "0.1.0")
}

func main() {

	r := gin.Default()
	ruta := "Usuarios.json"
	UsersL.file = ruta
	UsersL.loadUsers()
	r.POST("/singup", signupHandler)
	r.POST("/login", loginHandler)
	r.GET("/version", versionHandler)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
