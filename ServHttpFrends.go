package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type Frends struct {
	Source_id int `json:"source_id"`
	Target_id int `json:"target_id"`
}
type UpdateAge struct {
	NAge string `json:"new age"`
}
type Delete struct {
	Target_id int `json:"target_id"`
}

type User struct { // json   {"name" : "Rinkiller" , "age" : "20", "frends" : []int}
	Name   string `json:"name"`
	Age    string `json:"age"`
	Frends []int  `json:"frends"`
}

type service struct {
	store map[int]*User // [id] struct User
}

// 1. Сделайте обработчик создания пользователя. У пользователя должны быть следующие поля: имя, возраст и массив друзей. Пользователя необходимо сохранять в мапу. Пример запроса:
// POST /create HTTP/1.1
// Content-Type: application/json; charset=utf-8
// Host: localhost:8080
// {"name":"some name","age":"24","friends":[]}
// Данный запрос должен возвращать ID пользователя и статус 201.
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
		var id int = 1
		for {
			if _, ok := s.store[id]; ok == false {
				break
			}
			id++
		}
		s.store[id] = &u
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf("Пользователь %s зарегистрирован. ID пользователя: %d \n", u.Name, id)))
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

// 2. Сделайте обработчик, который делает друзей из двух пользователей. Например, если мы создали двух пользователей и нам вернулись их ID, то в запросе мы можем указать ID пользователя, который инициировал запрос на дружбу, и ID пользователя, который примет инициатора в друзья. Пример запроса:
// POST /make_friends HTTP/1.1
// Content-Type: application/json; charset=utf-8
// Host: localhost:8080
// {"source_id":1,"target_id":2}
// Данный запрос должен возвращать статус 200 и сообщение «username_1 и username_2 теперь друзья».
func (s *service) makeFrends(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		defer r.Body.Close()
		var f Frends
		if err := json.Unmarshal(content, &f); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		if (len(s.store)-1 < f.Source_id) || (len(s.store)-1 < f.Target_id) || (f.Source_id < 0) || (f.Target_id < 0) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		s.store[f.Source_id].Frends = append(s.store[f.Source_id].Frends, f.Target_id)
		s.store[f.Target_id].Frends = append(s.store[f.Target_id].Frends, f.Source_id)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("%s и %s теперь друзья \n", s.store[f.Target_id].Name, s.store[f.Source_id].Name)))
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

// Сделайте обработчик, который возвращает всех  пользователей. Пример запроса:
// GET /get
func (s *service) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		response := ""
		for _, u := range s.store {
			var frendsN []string
			for _, idF := range u.Frends {
				frendsN = append(frendsN, s.store[idF].Name)
			}
			response += fmt.Sprintf("Имя пользователя: %s, возраст: %s, имена друзей: %s \n", u.Name, u.Age, frendsN)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

// 3. Сделайте обработчик, который удаляет пользователя. Данный обработчик принимает ID пользователя и удаляет его из хранилища, а также стирает его из массива friends у всех его друзей. Пример запроса:
// DELETE /user HTTP/1.1
// Content-Type: application/json; charset=utf-8
// Host: localhost:8080
// {"target_id":"1"}
// Данный запрос должен возвращать 200 и имя удалённого пользователя.
func (s *service) deleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == "DELETE" {
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		defer r.Body.Close()
		var d Delete
		if err := json.Unmarshal(content, &d); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		if (len(s.store)-1 < d.Target_id) || (d.Target_id < 0) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		deleteUserName := s.store[d.Target_id].Name
		//Проверить и удалить id пользователя из всех его друзей
		for _, fid := range s.store[d.Target_id].Frends {
			for id, v := range s.store[fid].Frends {
				if d.Target_id == v {
					if id == 0 {
						if len(s.store[fid].Frends) == 1 {
							s.store[fid].Frends = []int{}
							break
						}
						s.store[fid].Frends = s.store[fid].Frends[id+1:]
						break
					}
					if id == len(s.store[fid].Frends)-1 {
						s.store[fid].Frends = s.store[fid].Frends[0:id]
						break
					}
					listFerst := s.store[fid].Frends[0:id]
					listOver := s.store[fid].Frends[id+1:]
					s.store[fid].Frends = append(listFerst, listOver...)
					break
				}

			}
		}
		delete(s.store, d.Target_id)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Пользователь %s удален из базы пользователей \n", deleteUserName)))
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

// 4. Сделайте обработчик, который возвращает всех друзей пользователя. Пример запроса:
// GET /friends/user_id HTTP/1.1
// Host: localhost:8080
// Connection: close
// После /friends/ указывается id пользователя, друзей которого мы хотим увидеть.
// -----------------------------------------
func (s *service) GetUser(w http.ResponseWriter, r *http.Request) { //BUGS page not found
	if r.Method == "GET" {
		id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/friends/"))
		if err != nil || id < 0 || len(s.store)-1 < id {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var frendsN []string
		for _, idF := range s.store[id].Frends {
			frendsN = append(frendsN, s.store[idF].Name)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Пользователь: %s,  дружит с : %s \n", s.store[id].Name, frendsN)))
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

// 5. Сделайте обработчик, который обновляет возраст пользователя. Пример запроса:
// PUT /user_id HTTP/1.1
// Content-Type: application/json; charset=utf-8
// Host: localhost:8080
// {"new age":"28"}
// Запрос должен возвращать 200 и сообщение «возраст пользователя успешно обновлён».
func (s *service) updateAge(w http.ResponseWriter, r *http.Request) { //BUGS page not found
	if r.Method == "PUT" {
		id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/"))
		if err != nil || id < 0 || len(s.store)-1 < id {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		defer r.Body.Close()
		var a UpdateAge
		if err := json.Unmarshal(content, &a); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		s.store[id].Age = a.NAge
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Возраст пользователя успешно обновлён \n"))
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func main() {
	mux := http.NewServeMux()
	srv := service{make(map[int]*User)}
	mux.HandleFunc("/create", srv.Create)
	mux.HandleFunc("/get", srv.GetAll)
	mux.HandleFunc("/make_friends", srv.makeFrends)
	mux.HandleFunc("/user", srv.deleteUser)
	mux.HandleFunc("/friends/", srv.GetUser)
	mux.HandleFunc("/", srv.updateAge)
	http.ListenAndServe("localhost:8080", mux)
}
