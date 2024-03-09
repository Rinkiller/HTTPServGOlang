package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type User struct { // json  {"name" : "Rinkiller" , "age" : 20} {"name" : "Rinkiller" , "age" : 20, map[string]string}
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type service struct {
	store map[string]*User
}

func (u *User) toString() string {
	return fmt.Sprintf("Name is %s, Age is %d \n", u.Name, u.Age)
}

func (s *service) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		defer r.Body.Close()
		var u User
		if err := json.Unmarshal(content, &u); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		s.store[u.Name] = &u
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("user was created " + u.Name))
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func (s *service) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		response := ""
		for _, user := range s.store {
			response += user.toString()
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func main() {
	mux := http.NewServeMux()
	srv := service{make(map[string]*User)}
	mux.HandleFunc("/create", srv.Create)
	mux.HandleFunc("/get", srv.GetAll)
	http.ListenAndServe("localhost:8080", mux)
}
