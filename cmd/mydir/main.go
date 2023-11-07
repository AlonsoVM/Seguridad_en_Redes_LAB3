package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

var UsersL UserList

type User struct {
	UserName string `json:"user"`
	Password string `json:"pass"`
	Salt     string `json:"salt"`
}

type UserList struct {
	Users []User `json:"users"`
}

func (UserL *UserList) loadUsers2(jsonFile string) {
	archivo, err := os.OpenFile(jsonFile, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Print("Error al abrir el archivo")
	}
	defer archivo.Close()
	decoder := json.NewDecoder(archivo)
	decoder.Decode(&UserL)
}

func loadUsers(jsonFile string) {

	archivo, err := os.OpenFile(jsonFile, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Print("Error al abrir el archivo")
	}
	defer archivo.Close()
	decoder := json.NewDecoder(archivo)
	decoder.Decode(&UsersL)
}

func saveUsers(jsonFile string, NewUser User) {
	archivo, err := os.OpenFile(jsonFile, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Print("Error al abrir el archivo")
	}
	defer archivo.Close()
	encoder := json.NewEncoder(archivo)
	UsersL.Users = append(UsersL.Users, NewUser)
	encoder.SetIndent("", " ")
	encoder.Encode(&UsersL)
}

func getUserPassword(UserName string) string {
	for _, userAux := range UsersL.Users {
		if userAux.UserName == UserName {
			return userAux.Password
		}
	}
	return ""
}

func getUserSalt(UserName string) string {
	for _, userAux := range UsersL.Users {
		if userAux.UserName == UserName {
			return userAux.Salt
		}
	}
	return ""
}

func UserExist(UserName string) bool {
	for i := 0; i < len(UsersL.Users); i++ {
		if UserName == UsersL.Users[i].UserName {
			return true
		}
	}
	return false
}

func createHashedPassword(salt string, pass string) string {
	tempPassword := fmt.Sprintf("%s%s", salt, pass)
	hash := sha256.New()
	return hex.EncodeToString(hash.Sum([]byte(tempPassword)))
}

func signupHandler(c *gin.Context) {
	rand.NewSource(time.Now().UnixMilli())
	var datosJson, _ = io.ReadAll(c.Request.Body)
	var UserAux User
	json.Unmarshal(datosJson, &UserAux)
	if UserExist(UserAux.UserName) {
		c.String(http.StatusConflict, "User exits")
		return
	}
	UserAux.Salt = fmt.Sprintf("%d", rand.Int())
	UserAux.Password = createHashedPassword(UserAux.Salt, UserAux.Password)
	saveUsers("Usuarios.json", UserAux)
	c.String(http.StatusOK, "User Added")

}

func loginHandler(c *gin.Context) {
	var datosJson, _ = io.ReadAll(c.Request.Body)
	var requestInfo User
	json.Unmarshal(datosJson, &requestInfo)
	Password := getUserPassword(requestInfo.UserName)
	Salt := getUserSalt(requestInfo.UserName)
	tempPass := createHashedPassword(Salt, requestInfo.Password)
	if Password != tempPass {
		c.JSON(http.StatusBadRequest, "Failed login bad password")
		return
	}
	c.JSON(http.StatusOK, "Loged")

}

func main() {

	r := gin.Default()
	ruta := "Usuarios.json"
	UsersL.loadUsers2(ruta)
	fmt.Print(UsersL.Users[1].UserName)
	r.POST("/singup", signupHandler)
	r.POST("/login", loginHandler)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
