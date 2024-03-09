package main

import (
	"io/ioutil"
	"net/http"
	"strings"
)

type service struct {
	store map[string]string
}

func (s *service) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		splitCont := strings.Split(string(content), " ")
		s.store[splitCont[0]] = string(content)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("user was created " + splitCont[0]))
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func (s *service) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		response := ""
		for _, user := range s.store {
			response += user + "\n"
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func main() {
	mux := http.NewServeMux()
	srv := service{make(map[string]string)}
	mux.HandleFunc("/create", srv.Create)
	mux.HandleFunc("/get", srv.GetAll)
	http.ListenAndServe("localhost:8080", mux)
}
