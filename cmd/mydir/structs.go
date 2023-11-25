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
	"os"
	"strings"
	"sync"
	"time"
)

type MemoryManager struct {
	StorageDir string
}

func (Mem *MemoryManager) getAllDoc(username string) (map[string]map[string]interface{}, error) {
	data := make(map[string]map[string]interface{})
	pathDir := fmt.Sprintf("%s/%s", Mem.StorageDir, username)
	dir, err := os.Open(pathDir)
	if err != nil {
		fmt.Println("Error opening the directory : ", pathDir)
		return data, err
	}
	defer dir.Close()

	files, err := dir.Readdirnames(0)
	if err != nil {
		fmt.Println("Error reading the entries of the directory : ", pathDir)
	}

	for i, file := range files {
		jsonData, _ := Mem.getInfo(username, file)
		data["id"+fmt.Sprint(i)] = jsonData
	}
	return data, nil
}

func (Mem *MemoryManager) saveInfo(username string, filename string, data []byte) (int, error) {
	pathDir := fmt.Sprintf("%s/%s/", Mem.StorageDir, username)
	err := os.MkdirAll(pathDir, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating directorys", err)
		return 0, err
	}
	pathfile := fmt.Sprintf("%s%s%s", pathDir, filename, ".json")
	_, err = os.Stat(pathfile)
	if err == nil {
		fmt.Println("The file already exits")
		return 0, &FileExits{filename}
	}
	archivo, err := os.OpenFile(pathfile, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Print("Error opening the file", filename)
		return 0, err
	}
	defer archivo.Close()
	bytesWritten, err := archivo.Write(data)
	if err != nil {
		fmt.Println("Error writting in the file :", filename)
		return 0, err
	}
	return bytesWritten, nil
}

func (Mem *MemoryManager) updateInfo(username string, filename string, data []byte) (int, error) {
	pathfile := fmt.Sprintf("%s/%s/%s%s", Mem.StorageDir, username, filename, ".json")
	archivo, err := os.OpenFile(pathfile, os.O_WRONLY, 0666)
	if err != nil {
		fmt.Print("Error opening the file", filename)
		return 0, err
	}
	defer archivo.Close()
	bytesWritten, err := archivo.Write(data)
	if err != nil {
		fmt.Println("Error writting in the file :", filename)
		return 0, err
	}
	return bytesWritten, nil
}

func (Mem *MemoryManager) removeInfo(username string, filename string) error {
	pathfile := fmt.Sprintf("%s/%s/%s%s", Mem.StorageDir, username, filename, ".json")
	err := os.Remove(pathfile)
	if err != nil {
		fmt.Print("Error removing", filename)
		return err
	}
	return nil
}

func (Mem *MemoryManager) getInfo(username string, filename string) (map[string]interface{}, error) {
	var pathfile string
	if strings.HasSuffix(filename, ".json") {
		pathfile = fmt.Sprintf("%s/%s/%s", Mem.StorageDir, username, filename)
	} else {
		pathfile = fmt.Sprintf("%s/%s/%s%s", Mem.StorageDir, username, filename, ".json")
	}
	archivo, err := os.OpenFile(pathfile, os.O_RDONLY, 0666)
	var jsonData map[string]interface{}
	if err != nil {
		fmt.Print("Error opening the file", filename)
		return jsonData, err
	}
	defer archivo.Close()
	decoder := json.NewDecoder(archivo)
	decoder.Decode(&jsonData)
	fmt.Println(jsonData)
	return jsonData, nil
}

type VolatileToken struct {
	Token    string
	userName string
	Time     time.Time
}

type TokenManager struct {
	VolatileTokens []VolatileToken
	mutex          sync.Mutex
}

func (tokenList *TokenManager) createToken() string {
	rand.NewSource(time.Now().UnixMilli())
	random := rand.Int63()
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, uint64(random))
	token := base64.StdEncoding.EncodeToString([]byte(bytes))
	return token
}

func (tokenList *TokenManager) getToken(username string) string {
	token := tokenList.createToken()
	tokenList.removeSingleToken(username)
	tokenList.saveToken(token, username)
	return token

}

func (tokenList *TokenManager) removeSingleToken(username string) {
	tokenList.mutex.Lock()
	for i, token := range tokenList.VolatileTokens {
		if token.userName == username {
			tokenList.VolatileTokens = append(tokenList.VolatileTokens[:i], tokenList.VolatileTokens[i+1:]...)
			fmt.Println("Removing token : ", token)
		}
	}
	tokenList.mutex.Unlock()
}

func (tokenList *TokenManager) saveToken(tempToken string, userName string) {
	var token VolatileToken
	token.Token = tempToken
	token.userName = userName
	token.Time = time.Now()
	tokenList.mutex.Lock()
	tokenList.VolatileTokens = append(tokenList.VolatileTokens, token)
	tokenList.mutex.Unlock()
}

func (tokenList *TokenManager) deleteOldTokens() {
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

func (tokenList *TokenManager) tokenExists(id string) bool {
	for _, token := range tokenList.VolatileTokens {
		if token.Token == id {
			return true
		}
	}
	return false
}

func (tokenList *TokenManager) getTokenOwner(id string) string {
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

type UserManager struct {
	Users []User `json:"users"`
	file  string
}

func (UserL *UserManager) createUser(body io.ReadCloser) (User, error) {
	rand.NewSource(time.Now().UnixMilli())
	var datosJson, _ = io.ReadAll(body)
	var UserAux User

	json.Unmarshal(datosJson, &UserAux)
	if UserL.UserExist(UserAux.UserName) {
		return UserAux, &UserExists{"The user already exits"}
	}

	UserAux.Salt = fmt.Sprintf("%d", rand.Int())
	UserAux.Password = createHashedPassword(UserAux.Salt, UserAux.Password)
	UserL.saveUsers(UserAux)
	return UserAux, nil
}

func createHashedPassword(salt string, pass string) string {
	tempPassword := fmt.Sprintf("%s%s", salt, pass)
	hash := sha256.New()
	return hex.EncodeToString(hash.Sum([]byte(tempPassword)))
}

func (UserL *UserManager) logUser(body io.ReadCloser) (User, error) {
	var datosJson, _ = io.ReadAll(body)
	var UserAux User
	json.Unmarshal(datosJson, &UserAux)
	if !UserL.UserExist(UserAux.UserName) {
		return UserAux, &UserNotExists{UserAux.UserName}
	}

	Password := UserL.getUserPassword(UserAux.UserName)
	Salt := UserL.getUserSalt(UserAux.UserName)
	tempPass := createHashedPassword(Salt, UserAux.Password)
	if Password != tempPass {
		return UserAux, &InvalidPassword{"The password provided is invalid"}
	}

	return UserAux, nil
}

func (UserL *UserManager) loadUsers() {
	archivo, err := os.OpenFile(UserL.file, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Print("Error opening the file", UserL.file)
	}
	defer archivo.Close()
	decoder := json.NewDecoder(archivo)
	decoder.Decode(&UserL)
}
func (UserL *UserManager) saveUsers(NewUser User) {
	archivo, err := os.OpenFile(UserL.file, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Print("Error opening the file", UserL.file)
	}
	defer archivo.Close()
	encoder := json.NewEncoder(archivo)
	UserL.Users = append(UserL.Users, NewUser)
	encoder.SetIndent("", " ")
	encoder.Encode(&UserL)
}
func (UserL *UserManager) getUserPassword(UserName string) string {
	for _, userAux := range UserL.Users {
		if userAux.UserName == UserName {
			return userAux.Password
		}
	}
	return ""
}

func (UserL *UserManager) getUserSalt(UserName string) string {
	for _, userAux := range UserL.Users {
		if userAux.UserName == UserName {
			return userAux.Salt
		}
	}
	return ""
}

func (UserL *UserManager) UserExist(UserName string) bool {
	for i := 0; i < len(UserL.Users); i++ {
		if UserName == UserL.Users[i].UserName {
			return true
		}
	}
	return false
}
