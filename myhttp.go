package main

import (
    "crypto/md5"
    "fmt"
	"net/http"
	"log"
	"io/ioutil"
)

func main() {
    MakeRequest()
}

func MakeRequest() {
	resp, err := http.Get("https://httpbin.org/get")
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

//	data := []byte(string(body))
    fmt.Printf("%x", md5.Sum(body))
}

