package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

type VolatileToken struct {
	Token string
	Time  time.Time
}

type VolatileTokenList struct {
	VolatileTokens []VolatileToken
	mutex          sync.Mutex
}

func (tokenList *VolatileTokenList) saveToken(tempToken string) {
	fmt.Println("Saving token")
	var token VolatileToken
	token.Token = tempToken
	token.Time = time.Now()
	tokenList.mutex.Lock()
	tokenList.VolatileTokens = append(tokenList.VolatileTokens, token)
	tokenList.mutex.Unlock()
}

func (tokenList *VolatileTokenList) deleteOldTokens() {
	for true {
		tokenList.mutex.Lock()
		for i, token := range tokenList.VolatileTokens {
			if time.Now().Sub(token.Time).Seconds() > 30 {
				tokenList.VolatileTokens = append(tokenList.VolatileTokens[:i], tokenList.VolatileTokens[i+1:]...)
				fmt.Println("Removing token : ", token.Token)
			}
		}
		tokenList.mutex.Unlock()
		time.Sleep(8 * time.Second)
		fmt.Println(tokenList.VolatileTokens)
	}
}

type User struct {
	UserName string `json:"username"`
	Password string `json:"password"`
	Salt     string `json:"salt"`
}

type UserList struct {
	Users []User `json:"users"`
	file  string
}

func (UserL *UserList) loadUsers() {
	archivo, err := os.OpenFile(UserL.file, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Print("Error al abrir el archivo")
	}
	defer archivo.Close()
	decoder := json.NewDecoder(archivo)
	decoder.Decode(&UserL)
}
func (UserL *UserList) saveUsers(NewUser User) {
	archivo, err := os.OpenFile(UserL.file, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Print("Error al abrir el archivo")
	}
	defer archivo.Close()
	encoder := json.NewEncoder(archivo)
	UsersL.Users = append(UsersL.Users, NewUser)
	encoder.SetIndent("", " ")
	encoder.Encode(&UsersL)
}
func (UserL *UserList) getUserPassword(UserName string) string {
	for _, userAux := range UsersL.Users {
		if userAux.UserName == UserName {
			return userAux.Password
		}
	}
	return ""
}

func (UserL *UserList) getUserSalt(UserName string) string {
	for _, userAux := range UsersL.Users {
		if userAux.UserName == UserName {
			return userAux.Salt
		}
	}
	return ""
}

func (UserL *UserList) UserExist(UserName string) bool {
	for i := 0; i < len(UsersL.Users); i++ {
		if UserName == UsersL.Users[i].UserName {
			return true
		}
	}
	return false
}
