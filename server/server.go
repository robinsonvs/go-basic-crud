package server

import (
	"crud/database"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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
	params := mux.Vars(r)

	ID, err := strconv.ParseUint(params["id"], 10, 32)
	if err != nil {
		w.Write([]byte("Failed on convert params to ID"))
		return
	}

	db, err := database.OpenConnection()
	if err != nil {
		w.Write([]byte("Failed on database connection"))
		return
	}
	defer db.Close()

	lineReturn, err := db.Query("select * from users where id = ?", ID)
	if err != nil {
		w.Write([]byte("Failed on find user ID"))
		return
	}
	defer lineReturn.Close()

	var user user
	if lineReturn.Next() {
		if err := lineReturn.Scan(&user.ID, &user.Name, &user.Mail); err != nil {
			w.Write([]byte("Failed on user scan"))
			return
		}
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		w.Write([]byte("Failed on convert user to json"))
		return
	}
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	ID, err := strconv.ParseUint(params["id"], 10, 32)
	if err != nil {
		w.Write([]byte("Fail on convert parameter to Integer"))
		return
	}

	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("Fail on read body"))
		return
	}

	var user user
	if err := json.Unmarshal(requestBody, &user); err != nil {
		w.Write([]byte("Fail on convert user to struc"))
		return
	}

	db, err := database.OpenConnection()
	if err != nil {
		w.Write([]byte("Failed on database connect"))
		return
	}
	defer db.Close()

	statement, err := db.Prepare("update users set name = ?, mail = ? where id = ?")
	if err != nil {
		w.Write([]byte("Fail on statement create"))
		return
	}
	defer statement.Close()

	if _, err := statement.Exec(user.Name, user.Mail, ID); err != nil {
		w.Write([]byte("Fail on update user"))
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
