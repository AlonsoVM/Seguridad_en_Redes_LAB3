package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

type MemoryManager struct {
	StorageDir string
}

func (Mem *MemoryManager) saveInfo(username string, filename string, data []byte) int {
	pathDir := fmt.Sprintf("%s/%s/", Mem.StorageDir, username)
	err := os.MkdirAll(pathDir, 0666)
	if err != nil {
		fmt.Println("Error creating directorys", err)
	}
	pathfile := fmt.Sprintf("%s%s", pathDir, filename)
	fmt.Println(pathfile)
	archivo, err := os.OpenFile(pathfile, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Print("Error opening the file", filename)
	}
	bytesWritten, err := archivo.Write(data)
	if err != nil {
		fmt.Println("Error writting in the file :", filename)
	}
	return bytesWritten
}

func (Mem *MemoryManager) getInfo(username string, filename string) interface{} {
	pathfile := fmt.Sprintf("%s/%s/%s", Mem.StorageDir, username, filename)
	archivo, err := os.OpenFile(pathfile, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Print("Error opening the file", filename)
	}
	decoder := json.NewDecoder(archivo)
	var jsonData interface{}
	decoder.Decode(&jsonData)
	return jsonData
}

type VolatileToken struct {
	Token    string
	userName string
	Time     time.Time
}

type VolatileTokenList struct {
	VolatileTokens []VolatileToken
	mutex          sync.Mutex
}

func (tokenList *VolatileTokenList) saveToken(tempToken string, userName string) {
	var token VolatileToken
	token.Token = tempToken
	token.userName = userName
	token.Time = time.Now()
	tokenList.mutex.Lock()
	tokenList.VolatileTokens = append(tokenList.VolatileTokens, token)
	tokenList.mutex.Unlock()
}

func (tokenList *VolatileTokenList) deleteOldTokens() {
	for true {
		tokenList.mutex.Lock()
		for i, token := range tokenList.VolatileTokens {
			if time.Now().Sub(token.Time).Seconds() > 120 {
				tokenList.VolatileTokens = append(tokenList.VolatileTokens[:i], tokenList.VolatileTokens[i+1:]...)
				fmt.Println("Removing token : ", token.Token)
			}
		}
		tokenList.mutex.Unlock()
		time.Sleep(8 * time.Second)
	}
}

func (tokenList *VolatileTokenList) tokenExists(id string) bool {
	for _, token := range tokenList.VolatileTokens {
		if token.Token == id {
			return true
		}
	}
	return false
}

func (tokenList *VolatileTokenList) getTokenOwner(id string) string {
	for _, token := range tokenList.VolatileTokens {
		if token.Token == id {
			return token.userName
		}
	}
	return ""
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
		fmt.Print("Error opening the file", UserL.file)
	}
	defer archivo.Close()
	decoder := json.NewDecoder(archivo)
	decoder.Decode(&UserL)
}
func (UserL *UserList) saveUsers(NewUser User) {
	archivo, err := os.OpenFile(UserL.file, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Print("Error opening the file", UserL.file)
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
