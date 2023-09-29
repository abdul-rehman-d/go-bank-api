package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ApiServer struct {
	listenAddress string
	store         Storage
}

func NewApiServer(listenAddress string, store Storage) (*ApiServer, error) {
	return &ApiServer{
		listenAddress,
		store,
	}, nil
}

func (server *ApiServer) Run() error {
	router := mux.NewRouter()
	router.Handle("/account", makeHttpHandleFunc(server.handleAccount))
	router.Handle("/account/{id}", makeHttpHandleFunc(server.handleGetAccount))

	fmt.Println("Server is running on: ", server.listenAddress)

	return http.ListenAndServe(server.listenAddress, router)
}

func (server *ApiServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "POST" {
		return server.handleCreateAccount(w, r)
	}
	if r.Method == "DELETE" {
		return server.handleDeleteAccount(w, r)
	}

	return WriteJSON(w, http.StatusMethodNotAllowed, ApiError{
		Error: fmt.Sprintf("method not allowed %s", r.Method),
	})
}

func (server *ApiServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "GET" {
		return fmt.Errorf("method not allowed %s", r.Method)
	}

	id := mux.Vars(r)["id"]

	if len(id) == 0 {
		return fmt.Errorf("param \"id\" is required")
	}

	parsedId, err := strconv.Atoi(id)

	if err != nil {
		return fmt.Errorf("id must be int")
	}

	acc, err := server.store.GetAccount(parsedId)

	if err != nil {
		if err.Error() == "not found" {
			return err
		}

		fmt.Printf("failed to get account with id %d: %s", parsedId, err.Error())
		return fmt.Errorf("something went wrong")
	}

	return WriteJSON(w, http.StatusOK, acc)
}

func (server *ApiServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	accReq := CreateAccountRequest{}
	if err := json.NewDecoder(r.Body).Decode(&accReq); err != nil {
		return fmt.Errorf("inavalid request body")
	}
	if len(accReq.FirstName) == 0 {
		return fmt.Errorf("first name is required")
	}
	if len(accReq.FirstName) < 3 {
		return fmt.Errorf("first name cannot be less than 3 characters")
	}
	if len(accReq.LastName) == 0 {
		return fmt.Errorf("last name is required")
	}
	if len(accReq.LastName) < 3 {
		return fmt.Errorf("last name cannot be less than 3 characters")
	}
	acc := NewAccount(accReq.FirstName, accReq.LastName)
	acc, err := server.store.CreateAccount(acc)
	if err != nil {
		return fmt.Errorf("failed to create account in db: %s", err.Error())
	}

	return WriteJSON(w, http.StatusOK, acc)
}

func (server *ApiServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (server *ApiServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// helpers
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type apiFunc = func(w http.ResponseWriter, r *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

func makeHttpHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{
				Error: err.Error(),
			})
		}
	}
}
