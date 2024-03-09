package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func writeJsonreqwest(url string, reqvest string, jsonString string) {
	data := []byte(jsonString)
	req, err := http.NewRequest(reqvest, url, bytes.NewBuffer(data)) //Post("http://localhost:8080/create")
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "aplication/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	fmt.Println("Response status: ", resp.Status)
	fmt.Println("Response Header: ", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response body: ", string(body))
}

func main() {

	url := "http://localhost:8080/create"
	jsonStr := `{"name":"Rinat","age":"40","frends":[]}`
	writeJsonreqwest(url, "POST", jsonStr)

	jsonStr = `{"name":"Alex","age":"32","frends":[]}`
	writeJsonreqwest(url, "POST", jsonStr)

	jsonStr = `{"name":"Uri","age":"18","frends":[]}`
	writeJsonreqwest(url, "POST", jsonStr)

	url = "http://localhost:8080/make_friends"
	jsonStr = `{"source_id":1,"target_id":2}`
	writeJsonreqwest(url, "POST", jsonStr)

	url = "http://localhost:8080/user"
	jsonStr = `{"target_id":1}`
	writeJsonreqwest(url, "DELETE", jsonStr)

	url = "http://localhost:8080/create"
	jsonStr = `{"name":"Oleg","age":"29","frends":[]}`
	writeJsonreqwest(url, "POST", jsonStr)

	url = "http://localhost:8080/2"
	jsonStr = `{"new age":"120"}`
	writeJsonreqwest(url, "PUT", jsonStr)
}
