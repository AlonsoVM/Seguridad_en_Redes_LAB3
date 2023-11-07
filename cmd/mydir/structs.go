package main

import (
	"encoding/json"
	"fmt"
	"os"
)

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
