package main

import (
	"BLABLA/config"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB

type user struct {
	UserID  int
	Balance float32
}

func getBalance(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		panic(err)
	}
	user := newUser(userID)

	data, err := json.Marshal(user)

	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)

}

func transfer(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		panic(err)
	}

	senderID := int(data["sender_id"].(float64))
	receiverID := int(data["receiver_id"].(float64))
	amount := data["amount"].(float32)

	sender := newUser(senderID)
	if sender.Balance < amount {
		fmt.Println("insufficient fund")
		return
	}

	receiver := newUser(receiverID)
	sender.Balance -= amount
	receiver.Balance += amount

	updateUser(sender)
	updateUser(receiver)

	w.WriteHeader(http.StatusOK)

}

func updateUser(user *user) {
	_, err := db.Exec("UPDATE BILLING.BALANCE SET BALANCE = ? WHERE ID = ?", user.Balance, user.UserID)
	if err != nil {
		panic(err)
	}

}

func newUser(UserId int) *user {
	var Balance float32

	row := db.QueryRow("SELECT BALANCE FROM BILLING.BALANCE WHERE ID = ?", UserId)

	err := row.Scan(&Balance)
	if err != nil {
		panic(err)
	}

	user := user{
		UserID:  UserId,
		Balance: Balance,
	}

	return &user

}

func CreditingFunds(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		panic(err)
	}

	userID := int(data["user_ID"].(float32))
	amount := data["amount"].(float32)

	user := newUser(userID)
	user.Balance += amount
	updateUser(user)

	w.WriteHeader(http.StatusOK)

}

func DebitingFunds(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		panic(err)
	}

	userID := int(data["user_ID"].(float32))
	amount := data["amount"].(float32)

	user := newUser(userID)
	user.Balance -= amount
	updateUser(user)

	w.WriteHeader(http.StatusOK)

}

func main() {
	r := mux.NewRouter()
	conf := config.NewConfig()
	var err error
	db, err = sql.Open("mysql", conf.UserDB+":"+conf.PasswordDB+"@/"+conf.NameDB)
	if err != nil {
		panic(err)
	}

	r.HandleFunc("/balance/{id:[0-9]+}", getBalance)
	r.HandleFunc("/transfer", transfer)
	r.HandleFunc("/crediting", CreditingFunds)
	r.HandleFunc("/debiting", DebitingFunds)

	http.ListenAndServe("localhost:"+conf.Port, nil)
}
