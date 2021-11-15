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

func FindUsers(w http.ResponseWriter, r *http.Request) {
	db, err := database.OpenConnection()
	if err != nil {
		w.Write([]byte("Failed on database connect"))
		return
	}
	defer db.Close()

	listReturn, err := db.Query("select * from users")
	if err != nil {
		w.Write([]byte("Failed on find users"))
		return
	}
	defer listReturn.Close()

	var users []user
	for listReturn.Next() {
		var user user

		if err := listReturn.Scan(&user.ID, &user.Name, &user.Mail); err != nil {
			w.Write([]byte("Failed on user scan"))
			return
		}

		users = append(users, user)
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(users); err != nil {
		w.Write([]byte("Failed on convert users to json"))
		return
	}

}

func FindUser(w http.ResponseWriter, r *http.Request) {

}
