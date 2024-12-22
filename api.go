package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddr string
	store      Storage
}

func WriteJson(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string
}

func makeHttpHandleFunc(fn apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			WriteJson(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/accounts", makeHttpHandleFunc(s.handleAccount))
	router.HandleFunc("/account/{id}", makeHttpHandleFunc(s.handleAccount))

	log.Println("Listening on", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)

}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {

	if r.Method == "GET" {
		return s.handleGetAccounts(w, r)
	}

	if r.Method == http.MethodPost {
		return s.handlePostAccount(w, r)
	}

	if r.Method == http.MethodPut {
		return s.handlePutAccount(w, r)
	}

	if r.Method == http.MethodDelete {
		return s.handleDeleteAccount(w, r)
	}

	return fmt.Errorf("Method not allowed %s", r.Method)
}

func (s *APIServer) handleGetAccountById(w http.ResponseWriter, r *http.Request) error {

	vars := mux.Vars(r)["id"]

	//account := NewAccount("John", "Doe")
	//db.get(vars["id"])

	fmt.Println(vars)

	return WriteJson(w, http.StatusOK, &Account{})
}

// GET ACCOUNTS
func (s *APIServer) handleGetAccounts(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetAccounts()

	if err != nil {
		return err
	}

	return WriteJson(w, http.StatusOK, accounts)
}

func (s *APIServer) handlePostAccount(w http.ResponseWriter, r *http.Request) error {
	createAccountReq := CreateAccountRequest{}

	if err := json.NewDecoder(r.Body).Decode(&createAccountReq); err != nil {
		return err
	}

	account := NewAccount(createAccountReq.FirstName, createAccountReq.LastName)

	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	return WriteJson(w, http.StatusCreated, account)

}

func (s *APIServer) handlePutAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}
