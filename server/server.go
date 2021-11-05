package server

import (
	"crud/database"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type user struct {
	ID   uint32 `json:"id"`
	Name string `json:"name"`
	Mail string `json:"mail"`
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	bodyRequest, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("Failed to read request body"))
		return
	}

	var user user
	if err = json.Unmarshal(bodyRequest, &user); err != nil {
		w.Write([]byte("Error on convert user"))
		return
	}

	db, err := database.OpenConnection()
	if err != nil {
		w.Write([]byte("Failed on database connect"))
		return
	}
	defer db.Close()

	statement, err := db.Prepare("insert into users (name, mail) values (?, ?)")
	if err != nil {
		w.Write([]byte("Failed on create statement"))
		return
	}
	defer statement.Close()

	insertUser, err := statement.Exec(user.Name, user.Mail)
	if err != nil {
		w.Write([]byte("Failed on exec statement"))
		return
	}

	newId, err := insertUser.LastInsertId()
	if err != nil {
		w.Write([]byte("Failed on get new ID"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Success on create user ID : %d", newId)))
}
