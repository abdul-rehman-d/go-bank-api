package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type ApiServer struct {
	listenAddress string
	// s *Storage
}

func NewApiServer(listenAddress string) (*ApiServer, error) {
	return &ApiServer{
		listenAddress,
	}, nil
}

func (server *ApiServer) Run() error {
	router := mux.NewRouter()
	router.Handle("/account", makeHttpHandleFunc(server.handleAccount))

	fmt.Println("Server is running on: ", server.listenAddress)

	return http.ListenAndServe(server.listenAddress, router)
}

func (server *ApiServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return server.handleGetAccount(w, r)
	}
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
	acc := NewAccount("mad", "man")
	return WriteJSON(w, http.StatusOK, acc)
}

func (server *ApiServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
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
