package main

import (
	"encoding/json"
	"fmt"
	//"io"
	"log"
	"net/http"
)

type Todo struct{
	UserID int `json:"userId"`
	ID int `json:"id"`
	Title string `json:"title"`
	Completed bool `json:"completed"`
}
func main() {

	url := "https://jsonplaceholder.typicode.com/todos/1"
	response, err := http.Get(url)

	if err!= nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK{
	//	bodyBytes,err := io.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
	// 	data := string(bodyBytes)
	// 	fmt.Println((data))
	todoItem := Todo{}
	// json.Unmarshal(bodyBytes,&todoItem)
	decoder := json.NewDecoder(response.Body)
	decoder.DisallowUnknownFields()
	if err:= decoder.Decode(&todoItem); err!=nil{
		log.Fatal("Decode e rror", err)
	}
	//fmt.Printf("Data from API: %+v", todoItem)
	//todo, err := json.Marshal(todoItem)
	todo, err := json.MarshalIndent(todoItem, "", "\t") //pirint pretty json strucvute wih MarshalIndent
	if err!=nil {
		log.Fatal("Encoding error:", err)
	}
	fmt.Println(string(todo))

	}
	


}