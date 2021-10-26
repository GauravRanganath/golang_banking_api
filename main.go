// main.go
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Article - Our struct for all articles
type Account struct {
	Id            string `json:"Id"`
	AccountHolder string `json:"AccountHolder"`
	Balance       int    `json:"Balance"`
}

var Accounts []Account

func returnAllAccounts(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllArticles")
	json.NewEncoder(w).Encode(Accounts)
}

func returnSingleAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	for _, account := range Accounts {
		if account.Id == key {
			json.NewEncoder(w).Encode(account)
		}
	}
}

func createNewAccount(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: createNewAccount")

	reqBody, _ := ioutil.ReadAll(r.Body)
	var account Account
	json.Unmarshal(reqBody, &account)

	Accounts = append(Accounts, account)

	json.NewEncoder(w).Encode(account)
}

func deleteAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	for index, account := range Accounts {
		if account.Id == id {
			Accounts = append(Accounts[:index], Accounts[index+1:]...)
		}
	}

}

func handleTransaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id1 := vars["id1"]
	id2 := vars["id2"]
	transferAmount, err := strconv.Atoi(vars["transferAmount"])

	if err != nil {
		fmt.Println("Enter a valid transfer amount")
		fmt.Println(transferAmount)
	} else {
		for index, account := range Accounts {
			if account.Id == id1 {
				Accounts[index].Balance = account.Balance - transferAmount
			}
		}

		for index, account := range Accounts {
			if account.Id == id2 {
				Accounts[index].Balance = account.Balance + transferAmount
			}
		}
	}

}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/accounts", returnAllAccounts).Methods("GET")
	myRouter.HandleFunc("/account", createNewAccount).Methods("POST")
	myRouter.HandleFunc("/account/{id}", deleteAccount).Methods("DELETE")
	myRouter.HandleFunc("/account/{id}", returnSingleAccount).Methods("GET")
	myRouter.HandleFunc("/account/transaction/{id1}/{id2}/{transferAmount}", handleTransaction).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func main() {
	Accounts = []Account{}
	handleRequests()
}
