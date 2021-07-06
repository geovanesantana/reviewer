package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/buger/jsonparser"
)

type Reviewer struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

type Reviewers struct {
	Reviewers []Reviewer `json:"list"`
}

type CurrentIndex struct {
	Current int64 `json:"current"`
}

func main() {

	list := getList()
	nextIndex := getNextIndex()

	if int(nextIndex) >= len(list.Reviewers) {
		nextIndex = 0
	}

	nextReviewer := list.Reviewers[nextIndex]

	fmt.Println("Starting request")
	body, err := json.Marshal(map[string]string{
		"text": fmt.Sprintf("<users/%s>", nextReviewer.Id),
	})

	if err != nil {
		log.Fatal(err)
	}

	url := "https://chat.googleapis.com/v1/spaces/AAAAWGALFXU/messages?key=AIzaSyDdI0hCZtE6vySjMm-WEfRq3CPzqKqqsHI&token=20AaV2TO_kqE455MJBeEGnzo8GosMvYj9eyI1IaMA04%3D"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	updateIndex(nextIndex)
	fmt.Println("Successful request")
}

func updateIndex(index int64) {
	data := CurrentIndex{Current: index}
	file, _ := json.MarshalIndent(data, "", "")

	err := ioutil.WriteFile("current.json", file, 0755)
	if err != nil {
		log.Fatal(err)
	}
}

func getNextIndex() int64 {
	file, err := os.Open("current.json")

	if err != nil {
		panic("Error opening JSON file")
	}

	defer file.Close()
	byteValue, _ := ioutil.ReadAll(file)

	current, err := jsonparser.GetInt(byteValue, "current")
	if err != nil {
		log.Fatal(err)
	}

	return current + 1
}

func getList() Reviewers {
	file, err := os.Open("./list.json")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	var list Reviewers
	byteValue, _ := ioutil.ReadAll(file)
	json.Unmarshal(byteValue, &list)

	return list
}
